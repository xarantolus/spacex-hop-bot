package consumer

import (
	"testing"

	"github.com/dghubble/go-twitter/twitter"
)

func TestProcessor_isQuestion(t *testing.T) {
	var p = new(Processor)

	tests := []struct {
		input string
		want  bool
	}{
		{"is this picture showing a lox tank?", false},
		{"@elonmusk here is a speculative question about superhevay?", true},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			var tw = twitter.Tweet{
				FullText: tt.input,
			}
			if got := p.isQuestion(&tw); got != tt.want {
				t.Errorf("Processor.isQuestion(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
