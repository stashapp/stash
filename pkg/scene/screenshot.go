package scene

import (
	"os"

	"github.com/stashapp/stash/pkg/ffmpeg"
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

	if err := encoder.Screenshot(probeResult, options); err != nil {
		logger.Warnf("[encoder] failure to generate screenshot: %v", err)
	}
}

type ScreenshotSetter interface {
	SetScreenshot(scene *models.Scene, imageData []byte) error
}

type PathsScreenshotSetter struct {
	Paths               *paths.Paths
	FileNamingAlgorithm models.HashAlgorithm
}

func (ss *PathsScreenshotSetter) SetScreenshot(scene *models.Scene, imageData []byte) error {
	checksum := scene.GetHash(ss.FileNamingAlgorithm)
	return SetScreenshot(ss.Paths, checksum, imageData)
}

func writeImage(path string, imageData []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(imageData)
	return err
}

func SetScreenshot(paths *paths.Paths, checksum string, imageData []byte) error {
	normalPath := paths.Scene.GetScreenshotPath(checksum)

	return writeImage(normalPath, imageData)
}
