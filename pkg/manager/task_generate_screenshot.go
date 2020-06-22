package manager

import (
	"context"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateScreenshotTask struct {
	Scene        models.Scene
	ScreenshotAt *float64
}

func (t *GenerateScreenshotTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	scenePath := t.Scene.Path
	probeResult, err := ffmpeg.NewVideoFile(instance.FFProbePath, scenePath)

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

	checksum := t.Scene.Checksum
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	// we'll generate the screenshot, grab the generated data and set it
	// in the database. We'll use SetSceneScreenshot to set the data
	// which also generates the thumbnail

	logger.Debugf("Creating screenshot for %s", scenePath)
	makeScreenshot(*probeResult, normalPath, 2, probeResult.Width, at)

	f, err := os.Open(normalPath)
	if err != nil {
		logger.Errorf("Error reading screenshot: %s", err.Error())
		return
	}
	defer f.Close()

	coverImageData, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Errorf("Error reading screenshot: %s", err.Error())
		return
	}

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	qb := models.NewSceneQueryBuilder()
	updatedTime := time.Now()
	updatedScene := models.ScenePartial{
		ID:        t.Scene.ID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	if err := SetSceneScreenshot(t.Scene.Checksum, coverImageData); err != nil {
		logger.Errorf("Error writing screenshot: %s", err.Error())
		tx.Rollback()
		return
	}

	// update the scene cover table
	if err := qb.UpdateSceneCover(t.Scene.ID, coverImageData, tx); err != nil {
		logger.Errorf("Error setting screenshot: %s", err.Error())
		tx.Rollback()
		return
	}

	// update the scene with the update date
	_, err = qb.Update(updatedScene, tx)
	if err != nil {
		logger.Errorf("Error updating scene: %s", err.Error())
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Error setting screenshot: %s", err.Error())
		return
	}
}
