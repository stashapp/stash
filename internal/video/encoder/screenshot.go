package encoder

import (
	"github.com/stashapp/stash/pkg/ffmpeg2"
	"github.com/stashapp/stash/pkg/video"
)

type ScreenshotOptions struct {
	OutputPath string
	Quality    int
	Time       float64
	Width      int
}

func Screenshot(encoder ffmpeg2.FFMpeg, input string, options ScreenshotOptions) error {
	ssOptions := video.ScreenshotOptions{
		OutputPath: options.OutputPath,
		OutputType: video.ScreenshotOutputTypeImage2,
		Quality:    options.Quality,
		Width:      options.Width,
	}

	args := video.ScreenshotTime(input, options.Time, ssOptions)

	return doGenerate(encoder, input, args)
}
