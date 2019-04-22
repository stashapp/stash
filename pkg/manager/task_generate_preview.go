package manager

import (
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"sync"
)

type GeneratePreviewTask struct {
	Scene models.Scene
}

func (t *GeneratePreviewTask) Start(wg *sync.WaitGroup,previewsCh chan<- struct{},errorCh chan<- struct{}) {
	defer wg.Done()

	videoFilename := t.videoFilename()
	imageFilename := t.imageFilename()
	if t.doesPreviewExist(t.Scene.Checksum) {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		errorCh <- struct{}{}
		return
	}

	generator, err := NewPreviewGenerator(*videoFile, videoFilename, imageFilename, instance.Paths.Generated.Screenshots)
	if err != nil {
		logger.Errorf("error creating preview generator: %s", err.Error())
		errorCh <- struct{}{}
		return
	}

	if err := generator.Generate(); err != nil {
		logger.Errorf("error generating preview: %s", err.Error())
		errorCh <- struct{}{}
		return
	}
	previewsCh <- struct{}{} // since we got here that means we generated a new preview
}

func (t *GeneratePreviewTask) doesPreviewExist(sceneChecksum string) bool {
	videoExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewPath(sceneChecksum))
	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewImagePath(sceneChecksum))
	return videoExists && imageExists
}

func (t *GeneratePreviewTask) videoFilename() string {
	return t.Scene.Checksum + ".mp4"
}

func (t *GeneratePreviewTask) imageFilename() string {
	return t.Scene.Checksum + ".webp"
}
