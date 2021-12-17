// This file defines basically everything the matcher does by specifying positive and negative keywords
package match

import "regexp"

type keywordMapping struct {
	from, to []string
}

// Note that all text here must be lowercase because the text is lowercased in the matching function
var (
	starshipKeywords = []string{
		"starship",
		"superheavy", "super heavy",
		"orbital launch tower", "orbital tower", "olt segment", "launch tower segment", "olp service tower", "olp tower",
		"orbital launch integration tower",
		"gse tank",
		"orbital launch table", "orbital table",
		"orbital launch pad", "orbital launch mount",
		"suborbital pad", "suborbital launch pad",
		"olp service tower",
		"orbital launch site",
		"launch tower arm", "mechazilla",
		"catch arms",
	}

	// moreSpecificKeywords are keywords that must be accompanied by at least one of the keywords mentioned in their slice.
	// This is useful for "raptor" (to make sure we only get engines) and some launch sites
	moreSpecificKeywords = []keywordMapping{
		// Engines
		{
			from: []string{"raptor"},
			to: []string{
				"starship", "vacuum", "sea-level", "sea level",
				"spacex", "mcgregor", "engine", "rb", "rc", "rvac",
				"launch site", "production site", "booster", "super heavy",
				"superheavy", "truck", "van", "raptorvan", "deliver", "sea level",
				"high bay", "nozzle", "tripod", "starbase",
			},
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
				[]string{"lc-49", "lc 49", "lauch complex 49", "lauch complex-49"},
				[]string{"lc-39a", "lc 39a", "lauch complex 39a", "lauch complex-39a"},
			),
			to: compose(starshipKeywords, generalSpaceXKeywords, []string{"ksc", "environmental assessment", "kennedy space center", "tower"}),
		},

		// Some words that are usually ambigious, but if combined with starship keywords they are fine
		{
			from: []string{"launch tower"},
			to:   starshipKeywords,
		},
	}
	// Helper slices that can be used for composing new keywords
	seaportKeywords       = []string{"sea launch", "oil", "rig"}
	generalSpaceXKeywords = []string{"spacex"}

	starshipMatchers = []*regexp.Regexp{
		// Starship SNx
		regexp.MustCompile(`\b((s-?\d{2,}\b)|(ship\s?\d{2,}\b)|(sn|starship|starship number)-?\s?\d+['’]?s?)`),
		// Booster BNx
		regexp.MustCompile(`(((?:#|\s|^)b\d{1,2}\b([^-]|$))|\b(bn|booster|booster number)(['’]|s)*\s?\d{1,3}['’]?s?\b)`),
		// Yes. I like watching tanks
		regexp.MustCompile(`\b(gse)\s?(?:tank|-)?\s?\d*\b`),
		// Raptor with a number
		regexp.MustCompile(`\b((?:raptor|raptor\s+engine|rvac|rb|rc)(?:\s+(?:center|centre|boost|vacuum))?(?:\s+engine)?\s*\d+)\b`),
	}

	closureTFRRegex = regexp.MustCompile("\b(?:closure|tfr|notmar|cryo|fts|scrub)")
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

		"starbasepulse": regexp.MustCompile(`(?:timelapse|time lapse)`),

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
		"electron", "blue origin", "neutron", "rocket lab", "rocketlab", "hungry hippo", "rklb", "falcon", "f9", "starlink",
		"tesla", "rivian", "giga press", "gigapress", "gigafactory", "openai", "boring", "hyperloop", "solarcity", "neuralink", "sls", "nasa_sls", "ula", "vulcan", "artemis", "rogozin",
		"virgingalactic", "virgin galactic", "virgin orbit", "virginorbit", "blueorigin", "boeing", "starliner", "soyuz", "soviet", "orion",
		"resilience", "shuttle", "challenger", "sts-51l", "sts-33", "new glenn", "china", "shenzhou", "india", "chinese", "japan", "space plane", "russia", "new shepard", "tsla", "dynetics", "hls",
		"ares", "titan", "ariane", "srb", "solid rocket booster", "terran", "relativity space", "relativityspace", "astra",
		"spaceshipthree", "spaceshiptwo", "spaceshipone", "vss enterprise", "starship enterprise", "archer", "sisko", "vss imagine", "galaxy note", "galaxy s", "bezos", "jeff who", "branson", "tory", "bruno",
		"masten", "centaur", "atlas v", "atlasv", "relativity", "northrop grumman", "northropgrumman", "bomber",
		"rookisaacman", "cygnus", "samsung", "angara", "firefly", "rolls-royce", "agrifood", "iot", "vs-50", "solid-propellant", "solid propellant",
		"são paulo", "sao paulo", "vlm-", "ac1", "arca", "ecorocket", "korea", "nuri",

		"roscosmos", "yenisey",

		// Blue Origins' Starship... kind of clone i guess?
		"jarvis", "glenn", "bob smith",

		"amazon", "kuiper", "nasaartemis", "isro",

		// e.g. crew-1, crew-2...
		"crew-", "crew dragon", "dragon", "crs", "dm-",

		"f22", "f-22", "jet", "b-52", "s-300",

		// Not interested in other stuff
		"doge", "coin", "btc", "fsd", "spce", "dogecoin", "crypto", "safemoon", "stock", "wall street", "wallstreet", "buffett",

		"no tfr",

		"volvo",

		// "super heavyweight" in olympics...
		"super heavyweight",

		"god", "the lord",

		// 3d models are nice, but we only care about "real info"
		"thanks", "thank you", "cheers", "render", "animat" /* ion/ed */, "3d", "model", "speculati" /*ng/on*/, "simulated", "print", "vfx", "not real", "photoshop",
		"art", "mission patch", "drawing", "board game", "starshipshuffle", "starship shuffle", "lego",
		"card game", "starship design", "daily_hopper", "daily hopper", "paper model", "papermodel",

		"8bitdo", "sn30",

		"your guess",

		"not starship", "non starship", "not about starship", "discord", "wonder if", "was wondering", "years ago", "year ago",

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

		"suprem", "aryan",

		"aircraft", "aerial refueling", "firepower",

		// "Star one brazilsat B4"
		"brasilsat", "star one",

		// Someone *really* named their delivery robot business "Starship"... why?
		"startup", "groceries", "robots", "starship robot",

		// And there's of course a wide range of products named S(numbers), which is annoying
		"tern", "gsd", "cargo bike",

		// Vitamin B2
		"vitamin",

		"ocisly", "jrti", "canaveral",

		"obetraveller", "ocean cam", "oceancam", "oceanscam", "paul",

		// I do not care about opinions on starship
		"agree", "disagree", "throwback to", "opinion", "imo", "imho", "i think", "mfw",

		"meme", "ratio", "apology", "drama", "petition to", "suck", "cursed", "uwu", "cult", "qwq", "reaction", "immigrant",

		"dearmoon", "dear moon", "inspiration4", "inspiration 4", "inspiration four", "rover", "alien",

		"sale", "buy", "shop", "store", "purchase", "shirt", "sweater", "giveaway", "give away", "retweet", "birthday", "download", "click", "tag", "discount",
		"pre-order", "merch", "vote", "podcast", "trending", "hater", "follow", "unfollow", "top friends", "plush", "black friday", "blackfriday", "newprofilepic",

		"child", "kid", "illegal", "nfl", "tiktok", "vax", "vacc", "shot", "shoot", "tik tok", "self harm", "sex", "cock", "s3x", "gspot", "g-spot", "fuck", "dick", "bullshit", "bikini",
		"booty", "cudd", "bathroom", "penis", "vagi", "furry",

		"patrons", "babylon", "boltup", "champion",

		"red bull", "browns",

		// Some conferences have a "stand B20", because why not trick this bot right?
		"booth", "stand b",

		"trump", "antifa", "biden", "riot", "taliban", "kill",

		// Things that are typical questions for polls. We cannot get polls using the Twitter v1 API, so this is kind of bad
		"feel about", "vs",

		// stuff that seems like starship, but isn't
		"starshipent", "monstax", "eshygazit", "wonho",

		// Account follows a sheriff
		"arrest", "violence ", "assault", "rape", "weapon", "victim", "murder", "crime", "investigat", "body", "memorial", "dead", "death", "cancer", "piss",

		"nonce", "pedo",

		"offend", "offensive", "fanboy", "fan boy", "fangirl", "fan girl",

		"covid", "corona",

		"shit", "anime", "manga", "bronco", "bae",

		"abortion", "roe v. wade", "roe v wade", "roe vs. wade", "roe vs wade",

		"starshipcongrss", "starshipcongress", "congress", "starflight academy",

		// Make sure we don't retweet anything related to horrible tragedies
		"9/11", "911", "twin tower", "wtc", "trade center", "die", "falling",
	}
)
