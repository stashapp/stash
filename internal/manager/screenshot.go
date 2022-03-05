package manager

import (
	"github.com/stashapp/stash/internal/video/encoder"
	"github.com/stashapp/stash/pkg/ffmpeg2"
	"github.com/stashapp/stash/pkg/logger"
)

// TODO - replace with scene.makeScreenshot
func makeScreenshot(ff ffmpeg2.FFMpeg, input string, outputPath string, quality int, width int, time float64) {
	options := encoder.ScreenshotOptions{
		OutputPath: outputPath,
		Quality:    quality,
		Time:       time,
		Width:      width,
	}

	if err := encoder.Screenshot(ff, input, options); err != nil {
		logger.Warnf("[encoder] failure to generate screenshot: %v", err)
	}
}
