package scanner

import "testing"

func TestScannerCreation(t *testing.T) {
	s := Scanner("liuli.uk", "http")
	if s.Domain != "liuli.uk" {
		t.Errorf("Assertion failed: Domain != liuli.uk")
	}
	if s.Protocol != "http" {
		t.Errorf("Assertion failed: Protocol != http")
	}
}

func TestFeed(t *testing.T) {
	s := Scanner("liuli.uk", "http")
	page := s.GetFeed(1)
	t.Logf("url: %s\nSucceeded: %t\nlen(Content): %d",
		page.Url, page.Succeeded, len(page.Content))
	if page.Url != "http://liuli.uk/wp/?feed=rss" {
		t.Errorf("Url is incorrect.")
	}
	page = s.GetFeed(2)
	t.Logf("url: %s\nSucceeded: %t\nlen(Content): %d",
		page.Url, page.Succeeded, len(page.Content))
	if page.Url != "http://liuli.uk/wp/?feed=rss&paged=2" {
		t.Errorf("Url is incorrect.")
	}
}

func TestPages(t *testing.T) {
	s := Scanner("liuli.uk", "http")
	channel := make(chan Page)
	s.GetPages([]string {
		"https://www.liuli.uk/wp/72842.html",
		"https://www.liuli.uk/wp/72837.html"}, channel)
	for i := 0; i < 2; i++ {
		page := <-channel
		t.Logf("url: %s\nSucceeded: %t\nlen(Content): %d",
			page.Url, page.Succeeded, len(page.Content))
	}
}