package generate

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/utils"
)

const (
	spriteScreenshotWidth = 160

	spriteRows   = 9
	spriteCols   = 9
	spriteChunks = spriteRows * spriteCols
)

func (g Generator) SpriteScreenshot(ctx context.Context, input string, seconds float64) (image.Image, error) {
	lockCtx := g.LockManager.ReadLock(ctx, input)
	defer lockCtx.Cancel()

	ssOptions := transcoder.ScreenshotOptions{
		OutputPath: "-",
		OutputType: transcoder.ScreenshotOutputTypeBMP,
		Width:      spriteScreenshotWidth,
	}

	args := transcoder.ScreenshotTime(input, seconds, ssOptions)

	return g.generateImage(lockCtx, args)
}

func (g Generator) SpriteScreenshotSlow(ctx context.Context, input string, frame int) (image.Image, error) {
	lockCtx := g.LockManager.ReadLock(ctx, input)
	defer lockCtx.Cancel()

	ssOptions := transcoder.ScreenshotOptions{
		OutputPath: "-",
		OutputType: transcoder.ScreenshotOutputTypeBMP,
		Width:      spriteScreenshotWidth,
	}

	args := transcoder.ScreenshotFrame(input, frame, ssOptions)

	return g.generateImage(lockCtx, args)
}

func (g Generator) generateImage(lockCtx *fsutil.LockContext, args ffmpeg.Args) (image.Image, error) {
	out, err := g.generateOutput(lockCtx, args)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(out))
	if err != nil {
		return nil, fmt.Errorf("decoding image from ffmpeg: %w", err)
	}

	return img, nil
}

func (g Generator) CombineSpriteImages(images []image.Image) image.Image {
	// Combine all of the thumbnails into a sprite image
	width := images[0].Bounds().Size().X
	height := images[0].Bounds().Size().Y
	canvasWidth := width * spriteCols
	canvasHeight := height * spriteRows
	montage := imaging.New(canvasWidth, canvasHeight, color.NRGBA{})
	for index := 0; index < len(images); index++ {
		x := width * (index % spriteCols)
		y := height * int(math.Floor(float64(index)/float64(spriteRows)))
		img := images[index]
		montage = imaging.Paste(montage, img, image.Pt(x, y))
	}

	return montage
}

func (g Generator) SpriteVTT(ctx context.Context, output string, spritePath string, stepSize float64) error {
	lockCtx := g.LockManager.ReadLock(ctx, spritePath)
	defer lockCtx.Cancel()

	return g.generateFile(lockCtx, g.ScenePaths, vttPattern, output, g.spriteVTT(spritePath, stepSize))
}

func (g Generator) spriteVTT(spritePath string, stepSize float64) generateFn {
	return func(lockCtx *fsutil.LockContext, tmpFn string) error {
		spriteImage, err := os.Open(spritePath)
		if err != nil {
			return err
		}
		defer spriteImage.Close()
		spriteImageName := filepath.Base(spritePath)
		image, _, err := image.DecodeConfig(spriteImage)
		if err != nil {
			return err
		}
		width := image.Width / spriteCols
		height := image.Height / spriteRows

		vttLines := []string{"WEBVTT", ""}
		for index := 0; index < spriteChunks; index++ {
			x := width * (index % spriteCols)
			y := height * int(math.Floor(float64(index)/float64(spriteRows)))
			startTime := utils.GetVTTTime(float64(index) * stepSize)
			endTime := utils.GetVTTTime(float64(index+1) * stepSize)

			vttLines = append(vttLines, startTime+" --> "+endTime)
			vttLines = append(vttLines, fmt.Sprintf("%s#xywh=%d,%d,%d,%d", spriteImageName, x, y, width, height))
			vttLines = append(vttLines, "")
		}
		vtt := strings.Join(vttLines, "\n")

		return os.WriteFile(tmpFn, []byte(vtt), 0644)
	}
}

// TODO - move all sprite generation code here
// WIP
// func (g Generator) Sprite(ctx context.Context, videoFile *ffmpeg.VideoFile, hash string) error {
// 	input := videoFile.Path
// 	if err := g.generateSpriteImage(ctx, videoFile, hash); err != nil {
// 		return fmt.Errorf("generating sprite image for %s: %w", input, err)
// 	}

// 	output := g.ScenePaths.GetSpriteVttFilePath(hash)
// 	if !g.Overwrite {
// 		if exists, _ := fsutil.FileExists(output); exists {
// 			return nil
// 		}
// 	}

// 	if err := g.generateFile(ctx, g.ScenePaths, vttPattern, output, g.spriteVtt(input, screenshotOptions{
// 		Time:    at,
// 		Quality: screenshotQuality,
// 		// default Width is video width
// 	})); err != nil {
// 		return err
// 	}

// 	logger.Debug("created screenshot: ", output)

// 	return nil
// }

// func (g Generator) generateSpriteImage(ctx context.Context, videoFile *ffmpeg.VideoFile, hash string) error {
// 	output := g.ScenePaths.GetSpriteImageFilePath(hash)
// 	if !g.Overwrite {
// 		if exists, _ := fsutil.FileExists(output); exists {
// 			return nil
// 		}
// 	}

// 	var images []image.Image
// 	var err error
// 	if options.VideoDuration > 0 {
// 		images, err = g.generateSprites(ctx, input, options.VideoDuration)
// 	} else {
// 		images, err = g.generateSpritesSlow(ctx, input, options.FrameCount)
// 	}

// 	if len(images) == 0 {
// 		return errors.New("images slice is empty")
// 	}

// 	montage, err := g.combineSpriteImages(images)
// 	if err != nil {
// 		return err
// 	}

// 	if err := imaging.Save(montage, output); err != nil {
// 		return err
// 	}

// 	logger.Debug("created sprite image: ", output)

// 	return nil
// }

// func useSlowSeek(videoFile *ffmpeg.VideoFile) (bool, error) {
// 	// For files with small duration / low frame count  try to seek using frame number intead of seconds
// 	// some files can have FrameCount == 0, only use SlowSeek if duration < 5
// 	if videoFile.Duration < 5 || (videoFile.FrameCount > 0 && videoFile.FrameCount <= int64(spriteChunks)) {
// 		if videoFile.Duration <= 0 {
// 			return false, fmt.Errorf("duration(%.3f)/frame count(%d) invalid", videoFile.Duration, videoFile.FrameCount)
// 		}

// 		logger.Warnf("[generator] video %s too short (%.3fs, %d frames), using frame seeking", videoFile.Path, videoFile.Duration, videoFile.FrameCount)
// 		return true, nil
// 	}
// }

// func (g Generator) combineSpriteImages(images []image.Image) (image.Image, error) {
// 	// Combine all of the thumbnails into a sprite image
// 	width := images[0].Bounds().Size().X
// 	height := images[0].Bounds().Size().Y
// 	canvasWidth := width * spriteCols
// 	canvasHeight := height * spriteRows
// 	montage := imaging.New(canvasWidth, canvasHeight, color.NRGBA{})
// 	for index := 0; index < len(images); index++ {
// 		x := width * (index % spriteCols)
// 		y := height * int(math.Floor(float64(index)/float64(spriteRows)))
// 		img := images[index]
// 		montage = imaging.Paste(montage, img, image.Pt(x, y))
// 	}

// 	return montage, nil
// }

// func (g Generator) generateSprites(ctx context.Context, input string, videoDuration float64) ([]image.Image, error) {
// 	logger.Infof("[generator] generating sprite image for %s", input)
// 	// generate `ChunkCount` thumbnails
// 	stepSize := videoDuration / float64(spriteChunks)

// 	var images []image.Image
// 	for i := 0; i < spriteChunks; i++ {
// 		time := float64(i) * stepSize

// 		img, err := g.spriteScreenshot(ctx, input, time)
// 		if err != nil {
// 			return nil, err
// 		}
// 		images = append(images, img)
// 	}

// 	return images, nil
// }

// func (g Generator) generateSpritesSlow(ctx context.Context, input string, frameCount int) ([]image.Image, error) {
// 	logger.Infof("[generator] generating sprite image for %s (%d frames)", input, frameCount)

// 	stepFrame := float64(frameCount-1) / float64(spriteChunks)

// 	var images []image.Image
// 	for i := 0; i < spriteChunks; i++ {
// 		// generate exactly `ChunkCount` thumbnails, using duplicate frames if needed
// 		frame := math.Round(float64(i) * stepFrame)
// 		if frame >= math.MaxInt || frame <= math.MinInt {
// 			return nil, errors.New("invalid frame number conversion")
// 		}

// 		img, err := g.spriteScreenshotSlow(ctx, input, int(frame))
// 		if err != nil {
// 			return nil, err
// 		}
// 		images = append(images, img)
// 	}

// 	return images, nil
// }

// func (g Generator) spriteScreenshot(ctx context.Context, input string, seconds float64) (image.Image, error) {
// 	ssOptions := transcoder.ScreenshotOptions{
// 		OutputPath: "-",
// 		OutputType: transcoder.ScreenshotOutputTypeBMP,
// 		Width:      spriteScreenshotWidth,
// 	}

// 	args := transcoder.ScreenshotTime(input, seconds, ssOptions)

// 	return g.generateImage(ctx, args)
// }

// func (g Generator) spriteScreenshotSlow(ctx context.Context, input string, frame int) (image.Image, error) {
// 	ssOptions := transcoder.ScreenshotOptions{
// 		OutputPath: "-",
// 		OutputType: transcoder.ScreenshotOutputTypeBMP,
// 		Width:      spriteScreenshotWidth,
// 	}

// 	args := transcoder.ScreenshotFrame(input, frame, ssOptions)

// 	return g.generateImage(ctx, args)
// }

// func (g Generator) spriteVTT(videoFile ffmpeg.VideoFile, spriteImagePath string, slowSeek bool) generateFn {
// 	return func(ctx context.Context, tmpFn string) error {
// 		logger.Infof("[generator] generating sprite vtt for %s", input)

// 		spriteImage, err := os.Open(spriteImagePath)
// 		if err != nil {
// 			return err
// 		}
// 		defer spriteImage.Close()
// 		spriteImageName := filepath.Base(spriteImagePath)
// 		image, _, err := image.DecodeConfig(spriteImage)
// 		if err != nil {
// 			return err
// 		}
// 		width := image.Width / spriteCols
// 		height := image.Height / spriteRows

// 		var stepSize float64
// 		if !slowSeek {
// 			nthFrame = g.NumberOfFrames / g.ChunkCount
// 			stepSize = float64(g.Info.NthFrame) / g.Info.FrameRate
// 		} else {
// 			// for files with a low framecount (<ChunkCount) g.Info.NthFrame can be zero
// 			// so recalculate from scratch
// 			stepSize = float64(videoFile.FrameCount-1) / float64(spriteChunks)
// 			stepSize /= g.Info.FrameRate
// 		}

// 		vttLines := []string{"WEBVTT", ""}
// 		for index := 0; index < spriteChunks; index++ {
// 			x := width * (index % spriteCols)
// 			y := height * int(math.Floor(float64(index)/float64(spriteRows)))
// 			startTime := utils.GetVTTTime(float64(index) * stepSize)
// 			endTime := utils.GetVTTTime(float64(index+1) * stepSize)

// 			vttLines = append(vttLines, startTime+" --> "+endTime)
// 			vttLines = append(vttLines, fmt.Sprintf("%s#xywh=%d,%d,%d,%d", spriteImageName, x, y, width, height))
// 			vttLines = append(vttLines, "")
// 		}
// 		vtt := strings.Join(vttLines, "\n")

// 		return os.WriteFile(tmpFn, []byte(vtt), 0644)
// 	}
// }
