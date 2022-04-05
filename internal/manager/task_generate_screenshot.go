package manager

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateScreenshotTask struct {
	Scene        models.Scene
	ScreenshotAt *float64
	txnManager   models.TransactionManager
}

func (t *GenerateScreenshotTask) Start(ctx context.Context) {
	scenePath := t.Scene.Path
	ffprobe := instance.FFProbe
	probeResult, err := ffprobe.NewVideoFile(scenePath, false)

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

	normalPath := instance.Paths.Scene.GetCoverPath(t.Scene.ID)

	logger.Debugf("Creating screenshot for %s", scenePath)
	pathErr := fsutil.EnsureDirAll(filepath.Dir(normalPath))
	if pathErr != nil {
		logger.Errorf("Cannot ensure path: %v", pathErr)
		return
	}
	makeScreenshot(*probeResult, normalPath, 2, probeResult.Width, at)

	// update the scene with the update date
	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()
		updatedTime := time.Now()
		updatedScene := models.ScenePartial{
			ID:        t.Scene.ID,
			UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
		}

		_, err = qb.Update(updatedScene)
		if err != nil {
			return fmt.Errorf("error updating scene: %v", err)
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}
