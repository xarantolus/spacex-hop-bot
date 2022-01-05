package match

import (
	"testing"
)

// This is the most general test, it just contains a bunch of
// possible tweet texts and defines whether we want them to be matched
func TestStarshipTextMatch(t *testing.T) {
	tests := []struct {
		text string
		want bool
	}{
		{"Chopsticks moving", false},
		{"Chopsticks moving\n\nhttp://nasaspaceflight.com/starbaselive", true},
		{"Chopsticks moving\n\nhttps://nasaspaceflight.com/starbaselive", true},
		{"Timelapse of today's chopstick test thus far.", false},
		{"Timelapse of today's chopstick test thus far.\nhttp://nasaspaceflight.com/starbaselive", true},
		{"Lots of frost on Booster 4, which now seems to be decreasing, potentially signifying we’re getting close to the conclusion of todays cryogenic proof test.\nCamera with flash @NASASpaceflight\nRed circle https://youtu.be/B1IbMBhococ\n\n#SpaceX @elonmusk", true},
		{"On @nasaspaceflight stream it looks like the road is open again, so no Cryo Proof today it seems.", true},
		{"Orbital launch table is venting! 📷 @NASASpaceflight", true},
		{"The orbital tank farm is venting.\n\n📷 @NASASpaceflight", true},
		{"LN2 loading has begun. Photos from @NASASpaceflight stream", true},
		{"GSE tank rolling down the road. Photos from @labpadre stream", true},
		{"Tomorrow at 5:06am ET, @SpaceX's 24th cargo resupply mission will lift off from Launch Complex 39A. Weather officials continue to predict a 30% chance of favorable weather conditions.\nTune in today at 12pm ET for our prelaunch news conference on @NASA TV: https://go.nasa.gov/3q9ByhW", false},
		{"Took a quick shot of the raptor engines delivered today", true},
		{"Took a quick shot of the raptor vacuum delivered today", true},
		{"Cryo proof coming up. #SpaceX #Starbase #Texas", true},
		{"Launch tower in LC-49 going up fast!", true},
		{"NASA Is Conducting An Environmental Assessment Of New SpaceX Proposal To Build A Starship Launch Site At Launch Complex-49 In Florida", true},
		{"I finally got to see for myself Ship 21’s nosecone", true},
		{"Peekaboo I see u #B5 #SpaceX", true},
		{"There’s more time to submit your comments on the Draft Programmatic Environmental Assessment for the proposed @SpaceX Starship/Super Heavy project in Boca Chica, Texas. Comment by Nov. 1. The new public meeting dates are Oct. 18th and 20th. Learn more at http://bit.ly/2YcScDe.", true},
		{"Raptor 63 being lifted up to the booster", true},
		{"SN10", true},
		{"Nasa moves forward with environmental assessment of LC-49", true},
		{"Starbase gse tanks are venting", true},
		{"BN10", true},
		{"A great animated video of a @SpaceX #Starship launch, orbit & landing using the new 'chopsticks' system", false},
		{"Starship SN10", true},
		{"SuperHeavy Booster", true},
		{"#starship", true},
		{"Already seeing some depress venting from Ship 20 here at the launch site!", true},
		{"deimos oil rig", true},
		{"S20 was transported to the suborbital pad b", true},
		{`
		SpaceX is targeting Wednesday, March 24 for launch of 60 Starlink satellites from Space Launch Complex 40 (SLC-40) at Cape Canaveral Space Force Station in Florida. The instantaneous window is at 4:28 a.m. EDT, or 8:28 UTC, and a backup opportunity is available on Thursday, March 25 at 4:06 a.m. EDT, or 8:06 UTC.
		
		The Falcon 9 first stage rocket booster supporting this mission previously supported launch of the GPS-III Space Vehicle 03 and Turksat 5A missions in addition to three Starlink missions. Following stage separation, SpaceX will land Falcon 9's first stage on the “Of Course I Still Love You” droneship, which will be located in the Atlantic Ocean. One half of Falcon 9's fairing supported the Sentinel-6A mission and the other supported a previous Starlink mission.
		`, false},
		{"I have received an Alert notice for tomorrow, July 19. Possible static fire attempt between noon and 10 p.m. on Booster B3.", true},

		// Oil platform names need at least a bit of context
		{"Starship will land on Deimos", true},
		{"SpaceX's Phobos launch platform", true},
		{"phobos & deimos in the port", true},
		{"deimos & phobos in the port", true},

		{"Pad announcement over the speakers: clearing pad for static fire", false},
		{"Lifting off next to the tower of Launch Complex 39A", false},
		{"starbase", false},
		{"Deimos in the Ocean", false},
		{"Phobos in the port", false},
		{"SpaceX will conduct a static fire later today", false},
		{"Samsung S22 Ultra", false},
		{"I mention Starship. $RKLB", false},
		{"The gse tanks are venting", false},
		{"Last week saw extensive work on Ship 20's TPS tiles, Booster 4 grew some engines, a GSE tank was tested, and some jets made an impressive flyover! Beyond Starbase, BO, China, and Astra all made launches, & Firefly Aerospace prepares their first flight!", false},
		{"Starship 20 and #Shenzhou12 ", false},
		{"Galaxy S22 Ultra", false},
		{"a known S-300 deployment", false},
		{"#NewProfilePic #SN15", false},
		{"Rolls-Royce chosen by U.S. for new B-52 engines in contract worth up to $2.6 bln", false},
		{"B-52", false},
		{"Bald Eagle in Canada flying over the water at the Canadian Raptor Conservancy by Fred Johns", false},
		{"No TFR posted for today", false},
		{"Unrelated doge coin tweet that also contains the keyword Starship", false},
		{"Unrelated tesla tweet", false},
		{"this tweet is not starship related", false},
		{"Starlink Mission", false},
		{"With a few B61 parachute deployment photos too!", false},
		{"B61 used ribbon parachute! Can see the ribbons clearly in this photo", false},
		{"The raptors were the most impressive part of the movie", false},
		{"Starship and Dogecoin", false},
		{"You want to upgrade# your music IQ sub to http://YouTube.com/a1madethebeat you won’t be sorry.  Find #A1madethebeat on all social media.\nIf you want to learn:\n- to make #beats\n- to record\n- to #mix your record\n- to use #mpc #akaiforce #S2400 #Ableton #logicprox #luna\n-#collab on a record", false},
		{"The nearest supernova since the one described by Kepler in 1604 is SN 1987A, around 168,000 LY distant. (Camera Hubble)", false},
	}

	matcher := NewStarshipMatcherForTests()
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := matcher.StarshipText(tt.text, antiStarshipKeywords, false); got != tt.want {
				t.Errorf("StarshipText(%q, antiStarshipKeywords) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestContainsStarshipAntiKeyword(t *testing.T) {
	t.Run("Single AntiKeyword", func(t *testing.T) {
		contains := ContainsStarshipAntiKeyword("The SLS is making progress faster than Starship")
		if !contains {
			t.Errorf("Expected antiKeyword 'SLS' to be detected, but wasn't")
		}
	})
}
