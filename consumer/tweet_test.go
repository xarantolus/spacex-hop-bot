package consumer

import (
	"testing"

	"github.com/xarantolus/spacex-hop-bot/match"
)

func TestBasicTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				text: "Tory announcing that Vulcan is heading to SLC-41. Potentially for a WDR (Wet Dress Rehearsal), at the very least fit checks.\n\nRemember, this vehicle actually has BE-4s, but not flight engines, thus a good while until a Static Fire test milestone.",
				acc:  "NASASpaceflight",
				want: false,
			},
			{
				text: "Good results from the Static Fire test for Falcon 9 B1060-8 ahead of Friday's launch.",
				acc:  "NASASpaceflight",
				want: false,
			},
			{
				text:     "Aborted Static Fire test, but no depress yet, so could be recycling.\n\n➡️https://youtube.com/watch?v=GP18t7ivstY",
				acc:      "NASASpaceflight",
				hasMedia: true,
				want:     true,
			},
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

			{
				acc:  "cnunezimages",
				text: "Hopper keeping watch 👀🔥🚀😎🤙",
				want: false,
				parent: &ttest{
					acc:      "cnunezimages",
					hasMedia: true,
					text:     "- Image Taken: December 29, 2021 - @elonmusk @spacex #Starbase #BocaChicaToMars #iCANimagine http://cnunezimages.com @SpaceIntellige3",
					want:     true,
				},
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
			// Someone asking a question below an elon tweet and getting an answer
			{
				text: "True, although it will look clean with close out panels installed. \n\nRaptor 2 has significant improvements in every way, but a complete design overhaul is necessary for the engine that can actually make life multiplanetary. It won’t be called Raptor.",
				acc:  "elonmusk",
				want: true,

				parent: &ttest{
					text: "Can't wait for Raptor 2, it's still a rat's nest up there.",
					want: true,
					parent: &ttest{
						text: "Random elon tweet",
						acc:  "elonmusk",
						want: false,
					},
				},
			},
			// Elon randomly answering tweets
			{
				acc:  "elonmusk",
				text: "All Raptor 2 tests going forward",
				want: true,
				parent: &ttest{
					text: "@SpaceX Raptor engine test last night in McGregor, Texas. The Raptor engine was tested on a horizontal test stand. #SpaceXtest \nFull Video: http://youtu.be/dCiEhBxTn7s",
					acc:  "photographer",
					want: true,
				},
			},
			// Top-level tweet
			{
				acc:  "elonmusk",
				text: "Tesla and Starship engines are currently the two hardest problems.",
				want: true,
			},

			{
				acc:  "elonmusk",
				text: "Each Raptor 1 engine above produces 185 metric tons of force. Raptor 2 just started production & will do 230+ tons or over half a million pounds of force.",
				want: true,
				parent: &ttest{
					acc:      "elonmusk",
					text:     "Starship Super Heavy engine steering test",
					hasMedia: true,
					want:     true,
				},
			},

			// Longer thread with questions
			{
				text: "Still aiming for booster 4 & Ship 20 for first orbital test flight (this is pure coincidence!)",
				acc:  "elonmusk",
				want: true,
				parent: &ttest{
					text: "Very interesting news about the upgrade to Ship's capability!\nWhich Booster+Ship combination are you aiming to fly the first orbital test with? Still Booster 4 and Ship 20, or use them only for ground testing?",
					acc:  "NASASpaceflight",
					want: true,
					parent: &ttest{
						acc:  "elonmusk",
						text: "Yup. Next booster will have 33 Raptor 2 engines, with 13 steering. \n\nShip is being upgraded to 9 engines (3 sea-level gimbaling, 6 vacuum fixed) with increased propellant load.",
						want: true,

						parent: &ttest{
							acc:  "NASASpaceflight",
							text: "Some sweet TVC (Thrust Vector Control) gimbal action from the Center 9 Raptor gang on the Booster.\n\nAnd that, ladies and gentlemen, is how the Booster steers.",
							want: true,
						},
					},
				},
			},
		},
	)
}

func TestTweetThreads(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				text: "The second static fire attempt of the day was aborted. Road has reopened. Another road closure is scheduled from 10 am to 6 pm central on Thursday if SpaceX wants to try again.",
				acc:  "nextspaceflight",
				want: true,
				parent: &ttest{
					acc:  "nextspaceflight",
					text: "LIVE: It appears that Ship 20 is going to make another attempt at a static fire\n\nhttps://youtu.be/GP18t7ivstY",
					want: true,
				},
			},

			{
				// Non-matching reply (due to no keywords)
				text:     "This is a reaction reply to the above tweet",
				hasMedia: true,
				want:     false,
				parent: &ttest{
					acc:  "nextspaceflight",
					text: "LIVE: It appears that Ship 20 is going to make another attempt at a static fire\n\nhttps://youtu.be/GP18t7ivstY",
					want: true,
				},
			},

			{
				// Reply that might match if it was not a reply.
				text:     "Ship 20 sure is beautiful today",
				acc:      "random_user",
				hasMedia: true,
				want:     false,
				parent: &ttest{
					acc:  "other_user",
					text: "Starship picture",
					want: true,
				},
			},

			{
				// Just a thread with one tweet with an description, then two images with non-matching description
				acc:      "NASASpaceflight",
				hasMedia: true,
				want:     true,
				text:     "From the beach",

				parent: &ttest{
					acc:      "NASASpaceflight",
					text:     "From the road",
					hasMedia: true,
					want:     true,

					parent: &ttest{
						text:     "Ship 20's on the test stand",
						acc:      "NASASpaceflight",
						hasMedia: true,

						want: true,
					},
				},
			},
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
			{
				acc:      "Random_Stranger",
				text:     "Here is an unrelated pic",
				hasMedia: true,
				want:     false,

				parent: &ttest{
					text:     "Great pace for Ship 20's test. Prop loading and a frost ring already. Great view from Mary (@BocaChicaGal)",
					acc:      "NASASpaceflight",
					hasMedia: true,

					want: true,
				},
			},
			{
				text:     "Methane Tank Fill Time!",
				acc:      "NASASpaceflight",
				hasMedia: true,
				want:     true,

				parent: &ttest{
					text:     "Well, this is as frosty as Booster 4's ever been.\nWe've moved into commentary mode on SBL, as the questions in chat are flying in. \nhttp://nasaspaceflight.com/starbaselive",
					acc:      "NASASpaceflight",
					hasMedia: true,
					want:     true,
				},
			},
			{
				text:     "Ship 20 just chilling:",
				hasMedia: true,
				acc:      "NASASpaceflight",
				want:     true,
				parent: &ttest{
					acc:         "NASASpaceflight",
					text:        "Prop loading for Ship 20 Static Fire Test 2! ➡️https://youtube.com/watch?v=GP18t7ivstY",
					want:        true,
					tweetSource: match.TweetSourceKnownList,
				},
			},
			{
				text: "The 2 LOX at the OTF 👇\n📈160th LOX delivery at the OTF\n(2/4) - Dec 29, 2021",
				acc:  "sb_deliveries",
				want: true,
				parent: &ttest{
					acc:  "sb_deliveries",
					text: "⛽ A lot of deliveries despite today’s long closure surprisingly!\n- 2 LOX to the Orbital Tank Farm\n- 4 LN2 to the Orbital Tank Farm\n- 2 LN2 to the Suborbital Tank Farm(1/4) - Dec 29, 2021 ",
					want: true,
				},
			},
		},
	)
}
