package analyzer

import (
	"regexp"
	"strings"
)

const (
	magnetPrefix string = "magnet:?xt=urn:btih:"
	regexTitle string = `(?s)<title>(.+?)</title>`
	regexContent string = `(?s)<div\s+class="entry-content"\s*?>` +
		`(.*?)</div><!--\s*?.entry-content\s*?-->`
	regexMagnet40 string = `(?s)[^0-9a-fA-F]([0-9a-fA-F]{40})[^0-9a-fA-F]`
	regexMagnet32 string = `(?s)[^0-9a-fA-F]([0-9a-fA-F]{32})[^0-9a-fA-F]`
	
	regexFeedItem string = `(?s)<item>.*?` +
		`<title>(.*?)</title>\s*?<link>(.*?)</link>` +
		`.*?</item>`
)

type FeedItem struct {
	Title string
	Link string
}

func GetMagnetLinks(text string) (result []string) {
	var r *regexp.Regexp
	var match []string
	var matchAll [][]string
	
	r = regexp.MustCompile(regexContent)
	match = r.FindStringSubmatch(text)
	if match != nil {
		text = match[1]
	} else {
		return make([]string, 0)
	}
	
	result = make([]string, 0)
	
	r = regexp.MustCompile(regexMagnet40)
	matchAll = r.FindAllStringSubmatch(text, -1)
	if matchAll != nil {
		for _, v := range matchAll {
			result = append(result, v[1])
		}
	}
	
	r = regexp.MustCompile(regexMagnet32)
	matchAll = r.FindAllStringSubmatch(text, -1)
	if matchAll != nil {
		for _, v := range matchAll {
			result = append(result, v[1])
		}
	}
	
	for i, v := range result {
		result[i] = strings.ToLower(magnetPrefix + v)
	}
	
	return
}

func GetPageTitle(text string) string {
	r := regexp.MustCompile(regexTitle)
	match := r.FindStringSubmatch(text)
	if match != nil {
		return strings.TrimSpace(match[1])
	} else {
		return ""
	}
}

func GetFeedItems(text string) (result []FeedItem) {
	r := regexp.MustCompile(regexFeedItem)
	matchAll := r.FindAllStringSubmatch(text, -1)
	
	result = make([]FeedItem, 0)
	if matchAll != nil {
		for _, v := range matchAll {
			if len(v) == 3 {
				result = append(result, FeedItem{v[1], v[2]})
			}
		}
	}
	
	return
}