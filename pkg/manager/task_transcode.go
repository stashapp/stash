package manager

import (
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"os"
	"sync"
)

type GenerateTranscodeTask struct {
	Scene models.Scene
}

func (t *GenerateTranscodeTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	videoCodec := t.Scene.VideoCodec.String
	if ffmpeg.IsValidCodec(videoCodec) {
		return
	}

	hasTranscode, _ := HasTranscode(&t.Scene)
	if hasTranscode {
		return
	}

	logger.Infof("[transcode] <%s> scene has codec %s", t.Scene.Checksum, t.Scene.VideoCodec.String)

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("[transcode] error reading video file: %s", err.Error())
		return
	}

	outputPath := instance.Paths.Generated.GetTmpPath(t.Scene.Checksum + ".mp4")
	options := ffmpeg.TranscodeOptions{
		OutputPath: outputPath,
	}
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)
	encoder.Transcode(*videoFile, options)
	if err := os.Rename(outputPath, instance.Paths.Scene.GetTranscodePath(t.Scene.Checksum)); err != nil {
		logger.Errorf("[transcode] error generating transcode: %s", err.Error())
		return
	}
	logger.Debugf("[transcode] <%s> created transcode: %s", t.Scene.Checksum, outputPath)
	return
}
