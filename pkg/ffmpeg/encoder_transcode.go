package ffmpeg

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func calculateTranscodeScale(probeResult VideoFile, maxTranscodeSize models.StreamingResolutionEnum) string {
	maxSize := 0
	switch maxTranscodeSize {
	case models.StreamingResolutionEnumLow:
		maxSize = 240
	case models.StreamingResolutionEnumStandard:
		maxSize = 480
	case models.StreamingResolutionEnumStandardHd:
		maxSize = 720
	case models.StreamingResolutionEnumFullHd:
		maxSize = 1080
	case models.StreamingResolutionEnumFourK:
		maxSize = 2160
	}

	// get the smaller dimension of the video file
	videoSize := probeResult.Height
	if probeResult.Width < videoSize {
		videoSize = probeResult.Width
	}

	// if our streaming resolution is larger than the video dimension
	// or we are streaming the original resolution, then just set the
	// input width
	if maxSize >= videoSize || maxSize == 0 {
		return "iw:-2"
	}

	// we're setting either the width or height
	// we'll set the smaller dimesion
	if probeResult.Width > probeResult.Height {
		// set the height
		return "-2:" + strconv.Itoa(maxSize)
	}

	return strconv.Itoa(maxSize) + ":-2"
}
