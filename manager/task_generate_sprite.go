package manager

import (
	"github.com/stashapp/stash/ffmpeg"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/models"
	"github.com/stashapp/stash/utils"
	"sync"
)

type GenerateSpriteTask struct {
	Scene models.Scene
}

func (t *GenerateSpriteTask) Start(wg *sync.WaitGroup) {
	if t.doesSpriteExist(t.Scene.Checksum) {
		wg.Done()
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.Paths.FixedPaths.FFProbe, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		wg.Done()
		return
	}

	imagePath := instance.Paths.Scene.GetSpriteImageFilePath(t.Scene.Checksum)
	vttPath := instance.Paths.Scene.GetSpriteVttFilePath(t.Scene.Checksum)
	generator, err := NewSpriteGenerator(*videoFile, imagePath, vttPath, 9, 9)
	if err != nil {
		logger.Errorf("error creating sprite generator: %s", err.Error())
		wg.Done()
		return
	}

	if err := generator.Generate(); err != nil {
		logger.Errorf("error generating sprite: %s", err.Error())
		wg.Done()
		return
	}

	wg.Done()
}

func (t *GenerateSpriteTask) doesSpriteExist(sceneChecksum string) bool {
	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetSpriteImageFilePath(sceneChecksum))
	vttExists, _ := utils.FileExists(instance.Paths.Scene.GetSpriteVttFilePath(sceneChecksum))
	return imageExists && vttExists
}