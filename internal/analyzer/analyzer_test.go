package analyzer

import "testing"

func TestMagnetLinks(t *testing.T) {
	result := GetMagnetLinks(testText)
	if len(result) < 2 {
		t.Errorf("Not all magnet links are matched.")
		t.Errorf("%v", result)
		return
	} else if len(result) > 2 {
		t.Errorf("Non-magnet links are matched.")
		t.Errorf("%v", result)
		return
	} else {
		if result[0] !=
			"magnet:?xt=urn:btih:012345678901234567890123456789abcdefabcd" {
			t.Errorf("Wrong magnet link is matched. Should be %s, got %s.",
				"magnet:?xt=urn:btih:012345678901234567890123456789abcdefabcd",
				result[0])
		}
		if result[1] !=
			"magnet:?xt=urn:btih:012345678901234567890123456789ab" {
			t.Errorf("Wrong magnet link is matched. Should be %s, got %s.",
				"magnet:?xt=urn:btih:012345678901234567890123456789ab",
				result[1])
		}
	}
}

func TestTitle(t *testing.T) {
	result := GetPageTitle(testText)
	if result == "" {
		t.Errorf("Cannot find the title.")
		return
	} else if result != "Testing" {
		t.Errorf("Wrong title is matched. Should be %s, got %s.",
			"Testing", result)
		return
	}
}

func TestUnicode(t *testing.T) {
	result := GetPageTitle(`<title>Testing | 琉璃神社 ★ HACG.me</title>`)
	if result == "" {
		t.Errorf("Cannot find the title.")
		return
	} else if result != "Testing | 琉璃神社 ★ HACG.me" {
		t.Errorf("Wrong title is matched. Should be %s, got %s.",
			"Testing | 琉璃神社 ★ HACG.me", result)
		return
	}
}

func TestFeedItems(t *testing.T) {
	result := GetFeedItems(testFeed)
	if len(result) != 2 {
		t.Errorf("Count of item is not correct.")
	}
	for _, item := range result {
		t.Logf("Link: %s", item)
		if item == "https://www.liuli.uk/wp" ||
			item != "https://www.liuli.uk/wp/12345.html" &&
				item != "https://www.liuli.uk/wp/67890.html" {
			t.Errorf("Wrong item is matched.")
		}
	}
}

const testText string = `<html>
<head>
<title>Testing</title>
</head>
<body>
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor 
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis 
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. 
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore 
eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt 
in culpa qui officia deserunt mollit anim id est laborum.
<div class="entry-content">
<p>Lorem ipsum dolor sit amet</p>
<p>Lorem ipsum dolor sit amet</p>
012345678901234567890123456789ABCDEFabcd
012345678901234567890123456789aB
<p>Lorem ipsum dolor sit amet</p>
<p>Lorem ipsum dolor sit amet</p>
</div><!-- .entry-content -->
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor 
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis 
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. 
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore 
eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt 
in culpa qui officia deserunt mollit anim id est laborum.
</body>
</html>`

const testFeed = `<?xml version="1.0" encoding="UTF-8"?><rss>
<channel>
<title>This title should not be matched</title>
<link>https://www.liuli.uk/wp</link>
<item>
	<title>Test 1</title>
	<link>https://www.liuli.uk/wp/12345.html</link>
</item>
<item>
	<title>Test 2</title>
	<link>https://www.liuli.uk/wp/67890.html</link>
</item>
</channel>
</rss>`
