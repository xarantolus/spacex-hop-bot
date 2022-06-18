package match

import (
	"strings"
)

// StarshipText returns whether the given text mentions starship
func (m *StarshipMatcher) StarshipText(text string, antiKeywords []string, skipMatchers bool) bool {
	text = strings.ToLower(text)

	// If we find ignored words, we ignore the tweet
	if _, contains := containsAntikeyword(antiKeywords, text); contains {
		return false
	}

	// else we check if there are any interesting keywords
	if _, contains := startsWithAny(text, starshipKeywords...); contains {
		return true
	}

	// Then we check for more "dynamic" words like "S20", "B4", etc.
	// If we input text with URLs, we skip matchers. This is because URLs often
	// contain random sequences of characters that can be picked up by these matchers
	if !skipMatchers {
		for _, r := range starshipMatchers {
			if r.MatchString(text) {
				return true
			}
		}
	}

	// Now we check for keywords that need additional keywords to be matched,
	// e.g. "raptor", "deimos" etc.
	for _, mapping := range moreSpecificKeywords {
		if mapping.matches(text) {
			return true
		}
	}

	return false
}

func ContainsStarshipAntiKeyword(text string) bool {
	_, contains := containsAntikeyword(antiStarshipKeywords, strings.ToLower(text))
	return contains
}

func containsAntikeyword(antiKeywords []string, text string) (word string, contains bool) {
	for _, antiRegex := range antiKeywordRegexes {
		if antiRegex.MatchString(text) {
			return "(antiKeywordRegex)" + antiRegex.String(), true
		}
	}

	for _, mapping := range moreSpecificAntiKeywords {
		if mapping.matches(text) {
			return "(moreSpecificAntiKeywords)", true
		}
	}

	return startsWithAny(text, antiKeywords...)
}
