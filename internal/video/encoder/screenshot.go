package encoder

import (
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
)

type ScreenshotOptions struct {
	OutputPath string
	Quality    int
	Time       float64
	Width      int
}

func Screenshot(encoder ffmpeg.FFMpeg, input string, options ScreenshotOptions) error {
	ssOptions := transcoder.ScreenshotOptions{
		OutputPath: options.OutputPath,
		OutputType: transcoder.ScreenshotOutputTypeImage2,
		Quality:    options.Quality,
		Width:      options.Width,
	}

	args := transcoder.ScreenshotTime(input, options.Time, ssOptions)

	return doGenerate(encoder, input, args)
}
