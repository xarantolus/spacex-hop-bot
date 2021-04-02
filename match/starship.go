package match

import (
	"log"
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
		"tesla", "openai", "boring", "hyperloop", "solarcity", "neuralink", "sls", "ula", "artemis",
		"virgingalactic", "virgin galactic", "blueorigin", "boeing", "starliner", "soyuz",

		"f22", "f-22", "jet",

		// Not interested in other stuff
		"doge", "fsd",

		"no tfr",

		// 3d models are nice, but we only care about "real info"
		"thanks", "thank you", "cheers", "render", "animation", "3d", "model", "speculation", "mysterious", "simulat", /* or/ed */

		"not starship", "non starship", "not about starship",

		// kerbal space program != "official" news
		"kerbal space program", "ksp",

		// not *that* kind of raptor
		"velociraptor", "jurassic", "cretaceous", "dino",

		"turtle", "dog", "cat",

		"ocisly", "canaveral",

		"bot",

		"dearmoon", "dear moon", "inspiration4",

		"sale", "buy", "shop", "store", "giveaway", "give away", "retweet", "birthday", "download", "click", "tag",
	}

	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`\bsn\d+\b`),
		// Booster BNx
		regexp.MustCompile(`\bbn\d+\b`),
	}

	closureTFRRegex = regexp.MustCompile("(?:closure|tfr|cryo|fts|scrub)")
	alertRegex      = regexp.MustCompile("(?:alert|static fire|closure|cryo|evacua|scrub)")

	// Users known to post better information that requires specific filtering
	specificUserMatchers = map[string]*regexp.Regexp{
		// One of the most important sources, gets alerted when the village has to evacuate for a flight
		"bocachicagal":    alertRegex,
		"starshipboca":    alertRegex,
		"bocachicamaria1": alertRegex,

		// These people likely tweet about test & launch stuff
		"rgvaerialphotos": closureTFRRegex,
		"bocaroad":        closureTFRRegex,
		"infographictony": closureTFRRegex,
		"spacex360":       closureTFRRegex,
		"bluemoondance74": closureTFRRegex,
		"nextspaceflight": closureTFRRegex,
		"tylerg1998":      closureTFRRegex,
		"nasaspaceflight": closureTFRRegex,
		"spacexboca":      closureTFRRegex,

		// Watches temporary flight restrictions
		"spacetfrs": regexp.MustCompile("(?:brownsville)"),

		// For Elon, we try to match anything that could be insider info
		"elonmusk": regexp.MustCompile("(?:booster|orbit|heavy|cryo|static fire|tower|ship|rud|engine|faa|starbase|boca chica|lox|liquid oxygen|methane|ch4|relight|fts|flip|cargo|lunar|tfr|fts|scrub|mach)"),
	}

	usersWithNoAntikeywords = map[string]bool{
		"elonmusk": true,
		"spacex":   true,
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

	text := tweet.Text()

	// We do not care about tweets that are timestamped with a text more than 24 hours ago
	// e.g. if someone posts a photo and then writes "took this on March 15, 2002"
	if d, ok := util.ExtractDate(text); ok && time.Since(d) > 24*time.Hour {
		return false
	}

	text = strings.ToLower(text)

	if strings.Contains(text, "patreon") && hasNoMedia(tweet) {
		return false
	}

	if isSatireAccount(tweet) {
		return false
	}

	// Now check if the text of the tweet matches what we're looking for.
	// if it's elon musk, then we don't check for anti-keywords
	if StarshipText(text, tweet.User != nil && usersWithNoAntikeywords[strings.ToLower(tweet.User.ScreenName)]) {
		return true
	}

	// Now check if we have a matcher for this specific user.
	// These users usually post high-quality information
	if tweet.User != nil {
		m, ok := specificUserMatchers[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			return m.MatchString(text)
		}
	}

	return false
}

func hasNoMedia(tweet *twitter.Tweet) bool {
	return (tweet.ExtendedEntities == nil || len(tweet.ExtendedEntities.Media) == 0) &&
		(tweet.Entities == nil || len(tweet.Entities.Media) == 0)
}

var (
	// See https://twitter.com/ULASeagull/status/1376913976362217472 and
	// https://twitter.com/i/lists/1357527189370130432 for a list of names
	satireNames = []string{}

	satireKeywords = []string{
		"parody", "joke",
	}
)

const SatireListID = 1377136100574064647

func LoadSatireList(client *twitter.Client) {
	users, _, err := client.Lists.Members(&twitter.ListsMembersParams{
		ListID: SatireListID,
		Count:  1000,
	})
	if err != nil {
		log.Println("[Twitter] Failed loading satire account list:", err.Error())
		return
	}

	for _, u := range users.Users {
		satireNames = append(satireNames, strings.ToLower(u.ScreenName))
	}
}

func isSatireAccount(tweet *twitter.Tweet) bool {
	if tweet.User == nil {
		return false
	}
	// If we know the user, it can't be satire
	_, known1 := specificUserMatchers[tweet.User.ScreenName]
	_, known2 := usersWithNoAntikeywords[tweet.User.ScreenName]
	if known1 || known2 {
		return false
	}

	username := strings.ToLower(tweet.User.ScreenName)

	for _, k := range satireNames {
		if username == k {
			return true
		}
	}

	desc := strings.ToLower(tweet.User.Description)
	for _, k := range satireKeywords {
		if strings.Contains(desc, k) {
			return true
		}
	}

	return false
}
