package manager

import (
	"os"
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateTranscodeTask struct {
	Scene     models.Scene
	Overwrite bool
	useMD5    bool
}

func (t *GenerateTranscodeTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	hasTranscode := HasTranscode(&t.Scene, t.useMD5)
	if !t.Overwrite && hasTranscode {
		return
	}

	var container ffmpeg.Container

	if t.Scene.Format.Valid {
		container = ffmpeg.Container(t.Scene.Format.String)
	} else { // container isn't in the DB
		// shouldn't happen unless user hasn't scanned after updating to PR#384+ version
		tmpVideoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
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

	if ffmpeg.IsStreamable(videoCodec, audioCodec, container) {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("[transcode] error reading video file: %s", err.Error())
		return
	}

	sceneHash := t.Scene.GetHash(t.useMD5)
	outputPath := instance.Paths.Generated.GetTmpPath(sceneHash + ".mp4")
	transcodeSize := config.GetMaxTranscodeSize()
	options := ffmpeg.TranscodeOptions{
		OutputPath:       outputPath,
		MaxTranscodeSize: transcodeSize,
	}
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

	if videoCodec == ffmpeg.H264 { // for non supported h264 files stream copy the video part
		if audioCodec == ffmpeg.MissingUnsupported {
			encoder.CopyVideo(*videoFile, options)
		} else {
			encoder.TranscodeAudio(*videoFile, options)
		}
	} else {
		if audioCodec == ffmpeg.MissingUnsupported {
			//ffmpeg fails if it trys to transcode an unsupported audio codec
			encoder.TranscodeVideo(*videoFile, options)
		} else {
			encoder.Transcode(*videoFile, options)
		}
	}

	if err := os.Rename(outputPath, instance.Paths.Scene.GetTranscodePath(sceneHash)); err != nil {
		logger.Errorf("[transcode] error generating transcode: %s", err.Error())
		return
	}

	logger.Debugf("[transcode] <%s> created transcode: %s", sceneHash, outputPath)
	return
}

// return true if transcode is needed
// used only when counting files to generate, doesn't affect the actual transcode generation
// if container is missing from DB it is treated as non supported in order not to delay the user
func (t *GenerateTranscodeTask) isTranscodeNeeded() bool {

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

	hasTranscode := HasTranscode(&t.Scene, t.useMD5)
	if !t.Overwrite && hasTranscode {
		return false
	}
	return true
}
