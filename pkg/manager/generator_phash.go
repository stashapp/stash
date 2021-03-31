package manager

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"sort"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
	"github.com/fvbommel/sortorder"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

type PhashGenerator struct {
	Info *GeneratorInfo

	VideoChecksum string
	Columns       int
	Rows          int
}

func NewPhashGenerator(videoFile ffmpeg.VideoFile, checksum string) (*PhashGenerator, error) {
	exists, err := utils.FileExists(videoFile.Path)
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
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

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
	for i := 0; i < chunkCount; i++ {
		time := offset + (float64(i) * stepSize)
		num := fmt.Sprintf("%.3d", i)
		filename := "phash_" + g.VideoChecksum + "_" + num + ".bmp"

		options := ffmpeg.ScreenshotOptions{
			OutputPath: instance.Paths.Generated.GetTmpPath(filename),
			Time:       time,
			Width:      160,
		}
		if err := encoder.Screenshot(g.Info.VideoFile, options); err != nil {
			return nil, err
		}
	}

	// Combine all of the thumbnails into a sprite image
	pattern := fmt.Sprintf("phash_%s_.+\\.bmp$", g.VideoChecksum)
	imagePaths, err := utils.MatchEntries(instance.Paths.Generated.Tmp, pattern)
	if err != nil {
		return nil, err
	}
	sort.Sort(sortorder.Natural(imagePaths))
	var images []image.Image
	for _, imagePath := range imagePaths {
		img, err := imaging.Open(imagePath)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

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

	for _, imagePath := range imagePaths {
		os.Remove(imagePath)
	}

	return montage, nil
}
