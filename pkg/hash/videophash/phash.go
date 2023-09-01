package videophash

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type phashConfig struct {
	screenshotSize int
	columns        int
	rows           int
}

func Generate(encoder *ffmpeg.FFMpeg, videoFile *models.VideoFile) (*uint64, error) {
	// Original algorithm hardcoded to 5 x 5 grid
	columns := 5

	// Per https://github.com/stashapp/stash/issues/3722 use a smaller amount of frames when
	// the length of the video is shorter. Aiming for around 7 seconds between frames.
	switch {
	case videoFile.Duration < 46:
		columns = 2
	case videoFile.Duration < 91:
		columns = 3
	case videoFile.Duration < 151:
		columns = 4
	}

	config := phashConfig{
		screenshotSize: 160,
		columns:        columns,
		rows:           columns,
	}

	sprite, err := generateSprite(encoder, videoFile, config)
	if err != nil {
		return nil, err
	}

	hash, err := goimagehash.PerceptionHash(sprite)
	if err != nil {
		return nil, fmt.Errorf("computing phash from sprite: %w", err)
	}
	hashValue := hash.GetHash()
	return &hashValue, nil
}

func generateSpriteScreenshot(encoder *ffmpeg.FFMpeg, input string, t float64, config phashConfig) (image.Image, error) {
	options := transcoder.ScreenshotOptions{
		Width:      config.screenshotSize,
		OutputPath: "-",
		OutputType: transcoder.ScreenshotOutputTypeBMP,
	}

	args := transcoder.ScreenshotTime(input, t, options)
	data, err := encoder.GenerateOutput(context.Background(), args, nil)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	return img, nil
}

func combineImages(images []image.Image, config phashConfig) image.Image {
	width := images[0].Bounds().Size().X
	height := images[0].Bounds().Size().Y
	canvasWidth := width * config.columns
	canvasHeight := height * config.rows
	montage := imaging.New(canvasWidth, canvasHeight, color.NRGBA{})
	for index := 0; index < len(images); index++ {
		x := width * (index % config.columns)
		y := height * int(math.Floor(float64(index)/float64(config.rows)))
		img := images[index]
		montage = imaging.Paste(montage, img, image.Pt(x, y))
	}

	return montage
}

func generateSprite(encoder *ffmpeg.FFMpeg, videoFile *models.VideoFile, config phashConfig) (image.Image, error) {
	logger.Infof("[generator] generating phash sprite for %s", videoFile.Path)

	// Generate sprite image offset by 5% on each end to avoid intro/outros
	chunkCount := config.columns * config.rows
	offset := 0.05 * videoFile.Duration
	stepSize := (0.9 * videoFile.Duration) / float64(chunkCount)
	var images []image.Image
	for i := 0; i < chunkCount; i++ {
		time := offset + (float64(i) * stepSize)

		img, err := generateSpriteScreenshot(encoder, videoFile.Path, time, config)
		if err != nil {
			return nil, fmt.Errorf("generating sprite screenshot: %w", err)
		}

		images = append(images, img)
	}

	// Combine all of the thumbnails into a sprite image
	if len(images) == 0 {
		return nil, fmt.Errorf("images slice is empty, failed to generate phash sprite for %s", videoFile.Path)
	}

	return combineImages(images, config), nil
}
