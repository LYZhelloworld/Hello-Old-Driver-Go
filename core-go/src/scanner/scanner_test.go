package scanner

import "testing"

func TestScannerCreation(t *testing.T) {
	s := Scanner(Config{"liuli.uk", "http"}, 10000, 10010)
	if s.config.Domain != "liuli.uk" {
		t.Errorf("Assertion failed: Domain != liuli.uk")
	}
	if s.config.Protocol != "http" {
		t.Errorf("Assertion failed: Protocol != http")
	}
	if s.start != 10000 {
		t.Errorf("Assertion failed: start != 10000")
	}
	if s.end != 10010 {
		t.Errorf("Assertion failed: end != 10010")
	}
}

func TestOnePage(t *testing.T) {
	s := Scanner(Config{"liuli.uk", "http"}, 72916, 72916)
	s.Start()
	page := <-s.Next
	
	t.Logf("url: %s", page.Url)
	t.Logf("Succeeded: %t", page.Succeeded)
	t.Logf("len(Content): %d", len(page.Content))
}

func TestMultiplePages(t *testing.T) {
	s := Scanner(Config{"liuli.uk", "http"}, 72910, 72919)
	s.Start()
	for i := 72910; i <= 72919; i++ {
		page := <-s.Next
		
		t.Logf("url: %s\nSucceeded: %t\nlen(Content): %d",
				page.Url, page.Succeeded, len(page.Content))
	}
}