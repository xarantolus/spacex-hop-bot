package match

import "testing"

func TestIsPadAnnouncement(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Just heard a pad announcement. Very hard to hear, sounded like some sort of pad operations. Could be some sort of testing?", true},
		{"Pad speakers: clearing everything for booster lift", true},
		{"Pad announcement over the speakers: clearing pad for static fire", true},

		{"LabPadre announced something", false},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := IsPadAnnouncement(tt.input); got != tt.want {
				t.Errorf("IsPadAnnouncement(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
