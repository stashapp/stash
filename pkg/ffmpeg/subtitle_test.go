package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTextSubtitleCodec(t *testing.T) {
	textCodecs := []string{"subrip", "srt", "ass", "ssa", "mov_text", "webvtt", "text", "ASS", " SubRip "}
	for _, c := range textCodecs {
		assert.True(t, IsTextSubtitleCodec(c), "expected %q to be text subtitle codec", c)
	}

	bitmapCodecs := []string{"hdmv_pgs_subtitle", "dvd_subtitle", "dvb_subtitle", "xsub", "pgssub", "", "h264"}
	for _, c := range bitmapCodecs {
		assert.False(t, IsTextSubtitleCodec(c), "expected %q not to be text subtitle codec", c)
	}
}

func TestGetTextSubtitleStreams(t *testing.T) {
	vf := &VideoFile{
		JSON: FFProbeJSON{
			Streams: []FFProbeStream{
				{Index: 0, CodecType: "video", CodecName: "h264"},
				{Index: 1, CodecType: "audio", CodecName: "aac"},
				{Index: 2, CodecType: "subtitle", CodecName: "subrip"},
				{Index: 3, CodecType: "subtitle", CodecName: "hdmv_pgs_subtitle"},
				{Index: 4, CodecType: "subtitle", CodecName: "ass"},
			},
		},
	}

	vf.SubtitleStreams = vf.getSubtitleStreams()
	assert.Len(t, vf.SubtitleStreams, 3)

	textStreams := vf.GetTextSubtitleStreams()
	assert.Len(t, textStreams, 2)
	assert.Equal(t, 2, textStreams[0].Index)
	assert.Equal(t, 4, textStreams[1].Index)
}
