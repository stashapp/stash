package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GeneratePreviewTask struct {
	Scene        models.Scene
	ImagePreview bool

	Options models.GeneratePreviewOptionsInput

	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
}

func (t *GeneratePreviewTask) GetDescription() string {
	return fmt.Sprintf("Generating preview for %s", t.Scene.Path)
}

func (t *GeneratePreviewTask) Start(ctx context.Context) {
	videoFilename := t.videoFilename()
	videoChecksum := t.Scene.GetHash(t.fileNamingAlgorithm)
	imageFilename := t.imageFilename()

	if !t.Overwrite && !t.required() {
		return
	}

	ffprobe := instance.FFProbe
	videoFile, err := ffprobe.NewVideoFile(t.Scene.Path, false)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	const generateVideo = true
	generator, err := NewPreviewGenerator(*videoFile, videoChecksum, videoFilename, imageFilename, instance.Paths.Generated.Screenshots, generateVideo, t.ImagePreview, t.Options.PreviewPreset.String())

	if err != nil {
		logger.Errorf("error creating preview generator: %s", err.Error())
		return
	}
	generator.Overwrite = t.Overwrite

	// set the preview generation configuration from the global config
	generator.Info.ChunkCount = *t.Options.PreviewSegments
	generator.Info.ChunkDuration = *t.Options.PreviewSegmentDuration
	generator.Info.ExcludeStart = *t.Options.PreviewExcludeStart
	generator.Info.ExcludeEnd = *t.Options.PreviewExcludeEnd
	generator.Info.Audio = config.GetInstance().GetPreviewAudio()

	if err := generator.Generate(); err != nil {
		logger.Errorf("error generating preview: %s", err.Error())
		return
	}
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

	videoExists, _ := fsutil.FileExists(instance.Paths.Scene.GetStreamPreviewPath(sceneChecksum))
	return videoExists
}

func (t *GeneratePreviewTask) doesImagePreviewExist(sceneChecksum string) bool {
	if sceneChecksum == "" {
		return false
	}

	imageExists, _ := fsutil.FileExists(instance.Paths.Scene.GetStreamPreviewImagePath(sceneChecksum))
	return imageExists
}

func (t *GeneratePreviewTask) videoFilename() string {
	return t.Scene.GetHash(t.fileNamingAlgorithm) + ".mp4"
}

func (t *GeneratePreviewTask) imageFilename() string {
	return t.Scene.GetHash(t.fileNamingAlgorithm) + ".webp"
}
