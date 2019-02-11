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
	defer wg.Done()

	if t.doesSpriteExist(t.Scene.Checksum) {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.StaticPaths.FFProbe, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	imagePath := instance.Paths.Scene.GetSpriteImageFilePath(t.Scene.Checksum)
	vttPath := instance.Paths.Scene.GetSpriteVttFilePath(t.Scene.Checksum)
	generator, err := NewSpriteGenerator(*videoFile, imagePath, vttPath, 9, 9)
	if err != nil {
		logger.Errorf("error creating sprite generator: %s", err.Error())
		return
	}

	if err := generator.Generate(); err != nil {
		logger.Errorf("error generating sprite: %s", err.Error())
		return
	}
}

func (t *GenerateSpriteTask) doesSpriteExist(sceneChecksum string) bool {
	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetSpriteImageFilePath(sceneChecksum))
	vttExists, _ := utils.FileExists(instance.Paths.Scene.GetSpriteVttFilePath(sceneChecksum))
	return imageExists && vttExists
}