package jobs

import "testing"

func Test_generateReviewTweetText(t *testing.T) {
	tests := []struct {
		arg      string
		withTags bool
		want     string
	}{
		{
			arg:      "The Environmental review is now complete.",
			withTags: true,
			want:     "Review update: The Environmental review is now complete.\n\n@elonmusk\nhttps://www.permits.performance.gov/permitting-project/spacex-starshipsuper-heavy-launch-vehicle-program-spacex-boca-chica-launch-site",
		},
		{
			arg:      `The target date of milestone "Issuance of a Final EA" of the "Environmental Assessment (EA)" has changed from "2021-12-31" to "2022-02-28"`,
			withTags: false,
			want:     `The target date of milestone "Issuance of a Final EA" of the "Environmental Assessment (EA)" has changed from "2021-12-31" to "2022-02-28"`,
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := generateReviewTweetText(tt.arg, tt.withTags); got != tt.want {
				t.Errorf("generateReviewTweetText() = %v, want %v", got, tt.want)
			}
		})
	}
}
