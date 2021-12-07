package jobs

import "testing"

func Test_generateReviewTweetText(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{"The Environmental review is now complete.", "Review update: The Environmental review is now complete.\n\n@elonmusk\nhttps://www.permits.performance.gov/permitting-project/spacex-starshipsuper-heavy-launch-vehicle-program-spacex-boca-chica-launch-site"},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := generateReviewTweetText(tt.arg); got != tt.want {
				t.Errorf("generateReviewTweetText() = %v, want %v", got, tt.want)
			}
		})
	}
}
