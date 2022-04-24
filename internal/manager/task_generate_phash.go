package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/hash/videophash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GeneratePhashTask struct {
	Scene               models.Scene
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
	txnManager          models.Repository
}

func (t *GeneratePhashTask) GetDescription() string {
	return fmt.Sprintf("Generating phash for %s", t.Scene.Path)
}

func (t *GeneratePhashTask) Start(ctx context.Context) {
	if !t.shouldGenerate() {
		return
	}

	ffprobe := instance.FFProbe
	videoFile, err := ffprobe.NewVideoFile(t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	hash, err := videophash.Generate(instance.FFMPEG, videoFile)
	if err != nil {
		logger.Errorf("error generating phash: %s", err.Error())
		logErrorOutput(err)
		return
	}

	if err := t.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		qb := t.txnManager.Scene
		hashValue := int64(*hash)
		scenePartial := models.ScenePartial{
			Phash: models.NewOptionalInt64(hashValue),
		}
		_, err := qb.UpdatePartial(ctx, t.Scene.ID, scenePartial)
		return err
	}); err != nil {
		logger.Error(err.Error())
	}
}

func (t *GeneratePhashTask) shouldGenerate() bool {
	return t.Overwrite || t.Scene.Phash == nil
}
