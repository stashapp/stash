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
	Config          SpriteGeneratorConfig
	SlowSeek        bool // use alternate seek function, very slow!

	Overwrite bool

	g *generate.Generator
}

// SpriteGeneratorConfig holds configuration for the SpriteGenerator
type SpriteGeneratorConfig struct {
	// MinimumSprites is the minimum number of sprites to generate, even if the video duration is short
	// SpriteInterval will be adjusted accordingly to ensure at least this many sprites are generated.
	// A value of 0 means no minimum, and the generator will use the provided SpriteInterval or
	// calculate it based on the video duration and MaximumSprites
	MinimumSprites int

	// MaximumSprites is the maximum number of sprites to generate, even if the video duration is long
	// SpriteInterval will be adjusted accordingly to ensure no more than this many sprites are generated
	// A value of 0 means no maximum, and the generator will use the provided SpriteInterval or
	// calculate it based on the video duration and MinimumSprites
	MaximumSprites int

	// SpriteInterval is the default interval in seconds between each sprite.
	// If MinimumSprites or MaximumSprites are set, this value will be adjusted accordingly
	// to ensure the desired number of sprites are generated
	// A value of 0 means the generator will calculate the interval based on the video duration and
	// the provided MinimumSprites and MaximumSprites
	SpriteInterval float64

	// SpriteSize is the size in pixels of the longest dimension of each sprite image.
	// The other dimension will be automatically calculated to maintain the aspect ratio of the video
	SpriteSize int
}

const (
	// DefaultSpriteAmount is the default number of sprites to generate if no configuration is provided
	// This corresponds to the legacy behavior of the generator, which generates 81 sprites at equal
	// intervals across the video duration
	DefaultSpriteAmount = 81

	// DefaultSpriteSize is the default size in pixels of the longest dimension of each sprite image
	// if no configuration is provided. This corresponds to the legacy behavior of the generator.
	DefaultSpriteSize = 160
)

var DefaultSpriteGeneratorConfig = SpriteGeneratorConfig{
	MinimumSprites: DefaultSpriteAmount,
	MaximumSprites: DefaultSpriteAmount,
	SpriteInterval: 0,
	SpriteSize:     DefaultSpriteSize,
}

// NewSpriteGenerator creates a new SpriteGenerator for the given video file and configuration
// It calculates the appropriate sprite interval and count based on the video duration and the provided configuration
func NewSpriteGenerator(videoFile ffmpeg.VideoFile, videoChecksum string, imageOutputPath string, vttOutputPath string, config SpriteGeneratorConfig) (*SpriteGenerator, error) {
	exists, err := fsutil.FileExists(videoFile.Path)
	if !exists {
		return nil, err
	}

	if videoFile.VideoStreamDuration <= 0 {
		s := fmt.Sprintf("video %s: duration(%.3f)/frame count(%d) invalid, skipping sprite creation", videoFile.Path, videoFile.VideoStreamDuration, videoFile.FrameCount)
		return nil, errors.New(s)
	}

	config.SpriteInterval = calculateSpriteInterval(videoFile, config)
	chunkCount := int(math.Ceil(videoFile.VideoStreamDuration / config.SpriteInterval))

	// adjust the chunk count to the next highest perfect square, to ensure the sprite image
	// is completely filled (no empty space in the grid) and the grid is as square as possible (minimizing the number of rows/columns)
	gridSize := generate.GetSpriteGridSize(chunkCount)
	newChunkCount := gridSize * gridSize

	if newChunkCount != chunkCount {
		logger.Debugf("[generator] adjusting chunk count from %d to %d to fit a %dx%d grid", chunkCount, newChunkCount, gridSize, gridSize)
		chunkCount = newChunkCount
	}

	if config.SpriteSize <= 0 {
		config.SpriteSize = DefaultSpriteSize
	}

	slowSeek := false

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
		Config:          config,
		SlowSeek:        slowSeek,
		g: &generate.Generator{
			Encoder:      instance.FFMpeg,
			FFMpegConfig: instance.Config,
			LockManager:  instance.ReadLockManager,
			ScenePaths:   instance.Paths.Scene,
		},
	}, nil
}

func calculateSpriteInterval(videoFile ffmpeg.VideoFile, config SpriteGeneratorConfig) float64 {
	// If a custom sprite interval is provided, start with that
	spriteInterval := config.SpriteInterval

	// If no custom interval is provided, calculate the interval based on the
	// video duration and minimum sprite count
	if spriteInterval <= 0 {
		minSprites := config.MinimumSprites
		if minSprites <= 0 {
			panic("invalid configuration: MinimumSprites must be greater than 0 if SpriteInterval is not set")
		}

		logger.Debugf("[generator] calculating sprite interval for video duration %.3fs with minimum sprites %d", videoFile.VideoStreamDuration, minSprites)
		return videoFile.VideoStreamDuration / float64(minSprites)
	}

	// Calculate the number of sprites that would be generated with the provided interval
	spriteCount := int(math.Ceil(videoFile.VideoStreamDuration / spriteInterval))

	// If the calculated sprite count is greater than the maximum, adjust the interval to meet the maximum
	if config.MaximumSprites > 0 && spriteCount > int(config.MaximumSprites) {
		spriteInterval = videoFile.VideoStreamDuration / float64(config.MaximumSprites)
		logger.Debugf("[generator] provided sprite interval %.1fs results in %d sprites, which exceeds the maximum of %d, adjusting interval to %.1fs", config.SpriteInterval, spriteCount, config.MaximumSprites, spriteInterval)
	}

	// If the calculated sprite count is less than the minimum, adjust the interval to meet the minimum
	if config.MinimumSprites > 0 && spriteCount < int(config.MinimumSprites) {
		spriteInterval = videoFile.VideoStreamDuration / float64(config.MinimumSprites)
		logger.Debugf("[generator] provided sprite interval %.1fs results in %d sprites, which is less than the minimum of %d, adjusting interval to %.1fs", config.SpriteInterval, spriteCount, config.MinimumSprites, spriteInterval)
	}

	return spriteInterval
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

	isPortrait := g.Info.VideoFile.Height > g.Info.VideoFile.Width

	if !g.SlowSeek {
		logger.Infof("[generator] generating sprite image for %s", g.Info.VideoFile.Path)
		// generate `ChunkCount` thumbnails
		stepSize := g.Info.VideoFile.VideoStreamDuration / float64(g.Info.ChunkCount)

		for i := 0; i < g.Info.ChunkCount; i++ {
			time := float64(i) * stepSize
			img, err := g.g.SpriteScreenshot(context.TODO(), g.Info.VideoFile.Path, time, g.Config.SpriteSize, isPortrait)
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

			img, err := g.g.SpriteScreenshotSlow(context.TODO(), g.Info.VideoFile.Path, int(frame), g.Config.SpriteSize)
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

	return g.g.SpriteVTT(context.TODO(), g.VTTOutputPath, g.ImageOutputPath, stepSize, g.Info.ChunkCount)
}

func (g *SpriteGenerator) imageExists() bool {
	exists, _ := fsutil.FileExists(g.ImageOutputPath)
	return exists
}

func (g *SpriteGenerator) vttExists() bool {
	exists, _ := fsutil.FileExists(g.VTTOutputPath)
	return exists
}
