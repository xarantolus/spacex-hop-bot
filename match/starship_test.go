package match

import (
	"strings"
	"testing"
)

func TestShipRegex(t *testing.T) {
	var shipMatch = starshipMatchers[0]

	var valid = []string{
		"sn10", "#sn10", "sn15", "sn 15", "starship s20",
		"starship number 15", "starship 15",
		"starship sn15s engines", "starship sn15's engines",
		"starship sn20?",
		"s300", "ship 20", "ship 20's nose", "ship 20’s nosecone section",
		"sn-11", "s-11",
	}

	var invalid = []string{"booster 10", "bn10", "b3496", "wordsn 10", "company's 20 cars", "company's 2021 report", "s3 dropping on netflix!",
		"u.s. to ship 4 mln covid-19 vaccine doses to nigeria, 5.66 mln to south africa"}

	for _, v := range valid {
		if shipMatch.FindString(v) == "" {
			t.Errorf("starshipMatchers[0] should have matched %q, but didn't", v)
		}
	}

	for _, i := range invalid {
		if shipMatch.FindString(i) != "" {
			t.Errorf("starshipMatchers[0] matched %q, but shouldn't have done that", i)
		}
	}
}

func TestBoosterRegex(t *testing.T) {
	var boostMatch = starshipMatchers[1]

	var valid = []string{"bn10", "bn1", "#b4", "bn 15", "booster b4",
		"booster number 15", "booster 15", "#bn4", "booster 15's engines",
		"booster number 15s engines", "booster 20’s", "booster 20's",
		"booster 3?", "bn-4", "b-4"}

	var invalid = []string{
		"starship 10", "b3496", "sn10", "wordbn 10",
		"company's 20 cars", "company's 2021 report",
		"booster 1049-11 arrives at the spacex dock",
		"eurocopter as.350-b2, is circling over cameron county",
	}

	for _, v := range valid {
		if boostMatch.FindString(v) == "" {
			t.Errorf("starshipMatchers[1] should have matched %q, but didn't", v)
		}
	}

	for _, i := range invalid {
		if boostMatch.FindString(i) != "" {
			t.Errorf("starshipMatchers[1] matched %q, but shouldn't have done that", i)
		}
	}
}

func TestGSERegex(t *testing.T) {
	var gseMatch = starshipMatchers[2]

	var valid = []string{"gse-5", "gse 5", "gse tank 5", "gse 5 tank", "gse tank"}

	var invalid = []string{"bn10", "bn1", "#b4", "bn 15", "booster b4",
		"booster number 15", "booster 15", "booster 15's engines",
		"booster number 15s engines", "booster 20’s", "booster 20's",
		"booster 3?", "starship 10", "b3496", "sn10", "wordbn 10", "company's 20 cars", "company's 2021 report"}

	for _, v := range valid {
		if gseMatch.FindString(v) == "" {
			t.Errorf("starshipMatchers[2] should have matched %q, but didn't", v)
		}
	}

	for _, i := range invalid {
		if gseMatch.FindString(i) != "" {
			t.Errorf("starshipMatchers[2] matched %q, but shouldn't have done that", i)
		}
	}
}

func TestRaptorRegex(t *testing.T) {
	var gseMatch = starshipMatchers[3]

	var valid = []string{"rvac 2", "rc 59", "raptor 2", "rb17", "rb9", "rc62",
		"raptor center 35", "raptor boost 35", "raptor vacuum 5", "raptor centre 35",
		"raptor engine boost 35", "raptor boost engine 35",
		"raptor engine vacuum 5", "raptor centre 35",
		"raptor vacuum 5", "raptor centre 35",
	}

	var invalid = []string{"bn10", "bn1", "#b4", "bn 15", "booster b4",
		"booster number 15", "booster 15", "booster 15's engines",
		"booster number 15s engines", "booster 20’s", "booster 20's",
		"booster 3?", "starship 10", "b3496", "sn10", "wordbn 10",
		"company's 20 cars", "company's 2021 report",
		"sn10", "#sn10", "sn15", "sn 15", "starship s20",
		"starship number 15", "starship 15",
		"starship sn15s engines", "starship sn15's engines",
		"starship sn20?",
		"raptor anyword 25",
		"s300", "ship 20", "ship 20's nose", "ship 20’s nosecone section",
		"sn-11", "s-11"}

	for _, v := range valid {
		if gseMatch.FindString(v) == "" {
			t.Errorf("starshipMatchers[3] should have matched %q, but didn't", v)
		}
	}

	for _, i := range invalid {
		if gseMatch.FindString(i) != "" {
			t.Errorf("starshipMatchers[3] matched %q, but shouldn't have done that", i)
		}
	}
}

func TestStarshipTextMatch(t *testing.T) {
	tests := []struct {
		text string
		want bool
	}{
		{"There’s more time to submit your comments on the Draft Programmatic Environmental Assessment for the proposed @SpaceX Starship/Super Heavy project in Boca Chica, Texas. Comment by Nov. 1. The new public meeting dates are Oct. 18th and 20th. Learn more at http://bit.ly/2YcScDe.", true},
		{"Rolls-Royce chosen by U.S. for new B-52 engines in contract worth up to $2.6 bln", false},
		{"Raptor 63 being lifted up to the booster", true},
		{"Bald Eagle in Canada flying over the water at the Canadian Raptor Conservancy by Fred Johns", false},
		{"No TFR posted for today", false},
		{"SN10", true},
		{"BN10", true},
		{"Starship SN10", true},
		{"SuperHeavy Booster", true},
		{"Unrelated doge coin tweet that also contains the keyword Starship", false},
		{"Unrelated tesla tweet", false},
		{"this tweet is not starship related", false},
		{"Starlink Mission", false},
		{`
SpaceX is targeting Wednesday, March 24 for launch of 60 Starlink satellites from Space Launch Complex 40 (SLC-40) at Cape Canaveral Space Force Station in Florida. The instantaneous window is at 4:28 a.m. EDT, or 8:28​ UTC, and a backup opportunity is available on Thursday, March 25 at 4:06 a.m. EDT, or 8:06​ UTC.

The Falcon 9 first stage rocket booster supporting this mission previously supported launch of the GPS-III Space Vehicle 03 and Turksat 5A missions in addition to three Starlink missions. Following stage separation, SpaceX will land Falcon 9’s first stage on the “Of Course I Still Love You” droneship, which will be located in the Atlantic Ocean. One half of Falcon 9’s fairing supported the Sentinel-6A mission and the other supported a previous Starlink mission.
`, false},
		{"I have received an Alert notice for tomorrow, July 19. Possible static fire attempt between noon and 10 p.m. on Booster B3.", true},
		{"Starship and Dogecoin", false},
		// Oil platform names need at least a bit of context
		{"Starship will land on Deimos", true},
		{"Deimos in the Ocean", false},
		{"SpaceX's Phobos launch platform", true},
		{"Phobos in the port", false},
		{"Samsung S22 Ultra", false},
		{"I mention Starship. $RKLB", false},
		{"Last week saw extensive work on Ship 20's TPS tiles, Booster 4 grew some engines, a GSE tank was tested, and some jets made an impressive flyover! Beyond Starbase, BO, China, and Astra all made launches, & Firefly Aerospace prepares their first flight!", false},
		{"Starship 20 and #Shenzhou12 ", false},
		{"Galaxy S22 Ultra", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := StarshipText(tt.text, antiStarshipKeywords); got != tt.want {
				t.Errorf("StarshipText(%q, antiStarshipKeywords) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}
func TestStreamTitles(t *testing.T) {
	tests := []struct {
		text string
		want bool
	}{
		{"Starlink Mission", false},
		{"Starship | SN11 | High-Altitude Flight Test", true},
		{"Starship | SN10 | High-Altitude Flight Recap", true},
		{"Starship | SN9 | High-Altitude Flight Test", true},
		{"Starship | SN8 | High-Altitude Flight Test", true},
		{"Starship SN20 & BN3: Orbital Flight Test", true},
		{"Starship | Starlink Mission", true},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := StarshipText(tt.text, nil); got != tt.want {
				t.Errorf("StarshipText(%q, nil) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestVariablesDuplicateKeywords(t *testing.T) {
	var words = make(map[string]bool)

	for _, k := range antiStarshipKeywords {
		_, ok := words[k]

		if ok {
			t.Errorf("Keyword %q is duplicated in antiStarshipKeywords slice", k)
		}

		words[k] = true
	}
}

func TestVariablesStringCase(t *testing.T) {
	for _, k := range starshipKeywords {
		if strings.ToLower(k) != k {
			t.Errorf("Keyword %q should be lowercase in starshipKeywords slice", k)
		}
	}
	for _, k := range antiStarshipKeywords {
		if strings.ToLower(k) != k {
			t.Errorf("Keyword %q should be lowercase in antiStarshipKeywords slice", k)
		}
	}
	for k, v := range userAntikeywordsOverwrite {
		if strings.ToLower(k) != k {
			t.Errorf("Account name %q should be lowercase in userAntikeywordsOverwrite map", k)
		}

		for _, s := range v {
			if strings.ToLower(s) != s {
				t.Errorf("Keyword %q should be lowercase in userAntikeywordsOverwrite slice for user %s", v, k)
			}
		}
	}

	for k := range specificUserMatchers {
		if strings.ToLower(k) != k {
			t.Errorf("Account name %q should be lowercase in specificUserMatchers map", k)
		}
	}
	for k := range hqMediaAccounts {
		if strings.ToLower(k) != k {
			t.Errorf("Account name %q should be lowercase in hqMediaAccounts map", k)
		}
	}

	for k := range veryImportantAccounts {
		if strings.ToLower(k) != k {
			t.Errorf("Account name %q should be lowercase in veryImportantAccounts map", k)
		}
	}
}

func Test_startsWithAnyGeneric(t *testing.T) {
	var searchedPrefixes = []string{"test", "best", "rest", "more than one word"}

	tests := []struct {
		argText string
		want    bool
	}{
		{"testing is nice", true},
		{"wrongprefixtesting is nice", false},
		{"the test keyword can be at any point in the string", true},
		{"we want to support more than one word", true},
		{"we want to support less than one word", false},

		{"the #test hashtag should still be recognized", true},
		{"also @test should work", true},
		{"#test at the beginning", true},
		{"@test should work", true},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := startsWithAny(tt.argText, searchedPrefixes...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_startsWithAnyStarship(t *testing.T) {
	tests := []struct {
		argText string
		want    bool
	}{
		{"KSP is my favourite game!", true},
		{"Project DogeCoin onto a Starship!", true},
		{"Starship reentering Kerbin's atmosphere", true},
		{"GSE Tank 6 rolling out", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := startsWithAny(strings.ToLower(tt.argText), antiStarshipKeywords...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_containsAnyGeneric(t *testing.T) {
	var searchedInfixes = []string{"test", "best", "rest", "more than one word"}

	tests := []struct {
		argText string
		want    bool
	}{
		{"testing is nice", true},
		{"wrongprefixtesting is nice", true},
		{"the test keyword can be at any point in the string", true},
		{"we want to support more than one word", true},
		{"we want to support less than one word", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := containsAny(tt.argText, searchedInfixes...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}
