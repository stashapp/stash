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
	useMD5    bool
}

func (t *GenerateSpriteTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if !t.Overwrite && !t.required() {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	sceneHash := t.Scene.GetHash(t.useMD5)
	imagePath := instance.Paths.Scene.GetSpriteImageFilePath(sceneHash)
	vttPath := instance.Paths.Scene.GetSpriteVttFilePath(sceneHash)
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

// required returns true if the sprite needs to be generated
func (t GenerateSpriteTask) required() bool {
	sceneHash := t.Scene.GetHash(t.useMD5)
	return !t.doesSpriteExist(sceneHash)
}

func (t *GenerateSpriteTask) doesSpriteExist(sceneChecksum string) bool {
	if sceneChecksum == "" {
		return false
	}

	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetSpriteImageFilePath(sceneChecksum))
	vttExists, _ := utils.FileExists(instance.Paths.Scene.GetSpriteVttFilePath(sceneChecksum))
	return imageExists && vttExists
}
