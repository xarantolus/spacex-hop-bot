package match

import (
	"regexp"
	"strings"
)

// Note that all text here must be lowercase because the text is lowercased in the matching function
var (
	starshipKeywords = []string{"starship", "superheavy", "raptor"}

	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`sn\d+`),
		// Booster BNx
		regexp.MustCompile(`bn\d+`),
	}
)

// Starship returns whether the given text mentions starship
func Starship(text string) bool {
	text = strings.ToLower(text)

	for _, k := range starshipKeywords {
		if strings.Contains(text, k) {
			return true
		}
	}

	for _, r := range starshipMatchers {
		if r.MatchString(text) {
			return true
		}
	}

	return false
}
