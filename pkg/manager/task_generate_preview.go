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
	"github.com/stashapp/stash/pkg/utils"
)

type GeneratePreviewTaskOperation int

const (
	PreviewTaskOpAll GeneratePreviewTaskOperation = iota
	PreviewTaskOpDefaultScreenshot
	PreviewTaskOpScreenshot
)

type GeneratePreviewTask struct {
	Scene        models.Scene
	Operation    GeneratePreviewTaskOperation
	ScreenshotAt float64
}

func (t *GeneratePreviewTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	videoFilename := t.videoFilename()
	imageFilename := t.imageFilename()

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	generator, err := NewPreviewGenerator(*videoFile, videoFilename, imageFilename, instance.Paths.Generated.Screenshots)
	if err != nil {
		logger.Errorf("error creating preview generator: %s", err.Error())
		return
	}

	switch t.Operation {
	case PreviewTaskOpAll:
		if !t.doesPreviewExist(t.Scene.Checksum) {
			err = generator.Generate()
		}
	case PreviewTaskOpDefaultScreenshot:
		err = generator.GenerateDefaultImage()
		if err == nil {
			err = t.makeScreenshots()
		}
	case PreviewTaskOpScreenshot:
		err = generator.GenerateImageAt(t.ScreenshotAt)
		if err == nil {
			err = t.makeScreenshots()
		}
	}

	if err != nil {
		logger.Errorf("error generating preview: %s", err.Error())
		return
	}
}

func (t *GeneratePreviewTask) makeScreenshots() error {
	checksum := t.Scene.Checksum
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	probeResult, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)

	if err != nil {
		return err
	}

	at := t.ScreenshotAt
	if t.Operation == PreviewTaskOpDefaultScreenshot {
		at = float64(probeResult.Duration) * 0.2
	}

	// we'll generate the screenshot, grab the generated data and set it
	// in the database. We'll use SetSceneScreenshot to set the data
	// which also generates the thumbnail
	makeScreenshot(*probeResult, normalPath, 2, probeResult.Width, at)

	f, err := os.Open(normalPath)
	if err != nil {
		return err
	}
	defer f.Close()

	coverImageData, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	qb := models.NewSceneQueryBuilder()
	updatedTime := time.Now()
	updatedScene := models.ScenePartial{
		ID:        t.Scene.ID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	updatedScene.Cover = &coverImageData
	err = SetSceneScreenshot(t.Scene.Checksum, coverImageData)
	_, err = qb.Update(updatedScene, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (t *GeneratePreviewTask) doesPreviewExist(sceneChecksum string) bool {
	videoExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewPath(sceneChecksum))
	imageExists, _ := utils.FileExists(instance.Paths.Scene.GetStreamPreviewImagePath(sceneChecksum))
	return videoExists && imageExists
}

func (t *GeneratePreviewTask) videoFilename() string {
	return t.Scene.Checksum + ".mp4"
}

func (t *GeneratePreviewTask) imageFilename() string {
	return t.Scene.Checksum + ".webp"
}
