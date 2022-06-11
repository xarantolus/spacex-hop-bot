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
				text:     "Tripod is chilling at McGregor! Possible firing coming up.\nhttp://nsf.live/mcgregor",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Something about the SpaceX Starship and SuperHeavy Project",
				acc:  "FAANews",
				want: true,
			},
			{
				text:     "Raptor 2 from different angles",
				hasMedia: true,
				want:     true,
			},
			{
				text: "The Booster Quick Disconnect did a high speed retraction test at 11:06 local time.",
				want: true,
			},
			{
				// this is the test user ID; we don't want to retweet our own tweets
				userID: testBotSelfUserID,
				text:   "S20 standing on the pad",
				want:   false,
			},
			{
				text: "S20 standing on the pad",
				want: true,
			},
			{
				text: "Last month the FCC asked SpaceX a series of questions about their next generation Starlink constellation, Starlink v2.0, or \"Gen2\". SpaceX responded back and affirms they will definitely launch it using Starship and could be ready as soon as (Future Date)",
				want: true,
			},
			{
				text: "Raptor fired on McGregor Live 5 minutes ago",
				want: true,
			},
			{
				text: "Merlin fired on McGregor Live 5 minutes ago",
				want: false,
			},
			{
				text: "STARBASE big voice announcement: \"Overhead drone operations will occur for the next hour.\" Get ready for B4 lift off the orbital launch mount! #Starbase #Starship #SpaceX",
				want: true,
			},
			{
				text: "NASA has selected Starship for an additional mission to the Moon with astronauts as part of the Artemis program! http://nasa.gov/press-release/nasa-provides-update-to-astronaut-moon-lander-plans-under-artemis",
				acc:  "SpaceX",
				want: true,

				quoted: &ttest{
					acc:      "NASAArtemis",
					text:     "Artemis III astronauts will land on the surface aboard a @SpaceX Starship Human Landing System. These new opportunities are for missions beyond #Artemis III.",
					hasMedia: true,
					want:     true,
				},
			},
			{
				text: "I'm LIVE from the Starbase build site, tune in: https://youtube.com/watch?v=3195jmsakdfj",
				want: true,
			},
			{
				text:     "The launch tower at Starbase will help stack Starship and catch the Super Heavy rocket booster",
				acc:      "SpaceX",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Likely the last cryo proof before the orbital test flight",
				want: true,
			},
			{
				text:     "Full stack #SpaceX",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Full stack #starship #sn20 #bn4 #Starbase",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Full stack is imminent",
				location: match.BocaChicaBeachPlaceID,
				want:     true,
			},
			{
				text:     "Hoping to see a full stack again, potentially today!\n\nhttps://nasaspaceflight.com/starbaselive",
				hasMedia: true,
				want:     true,
				quoted: &ttest{
					text:     "Mighty fine morning… 😎🚀 - @NASASpaceflight",
					location: match.SpaceXLaunchSiteID,
					hasMedia: true,
					want:     true,
				},
			},
			{
				text: "(1/9) Lets look at what is still remaining to complete the #Widebay now that this poll has ended. \n\nIt helps to do some comparisons between the existing #Highbay and the new #Widebay. The best place to start is with the #BridgeCranes\n\n📷:@CSI_Starbase",
				want: true,
			},
			{
				text:     "Set the controls for the heart of the sun. @SpaceX #Starship SuperHeavy armed stacking",
				location: "3309acacf870f6f5", // Matamoros, Tamaulipas
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Looks like progress on the launch tower at the cape is moving along. Already have at least two sides to one section up.",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Hey @CSI_Starbase, looks like progress on the launch tower at the cape is moving right along. Already have at least two sides to one section up.",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Ground breaking of Phase 1 of StarFactory at StarBase. Building is expected to be over 300,000sq. ft. by 60' tall and will replace all the large production tents. \nBuilding will extend almost to fence of Hwy 4!",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "There is only one kind of tower I know of that has steel segments that look like this...👀\nSeen traversing the NASA Causeway headed toward KSC this afternoon were likely among the first parts of a Starship orbital integration tower arriving in Florida🚀",
				acc:      "TrevorMahlmann",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Deimos in the port",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Phobos in the port",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Another piece added to Phobos @nasaspaceflight",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Another piece added to Deimos @nasaspaceflight",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "The SpaceX Deimos rig is moving and departing Port of Brownsville!",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Deimos update pt. 69\nDeimos is now pierside for refit and generator work for 2 weeks before departure. 4 large pressure storage vessels arrived by barge today. @elonmusk",
				want: true,
			},
			{
				text:     "SpaceX’s Phobos is actively being worked on before it becomes a Starship sea launching platform https://spaceexplored.com/guides/phobos-starship/",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Phobos in the port",
				location: match.PascagoulaPlaceID,
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Here are some helicopter shots I took today of the new #SpaceX Roberts Road site & the work on #Starship construction @ 39A. Notice all the new land clearing the new structures & the buildout of Hanger X for F9 refurbishment. In coop w/@FarryFaz. #NASA",
				acc:      "GregScott_photo",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Flying around SpaceX’s Hanger X on Roberts Road today with @GregScott_photo",
				acc:      "FarryFaz",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Spaceport Deimos (named after a martian Moon) on the move",
				want: true,
			},
			{
				text:     "Chopsticks are going up! Will they be used to remove the booster from the orbital launch mount?\nhttps://nasaspaceflight.com/starbaselive",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Progress MS-19/80P: The Progress MS-19/80P cargo ship is in the final stages of a 2-day rendezvous with the International Space Station; NASA TV is providing live coverage: https://youtube.com/...",
				want: false,
			},
			{
				text: "CraneX lift within the next 10 minutes, watch with @jessica_kirsh: https://youtube.com/...",
				want: true,
			},
			{
				text: "CraneX photo by @BocaChicaGal",
				want: true,
			},
			{
				text: "High bay stacking continues",
				want: true,
			},
			{
				text: "New grid fins arrived at the build site",
				want: true,
			},
			{
				text: "New booster grid fin arrived",
				want: true,
			},
			{
				text:     "Super Heavy Grid Fins.\n\n#SpaceX\n\n📸 for @Teslarati",
				hasMedia: true,
				want:     true,
			},
			{
				text: "SuperHeavy standing still",
				want: true,
			},
			{
				text: "Looks like progress on the Deimos sea-launch platform",
				want: true,
			},
			{
				text: "So I’ve checked in with what’s happened in the US whilst I slept and they have a fully stacked Starship now - and, apparently, some kind of soup police?",
				want: false,
			},
			{
				text: "Stage Zero - Immensely Complex! (and I freaking love it) Check out the Ship QD Arm time-lapse. Watch live at here for this awesome view. https://youtu.be/7zsl4q6fwfQ",
				want: true,
			},
			{
				text: "The armchair scientists in the NSF chat get worse, always someone like: spaceX can’t stack S20 because I over boiled my eggs and burnt my toast therefore the FAA won’t let them stack it!",
				want: false,
			},
			{
				text: "full stack @nasaspaceflight",
				want: true,
			},
			{
				text: "full stack @spacex",
				want: true,
			},
			{
				text: "PA Announcement just now: “attention on the pad, we’re 15 minutes away from ship proof.” @NASASpaceflight",
				want: true,
			},
			{
				text: "This has nothing to do with Starships...its just amazing...",
				want: false,
			},
			{
				text: "They are currently testing the catch arms at the launch site",
				want: true,
			},
			{
				text: "The ship lift points have been deployed on the chopsticks",
				want: true,
			},
			{
				text: "Ship now attached to the chopsticks",
				want: true,
			},
			{
				text: "Booster now standing next to the launch tower",
				want: true,
			},
			{
				text: "Second day of the presentation week. 2 days remain. Thanks to @BocaChicaGal we are getting some amazing views of the launch site. Crossing fingers for a stacking today!\n\nWatch all the action the next few days at: http://nasaspaceflight.com/starbaselive @NASASpaceflight",
				want: true,
			},
			{
				text: "#ataresults Congrats to B18 Doubles champs Name and Name!! 😉\n\nBoth Name and Name also got 3rd in their respective B18 draws.",
				want: false,
			},
			{
				text: "Someone got DISSED by @NASASpaceflight / Chris Bergin on a public form for posting StarshipGazer and Labpadre views..\n\nIt’s obvious they care about the money more than anything else..",
				want: false,
			},
			{
				text: "So far Starship Troopers is like Fascist Degrassi and it’s brilliant",
				want: false,
			},
			{
				text: "Starship and 9/11 in the same tweet",
				want: false,
			},
			{
				text: "Compressed 24 hours of remote 4k video into 30 secs. on SpaceX's landing pad for NROL-87 mission. @NASASpaceflight @SpaceX",
				want: false,
			},
			{
				text: "Next raptor delivery seen on @nasaspaceflight cam",
				want: true,
			},
			{
				text: "In case anybody cares, @RoyalCaribbean has yet to respond to my multiple requests for comment on the Harmony of the Seas range violation. Sent an initial inquiry immediately after yesterday's scrub.",
				acc:  "nextspaceflight",
				want: false,
			},
			{
				text: "The new #SpaceX facilities at #NASA's Roberts Rd in Cape Canaveral is at full go. Land has been cleared, a new booster refurbishment building is being completed and lots more unknown structures are in progress. Stay in tune for more updates as it progresses @elonmusk @MarcusHouse",
				want: true,
			},
			{
				text:     "Robert's Road update!",
				acc:      "FarryFaz",
				hasMedia: true,
				want:     true,
			},
			{
				text: "A booster loadspreader is being lifted. Don't panic yet, this could be a sign of depressurization.\n\n📷 @NASASpaceflight",
				want: true,
			},
			{
				text: "The load spreader is up.\n\nhttp://nasaspaceflight.com/starbaselive",
				want: true,
			},
			{
				text: "Fresh out of YC S21, Epsilon3 raises seed round to continue modernizing space and launch operations. https://buff.ly/3rUsIFn",
				want: false,
			},
			{
				text:     "Progress on the HLS Starship variant",
				hasMedia: true,
				want:     true,
			},
			{
				text: "It looks like Ship 22’s aft section was moved into the mid bay following speculation that it was set to be scrapped.\n\n📸: @LabPadre",
				acc:  "spacex360",
				want: true,
			},
			{
				text:     "Starship is simply beautiful",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Stop simping for elon just because you like Starship",
				want: false,
			},
			{
				text: "Booster cryo proof coming up!",
				want: true,
			},
			{
				text: "Booster cryoproof coming up!",
				want: true,
			},
			{
				text: "New LN2 tanker spotted at the #OTF",
				want: true,
			},
			{
				text: "New LN2 tanker spotted at the #OrbitalTankFarm",
				want: true,
			},
			{
				text:     "Sea-level raptors",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Sealevel raptors",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Olmos Park restaurant Glass and Plate Restaurant closed ‘due to lack of employees’\nhttps://www.expressnews.com/food/restaurants/article/Olmos-Park-restaurant-Glass-and-Plate-Restaurant-16766960.php",
				acc:  "ExpressNews",
				want: false,
			},
			{
				text:     "SpaceX is testing the lift arms strength with giant bags. So cool.",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Chopsticks did several mini raises this afternoon. Can't wait for a complete raise/swing/open exercise soon.\n#Starbase  #Starship  #SpaceX\n 📸 Me for WAI Media @felixschlang",
				acc:      "CosmicalChief",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Mechazilla's Chopsticks received extra attention on Sunday, with what looks like a lifting bar fit check.\n\nMeanwhile, Booster 3 slice and dice operations continue.\n\nMary (@BocaChicaGal) views:\nhttp://nasaspaceflight.com/starbaselive",
				acc:  "nasaspaceflight",
				want: true,
			},
			{
				text: "My starlink dish arrived",
				want: false,
			},
			{
				text: "https://shop.blueorigin.com/collections/new/products/new-glenn-108th-scale",
			},
			{
				text:        "@elonmusk Hey I've got a great idea about what you can use as a mass simulator for the first Starship Orbital Demo! This! https://shop.blueorigin.com/collections/new/products/new-glenn-108th-scale\n\nBlue Origin can finally say they went orbital with New Glenn then!",
				tweetSource: match.TweetSourceLocationStream,
				want:        false,
			},
			{
				text: "A road closure for Starbase is now active",
				want: true,
			},
			{
				text:     "As the world turns here at Starbase, Texas @SpaceX continues to push engineering limits to the moon! Watch the full time lapse here: https://youtu.be/tiuLI8t5JWU #SpaceX #Starbase #Texas #Starship",
				acc:      "LabPadre",
				hasMedia: true,
				want:     true,
			},
			{
				text: "Timelapse of today's chopstick test thus far.\nhttp://nasaspaceflight.com/starbaselive",
				acc:  "nextspaceflight",
				want: true,
			},
			{
				text: "Timelapse of the chopsticks slowly on the move in Starbase.\n\nhttp://nasaspaceflight.com/starbaselive",
				acc:  "nextspaceflight",
				want: true,
			},
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
				text: "Orbital Launch Integration Tower catching arms have begun their first vertical move visible on @LabPadre rover cam 1 & 2. @elonmusk @SpaceX",
				want: true,
			},
			{
				text: "OLIT catching arms have begun their first vertical move visible on @LabPadre rover cam 1 & 2. @elonmusk @SpaceX",
				want: true,
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
				text:        "Just heard over the SpaceX PA system that S24 will be lifted shortly!",
				acc:         "locationstream+location",
				location:    match.SpaceXLaunchSiteID,
				want:        true,
				tweetSource: match.TweetSourceLocationStream,
			},
			{
				text:        "Just heard over the SpaceX PA system that S24 will be lifted shortly!",
				acc:         "locationstream+location+media",
				location:    match.SpaceXLaunchSiteID,
				hasMedia:    true,
				want:        true,
				tweetSource: match.TweetSourceLocationStream,
			},
			{
				text:     "Just heard over the SpaceX PA system that S24 will be lifted shortly!",
				acc:      "location+media",
				location: match.SpaceXLaunchSiteID,
				hasMedia: true,
				want:     true,
			},
			{
				// Even ignored accounts are allowed if they tweet from an important place
				text:           "Starship!",
				accDescription: "Dogecoin & crypto fan",
				location:       match.SpaceXBuildSiteID,
				hasMedia:       true,
				want:           true,
			},
			{
				// Even ignored accounts are allowed if they tweet from an important place
				text:     "Took this image, it's similar to my 3D renders!",
				userID:   match.TestIgnoredUserID,
				location: match.SpaceXBuildSiteID,
				hasMedia: true,
				want:     true,
			},
			{
				// It contains an antikeyword (dogecoin), but because it's from a SpaceX site and has media we retweet it anyways
				text:     "Full Stack #DogecoinToTheMoon",
				location: match.SpaceXLaunchSiteID,
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Ship 20 prepares for stacking.",
				location: match.SpaceXLaunchSiteID,
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Starbase flyover @leifviper, Stroker, @rookisaacman @slickf16 @SpaceX",
				acc:      "KiddPoteet",
				location: match.SpaceXBuildSiteID,
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Catch arm lift tests underway! 🦾",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "Full stack",
				location: match.SpaceXLaunchSiteID,
				hasMedia: true,
				want:     true,
			},
			{
				location: match.SpaceXLaunchSiteID,
				hasMedia: true,
				want:     true,
			},
			{
				text: "#BREAKING Ooooh. Launchpad giant voice just announced they are clearing the launch tower and launch mount. Maybe we will be getting some chopstick heavy lifting going today!!! #Starbase #Starship #SpaceX",
				want: true,
			},
			{
				text:     "#BREAKING Ooooh. Launchpad giant voice just announced they are clearing the launch tower and launch mount. Maybe we will be getting some chopstick heavy lifting going today!!! #Starbase #Starship #SpaceX",
				want:     true,
				location: match.SpaceXLaunchSiteID,
			},
			{
				acc:      "AdamCuker",
				text:     "SpaceX Raptor 2.0 rocket engine test last night in McGregor, Texas @SpaceX. The video is dark due to dense fog. This was the loudest I've ever heard it. Residences could hear this test over 30+ miles away!#SpaceXtest (Incredible Roar)\nRaptor 2 Test Video: https://youtu.be/BKR3WE55cQ8",
				hasMedia: true,
				want:     true,
			},
			{
				text:     "New booster standing tall in the setting sun tonight @SpaceX #McGregor Rocket  B1072 @elonmusk",
				acc:      "jswartzphoto",
				hasMedia: true,
				location: match.SpaceXMcGregorPlaceID,
				want:     false,
			},
			{
				text:     "Alright SpaceXers. This a FH center?\n @SpeedyPatriot13 @BoosterSpX @NolanTrees\n\n#Falcon9 #FalconHeavy #McGregor #SpaceX",
				acc:      "jswartzphoto",
				hasMedia: true,
				location: match.SpaceXMcGregorPlaceID,
				want:     false,
			},
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
			// Tweet threads with questions
			{
				text: "How many S20 cryogenic pressure test(s)?",
				acc:  "Starship_Sults",
				want: false,

				parent: &ttest{
					text: "How many S20 static fires?",
					acc:  "Starship_Sults",
					want: false,
				},
			},
			{
				text:     "How many S20 cryogenic pressure test(s)?",
				acc:      "Starship_Sults",
				hasMedia: true,
				want:     true,

				parent: &ttest{
					text:     "How many S20 static fires?",
					acc:      "Starship_Sults",
					hasMedia: true,
					want:     true,
				},
			},
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
				text: "Mars colonial transporter",
				acc:  "elonmusk",
				want: true,

				parent: &ttest{
					text: "Absolutely bonkers! I can’t wait to hear the rumble of a rocket that’s over twice as powerful as the Saturn V 😍",
					acc:  "Erdayastronaut",
					want: true,

					quoted: &ttest{
						text: "33 Raptor rocket engines, each producing 230 metric tons of force",
						acc:  "elonmusk",
						want: true,
					},
				},
			},
			{
				parent: &ttest{
					text: "What's going on at starbase?",
					want: true,
				},

				text: "Rocket seems fine",
				acc:  "elonmusk",
				want: true,
			},
			// Top-level tweet
			{
				acc:  "elonmusk",
				text: "Tesla and Starship engines are currently the two hardest problems.",
				want: true,
			},
			{
				text: "I usually drive an alpha build, but switch to beta right before release so I know what Tesla owners are getting",
				acc:  "elonmusk",
				want: false,

				parent: &ttest{
					text: "What version are you driving",
					acc:  "teslaownersSV",
					want: false,

					parent: &ttest{
						text: "This is pretty good. 10.12 will have major improvements for tricky unprotected lefts & heavy traffic in general. We’re also making good progress with single stack.",
						acc:  "elonmusk",
						want: false,

						parent: &ttest{
							text: "#FSDBeta 10.11.1 has huge improvements. Best build so far. @elonmusk",
							acc:  "teslaownersSV",
							want: false,
						},
					},
				},
			},
			{
				text: "♥️♥️ NASA ♥️♥️",
				acc:  "elonmusk",
				want: true,

				quoted: &ttest{
					text:     "Artemis III astronauts will touch down on the Moon aboard a @SpaceX Starship Human Landing System. We will be asking U.S. companies to develop astronaut Moon landers for @NASAArtemis missions beyond #Artemis III: https://go.nasa.gov/3IxKUuL",
					acc:      "NASA",
					hasMedia: true,
					want:     true,
				},
			},
			// This is a real tweet
			{
				text: "16 story tall rocket, traveling several times faster than a bullet, backflips & fires engines to return to launch site",
				acc:  "elonmusk",
				want: false,

				parent: &ttest{
					text: "View of Falcon 9's stage separation from ground cameras",
					acc:  "SpaceX",
					want: false,
				},
			},
			// Same thing, but for starship
			{
				text: "16 story tall rocket, traveling several times faster than a bullet, backflips & fires engines to return to launch site",
				acc:  "elonmusk",
				want: true,

				parent: &ttest{
					text: "View of Starship stage separation from ground cameras",
					acc:  "SpaceX",
					want: true,
				},
			},

			// Starship tweet with a follow up of an unrelated question
			{
				text: "Will be ready in Q4.",
				acc:  "elonmusk",
				want: false,

				parent: &ttest{
					// Unrelated question about tesla
					text: "When will FSD beta be available for all?",
					want: false,

					parent: &ttest{
						text: "Starship tweet",
						acc:  "elonmusk",
						want: true,
					},
				},
			},
			// Elon answering to an ignored account (e.g. a 3d animation)
			{
				text: "Yes on both counts. That would be a great outcome for civilization.",
				acc:  "elonmusk",
				want: true,

				parent: &ttest{
					text: "So Mars will eventually get its own Mechazilla?\n\nWould there be any value in eventually building Super Heavy Boosters on Mars as a launch platform for outer solar system missions?",
					want: true,

					parent: &ttest{
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

								want: true,
							},
						},
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

					want: true,
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
				text: "Later this year, remaining fussy bits will be gone, allowing deletion of shroud",
				acc:  "elonmusk",
				want: true,
				parent: &ttest{
					text:     "Raptor V1 vs Raptor V2.  Greatly simplified whilst increasing thrust. Costs half as much.",
					acc:      "nasaspaceflight",
					hasMedia: true,
					want:     true,
				},
			},
			{
				text: "ASAP: Other Starship HLS risks \"include things like software and hardware integration, flight rate, hardware turnaround times and reuse.\"\n\"NASA is working on all of those and trying to make sure that it's comfortable with the approach that is being proposed by SpaceX.\"",
				acc:  "thesheetztweetz",
				want: true,

				parent: &ttest{
					text: "ASAP also identified landing technologies accuracy / stability / hazard avoidance as another top risk to HLS Starship.\nSome of the mitigations include the fact that there will be uncrewed test landings prior to the first human landing.\"",
					acc:  "thesheetztweetz",
					want: true,

					parent: &ttest{
						text: "ASAP identified the required Starship cryo-fluid transfer/management of refueling as a top risk to the HLS program.",
						acc:  "thesheetztweetz",
						want: true,

						parent: &ttest{
							text: "ASAP: SpaceX also gave NASA a \"good understanding of some of the challenges\" the company is having with Raptor engine production.",
							acc:  "thesheetztweetz",
							want: true,

							parent: &ttest{
								text: "ASAP: \"NASA also conducted some site visits to Boca Chica and Hawthorne that indicated there's been significant progress in the overall production of Starship and HLS.\"",
								acc:  "thesheetztweetz",
								want: true,

								parent: &ttest{
									text:     "NASA's Aerospace Safety Advisory Panel says SpaceX this month provided the agency with \"an integrated master schedule\" on HLS Starship development.",
									acc:      "thesheetztweetz",
									want:     true,
									hasMedia: true,
								},
							},
						},
					},
				},
			},
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

func TestHQMediaTweet(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				acc:  "cnunezimages",
				text: "Tweet without media",
				want: false,
			},
			{
				acc:      "cnunezimages",
				text:     "Tweet with media",
				hasMedia: true,
				want:     true,
			},
			{
				acc:      "cnunezimages",
				text:     "Glowing 20 - @elonmusk @spacex\n#Starbase  #BocaChicaToMars #iCANimagine http://cnunezimages.com @SpaceIntellige3\n_____________________________________\n- Image Taken: " + time.Now().Format("January 2, 2006") + " -",
				hasMedia: true,
				want:     true,
			},
		},
	)
}

func TestAdTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				// If it's without an buy link, we retweet it
				text: "New Starship S20 stuff!",
				want: true,
			},
			{
				// It's clearly trying to sell something
				text: "New Starship S20 stuff!\n\nhttps://www.etsy.com/lang/listing/19835918395819385/whatever",
				want: false,
			},
			{
				text:     "Chopsticks Can Accomplish Anything!\nRight @elonmusk?\n\nhttps://etsy.com/ca/listing/1155335887",
				hasMedia: true,
				want:     false,
			},
			{
				text:     "LAST CHANCE! 20% OFF ENDS TONIGHT at 11:59PT!\n\nhttp://etsy.com/shop/thelaunchpadshop\n\n#blackfriday #spacex #nasa #space #starship #starbase",
				hasMedia: true,
				want:     false,
			},
		},
	)
}

func TestIgnoredTweets(t *testing.T) {
	testStarshipRetweets(t,
		[]ttest{
			{
				text: "Sri Lanka's central bank has secured foreign exchange to pay for fuel and cooking gas shipments that will ease crippling shortages, its governor said, but police fired tear gas and water canon to push back student protesters",
				want: false,
			},
			{
				text: "Elon Musk says Twitter deal “cannot move forward” until proof of number of fake accounts is provided",
				want: false,
			},
			{
				// Contains "STF" (suborbital tank farm) and "open"
				text: "Fantastic news!🥳 MIRI, the UK's main contribution to the @ESA_Webb, has opened its eye to the sky! 💫 MIRI's painstaking alignment process was supported by scientists and engineers from @ukatc and RAL Space. Huge congratulations to everyone involved!👏 👉https://ralspace.stfc.ac.uk/Pages/Webb%E2%80%99s-coolest-instrument-captures-first-star.aspx",
				want: false,
			},
			{
				text: "Amber Heard describes using make-up and 'super heavy, red matte lipstick' to conceal injuries before appearing on talk show",
			},
			{
				text: "#DYK our Orbital Outpost SN-5000 service module offers free-flyer, logistics services for low-Earth orbit & cislunar destinations? It can carry 6,500 lbs. of pressurized & 3,500 lbs. of unpressurized cargo with 3 external mounting locations! #SpaceSymposium #TeamSNC",
				acc:  "SierraNevCorp",
				want: false,
			},
			{
				text: "Buckle up for another stacked cast this Tuesday morning 🥞🥞🥞\n💸Elon Musk named to $TWTR's board of directors\n🔥@WarnerMedia\nCEO @jasonkilar joins for an exclusive interview\n+ $AMZN looks to take on Starlink with satellite-based internet",
				want: false,
			},
			{
				text: "Someone else brought it up in a conversation and just...",
				want: false,
				parent: &ttest{
					text: "So did we all collectively forget that dearMoon was Yusaku Maezawa’s lunar “will you be my girlfriend” competition for a month or so 2 years ago because I sure as hell did",
					want: false,
				},
			},
			{
				text:     "FORGET CAPE CANAVERAL…\nLook what happened at Cape Cornwall last week\nProud to share @RealHomerHickam’s  #RocketBoys story\nAnd how a dad, @RookIsaacman went to space\nAnd look… this dad’s got one of @ElonMusk’s StarShip on his shoulders!\nLaunching the #OurMillion22👩‍👩‍👦‍👦 appeal",
				hasMedia: true,
				want:     false,
			},
			{
				text:     "LAST NIGHT: The West bound lane of HWY 100 was temporarily closed last night due to a vehicle on fire.  Teenagers driving back from SPI noticed their car was emitting smoke. They pulled over, got out of the vehicle before the vehicle became engulfed in flames.",
				acc:      "SheriffGarza",
				hasMedia: true,
				want:     false,
			},
			{
				text:     "Cape Canaveral's iconic Missile Row played a vital role in the development of US rocketry. After many years of silence, rockets are returning to its historic launch pads. You can read about its history and its future over at @NASASpaceflight",
				hasMedia: true,
				want:     false,
			},
			{
				text:     "Finally, at long last we have the @StackUpDotOrg\nRex Brickheads that we gave away during our #CallToArms last year! \nI’ll be getting these shipped out/delivered to their proper homes soon. 👀\nIf you won one, keep your eyes out on your mail!",
				hasMedia: true,
				want:     false,
			},
			{
				text: "Bob and fast-boat Maverick departed Port Canaveral a short while ago to support the Transporter-4 mission.\nhttp://nasaspaceflight.com/fleetcam",
				want: false,
			},
			{
				text: "C'mon @thedrive, NASA's Artemis I is sitting on the launch pad right now, and 3 more human rated @NASA_Orion crew modules are in production right which will all send humans to the Moon.",
				want: false,
			},
			{
				text: "The Apollo 16 space vehicle bracketed with the Launch Umbilical Tower (LUT) to the left and the Mobile Service Structure (MSS) to the right during March, 1972. The space vehicle was rolled out to the Launch Pad twice due to  spacecraft repairs. #Apollo16 #Apollo50",
				want: false,
			},
			{
				text: "Also, cryogenics are generally terrible for ballistic missile systems.",
				want: false,
				acc:  "nextspaceflight",

				parent: &ttest{
					text: "I am sorry, but this excuse is total BS. It is industry standard to broadcast the primary countdown loop. Pretty much all of the U.S. launch providers do it, and NASA did it during Shuttle. If you are worried about ITAR, you make the callout on a different loop.",
					want: false,
					acc:  "nextspaceflight",

					quoted: &ttest{
						text: "NASA's Tom Whitmeyer says press won't have access to countdown loops for the Space Launch System's wet dress rehearsal next week (breaking frm tradition) because of ITAR concerns and fears that adversaries will glean cryogenic timing info for clues into ballistic missile systems.",

						want: false,
					},
				},
			},
			{
				text: "On March 26, the Solar Orbiter spacecraft completed its closest pass of the Sun yet — passing just 48 million kilometers above its surface. On its trek to perihelion, Solar Orbiter took some of the highest resolution images of the Sun ever taken.\n\nARTICLE:\nhttps://www.nasaspaceflight.com/2022/03/solar-orbiter-close-pass/",
				want: false,
			},
			{
				text: "... in the later seasons of TNG ... Becky’s protest of Romeo and Jules during that whole homophobia story arc in S12",
				want: false,
			},
			{
				text: "finally catching up on s10 of call the midwife",
				want: false,
			},
			{
				accDescription: "3d artist",
				text:           "Starship orbital test flight",
				hasMedia:       true,
				want:           false,
			},
			{
				text:     "Fin whales for @_! There were so many and they were so close to the ship",
				hasMedia: true,
				want:     false,
			},
			{
				text:     "This happened b4",
				hasMedia: true,
				want:     false,
			},
			{
				text: "Doug arrived at Port Canaveral just after 2am this morning with the fairing from the Starlink 4-10 mission.\n\nhttp://nasaspaceflight.com/fleetcam",
				want: false,
			},
			{
				text: "SpaceX to launch AST SpaceMobile's first orbital cell towers https://teslarati.com/spacex-ast-spacemobile-bluebird-launch-contract/ by @13ericralph31",
				want: false,
			},
			{
				text: "Two retractions closer to rollout day for #Artemis I! Platforms D & E are retracted inside High Bay 3 of the Vehicle Assembly Building at @NASAKennedy. Look at our Moon rocket revealing itself.",
				acc:  "NASA_SLS",
				want: false,
			},
			{
				text: "A famous marketing stunt inspired Toronto-based SpaceRyde's founders to create an out of this world innovation.🎈 Now, working with MDA's LaunchPad Program, they are ready for lift-off.",
				want: false,
			},
			{
				text: "Woah, we're halfway there 🎶🚀\n\nAs of right now, half of the platforms surrounding @NASA_SLS and @NASA_Orion have been retracted in High Bay 3 of the Vehicle Assembly Building. Who else is ready to see this massive Moon rocket roll out to Launch Complex 39B?",
				want: false,
			},
			{
				text: "As an example, here is S00012 Vanguard 2 rocket, where you can see oscillations in the solution during 1961, and the correction (red original, blue corrected) in 1962-1964",
				want: false,
			},
			{
				text:     "Yet another beautiful spacex start from the cape",
				hasMedia: true,
				want:     false,
			},
			{
				text:     "Check out this time lapse showing the retraction of half of Platform C in High Bay 3 of the Vehicle Assembly Building today at @NASAKennedy. On March 17, @NASA_SLS & @NASA_Orion will roll out to Launch Pad 39B for wet dress rehearsal for @NASAArtemis I.",
				acc:      "NASAGroundSys",
				hasMedia: true,
				want:     false,
			},
			{
				text:     "Half of Platform C in High Bay 3 of the Vehicle Assembly Building is now retracted, continuing to reveal more of @NASA_SLS & @NASA_Orion.",
				acc:      "NASAGroundSys",
				hasMedia: true,
				want:     false,
			},
			{
				text: "Sagittarius B2 is a giant molecular cloud at the center of the Milky Way, and it's made of alcohol. An ester, ethyl formate is also responsible for the flavour of raspberries, leading some articles to postulate the cloud is smelling of ‘raspberry rum’ https://buff.ly/3vEXeEa",
				want: false,
			},
			{
				text:        "“Pathfinders” in the build yard @ Smokey’s Outpost…#SpaceX #Starbase #BocaChicaToMars #B4 #S20 #OLM #Bocachica #ElonMusk @elonmusk #DogecoinToTheMoon #Moon #Mars #Launchpad #modeling",
				tweetSource: match.TweetSourceLocationStream,
				location:    "somelocation",
				want:        false,
			},
			{
				text: "In just under an hour, Starlink 4-9 is set to lift off from LC-39A.\nLaunch is scheduled for 9:25 AM ET.\nBooster 1060 will be making it's 11th flight.\n📷: Me for @SuperclusterHQ",
				want: false,
			},
			{
				text:     "#GOEST lifts off from the pad!\nHere's a view from one of my remote cameras at the launch pad.",
				hasMedia: true,
				want:     false,
			},
			{
				text: "Watch a live view of United Launch Alliance’s Atlas 5 rocket on its launch pad at Cape Canaveral, awaiting liftoff… https://t.co/yqxePuGbld",
				want: false,
			},
			{
				text: "https://starshipsls.wixsite.com/futureastronaut/post/spacex-to-launch-starlink-4-8-with-49-more-starlink-satellites Tomorrow, @SpaceX will launch 49 more Starlink satellites on the Starlink 4-8 mission. Find out more in my new article.",
				want: false,
			},
			{
				text:     "S15 in the wild #ForzaHorizon5",
				hasMedia: true,
				want:     false,
			},
			{
				text: "Either a BUK or S300 was active in Kyiv tonight, engaged targets https://t.co/XpAXN1ra6B",
				want: false,
			},
			{
				text: "RS-25 ignition! Static Fire on the A-1 test stand at Stennis!",
				want: false,
			},
			{
				text: "SpaceX has confirmed separation  and a nominal orbital insertion of the Group 4-8 Starlink stack launched this morning at 9:44am EST on a Falcon 9 this morning from SLC-40 at Cape Canaveral Space Force Station.\n\nThis successfully concludes today's mission.",
				want: false,
			},
			{
				text: "Caught a pic of Deimos next to Mars!",
				want: false,
			},

			{
				text: "On Friday, February 18, 2022, Sheriff's Deputies responded to an address in Olmito in reference to shots fired. While en route, information was given that the suspect was in the Villa Los Pinos Subdivision and had shot in the direction of a victim.",
				acc:  "CameronCountySO",
				want: false,
			},
			{
				text: "Finally high-speed Internet in the middle of the #TexasHillCountry ❗️👏🏻👏🏻 Thanks @elonmusk @SpaceX",
				want: false,
			},
			{
				text: "USAF B52 CHIEF11 visible again over eastern mediterranean",
				want: false,
			},
			{
				text: "Saturn V rolling out of High bay 3",
				want: false,
			},
			{
				text: "A second S-400 bn has also been identified",
				want: false,
			},
			{
				text: "What a day!!!! It has already seen a #Starship full-stack, next will be a launch from Kourou(@OneWeb), then I cross my fingers for @Astra's 4th attempt to launch and the big finale will be @elonmusk's update on the Starship program this evening!",
				want: false,
			},
			{
				text: "Starship SN15\n\nGet it on https://opensea.io/some/link",
				want: false,
			},
			{
				text: "Starship NFT dropping soon!",
				want: false,
			},
			{
				text: "Starship on OpenSea now available!!!",
				want: false,
			},
			{
				text: "This would be an 11-day turnaround for Pad 39A.",
				acc:  "same_user",
				want: false,

				parent: &ttest{
					text:     "Per @NASASpaceflight forum members analyzing FAA & NOTAM alerts, it looks like SpaceX's second Starlink launch of the year/month (Starlink 4-6) is probably scheduled NET ~9pm EST, January 17th!\n\nAnd my guess was only off by one day!\nhttps://forum.nasaspaceflight.com/index.php?topi",
					hasMedia: true,
					acc:      "same_user",

					want: false,
				},
			},
			{
				text:     "Cutting B83 makes sense; it’s not practically useable (and probably not lawfully & feasibly against too many targets). Good legacy move too: Biden becomes the president to cut the last megaton-class weapon in the US arsenal.",
				hasMedia: true,

				want: false,
			},
		},
	)
}
