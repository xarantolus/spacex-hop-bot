package match

import (
	"strings"
	"testing"
)

func Test_startsWithAnyGeneric(t *testing.T) {
	var searchedPrefixes = []string{"test", "best", "rest", "more than one word"}

	tests := []struct {
		argText string
		want    bool
	}{
		{"testing is nice", true},
		{"wrongprefixtesting is nice", false},
		{"the test keyword can be at any point in the string", true},
		{"we want to support more than one word", true},
		{"we want to support less than one word", false},

		{"the #test hashtag should still be recognized", true},
		{"also @test should work", true},
		{"#test at the beginning", true},
		{"@test should work", true},
		{"this \"test\" seems ok", true},
		{"this \"tes\"t seems ok", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := startsWithAny(tt.argText, searchedPrefixes...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_startsWithAnyStarshipAntikeywords(t *testing.T) {
	tests := []struct {
		argText string
		want    bool
	}{
		{"road closure with no information where it is", false},
		{"", false},
		{"Starship-Orion is a good idea", true},
		{"KSP is my favourite game!", true},
		{"What if I add many spaces before       KSP", true},
		{"-   .  +  - . k . -  .  .. -   . . -.- .", false},
		{"-   .  +  - . ksp . -  .  .. -   . . -.- .", true},
		{"Project DogeCoin onto a Starship!", true},
		{"Starship reentering Kerbin's atmosphere", true},
		{"GSE Tank 6 rolling out", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := startsWithAny(strings.ToLower(tt.argText), antiStarshipKeywords...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_containsAnyGeneric(t *testing.T) {
	var searchedInfixes = []string{"test", "best", "rest", "more than one word"}

	tests := []struct {
		argText string
		want    bool
	}{
		{"testing is nice", true},
		{"wrongprefixtesting is nice", true},
		{"the test keyword can be at any point in the string", true},
		{"we want to support more than one word", true},
		{"we want to support less than one word", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := containsAny(tt.argText, searchedInfixes...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}
