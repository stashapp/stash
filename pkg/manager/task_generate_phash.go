package manager

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GeneratePhashTask struct {
	Scene               models.Scene
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
	txnManager          models.TransactionManager
}

func (t *GeneratePhashTask) GetDescription() string {
	return fmt.Sprintf("Generating phash for %s", t.Scene.Path)
}

func (t *GeneratePhashTask) Start(ctx context.Context) {
	if !t.shouldGenerate() {
		return
	}

	ffprobe := instance.FFProbe
	videoFile, err := ffprobe.NewVideoFile(t.Scene.Path, false)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	generator, err := NewPhashGenerator(*videoFile, sceneHash)

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

func (t *GeneratePhashTask) shouldGenerate() bool {
	return t.Overwrite || !t.Scene.Phash.Valid
}
