package manager

import (
	"github.com/remeh/sizedwaitgroup"

	"context"
	"database/sql"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GeneratePhashTask struct {
	Scene               models.Scene
	fileNamingAlgorithm models.HashAlgorithm
	txnManager          models.TransactionManager
}

func (t *GeneratePhashTask) Start(wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()

	if !t.shouldGenerate() {
		return
	}

	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	imagePath := instance.Paths.Scene.GetSpriteImageFilePath(sceneHash)
	generator, err := NewPhashGenerator(imagePath)

	if err != nil {
		logger.Errorf("error creating phash generator: %s", err.Error())
		return
	}
	hash, err := generator.Generate()
	if err != nil {
		logger.Errorf("error generating phash: %s", err.Error())
		return
	}

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()
		hashValue := sql.NullInt64{Int64: int64(*hash), Valid: true}
		scenePartial := models.ScenePartial{
			ID:    t.Scene.ID,
			Phash: &hashValue,
		}
		_, err := qb.Update(scenePartial)
		return err
	}); err != nil {
		logger.Error(err.Error())
	}
}

func (t *GeneratePhashTask) doesSpriteExist(sceneChecksum string) bool {
	if sceneChecksum == "" {
		return false
	}

	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetSpriteImageFilePath(sceneChecksum))
	return imageExists
}

func (t *GeneratePhashTask) shouldGenerate() bool {
	if !t.Scene.Phash.Valid {
		sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
		return t.doesSpriteExist(sceneHash)
	}
	return false
}
