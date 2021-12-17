package match

import (
	"strings"
	"testing"
)

func TestVariablesFirstIsAlphabet(t *testing.T) {
	// the startsWithAny function assumes that every antiKeyword starts with a letter between a and z
	for _, k := range antiStarshipKeywords {
		if !isAlphanumerical(rune(k[0])) {
			t.Errorf("Keyword %q in antiStarshipKeywords slice does not start with an alphanumerical character", k)
		}
	}
}

// since text in the StarshipText function is lowercase, we must make sure that all keywords are lowercase too
func TestVariablesStringCase(t *testing.T) {
	for _, k := range starshipKeywords {
		if strings.ToLower(k) != k {
			t.Errorf("Keyword %q should be lowercase in starshipKeywords slice", k)
		}
	}
	for i, kws := range moreSpecificKeywords {
		for _, k := range kws.from {
			if strings.ToLower(k) != k {
				t.Errorf("Keyword %q should be lowercase in moreSpecificKeywords[%d] 'from' mapping", k, i)
			}
		}

		for _, k := range kws.to {
			if strings.ToLower(k) != k {
				t.Errorf("Keyword %q should be lowercase in moreSpecificKeywords[%d] 'to' mapping", k, i)
			}
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

// Make sure we didn't forget specifying a "from" or "to" attribute
func TestMoreSpecificLength(t *testing.T) {
	for i, mapping := range moreSpecificKeywords {
		if len(mapping.from) == 0 {
			t.Errorf("moreSpecificKeywords[%d].from has length 0", i)
		}
		if len(mapping.to) == 0 {
			t.Errorf("moreSpecificKeywords[%d].to has length 0", i)
		}
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

// Tests for regexes

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
		"booster 3?", "booster's 4 and 5", "boosters 4 and 5"}

	var invalid = []string{
		"starship 10", "b3496", "sn10", "wordbn 10",

		// These would be nice, but there are many satellites and other stuff that is named B-somenumber, which makes it annoying
		"bn-4", "b-4",

		"company's 20 cars", "company's 2021 report",
		"booster 1049-11 arrives at the spacex dock",
		"somethingb3", "sb8",
		"eurocopter as.350-b2, is circling over cameron county",
		"f-35a completes final inert drop test of new b61-12 nuclear bomb",
		"https://example.com/somelinkthatincludesb3asboostername",
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
