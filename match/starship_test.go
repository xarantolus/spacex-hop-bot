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
		"s300",
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
		"booster number 15", "booster 15", "booster 15's engines",
		"booster number 15s engines", "booster 20’s", "booster 20's",
		"booster 3?"}

	var invalid = []string{"starship 10", "b3496", "sn10", "wordbn 10", "company's 20 cars", "company's 2021 report"}

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

func TestStarshipAntiKeywords(t *testing.T) {
	tests := []struct {
		text string
		want bool
	}{
		{"No TFR posted for today", false},
		{"SN10", true},
		{"BN10", true},
		{"Starship SN10", true},
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
func TestVariables(t *testing.T) {
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
}
