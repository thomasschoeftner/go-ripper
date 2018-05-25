package metainfo

import "testing"

func TestIsVideo(t *testing.T) {
	x := &VideoMetaInfo{}
	if !IsVideo(x) {
		t.Error("video meta info was expected to qualify as video, but failed to")
	}
}