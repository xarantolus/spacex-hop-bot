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
		"launch tower arm", "mechazilla",
	}

	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`\b((s-?\d{2,}\b)|(ship\s?\d{2,}\b)|(sn|starship|starship number)-?\s?\d+['’]?s?)`),
		// Booster BNx
		regexp.MustCompile(`(((?:#|\s|^)b-?\d{1,2}\b)|\b(bn|booster|booster number)-?\s?\d{1,3}['’]?s?\b)`),
		// Yes. I like watching tanks
		regexp.MustCompile(`\b(gse)\s?(?:tank|-)?\s?\d*\b`),
		// Raptor with a number
		regexp.MustCompile(`\b((?:raptor|raptor\s+engine|rvac|rb|rc)\s*\d+)\b`),
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

		// Watches temporary flight restrictions
		"spacetfrs": regexp.MustCompile("(?:brownsville)"),

		// For Elon, we try to match anything that could be insider info
		"elonmusk": regexp.MustCompile("(?:booster|heavy|cryo|static fire|tower|ship|rud|faa|starbase|boca chica|lox|liquid oxygen|methane|ch4|relight|fts|flip|cargo|lunar|tfr|fts|scrub|flap)"),
	}

	userAntikeywordsOverwrite = map[string][]string{
		"elonmusk": {"tesla", "model s", "model 3", "model x", "model y", "car", "giga", "falcon", "boring company", "tunnel", "loop", "doge", "ula", "tonybruno", "jeff", "fsd", "giga berlin", "giga factory", "gigafactory", "giga press"},
	}

	hqMediaAccounts = map[string]bool{
		"starshipgazer": true,
		"cnunezimages":  true,
	}

	veryImportantAccounts = map[string]bool{
		"elonmusk": true,
		"spacex":   true,
	}

	antiStarshipKeywords = []string{
		"electron", "blue origin", "neutron", "rocket lab", "rocketlab", "rklb", "falcon", "f9", "starlink",
		"tesla", "giga press", "gigapress", "gigafactory", "openai", "boring", "hyperloop", "solarcity", "neuralink", "sls", "nasa_sls", "ula", "vulcan", "artemis", "rogozin",
		"virgingalactic", "virgin galactic", "virgin orbit", "virginorbit", "blueorigin", "boeing", "starliner", "soyuz", "orion",
		"resilience", "shuttle", "challenger", "sts-51l", "sts-33", "new glenn", "china", "shenzhou", "india", "chinese", "japan", "space plane", "russia", "new shepard", "tsla", "dynetics", "hls",
		"ares", "titan", "ariane", "srb", "solid rocket booster", "terran", "relativity space", "relativityspace", "astra",
		"spaceshipthree", "spaceshiptwo", "spaceshipone", "vss enterprise", "starship enterprise", "vss imagine", "samsung", "galaxy s", "bezos", "jeff who", "branson", "tory", "bruno",
		"masten", "centaur", "atlas v", "relativity", "northrop grumman", "northropgrumman", "bomber",
		"rookisaacman", "cygnus", "samsung", "angara", "firefly",

		"roscosmos", "yenisey",

		// Blue Origins' Starship... kind of clone i guess?
		"jarvis", "glenn", "bob smith",

		"amazon", "kuiper", "nasaartemis", "isro",

		// e.g. crew-1, crew-2...
		"crew-", "crew dragon", "dragon", "crs",

		"f22", "f-22", "jet",

		// Not interested in other stuff
		"doge", "coin", "btc", "fsd", "spce", "dogecoin", "crypto",

		"no tfr",

		"volvo",

		// "super heavyweight" in olympics...
		"super heavyweight",

		"god", "the lord",

		// 3d models are nice, but we only care about "real info"
		"thanks", "thank you", "cheers", "render", "animation", "3d", "model", "speculati" /*ng/on*/, "simulated", "print", "vfx", "not real", "photoshop",
		"artwork", "artist", "#art", "mission patch", "drawing", "board game", "starshipshuffle", "starship shuffle",
		"card game", "starship design", "daily_hopper", "daily hopper",

		"not starship", "non starship", "not about starship", "discord", "wonder if", "was wondering",

		// kerbal space program, games, star wars != "official" news
		"kerbal space program", "ksp", "no mans sky", "nomanssky", "no man’s sky", "no man's sky", "kerbals", "pocket rocket", "pocketrocket", "simplerockets",
		"star trek", "startrek", "starcitizen", "star citizen", "battle droid", "b1-series", "civil war", "jabba the hutt", "sfs", "space flight simulator",
		"rocket explorer",

		// KSP planets, moons, stars etc.
		"moho", "gilly", "kerbin", "mun", "minmus", "duna", "jool", "laythe", "vall",
		"tylo", "bop", "pol", "dres", "eeloo", "kerbol",

		// movies
		"the martian",

		// not *that* kind of raptor
		"velociraptor", "jurassic", "cretaceous", "dino",

		"aircraft", "aerial refueling", "firepower",

		// "Star one brazilsat B4"
		"brasilsat", "star one",

		// Someone *really* named their delivery robot business "Starship"... why?
		"startup", "groceries", "robots", "starship robot",

		// And there's of course a wide range of products named S(numbers), which is annoying
		"tern", "gsd", "cargo bike",

		"ocisly", "jrti", "canaveral",

		"obetraveller", "ocean cam", "oceancam", "oceanscam", "paul",

		// I do not care about opinions on starship
		"agree", "disagree", "throwback to",

		"meme", "ratio", "apology", "drama", "petition to", "suck", "cursed", "uwu", "cult", "qwq", "reaction", "immigrant",

		"dearmoon", "dear moon", "inspiration4", "inspiration 4", "inspiration four", "rover", "alien",

		"sale", "buy", "shop", "store", "giveaway", "give away", "retweet", "birthday", "download", "click", "tag", "discount",
		"pre-order", "merch", "vote", "podcast", "trending", "hater", "follow", "unfollow", "top friends",

		"child", "illegal", "nfl", "tiktok", "tik tok", "self harm", "sex", "cock", "s3x", "gspot", "g-spot", "fuck", "dick", "bullshit", "bikini",

		"red bull",

		"trump", "antifa", "biden", "riot", "taliban", "kill",

		// Things that are typical questions for polls. We cannot get polls using the Twitter v1 API, so this is kind of bad
		"feel about", "vs",

		// stuff that seems like starship, but isn't
		"starshipent", "monstax", "eshygazit", "wonho",

		// Account follows a sheriff
		"arrest", "violence ", "assault", "rape", "weapon", "victim", "murder", "crime", "investigat", "body", "memorial", "dead", "death", "cancer",

		"nonce", "pedo",

		"covid", "corona",

		// Make sure we don't retweet anything related to horrible tragedies
		"9/11", "911", "twin tower", "wtc", "trade center", "die", "death", "falling",
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
	// https://twitter.com/places/07d9f642af482000
	SpaceXMcGregorPlaceID = "07d9f642af482000"

	// Other places around the area:
	// "Isla Blanca Park": https://twitter.com/places/11dca9a728950001
	// "South Padre Island, TX": https://twitter.com/places/1d1f665883989434
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
	if strings.Contains(text, "raptor") && (startsWithAny(text, "starship", "vacuum", "sea-level", "sea level", "spacex", "mcgregor", "engine", "rb", "rc", "rvac", "launch site", "production site", "booster", "super heavy", "superheavy", "truck", "van", "raptorvan", "deliver", "sea level", "high bay")) {
		return true
	}

	// The phobos and deimos oil rigs that will be used as sea-spaceports
	if containsAny(text, "deimos", "phobos") && containsAny(text, "spacex", "starship", "super heavy", "superheavy", "sea launch", "oil") {
		return true
	}

	return false
}

// The faceRatio of a tweet is the number of faces in all images (or video thumbnails) divided by the number of images in the tweet
const maxFaceRatio = 1.1

var faceDetector = NewFaceDetector()

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

	// We ignore certain (e.g. satire, artist) accounts
	if tweet.User != nil {
		if _, important := veryImportantAccounts[strings.ToLower(tweet.User.Name)]; !important && IsOrMentionsIgnoredAccount(&tweet.Tweet) {
			return false
		}
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

	var containsBadWords = containsAntikeyword(antiKeywords, text)

	// If the tweet is tagged with Starbase as location, we just retweet it
	if !containsBadWords && IsAtSpaceXSite(&tweet.Tweet) {
		fr := faceDetector.FaceRatio(&tweet.Tweet)
		log.Printf("[FaceRatio] %s: %f\n", util.TweetURL(&tweet.Tweet), fr)
		return fr <= maxFaceRatio
	}

	// If the tweet mentions raptor without images, we still retweet it.
	// This is mostly for tweets from SpaceX McGregor
	if !containsBadWords && strings.Contains(text, "raptor") && IsAtSpaceXSite(&tweet.Tweet) {
		fr := faceDetector.FaceRatio(&tweet.Tweet)
		log.Printf("[FaceRatio] %s: %f\n", util.TweetURL(&tweet.Tweet), fr)
		return fr <= maxFaceRatio
	}

	// Now check if it mentions too many people
	if strings.Count(text, "@") > 5 {
		return false
	}

	// Check if the text matches
	if StarshipText(text, antiKeywords) {
		fr := faceDetector.FaceRatio(&tweet.Tweet)
		log.Printf("[FaceRatio] %s: %f\n", util.TweetURL(&tweet.Tweet), fr)
		return fr <= maxFaceRatio
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

func IsAtSpaceXSite(tweet *twitter.Tweet) bool {
	return tweet.Place != nil && (tweet.Place.ID == StarbasePlaceID ||
		tweet.Place.ID == SpaceXLaunchSiteID || tweet.Place.ID == SpaceXBuildSiteID ||
		tweet.Place.ID == SpaceXMcGregorPlaceID)
}

func hasMedia(tweet *twitter.Tweet) bool {
	return tweet.Entities != nil && len(tweet.Entities.Media) > 0 || tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0
}

func ContainsStarshipAntiKeyword(text string) bool {
	return containsAntikeyword(antiStarshipKeywords, strings.ToLower(text))
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

		for currentIndex < len(text) && (unicode.IsSpace(rune(text[currentIndex])) || strings.ContainsAny(string(rune(text[currentIndex])), "#@$")) {
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
