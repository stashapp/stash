package manager

import (
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GeneratePreviewTask struct {
	Scene        models.Scene
	ImagePreview bool

	Options models.GeneratePreviewOptionsInput

	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
}

func (t *GeneratePreviewTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	videoFilename := t.videoFilename()
	imageFilename := t.imageFilename()
	if !t.Overwrite && !t.required() {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	const generateVideo = true
	generator, err := NewPreviewGenerator(*videoFile, videoFilename, imageFilename, instance.Paths.Generated.Screenshots, generateVideo, t.ImagePreview, t.Options.PreviewPreset.String())
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

	videoExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewPath(sceneChecksum))
	return videoExists
}

func (t *GeneratePreviewTask) doesImagePreviewExist(sceneChecksum string) bool {
	if sceneChecksum == "" {
		return false
	}

	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewImagePath(sceneChecksum))
	return imageExists
}

func (t *GeneratePreviewTask) videoFilename() string {
	return t.Scene.GetHash(t.fileNamingAlgorithm) + ".mp4"
}

func (t *GeneratePreviewTask) imageFilename() string {
	return t.Scene.GetHash(t.fileNamingAlgorithm) + ".webp"
}
