package match

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

type slice struct {
	name  string
	slice []string
}

// getAllSlices returns a list of all keyword slices (with their names) defined in the keyword file
func getAllSlices() (res []slice) {
	res = append(res,
		[]slice{
			{
				name:  "antiStarshipKeywords",
				slice: antiStarshipKeywords,
			},
			{
				name:  "starshipKeywords",
				slice: starshipKeywords,
			},
			{
				name:  "ignoredAccountDescriptionKeywords",
				slice: ignoredAccountDescriptionKeywords,
			},
		}...,
	)

	for i, v := range moreSpecificKeywords {
		res = append(res, slice{
			name:  fmt.Sprintf("moreSpecificKeywords[%d].from", i),
			slice: v.from,
		})
		res = append(res, slice{
			name:  fmt.Sprintf("moreSpecificKeywords[%d].to", i),
			slice: v.to,
		})
		res = append(res, slice{
			name:  fmt.Sprintf("moreSpecificKeywords[%d].antiKeywords", i),
			slice: v.antiKeywords,
		})
	}

	for u, v := range userAntikeywordsOverwrite {
		res = append(res, slice{
			name:  fmt.Sprintf("userAntikeywordsOverwrite[%q]", u),
			slice: v,
		})
	}

	return
}

func TestVariablesDuplicateKeywords(t *testing.T) {
	for _, v := range getAllSlices() {
		var words = make(map[string]bool)

		for _, k := range v.slice {
			_, ok := words[k]

			if ok {
				t.Errorf("Keyword %q is duplicated in %s slice", k, v.name)
			}

			words[k] = true
		}
	}
}

// The first character of any keyword matched with startsWithAny *must* be alphanumerical, as
// startsWithAny considers that the start of a word. So if they keyword was "#test" and the text is "#test",
// startsWithAny would only check if "test" starts with "#test", which of course doesn't work
func TestVariablesFirstIsAlphabet(t *testing.T) {
	for _, v := range getAllSlices() {
		for _, k := range v.slice {
			if !isAlphanumerical(rune(k[0])) {
				t.Errorf("Keyword %q in %s slice does not start with an alphanumerical character", k, v.name)
			}
		}
	}
}

// since text in the StarshipText function is lowercase, we must make sure that all keywords are lowercase too
func TestVariablesStringCase(t *testing.T) {
	for _, v := range getAllSlices() {
		for _, k := range v.slice {
			if strings.ToLower(k) != k {
				t.Errorf("Keyword %q should be lowercase in %s slice", k, v.name)
			}
		}
	}

	for k := range userAntikeywordsOverwrite {
		if strings.ToLower(k) != k {
			t.Errorf("Account name %q should be lowercase in userAntikeywordsOverwrite map", k)
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
	var testLen = func(mappingList []keywordMapping, name string) {
		for i, mapping := range mappingList {
			if len(mapping.from) == 0 {
				t.Errorf("%s[%d].from must not have length 0", name, i)
			}
			if len(mapping.to) == 0 {
				t.Errorf("%s[%d].to must not have length 0", name, i)
			}
		}
	}

	testLen(moreSpecificKeywords, "moreSpecificKeywords")
	testLen(moreSpecificAntiKeywords, "moreSpecificAntiKeywords")
}

func TestMoreSpecificMistakes(t *testing.T) {
	var testMappings = func(mapping []keywordMapping, name string) {
		for i, mapping := range mapping {
			if containsAll(mapping.to, starshipKeywords) {
				t.Errorf("%s[%d].to is composed with starshipKeywords, but that doesn't work and should be removed", name, i)
			}
			if containsAll(mapping.from, starshipKeywords) {
				t.Errorf("%s[%d].from is composed with starshipKeywords, but that doesn't work and should be removed", name, i)
			}
			if containsAll(mapping.from, antiStarshipKeywords) {
				t.Errorf("%s[%d].from is composed with antiStarshipKeywords, but that doesn't work and should be removed", name, i)
			}
			if containsAll(mapping.from, antiStarshipKeywords) {
				t.Errorf("%s[%d].from is composed with antiStarshipKeywords, but that doesn't work and should be removed", name, i)
			}
		}
	}
	testMappings(moreSpecificKeywords, "moreSpecificKeywords")
	testMappings(moreSpecificAntiKeywords, "moreSpecificAntiKeywords")
}

// Tests for regexes

func helpTestRegex(t *testing.T, regex *regexp.Regexp, regexName string, valid, invalid []string) {
	t.Helper()

	for _, v := range valid {
		if regex.FindString(v) == "" {
			t.Errorf("%s should have matched %q, but didn't", regexName, v)
		}
	}

	for _, i := range invalid {
		if regex.FindString(i) != "" {
			t.Errorf("%s matched %q, but shouldn't have done that", regexName, i)
		}
	}
}

func TestShipRegex(t *testing.T) {
	var valid = []string{
		"sn10", "#sn10", "sn15", "sn 15", "starship s20",
		"starship number 15", "starship 15",
		"starship sn15s engines", "starship sn15's engines",
		"starship sn20?",
		"ship 20", "ship 20's nose", "ship 20’s nosecone section",
		"sn-11",
	}

	var invalid = []string{"booster 10", "bn10", "b3496", "wordsn 10", "company's 20 cars", "company's 2021 report", "s3 dropping on netflix!",
		"u.s. to ship 4 mln covid-19 vaccine doses to nigeria, 5.66 mln to south africa", "s-11", "s70414937", "s300"}

	helpTestRegex(t, starshipMatchers[0], "starshipMatchers[0]", valid, invalid)
}

func TestBoosterRegex(t *testing.T) {
	var valid = []string{"bn10", "bn1", "#b4", "bn 15", "booster b4",
		"booster number 15", "booster 15", "#bn4", "booster 15's engines",
		"booster number 15s engines", "booster 20’s", "booster 20's",
		"booster 3?", "booster's 4 and 5", "boosters 4 and 5"}

	var invalid = []string{
		"starship 10", "b3496", "sn10", "wordbn 10",

		// These would be nice, but there are many satellites and other stuff that is named B-somenumber, which makes it annoying
		"bn-4", "b-4",

		// Falcon 9 booster name
		"b1051", "b1072",

		"company's 20 cars", "company's 2021 report",
		"booster 1049-11 arrives at the spacex dock",
		"somethingb3", "sb8", "web4",
		"eurocopter as.350-b2, is circling over cameron county",
		"f-35a completes final inert drop test of new b61-12 nuclear bomb",
		"https://example.com/somelinkthatincludesb3asboostername",
	}

	helpTestRegex(t, starshipMatchers[1], "starshipMatchers[1]", valid, invalid)
}

func TestGSERegex(t *testing.T) {
	var valid = []string{"gse-5", "gse 5", "gse 3", "gse tank 5", "gse 5 tank"}

	var invalid = []string{"bn10", "bn1", "#b4", "bn 15", "booster b4",
		"booster number 15", "booster 15", "booster 15's engines",
		"booster number 15s engines", "booster 20’s", "booster 20's",
		"booster 3?", "starship 10", "b3496", "sn10", "wordbn 10", "company's 20 cars", "company's 2021 report",
		"gse tank"}

	helpTestRegex(t, starshipMatchers[2], "starshipMatchers[2]", valid, invalid)
}

func TestRaptorRegex(t *testing.T) {
	var valid = []string{"rvac 2", "rc 59", "raptor 2", "rb17", "rb9", "rc62",
		"raptor center 35", "raptor boost 35", "raptor vacuum 5", "raptor centre 35",
		"raptor engine boost 35", "raptor boost engine 35",
		"raptor engine vacuum 5", "raptor centre 35",
		"raptor vacuum 5", "raptor centre 35",
		"raptor v1", "raptor v2",
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

	helpTestRegex(t, starshipMatchers[3], "starshipMatchers[3]", valid, invalid)
}

func TestAlertRegex(t *testing.T) {
	var valid = []string{
		"have received an alert notice",
		"static fire will be attempted later today",
		"cryo proof upcoming",
		"Spacex is clearing the pad",
		"pad cleared",
	}

	var invalid = []string{"booster 10", "bn10", "b3496", "wordsn 10", "company's 20 cars", "company's 2021 report", "s3 dropping on netflix!",
		"u.s. to ship 4 mln covid-19 vaccine doses to nigeria, 5.66 mln to south africa", ""}

	helpTestRegex(t, alertRegex, "alertRegex", valid, invalid)
}

func TestClosureTFRRegex(t *testing.T) {
	var valid = []string{
		"fts is installed",
		"new notmar posted",
		"cryo proof upcoming",
	}

	var invalid = []string{"booster 10", "bn10", "b3496", "wordsn 10", "company's 20 cars", "company's 2021 report", "s3 dropping on netflix!",
		"u.s. to ship 4 mln covid-19 vaccine doses to nigeria, 5.66 mln to south africa", "",
		"have received an alert notice",
		"static fire will be attempted later today",
	}

	helpTestRegex(t, closureTFRRegex, "closureTFRRegex", valid, invalid)
}

func TestFalcon9BoosterRegex(t *testing.T) {
	var valid = []string{
		"booster 1021",
		"b1072",
		"booster b1021",
		"booster 1050",
	}

	var invalid = []string{"booster 10", "bn10", "b3496", "wordsn 10", "company's 20 cars", "company's 2021 report", "s3 dropping on netflix!",
		"", "b4", "notbooster 1050", "n1025",
	}

	helpTestRegex(t, antiKeywordRegexes[0], "antiKeywordRegexes[0]", valid, invalid)
}

func TestNonStarshipRegex(t *testing.T) {
	var valid = []string{"s00012 vanguard", "s70414937", "s203198", "this is s123", "Someone will ship 500 tons of cargo to some place"}

	var invalid = []string{"s10", "s 10", "sn10", "sn 10", "bn10", "b3496", "wordsn 10", "company's 20 cars", "company's 2021 report", "s3 dropping on netflix!",
		"", "b4", "notbooster 1050", "n1025", "booster 1021", "b1072", "booster b1021", "booster 1050", "atlas5",
		"last nights 110 sec Raptor 2 Tripod test.   Evening Tripod tests are the best!",
	}

	helpTestRegex(t, antiKeywordRegexes[1], "antiKeywordRegexes[1]", valid, invalid)
}

func TestAgeRegex(t *testing.T) {
	var valid = []string{"5yo", "5y/o", "5 y/o", "5 year old"}

	var invalid = []string{"yo", " yo"}

	helpTestRegex(t, antiKeywordRegexes[2], "antiKeywordRegexes[2]", valid, invalid)
}

// containsAll returns if subset is a subset of set
func containsAll(subset, set []string) bool {
	var asmap = map[string]bool{}
	for _, s := range subset {
		asmap[s] = true
	}

	for _, s := range set {
		if asmap[s] == false {
			return false
		}
	}

	return true
}
