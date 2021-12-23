package match

type StarshipMatcher struct {
	*Ignorer
}

func NewStarshipMatcher(ignoredUsers *Ignorer) *StarshipMatcher {
	return &StarshipMatcher{
		ignoredUsers,
	}
}
