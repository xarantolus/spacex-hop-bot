package match

import (
	"log"
	"regexp"
	"strings"
	"time"
	"unicode"

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
		"orbital launch pad", "orbital launch mount",
		"olp service tower",
	}

	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`\b((s\d{2,}\b)|(sn|starship|starship number)\s?\d+['’]?s?)`),
		// Booster BNx
		regexp.MustCompile(`\b((b\d{1,2}\b)|(bn|booster|booster number)\s?\d+['’]?s?)`),
		// Yes. I like watching tanks
		regexp.MustCompile(`\b(gse)\s?(?:tank|-)?\s?\d*\b`),
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

		"sheriffgarza": regexp.MustCompile(`(?:close|closure|spacex)`),

		"austinbarnard45": regexp.MustCompile("(?:day in Texas)"),

		// Watches temporary flight restrictions
		"spacetfrs": regexp.MustCompile("(?:brownsville)"),

		// For Elon, we try to match anything that could be insider info
		"elonmusk": regexp.MustCompile("(?:booster|heavy|cryo|static fire|tower|ship|rud|faa|starbase|boca chica|lox|liquid oxygen|methane|ch4|relight|fts|flip|cargo|lunar|tfr|fts|scrub|mach)"),
	}

	userAntikeywordsOverwrite = map[string][]string{
		"elonmusk": {"tesla", "model s", "model 3", "model x", "model y", "car", "giga", "falcon", "boring company", "tunnel", "loop", "doge", "ula", "tonybruno", "jeff", "fsd"},
	}

	hqMediaAccounts = map[string]bool{
		"starshipgazer": true,
	}

	antiStarshipKeywords = []string{
		"electron", "blue origin", "neutron", "rocket lab", "rocketlab", "falcon", "f9", "starlink",
		"tesla", "openai", "boring", "hyperloop", "solarcity", "neuralink", "sls", "ula", "vulcan", "artemis",
		"virgingalactic", "virgin galactic", "virgin orbit", "virginorbit", "blueorigin", "boeing", "starliner", "soyuz", "orion",
		"resilience", "shuttle", "new glenn", "china", "chinese", "russia", "new shepard", "tsla", "dynetics", "hls",
		"ares", "titan", "ariane", "srb", "solid rocket booster", "terran", "relativity space", "relativityspace", "astra",
		"spaceshipthree", "spaceshiptwo", "spaceshipone", "vss enterprise", "vss imagine", "samsung", "bezos", "branson",
		"masten", "centaur",

		"amazon", "kuiper",

		// e.g. crew-1, crew-2...
		"crew-", "crew dragon", "dragon",

		"f22", "f-22", "jet",

		// Not interested in other stuff
		"doge", "coin", "btc", "fsd", "spce", "dogecoin",

		"no tfr",

		// "super heavyweight" in olympics...
		"super heavyweight",

		// 3d models are nice, but we only care about "real info"
		"thanks", "thank you", "cheers", "render", "animation", "3d", "model", "speculati" /*ng/on*/, "simulated", "print", "vfx", "not real", "photoshop",
		"artwork", "artist", "#art",

		"not starship", "non starship", "not about starship", "discord",

		// kerbal space program, games, star wars != "official" news
		"kerbal space program", "ksp", "no mans sky", "nomanssky", "no man’s sky", "no man's sky", "kerbals", "pocket rocket", "pocketrocket",
		"star trek", "startrek", "starcitizen", "star citizen", "battle droid", "b1-series", "civil war", "jabba the hutt",

		// KSP planets, moons, stars etc.
		"moho", "gilly", "kerbin", "mun", "minmus", "duna", "jool", "laythe", "vall",
		"tylo", "bop", "pol", "dres", "eeloo", "kerbol",

		// not *that* kind of raptor
		"velociraptor", "jurassic", "cretaceous", "dino",

		"aircraft", "aerial refueling", "firepower",

		// "Star one brazilsat B4"
		"brasilsat", "star one",

		// Someone *really* named their delivery robot business "Starship"... why?
		"delivery", "startup", "groceries", "robots", "starship robot",

		"ocisly", "jrti", "canaveral",

		"meme", "suck", "cursed", "uwu", "cult", "qwq", "reaction", "immigrants",

		"dearmoon", "dear moon", "inspiration4", "rover", "alien",

		"sale", "buy", "shop", "store", "giveaway", "give away", "retweet", "birthday", "download", "click", "tag", "discount",
		"follow", "pre-order", "merch", "vote", "podcast", "trending",

		"child", "illegal", "nfl", "tiktok", "tik tok", "self harm", "sex", "cock", "s3x", "gspot",

		// stuff that seems like starship, but isn't
		"starshipent", "monstax", "eshygazit", "wonho",

		// Account follows a sheriff
		"assault", "rape", "deadly", "weapon", "victim", "murder", "crime", "investigat", "body", "memorial",
	}
)

const (
	// TODO: find IDs for "Mesa del Gavilan", Stargate and generally places around/between the site.
	// The data seems to come from foursquare, but the IDs are *not* the same on both services

	// https://twitter.com/places/124bed061054f000
	SpaceXBuildSiteID = "124bed061054f000"
	// https://twitter.com/places/124cb6de55957000
	SpaceXLaunchSiteID = "124cb6de55957000"
	// https://twitter.com/places/1380f3b60f972001
	StarbasePlaceID = "1380f3b60f972001"
)

// StarshipText returns whether the given text mentions starship
func StarshipText(text string, antiKeywords []string) bool {
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

	// Raptor has more than one meaning, so we need to be more careful
	if strings.Contains(text, "raptor") && containsAny(text, "starship", "vacuum", "spacex", "mcgregor", "engine", "rb", "rc", "rvac", "raptorvan", "launch site", "production site", "booster", "super heavy", "superheavy", "truck") {
		return true
	}

	// The phobos and deimos oil rigs that will be used as sea-spaceports
	if containsAny(text, "deimos", "phobos") && containsAny(text, "spacex", "starship", "super heavy", "superheavy", "sea launch", "oil") {
		return true
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

	// We ignore certain satire accounts
	if isIgnoredAccount(&tweet.Tweet) {
		return false
	}

	// Now check if the text of the tweet matches what we're looking for.
	text = strings.ToLower(text)

	// Depending on the user, we use different antiKeywords
	antiKeywords := antiStarshipKeywords
	if tweet.User != nil {
		ak, ok := userAntikeywordsOverwrite[strings.ToLower(tweet.User.ScreenName)]
		if ok {
			antiKeywords = ak
		}
	}

	// Check if the text matches
	if StarshipText(text, antiKeywords) {
		return true
	}

	var containsBadWords = containsAntikeyword(antiKeywords, text)

	// If the tweet is tagged with Starbase as location, we just retweet it
	// TODO: Maybe only if it has media, not sure
	if tweet.Place != nil && !containsBadWords && (tweet.Place.ID == StarbasePlaceID || tweet.Place.ID == SpaceXLaunchSiteID || tweet.Place.ID == SpaceXBuildSiteID) {
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
	return startsWithAny(text, words...)
}

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

		for currentIndex < len(text) && (unicode.IsSpace(rune(text[currentIndex])) || rune(text[currentIndex]) == '#' || rune(text[currentIndex]) == '@') {
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
