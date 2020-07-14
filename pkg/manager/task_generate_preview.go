package manager

import (
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GeneratePreviewTask struct {
	Scene         models.Scene
	ImagePreview  bool
	PreviewPreset string
	useMD5        bool
}

func (t *GeneratePreviewTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	videoFilename := t.videoFilename()
	imageFilename := t.imageFilename()
	if !t.required() {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	sceneHash := t.Scene.GetHash(t.useMD5)
	videoExists := t.doesVideoPreviewExist(sceneHash)
	generator, err := NewPreviewGenerator(*videoFile, videoFilename, imageFilename, instance.Paths.Generated.Screenshots, !videoExists, t.ImagePreview, t.PreviewPreset)
	if err != nil {
		logger.Errorf("error creating preview generator: %s", err.Error())
		return
	}

	if err := generator.Generate(); err != nil {
		logger.Errorf("error generating preview: %s", err.Error())
		return
	}
}

func (t GeneratePreviewTask) required() bool {
	sceneHash := t.Scene.GetHash(t.useMD5)
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
	return t.Scene.GetHash(t.useMD5) + ".mp4"
}

func (t *GeneratePreviewTask) imageFilename() string {
	return t.Scene.GetHash(t.useMD5) + ".webp"
}
