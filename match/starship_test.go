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
		// TODO: Add test capses.
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := StarshipText(tt.text); got != tt.want {
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
	for k, _ := range specificUserMatchers {
		if strings.ToLower(k) != k {
			t.Errorf("Account name %q should be lowercase in specificUserMatchers map", k)
		}
	}
}
