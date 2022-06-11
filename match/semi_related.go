package match

import "strings"

var starshipRelatedWhenElonReplies = compose(
	starshipKeywords,
	placesKeywords,
	nonSpecificKeywords,
	testCampaignKeywords,
)

func ElonReplyIsStarshipRelated(text string) bool {
	text = strings.ToLower(text)

	return startsWithAny(text, starshipRelatedWhenElonReplies...)
}
