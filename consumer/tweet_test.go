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
