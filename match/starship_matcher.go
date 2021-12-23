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

func NewStarshipMatcherForTests() *StarshipMatcher {
	return &StarshipMatcher{
		Ignorer: &Ignorer{
			list:     bot.ListMembers(nil, "test"),
			keywords: nil,
		},
	}
}
