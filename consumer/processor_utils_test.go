package consumer

import (
	"testing"

	"github.com/dghubble/go-twitter/twitter"
)

func TestProcessor_isQuestion(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Here is a picture of the booster", false},
		{"is this picture showing a lox tank?", true},
		{"@elonmusk here is a speculative question about superheavy?", true},
		{"What is going on here?\nhttps://www.youtube.com/watch?v=GP18t7ivstY", true},
		{"No questions here!\nhttps://www.youtube.com/watch?v=GP18t7ivstY", false},
		{"https://www.youtube.com/watch?v=GP18t7ivstY", false},
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
		{"Hi @someone", false},
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

func Test_isReactionGIF(t *testing.T) {
	tests := []struct {
		tweet *twitter.Tweet
		want  bool
	}{
		{
			tweet: &twitter.Tweet{
				ExtendedEntities: &twitter.ExtendedEntity{
					Media: []twitter.MediaEntity{
						{
							Type: "video/gif",
						},
					},
				},
			},
			want: true,
		},
		{
			tweet: &twitter.Tweet{
				ExtendedEntities: &twitter.ExtendedEntity{
					Media: []twitter.MediaEntity{
						{
							Type: "image/jpeg",
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := isReactionGIF(tt.tweet); got != tt.want {
				t.Errorf("isReactionGIF() = %v, want %v", got, tt.want)
			}
		})
	}
}
