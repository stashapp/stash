package generate

import (
	"context"

	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

const (
	// thumbnailWidth   = 320
	// thumbnailQuality = 5

	screenshotQuality = 2

	screenshotDurationProportion = 0.2
)

type ScreenshotOptions struct {
	At *float64
}

func (g Generator) Screenshot(ctx context.Context, input string, videoWidth int, videoDuration float64, options ScreenshotOptions) ([]byte, error) {
	lockCtx := g.LockManager.ReadLock(ctx, input)
	defer lockCtx.Cancel()

	logger.Infof("Creating screenshot for %s", input)

	at := screenshotDurationProportion * videoDuration
	if options.At != nil {
		at = *options.At
	}

	ret, err := g.generateBytes(lockCtx, g.ScenePaths, jpgPattern, g.screenshot(input, screenshotOptions{
		Time:    at,
		Quality: screenshotQuality,
		// default Width is video width
	}))
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type screenshotOptions struct {
	Time    float64
	Width   int
	Quality int
}

func (g Generator) screenshot(input string, options screenshotOptions) generateFn {
	return func(lockCtx *fsutil.LockContext, tmpFn string) error {
		ssOptions := transcoder.ScreenshotOptions{
			OutputPath: tmpFn,
			OutputType: transcoder.ScreenshotOutputTypeImage2,
			Quality:    options.Quality,
			Width:      options.Width,
		}

		args := transcoder.ScreenshotTime(input, options.Time, ssOptions)

		return g.generate(lockCtx, args)
	}
}
