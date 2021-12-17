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

func Test_startsWithAnyStarship(t *testing.T) {
	tests := []struct {
		argText string
		want    bool
	}{
		{"KSP is my favourite game!", true},
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
