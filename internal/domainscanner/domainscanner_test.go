package domainscanner

import "testing"

func TestGetDomain(t *testing.T) {
	result := GetDomain()
	t.Logf(result)
	if result == "" {
		t.Errorf("Result is empty.")
	}
}
