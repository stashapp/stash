package generate

import (
	"strconv"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
)

// svtAV1Presets maps x264 preset names onto SVT-AV1's numeric presets (0 slowest, 13 fastest).
var svtAV1Presets = map[string]string{
	"ultrafast": "12",
	"superfast": "11",
	"veryfast":  "10",
	"faster":    "9",
	"fast":      "8",
	"medium":    "7",
	"slow":      "6",
	"slower":    "5",
	"veryslow":  "4",
}

var videoEncoders = map[models.GeneratedVideoCodec]ffmpeg.VideoCodec{
	models.GeneratedVideoCodecH264: ffmpeg.VideoCodecLibX264,
	models.GeneratedVideoCodecHevc: ffmpeg.VideoCodecLibX265,
	models.GeneratedVideoCodecAv1:  ffmpeg.VideoCodecSVTAV1,
}

// VideoEncoder returns the software encoder used for codec. Unset or unknown codecs use H.264.
func VideoEncoder(codec models.GeneratedVideoCodec) ffmpeg.VideoCodec {
	if enc, ok := videoEncoders[codec]; ok {
		return enc
	}
	return ffmpeg.VideoCodecLibX264
}

// videoEncoderArgs returns the encoder for codec and its output args. Callers speak x264:
// preset and crf are translated so each codec lands at roughly the same visual quality.
func videoEncoderArgs(codec models.GeneratedVideoCodec, preset string, crf int) (ffmpeg.VideoCodec, ffmpeg.Args) {
	enc := VideoEncoder(codec)
	args := ffmpeg.Args{"-pix_fmt", "yuv420p"}

	switch enc {
	case ffmpeg.VideoCodecLibX265:
		// x265 veryslow is ~10x slower than x264 veryslow (33s per 20s marker clip); slow keeps bulk generation usable
		if preset == "veryslow" {
			preset = "slow"
		}
		// x265 CRF 28 ~ x264 CRF 23 per the ffmpeg wiki; hvc1 tag is required for Safari playback
		args = append(args, "-preset", preset, "-crf", strconv.Itoa(crf+5), "-tag:v", "hvc1", "-x265-params", "log-level=error")
	case ffmpeg.VideoCodecSVTAV1:
		p, ok := svtAV1Presets[preset]
		if !ok {
			p = svtAV1Presets["medium"]
		}
		args = append(args, "-preset", p, "-crf", strconv.Itoa(crf+12))
	default:
		args = append(args, "-profile:v", "high", "-level", "4.2", "-preset", preset, "-crf", strconv.Itoa(crf))
	}

	return enc, args
}
