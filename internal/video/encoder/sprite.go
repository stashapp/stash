package encoder

import (
	"image"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
)

type SpriteScreenshotOptions struct {
	Width int
}

func SpriteScreenshot(encoder ffmpeg.FFMpeg, input string, seconds float64, options SpriteScreenshotOptions) (image.Image, error) {
	ssOptions := transcoder.ScreenshotOptions{
		OutputPath: "-",
		OutputType: transcoder.ScreenshotOutputTypeBMP,
		Width:      options.Width,
	}

	args := transcoder.ScreenshotTime(input, seconds, ssOptions)

	return doGenerateImage(encoder, input, args)
}

func SpriteScreenshotSlow(encoder ffmpeg.FFMpeg, input string, frame int, options SpriteScreenshotOptions) (image.Image, error) {
	ssOptions := transcoder.ScreenshotOptions{
		OutputPath: "-",
		OutputType: transcoder.ScreenshotOutputTypeBMP,
		Width:      options.Width,
	}

	args := transcoder.ScreenshotFrame(input, frame, ssOptions)

	return doGenerateImage(encoder, input, args)
}
