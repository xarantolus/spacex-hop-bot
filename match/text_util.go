package match

import (
	"log"
	"strings"
	"unicode"
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
	var iterations = 0

	var currentIndex = 0

	for {
		iterations++

		for currentIndex < len(text) && !isAlphanumerical(rune(text[currentIndex])) {
			currentIndex++
		}

		for _, w := range words {
			if strings.HasPrefix(text[currentIndex:], w) {
				return true
			}
		}

		// Now skip to the next space character
		for currentIndex < len(text) && !unicode.IsSpace(rune(text[currentIndex])) {
			currentIndex++
		}

		if currentIndex == len(text) {
			break
		}

		if iterations > 1000 {
			log.Printf("Input text %q causes containsAny to loop longer than expected", text)
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
