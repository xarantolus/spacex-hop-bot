package match

import (
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/util"
)

// Note that all text here must be lowercase because the text is lowercased in the matching function
var (
	starshipKeywords = []string{"starship", "superheavy", "raptor", "super heavy"}

	antiStarshipKeywords = []string{
		"electron", "blue origin", "neutron", "rocket lab", "rocketlab", "falcon", "starlink",
		"tesla", "openai", "boring", "hyperloop", "solarcity", "neuralink",

		// Not interested in other stuff
		"doge", "fsd",

		"no tfr",

		// 3d models are nice, but we only care about "real info"
		"thanks", "thank you", "cheers", "render", "animation", "3d",

		"not starship", "non starship", "not about starship",

		// kerbal space program != "official" news
		"kerbal space program", "ksp",

		// not *that* kind of raptor
		"velociraptor", "jurassic", "cretaceous", "dino",

		"sale", "buy", "shop", "giveaway", "give away", "retweet", "birthday",
	}

	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`sn\d+`),
		// Booster BNx
		regexp.MustCompile(`bn\d+`),
	}

	closureTFRRegex = regexp.MustCompile("(?:closure|tfr|cryo|fts)")
	// Users known to post better information that requires specific filtering
	specificUserMatchers = map[string]*regexp.Regexp{
		"bocachicagal":    regexp.MustCompile("(?:alert|static fire|closure|cryo|evacua)"),
		"rgvaerialphotos": closureTFRRegex,
		"bocaroad":        closureTFRRegex,
		"infographictony": closureTFRRegex,
		"spacex360":       closureTFRRegex,
		"bluemoondance74": closureTFRRegex,
		"nextspaceflight": closureTFRRegex,
		"tylerg1998":      closureTFRRegex,
		"spacetfrs":       regexp.MustCompile("(?:brownsville)"),

		// For Elon, we try to match anything that could be insider info
		"elonmusk": regexp.MustCompile("(?:booster|orbit|cryo|static fire|tower|ship|rud|engine)"),
	}
)

// StarshipText returns whether the given text mentions starship
func StarshipText(text string, ignoreBlocklist bool) bool {

	text = strings.ToLower(text)

	if !ignoreBlocklist {
		for _, k := range antiStarshipKeywords {
			if strings.Contains(text, k) {
				return false
			}
		}
	}

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
	// Ignore OLD tweets
	if d, err := tweet.CreatedAtTime(); err == nil && time.Since(d) > 24*time.Hour {
		return false
	}

	// We do not care about tweets that are timestamped with a text more than 24 hours ago
	// e.g. if someone posts a photo and then writes "took this on March 15, 2002"
	if d, ok := util.ExtractDate(tweet.FullText); ok && time.Since(d) > 24*time.Hour {
		return false
	}

	// Now check if the text of the tweet matches what we're looking for.
	// if it's elon musk, then we don't check for anti-keywords
	if StarshipText(tweet.FullText, tweet.User != nil && tweet.User.ScreenName == "elonmusk") {
		return true
	}

	// Now check if we have a matcher for this specific user.
	// These users usually post high-quality information
	if tweet.User != nil {
		m, ok := specificUserMatchers[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			return m.MatchString(strings.ToLower(tweet.FullText))
		}
	}

	return false
}
