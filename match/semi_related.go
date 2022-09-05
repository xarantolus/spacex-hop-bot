package match

import "strings"

var starshipRelatedWhenElonReplies = compose(
	starshipKeywords,
	placesKeywords,
	nonSpecificKeywords,
	testCampaignKeywords,
)

var notStarshipRelatedWhenElonReplies = compose()

func ElonReplyIsStarshipRelated(text string) bool {
	text = strings.ToLower(text)

	if _, notRelated := startsWithAny(text, notStarshipRelatedWhenElonReplies...); notRelated {
		return false
	}

	_, contains := startsWithAny(text, starshipRelatedWhenElonReplies...)
	return contains
}
