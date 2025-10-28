package manager

import (
	"context"
	"errors"
	"fmt"
	"image"
	"math"

	"github.com/disintegration/imaging"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/scene/generate"
)

type SpriteGenerator struct {
	Info *generatorInfo

	VideoChecksum   string
	ImageOutputPath string
	VTTOutputPath   string
	Rows            int
	Columns         int
	SlowSeek        bool // use alternate seek function, very slow!

	Overwrite bool

	g *generate.Generator
}

func NewSpriteGenerator(videoFile ffmpeg.VideoFile, videoChecksum string, imageOutputPath string, vttOutputPath string, rows int, cols int) (*SpriteGenerator, error) {
	exists, err := fsutil.FileExists(videoFile.Path)
	if !exists {
		return nil, err
	}
	slowSeek := false
	chunkCount := rows * cols

	// For files with small duration / low frame count  try to seek using frame number intead of seconds
	if videoFile.VideoStreamDuration < 5 || (0 < videoFile.FrameCount && videoFile.FrameCount <= int64(chunkCount)) { // some files can have FrameCount == 0, only use SlowSeek  if duration < 5
		if videoFile.VideoStreamDuration <= 0 {
			s := fmt.Sprintf("video %s: duration(%.3f)/frame count(%d) invalid, skipping sprite creation", videoFile.Path, videoFile.VideoStreamDuration, videoFile.FrameCount)
			return nil, errors.New(s)
		}
		logger.Warnf("[generator] video %s too short (%.3fs, %d frames), using frame seeking", videoFile.Path, videoFile.VideoStreamDuration, videoFile.FrameCount)
		slowSeek = true
		// do an actual frame count of the file ( number of frames = read frames)
		ffprobe := GetInstance().FFProbe
		fc, err := ffprobe.GetReadFrameCount(videoFile.Path)
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
		g: &generate.Generator{
			Encoder:      instance.FFMpeg,
			FFMpegConfig: instance.Config,
			LockManager:  instance.ReadLockManager,
			ScenePaths:   instance.Paths.Scene,
		},
	}, nil
}

func (g *SpriteGenerator) Generate() error {
	if err := g.generateSpriteImage(); err != nil {
		return err
	}
	if err := g.generateSpriteVTT(); err != nil {
		return err
	}
	return nil
}

func (g *SpriteGenerator) generateSpriteImage() error {
	if !g.Overwrite && g.imageExists() {
		return nil
	}

	var images []image.Image

	if !g.SlowSeek {
		logger.Infof("[generator] generating sprite image for %s", g.Info.VideoFile.Path)
		// generate `ChunkCount` thumbnails
		stepSize := g.Info.VideoFile.VideoStreamDuration / float64(g.Info.ChunkCount)

		for i := 0; i < g.Info.ChunkCount; i++ {
			time := float64(i) * stepSize

			img, err := g.g.SpriteScreenshot(context.TODO(), g.Info.VideoFile.Path, time)
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

			img, err := g.g.SpriteScreenshotSlow(context.TODO(), g.Info.VideoFile.Path, int(frame))
			if err != nil {
				return err
			}
			images = append(images, img)
		}

	}

	if len(images) == 0 {
		return fmt.Errorf("images slice is empty, failed to generate sprite images for %s", g.Info.VideoFile.Path)
	}

	return imaging.Save(g.g.CombineSpriteImages(images), g.ImageOutputPath)
}

func (g *SpriteGenerator) generateSpriteVTT() error {
	if !g.Overwrite && g.vttExists() {
		return nil
	}
	logger.Infof("[generator] generating sprite vtt for %s", g.Info.VideoFile.Path)

	var stepSize float64
	if !g.SlowSeek {
		stepSize = float64(g.Info.NthFrame) / g.Info.FrameRate
	} else {
		// for files with a low framecount (<ChunkCount) g.Info.NthFrame can be zero
		// so recalculate from scratch
		stepSize = float64(g.Info.VideoFile.FrameCount-1) / float64(g.Info.ChunkCount)
		stepSize /= g.Info.FrameRate
	}

	return g.g.SpriteVTT(context.TODO(), g.VTTOutputPath, g.ImageOutputPath, stepSize)
}

func (g *SpriteGenerator) imageExists() bool {
	exists, _ := fsutil.FileExists(g.ImageOutputPath)
	return exists
}

func (g *SpriteGenerator) vttExists() bool {
	exists, _ := fsutil.FileExists(g.VTTOutputPath)
	return exists
}
