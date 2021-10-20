package ffmpeg

import (
	"fmt"
	"image"
	"strings"
)

type SpriteScreenshotOptions struct {
	Time  float64
	Width int
}

func (e *Encoder) SpriteScreenshot(probeResult VideoFile, options SpriteScreenshotOptions) (image.Image, error) {
	args := []string{
		"-v", "error",
		"-ss", fmt.Sprintf("%v", options.Time),
		"-i", probeResult.Path,
		"-vframes", "1",
		"-vf", fmt.Sprintf("scale=%v:-1", options.Width),
		"-c:v", "bmp",
		"-f", "rawvideo",
		"-",
	}
	data, err := e.run(probeResult.Path, args, nil)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(data)

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, err
}
