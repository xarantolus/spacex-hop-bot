package match

import (
	"strings"
	"testing"
)

func TestStarship(t *testing.T) {
	tests := []struct {
		text string
		want bool
	}{
		{"No TFR posted for today", false},
		{"SN10", true},
		{"BN10", true},
		{"Starship SN10", true},
		{"Starship SN10", true},
		{"Unrelated doge coin tweet that also contains the keyword Starship", false},
		{"Unrelated tesla tweet", false},
		{"this tweet is not starship related", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := StarshipText(tt.text, false); got != tt.want {
				t.Errorf("StarshipText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariables(t *testing.T) {
	for _, k := range starshipKeywords {
		if strings.ToLower(k) != k {
			t.Errorf("Keyword %q should be lowercase in starshipKeywords slice", k)
		}
	}
	for k := range specificUserMatchers {
		if strings.ToLower(k) != k {
			t.Errorf("Account name %q should be lowercase in specificUserMatchers map", k)
		}
	}
}
