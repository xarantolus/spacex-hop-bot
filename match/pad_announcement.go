package match

import "strings"

func IsPadAnnouncement(text string) bool {
	tl := strings.ToLower(text)

	return startsWithAny(tl, "launchpad", "pad", "starbase", "boca chica", "bocachica", "build site") && (startsWithAny(tl, "announce", "speaker", "clear"))
}
