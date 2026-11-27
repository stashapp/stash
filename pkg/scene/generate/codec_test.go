package generate

import (
	"slices"
	"testing"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
)

func TestVideoEncoderArgs(t *testing.T) {
	tests := []struct {
		codec   models.GeneratedVideoCodec
		preset  string
		wantEnc ffmpeg.VideoCodec
		want    []string
	}{
		{"", "slow", ffmpeg.VideoCodecLibX264, []string{"-profile:v", "high", "-preset", "slow", "-crf", "21"}},
		{models.GeneratedVideoCodecH264, "slow", ffmpeg.VideoCodecLibX264, []string{"-preset", "slow", "-crf", "21"}},
		{models.GeneratedVideoCodecHevc, "slow", ffmpeg.VideoCodecLibX265, []string{"-preset", "slow", "-crf", "26", "-tag:v", "hvc1"}},
		{models.GeneratedVideoCodecHevc, "veryslow", ffmpeg.VideoCodecLibX265, []string{"-preset", "veryslow"}}, // user's preview preset is honoured
		{models.GeneratedVideoCodecH264, "veryslow", ffmpeg.VideoCodecLibX264, []string{"-preset", "veryslow"}},
		{models.GeneratedVideoCodecAv1, "veryslow", ffmpeg.VideoCodecSVTAV1, []string{"-preset", "4", "-crf", "33"}},
		{models.GeneratedVideoCodecAv1, "bogus", ffmpeg.VideoCodecSVTAV1, []string{"-preset", "7"}},
	}

	for _, tt := range tests {
		enc, args := videoEncoderArgs(tt.codec, tt.preset, 21)
		if enc != tt.wantEnc {
			t.Errorf("%q: encoder = %v, want %v", tt.codec, enc, tt.wantEnc)
		}
		// check each wanted flag/value pair appears adjacent in args
		for i := 0; i+1 < len(tt.want); i += 2 {
			idx := slices.Index(args, tt.want[i])
			if idx < 0 || idx+1 >= len(args) || args[idx+1] != tt.want[i+1] {
				t.Errorf("%q: args %v missing %s %s", tt.codec, args, tt.want[i], tt.want[i+1])
			}
		}
		if tt.codec != models.GeneratedVideoCodecH264 && tt.codec != "" && slices.Contains(args, "-profile:v") {
			t.Errorf("%q: x264-only -profile:v leaked into args %v", tt.codec, args)
		}
	}
}
