package match

import "strings"

var padMappings = []keywordMapping{
	{
		from: ignoreSpaces([]string{"launchpad", "pad", "starbase", "boca chica", "bocachica", "build site"}),
		to:   ignoreSpaces([]string{"announce", "speaker", "clear"}),
	},
	{
		from: ignoreSpaces([]string{"announce", "speaker", "pa system", "pad"}),
		to:   ignoreSpaces([]string{"lift", "clear"}),
	},
}

func IsPadAnnouncement(text string) bool {
	tl := strings.ToLower(text)

	for _, m := range padMappings {
		if m.matches(tl) {
			return true
		}
	}
	return false
}
