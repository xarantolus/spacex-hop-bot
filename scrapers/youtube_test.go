package scrapers

import "testing"

// Make sure URL generation works correctly
func TestLiveVideo_URL(t *testing.T) {
	tests := []struct {
		lv   *LiveVideo
		want string
	}{
		{
			lv: &LiveVideo{
				VideoID: "gA6ppby3JC8",
			},
			want: "https://www.youtube.com/watch?v=gA6ppby3JC8",
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := tt.lv.URL(); got != tt.want {
				t.Errorf("LiveVideo.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}
