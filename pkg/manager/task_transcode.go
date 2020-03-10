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
	Scene models.Scene
}

func (t *GenerateTranscodeTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	hasTranscode, _ := HasTranscode(&t.Scene)
	if hasTranscode {
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
	if ffmpeg.IsValidCodec(videoCodec) && ffmpeg.IsValidCombo(videoCodec, container) {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("[transcode] error reading video file: %s", err.Error())
		return
	}

	outputPath := instance.Paths.Generated.GetTmpPath(t.Scene.Checksum + ".mp4")
	transcodeSize := config.GetMaxTranscodeSize()
	options := ffmpeg.TranscodeOptions{
		OutputPath:       outputPath,
		MaxTranscodeSize: transcodeSize,
	}
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)

	if videoCodec == "h264" { // for non supported h264 files stream copy the video part
		encoder.TranscodeAudio(*videoFile, options)
	} else {
		encoder.Transcode(*videoFile, options)
	}

	if err := os.Rename(outputPath, instance.Paths.Scene.GetTranscodePath(t.Scene.Checksum)); err != nil {
		logger.Errorf("[transcode] error generating transcode: %s", err.Error())
		return
	}

	logger.Debugf("[transcode] <%s> created transcode: %s", t.Scene.Checksum, outputPath)
	return
}

// return true if transcode is needed
// used only when counting files to generate, doesn't affect the actual transcode generation
// if container is missing from DB it is treated as non supported in order not to delay the user
func (t *GenerateTranscodeTask) isTranscodeNeeded() bool {

	videoCodec := t.Scene.VideoCodec.String
	container := ""
	if t.Scene.Format.Valid {
		container = t.Scene.Format.String
	}

	if ffmpeg.IsValidCodec(videoCodec) && ffmpeg.IsValidCombo(videoCodec, ffmpeg.Container(container)) {
		return false
	}

	hasTranscode, _ := HasTranscode(&t.Scene)
	if hasTranscode {
		return false
	}
	return true
}
