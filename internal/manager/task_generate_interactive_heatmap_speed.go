package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type GenerateInteractiveHeatmapSpeedTask struct {
	Scene               models.Scene
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
	TxnManager          models.Repository
}

func (t *GenerateInteractiveHeatmapSpeedTask) GetDescription() string {
	return fmt.Sprintf("Generating heatmap and speed for %s", t.Scene.Path)
}

func (t *GenerateInteractiveHeatmapSpeedTask) Start(ctx context.Context) {
	if !t.shouldGenerate() {
		return
	}

	videoChecksum := t.Scene.GetHash(t.fileNamingAlgorithm)
	funscriptPath := scene.GetFunscriptPath(t.Scene.Path)
	heatmapPath := instance.Paths.Scene.GetInteractiveHeatmapPath(videoChecksum)

	generator := NewInteractiveHeatmapSpeedGenerator(funscriptPath, heatmapPath)

	err := generator.Generate()

	if err != nil {
		logger.Errorf("error generating heatmap: %s", err.Error())
		return
	}

	median := generator.InteractiveSpeed

	var s *models.Scene

	if err := t.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		var err error
		s, err = t.TxnManager.Scene.FindByPath(ctx, t.Scene.Path)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	if err := t.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		qb := t.TxnManager.Scene
		scenePartial := models.ScenePartial{
			InteractiveSpeed: models.NewOptionalInt(median),
		}
		_, err := qb.UpdatePartial(ctx, s.ID, scenePartial)
		return err
	}); err != nil {
		logger.Error(err.Error())
	}

}

func (t *GenerateInteractiveHeatmapSpeedTask) shouldGenerate() bool {
	if !t.Scene.Interactive {
		return false
	}
	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	return !t.doesHeatmapExist(sceneHash) || t.Overwrite
}

func (t *GenerateInteractiveHeatmapSpeedTask) doesHeatmapExist(sceneChecksum string) bool {
	if sceneChecksum == "" {
		return false
	}

	imageExists, _ := fsutil.FileExists(instance.Paths.Scene.GetInteractiveHeatmapPath(sceneChecksum))
	return imageExists
}
