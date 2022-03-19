package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateTranscodeTask struct {
	Scene               models.Scene
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm

	// is true, generate even if video is browser-supported
	Force bool
}

func (t *GenerateTranscodeTask) GetDescription() string {
	return fmt.Sprintf("Generating transcode for %s", t.Scene.Path)
}

func (t *GenerateTranscodeTask) Start(ctc context.Context) {
	hasTranscode := HasTranscode(&t.Scene, t.fileNamingAlgorithm)
	if !t.Overwrite && hasTranscode {
		return
	}

	ffprobe := instance.FFProbe
	var container ffmpeg.Container

	if t.Scene.Format.Valid {
		container = ffmpeg.Container(t.Scene.Format.String)
	} else { // container isn't in the DB
		// shouldn't happen unless user hasn't scanned after updating to PR#384+ version
		tmpVideoFile, err := ffprobe.NewVideoFile(t.Scene.Path, false)
		if err != nil {
			logger.Errorf("[transcode] error reading video file: %s", err.Error())
			return
		}

		container = ffmpeg.MatchContainer(tmpVideoFile.Container, t.Scene.Path)
	}

	videoCodec := t.Scene.VideoCodec.String
	audioCodec := ffmpeg.MissingUnsupported
	if t.Scene.AudioCodec.Valid {
		audioCodec = ffmpeg.AudioCodec(t.Scene.AudioCodec.String)
	}

	if !t.Force && ffmpeg.IsStreamable(videoCodec, audioCodec, container) {
		return
	}

	videoFile, err := ffprobe.NewVideoFile(t.Scene.Path, false)
	if err != nil {
		logger.Errorf("[transcode] error reading video file: %s", err.Error())
		return
	}

	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	outputPath := instance.Paths.Generated.GetTmpPath(sceneHash + ".mp4")
	transcodeSize := config.GetInstance().GetMaxTranscodeSize()
	options := ffmpeg.TranscodeOptions{
		OutputPath:       outputPath,
		MaxTranscodeSize: transcodeSize,
	}
	encoder := instance.FFMPEG

	if videoCodec == ffmpeg.H264 { // for non supported h264 files stream copy the video part
		if audioCodec == ffmpeg.MissingUnsupported {
			encoder.CopyVideo(*videoFile, options)
		} else {
			encoder.TranscodeAudio(*videoFile, options)
		}
	} else {
		if audioCodec == ffmpeg.MissingUnsupported {
			// ffmpeg fails if it trys to transcode an unsupported audio codec
			encoder.TranscodeVideo(*videoFile, options)
		} else {
			encoder.Transcode(*videoFile, options)
		}
	}

	if err := fsutil.SafeMove(outputPath, instance.Paths.Scene.GetTranscodePath(sceneHash)); err != nil {
		logger.Errorf("[transcode] error generating transcode: %s", err.Error())
		return
	}

	logger.Debugf("[transcode] <%s> created transcode: %s", sceneHash, outputPath)
}

// return true if transcode is needed
// used only when counting files to generate, doesn't affect the actual transcode generation
// if container is missing from DB it is treated as non supported in order not to delay the user
func (t *GenerateTranscodeTask) isTranscodeNeeded() bool {
	hasTranscode := HasTranscode(&t.Scene, t.fileNamingAlgorithm)
	if !t.Overwrite && hasTranscode {
		return false
	}

	if t.Force {
		return true
	}

	videoCodec := t.Scene.VideoCodec.String
	container := ""
	audioCodec := ffmpeg.MissingUnsupported
	if t.Scene.AudioCodec.Valid {
		audioCodec = ffmpeg.AudioCodec(t.Scene.AudioCodec.String)
	}

	if t.Scene.Format.Valid {
		container = t.Scene.Format.String
	}

	if ffmpeg.IsStreamable(videoCodec, audioCodec, ffmpeg.Container(container)) {
		return false
	}

	return true
}
