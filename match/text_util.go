package match

import (
	"log"
	"strings"
)

// containsAny checks whether any of words is *anywhere* in the text
func containsAny(text string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(text, w) {
			return true
		}
	}
	return false
}

// startsWithAny checks whether any of words is the start of a sequence of words in the text
func startsWithAny(text string, words ...string) bool {
	if len(words) == 0 {
		return false
	}

	var iterations = 0

	var currentIndex = 0

	for {
		iterations++

		var nextIndexOffset = strings.IndexFunc(text[currentIndex:], isAlphanumerical)
		if nextIndexOffset < 0 {
			break
		}
		currentIndex += nextIndexOffset

		for _, w := range words {
			if strings.HasPrefix(text[currentIndex:], w) {
				return true
			}
		}

		// Now skip to the next non-alphanumerical character
		nextIndexOffset = strings.IndexFunc(text[currentIndex:], func(r rune) bool {
			return !isAlphanumerical(r)
		})
		if nextIndexOffset < 0 {
			break
		}
		currentIndex += nextIndexOffset

		if iterations > 1000 {
			log.Printf("Input text %q causes startsWithAny to loop longer than expected", text)
			return false
		}
	}

	return false
}

func isAlphanumerical(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

// compose calculates the union of the given sets of strings, eliminating duplicates
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

// ignoreSpaces returns an array with exactly the words in words,
// but it also generates additional words by removing all spaces
func ignoreSpaces(words []string) (result []string) {
	var dedup = map[string]bool{}

	for _, w := range words {
		if !dedup[w] {
			dedup[w] = true
			result = append(result, w)
		}

		split := strings.Fields(w)
		if len(split) == 1 {
			if !dedup[w] {
				dedup[w] = true
				result = append(result, w)
			}
			continue
		}

		nw := strings.Join(split, "")
		if !dedup[nw] {
			dedup[nw] = true
			result = append(result, nw)
		}

		nw = strings.Join(split, "-")
		if !dedup[nw] {
			dedup[nw] = true
			result = append(result, nw)
		}

		nw = strings.Join(split, "_")
		if !dedup[nw] {
			dedup[nw] = true
			result = append(result, nw)
		}

		if len(split) >= 3 {
			nw = split[0] + strings.Join(split[1:], " ")
			if !dedup[nw] {
				dedup[nw] = true
				result = append(result, nw)
			}
			nw = strings.Join(split[:2], " ") + split[2]
			if !dedup[nw] {
				dedup[nw] = true
				result = append(result, nw)
			}
		}
	}

	return
}
