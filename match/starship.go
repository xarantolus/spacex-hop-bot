package match

import (
	"regexp"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

// Note that all text here must be lowercase because the text is lowercased in the matching function
var (
	starshipKeywords = []string{"starship", "superheavy", "raptor", "super heavy"}

	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`sn\d+`),
		// Booster BNx
		regexp.MustCompile(`bn\d+`), // TODO: "Booster 1"
	}

	closureTFRRegex = regexp.MustCompile("(?:closure|TFR|cryo|FTS)")
	// Users known to post better information that requires specific filtering
	specificUserMatchers = map[string]*regexp.Regexp{
		"bocachicagal":    regexp.MustCompile("(?:alert|static fire)|(?:closure|cryo|evacua)"),
		"rgvaerialphotos": closureTFRRegex,
		"bocaroad":        closureTFRRegex,
		"infographictony": closureTFRRegex,
		"spacex360":       closureTFRRegex,
		"bluemoondance74": closureTFRRegex,
		"nextspaceflight": closureTFRRegex,
	}
)

// StarshipText returns whether the given text mentions starship
func StarshipText(text string) bool {
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

// StarshipTweet returns whether the given tweet mentions starship. It also includes custom matchers for certain users
func StarshipTweet(tweet *twitter.Tweet) bool {
	if StarshipText(tweet.FullText) {
		return true
	}

	if tweet.User != nil {
		m, ok := specificUserMatchers[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			return m.MatchString(strings.ToLower(tweet.FullText))
		}
	}

	return false
}
