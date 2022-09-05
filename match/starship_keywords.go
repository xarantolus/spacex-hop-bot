// This file defines basically everything the matcher does by specifying positive and negative keywords
package match

import "regexp"

// keywordMapping basically defines two sets of keywords.
// if at least one keyword from `from` and one from `to` is matched,
// then the match is positive
type keywordMapping struct {
	from, to, antiKeywords []string
}

func (mapping *keywordMapping) matches(text string) bool {
	_, ok := startsWithAny(text, mapping.from...)
	if !ok {
		return false
	}
	_, ok = startsWithAny(text, mapping.to...)
	if !ok {
		return false
	}

	_, ok = startsWithAny(text, mapping.antiKeywords...)
	return !ok
}

// Note that all text here must be lowercase because the text is lowercased in the matching function
var (
	// If at least one of these keywords is present in a tweet (and no antiKeywords are),
	// we should retweet
	starshipKeywords = ignoreSpaces([]string{
		"starship",
		"super heavy",

		"star factory",

		"orbital launch tower", "orbital tower", "olt segment",
		"launch tower segment", "olp service tower", "olp tower",
		"orbital launch integration tower", "launch tower arm",

		"wide bay", "mega bay", "high bay",

		"orbital tank farm",

		"orbital launch table", "orbital table",
		"orbital launch pad", "orbital launch mount",
		"suborbital pad", "suborbital launch pad",
		"olp service tower",
		"orbital launch site",
	})

	// starshipMatchers are more specific regexes that act like starshipKeywords
	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`\b((s\d{2}\b)|(ship\s?\d{2}\b)|(sn-?|starship|starship number)\s?\d['’]?s?)`),
		// Booster BNx
		regexp.MustCompile(`(((?:#|\s|^)b\d{1,2}\b([^-]|$))|\b(bn|booster|booster number)(['’]|s)*\s?\d{1,3}['’]?s?\b)`),
		// Yes. I like watching tanks
		regexp.MustCompile(`\b(gse)\s?(?:tank|-)?\s?\d+\b`),
		// Raptor with a number
		regexp.MustCompile(`\b((?:raptor|raptor\s+engine|rvac|rb|rc)(?:\s+(?:center|centre|boost|vacuum))?(?:\s+engine)?\s*v?\d+)\b`),
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
			to: compose(liveStreams, ignoreSpaces([]string{
				"starship", "vacuum", "sea level",
				"spacex", "mcgregor", "engine", "rb", "rc", "rvac",
				"launch site", "production", "booster", "super heavy",
				"superheavy", "truck", "van", "raptorvan", "deliver",
				"flare", "high bay", "nozzle", "tripod", "starbase", "static fire",
			})),
		},
		{
			from: ignoreSpaces([]string{"mc gregor"}),
			to:   ignoreSpaces([]string{"tri pod", "raptor"}),
		},

		// Stuff noticed on live streams
		{
			from: ignoreSpaces(compose(liveStreams,
				[]string{"orbital tank farm", "otf"},
				[]string{"suborbital tank farm", "stf"},
				[]string{"olm", "olt", "olit"},
				[]string{"stage zero"},
				[]string{"ols"},
				[]string{"booster", "orbital test flight", "orbital flight test"},
				liveStreams, placesKeywords,
			)),
			to: ignoreSpaces([]string{
				"methane", "tanker", "lox", "ch4", "lch4", "ln2", "frost", "fire", "vent",
				"argon", "road", "highway", "hwy", "qd", "sqd", "bqd", "quick disconnect",
				"raptor", "crane x",
			}),
		},

		{
			from: compose(placesKeywords, sitesKeywords),
			to:   liveStreams,
		},

		// Cranes lifting stuff like boosters etc.
		{
			from:         ignoreSpaces([]string{"crane x", "liebherr lr", "grid fin", "fin ", "fins ", "flap ", "flaps"}),
			to:           compose(liveStreams, generalSpaceXKeywords, nonSpecificKeywords, placesKeywords, sitesKeywords),
			antiKeywords: ignoreSpaces([]string{"whale", "cold gas"}),
		},

		// Ground infrastructure
		{
			from: ignoreSpaces([]string{"gse tank"}),
			to:   compose(nonSpecificKeywords, generalSpaceXKeywords),
		},

		// Testing activity
		{
			from:         ignoreSpaces([]string{"cryo proof", "proof", "stack"}),
			to:           compose(nonSpecificKeywords, generalSpaceXKeywords, liveStreams),
			antiKeywords: []string{"twitter"},
		},
		{
			from: []string{"road closure", "temporary flight restriction", "tfr "},
			to:   compose(nonSpecificKeywords, generalSpaceXKeywords),
		},

		// Seaports/Oil rigs that might be used for launches/landings?
		{
			from: []string{"deimos"},
			to:   compose(seaportKeywords, generalSpaceXKeywords, []string{"phobos"}, liveStreams),
		},
		{
			from: []string{"phobos"},
			to:   compose(seaportKeywords, generalSpaceXKeywords, []string{"deimos"}, liveStreams),
		},

		{
			from: ignoreSpaces([]string{"aerial shots", "fly over"}),
			to:   compose(placesKeywords, sitesKeywords),
		},

		// New launch pads at different locations
		{
			from: compose(
				ignoreSpaces([]string{"lc-49", "lc 49", "launch complex 49", "launch complex-49"}),

				// Don't match this one as it's currently in use and I have no idea how to differentiate starship tweets from falcon ones
				// []string{"lc-39a", "lc 39a", "launch complex 39a", "launch complex-39a"},
				generalSpaceXKeywords,
			),
			to: compose([]string{"environmental assessment", "launch tower"}),
		},

		// Launch tower
		{
			from: compose([]string{"mechazilla", "olit"}),
			to: compose(
				placesKeywords, sitesKeywords, nonSpecificKeywords,
				[]string{"qd", "sqd", "bqd", "arm", "catch", "lift"},
			),
		},
		{
			from: compose([]string{"tower"}),
			to: compose(
				placesKeywords, sitesKeywords, nonSpecificKeywords,
				[]string{"qd", "sqd", "bqd", "arm", "catch"},
			),
		},
		{
			from: compose([]string{"lift arms"}),
			to:   compose(nonSpecificKeywords, generalSpaceXKeywords),
		},

		// Load spreader
		{
			from: compose(ignoreSpaces([]string{"load spreader"})),
			to: compose(
				placesKeywords, sitesKeywords, liveStreams,
			),
		},

		// It looks like stuff is happening in Cape Canaveral at SpaceX Roberts Road. Not 100% though
		{
			from: compose(
				ignoreSpaces([]string{
					"roberts road", "robert's road", "robert road",
					"roberts rd", "robert's road", "robert rd",
				}),
			),
			to: compose(generalSpaceXKeywords, placesKeywords, sitesKeywords, liveStreams,
				ignoreSpaces([]string{"update", "olit", "launch tower"}),
			),
		},
		{
			from: compose(
				[]string{"cape", "canaveral"},
			),
			to: compose(placesKeywords, sitesKeywords, liveStreams,
				ignoreSpaces([]string{"update", "olit", "launch tower", "tower segment"}),
			),
			antiKeywords: ignoreSpaces([]string{"boat"}),
		},

		// Some words that are usually ambigious, but if combined with starship keywords they are fine
		{
			from: ignoreSpaces([]string{"launch tower", "launch pad", "launch mount", "chop stick", "chop stix", "catch arm"}),
			to:   compose(seaportKeywords, placesKeywords, sitesKeywords, liveStreams, nonSpecificKeywords), // Launch tower arm lift/load tests

		},
		{
			from: ignoreSpaces([]string{"launch mount", "chop stick", "chop stix", "catch arm", "can crusher"}),
			to:   compose([]string{"lift", "load"}, seaportKeywords, placesKeywords, sitesKeywords, liveStreams, nonSpecificKeywords),
		},

		{
			from: placesKeywords,
			to:   sitesKeywords,
		},
	}

	locationKeywords = map[string][]string{
		PascagoulaPlaceID: ignoreSpaces([]string{
			"sea launch", "phobos", "deimos",
		}),
	}

	// Helper slices that can be used for composing new keywords
	seaportKeywords       = ignoreSpaces([]string{"sea launch", "port", "oil", "rig"})
	placesKeywords        = ignoreSpaces([]string{"starbase", "boca chica"})
	sitesKeywords         = ignoreSpaces([]string{"launch site", "build site", "production site"})
	nonSpecificKeywords   = compose(ignoreSpaces([]string{"ship", "booster", "stage zero", "orbital test flight", "orbital flight test", "orbital flight"}), liveStreams, placesKeywords)
	generalSpaceXKeywords = ignoreSpaces([]string{"spacex", "space port", "elon", "musk", "gwynne", "shotwell"})
	testCampaignKeywords  = ignoreSpaces([]string{"static fire", "cryo", "detank", "can crusher", "test stand", "full stack"})
	liveStreams           = ignoreSpaces([]string{
		// 24/7 live camera views are often mentioned when something is shown on a screenshot
		"lab padre", "nasa space flight",
		"mc gregor live", "star base live",

		// Other streamers
		"jessica kirsh", "boca chica gal", "starship gazer",
	})

	// Regexes for road closures and testing activity
	closureTFRRegex = regexp.MustCompile(`\b(?:closure|tfr|notmar|cryo|fts)`)
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
			regexp.MustCompile(`(?:\s|^)(?:booster|super heavy|cryo|static fire|tower|ship|rud|faa|starbase|boca chica|lox|liquid oxygen|methane|ch4|relight|fts|cargo|lunar|tfr|scrub|flap|starship)\b`),
			// Try to match things for orbital flight tests
			regexp.MustCompile(`(?:orbit(?:.|\s)+(flight test|test flight)|(flight test|test flight)(?:.|\s)+orbit)`),
		},
		"spacex": {regexp.MustCompile(`(starship)`)},
	}

	userAntikeywordsOverwrite = map[string][]string{
		"elonmusk": {"tesla", "model s", "model 3", "model x", "model y", "car", "giga", "falcon", "boring company", "tunnel", "loop", "doge", "ula", "tonybruno", "jeff", "fsd", "giga berlin", "giga factory", "gigafactory", "giga press", "traffic", "alpha", "beta"},

		"spacex":  {},
		"faanews": {},

		// NASA Accounts that sometimes tweet about Starship don't need any antiKeywords - they are "allowed"
		// to mention Starship together with e.g. Orion (which would be ignored if not for these overrides).
		"nasa":          {"high bay"},
		"nasajpl":       {"high bay"},
		"nasa_marshall": {"high bay"},
		"nasa_gateway":  {"high bay"},
		"nasaartemis":   {"high bay"},
		"nasakennedy":   {"high bay", "starliner", "boeing"},
		"nasagoddard":   {"high bay"},
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
		"faanews":  true,
	}

	// If an account has any of these words in its description, we don't retweet tweets from it
	ignoredAccountDescriptionKeywords = ignoreSpaces([]string{
		// Parody accounts
		"parody", "joke",

		// 3D artists
		"artist", "blender", "3d", "vfx", "render", "animat", /* e/ion */

		// Sports stuff
		"nhl",

		"cum", "only fans",

		"crypto",
	})

	antiKeywordRegexes = []*regexp.Regexp{
		/* Falcon 9 booster numbers all start with 10 */
		regexp.MustCompile(`\b((?:b|booster)\s*10\d{2})\b`),

		/* Things that might look like a starship serial number, but aren't */
		regexp.MustCompile(`(^|\s)((?:s|ship)\s*\d{3,})\b`),

		// Things like "7yo", "8 yo", "15y/o"
		regexp.MustCompile(`\b(\d{1,2})\s*(?:yo|y/o|years?\s*old)\b`),
	}

	// If a tweet contains any of these keywords, it will not be retweeted. This is a way of filtering out *non-starship* stuff
	antiStarshipKeywords = []string{
		"electron", "blue origin", "neutron", "rocket lab", "rocketlab", "hungry hippo", "rklb", "falcon", "fairing half", "merlin", "m1d",
		"tesla ", "rivian", "giga press", "gigapress", "gigafactory", "openai", "boring", "hyperloop", "solarcity", "neuralink",
		"sls", "space launch system", "nasa_sls", "nasa_orion", "vehicle assembly building", "high bay 3", "vab", "ula", "united launch alliance", "vulcan", "rogozin",
		"virgingalactic", "virgin galactic", "virgin orbit", "virginorbit", "blueorigin", "boeing", "starliner", "soyuz", "soviet",
		"resilience", "shuttle", "challenger", "sts-51l", "sts-33", "new glenn", "china", "long march", "casc", "shenzhou", "india", "chinese", "japan", "space plane", "russia", "new shepard", "tsla", "dynetics",
		"ares", "titan", "ariane", "srb", "solid rocket booster", "terran", "relativity space", "relativityspace", "astra", "lv0",
		"spaceshipthree", "spaceshiptwo", "spaceshipone", "vss enterprise", "starship enterprise", "archer", "sisko", "vss imagine",
		"galaxy note", "galaxy s", "bezos", "jeff who", "branson", "tory", "bruno", "rp-1", "rp1", "biofuel", "bio fuel",
		"masten", "centaur", "atlas", "relativity", "northrop grumman", "northropgrumman", "bomber", "national team",
		"orbex", "rfa", "isar", "oneweb", "antares", "vega", "usaf b", "ms-", "starshipsls",
		"cygnus", "samsung", "s22 ultra", "angara", "firefly", "rolls-royce", "agrifood", "iot", "vs-50", "solid-propellant", "solid propellant",
		"são paulo", "sao paulo", "vlm-", "ac1", "arca", "ecorocket", "korea", "nuri", "mars rover", "perseverance", "curiosity", "ingenuity", "zhurong",
		"skoltech", "bmw",

		"launch umbilical tower", "mobile service structure", "appollo",

		"roscosmos", "yenisey",

		"twitter deal",

		"hubble", "nasahubble",

		// Blue Origins' Starship... kind of clone i guess?
		"jarvis", "glenn", "bob smith",

		"be4", "be-4", "be 4 engine",

		"war time", "wartime", "long range strike", "kyiv", "ukrain", "missile", "putin", "first strike",
		"call to arms", "calltoarms",

		"amazon", "kuiper", "isro",

		"yankees", "dodgers",

		"radian raptor",

		"7news",

		// e.g. crew-1, crew-2...
		"crew-", "crew dragon", "dragon", "crs", "dm-",

		"f22", "f-22", "jet", "b-52", "s-300", "f-1", "b52", "b350", "rs-25", "stennis",

		"b16 doubl",

		"seed round", "yc s", "not a starship",

		// Not interested in other stuff
		"doge", "babydoge", "coin", "btc", "fsd", "spce", "dogecoin", "crypto", "nft", "mint", "opensea",
		"safemoon", "stock", "wall street", "wallstreet", "buffett", "metaverse", "terra", "twtr", "board of director",
		"scam", "shill",

		"no tfr",

		// Usually mentions something like "Just started binge watching S30 of show xyz"
		"binge watch",

		"volvo",

		// "super heavyweight" in olympics...
		"super heavyweight",

		"parachute",

		"supernova",

		"kawai", "anthem", "katy perry",

		"god", "the lord", "pray",

		"firefight", "texaswildfire", "wildfire", "on fire", "engulfed in flames",

		"xanda",

		// 3d models are nice, but we only care about "real info"
		"render", "animat" /* ion/ed */, "3d", "model", "simulated", "print", "vfx", "not real", "photoshop",
		"art ", "artist", "mission patch", "drawing", "board game", "starshipshuffle", "starship shuffle", "lego",
		"card game", "starship design", "daily_hopper", "daily hopper", "paper model", "papermodel", "toy",

		"watercolor",

		"fantasy 4 cards", "fantasy cards",

		"8bitdo", "sn30", "ps-5", "ps5", "ssd", "sony",

		"gaming",

		"son", "daughter",

		"your guess",

		"not starship", "non starship", "not about starship", "discord", "wonder if", "was wondering", "years ago", "year ago",
		"nothing to do with starship", "not related to starship", "unrelated to starship",

		// kerbal space program, games, star wars != "official" news
		"kerbal space program", "ksp", "no mans sky", "nomanssky", "no man’s sky", "no man's sky", "kerbals", "pocket rocket", "pocketrocket", "simplerockets",
		"star trek", "startrek", "starcitizen", "star citizen", "battle droid", "b1-series", "civil war", "jabba the hutt", "sfs", "space flight simulator",
		"rocket explorer", "forza", "star wars", "starwars",

		"tax",

		// RC3 Seabee
		"seabee",

		// KSP planets, moons, stars etc.
		"moho", "gilly", "kerbin", "mun", "minmus", "duna", "jool", "laythe", "vall",
		"tylo", "bop", "dres", "eeloo", "kerbol",

		// Zodiac signs
		"aries", "taurus", "gemini", "cancer", "virgo", "libra", "scorpio", "sagittarius", "capricorn", "aquarius", "pisces",

		// movies
		"the martian", "starship trooper",

		// not *that* kind of raptor
		"velociraptor", "jurassic", "cretaceous", "dino",

		"ourmillion22",

		"kitty hawk",

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

		"ocisly", "jrti", "asog",

		"obetraveller", "ocean cam", "oceancam", "oceanscam", "paul",

		// I do not care about opinions on starship
		"agree", "disagree", "throwback to", "opinion", "imo", "imho", "i think", "mfw", "vibe", "dream", "laughs in",
		"armchair", "arm chair", "petition", "trust me bro",

		"gorgeous girl",

		"meme", "ratio", "apology", "drama", "petition to", "suck", "cursed", "uwu", "cult", "qwq", "reaction", "immigrant",

		"alien",

		"sale", "buy", "gift", "shop", "store", "purchase", "shirt", "sweater", "giveaway", "give away", "retweet", "birthday", "discount",
		"pre-order", "merch", "vote", "podcast", "trending", "hater", "unfollow", "top friends", "plush", "black friday", "blackfriday", "newprofilepic",
		"retweet if",

		"child", "kid", "parenthood",

		"illegal", "nfl", "nhl", "draw", "vax", "vacc", "booster shot", "shoot", "tik tok", "self harm", "sex", "cock", "s3x", "gspot", "g-spot", "fuck", "dick", "bullshit", "bikini",
		"booty", "cudd", "bathroom", "penis", "vagi", "furry", "stroking", "fap", "chick", "doggy", "only fans",

		"simp ", "simping",

		"belarus", "battalion",

		"stfu", "jerk", "thunderf00t", "thunderfoot", "common sense skeptic", "rambl",

		"patrons", "babylon", "boltup", "champion",

		"red bull",

		"tier list",

		// Annoying elon musk quotes
		"consciousness",

		// Some conferences have a "stand B20", because why not trick this bot right?
		"booth", "stand b",

		"trump", "antifa", "communism", "biden", "riot", "taliban", "kill", "beat", "ideology", "gender",

		"surgery", "emergency",

		"homopho", "hetero ", "cis ", "season",

		// Starts with "olm", which tricks the matcher
		"olmos",

		// Things that are typical questions for polls. We cannot get polls using the Twitter v1 API, so this is kind of bad
		"feel about",

		// stuff that seems like starship, but isn't
		"starshipent", "monstax", "eshygazit", "wonho",

		// Account follows a sheriff
		"arrest", "violence ", "assault", "rape", "weapon", "victim", "murder", "crime", "body", "nigg", "memorial", "dead", "death", "suicide", "piss", "wwii", "ww ii", "wwll", "ww ll",
		"abus", "gun", "injur",

		"nonce", "pedo",

		"bomb", "arsenal",

		"hospital", "midwife", "housewife", "baby face", "babyface",

		"offend", "offensive", "fanboy", "fan boy", "fangirl", "fan girl",

		"covid", "corona", "rona", "omicron", "tests positive", "positive test", "cdc",

		"diss", "shit", "anime", "manga", "bronco", "bae", "facist", "fascist",

		"abortion", "roe v. wade", "roe v wade", "roe vs. wade", "roe vs wade",

		"starshipcongrss", "starshipcongress", "congress", "starflight academy",

		// Sometimes the bot get confused because of "eiffel tower"
		"eiffel",

		// Make sure we don't retweet anything related to horrible tragedies
		"9/11", "911", "twin tower", "wtc", "trade center", "die", "falling",
	}

	moreSpecificAntiKeywords = []keywordMapping{
		{
			from: []string{"starlink"},
			to:   []string{"doug", "bob", "vessel", "fairing"},
		},
	}
)
