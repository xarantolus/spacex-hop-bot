package match

import "github.com/xarantolus/spacex-hop-bot/bot"

type StarshipMatcher struct {
	*Ignorer
}

func NewStarshipMatcher(ignoredUsers *Ignorer) *StarshipMatcher {
	return &StarshipMatcher{
		ignoredUsers,
	}
}

var TestIgnoredUserID int64 = 1983513

func NewStarshipMatcherForTests() *StarshipMatcher {
	return &StarshipMatcher{
		Ignorer: &Ignorer{
			list:     bot.ListMembersForTests(TestIgnoredUserID),
			keywords: ignoredAccountDescriptionKeywords,
		},
	}
}
