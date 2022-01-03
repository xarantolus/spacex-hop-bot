package bot

import (
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/util"
)

type UserList struct {
	listIDs []int64
	c       *twitter.Client

	purpose string

	mlock      sync.RWMutex
	members    map[int64]bool
	lastUpdate time.Time
}

func (l *UserList) ContainedIDs() (list []int64) {
	l.mlock.RLock()
	defer l.mlock.RUnlock()

	for mid := range l.members {
		list = append(list, mid)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	return list
}

func (l *UserList) ContainsByID(id int64) bool {
	l.update()

	l.mlock.RLock()
	defer l.mlock.RUnlock()

	return l.members[id]
}

// TweetAssociatedWithAny returns whether the given tweet is by, or mentions a user that is in this UserList
func (l *UserList) TweetAssociatedWithAny(tweet *twitter.Tweet) bool {
	if tweet == nil {
		return false
	}

	if tweet.User != nil && l.ContainsByID(tweet.User.ID) {
		return true
	}

	if tweet.QuotedStatus != nil && l.TweetAssociatedWithAny(tweet.QuotedStatus) {
		return true
	}

	if tweet.RetweetedStatus != nil && l.TweetAssociatedWithAny(tweet.RetweetedStatus) {
		return true
	}

	if tweet.Entities != nil {
		for _, m := range tweet.Entities.UserMentions {
			if l.ContainsByID(m.ID) {
				return true
			}
		}
	}

	if tweet.ExtendedTweet != nil && tweet.ExtendedTweet.Entities != nil {
		for _, m := range tweet.ExtendedTweet.Entities.UserMentions {
			if l.ContainsByID(m.ID) {
				return true
			}
		}
	}

	return false
}

func (l *UserList) update() {
	l.mlock.RLock()
	shouldUpdate := time.Since(l.lastUpdate) > 90*time.Minute
	l.mlock.RUnlock()

	if !shouldUpdate {
		return
	}

	l.mlock.Lock()
	defer l.mlock.Unlock()

	if len(l.listIDs) == 0 {
		return
	}

	l.members = make(map[int64]bool)

	for _, listID := range l.listIDs {
		list, _, err := l.c.Lists.Members(&twitter.ListsMembersParams{
			ListID: listID,
			Count:  1000,
		})
		if util.LogError(err, "loading list members for list "+strconv.FormatInt(listID, 10)) || list == nil {
			continue
		}

		for _, user := range list.Users {
			l.members[user.ID] = true
		}
	}

	log.Printf("[List] Updated list and loaded %d %s users\n", len(l.members), l.purpose)

	l.lastUpdate = time.Now()
}

// ListMembers loads a list of all users from the lists with the given ID
func ListMembers(c *twitter.Client, purpose string, listIDs ...int64) (membersMap *UserList) {
	membersMap = &UserList{
		c:       c,
		listIDs: listIDs,
		purpose: purpose,
	}

	membersMap.update()

	return
}

func ListMembersForTests(userIDs ...int64) *UserList {
	var membersMap = make(map[int64]bool)
	for _, mid := range userIDs {
		membersMap[mid] = true
	}
	return &UserList{
		members: membersMap,
		purpose: "test",
	}
}
