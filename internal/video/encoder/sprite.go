package encoder

import (
	"image"

	"github.com/stashapp/stash/pkg/ffmpeg2"
)

type SpriteScreenshotOptions struct {
	Width int
}

func SpriteScreenshot(encoder ffmpeg2.FFMpeg, input string, seconds float64, options SpriteScreenshotOptions) (image.Image, error) {
	ssOptions := ffmpeg2.ScreenshotOptions{
		OutputPath: "-",
		OutputType: ffmpeg2.ScreenshotOutputTypeBMP,
		Width:      options.Width,
	}

	args := ffmpeg2.ScreenshotTime(input, seconds, ssOptions)

	return doGenerateImage(encoder, input, args)
}

func SpriteScreenshotSlow(encoder ffmpeg2.FFMpeg, input string, frame int, options SpriteScreenshotOptions) (image.Image, error) {
	ssOptions := ffmpeg2.ScreenshotOptions{
		OutputPath: "-",
		OutputType: ffmpeg2.ScreenshotOutputTypeBMP,
		Width:      options.Width,
	}

	args := ffmpeg2.ScreenshotFrame(input, frame, ssOptions)

	return doGenerateImage(encoder, input, args)
}
