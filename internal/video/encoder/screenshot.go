package encoder

import (
	"github.com/stashapp/stash/pkg/ffmpeg2"
)

type ScreenshotOptions struct {
	OutputPath string
	Quality    int
	Time       float64
	Width      int
}

func Screenshot(encoder ffmpeg2.FFMpeg, input string, options ScreenshotOptions) error {
	ssOptions := ffmpeg2.ScreenshotOptions{
		OutputPath: options.OutputPath,
		OutputType: ffmpeg2.ScreenshotOutputTypeImage2,
		Quality:    options.Quality,
		Width:      options.Width,
	}

	args := ffmpeg2.ScreenshotTime(input, options.Time, ssOptions)

	return doGenerate(encoder, input, args)
}
