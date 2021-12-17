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
		{"Launch tower in LC-49 going up fast!", true},
		{"NASA Is Conducting An Environmental Assessment Of New SpaceX Proposal To Build A Starship Launch Site At Launch Complex-49 In Florida", true},
		{"I finally got to see for myself Ship 21’s nosecone", true},
		{"Peekaboo I see u #B5 #SpaceX", true},
		{"There’s more time to submit your comments on the Draft Programmatic Environmental Assessment for the proposed @SpaceX Starship/Super Heavy project in Boca Chica, Texas. Comment by Nov. 1. The new public meeting dates are Oct. 18th and 20th. Learn more at http://bit.ly/2YcScDe.", true},
		{"Rolls-Royce chosen by U.S. for new B-52 engines in contract worth up to $2.6 bln", false},
		{"B-52", false},
		{"Raptor 63 being lifted up to the booster", true},
		{"Bald Eagle in Canada flying over the water at the Canadian Raptor Conservancy by Fred Johns", false},
		{"No TFR posted for today", false},
		{"SN10", true},
		{"BN10", true},
		{"A great animated video of a @SpaceX #Starship launch, orbit & landing using the new 'chopsticks' system", false},
		{"Starship SN10", true},
		{"SuperHeavy Booster", true},
		{"Unrelated doge coin tweet that also contains the keyword Starship", false},
		{"Unrelated tesla tweet", false},
		{"this tweet is not starship related", false},
		{"Starlink Mission", false},
		{`
SpaceX is targeting Wednesday, March 24 for launch of 60 Starlink satellites from Space Launch Complex 40 (SLC-40) at Cape Canaveral Space Force Station in Florida. The instantaneous window is at 4:28 a.m. EDT, or 8:28 UTC, and a backup opportunity is available on Thursday, March 25 at 4:06 a.m. EDT, or 8:06 UTC.

The Falcon 9 first stage rocket booster supporting this mission previously supported launch of the GPS-III Space Vehicle 03 and Turksat 5A missions in addition to three Starlink missions. Following stage separation, SpaceX will land Falcon 9's first stage on the “Of Course I Still Love You” droneship, which will be located in the Atlantic Ocean. One half of Falcon 9's fairing supported the Sentinel-6A mission and the other supported a previous Starlink mission.
`, false},
		{"I have received an Alert notice for tomorrow, July 19. Possible static fire attempt between noon and 10 p.m. on Booster B3.", true},
		{"Starship and Dogecoin", false},
		// Oil platform names need at least a bit of context
		{"Starship will land on Deimos", true},
		{"Deimos in the Ocean", false},
		{"SpaceX's Phobos launch platform", true},
		{"Phobos in the port", false},
		{"Samsung S22 Ultra", false},
		{"I mention Starship. $RKLB", false},
		{"Last week saw extensive work on Ship 20's TPS tiles, Booster 4 grew some engines, a GSE tank was tested, and some jets made an impressive flyover! Beyond Starbase, BO, China, and Astra all made launches, & Firefly Aerospace prepares their first flight!", false},
		{"Starship 20 and #Shenzhou12 ", false},
		{"Galaxy S22 Ultra", false},
		{"a known S-300 deployment", false},
		{"#NewProfilePic #SN15", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := StarshipText(tt.text, antiStarshipKeywords); got != tt.want {
				t.Errorf("StarshipText(%q, antiStarshipKeywords) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}
