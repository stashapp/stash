package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateInteractiveHeatmapSpeedTask struct {
	Scene               models.Scene
	Overwrite           bool
	fileNamingAlgorithm models.HashAlgorithm
	TxnManager          Repository
}

func (t *GenerateInteractiveHeatmapSpeedTask) GetDescription() string {
	return fmt.Sprintf("Generating heatmap and speed for %s", t.Scene.Path())
}

func (t *GenerateInteractiveHeatmapSpeedTask) Start(ctx context.Context) {
	if !t.shouldGenerate() {
		return
	}

	videoChecksum := t.Scene.GetHash(t.fileNamingAlgorithm)
	funscriptPath := video.GetFunscriptPath(t.Scene.Path())
	heatmapPath := instance.Paths.Scene.GetInteractiveHeatmapPath(videoChecksum)

	generator := NewInteractiveHeatmapSpeedGenerator(funscriptPath, heatmapPath)

	err := generator.Generate()

	if err != nil {
		logger.Errorf("error generating heatmap: %s", err.Error())
		return
	}

	median := generator.InteractiveSpeed

	if err := t.TxnManager.WithTxn(ctx, func(ctx context.Context) error {
		primaryFile := t.Scene.PrimaryFile()
		primaryFile.InteractiveSpeed = &median
		qb := t.TxnManager.File
		return qb.Update(ctx, primaryFile)
	}); err != nil {
		logger.Error(err.Error())
	}

}

func (t *GenerateInteractiveHeatmapSpeedTask) shouldGenerate() bool {
	primaryFile := t.Scene.PrimaryFile()
	if primaryFile == nil || !primaryFile.Interactive {
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
