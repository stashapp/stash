package ffmpeg

import (
	"fmt"
	"image"
	"strings"
)

type SpriteScreenshotOptions struct {
	Time  float64
	Frame int
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

// SpriteScreenshotSlow uses the select filter to get a single frame from a videofile instead of seeking
// It is very slow and should only be used for files with very small duration in secs /  frame count
func (e *Encoder) SpriteScreenshotSlow(probeResult VideoFile, options SpriteScreenshotOptions) (image.Image, error) {
	args := []string{
		"-v", "error",
		"-i", probeResult.Path,
		"-vsync", "0", // do not create/drop frames
		"-vframes", "1",
		"-vf", fmt.Sprintf("select=eq(n\\,%d),scale=%v:-1", options.Frame, options.Width), // keep only frame number options.Frame
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
