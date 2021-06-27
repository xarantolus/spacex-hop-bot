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
	// we also match Raptor, but only if either "SpaceX", "Engine" or "McGregor" is mentioned
	starshipKeywords = []string{
		"starship",
		"superheavy", "super heavy",
		"orbital launch tower", "orbital tower", "olt segment", "launch tower segment", "olp service tower", "olp tower",
		"orbital launch integration tower",
		"gse tank",
		"orbital launch table", "orbital table",
		"orbital launch pad",
	}

	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`\bsn\d+\b`),
		// Booster BNx
		regexp.MustCompile(`\bbn\d+\b`),
	}

	closureTFRRegex = regexp.MustCompile("\b(?:closure|tfr|cryo|fts|scrub)")
	alertRegex      = regexp.MustCompile("\b(?:alert|static fire|closure|cryo|evac|scrub)")

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

		"sheriffgarza": regexp.MustCompile(`(?:close)`),

		"austinbarnard45": regexp.MustCompile("(?:day in Texas)"),

		// Watches temporary flight restrictions
		"spacetfrs": regexp.MustCompile("(?:brownsville)"),

		// For Elon, we try to match anything that could be insider info
		"elonmusk": regexp.MustCompile("(?:booster|heavy|cryo|static fire|tower|ship|rud|faa|starbase|boca chica|lox|liquid oxygen|methane|ch4|relight|fts|flip|cargo|lunar|tfr|fts|scrub|mach)"),
	}

	userAntikeywordsOverwrite = map[string][]string{
		"elonmusk": {"tesla", "model s", "model 3", "model x", "model y", "car", "giga", "falcon", "boring company", "tunnel", "loop", "doge"},
	}

	hqMediaAccounts = map[string]bool{
		"starshipgazer": true,
	}

	antiStarshipKeywords = []string{
		"electron", "blue origin", "neutron", "rocket lab", "rocketlab", "falcon", "f9", "starlink",
		"tesla", "openai", "boring", "hyperloop", "solarcity", "neuralink", "sls", "ula", "vulcan", "artemis",
		"virgingalactic", "virgin galactic", "blueorigin", "boeing", "starliner", "soyuz", "orion",
		"resilience", "shuttle", "new glenn", "china", "chinese", "russia", "new shepard", "tsla", "dynetics", "hls",
		"ares", "titan", "ariane", "srb", "solid rocket booster", "terran", "relativity space", "relativityspace",

		// e.g. crew-1, crew-2...
		"crew-", "crew dragon",

		"f22", "f-22", "jet",

		// Not interested in other stuff
		"doge", "coin", "btc", "fsd", "spce",

		"no tfr",

		// 3d models are nice, but we only care about "real info"
		"thanks", "thank you", "cheers", "render", "animation", "3d", "model", "speculation", "mysterious", "simulat" /* or/ed */, "print", "vfx",

		"not starship", "non starship", "not about starship", "discord",

		// kerbal space program != "official" news
		"kerbal space program", "ksp", "no mans sky", "nomanssky", "no man’s sky", "no man's sky", "kerbals", "pocket rocket", "pocketrocket",
		"star trek", "startrek", "starcitizen", "star citizen",

		// not *that* kind of raptor
		"velociraptor", "jurassic", "cretaceous", "dino",

		// Someone *really* named their delivery robot business "Starship"... why?
		"delivery", "startup", "groceries", "robots", "starship robot",

		"ocisly", "jrti", "canaveral",

		"meme", "bot", "suck", "cursed", "uwu", "qwq", "reaction", "immigrants",

		"dearmoon", "dear moon", "inspiration4", "rover", "alien",

		"sale", "buy", "shop", "store", "giveaway", "give away", "retweet", "birthday", "download", "click", "tag", "discount",
		"follow", "pre-order", "merch", "vote", "podcast", "trending",

		"child", "illegal", "nfl", "tiktok", "tik tok", "self harm",

		// stuff that seems like starship, but isn't
		"starshipent", "monstax", "eshygazit", "wonho",

		// Account follows a sheriff
		"assault", "rape", "deadly", "weapon", "victim", "murder", "crime", "investigat", "body", "memorial",
	}
)

// StarshipText returns whether the given text mentions starship
func StarshipText(text string, antiKeywords []string) bool {

	text = strings.ToLower(text)

	if containsAntikeyword(antiKeywords, text) {
		return false
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
func StarshipTweet(tweet TweetWrapper) bool {
	// Ignore OLD tweets
	if d, err := tweet.CreatedAtTime(); err == nil && time.Since(d) > 24*time.Hour {
		return false
	}

	text := tweet.Text()

	// We do not care about tweets that are timestamped with a text more than 24 hours ago
	// e.g. if someone posts a photo and then writes "took this on March 15, 2002"
	if d, ok := util.ExtractDate(text); ok && time.Since(d) > 48*time.Hour {
		return false
	}

	text = strings.ToLower(text)

	if isIgnoredAccount(&tweet.Tweet) {
		return false
	}

	// Now check if the text of the tweet matches what we're looking for.

	antiKeywords := antiStarshipKeywords
	if tweet.User != nil {
		ak, ok := userAntikeywordsOverwrite[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			antiKeywords = ak
		}
	}

	if StarshipText(text, antiKeywords) {
		return true
	}

	// Raptor has more than one meaning, so we need to be more careful
	if !containsAntikeyword(antiKeywords, text) && strings.Contains(text, "raptor") && (strings.Contains(text, "starship") || strings.Contains(text, "spacex") || strings.Contains(text, "mcgregor") || strings.Contains(text, "engine")) {
		return true
	}

	if !containsAntikeyword(antiKeywords, text) && (strings.Contains(text, "starbase") || strings.Contains(text, "boca chica")) && (strings.Contains(text, "launch site") || strings.Contains(text, "launch tower")) {
		return true
	}

	// Now check if we have a matcher for this specific user.
	// These users usually post high-quality information
	if tweet.User != nil {
		m, ok := specificUserMatchers[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			return m.MatchString(text)
		}

		// There are some accounts that always post high-quality pictures and videos.
		// For them we retweet *everything* that has media
		if hqMediaAccounts[strings.ToLower(tweet.User.ScreenName)] {
			return hasMedia(&tweet.Tweet)
		}
	}

	return false
}

func hasMedia(tweet *twitter.Tweet) bool {
	return tweet.Entities != nil && len(tweet.Entities.Media) > 0 || tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0
}

func containsAntikeyword(words []string, text string) bool {
	for _, k := range words {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}
