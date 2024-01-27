package manager

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type generatorInfo struct {
	ChunkCount     int
	FrameRate      float64
	NumberOfFrames int

	// NthFrame used for sprite generation
	NthFrame int

	VideoFile ffmpeg.VideoFile
}

func newGeneratorInfo(videoFile ffmpeg.VideoFile) (*generatorInfo, error) {
	exists, err := fsutil.FileExists(videoFile.Path)
	if !exists {
		logger.Errorf("video file not found")
		return nil, err
	}

	generator := &generatorInfo{VideoFile: videoFile}
	return generator, nil
}

func (g *generatorInfo) calculateFrameRate(videoStream *ffmpeg.FFProbeStream) error {
	var framerate float64
	if g.VideoFile.FrameRate == 0 {
		framerate, _ = strconv.ParseFloat(videoStream.RFrameRate, 64)
	} else {
		framerate = g.VideoFile.FrameRate
	}

	numberOfFrames, _ := strconv.Atoi(videoStream.NbFrames)

	if numberOfFrames == 0 && isValidFloat64(framerate) && g.VideoFile.VideoStreamDuration > 0 { // TODO: test
		numberOfFrames = int(framerate * g.VideoFile.VideoStreamDuration)
	}

	// If we are missing the frame count or frame rate then seek through the file and extract the info with regex
	if numberOfFrames == 0 || !isValidFloat64(framerate) {
		info, err := instance.FFMpeg.CalculateFrameRate(context.TODO(), &g.VideoFile)
		if err != nil {
			logger.Errorf("error calculating frame rate: %v", err)
		} else {
			if numberOfFrames == 0 {
				numberOfFrames = info.NumberOfFrames
			}
			if !isValidFloat64(framerate) {
				framerate = info.FrameRate
			}
		}
	}

	// Something seriously wrong with this file
	if numberOfFrames == 0 || !isValidFloat64(framerate) {
		logger.Errorf(
			"number of frames or framerate is 0.  nb_frames <%s> framerate <%f> duration <%f>",
			videoStream.NbFrames,
			framerate,
			g.VideoFile.VideoStreamDuration,
		)
	}

	g.FrameRate = framerate
	g.NumberOfFrames = numberOfFrames

	return nil
}

// isValidFloat64 ensures the given value is a valid number (not NaN) which is not equal to 0
func isValidFloat64(value float64) bool {
	return !math.IsNaN(value) && value != 0
}

func (g *generatorInfo) configure() error {
	videoStream := g.VideoFile.VideoStream
	if videoStream == nil {
		return fmt.Errorf("missing video stream")
	}

	if err := g.calculateFrameRate(videoStream); err != nil {
		return err
	}

	// #2250 - ensure ChunkCount is valid
	if g.ChunkCount < 1 {
		logger.Warnf("[generator] Segment count (%d) must be > 0. Using 1 instead.", g.ChunkCount)
		g.ChunkCount = 1
	}

	g.NthFrame = g.NumberOfFrames / g.ChunkCount

	return nil
}
