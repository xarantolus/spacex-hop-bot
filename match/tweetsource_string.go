// Code generated by "stringer -type=TweetSource"; DO NOT EDIT.

package match

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TweetSourceUnknown-0]
	_ = x[TweetSourceLocationStream-1]
	_ = x[TweetSourceKnownList-2]
	_ = x[TweetSourceTimeline-3]
	_ = x[TweetSourceTrustedUser-4]
}

const _TweetSource_name = "TweetSourceUnknownTweetSourceLocationStreamTweetSourceKnownListTweetSourceTimelineTweetSourceTrustedUser"

var _TweetSource_index = [...]uint8{0, 18, 43, 63, 82, 104}

func (i TweetSource) String() string {
	if i < 0 || i >= TweetSource(len(_TweetSource_index)-1) {
		return "TweetSource(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TweetSource_name[_TweetSource_index[i]:_TweetSource_index[i+1]]
}
