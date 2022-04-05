package manager

import (
	"context"
	"database/sql"
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
	TxnManager          models.TransactionManager
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

	median := sql.NullInt64{
		Int64: generator.InteractiveSpeed,
		Valid: true,
	}

	var s *models.Scene

	if err := t.TxnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		var err error
		s, err = r.Scene().FindByPath(t.Scene.Path)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	if err := t.TxnManager.WithTxn(ctx, func(r models.Repository) error {
		qb := r.Scene()
		scenePartial := models.ScenePartial{
			ID:               s.ID,
			InteractiveSpeed: &median,
		}
		_, err := qb.Update(scenePartial)
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
