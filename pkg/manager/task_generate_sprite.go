package manager

import (
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GenerateSpriteTask struct {
	Scene     models.Scene
	Overwrite bool
}

func (t *GenerateSpriteTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.doesSpriteExist(t.Scene.Checksum) && !t.Overwrite {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
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
	generator.Overwrite = t.Overwrite

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
