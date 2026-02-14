package manager

import (
	"context"
	"errors"
	"fmt"
	"image"
	"math"

	"github.com/disintegration/imaging"

	"github.com/stashapp/stash/internal/manager/config"
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
	SpriteInterval  float64
	SlowSeek        bool // use alternate seek function, very slow!

	Overwrite bool

	g *generate.Generator
}

func NewSpriteGenerator(videoFile ffmpeg.VideoFile, videoChecksum string, imageOutputPath string, vttOutputPath string) (*SpriteGenerator, error) {
	exists, err := fsutil.FileExists(videoFile.Path)
	if !exists {
		return nil, err
	}
	slowSeek := false

	config := config.GetInstance()

	spriteCount := int64(81)
	minimumSpriteCount := 0
	spriteInterval := float64(0)
	if config.GetUseCustomSpriteGeneration() {
		spriteInterval = float64(config.GetSpriteInterval())
		minimumSpriteCount = config.GetMinimumSprites()
		spriteCount = int64(math.Ceil(videoFile.VideoStreamDuration / float64(spriteInterval)))
	} else {
		spriteInterval = videoFile.VideoStreamDuration / float64(spriteCount)
	}

	// For files with small duration / low frame count  try to seek using frame number intead of seconds
	if (spriteCount < 1) || (videoFile.VideoStreamDuration < spriteInterval*float64(minimumSpriteCount)) || (0 < videoFile.FrameCount && videoFile.FrameCount <= spriteCount) { // some files can have FrameCount == 0
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
	if spriteCount > int64(minimumSpriteCount) {
		generator.ChunkCount = int(spriteCount)
	} else {
		generator.ChunkCount = minimumSpriteCount
	}
	if err := generator.configure(); err != nil {
		return nil, err
	}

	return &SpriteGenerator{
		Info:            generator,
		VideoChecksum:   videoChecksum,
		ImageOutputPath: imageOutputPath,
		VTTOutputPath:   vttOutputPath,
		SpriteInterval:  spriteInterval,
		SlowSeek:        slowSeek,
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

		time := 0.
		for time <= g.Info.VideoFile.VideoStreamDuration {
			img, err := g.g.SpriteScreenshot(context.TODO(), g.Info.VideoFile.Path, time)
			if err != nil {
				return err
			}
			images = append(images, img)
			time += float64(g.SpriteInterval)
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
		stepSize = float64(g.SpriteInterval)
	} else {
		// for files with a low framecount (<ChunkCount) g.Info.NthFrame can be zero
		// so recalculate from scratch
		stepSize = float64(g.Info.VideoFile.FrameCount-1) / float64(g.Info.ChunkCount)
		stepSize /= g.Info.FrameRate
	}

	return g.g.SpriteVTT(context.TODO(), g.VTTOutputPath, g.ImageOutputPath, stepSize, g.Info.VideoFile.VideoStreamDuration, g.Info.ChunkCount)
}

func (g *SpriteGenerator) imageExists() bool {
	exists, _ := fsutil.FileExists(g.ImageOutputPath)
	return exists
}

func (g *SpriteGenerator) vttExists() bool {
	exists, _ := fsutil.FileExists(g.VTTOutputPath)
	return exists
}
