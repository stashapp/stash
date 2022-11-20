package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene/generate"
)

type GenerateScreenshotTask struct {
	Scene        models.Scene
	ScreenshotAt *float64
	txnManager   Repository
}

func (t *GenerateScreenshotTask) Start(ctx context.Context) {
	scenePath := t.Scene.Path

	videoFile := t.Scene.Files.Primary()
	if videoFile == nil {
		return
	}

	var at float64
	if t.ScreenshotAt == nil {
		at = float64(videoFile.Duration) * 0.2
	} else {
		at = *t.ScreenshotAt
	}

	// we'll generate the screenshot, grab the generated data and set it
	// in the database.

	logger.Debugf("Creating screenshot for %s", scenePath)

	g := generate.Generator{
		Encoder:      instance.FFMPEG,
		FFMpegConfig: instance.Config,
		LockManager:  instance.ReadLockManager,
		ScenePaths:   instance.Paths.Scene,
		Overwrite:    true,
	}

	coverImageData, err := g.Screenshot(context.TODO(), videoFile.Path, videoFile.Width, videoFile.Duration, generate.ScreenshotOptions{
		At: &at,
	})
	if err != nil {
		logger.Errorf("Error generating screenshot: %v", err)
		logErrorOutput(err)
		return
	}

	if err := t.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		qb := t.txnManager.Scene
		updatedScene := models.NewScenePartial()

		// update the scene cover table
		if err := qb.UpdateCover(ctx, t.Scene.ID, coverImageData); err != nil {
			return fmt.Errorf("error setting screenshot: %v", err)
		}

		// update the scene with the update date
		_, err = qb.UpdatePartial(ctx, t.Scene.ID, updatedScene)
		if err != nil {
			return fmt.Errorf("error updating scene: %v", err)
		}

		return nil
	}); err != nil && ctx.Err() == nil {
		logger.Error(err.Error())
	}
}
