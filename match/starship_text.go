package match

import (
	"strings"
)

// StarshipText returns whether the given text mentions starship
func (m *StarshipMatcher) StarshipText(text string, antiKeywords []string) bool {
	text = strings.ToLower(text)

	// If we find ignored words, we ignore the tweet
	if containsAntikeyword(antiKeywords, text) {
		return false
	}

	// else we check if there are any interesting keywords
	if containsAny(text, starshipKeywords...) {
		return true
	}

	// Then we check for more "dynamic" words like "S20", "B4", etc.
	for _, r := range starshipMatchers {
		if r.MatchString(text) {
			return true
		}
	}

	// Now we check for keywords that need additional keywords to be matched,
	// e.g. "raptor", "deimos" etc.
	for _, mapping := range moreSpecificKeywords {
		if startsWithAny(text, mapping.from...) && startsWithAny(text, mapping.to...) && !startsWithAny(text, mapping.antiKeywords...) {
			return true
		}
	}

	return false
}

func ContainsStarshipAntiKeyword(text string) bool {
	return containsAntikeyword(antiStarshipKeywords, strings.ToLower(text))
}

func containsAntikeyword(words []string, text string) bool {
	return startsWithAny(text, words...)
}

func compose(s ...[]string) (res []string) {
	var dedup = map[string]bool{}
	for _, v := range s {
		for _, k := range v {
			if !dedup[k] {
				res = append(res, k)
				dedup[k] = true
			}
		}
	}
	return res
}
