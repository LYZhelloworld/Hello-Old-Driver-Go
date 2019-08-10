package scanner

import (
	"io/ioutil"
	"fmt"
	"net/http"
)

const channelSize int = 128

type scanner struct {
	config Config
	viewedUrls map[string]bool
	start int
	end int
	currentId int
	isFinished bool
	
	Next chan Page
}

type Page struct {
	Url string
	Succeeded bool
	Content string
}

type Config struct {
	Domain string
	Protocol string
}

func Scanner(config Config, start int, end int) *scanner {
	s := new(scanner)
	s.config = config
	s.viewedUrls = make(map[string]bool)
	s.start = start
	s.end = end
	s.Next = make(chan Page, channelSize)
	s.currentId = start
	s.isFinished = false
	return s
}

func (s scanner) LoadVisitedUrls(resourceList []string) {
	for _, v := range resourceList {
		s.viewedUrls[v] = true
	}
	return
}

func (s scanner) Start() {
	for i := s.start; i <= s.end; i++ {
		poolId := i - s.start
		go s.visit(i, poolId)
	}
}

func (s scanner) visit(id int, poolId int) {
	if id > s.end {
		s.Next <- Page{"", false, ""}
		return
	}
	
	url := generateUrl(s.config.Protocol, s.config.Domain, id)
	if isVisited(url, s.viewedUrls) {
		// Skip visited urls
		s.Next <- Page{url, false, ""}
		return
	}
	
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		s.Next <- Page{url, false, ""}
		return
	}
	req.Header.Add("Accept-Charset", "utf-8")
	resp, err := client.Do(req)
	if err != nil {
		s.Next <- Page{url, false, ""}
		return
	}
	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Next <- Page{url, false, ""}
		return
	}
	
	s.Next <- Page{url, true, string(body)}
	return
}

func generateUrl(protocol string, domain string, id int) string {
	return fmt.Sprintf("%s://%s/wp/%d.html", protocol, domain, id)
}

func isVisited(key string, list map[string]bool) bool {
	v, ok := list[key]
	return ok && (v == true)
}