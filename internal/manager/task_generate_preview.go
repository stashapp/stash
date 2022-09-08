package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene/generate"
)

type GeneratePreviewTask struct {
	Scene        models.Scene
	ImagePreview bool

	Options generate.PreviewOptions

	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm

	generator *generate.Generator
}

func (t *GeneratePreviewTask) GetDescription() string {
	return fmt.Sprintf("Generating preview for %s", t.Scene.Path)
}

func (t *GeneratePreviewTask) Start(ctx context.Context) {
	if !t.Overwrite && !t.required() {
		return
	}

	ffprobe := instance.FFProbe
	videoFile, err := ffprobe.NewVideoFile(t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %v", err)
		return
	}

	videoChecksum := t.Scene.GetHash(t.fileNamingAlgorithm)

	if err := t.generateVideo(videoChecksum, videoFile.Duration); err != nil {
		logger.Errorf("error generating preview: %v", err)
		logErrorOutput(err)
		return
	}

	if t.ImagePreview {
		if err := t.generateWebp(videoChecksum); err != nil {
			logger.Errorf("error generating preview webp: %v", err)
			logErrorOutput(err)
		}
	}
}

func (t GeneratePreviewTask) generateVideo(videoChecksum string, videoDuration float64) error {
	videoFilename := t.Scene.Path

	if err := t.generator.PreviewVideo(context.TODO(), videoFilename, videoDuration, videoChecksum, t.Options, true); err != nil {
		logger.Warnf("[generator] failed generating scene preview, trying fallback")
		if err := t.generator.PreviewVideo(context.TODO(), videoFilename, videoDuration, videoChecksum, t.Options, true); err != nil {
			return err
		}
	}

	return nil
}

func (t GeneratePreviewTask) generateWebp(videoChecksum string) error {
	videoFilename := t.Scene.Path
	return t.generator.PreviewWebp(context.TODO(), videoFilename, videoChecksum)
}

func (t GeneratePreviewTask) required() bool {
	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	videoExists := t.doesVideoPreviewExist(sceneHash)
	imageExists := !t.ImagePreview || t.doesImagePreviewExist(sceneHash)
	return !imageExists || !videoExists
}

func (t *GeneratePreviewTask) doesVideoPreviewExist(sceneChecksum string) bool {
	if sceneChecksum == "" {
		return false
	}

	videoExists, _ := fsutil.FileExists(instance.Paths.Scene.GetVideoPreviewPath(sceneChecksum))
	return videoExists
}

func (t *GeneratePreviewTask) doesImagePreviewExist(sceneChecksum string) bool {
	if sceneChecksum == "" {
		return false
	}

	imageExists, _ := fsutil.FileExists(instance.Paths.Scene.GetWebpPreviewPath(sceneChecksum))
	return imageExists
}
