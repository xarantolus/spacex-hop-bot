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

		{"Starship S16 Orbital flight attempt", "S16"},
	}

	for _, d := range testdata {
		if res := ShipNameRegex.FindString(d.in); res != d.out {
			t.Errorf("Expected ShipNameRegex to match %q with input %q, but got %q", d.out, d.in, res)
		}
	}
}
