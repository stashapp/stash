package transcoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubtitleExtractArgs(t *testing.T) {
	args := SubtitleExtract("/videos/test.mkv", 2, "/tmp/out.vtt")
	assert.Equal(t, []string{
		"-v", "error",
		"-y",
		"-i", "/videos/test.mkv",
		"-map", "0:2",
		"-vn",
		"-an",
		"-f", "webvtt",
		"/tmp/out.vtt",
	}, args.Args())
}
