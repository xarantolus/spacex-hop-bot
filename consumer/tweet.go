package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/xarantolus/spacex-hop-bot/match"
	"github.com/xarantolus/spacex-hop-bot/util"
	"mvdan.cc/xurls/v2"
)

// Process handles and retweets
type Processor struct {
	debug bool

	client *twitter.Client

	selfUser *twitter.User

	// map[URL]last Retweet time
	seenLinks map[string]time.Time

	seenTweets map[int64]bool

	spacePeopleListID      int64
	spacePeopleListMembers map[int64]bool
}

const (
	articlesFilename = "articles.json"
)

// NewProcessor returns a new processor with the given options
func NewProcessor(debug bool, client *twitter.Client, selfUser *twitter.User, spacePeopleListID int64) *Processor {
	p := &Processor{
		debug: debug,

		client:   client,
		selfUser: selfUser,

		seenLinks: make(map[string]time.Time),

		spacePeopleListID: spacePeopleListID,

		seenTweets:             make(map[int64]bool),
		spacePeopleListMembers: make(map[int64]bool),
	}

	util.LogError(util.LoadJSON(articlesFilename, &p.seenLinks), "loading links")

	return p
}

// Tweet processes the given tweet and checks whether it should be retweeted.
// Tweets that have already been seen are ignored. It is not safe for concurrent use.
func (p *Processor) Tweet(tweet match.TweetWrapper) {
	// So now we got a tweet. There are three categories that interest us:
	// 1. Elon Musk drops insider info about starship, e.g. as a reply.
	//    We do not care about his other tweets, so we check if any tweet
	//    in the reply chain matches stuff about starship
	// 2. We find a retweet
	// 3. We find a quoted tweet
	// 4. We find a tweet that is about starship

	if p.seenTweets[tweet.ID] || tweet.Retweeted {
		return
	}

	// Skip our own tweets
	if tweet.User != nil && tweet.User.ID == p.selfUser.ID {
		return
	}

	switch {
	case tweet.User != nil && tweet.User.ScreenName == "elonmusk":
		// When elon drops starship info, we want to retweet it.
		// We basically detect if the thread/tweet is about starship and
		// retweet everything that is appropriate
		p.thread(&tweet.Tweet)
	case tweet.QuotedStatus != nil:
		// If someone quotes a tweet, we only check the tweet that was quoted.
		p.Tweet(match.TweetWrapper{TweetSource: tweet.TweetSource, Tweet: *tweet.QuotedStatus})
	case tweet.RetweetedStatus != nil:
		p.Tweet(match.TweetWrapper{TweetSource: tweet.TweetSource, Tweet: *tweet.RetweetedStatus})
	case tweet.QuotedStatusID != 0 && !p.hasMedia(&tweet.Tweet):
		// Quoted tweets should be skipped, at least if they don't have media
	case match.StarshipTweet(tweet) && !p.isReply(&tweet.Tweet) && !p.isQuestion(&tweet.Tweet) && !p.isReactionGIF(&tweet.Tweet):
		// If the tweet itself is about starship, we retweet it
		// We already filtered out replies, which is important because we don't want to
		// retweet every question someone posts under an elon post, only those that
		// elon responded to.
		// Then we also filter out all tweets that tag elon musk, e.g. there could be someone
		// just tweeting something like "Do you think xyz... @elonmusk"

		// Filter out non-english tweets
		if tweet.Lang != "" && tweet.Lang != "en" && tweet.Lang != "und" {
			log.Println("Skipped", util.TweetURL(&tweet.Tweet), "because of language ", tweet.Lang)
			break
		}

		// Ignore links if there aren't any images
		if p.shouldIgnoreLink(&tweet.Tweet) && !p.hasMedia(&tweet.Tweet) {
			log.Println("Ignoring", util.TweetURL(&tweet.Tweet), "because of a link we ignore")

			break
		}

		if p.mentionsTooMany(&tweet.Tweet) {
			log.Println("Ignoring", util.TweetURL(&tweet.Tweet), "because it mentions too many people")
			break
		}

		// Depending on the tweet source, we require media
		if tweet.TweetSource == match.TweetSourceLocationStream {
			if p.hasMedia(&tweet.Tweet) {
				p.retweet(&tweet.Tweet, "normal + location media", tweet.TweetSource)
			} else {
				log.Println("[Twitter] Not retweeting", util.TweetURL(&tweet.Tweet), "because it's from the location stream and has no media")
			}
		} else {
			p.retweet(&tweet.Tweet, "normal matcher", tweet.TweetSource)
		}
	}

	p.seenTweets[tweet.ID] = true
}

func (p *Processor) mentionsTooMany(t *twitter.Tweet) bool {
	return strings.Count(t.Text(), "@") > 3
}

// isReply returns if the given tweet is a reply to another user
func (p *Processor) isReply(t *twitter.Tweet) bool {
	if t.QuotedStatusID != 0 {
		return true
	}

	if t.User == nil || t.InReplyToStatusID == 0 {
		return false
	}

	if t.User.ID != t.InReplyToUserID {
		return true
	}

	t, _, err := p.client.Statuses.Show(t.InReplyToStatusID, &twitter.StatusShowParams{
		TweetMode: "extended",
	})
	if err != nil {
		// If something goes wrong, we just assume it is a reply
		return true
	}

	return p.isReply(t)
}

func (p *Processor) isQuestion(tweet *twitter.Tweet) bool {
	return strings.Index(strings.ToLower(tweet.Text()), "@") < strings.Index(tweet.Text(), "?")
}

func (p *Processor) isReactionGIF(tweet *twitter.Tweet) bool {
	if tweet.ExtendedEntities == nil || len(tweet.ExtendedEntities.Media) != 1 {
		return false
	}

	// Type of a GIF is animated_gif
	return strings.Contains(tweet.ExtendedEntities.Media[0].Type, "gif")
}

func (p *Processor) hasMedia(tweet *twitter.Tweet) bool {
	return tweet.Entities != nil && len(tweet.Entities.Media) > 0 || tweet.ExtendedEntities != nil && len(tweet.ExtendedEntities.Media) > 0
}

// retweet retweets the given tweet, but if it fails it doesn't care
func (p *Processor) retweet(tweet *twitter.Tweet, reason string, source match.TweetSource) {
	// don't retweet anything in debug mode
	if p.debug {
		log.Printf("Not retweeting %s because we're in debug mode", util.TweetURL(tweet))
		return
	}

	// If we have already retweeted a tweet, we don't try to do it again, that just leads to errors
	if tweet.Retweeted || tweet.RetweetedStatus != nil && tweet.RetweetedStatus.Retweeted {
		return
	}

	_, _, err := p.client.Statuses.Retweet(tweet.ID, nil)
	if err != nil {
		util.LogError(err, "retweeting "+util.TweetURL(tweet))
		return
	}

	// save tweet so we can reproduce why it was matched
	p.saveTweet(tweet)

	// Add the user to our space people list
	// We ignore those from the location stream as they might not always tweet about starship
	if source != match.TweetSourceLocationStream {
		p.addSpaceMember(tweet)
	}

	twurl := util.TweetURL(tweet)
	log.Printf("[Twitter] Retweeted %s (%s - %s)", twurl, reason, source.String())

	// Setting Retweeted can help thread to detect that it should stop
	tweet.Retweeted = true
}

var (
	ignoredHosts = map[string]bool{
		"patreon.com":  true,
		"gofundme.com": true,
	}
	urlRegex *regexp.Regexp
)

func init() {
	var err error
	urlRegex, err = xurls.StrictMatchingScheme("https|http")
	if err != nil {
		panic("parsing URL regex: " + err.Error())
	}
}

// shouldIgnoreLink returns whether this tweet should be ignored because of a linked article
func (p *Processor) shouldIgnoreLink(tweet *twitter.Tweet) (ignore bool) {

	// Get the text *with* URLs
	var textWithURLs = tweet.SimpleText
	if textWithURLs == "" {
		textWithURLs = tweet.FullText
	}

	// Find all URLs
	urls := urlRegex.FindAllString(textWithURLs, -1)

	// Now check if any of these URLs is ignored
	for _, u := range urls {
		u = util.FindCanonicalURL(u)

		parsed, err := url.ParseRequestURI(u)
		if err != nil {
			log.Println("Cannot parse URL:", err.Error())
			continue
		}

		// Check if the host is ignored
		host := strings.ToLower(strings.TrimPrefix("www.", parsed.Hostname()))
		if ignoredHosts[host] {
			return true
		}

		// If we retweeted this link in the last 24 hours, we should
		// definitely ignore it
		lastRetweetTime, ok := p.seenLinks[u]
		if ok && time.Since(lastRetweetTime) < 24*time.Hour {
			return true
		}

		// Mark this link as seen, but allow a retweet
		p.seenLinks[u] = time.Now()

		// Now save it to make sure we still know after a restart
		util.LogError(util.SaveJSON(articlesFilename, p.seenLinks), "saving links")

		return false
	}

	return false
}

// saveTweet appends the given tweet to a JSON file for later inspections, especially in case of wrong retweets
func (p *Processor) saveTweet(tweet *twitter.Tweet) {
	f, err := os.OpenFile("retweeted.ndjson", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		util.LogError(err, "open tweet file")
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(tweet)
	if err != nil {
		util.LogError(err, "encoding tweet json")
	}
}

// addSpaceMember adds the user of the given tweet to the space people list
func (p *Processor) addSpaceMember(tweet *twitter.Tweet) {
	if tweet.User == nil || p.spacePeopleListMembers[tweet.User.ID] {
		return
	}

	// Idea: We make the list private, add the member and then make it public again.
	// That way they are not notified/annoyed
	defer p.client.Lists.Update(&twitter.ListsUpdateParams{
		ListID: p.spacePeopleListID,
		Mode:   "public",
	})
	// Set the list to private before updating
	p.client.Lists.Update(&twitter.ListsUpdateParams{
		ListID: p.spacePeopleListID,
		Mode:   "private",
	})

	p.spacePeopleListMembers[tweet.User.ID] = true

	_, err := p.client.Lists.MembersCreate(&twitter.ListsMembersCreateParams{
		ListID: p.spacePeopleListID,
		UserID: tweet.User.ID,
	})
	util.LogError(err, fmt.Sprintf("adding %s to list", tweet.User.ScreenName))
}

// thread processes tweet threads and retweets everything on-topic.
// This is useful because Elon Musk often replies to people that quote tweeted/asked a questions on his tweets
// See this for example: https://twitter.com/elonmusk/status/1372826575293583366
// or here: https://twitter.com/elonmusk/status/1372725108909957121
func (p *Processor) thread(tweet *twitter.Tweet) (didRetweet bool) {
	if tweet == nil {
		// Just in case
		return false
	}

	// Was that tweet interesting the last time we saw it?
	// If yes, then we should probably retweet the next stuff.
	// If not, we can stop here because it won't get any better
	// (we already checked the last time if it's good)
	if p.seenTweets[tweet.ID] || tweet.Retweeted {
		return tweet.Retweeted
	}
	p.seenTweets[tweet.ID] = true

	// First process the rest of the thread
	if tweet.InReplyToStatusID != 0 {
		// Ok, there was a reply. Check if we can do something with that
		parent, _, err := p.client.Statuses.Show(tweet.InReplyToStatusID, &twitter.StatusShowParams{
			IncludeEntities: twitter.Bool(true),
			TweetMode:       "extended",
		})
		util.LogError(err, "tweet reply status fetch (thread)")

		// If we have a matching tweet thread
		if parent != nil && p.thread(parent) {
			p.seenTweets[parent.ID] = true
			p.retweet(parent, "thread: matched parent", match.TweetSourceUnknown)
			didRetweet = true
		}
	}

	// A quoted tweet. Let's see if there's anything interesting
	if tweet.QuotedStatusID != 0 && tweet.QuotedStatus != nil {
		return p.thread(tweet.QuotedStatus)
	}

	realTweet := tweet
	if tweet.RetweetedStatus != nil {
		realTweet = tweet.RetweetedStatus
	}

	// Now actually match the tweet
	if didRetweet || match.StarshipTweet(match.TweetWrapper{TweetSource: match.TweetSourceUnknown, Tweet: *realTweet}) {

		p.retweet(tweet, "thread: matched", match.TweetSourceUnknown)

		return true
	}

	return didRetweet
}
