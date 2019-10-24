package main

import (
	"flag"
	"fmt"
	"llss/internal/analyzer"
	"llss/internal/domainscanner"
	"llss/internal/scanner"
	"os"
	"path/filepath"
	"strconv"
)

const channelSize int = 16

func main() {
	protocol, domain, p, ok := checkArgs()
	if !ok {
		return
	}

	s := scanner.Scanner{Domain: domain, Protocol: protocol}
	page := s.GetFeed(p)
	if page.Succeeded {
		links := analyzer.GetFeedItems(page.Content)
		resultChan := make(chan scanner.Page, channelSize)
		s.GetPages(links, resultChan)
		for i := 0; i < len(links); i++ {
			page = <-resultChan
			if !page.Succeeded {
				continue
			}
			title := analyzer.GetPageTitle(page.Content)
			magnets := analyzer.GetMagnetLinks(page.Content)
			if len(magnets) == 0 {
				continue
			}
			print(title, magnets)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Cannot get RSS Feed page. Exit.")
		return
	}
}

func print(title string, magnets []string) {
	fmt.Println(title)
	for _, magnet := range magnets {
		fmt.Printf("  %s\n", magnet)
	}
}

func checkArgs() (protocol string, domain string, page int, ok bool) {
	var (
		help bool
		err  error
	)

	flag.StringVar(&protocol, "p", "https",
		"Protocol used to get the feed, like \"http\" or \"https\"")
	flag.StringVar(&domain, "d", "",
		"Domain used to get the feed")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.Usage = usage

	flag.Parse()
	if help {
		flag.Usage()
		ok = false
		return
	}

	if domain == "" {
		domain = domainscanner.GetDomain()
		if domain == "" {
			fmt.Fprintf(os.Stderr, "Error when retrieving domain. Stop.\n")
			ok = false
			return
		}
	}

	args := flag.Args()
	if len(args) == 0 {
		page = 1
	} else {
		page, err = strconv.Atoi(args[0])
		if err != nil ||
			(err == nil && page <= 0) {
			fmt.Fprintf(os.Stderr, "Invalid page number\n")
			flag.Usage()
			ok = false
			return
		}
	}

	ok = true
	return
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-h] [-p protocol] [-d domain] [page]\n",
		filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  page: the page number of result. "+
		"Must be greater than 0. Default value is 1.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}
