package manager

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type PhashGenerator struct {
	Info *GeneratorInfo

	VideoChecksum string
	Columns       int
	Rows          int
}

func NewPhashGenerator(videoFile ffmpeg.VideoFile, checksum string) (*PhashGenerator, error) {
	exists, err := fsutil.FileExists(videoFile.Path)
	if !exists {
		return nil, err
	}

	generator, err := newGeneratorInfo(videoFile)
	if err != nil {
		return nil, err
	}

	return &PhashGenerator{
		Info:          generator,
		VideoChecksum: checksum,
		Columns:       5,
		Rows:          5,
	}, nil
}

func (g *PhashGenerator) Generate() (*uint64, error) {
	encoder := instance.FFMPEG

	sprite, err := g.generateSprite(&encoder)
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

func (g *PhashGenerator) generateSprite(encoder *ffmpeg.Encoder) (image.Image, error) {
	logger.Infof("[generator] generating phash sprite for %s", g.Info.VideoFile.Path)

	// Generate sprite image offset by 5% on each end to avoid intro/outros
	chunkCount := g.Columns * g.Rows
	offset := 0.05 * g.Info.VideoFile.Duration
	stepSize := (0.9 * g.Info.VideoFile.Duration) / float64(chunkCount)
	var images []image.Image
	for i := 0; i < chunkCount; i++ {
		time := offset + (float64(i) * stepSize)

		options := ffmpeg.SpriteScreenshotOptions{
			Time:  time,
			Width: 160,
		}
		img, err := encoder.SpriteScreenshot(g.Info.VideoFile, options)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	// Combine all of the thumbnails into a sprite image
	if len(images) == 0 {
		return nil, fmt.Errorf("images slice is empty, failed to generate phash sprite for %s", g.Info.VideoFile.Path)
	}
	width := images[0].Bounds().Size().X
	height := images[0].Bounds().Size().Y
	canvasWidth := width * g.Columns
	canvasHeight := height * g.Rows
	montage := imaging.New(canvasWidth, canvasHeight, color.NRGBA{})
	for index := 0; index < len(images); index++ {
		x := width * (index % g.Columns)
		y := height * int(math.Floor(float64(index)/float64(g.Rows)))
		img := images[index]
		montage = imaging.Paste(montage, img, image.Pt(x, y))
	}

	return montage, nil
}
