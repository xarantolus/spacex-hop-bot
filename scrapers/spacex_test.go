package scrapers

import (
	"testing"
)

func TestNameRegex(t *testing.T) {
	var testdata = []struct {
		in, out string
	}{
		// Prototype naming convention
		{"SN20", "SN20"},
		{"BN2", "BN2"},

		// New naming convention
		{"B2", "B2"},
		{"S20", "S20"},

		// Full names
		{"Starship 20", "Starship 20"},
		{"StarShip 20", "StarShip 20"},

		{"Booster 2", "Booster 2"},

		// In a text
		{"Starship S16 Orbital flight attempt", "S16"},

		// Should *not* be matched
		{"Starship", ""},
	}

	for _, d := range testdata {
		if res := ShipNameRegex.FindString(d.in); res != d.out {
			t.Errorf("Expected ShipNameRegex to match %q with input %q, but got %q", d.out, d.in, res)
		}
	}
}
