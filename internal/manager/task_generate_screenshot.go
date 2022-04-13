package manager

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scene/generate"
)

type GenerateScreenshotTask struct {
	Scene               models.Scene
	ScreenshotAt        *float64
	fileNamingAlgorithm models.HashAlgorithm
	txnManager          models.Repository
}

func (t *GenerateScreenshotTask) Start(ctx context.Context) {
	scenePath := t.Scene.Path
	ffprobe := instance.FFProbe
	probeResult, err := ffprobe.NewVideoFile(scenePath)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	var at float64
	if t.ScreenshotAt == nil {
		at = float64(probeResult.Duration) * 0.2
	} else {
		at = *t.ScreenshotAt
	}

	checksum := t.Scene.GetHash(t.fileNamingAlgorithm)
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	// we'll generate the screenshot, grab the generated data and set it
	// in the database. We'll use SetSceneScreenshot to set the data
	// which also generates the thumbnail

	logger.Debugf("Creating screenshot for %s", scenePath)

	g := generate.Generator{
		Encoder:     instance.FFMPEG,
		LockManager: instance.ReadLockManager,
		ScenePaths:  instance.Paths.Scene,
		Overwrite:   true,
	}

	if err := g.Screenshot(context.TODO(), probeResult.Path, checksum, probeResult.Width, probeResult.Duration, generate.ScreenshotOptions{
		At: &at,
	}); err != nil {
		logger.Errorf("Error generating screenshot: %v", err)
		logErrorOutput(err)
		return
	}

	f, err := os.Open(normalPath)
	if err != nil {
		logger.Errorf("Error reading screenshot: %s", err.Error())
		return
	}
	defer f.Close()

	coverImageData, err := io.ReadAll(f)
	if err != nil {
		logger.Errorf("Error reading screenshot: %s", err.Error())
		return
	}

	if err := t.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		qb := t.txnManager.Scene
		updatedTime := time.Now()
		updatedScene := models.ScenePartial{
			ID:        t.Scene.ID,
			UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
		}

		if err := scene.SetScreenshot(instance.Paths, checksum, coverImageData); err != nil {
			return fmt.Errorf("error writing screenshot: %v", err)
		}

		// update the scene cover table
		if err := qb.UpdateCover(ctx, t.Scene.ID, coverImageData); err != nil {
			return fmt.Errorf("error setting screenshot: %v", err)
		}

		// update the scene with the update date
		_, err = qb.Update(ctx, updatedScene)
		if err != nil {
			return fmt.Errorf("error updating scene: %v", err)
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}
