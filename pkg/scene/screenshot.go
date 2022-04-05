package scene

import (
	"path/filepath"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"

	// needed to decode other image formats
	_ "image/gif"
	_ "image/png"
)

type screenshotter interface {
	Screenshot(probeResult ffmpeg.VideoFile, options ffmpeg.ScreenshotOptions) error
}

func makeScreenshot(encoder screenshotter, probeResult ffmpeg.VideoFile, outputPath string, quality int, width int, time float64) {
	options := ffmpeg.ScreenshotOptions{
		OutputPath: outputPath,
		Quality:    quality,
		Time:       time,
		Width:      width,
	}

	if err := fsutil.EnsureDirAll(filepath.Dir(outputPath)); err != nil {
		logger.Warnf("[encoder] failure to generate screenshot: %v", err)
		return
	}
	if err := encoder.Screenshot(probeResult, options); err != nil {
		logger.Warnf("[encoder] failure to generate screenshot: %v", err)
	}
}

type ScreenshotSetter interface {
	SetScreenshot(scene *models.Scene, imageData []byte) error
}

type PathsScreenshotSetter struct {
	Paths *paths.Paths
}

func (ss *PathsScreenshotSetter) SetScreenshot(scene *models.Scene, imageData []byte) error {
	return SetScreenshot(ss.Paths, scene.ID, imageData)
}

func writeImage(path string, imageData []byte) error {
	return fsutil.WriteFile(path, imageData)
}

func SetScreenshot(paths *paths.Paths, sceneID int, imageData []byte) error {
	normalPath := paths.Scene.GetCoverPath(sceneID)

	return writeImage(normalPath, imageData)
}
