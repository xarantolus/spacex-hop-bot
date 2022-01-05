package consumer

import (
	"testing"
	"time"

	"github.com/xarantolus/spacex-hop-bot/match"
)

func TestBasicTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				text: "Chopsticks\n@LabPadre",
				want: true,
			},
			{
				text: "CraneX lifting a Ship",
				want: true,
			},
			{
				text: "SpaceX crane went for a stroll to hook up with Booster 4.\n\nhttp://nasaspaceflight.com/starbaselive",
				acc:  "nextspaceflight",
				want: true,
			},
			{
				text: "Gorgeous gorgeous girls want a starship orbital test flight.",
				want: false,
			},
			{
				text: "Timelapse of the chopsticks slowly on the move in Starbase.\n\nhttp://nasaspaceflight.com/starbaselive",
				acc:  "nextspaceflight",
				want: true,
			},
			{
				text: "Orbital Launch Tower catching arms have begun their first vertical move visible on @LabPadre rover cam 1 & 2. @elonmusk @SpaceX",
				want: true,
			},
			{
				text: "Starship-Orion is a good idea",
				want: false,
			},
			{
				text: "In 2022, we are likely to see the debuts orbital flights of the two most powerful rockets ever: SpaceX’s Starship and Boeing/NASA’s SLS.\nEasily the most significant launch vehicles since the Saturn V of the Apollo era.\nHappy New Year! Many exciting days (& launches Rocket ) to come",
				want: false,
			},
			{
				text: "Talking to every goth rocker chick in the solar meatspace until I find one with a starship guidance chip still containing the coordinates for a disused dive bar named Pair-a-Dice that shares its orbit with the 13th planet where my family's DNA backup chip is under a floor tile.",
				acc:  "swiftonsecurity",
				want: false,
			},
			{
				text: "Following yesterday's focus on Ship 20, attention has switched to Booster 4 on the Orbital Launch Mount.\n\nThe SpaceX LR11000 crawler crane has been hooked up to the booster.\n\nMary (@BocaChicaGal) is already out there, with the enhanced view:\n\nhttp://nasaspaceflight.com/starbaselive",
				acc:  "nasaspaceflight",
				want: true,
			},

			{
				text: "Morning to Morning 24 hour timelapse (Dec 29 through Dec 30) of @NASASpaceflight's Starbase Live camera at https://nasaspaceflight.com/starbaselive (click to watch live)\n\nSN20 static fire day!\n\n#BocaChicaToMars #SpaceX #Starship",
				acc:  "StarbasePulse",

				want: true,
			},
			{
				acc:  "TheFavoritist",
				text: "Ship 20 static fired its Raptor engines again as SpaceX progresses toward the first orbital test flight. After one successful firing, Ship 20 aborted a second attempt.\nVideo from @BocaChicaGal and the NSF Robots. Edited by @Patrick_Colqu\n📺 https://youtu.be/_12ePNH0wTc",
				want: true,
			},
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
					text:     "- Image Taken: " + time.Now().Format("Monday, January 02, 2006") + " - @elonmusk @spacex #Starbase #BocaChicaToMars #iCANimagine http://cnunezimages.com @SpaceIntellige3",
					want:     true,
				},
			},
			{
				text: "Primary Date: Road Closure Scheduled for " + time.Now().Format("Monday, January 02, 2006") + " from 10:00 a.m. to 8:00 p.m.",
				acc:  "BocaRoad",
				want: true,
			},

			// Road closures
			{
				text: "Secondary Date: Road Closure Scheduled Extended for " + time.Now().Format("Monday, January 02, 2006") + " from 10:00 a.m. to 8:00 p.m.",
				acc:  "BocaRoad",
				want: true,
			},

			{
				text: "Booster 4 lifting soon. Can watch it LIVE on my YouTube stream from a unique angle filming with a professional camera:\nhttps://www.youtube.com/watch?v=yV48vHXNkNA",
				acc:  "starshipgazer",
				want: true,
			},

			{
				text: "Booster QD (Quick Disconnect) detached, retracted, and hood closed.\n\nVery cool to watch that in action.\n\nA bit of a timelapse via Mary (@BocaChicaGal)'s view: http://nasaspaceflight.com/starbaselive",
				acc:  "nasaspaceflight",
				want: true,
			},

			{
				text: "This great image from @NASAHubble shows the rich and diverse collection of galaxies in the cluster Abell S0740. The cluster is more than 450 million light-years away.",
				want: false,
			},
			{
				// This is about cars
				text: "*laughs in B5 S4 where everything maintenance or repair wise required the engine out or at least the entire front clip off*",
				want: false,
			},

			{
				text: "Humans for scale. \n\nOne human got to touch a Raptor nozzle #jealous \n\nhttp://nasaspaceflight.com/starbaselive",
				want: true,
			},
		},
	)
}

func TestLocationTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				text:     "Might we see Booster lift back off the OLM? CraneX is in place, load spreader attached, booster stand is on hand and ready. #Starbase #SpaceX",
				location: match.SpaceXLaunchSiteID,
				hasMedia: true,
				want:     true,
			},
			{
				text:        "Might we see Booster lift back off the OLM? CraneX is in place, load spreader attached, booster stand is on hand and ready. #Starbase #SpaceX",
				tweetSource: match.TweetSourceLocationStream,
				hasMedia:    true,
				want:        true,
			},
			{
				text:        "Might we see Booster lift back off the OLM? CraneX is in place, load spreader attached, booster stand is on hand and ready. #Starbase #SpaceX",
				location:    match.SpaceXLaunchSiteID,
				tweetSource: match.TweetSourceLocationStream,
				hasMedia:    true,
				want:        true,
			},

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
				// Announcement without media
				text:        "Pad announcement over the speakers: clearing pad for static fire",
				location:    match.StarbasePlaceID,
				tweetSource: match.TweetSourceLocationStream,
				want:        true,
			},
			{
				text:        "Pad announcement over the speakers: clearing pad for static fire",
				location:    match.StarbasePlaceID,
				tweetSource: match.TweetSourceLocationStream,
				hasMedia:    true,
				want:        true,
			},

			// Pad announcements with questions are allowed (if at the place)
			{
				text: "Just heard a pad announcement. Very hard to hear, sounded like some sort of pad operations. Could be some sort of testing?",
				want: false,
			},
			{
				text:     "Just heard a pad announcement. Very hard to hear, sounded like some sort of pad operations. Could be some sort of testing?",
				want:     true,
				location: match.SpaceXLaunchSiteID,
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
			// Elon answering to an ignored account (e.g. a 3d animation)
			{
				text: "And ship will be caught by Mechazilla too. As with booster, no landing legs. Those are only needed for moon & Mars until there is local infrastructure.",
				acc:  "elonmusk",
				want: true,

				parent: &ttest{
					text: "Pretty close. Booster & arms will move faster. QD arm will steady booster for ship mate.",
					acc:  "elonmusk",
					want: true,

					parent: &ttest{
						text:     "Mechazilla <1 Hour Turnaround.\n#SpaceX #Starship @elonmusk",
						acc:      "ErcXspace",
						hasMedia: true,
						userID:   match.TestIgnoredUserID,

						// This *should* probably be true, but would require a bigger rewrite of the thread logic
						want: false,
					},
				},
			},
			{
				text: "Pretty close. Booster & arms will move faster. QD arm will steady booster for ship mate.",
				acc:  "elonmusk",
				want: true,

				parent: &ttest{
					text:     "Mechazilla <1 Hour Turnaround.\n#SpaceX #Starship @elonmusk",
					acc:      "ErcXspace",
					hasMedia: true,
					userID:   match.TestIgnoredUserID,

					want: false,
				},
			},
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

			{
				text: "The Starship fleet is designed to achieve over 1000 times more payload to orbit than all other rockets on Earth combined.\n\nAlmost no one understands this.",
				acc:  "elonmusk",
				want: true,

				parent: &ttest{
					text: "True",
					acc:  "elonmusk",
					want: true,

					parent: &ttest{
						text: "What’s perhaps most crazy is a single Starship / SuperHeavy launch could put everything launched in that quarter into orbit in a single launch… now that's impressive.",
						acc:  "Erdayastronaut",
						want: true,

						parent: &ttest{
							text: "Actually, 41 tons for SpaceX in Q3 & aiming for 80 tons in Q4. That said, China launch mass to orbit is extremely impressive.",
							acc:  "elonmusk",

							want: false,

							parent: &ttest{
								text: "China led Q3 in both the number and payload mass of orbital rocket launches, according to the latest @BryceSpaceTech report.\n\nKilograms of mass launched:\n\nCASC 45,010\nSpaceX 32,634\nArianespace 25,881\nRoscosmos 20,500\nNorthrop Grumman 5,358\nULA 2,888\n\nhttps://brycetech.com/briefing",
								acc:  "thesheetztweetz",
								want: false,
							},
						},
					},
				},
			},

			{
				text: "Orbital flight test of the largest rocket ever soon!",
				acc:  "elonmusk",
				want: true,
			},
			{
				text: "Orbital test flight of the most capable rocket ever!",
				acc:  "elonmusk",
				want: true,
			},
			{
				text: "Launching the test flight from the cape should work.\n\nHope to reach orbit!",
				acc:  "elonmusk",
				want: true,
			},

			{
				text:     "129 Orbital Flights",
				acc:      "elonmusk",
				hasMedia: true,
				want:     false,
			},
			{
				text: "Construction of Starship orbital launch pad at the Cape has begun",
				acc:  "elonmusk",
				want: true,
			},

			{
				text: "Yes",
				acc:  "elonmusk",
				want: true,

				parent: &ttest{
					text: "Still at 39A?",
					acc:  "NASASpaceflight",
					want: true,

					parent: &ttest{
						text: "Construction of Starship orbital launch pad at the Cape has begun",
						acc:  "elonmusk",
						want: true,
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
				text: "It looks like they have also released all of the chains holding #Mechazilla back. Might see this monster flex its arms tonight if we get lucky @elonmusk #SpaceX",
				want: true,
				acc:  "same_user",
				parent: &ttest{
					text: "Meet the #LukeBeamWalkers of #Starbase, TX.  Here are a few of my favorite moments from @StarshipGazer's stream today. Removing this scaffolding is one of the last remaining items before #Mechazilla can start performing curls with #Booster4 and #Ship20\n\nhttps://youtube.com/watch?v=wAzC07",
					acc:  "same_user",
					want: true,
				},
			},
			{
				text:     "The flarestack for the ground-based Raptor engine stands (horizontal and vertical) was busy testing through numerous ignite-increase/decrease-extinguish-repeat sequences.",
				acc:      "bluemoondance74",
				hasMedia: true,

				want: true,

				parent: &ttest{
					text: "An open Merlin Vacuum engine test bay, part of the Multi-Merlin stand/Small Site, is prepared for static fire testing (left); and a closed bay (right).",
					acc:  "bluemoondance74",

					want: false,

					parent: &ttest{
						text: "Even during the holiday, testing at @SpaceX’s McGregor facility has continued.\nOn my last visit, test preps were being made at the Merlin Vacuum engine bay, along with flarestack testing at the ground-based Raptor stands.\n(Roars have been heard daily- w/ 2 mega rumbles today!)",
						acc:  "bluemoondance74",

						want: false,
					},
				},
			},

			{
				text: "Musk recently tweeted that only Raptor 2s are being delivered to McGregor and tested from now on :)",
				want: false,
				parent: &ttest{
					text: "How did you guess that? ",
					acc:  "other_user",
					want: false,

					parent: &ttest{
						text: "If Elon is to be believed, this is a Raptor 2 static fire :D",
						want: false,
						quoted: &ttest{
							text:     "Raptor engine roar 🔥🚀✨\n@NASASpaceflight #SpaceXTests",
							acc:      "bluemoondance74",
							hasMedia: true,
							want:     true,

							parent: &ttest{
								text:     "#McGregorTX",
								acc:      "bluemoondance74",
								hasMedia: true,
							},
						},
					},
				},
			},
			{
				text: "Actually elon confirmed the new engine configuration for Starship 29:",
				acc:  "random_user",
				want: false,
				quoted: &ttest{
					text: "Starship 29 will have more engines in the future",
					acc:  "elonmusk",
					want: true,
				},
				parent: &ttest{
					text: "There isn't any information on how many engines Starship 29 will have",
					acc:  "other_user",
					want: true,
				},
			},
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
				text: "The 2 LOX at the OTF 👇\n📈160th LOX delivery at the OTF\n(2/4) - " + time.Now().Format("January 02, 2006"),
				acc:  "sb_deliveries",
				want: true,
				parent: &ttest{
					acc:  "sb_deliveries",
					text: "⛽ A lot of deliveries despite today’s long closure surprisingly!\n- 2 LOX to the Orbital Tank Farm\n- 4 LN2 to the Orbital Tank Farm\n- 2 LN2 to the Suborbital Tank Farm(1/4) - " + time.Now().Format("January 02, 2006"),
					want: true,
				},
			},
			{
				acc:      "RGVaerialphotos",
				hasMedia: true,
				want:     true,
				parent: &ttest{
					text:     "Ship 21 Nose Cone",
					acc:      "RGVaerialphotos",
					want:     true,
					hasMedia: true,
				},
			},
		},
	)
}

func TestQuotedTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				text:     "The most used word in the Starbase World right now is \"Chopsticks\". \nThey moved a little last night (https://twitter.com/nextspaceflight/status/1477869030094479365…), so hopefully, we'll see some more action soon!\n\nMary (@BocaChicaGal) with the cool view:\nhttp://nasaspaceflight.com/starbaselive",
				acc:      "nasaspaceflight",
				hasMedia: true,

				want: true,

				quoted: &ttest{
					acc:      "nextspaceflight",
					text:     "Timelapse of the chopsticks slowly on the move in Starbase.",
					hasMedia: true,

					want: true,
				},
			},
			{
				// If someone quotes their own tweet with more media, we want to retweet it
				text: "Even more pics of Starship S20",
				acc:  "same_user",

				hasMedia: true,
				want:     true,

				quoted: &ttest{
					text: "Starship SN20 being lifted on top of B4",
					acc:  "same_user",

					hasMedia: true,
					want:     true,
				},
			},
			{
				// But we don't want it if there's no media
				text: "Seeing S20 was epic!",
				acc:  "same_user",

				want: false,

				quoted: &ttest{
					text: "Starship SN20 being lifted on top of B4",
					acc:  "same_user",

					hasMedia: true,
					want:     true,
				},
			},
			{
				text:     "Another picture of S20",
				acc:      "other_user",
				hasMedia: true,
				want:     true,

				quoted: &ttest{
					text: "Starship SN20 being lifted on top of B4",
					acc:  "same_user",

					hasMedia: true,
					want:     true,
				},
			},

			{
				text: "😂",
				acc:  "random_user",
				want: false,

				quoted: &ttest{
					text: "Starship S20 is looking interesting today",
					want: true,
				},
			},
			{
				text: "Nice render!",
				acc:  "random_user",
				want: false,

				quoted: &ttest{
					acc:  "Starship 20 in orbit",
					want: false,
				},
			},
			{
				text: "Nice info here!",
				acc:  "random_user",
				want: false,

				quoted: &ttest{
					text:     "Starship SN20 being lifted on top of B4",
					hasMedia: true,
					want:     true,
				},
			},
		},
	)
}
