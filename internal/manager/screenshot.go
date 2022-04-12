package manager

import (
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
)

func makeScreenshot(probeResult ffmpeg.VideoFile, outputPath string, quality int, width int, time float64) {
	encoder := instance.FFMPEG
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
