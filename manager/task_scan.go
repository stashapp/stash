package manager

import (
	"context"
	"database/sql"
	"github.com/stashapp/stash/database"
	"github.com/stashapp/stash/ffmpeg"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/models"
	"github.com/stashapp/stash/utils"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type ScanTask struct {
	FilePath string
}

func (t *ScanTask) Start(wg *sync.WaitGroup) {
	if filepath.Ext(t.FilePath) == ".zip" {
		t.scanGallery()
	} else {
		t.scanScene()
	}

	wg.Done()
}

func (t *ScanTask) scanGallery() {
	qb := models.NewGalleryQueryBuilder()
	gallery, _ := qb.FindByPath(t.FilePath)
	if gallery != nil {
		// We already have this item in the database, keep going
		return
	}

	checksum, err := t.calculateChecksum()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)
	gallery, _ = qb.FindByChecksum(checksum, tx)
	if gallery != nil {
		logger.Infof("%s already exists.  Updating path...", t.FilePath)
		gallery.Path = t.FilePath
		_, err = qb.Update(*gallery, tx)
	} else {
		logger.Infof("%s doesn't exist.  Creating new item...", t.FilePath)
		currentTime := time.Now()
		newGallery := models.Gallery{
			Checksum:  checksum,
			Path:      t.FilePath,
			CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		}
		_, err = qb.Create(newGallery, tx)
	}

	if err != nil {
		logger.Error(err.Error())
		_ = tx.Rollback()
	} else if err := tx.Commit(); err != nil {
		logger.Error(err.Error())
	}
}

func (t *ScanTask) scanScene() {
	videoFile, err := ffmpeg.NewVideoFile(instance.StaticPaths.FFProbe, t.FilePath)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	qb := models.NewSceneQueryBuilder()
	scene, _ := qb.FindByPath(t.FilePath)
	if scene != nil {
		// We already have this item in the database, keep going
		return
	}

	checksum, err := t.calculateChecksum()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	t.makeScreenshots(*videoFile, checksum)

	scene, _ = qb.FindByChecksum(checksum)
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)
	if scene != nil {
		logger.Infof("%s already exists.  Updating path...", t.FilePath)
		scene.Path = t.FilePath
		_, err = qb.Update(*scene, tx)
	} else {
		logger.Infof("%s doesn't exist.  Creating new item...", t.FilePath)
		currentTime := time.Now()
		newScene := models.Scene{
			Checksum:   checksum,
			Path:       t.FilePath,
			Duration:   sql.NullFloat64{Float64: videoFile.Duration, Valid: true},
			VideoCodec: sql.NullString{String: videoFile.VideoCodec, Valid: true},
			AudioCodec: sql.NullString{String: videoFile.AudioCodec, Valid: true},
			Width:      sql.NullInt64{Int64: int64(videoFile.Width), Valid: true},
			Height:     sql.NullInt64{Int64: int64(videoFile.Height), Valid: true},
			Framerate:  sql.NullFloat64{Float64: videoFile.FrameRate, Valid: true},
			Bitrate:    sql.NullInt64{Int64: videoFile.Bitrate, Valid: true},
			Size:       sql.NullString{String: strconv.Itoa(int(videoFile.Size)), Valid: true},
			CreatedAt:  models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:  models.SQLiteTimestamp{Timestamp: currentTime},
		}
		_, err = qb.Create(newScene, tx)
	}

	if err != nil {
		logger.Error(err.Error())
		_ = tx.Rollback()
	} else if err := tx.Commit(); err != nil {
		logger.Error(err.Error())
	}
}

func (t *ScanTask) makeScreenshots(probeResult ffmpeg.VideoFile, checksum string) {
	thumbPath := instance.Paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	thumbExists, _ := utils.FileExists(thumbPath)
	normalExists, _ := utils.FileExists(normalPath)
	if thumbExists && normalExists {
		logger.Debug("Screenshots already exist for this path... skipping")
		return
	}

	t.makeScreenshot(probeResult, thumbPath, 5, 320)
	t.makeScreenshot(probeResult, normalPath, 2, probeResult.Width)
}

func (t *ScanTask) makeScreenshot(probeResult ffmpeg.VideoFile, outputPath string, quality int, width int) {
	encoder := ffmpeg.NewEncoder(instance.StaticPaths.FFMPEG)
	options := ffmpeg.ScreenshotOptions{
		OutputPath: outputPath,
		Quality:    quality,
		Time:       float64(probeResult.Duration) * 0.2,
		Width:      width,
	}
	encoder.Screenshot(probeResult, options)
}

func (t *ScanTask) calculateChecksum() (string, error) {
	logger.Infof("%s not found.  Calculating checksum...", t.FilePath)
	checksum, err := utils.MD5FromFilePath(t.FilePath)
	if err != nil {
		return "", err
	}
	logger.Debugf("Checksum calculated: %s", checksum)
	return checksum, nil
}
