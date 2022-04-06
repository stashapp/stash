package scene

import (
	"io/ioutil"
	"path/filepath"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
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

type CoverGetter interface {
	GetCover(sceneID int) ([]byte, error)
}

type CoverSetter interface {
	SetCover(sceneID int, imageData []byte) error
}

type CoverGetterSetter interface {
	CoverGetter
	CoverSetter
}

type FileWriter interface {
	WriteFile(path string, file []byte) error
}

type PathsCover struct {
	Paths      *paths.Paths
	FileWriter FileWriter
}

func (ss *PathsCover) GetCover(sceneID int) ([]byte, error) {
	normalPath := ss.Paths.Scene.GetCoverPath(sceneID)
	// if the file doesn't exist, return nil
	if exists, _ := fsutil.FileExists(normalPath); !exists {
		return nil, nil
	}
	return ioutil.ReadFile(normalPath)
}

func (ss *PathsCover) SetCover(sceneID int, imageData []byte) error {
	normalPath := ss.Paths.Scene.GetCoverPath(sceneID)
	if err := ss.FileWriter.WriteFile(normalPath, imageData); err != nil {
		return err
	}
	return nil
}
