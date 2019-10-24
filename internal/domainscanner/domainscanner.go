package domainscanner

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const url = "http://acg.gy/"
const re = "<a href=\"https?://([a-zA-Z0-9\\.]*)\">"

// GetDomain gets the current domain name used
func GetDomain() string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	req.Header.Add("Accept-Charset", "utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	content := string(body)

	r := regexp.MustCompile(re)
	match := r.FindStringSubmatch(content)
	if match != nil {
		return strings.TrimSpace(match[1])
	} else {
		return ""
	}
}
