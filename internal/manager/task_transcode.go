package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene/generate"
)

type GenerateTranscodeTask struct {
	Scene               models.Scene
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm

	// is true, generate even if video is browser-supported
	Force bool

	g *generate.Generator
}

func (t *GenerateTranscodeTask) GetDescription() string {
	return fmt.Sprintf("Generating transcode for %s", t.Scene.Path)
}

func (t *GenerateTranscodeTask) Start(ctx context.Context) {
	hasTranscode := HasTranscode(&t.Scene, t.fileNamingAlgorithm)
	if !t.Overwrite && hasTranscode {
		return
	}

	f := t.Scene.Files.Primary()

	ffprobe := instance.FFProbe
	var container ffmpeg.Container

	var err error
	container, err = GetVideoFileContainer(f)
	if err != nil {
		logger.Errorf("[transcode] error getting scene container: %s", err.Error())
		return
	}

	var videoCodec string

	if f.VideoCodec != "" {
		videoCodec = f.VideoCodec
	}

	audioCodec := ffmpeg.MissingUnsupported
	if f.AudioCodec != "" {
		audioCodec = ffmpeg.ProbeAudioCodec(f.AudioCodec)
	}

	if !t.Force && ffmpeg.IsStreamable(videoCodec, audioCodec, container) == nil {
		return
	}

	// TODO - move transcode generation logic elsewhere

	videoFile, err := ffprobe.NewVideoFile(f.Path)
	if err != nil {
		logger.Errorf("[transcode] error reading video file: %s", err.Error())
		return
	}

	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	transcodeSize := config.GetInstance().GetMaxTranscodeSize()

	w, h := videoFile.TranscodeScale(transcodeSize.GetMaxResolution())

	// if scale is being set, then we can't use stream copy
	scaleSet := w == 0 && h == 0

	if scaleSet && videoCodec == ffmpeg.H264 { // for non supported h264 files stream copy the video part
		if audioCodec == ffmpeg.MissingUnsupported {
			err = t.g.TranscodeCopyVideo(ctx, videoFile.Path, sceneHash)
		} else {
			err = t.g.TranscodeAudio(ctx, videoFile.Path, sceneHash)
		}
	} else {
		options := generate.TranscodeOptions{
			Width:  w,
			Height: h,
		}

		if audioCodec == ffmpeg.MissingUnsupported {
			// ffmpeg fails if it tries to transcode an unsupported audio codec
			err = t.g.TranscodeVideo(ctx, videoFile.Path, sceneHash, options)
		} else {
			err = t.g.Transcode(ctx, videoFile.Path, sceneHash, options)
		}
	}

	if err != nil {
		logger.Errorf("[transcode] error generating transcode: %v", err)
		return
	}
}

// return true if transcode is needed
// used only when counting files to generate, doesn't affect the actual transcode generation
// if container is missing from DB it is treated as non supported in order not to delay the user
func (t *GenerateTranscodeTask) required() bool {
	f := t.Scene.Files.Primary()
	if f == nil {
		return false
	}

	hasTranscode := HasTranscode(&t.Scene, t.fileNamingAlgorithm)
	if !t.Overwrite && hasTranscode {
		return false
	}

	if t.Force {
		return true
	}

	var videoCodec string
	if f.VideoCodec != "" {
		videoCodec = f.VideoCodec
	}
	container := ""
	audioCodec := ffmpeg.MissingUnsupported
	if f.AudioCodec != "" {
		audioCodec = ffmpeg.ProbeAudioCodec(f.AudioCodec)
	}

	if f.Format != "" {
		container = f.Format
	}

	if ffmpeg.IsStreamable(videoCodec, audioCodec, ffmpeg.Container(container)) == nil {
		return false
	}

	return true
}
