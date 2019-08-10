package analyzer

import (
	"regexp"
	"strings"
)

const (
	magnetPrefix string = "magnet:?xt=urn:btih:"
	regexTitle string = `(?s)<title>(.+?)</title>`
	regexContent string = `(?s)<div(?:\s+?)class="entry-content"(?:\s*?)>` +
		`(.*?)</div><!--(?:\s*?).entry-content(?:\s*?)-->`
	regexMagnet40 string = `(?s)[^0-9a-fA-F]([0-9a-fA-F]{40})[^0-9a-fA-F]`
	regexMagnet32 string = `(?s)[^0-9a-fA-F]([0-9a-fA-F]{32})[^0-9a-fA-F]`
)

type Analyzer struct {
	ContentOnly bool
}

func (a Analyzer) GetMagnetLinks(text string) (result []string) {
	var r *regexp.Regexp
	var match [][]string
	
	if a.ContentOnly {
		r = regexp.MustCompile(regexContent)
		match := r.FindStringSubmatch(text)
		if match != nil {
			text = match[1]
		} else {
			return make([]string, 0)
		}
	}
	
	result = make([]string, 0)
	
	r = regexp.MustCompile(regexMagnet40)
	match = r.FindAllStringSubmatch(text, -1)
	if match != nil {
		for _, v := range match {
			result = append(result, v[1])
		}
	}
	
	r = regexp.MustCompile(regexMagnet32)
	match = r.FindAllStringSubmatch(text, -1)
	if match != nil {
		for _, v := range match {
			result = append(result, v[1])
		}
	}
	
	for i, v := range result {
		result[i] = strings.ToLower(magnetPrefix + v)
	}
	
	return
}

func (a Analyzer) GetPageTitle(text string) string {
	r := regexp.MustCompile(regexTitle)
	match := r.FindStringSubmatch(text)
	if match != nil {
		return strings.TrimSpace(match[1])
	} else {
		return ""
	}
}