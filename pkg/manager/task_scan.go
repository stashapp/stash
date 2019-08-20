package manager

import (
	"context"
	"database/sql"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
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
		exists, _ := utils.FileExists(gallery.Path)
		if exists {
			logger.Infof("%s already exists.  Duplicate of %s ", t.FilePath, gallery.Path)
		} else {

			logger.Infof("%s already exists.  Updating path...", t.FilePath)
			gallery.Path = t.FilePath
			_, err = qb.Update(*gallery, tx)
		}
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
	qb := models.NewSceneQueryBuilder()
	scene, _ := qb.FindByPath(t.FilePath)
	if scene != nil {
		// We already have this item in the database, keep going
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath)
	if err != nil {
		logger.Error(err.Error())
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
		exists, _ := utils.FileExists(scene.Path)
		if exists {
			logger.Infof("%s already exists.  Duplicate of %s ", t.FilePath, scene.Path)
		} else {
			logger.Infof("%s already exists.  Updating path...", t.FilePath)
			scene.Path = t.FilePath
			_, err = qb.Update(*scene, tx)
		}
	} else {
		logger.Infof("%s doesn't exist.  Creating new item...", t.FilePath)
		currentTime := time.Now()
		newScene := models.Scene{
			Checksum:   checksum,
			Path:       t.FilePath,
			Title:      sql.NullString{String: videoFile.Title, Valid: true},
			Details:    sql.NullString{String: videoFile.Comment, Valid: true},
			Date:       models.SQLiteDate{String: videoFile.CreationTime.Format("2006-01-02")},
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
	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)
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
