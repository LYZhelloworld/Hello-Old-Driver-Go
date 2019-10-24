package scanner

import "testing"

func TestScannerCreation(t *testing.T) {
	s := Scanner{"liuli.be", "http"}
	if s.Domain != "liuli.be" {
		t.Errorf("Assertion failed: Domain != liuli.be")
	}
	if s.Protocol != "http" {
		t.Errorf("Assertion failed: Protocol != http")
	}
}

func TestFeed(t *testing.T) {
	s := Scanner{"liuli.be", "http"}
	page := s.GetFeed(1)
	t.Logf("url: %s\nSucceeded: %t\nlen(Content): %d",
		page.URL, page.Succeeded, len(page.Content))
	if page.URL != "http://liuli.be/wp/?feed=rss" {
		t.Errorf("Url is incorrect.")
	}
	page = s.GetFeed(2)
	t.Logf("url: %s\nSucceeded: %t\nlen(Content): %d",
		page.URL, page.Succeeded, len(page.Content))
	if page.URL != "http://liuli.be/wp/?feed=rss&paged=2" {
		t.Errorf("Url is incorrect.")
	}
}

func TestPages(t *testing.T) {
	s := Scanner{"liuli.be", "http"}
	channel := make(chan Page)
	s.GetPages([]string{
		"https://www.liuli.be/wp/72842.html",
		"https://www.liuli.be/wp/72837.html"}, channel)
	for i := 0; i < 2; i++ {
		page := <-channel
		t.Logf("url: %s\nSucceeded: %t\nlen(Content): %d",
			page.URL, page.Succeeded, len(page.Content))
	}
}
