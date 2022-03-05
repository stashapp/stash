package encoder

import (
	"image"

	"github.com/stashapp/stash/pkg/ffmpeg2"
	"github.com/stashapp/stash/pkg/video"
)

type SpriteScreenshotOptions struct {
	Width int
}

func SpriteScreenshot(encoder ffmpeg2.FFMpeg, input string, seconds float64, options SpriteScreenshotOptions) (image.Image, error) {
	ssOptions := video.ScreenshotOptions{
		OutputPath: "-",
		OutputType: video.ScreenshotOutputTypeBMP,
		Width:      options.Width,
	}

	args := video.ScreenshotTime(input, seconds, ssOptions)

	return doGenerateImage(encoder, input, args)
}

func SpriteScreenshotSlow(encoder ffmpeg2.FFMpeg, input string, frame int, options SpriteScreenshotOptions) (image.Image, error) {
	ssOptions := video.ScreenshotOptions{
		OutputPath: "-",
		OutputType: video.ScreenshotOutputTypeBMP,
		Width:      options.Width,
	}

	args := video.ScreenshotFrame(input, frame, ssOptions)

	return doGenerateImage(encoder, input, args)
}
