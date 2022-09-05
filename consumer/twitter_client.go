package consumer

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

type TwitterClient interface {
	LoadStatus(tweetID int64) (*twitter.Tweet, error)

	AddListMember(listID int64, userID int64) (err error)

	Retweet(*twitter.Tweet) error
	UnRetweet(tweetID int64) error

	Tweet(text string, inReplyToID *int64) (*twitter.Tweet, error)
}

type NormalTwitterClient struct {
	Client *twitter.Client
	Debug  bool
}

func (n *NormalTwitterClient) UnRetweet(tweetID int64) error {
	_, _, err := n.Client.Statuses.Unretweet(tweetID, nil)
	return err
}

func (n *NormalTwitterClient) LoadStatus(tweetID int64) (tweet *twitter.Tweet, err error) {
	tweet, _, err = n.Client.Statuses.Show(tweetID, &twitter.StatusShowParams{
		IncludeEntities: twitter.Bool(true),
		TweetMode:       "extended",
	})

	return
}

func (n *NormalTwitterClient) AddListMember(listID int64, userID int64) (err error) {
	// Idea: We make the list private, add the member and then make it public again.
	// That way they are not notified/annoyed
	defer n.Client.Lists.Update(&twitter.ListsUpdateParams{
		ListID: listID,
		Mode:   "public",
	})
	// Set the list to private before updating
	n.Client.Lists.Update(&twitter.ListsUpdateParams{
		ListID: listID,
		Mode:   "private",
	})

	_, err = n.Client.Lists.MembersCreate(&twitter.ListsMembersCreateParams{
		ListID: listID,
		UserID: userID,
	})

	return
}

func (n *NormalTwitterClient) Retweet(tweet *twitter.Tweet) error {
	if n.Debug {
		return fmt.Errorf("not retweeting tweets in debug mode")
	}

	_, _, err := n.Client.Statuses.Retweet(tweet.ID, nil)

	return err
}

func (n *NormalTwitterClient) Tweet(text string, inReplyToID *int64) (t *twitter.Tweet, err error) {
	if n.Debug {
		err = fmt.Errorf("not tweeting in debug mode")
		return
	}

	var updateParams *twitter.StatusUpdateParams
	if inReplyToID != nil && *inReplyToID != 0 {
		updateParams = &twitter.StatusUpdateParams{
			InReplyToStatusID: *inReplyToID,
		}
	}

	t, _, err = n.Client.Statuses.Update(text, updateParams)
	return
}
