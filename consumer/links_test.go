package consumer

import "testing"

func Test_isImportantURL(t *testing.T) {
	tests := []struct {
		uri           string
		wantImportant bool
	}{
		{"http://nasaspaceflight.com/starbaselive", true},
		{"https://nasaspaceflight.com/starbaselive", true},
		{"https://www.nasaspaceflight.com/starbaselive", true},

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
