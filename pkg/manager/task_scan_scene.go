package manager

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

func (t *ScanTask) scanScene() *models.Scene {
	logError := func(err error) *models.Scene {
		logger.Error(err.Error())
		return nil
	}

	var retScene *models.Scene
	var s *models.Scene

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		s, err = r.Scene().FindByPath(t.FilePath)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return nil
	}

	scanner := file.Scanner{
		Hasher:          &file.FSHasher{},
		CalculateOSHash: true,
		CalculateMD5:    t.fileNamingAlgorithm == models.HashAlgorithmMd5 || t.calculateMD5,
	}

	if s != nil {
		scanned, err := scanner.ScanExisting(s, t.FilePath, t.FileInfo)
		if err != nil {
			return logError(err)
		}

		if err := t.scanSceneExisting(s, scanned); err != nil {
			return logError(err)
		}

		return nil
	}

	file, err := scanner.ScanNew(t.FilePath, t.FileInfo)
	if err != nil {
		return logError(err)
	}

	retScene, err = t.scanSceneNew(file)
	if err != nil {
		return logError(err)
	}

	return retScene
}

func (t *ScanTask) scanSceneExisting(s *models.Scene, scanned *file.Scanned) (err error) {
	interactive := t.getInteractive()

	config := config.GetInstance()
	oldHash := s.GetHash(config.GetVideoFileNamingAlgorithm())
	changed := false

	if scanned.ContentsChanged() {
		logger.Infof("%s has been updated: rescanning", t.FilePath)

		videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
		if err != nil {
			return err
		}

		t.videoFileToScene(s, videoFile)
		changed = true
	} else if scanned.FileUpdated() || s.Interactive != interactive {
		logger.Infof("Updated scene file %s", t.FilePath)

		// update fields as needed
		s.SetFile(*scanned.New)
		changed = true
	}

	// check for container
	if !s.Format.Valid {
		videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
		if err != nil {
			return err
		}
		container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)
		logger.Infof("Adding container %s to file %s", container, t.FilePath)
		s.Format = models.NullString(string(container))
		changed = true
	}

	if changed {
		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			qb := r.Scene()

			// ensure no clashes of hashes
			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
				dupe, _ := qb.FindByChecksum(s.Checksum.String)
				if dupe != nil {
					return fmt.Errorf("MD5 for file %s is the same as that of %s", t.FilePath, dupe.Path)
				}
			}

			if scanned.New.OSHash != "" && scanned.Old.OSHash != scanned.New.OSHash {
				dupe, _ := qb.FindByOSHash(scanned.New.OSHash)
				if dupe != nil {
					return fmt.Errorf("OSHash for file %s is the same as that of %s", t.FilePath, dupe.Path)
				}
			}

			_, err := qb.UpdateFull(*s)
			return err
		}); err != nil {
			return err
		}

		// Migrate any generated files if the hash has changed
		newHash := s.GetHash(config.GetVideoFileNamingAlgorithm())
		if newHash != oldHash {
			MigrateHash(oldHash, newHash)
		}

		GetInstance().PluginCache.ExecutePostHooks(t.ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
	}

	// We already have this item in the database
	// check for thumbnails, screenshots
	t.makeScreenshots(nil, s.GetHash(t.fileNamingAlgorithm))

	return nil
}

func (t *ScanTask) videoFileToScene(s *models.Scene, videoFile *ffmpeg.VideoFile) {
	container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)

	s.Duration = sql.NullFloat64{Float64: videoFile.Duration, Valid: true}
	s.VideoCodec = sql.NullString{String: videoFile.VideoCodec, Valid: true}
	s.AudioCodec = sql.NullString{String: videoFile.AudioCodec, Valid: true}
	s.Format = sql.NullString{String: string(container), Valid: true}
	s.Width = sql.NullInt64{Int64: int64(videoFile.Width), Valid: true}
	s.Height = sql.NullInt64{Int64: int64(videoFile.Height), Valid: true}
	s.Framerate = sql.NullFloat64{Float64: videoFile.FrameRate, Valid: true}
	s.Bitrate = sql.NullInt64{Int64: videoFile.Bitrate, Valid: true}
	s.Size = sql.NullString{String: strconv.FormatInt(videoFile.Size, 10), Valid: true}
}

func (t *ScanTask) scanSceneNew(file *models.File) (retScene *models.Scene, err error) {
	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
	if err != nil {
		return nil, err
	}

	// Override title to be filename if UseFileMetadata is false
	if !t.UseFileMetadata {
		videoFile.SetTitleFromPath(t.StripFileExtension)
	}

	checksum := file.Checksum
	oshash := file.OSHash

	// check for scene by checksum and oshash - MD5 should be
	// redundant, but check both
	var s *models.Scene
	t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		qb := r.Scene()
		if checksum != "" {
			s, _ = qb.FindByChecksum(checksum)
		}

		if s == nil {
			s, _ = qb.FindByOSHash(oshash)
		}

		return nil
	})

	sceneHash := oshash

	if t.fileNamingAlgorithm == models.HashAlgorithmMd5 {
		sceneHash = checksum
	}

	t.makeScreenshots(videoFile, sceneHash)
	interactive := t.getInteractive()

	if s != nil {
		exists, _ := utils.FileExists(s.Path)
		if !t.CaseSensitiveFs {
			// #1426 - if file exists but is a case-insensitive match for the
			// original filename, then treat it as a move
			if exists && strings.EqualFold(t.FilePath, s.Path) {
				exists = false
			}
		}

		if exists {
			logger.Infof("%s already exists. Duplicate of %s", t.FilePath, s.Path)
		} else {
			logger.Infof("%s already exists. Updating path...", t.FilePath)
			scenePartial := models.ScenePartial{
				ID:          s.ID,
				Path:        &t.FilePath,
				Interactive: &interactive,
			}
			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				_, err := r.Scene().Update(scenePartial)
				return err
			}); err != nil {
				return nil, err
			}

			GetInstance().PluginCache.ExecutePostHooks(t.ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
		}
	} else {
		logger.Infof("%s doesn't exist. Creating new item...", t.FilePath)
		currentTime := time.Now()
		newScene := models.Scene{
			Checksum: sql.NullString{String: checksum, Valid: checksum != ""},
			OSHash:   sql.NullString{String: oshash, Valid: oshash != ""},
			Path:     t.FilePath,
			FileModTime: models.NullSQLiteTimestamp{
				Timestamp: file.FileModTime,
				Valid:     true,
			},
			Title:       sql.NullString{String: videoFile.Title, Valid: true},
			CreatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			Interactive: interactive,
		}

		t.videoFileToScene(&newScene, videoFile)

		if t.UseFileMetadata {
			newScene.Details = sql.NullString{String: videoFile.Comment, Valid: true}
			newScene.Date = models.SQLiteDate{String: videoFile.CreationTime.Format("2006-01-02")}
		}

		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			var err error
			retScene, err = r.Scene().Create(newScene)
			return err
		}); err != nil {
			return nil, err
		}

		GetInstance().PluginCache.ExecutePostHooks(t.ctx, retScene.ID, plugin.SceneCreatePost, nil, nil)
	}

	return retScene, nil
}

func (t *ScanTask) makeScreenshots(probeResult *ffmpeg.VideoFile, checksum string) {
	thumbPath := instance.Paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	thumbExists, _ := utils.FileExists(thumbPath)
	normalExists, _ := utils.FileExists(normalPath)

	if thumbExists && normalExists {
		return
	}

	if probeResult == nil {
		var err error
		probeResult, err = ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)

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

func (t *ScanTask) getInteractive() bool {
	_, err := os.Stat(utils.GetFunscriptPath(t.FilePath))
	return err == nil
}
