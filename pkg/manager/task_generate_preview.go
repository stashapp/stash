package manager

import (
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GeneratePreviewTask struct {
	Scene         models.Scene
	ImagePreview  bool
	PreviewPreset string
	Overwrite     bool
}

func (t *GeneratePreviewTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	videoFilename := t.videoFilename()
	imageFilename := t.imageFilename()
	videoExists := t.doesVideoPreviewExist(t.Scene.Checksum)
	if !t.Overwrite && ((!t.ImagePreview || t.doesImagePreviewExist(t.Scene.Checksum)) && videoExists) {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	generator, err := NewPreviewGenerator(*videoFile, videoFilename, imageFilename, instance.Paths.Generated.Screenshots, true, t.ImagePreview, t.PreviewPreset)
	if err != nil {
		logger.Errorf("error creating preview generator: %s", err.Error())
		return
	}
	generator.Overwrite = t.Overwrite

	// set the preview generation configuration from the global config
	generator.Info.ChunkCount = config.GetPreviewSegments()
	generator.Info.ChunkDuration = config.GetPreviewSegmentDuration()
	generator.Info.ExcludeStart = config.GetPreviewExcludeStart()
	generator.Info.ExcludeEnd = config.GetPreviewExcludeEnd()

	if err := generator.Generate(); err != nil {
		logger.Errorf("error generating preview: %s", err.Error())
		return
	}
}

func (t *GeneratePreviewTask) doesVideoPreviewExist(sceneChecksum string) bool {
	videoExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewPath(sceneChecksum))
	return videoExists
}

func (t *GeneratePreviewTask) doesImagePreviewExist(sceneChecksum string) bool {
	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewImagePath(sceneChecksum))
	return imageExists
}

func (t *GeneratePreviewTask) videoFilename() string {
	return t.Scene.Checksum + ".mp4"
}

func (t *GeneratePreviewTask) imageFilename() string {
	return t.Scene.Checksum + ".webp"
}
