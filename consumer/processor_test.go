package consumer

import (
	"testing"

	"github.com/dghubble/go-twitter/twitter"
)

// TODO: Should have retweeted https://twitter.com/NASASpaceflight/status/1473717315371941902

func TestProcessor_isQuestion(t *testing.T) {

	tests := []struct {
		input string
		want  bool
	}{
		{"Here is a picture of the booster", false},
		{"is this picture showing a lox tank?", true},
		{"@elonmusk here is a speculative question about superhevay?", true},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			var tw = twitter.Tweet{
				FullText: tt.input,
			}
			if got := isQuestion(&tw); got != tt.want {
				t.Errorf("Processor.isQuestion(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func Test_isTagsOnly(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Hello World!", false},
		{"Hello #Starbase, nice #Starships!", false},

		{"#Starbase #Starbase #SpaceX #Starship @elonmusk", true},
		{"#Starship #Starbase", true},
		{"@elonmusk", true},
		{"@elonmusk #starship", true},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := isTagsOnly(tt.input); got != tt.want {
				t.Errorf("isTagsOnly(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
