package consumer

import (
	"testing"

	"github.com/xarantolus/spacex-hop-bot/match"
)

func TestBasicTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				// this is the test user ID; we don't want to retweet our own tweets
				userID: 5,
				text:   "S20 standing on the pad",
				want:   false,
			},
			{
				text: "S20 standing on the pad",
				want: true,
			},
			{
				text: "Unrelated",
				want: false,
			},
			{
				text: "Road closure with no information where it is",
				want: false,
			},

			{
				text: "Road closure with no information where it is, but trusted account",
				acc:  "nextspaceflight",
				want: true,
			},

			// If we have a tweet that only contains (hash)tags, it should only retweeted if it has media
			{
				text:     "#Starbase #Starbase #SpaceX #Starship @elonmusk",
				hasMedia: true,
				want:     true,
			},
			{
				text: "#Starbase #Starbase #SpaceX #Starship @elonmusk",
				want: false,
			},
		},
	)
}

func TestLocationTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			// If it explicitly mentions a starship, then no need for location
			{
				text: "Pad announcement over the speakers: clearing pad for S20 static fire",
				want: true,
			},
			// Here we have the same tweet, but one with a good location
			{
				text: "Pad announcement over the speakers: clearing pad for static fire",
				want: false,
			},
			{
				text:     "Pad announcement over the speakers: clearing pad for static fire",
				location: "random place",
				want:     false,
			},
			{
				text:     "Pad announcement over the speakers: clearing pad for static fire",
				location: match.StarbasePlaceID,
				hasMedia: true,
				want:     true,
			},

			{
				text:        "Pad announcement over the speakers: clearing pad for static fire",
				want:        false,
				tweetSource: match.TweetSourceLocationStream,
			},
			{
				text:        "Pad announcement over the speakers: clearing pad for static fire",
				location:    "random place",
				tweetSource: match.TweetSourceLocationStream,
				want:        false,
			},
			{
				text:        "Pad announcement over the speakers: clearing pad for static fire",
				location:    match.StarbasePlaceID,
				tweetSource: match.TweetSourceLocationStream,
				hasMedia:    true,
				want:        true,
			},

			// However, we don't want *any* tweet from starbase etc.
			{
				text:        "Drinking some coffee at the beach",
				location:    match.StarbasePlaceID,
				tweetSource: match.TweetSourceLocationStream,
				want:        false,
			},
		},
	)
}

func TestQuestionTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			// Questions only if we have media or are at the spacex locations
			{
				acc:  "considercosmos",
				text: "Super Heavy is now hooked up to @SpaceX crane...\nWill we see a booster 4 lift soon?",
				want: false,
			},
			{
				acc:      "considercosmos",
				text:     "Super Heavy is now hooked up to @SpaceX crane...\nWill we see a booster 4 lift soon?",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Super Heavy is now hooked up to @SpaceX crane...\nWill we see a booster 4 lift soon?",
				location: match.SpaceXBuildSiteID,
				want:     true,
			},
		},
	)
}

func TestElonTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				acc:  "elonmusk",
				text: "Wow, working on this problem has soaked up a lot of my time & brain cycles over the past ~7 years! This and Starship engines are currently the two hardest problems.",
				want: true,
			},
		},
	)
}

func TestTweetThreads(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				acc:      "NASASpaceflight",
				text:     "Standing by for siren!",
				hasMedia: true,
				want:     true,

				parent: &ttest{
					text:     "Great pace for Ship 20's test. Prop loading and a frost ring already. Great view from Mary (@BocaChicaGal)",
					acc:      "NASASpaceflight",
					hasMedia: true,

					want: true,
				},
			},
		},
	)
}
