package manager

import (
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateScreenshotTask struct {
	Scene        models.Scene
	ScreenshotAt *float64
}

func (t *GenerateScreenshotTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	scenePath := t.Scene.Path
	probeResult, err := ffmpeg.NewVideoFile(instance.FFProbePath, scenePath)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	var at float64
	if t.ScreenshotAt == nil {
		at = float64(probeResult.Duration) * 0.2
	} else {
		at = *t.ScreenshotAt
	}

	checksum := t.Scene.Checksum
	thumbPath := instance.Paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	logger.Debugf("Creating thumbnail for %s", scenePath)
	makeScreenshot(*probeResult, thumbPath, 5, 320, at)

	logger.Debugf("Creating screenshot for %s", scenePath)
	makeScreenshot(*probeResult, normalPath, 2, probeResult.Width, at)
}
