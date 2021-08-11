package manager

import (
	"github.com/remeh/sizedwaitgroup"

	"context"
	"database/sql"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GeneratePhashTask struct {
	Scene               models.Scene
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
	txnManager          models.TransactionManager
}

func (t *GeneratePhashTask) Start(wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()

	if !t.shouldGenerate() {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path, false)
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
