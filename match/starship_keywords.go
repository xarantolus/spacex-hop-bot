// This file defines basically everything the matcher does by specifying positive and negative keywords
package match

import "regexp"

// keywordMapping basically defines two sets of keywords.
// if at least one keyword from `from` and one from `to` is matched,
// then the match is positive
type keywordMapping struct {
	from, to, antiKeywords []string
}

// Note that all text here must be lowercase because the text is lowercased in the matching function
var (
	// If at least one of these keywords is present in a tweet (and no antiKeywords are),
	// we should retweet
	starshipKeywords = []string{
		"starship",
		"superheavy", "super heavy",

		"orbital launch tower", "orbital tower", "olt segment",
		"launch tower segment", "olp service tower", "olp tower",
		"orbital launch integration tower", "launch tower arm",

		"wide bay", "high bay",

		"orbital tank farm",

		"orbital launch table", "orbital table",
		"orbital launch pad", "orbital launch mount",
		"suborbital pad", "suborbital launch pad",
		"olp service tower",
		"orbital launch site",
	}

	// starshipMatchers are more specific regexes that act like starshipKeywords
	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`\b((s-?\d{2,}\b)|(ship\s?\d{2,}\b)|(sn|starship|starship number)-?\s?\d+['’]?s?)`),
		// Booster BNx
		regexp.MustCompile(`(((?:#|\s|^)b\d{1,2}\b([^-]|$))|\b(bn|booster|booster number)(['’]|s)*\s?\d{1,3}['’]?s?\b)`),
		// Yes. I like watching tanks
		regexp.MustCompile(`\b(gse)\s?(?:tank|-)?\s?\d+\b`),
		// Raptor with a number
		regexp.MustCompile(`\b((?:raptor|raptor\s+engine|rvac|rb|rc)(?:\s+(?:center|centre|boost|vacuum))?(?:\s+engine)?\s*\d+)\b`),
	}

	// moreSpecificKeywords are keywords that must be accompanied by at least one of the keywords mentioned in their slice.
	// This is useful for "raptor" (to make sure we only get engines) and some launch sites
	// The compose() function can be used to combine multiple slices.
	// It does NOT make sense to put starshipKeywords into any of these slices, because if
	// we reach the point where we look for more specific keywords, none of the starshipKeywords has matched
	moreSpecificKeywords = []keywordMapping{
		// Engines
		{
			from: []string{"raptor"},
			to: []string{
				"starship", "vacuum", "sea-level", "sea level",
				"spacex", "mcgregor", "engine", "rb", "rc", "rvac",
				"launch site", "production site", "booster", "super heavy",
				"superheavy", "truck", "van", "raptorvan", "deliver", "sea level",
				"flare", "high bay", "nozzle", "tripod", "starbase", "static fire",
			},
		},

		// Stuff noticed on live streams
		{
			from: compose(liveStreams,
				[]string{"orbital tank farm", "otf"},
				[]string{"suborbital tank farm", "stf"},
				[]string{"olm", "olt", "olit"},
			),
			to: []string{
				"methane", "tank", "lox", "ch4", "lch4", "ln2", "frost", "fire", "vent",
				"argon", "pad", "road", "highway", "close", "open", "qd", "quick disconnect",
				"raptor", "cranex",
			},
		},

		// Cranes lifting stuff like boosters etc.
		{
			from: compose([]string{"cranex", "liebherr lr"}),
			to:   compose(liveStreams, generalSpaceXKeywords, nonSpecificKeywords),
		},

		// Ground infrastructure
		{
			from: []string{"gse tank"},
			to:   compose(nonSpecificKeywords, generalSpaceXKeywords),
		},

		// Testing activity
		{
			from: []string{"cryogenic proof", "cryo proof", "cryoproof"},
			to:   compose(nonSpecificKeywords, generalSpaceXKeywords),
		},
		{
			from: []string{"road closure", "temporary flight restriction", "tfr "},
			to:   compose(nonSpecificKeywords, generalSpaceXKeywords),
		},

		// Seaports/Oil rigs that might be used for launches/landings?
		{
			from: []string{"deimos"},
			to:   compose(seaportKeywords, generalSpaceXKeywords, []string{"phobos"}),
		},
		{
			from: []string{"phobos"},
			to:   compose(seaportKeywords, generalSpaceXKeywords, []string{"deimos"}),
		},

		// New launch pads at different locations
		{
			from: compose(
				[]string{"lc-49", "lc 49", "launch complex 49", "launch complex-49"},

				// Don't match this one as it's currently in use and I have no idea how to differentiate starship tweets from falcon ones
				// []string{"lc-39a", "lc 39a", "launch complex 39a", "launch complex-39a"},
				generalSpaceXKeywords,
			),
			to: compose([]string{"environmental assessment", "tower"}),
		},

		// Launch tower
		{
			from: compose([]string{"mechazilla", "olit"}),
			to: compose(
				placesKeywords, nonSpecificKeywords,
				[]string{"qd", "arm", "catch"},
			),
		},

		// Some words that are usually ambigious, but if combined with starship keywords they are fine
		{
			from: []string{"launch tower", "launch pad", "launch mount", "chopstick", "chop stick", "catch arms"},
			to:   compose(seaportKeywords, placesKeywords, liveStreams),
		},
	}

	// Helper slices that can be used for composing new keywords
	seaportKeywords       = []string{"sea launch", "oil", "rig"}
	placesKeywords        = []string{"starbase", "boca chica"}
	nonSpecificKeywords   = compose([]string{"ship", "booster"}, liveStreams, placesKeywords)
	generalSpaceXKeywords = []string{"spacex"}
	liveStreams           = []string{
		// 24/7 live camera views are often mentioned when something is shown on a screenshot
		"labpadre", "nasaspaceflight",
		// Other streamers
		"jessica_kirsh", "bocachicagal", "starship gazer", "starshipgazer",
	}

	// Regexes for road closures and testing activity
	closureTFRRegex = regexp.MustCompile(`\b(?:closure|tfr|notmar|cryo|fts|scrub)`)
	alertRegex      = regexp.MustCompile(`\b(?:alert|static fire|closure|cryo|evac|scrub|pad.*clear|clear.*pad)`)

	// Users that are known to post better information that requires less filtering.
	// The regexes are combined as OR, which means that only one has to match for a successful match
	specificUserMatchers = map[string][]*regexp.Regexp{
		// One of the most important sources, gets alerted when the village has to evacuate for a flight
		"bocachicagal":    {alertRegex, closureTFRRegex},
		"starshipboca":    {alertRegex, closureTFRRegex},
		"bocachicamaria1": {alertRegex, closureTFRRegex},

		// Photographers usually at the place
		"austindesisto": {alertRegex, closureTFRRegex},
		"starshipgazer": {alertRegex, closureTFRRegex},

		// These people likely tweet about test & launch stuff
		"nasaspaceflight": {closureTFRRegex, alertRegex},
		"spacex360":       {closureTFRRegex, alertRegex},
		"rgvaerialphotos": {closureTFRRegex},
		"bocaroad":        {closureTFRRegex},
		"bluemoondance74": {closureTFRRegex},
		"nextspaceflight": {closureTFRRegex},
		"tylerg1998":      {closureTFRRegex},
		"spacexboca":      {closureTFRRegex},

		"sheriffgarza": {regexp.MustCompile(`(?:close|closure|spacex)`)},

		// Always retweet the timelapse by this bot
		"starbasepulse": {regexp.MustCompile(`(?:timelapse|time lapse)`)},

		// Watches temporary flight restrictions
		"spacetfrs": {regexp.MustCompile("(?:brownsville)")},

		// For Elon, we try to match anything that could be insider info
		"elonmusk": {
			regexp.MustCompile(`(?:booster|cryo|static fire|tower|ship|rud|faa|starbase|boca chica|lox|liquid oxygen|methane|ch4|relight|fts|flip|cargo|lunar|tfr|fts|scrub|flap)`),
			// Try to match things for orbital flight tests
			regexp.MustCompile(`(?:orbit(?:.|\s)+(flight test|test flight)|(flight test|test flight)(?:.|\s)+orbit)`),
		},
	}

	userAntikeywordsOverwrite = map[string][]string{
		"elonmusk": {"tesla", "model s", "model 3", "model x", "model y", "car", "giga", "falcon", "boring company", "tunnel", "loop", "doge", "ula", "tonybruno", "jeff", "fsd", "giga berlin", "giga factory", "gigafactory", "giga press"},

		// NASA Accounts that sometimes tweet about Starship don't need any antiKeywords - they are "allowed"
		// to mention Starship together with e.g. Orion (which would be ignored if not for these overrides).
		"nasa":          {},
		"nasajpl":       {},
		"nasa_marshall": {},
		"nasa_gateway":  {},
		"nasaartemis":   {},
		"nasakennedy":   {},
		"nasagoddard":   {},
	}

	// Accounts that post only Starship photos - so if they post a picture, they
	// are retweeted automatically
	hqMediaAccounts = map[string]bool{
		"starshipgazer": true,
		"cnunezimages":  true,
	}

	// Accounts that are *never* considered satire accounts, even if they were on a list of these accounts
	veryImportantAccounts = map[string]bool{
		"elonmusk": true,
		"spacex":   true,
	}

	// If an account has any of these words in its description, we don't retweet tweets from it
	ignoredAccountDescriptionKeywords = []string{
		// Parody accounts
		"parody", "joke",

		// 3D artists
		"blender", "3d", "vfx", "render", "animat", /* e/ion */

		// Sports stuff
		"nhl",
	}

	// If a tweet contains any of these keywords, it will not be retweeted. This is a way of filtering out *non-starship* stuff
	antiStarshipKeywords = []string{
		"electron", "blue origin", "neutron", "rocket lab", "rocketlab", "hungry hippo", "rklb", "falcon", "merlin", "m1d", "f9", "starlink",
		"tesla", "rivian", "giga press", "gigapress", "gigafactory", "openai", "boring", "hyperloop", "solarcity", "neuralink", "sls", "nasa_sls", "ula", "vulcan", "artemis", "rogozin",
		"virgingalactic", "virgin galactic", "virgin orbit", "virginorbit", "blueorigin", "boeing", "starliner", "soyuz", "soviet", "orion",
		"resilience", "shuttle", "challenger", "sts-51l", "sts-33", "new glenn", "china", "shenzhou", "india", "chinese", "japan", "space plane", "russia", "new shepard", "tsla", "dynetics", "hls",
		"ares", "titan", "ariane", "srb", "solid rocket booster", "terran", "relativity space", "relativityspace", "astra",
		"spaceshipthree", "spaceshiptwo", "spaceshipone", "vss enterprise", "starship enterprise", "archer", "sisko", "vss imagine", "galaxy note", "galaxy s", "bezos", "jeff who", "branson", "tory", "bruno",
		"masten", "centaur", "atlas v", "atlasv", "relativity", "northrop grumman", "northropgrumman", "bomber",
		"rookisaacman", "cygnus", "samsung", "angara", "firefly", "rolls-royce", "agrifood", "iot", "vs-50", "solid-propellant", "solid propellant",
		"são paulo", "sao paulo", "vlm-", "ac1", "arca", "ecorocket", "korea", "nuri", "mars rover", "perseverance", "curiosity", "ingenuity", "zhurong",

		"roscosmos", "yenisey",

		"hubble", "nasahubble",

		// Blue Origins' Starship... kind of clone i guess?
		"jarvis", "glenn", "bob smith",

		"be4", "be-4", "be 4 engine",

		"amazon", "kuiper", "nasaartemis", "isro",

		"7news", "bay light",

		// e.g. crew-1, crew-2...
		"crew-", "crew dragon", "dragon", "crs", "dm-",

		"f22", "f-22", "jet", "b-52", "s-300",

		// Not interested in other stuff
		"doge", "coin", "btc", "fsd", "spce", "dogecoin", "crypto", "safemoon", "stock", "wall street", "wallstreet", "buffett",

		"no tfr",

		"volvo",

		// "super heavyweight" in olympics...
		"super heavyweight",

		"parachute",

		"supernova",

		"god", "the lord",

		// 3d models are nice, but we only care about "real info"
		"thanks", "thank you", "cheers", "render", "animat" /* ion/ed */, "3d", "model", "speculati" /*ng/on*/, "simulated", "print", "vfx", "not real", "photoshop",
		"art", "mission patch", "drawing", "board game", "starshipshuffle", "starship shuffle", "lego",
		"card game", "starship design", "daily_hopper", "daily hopper", "paper model", "papermodel",

		"8bitdo", "sn30",

		"gaming",

		"your guess",

		"not starship", "non starship", "not about starship", "discord", "wonder if", "was wondering", "years ago", "year ago",

		// kerbal space program, games, star wars != "official" news
		"kerbal space program", "ksp", "no mans sky", "nomanssky", "no man’s sky", "no man's sky", "kerbals", "pocket rocket", "pocketrocket", "simplerockets",
		"star trek", "startrek", "starcitizen", "star citizen", "battle droid", "b1-series", "civil war", "jabba the hutt", "sfs", "space flight simulator",
		"rocket explorer",

		"tax",

		// KSP planets, moons, stars etc.
		"moho", "gilly", "kerbin", "mun", "minmus", "duna", "jool", "laythe", "vall",
		"tylo", "bop", "pol", "dres", "eeloo", "kerbol",

		// movies
		"the martian",

		// not *that* kind of raptor
		"velociraptor", "jurassic", "cretaceous", "dino",

		"suprem", "aryan",

		"aircraft", "aerial refueling", "firepower",

		// "Star one brazilsat B4"
		"brasilsat", "star one",

		// Someone *really* named their delivery robot business "Starship"... why?
		"startup", "groceries", "delivery robots", "starship robot",

		// And there's of course a wide range of products named S(numbers), which is annoying
		"tern", "gsd", "cargo bike",

		// Vitamin B2
		"vitamin",

		"ocisly", "jrti", "asog", "canaveral",

		"obetraveller", "ocean cam", "oceancam", "oceanscam", "paul",

		// I do not care about opinions on starship
		"agree", "disagree", "throwback to", "opinion", "imo", "imho", "i think", "mfw", "vibe", "dream", "laughs in",

		"gorgeous girl",

		"meme", "ratio", "apology", "drama", "petition to", "suck", "cursed", "uwu", "cult", "qwq", "reaction", "immigrant",

		"dearmoon", "dear moon", "inspiration4", "inspiration 4", "inspiration four", "alien",

		"sale", "buy", "shop", "store", "purchase", "shirt", "sweater", "giveaway", "give away", "retweet", "birthday", "discount",
		"pre-order", "merch", "vote", "podcast", "trending", "hater", "follower", "unfollow", "top friends", "plush", "black friday", "blackfriday", "newprofilepic",

		"child", "kid", "illegal", "nfl", "tiktok", "vax", "vacc", "booster shot", "shoot", "tik tok", "self harm", "sex", "cock", "s3x", "gspot", "g-spot", "fuck", "dick", "bullshit", "bikini",
		"booty", "cudd", "bathroom", "penis", "vagi", "furry", "strok", "fap", "chick",

		"stfu",

		"patrons", "babylon", "boltup", "champion",

		"red bull", "browns",

		"tier list",

		// Annoying elon musk quotes
		"consciousness",

		// Some conferences have a "stand B20", because why not trick this bot right?
		"booth", "stand b",

		"trump", "antifa", "biden", "riot", "taliban", "kill", "beat", "ideology", "gender",

		"surgery", "emergency",

		// Things that are typical questions for polls. We cannot get polls using the Twitter v1 API, so this is kind of bad
		"feel about", "vs",

		// stuff that seems like starship, but isn't
		"starshipent", "monstax", "eshygazit", "wonho",

		// Account follows a sheriff
		"arrest", "violence ", "assault", "rape", "weapon", "victim", "murder", "crime", "investigat", "body", "nigg", "memorial", "dead", "death", "cancer", "piss",
		"abus",

		"nonce", "pedo",

		"bomb",

		"hospital",

		"offend", "offensive", "fanboy", "fan boy", "fangirl", "fan girl",

		"covid", "corona", "rona", "omicron", "tests positive", "positive test", "cdc",

		"shit", "anime", "manga", "bronco", "bae",

		"abortion", "roe v. wade", "roe v wade", "roe vs. wade", "roe vs wade",

		"starshipcongrss", "starshipcongress", "congress", "starflight academy",

		// Make sure we don't retweet anything related to horrible tragedies
		"9/11", "911", "twin tower", "wtc", "trade center", "die", "falling",
	}
)
