package scanner

import (
	"io/ioutil"
	"fmt"
	"net/http"
)

type Scanner struct {
	Domain string
	Protocol string
}

type Page struct {
	Url string
	Succeeded bool
	Content string
}

func (s Scanner) GetFeed(page int) Page {
	var url string
	if page == 1 {
		url = fmt.Sprintf("%s://%s/wp/?feed=rss", s.Protocol, s.Domain)
	} else if page > 1 {
		url = fmt.Sprintf("%s://%s/wp/?feed=rss&paged=%d", s.Protocol, s.Domain, page)
	} else {
		// Error
		return Page{url, false, ""}
	}
	
	resultChan := make(chan Page)
	go visit(url, resultChan)
	result := <-resultChan
	return result
}

func (s Scanner) GetPages(urls []string, channel chan Page) {
	for _, url := range urls {
		go visit(url, channel)
	}
	return
}

func visit(url string, channel chan Page) {
	result, ok := get(url)
	channel <- Page{url, ok, result}
}

func get(url string) (string, bool) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", false
	}
	req.Header.Add("Accept-Charset", "utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}

	return string(body), true
}