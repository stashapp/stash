package manager

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

type SpriteGenerator struct {
	Info *GeneratorInfo

	VideoChecksum   string
	ImageOutputPath string
	VTTOutputPath   string
	Rows            int
	Columns         int
	SlowSeek        bool // use alternate seek function, very slow!

	Overwrite bool
}

func NewSpriteGenerator(videoFile ffmpeg.VideoFile, videoChecksum string, imageOutputPath string, vttOutputPath string, rows int, cols int) (*SpriteGenerator, error) {
	exists, err := fsutil.FileExists(videoFile.Path)
	if !exists {
		return nil, err
	}
	slowSeek := false
	chunkCount := rows * cols

	// For files with small duration / low frame count  try to seek using frame number intead of seconds
	if videoFile.Duration < 5 || (0 < videoFile.FrameCount && videoFile.FrameCount <= int64(chunkCount)) { // some files can have FrameCount == 0, only use SlowSeek  if duration < 5
		if videoFile.Duration <= 0 {
			s := fmt.Sprintf("video %s: duration(%.3f)/frame count(%d) invalid, skipping sprite creation", videoFile.Path, videoFile.Duration, videoFile.FrameCount)
			return nil, errors.New(s)
		}
		logger.Warnf("[generator] video %s too short (%.3fs, %d frames), using frame seeking", videoFile.Path, videoFile.Duration, videoFile.FrameCount)
		slowSeek = true
		// do an actual frame count of the file ( number of frames = read frames)
		ffprobe := GetInstance().FFProbe
		fc, err := ffprobe.GetReadFrameCount(&videoFile)
		if err == nil {
			if fc != videoFile.FrameCount {
				logger.Warnf("[generator] updating framecount (%d) for %s with read frames count (%d)", videoFile.FrameCount, videoFile.Path, fc)
				videoFile.FrameCount = fc
			}
		}
	}

	generator, err := newGeneratorInfo(videoFile)
	if err != nil {
		return nil, err
	}
	generator.ChunkCount = chunkCount
	if err := generator.configure(); err != nil {
		return nil, err
	}

	return &SpriteGenerator{
		Info:            generator,
		VideoChecksum:   videoChecksum,
		ImageOutputPath: imageOutputPath,
		VTTOutputPath:   vttOutputPath,
		Rows:            rows,
		SlowSeek:        slowSeek,
		Columns:         cols,
	}, nil
}

func (g *SpriteGenerator) Generate() error {
	encoder := instance.FFMPEG

	if err := g.generateSpriteImage(&encoder); err != nil {
		return err
	}
	if err := g.generateSpriteVTT(&encoder); err != nil {
		return err
	}
	return nil
}

func (g *SpriteGenerator) generateSpriteImage(encoder *ffmpeg.Encoder) error {
	if !g.Overwrite && g.imageExists() {
		return nil
	}

	var images []image.Image

	if !g.SlowSeek {
		logger.Infof("[generator] generating sprite image for %s", g.Info.VideoFile.Path)
		// generate `ChunkCount` thumbnails
		stepSize := g.Info.VideoFile.Duration / float64(g.Info.ChunkCount)

		for i := 0; i < g.Info.ChunkCount; i++ {
			time := float64(i) * stepSize

			options := ffmpeg.SpriteScreenshotOptions{
				Time:  time,
				Width: 160,
			}

			img, err := encoder.SpriteScreenshot(g.Info.VideoFile, options)

			if err != nil {
				return err
			}
			images = append(images, img)
		}
	} else {
		logger.Infof("[generator] generating sprite image for %s (%d frames)", g.Info.VideoFile.Path, g.Info.VideoFile.FrameCount)

		stepFrame := float64(g.Info.VideoFile.FrameCount-1) / float64(g.Info.ChunkCount)

		for i := 0; i < g.Info.ChunkCount; i++ {
			// generate exactly `ChunkCount` thumbnails, using duplicate frames if needed
			frame := math.Round(float64(i) * stepFrame)
			if frame >= math.MaxInt || frame <= math.MinInt {
				return errors.New("invalid frame number conversion")
			}
			options := ffmpeg.SpriteScreenshotOptions{
				Frame: int(frame),
				Width: 160,
			}
			img, err := encoder.SpriteScreenshotSlow(g.Info.VideoFile, options)
			if err != nil {
				return err
			}
			images = append(images, img)
		}

	}

	if len(images) == 0 {
		return fmt.Errorf("images slice is empty, failed to generate sprite images for %s", g.Info.VideoFile.Path)
	}
	// Combine all of the thumbnails into a sprite image
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
	if !g.Overwrite && g.vttExists() {
		return nil
	}
	logger.Infof("[generator] generating sprite vtt for %s", g.Info.VideoFile.Path)

	spriteImage, err := os.Open(g.ImageOutputPath)
	if err != nil {
		return err
	}
	defer spriteImage.Close()
	spriteImageName := filepath.Base(g.ImageOutputPath)
	image, _, err := image.DecodeConfig(spriteImage)
	if err != nil {
		return err
	}
	width := image.Width / g.Columns
	height := image.Height / g.Rows

	var stepSize float64
	if !g.SlowSeek {
		stepSize = float64(g.Info.NthFrame) / g.Info.FrameRate
	} else {
		// for files with a low framecount (<ChunkCount) g.Info.NthFrame can be zero
		// so recalculate from scratch
		stepSize = float64(g.Info.VideoFile.FrameCount-1) / float64(g.Info.ChunkCount)
		stepSize /= g.Info.FrameRate
	}

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

	return os.WriteFile(g.VTTOutputPath, []byte(vtt), 0644)
}

func (g *SpriteGenerator) imageExists() bool {
	exists, _ := fsutil.FileExists(g.ImageOutputPath)
	return exists
}

func (g *SpriteGenerator) vttExists() bool {
	exists, _ := fsutil.FileExists(g.VTTOutputPath)
	return exists
}
