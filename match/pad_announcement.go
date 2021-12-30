package match

import "strings"

func IsPadAnnouncement(text string) bool {
	tl := strings.ToLower(text)

	return startsWithAny(tl, "pad") && (startsWithAny(tl, "announce") || startsWithAny(tl, "speaker"))
}
