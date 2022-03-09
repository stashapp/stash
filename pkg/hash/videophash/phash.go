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
	"github.com/stashapp/stash/pkg/ffmpeg2"
	"github.com/stashapp/stash/pkg/logger"
)

const (
	screenshotSize = 160
	columns        = 5
	rows           = 5
)

// isValidFloat64 ensures the given value is a valid number (not NaN) which is not equal to 0
func isValidFloat64(value float64) bool {
	return !math.IsNaN(value) && value != 0
}

func Generate(encoder ffmpeg2.FFMpeg, videoFile *ffmpeg.VideoFile) (*uint64, error) {
	sprite, err := generateSprite(encoder, videoFile)
	if err != nil {
		return nil, err
	}

	hash, err := goimagehash.PerceptionHash(sprite)
	if err != nil {
		return nil, err
	}
	hashValue := hash.GetHash()
	return &hashValue, nil
}

func generateSpriteScreenshot(encoder ffmpeg2.FFMpeg, input string, t float64) (image.Image, error) {
	options := ffmpeg2.ScreenshotOptions{
		Width:      screenshotSize,
		OutputPath: "-",
		OutputType: ffmpeg2.ScreenshotOutputTypeBMP,
	}

	args := ffmpeg2.ScreenshotTime(input, t, options)
	data, err := ffmpeg2.GenerateOutput(context.Background(), encoder, args)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func combineImages(images []image.Image) image.Image {
	width := images[0].Bounds().Size().X
	height := images[0].Bounds().Size().Y
	canvasWidth := width * columns
	canvasHeight := height * rows
	montage := imaging.New(canvasWidth, canvasHeight, color.NRGBA{})
	for index := 0; index < len(images); index++ {
		x := width * (index % columns)
		y := height * int(math.Floor(float64(index)/float64(rows)))
		img := images[index]
		montage = imaging.Paste(montage, img, image.Pt(x, y))
	}

	return montage
}

func generateSprite(encoder ffmpeg2.FFMpeg, videoFile *ffmpeg.VideoFile) (image.Image, error) {
	logger.Infof("[generator] generating phash sprite for %s", videoFile.Path)

	// Generate sprite image offset by 5% on each end to avoid intro/outros
	chunkCount := columns * rows
	offset := 0.05 * videoFile.Duration
	stepSize := (0.9 * videoFile.Duration) / float64(chunkCount)
	var images []image.Image
	for i := 0; i < chunkCount; i++ {
		time := offset + (float64(i) * stepSize)

		img, err := generateSpriteScreenshot(encoder, videoFile.Path, time)
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	// Combine all of the thumbnails into a sprite image
	if len(images) == 0 {
		return nil, fmt.Errorf("images slice is empty, failed to generate phash sprite for %s", videoFile.Path)
	}

	return combineImages(images), nil
}
