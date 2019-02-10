package manager

import (
	"github.com/stashapp/stash/ffmpeg"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/models"
	"github.com/stashapp/stash/utils"
	"sync"
)

type GeneratePreviewTask struct {
	Scene models.Scene
}

func (t *GeneratePreviewTask) Start(wg *sync.WaitGroup) {
	videoFilename := t.videoFilename()
	imageFilename := t.imageFilename()
	if t.doesPreviewExist(videoFilename, imageFilename) {
		wg.Done()
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.Paths.FixedPaths.FFProbe, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		wg.Done()
		return
	}

	generator, err := NewPreviewGenerator(*videoFile, videoFilename, imageFilename, instance.Paths.Generated.Screenshots)
	if err != nil {
		logger.Errorf("error creating preview generator: %s", err.Error())
		wg.Done()
		return
	}

	if err := generator.Generate(); err != nil {
		logger.Errorf("error generating preview: %s", err.Error())
		wg.Done()
		return
	}

	wg.Done()
}

func (t *GeneratePreviewTask) doesPreviewExist(videoFilename string, imageFilename string) bool {
	videoExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewPath(videoFilename))
	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewImagePath(imageFilename))
	return videoExists && imageExists
}

func (t *GeneratePreviewTask) videoFilename() string {
	return t.Scene.Checksum + ".mp4"
}

func (t *GeneratePreviewTask) imageFilename() string {
	return t.Scene.Checksum + ".webp"
}