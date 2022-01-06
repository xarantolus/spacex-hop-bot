package consumer

import (
	"strings"
	"testing"
)

func TestImportantURLMap(t *testing.T) {
	for host := range importantURLs {
		t.Run(host, func(t *testing.T) {
			if strings.ToLower(host) != host {
				t.Errorf("host %q must be lowercase in importantURLs", host)
			}
			if strings.HasPrefix(host, "www.") {
				t.Errorf("host %q must not start with 'www.' in importantURLs", host)
			}
		})
	}
}

func Test_isImportantURL(t *testing.T) {
	tests := []struct {
		uri           string
		wantImportant bool
	}{
		{"http://nasaspaceflight.com/starbaselive", true},
		{"https://nasaspaceflight.com/starbaselive", true},
		{"https://www.nasaspaceflight.com/starbaselive", true},

		{"http://cnunezimages.com", true},
		{"http://cnunezimages.com/", true},
		{"http://cnunezimages.com/any-link-really", true},

		{"https://www.cameroncountytx.gov/spacex/", true},
		{"https://cameroncountytx.gov/spacex/", true},

		{"https://cameroncountytx.gov/", false},

		{"Twitter dot com", false},
		{"https://twitter.comcom", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if gotImportant := isImportantURL(tt.uri); gotImportant != tt.wantImportant {
				t.Errorf("isImportantURL(%q) = %v, want %v", tt.uri, gotImportant, tt.wantImportant)
			}
		})
	}
}
