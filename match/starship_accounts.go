package match

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

func IsImportantAcount(u *twitter.User) bool {
	return u != nil && veryImportantAccounts[strings.ToLower(u.ScreenName)]
}
