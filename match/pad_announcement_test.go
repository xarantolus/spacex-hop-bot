package match

import "testing"

func TestIsPadAnnouncement(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"STARBASE big voice announcement: \"Overhead drone operations will occur for the next hour.\" Get ready for B4 lift off the orbital launch mount! #Starbase #Starship #SpaceX", true},
		{"PA: 15 minutes away from clearing the orbital pad for ship proof.", true},
		{"PA Announcement just now: “attention on the pad, we’re 15 minutes away from ship proof.” @NASASpaceflight", true},
		{"Just heard a pad announcement. Very hard to hear, sounded like some sort of pad operations. Could be some sort of testing?", true},
		{"Pad speakers: clearing everything for booster lift", true},
		{"Pad announcement over the speakers: clearing pad for static fire", true},
		{"#BREAKING Ooooh. Launchpad giant voice just announced they are clearing the launch tower and launch mount. Maybe we will be getting some chopstick heavy lifting going today!!! #Starbase #Starship #SpaceX", true},
		{"LabPadre announced something", false},
		{"Just heard over the SpaceX PA system that S24 will be lifted shortly!", true},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := IsPadAnnouncement(tt.input); got != tt.want {
				t.Errorf("IsPadAnnouncement(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
