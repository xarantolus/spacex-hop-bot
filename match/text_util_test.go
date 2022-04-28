package match

import (
	"reflect"
	"strings"
	"testing"
)

func Test_startsWithAnyGeneric(t *testing.T) {
	var searchedPrefixes = []string{"test", "best", "rest", "more than one word"}

	tests := []struct {
		argText string
		want    bool
	}{
		{"testing is nice", true},
		{"wrongprefixtesting is nice", false},
		{"the test keyword can be at any point in the string", true},
		{"we want to support more than one word", true},
		{"we want to support less than one word", false},

		{"the #test hashtag should still be recognized", true},
		{"also @test should work", true},
		{"#test at the beginning", true},
		{"@test should work", true},
		{"this \"test\" seems ok", true},
		{"this \"tes\"t seems ok", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := startsWithAny(tt.argText, searchedPrefixes...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_startsWithAnyStarshipAntikeywords(t *testing.T) {
	tests := []struct {
		argText string
		want    bool
	}{
		{"road closure with no information where it is", false},
		{"https://shop.blueorigin.com/collections/new/products/new-glenn-108th-scale", true},
		{"", false},
		{"Starship-SLS is a good idea", true},
		{"KSP is my favourite game!", true},
		{"What if I add many spaces before       KSP", true},
		{"-   .  +  - . k . -  .  .. -   . . -.- .", false},
		{"-   .  +  - . ksp . -  .  .. -   . . -.- .", true},
		{"Project DogeCoin onto a Starship!", true},
		{"Starship reentering Kerbin's atmosphere", true},
		{"GSE Tank 6 rolling out", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := startsWithAny(strings.ToLower(tt.argText), antiStarshipKeywords...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_containsAnyGeneric(t *testing.T) {
	var searchedInfixes = []string{"test", "best", "rest", "more than one word"}

	tests := []struct {
		argText string
		want    bool
	}{
		{"testing is nice", true},
		{"wrongprefixtesting is nice", true},
		{"the test keyword can be at any point in the string", true},
		{"we want to support more than one word", true},
		{"we want to support less than one word", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := containsAny(tt.argText, searchedInfixes...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_compose(t *testing.T) {
	var all = []string{"test"}

	tests := []struct {
		arg  [][]string
		want []string
	}{
		{
			arg:  [][]string{all, {"another"}, {"3"}},
			want: []string{"test", "another", "3"},
		},
		{
			arg:  [][]string{all, {"another", "duplicate"}, {"duplicate"}, {"3"}},
			want: []string{"test", "another", "duplicate", "3"},
		},
		{
			arg:  [][]string{all, {"duplicate", "duplicate"}, {"duplicate"}, {"3"}},
			want: []string{"test", "duplicate", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := compose(tt.arg...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compose(%v) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func Test_ignoreSpaces(t *testing.T) {
	tests := []struct {
		arg        []string
		wantResult []string
	}{
		{
			arg: []string{"orbital launch tower", "orbital tower"},
			wantResult: []string{
				"orbital launch tower", "orbitallaunchtower", "orbital-launch-tower", "orbital_launch_tower", "orbitallaunch tower", "orbital launchtower",

				"orbital tower", "orbitaltower", "orbital-tower", "orbital_tower",
			},
		},
		{
			arg:        []string{"a b c"},
			wantResult: []string{"a b c", "abc", "a-b-c", "a_b_c", "ab c", "a bc"},
		},
		{
			arg:        []string{"a", "b", "c"},
			wantResult: []string{"a", "b", "c"},
		},
		{
			arg:        []string{"a b", "c"},
			wantResult: []string{"a b", "ab", "a-b", "a_b", "c"},
		},
		{
			arg:        []string{"starship", "superheavy", "super heavy"},
			wantResult: []string{"starship", "superheavy", "super heavy", "super-heavy", "super_heavy"},
		},
		{
			arg:        []string{"sea level"},
			wantResult: []string{"sea level", "sealevel", "sea-level", "sea_level"},
		},
		{
			arg:        []string{"mc gregor live"},
			wantResult: []string{"mc gregor live", "mcgregorlive", "mc-gregor-live", "mc_gregor_live", "mcgregor live", "mc gregorlive"},
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if gotResult := ignoreSpaces(tt.arg); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ignoreSpaces(%v) = %v, want %v", tt.arg, gotResult, tt.wantResult)
			}
		})
	}
}
