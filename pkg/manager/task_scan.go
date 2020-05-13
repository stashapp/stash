package manager

import (
	"context"
	"database/sql"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type ScanTask struct {
	FilePath        string
	UseFileMetadata bool
}

func (t *ScanTask) Start(wg *sync.WaitGroup) {
	if isGallery(t.FilePath) || isGalleryFolder(t.FilePath) {
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
	var checksum string
	var err error
	if models.IsPathArchive(t.FilePath) {
		checksum, err = t.calculateChecksum()
	} else {
		checksum, err = models.ChecksumFromDirPath(t.FilePath)
	}

	ok, err := utils.IsZipFileUncompressed(t.FilePath)
	if err == nil && !ok {
		logger.Warnf("%s is using above store (0) level compression.", t.FilePath)
	}
	if err != nil {
		logger.Error(err.Error())
		return
	}

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)
	gallery, _ = qb.FindByChecksum(checksum, tx)
	if gallery != nil {
		var exists bool
		if models.IsPathArchive(t.FilePath) {
			exists, _ = utils.FileExists(gallery.Path)
		} else { // on image folder , folder might still be there but without same images
			checksumRe, err := models.ChecksumFromDirPath(gallery.Path)
			exists = (err != nil && checksumRe == checksum)
		}
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

// associates a gallery to a scene with the same basename
func (t *ScanTask) associateGallery(wg *sync.WaitGroup) {
	qb := models.NewGalleryQueryBuilder()
	gallery, _ := qb.FindByPath(t.FilePath)
	if gallery == nil {
		// shouldn't happen , associate is run after scan is finished
		logger.Errorf("associate: gallery %s not found in DB", t.FilePath)
		wg.Done()
		return
	}

	if !gallery.SceneID.Valid { // gallery has no SceneID
		basename := t.FilePath
		if models.IsPathArchive(t.FilePath) {
			basename = strings.TrimSuffix(t.FilePath, filepath.Ext(t.FilePath))
		}
		var relatedFiles []string
		for _, ext := range extensionsToScan { // make a list of media files that can be related to the gallery
			related := basename + "." + ext
			if !isGallery(related) { //exclude gallery extensions from the related files
				relatedFiles = append(relatedFiles, related)
			}
		}
		for _, scenePath := range relatedFiles {
			qbScene := models.NewSceneQueryBuilder()
			scene, _ := qbScene.FindByPath(scenePath)
			if scene != nil { // found related Scene
				logger.Infof("associate: Gallery %s is related to scene: %d", t.FilePath, scene.ID)

				gallery.SceneID.Int64 = int64(scene.ID)
				gallery.SceneID.Valid = true

				ctx := context.TODO()
				tx := database.DB.MustBeginTx(ctx, nil)

				_, err := qb.Update(*gallery, tx)
				if err != nil {
					logger.Errorf("associate: Error updating gallery sceneId %s", err)
					_ = tx.Rollback()
				} else if err := tx.Commit(); err != nil {
					logger.Error(err.Error())
				}

				break // since a gallery can have only one related scene
				// only first found is associated
			}

		}

	}
	wg.Done()
}

func (t *ScanTask) scanScene() {
	qb := models.NewSceneQueryBuilder()
	scene, _ := qb.FindByPath(t.FilePath)
	if scene != nil {
		// We already have this item in the database
		//check for thumbnails,screenshots
		t.makeScreenshots(nil, scene.Checksum)

		//check for container
		if !scene.Format.Valid {
			videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)
			logger.Infof("Adding container %s to file %s", container, t.FilePath)

			ctx := context.TODO()
			tx := database.DB.MustBeginTx(ctx, nil)
			err = qb.UpdateFormat(scene.ID, string(container), tx)
			if err != nil {
				logger.Error(err.Error())
				_ = tx.Rollback()
			} else if err := tx.Commit(); err != nil {
				logger.Error(err.Error())
			}

		}
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)

	// Override title to be filename if UseFileMetadata is false
	if !t.UseFileMetadata {
		videoFile.SetTitleFromPath()
	}

	checksum, err := t.calculateChecksum()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	t.makeScreenshots(videoFile, checksum)

	scene, _ = qb.FindByChecksum(checksum)
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)
	if scene != nil {
		exists, _ := utils.FileExists(scene.Path)
		if exists {
			logger.Infof("%s already exists.  Duplicate of %s ", t.FilePath, scene.Path)
		} else {
			logger.Infof("%s already exists.  Updating path...", t.FilePath)
			scenePartial := models.ScenePartial{
				ID:   scene.ID,
				Path: &t.FilePath,
			}
			_, err = qb.Update(scenePartial, tx)
		}
	} else {
		logger.Infof("%s doesn't exist.  Creating new item...", t.FilePath)
		currentTime := time.Now()
		newScene := models.Scene{
			Checksum:   checksum,
			Path:       t.FilePath,
			Title:      sql.NullString{String: videoFile.Title, Valid: true},
			Duration:   sql.NullFloat64{Float64: videoFile.Duration, Valid: true},
			VideoCodec: sql.NullString{String: videoFile.VideoCodec, Valid: true},
			AudioCodec: sql.NullString{String: videoFile.AudioCodec, Valid: true},
			Format:     sql.NullString{String: string(container), Valid: true},
			Width:      sql.NullInt64{Int64: int64(videoFile.Width), Valid: true},
			Height:     sql.NullInt64{Int64: int64(videoFile.Height), Valid: true},
			Framerate:  sql.NullFloat64{Float64: videoFile.FrameRate, Valid: true},
			Bitrate:    sql.NullInt64{Int64: videoFile.Bitrate, Valid: true},
			Size:       sql.NullString{String: strconv.Itoa(int(videoFile.Size)), Valid: true},
			CreatedAt:  models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:  models.SQLiteTimestamp{Timestamp: currentTime},
		}

		if t.UseFileMetadata {
			newScene.Details = sql.NullString{String: videoFile.Comment, Valid: true}
			newScene.Date = models.SQLiteDate{String: videoFile.CreationTime.Format("2006-01-02")}
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

func (t *ScanTask) makeScreenshots(probeResult *ffmpeg.VideoFile, checksum string) {
	thumbPath := instance.Paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	thumbExists, _ := utils.FileExists(thumbPath)
	normalExists, _ := utils.FileExists(normalPath)

	if thumbExists && normalExists {
		logger.Debug("Screenshots already exist for this path... skipping")
		return
	}

	if probeResult == nil {
		var err error
		probeResult, err = ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath)

		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Infof("Regenerating images for %s", t.FilePath)
	}

	at := float64(probeResult.Duration) * 0.2

	if !thumbExists {
		logger.Debugf("Creating thumbnail for %s", t.FilePath)
		makeScreenshot(*probeResult, thumbPath, 5, 320, at)
	}

	if !normalExists {
		logger.Debugf("Creating screenshot for %s", t.FilePath)
		makeScreenshot(*probeResult, normalPath, 2, probeResult.Width, at)
	}
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

func (t *ScanTask) doesPathExist() bool {
	if isGallery(t.FilePath) || isGalleryFolder(t.FilePath) {
		qb := models.NewGalleryQueryBuilder()
		gallery, _ := qb.FindByPath(t.FilePath)
		if gallery != nil {
			return true
		}
	} else {
		qb := models.NewSceneQueryBuilder()
		scene, _ := qb.FindByPath(t.FilePath)
		if scene != nil {
			return true
		}
	}
	return false
}
