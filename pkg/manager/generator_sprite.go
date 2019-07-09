package manager

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/disintegration/imaging"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
)

type SpriteGenerator struct {
	Info *GeneratorInfo

	ImageOutputPath string
	VTTOutputPath   string
	Rows            int
	Columns         int
}

func NewSpriteGenerator(videoFile ffmpeg.VideoFile, imageOutputPath string, vttOutputPath string, rows int, cols int) (*SpriteGenerator, error) {
	exists, err := utils.FileExists(videoFile.Path)
	if !exists {
		return nil, err
	}
	generator, err := newGeneratorInfo(videoFile)
	if err != nil {
		return nil, err
	}
	generator.ChunkCount = rows * cols
	if err := generator.configure(); err != nil {
		return nil, err
	}

	return &SpriteGenerator{
		Info:            generator,
		ImageOutputPath: imageOutputPath,
		VTTOutputPath:   vttOutputPath,
		Rows:            rows,
		Columns:         cols,
	}, nil
}

func (g *SpriteGenerator) Generate() error {
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

	if err := g.generateSpriteImage(&encoder); err != nil {
		return err
	}
	if err := g.generateSpriteVTT(&encoder); err != nil {
		return err
	}
	return nil
}

func (g *SpriteGenerator) generateSpriteImage(encoder *ffmpeg.Encoder) error {
	if g.imageExists() {
		return nil
	}
	logger.Infof("[generator] generating sprite image for %s", g.Info.VideoFile.Path)

	// Create `this.chunkCount` thumbnails in the tmp directory
	stepSize := int(g.Info.VideoFile.Duration / float64(g.Info.ChunkCount))
	for i := 0; i < g.Info.ChunkCount; i++ {
		time := i * stepSize
		num := fmt.Sprintf("%.3d", i)
		filename := "thumbnail" + num + ".jpg"

		options := ffmpeg.ScreenshotOptions{
			OutputPath: instance.Paths.Generated.GetTmpPath(filename),
			Time:       float64(time),
			Width:      160,
		}
		encoder.Screenshot(g.Info.VideoFile, options)
	}

	// Combine all of the thumbnails into a sprite image
	globPath := filepath.Join(instance.Paths.Generated.Tmp, "thumbnail*.jpg")
	imagePaths, _ := doublestar.Glob(globPath)
	utils.NaturalSort(imagePaths)
	var images []image.Image
	for _, imagePath := range imagePaths {
		img, err := imaging.Open(imagePath)
		if err != nil {
			return err
		}
		images = append(images, img)
	}

	if len(images) == 0 {
		return fmt.Errorf("images slice is empty, failed to generate sprite images for %s", g.Info.VideoFile.Path)
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

	return imaging.Save(montage, g.ImageOutputPath)
}

func (g *SpriteGenerator) generateSpriteVTT(encoder *ffmpeg.Encoder) error {
	if g.vttExists() {
		return nil
	}
	logger.Infof("[generator] generating sprite vtt for %s", g.Info.VideoFile.Path)

	spriteImage, err := imaging.Open(g.ImageOutputPath)
	if err != nil {
		return err
	}
	spriteImageName := filepath.Base(g.ImageOutputPath)
	width := spriteImage.Bounds().Size().X / g.Columns
	height := spriteImage.Bounds().Size().Y / g.Rows

	stepSize := float64(g.Info.NthFrame) / g.Info.FrameRate

	vttLines := []string{"WEBVTT", ""}
	for index := 0; index < g.Info.ChunkCount; index++ {
		x := width * (index % g.Columns)
		y := height * int(math.Floor(float64(index)/float64(g.Rows)))
		startTime := utils.GetVTTTime(float64(index) * stepSize)
		endTime := utils.GetVTTTime(float64(index+1) * stepSize)

		vttLines = append(vttLines, startTime+" --> "+endTime)
		vttLines = append(vttLines, fmt.Sprintf("%s#xywh=%d,%d,%d,%d", spriteImageName, x, y, width, height))
		vttLines = append(vttLines, "")
	}
	vtt := strings.Join(vttLines, "\n")

	return ioutil.WriteFile(g.VTTOutputPath, []byte(vtt), 0644)
}

func (g *SpriteGenerator) imageExists() bool {
	exists, _ := utils.FileExists(g.ImageOutputPath)
	return exists
}

func (g *SpriteGenerator) vttExists() bool {
	exists, _ := utils.FileExists(g.VTTOutputPath)
	return exists
}
