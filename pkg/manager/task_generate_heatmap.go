package manager

import (
	"context"
	"database/sql"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GenerateHeatmapTask struct {
	Scene               models.Scene
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
	TxnManager          models.TransactionManager
}

func (t *GenerateHeatmapTask) Start(wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()

	videoChecksum := t.Scene.GetHash(t.fileNamingAlgorithm)
	funscriptPath := utils.GetFunscriptPath(t.Scene.Path)
	heatmapPath := instance.Paths.Scene.GetHeatmapPath(videoChecksum)

	if !t.Overwrite && !t.required() {
		return
	}

	generator := NewHeatmapGenerator(funscriptPath, heatmapPath)

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

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		s, err = r.Scene().FindByPath(t.Scene.Path)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
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

func (t *GenerateHeatmapTask) required() bool {
	if !t.Scene.Interactive {
		return false
	}
	sceneHash := t.Scene.GetHash(t.fileNamingAlgorithm)
	return !t.doesHeatmapExist(sceneHash)
}

func (t *GenerateHeatmapTask) doesHeatmapExist(sceneChecksum string) bool {
	if sceneChecksum == "" {
		return false
	}

	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetHeatmapPath(sceneChecksum))
	return imageExists
}
