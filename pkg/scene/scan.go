package scene

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
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

const mutexType = "scene"

type videoFileCreator interface {
	NewVideoFile(path string, stripFileExtension bool) (*ffmpeg.VideoFile, error)
}

type Scanner struct {
	file.Scanner

	StripFileExtension  bool
	UseFileMetadata     bool
	FileNamingAlgorithm models.HashAlgorithm

	Ctx              context.Context
	CaseSensitiveFs  bool
	TxnManager       models.TransactionManager
	Paths            *paths.Paths
	Screenshotter    screenshotter
	VideoFileCreator videoFileCreator
	PluginCache      *plugin.Cache
	MutexManager     *utils.MutexManager
}

func FileScanner(hasher file.Hasher, fileNamingAlgorithm models.HashAlgorithm, calculateMD5 bool) file.Scanner {
	return file.Scanner{
		Hasher:          hasher,
		CalculateOSHash: true,
		CalculateMD5:    fileNamingAlgorithm == models.HashAlgorithmMd5 || calculateMD5,
	}
}

func (scanner *Scanner) ScanExisting(existing file.FileBased, file file.SourceFile) (err error) {
	scanned, err := scanner.Scanner.ScanExisting(existing, file)
	if err != nil {
		return err
	}

	s := existing.(*models.Scene)

	path := scanned.New.Path
	interactive := getInteractive(path)

	oldHash := s.GetHash(scanner.FileNamingAlgorithm)
	changed := false

	var videoFile *ffmpeg.VideoFile

	if scanned.ContentsChanged() {
		logger.Infof("%s has been updated: rescanning", path)

		s.SetFile(*scanned.New)

		videoFile, err = scanner.VideoFileCreator.NewVideoFile(path, scanner.StripFileExtension)
		if err != nil {
			return err
		}

		videoFileToScene(s, videoFile)
		changed = true
	} else if scanned.FileUpdated() || s.Interactive != interactive {
		logger.Infof("Updated scene file %s", path)

		// update fields as needed
		s.SetFile(*scanned.New)
		changed = true
	}

	// check for container
	if !s.Format.Valid {
		if videoFile == nil {
			videoFile, err = scanner.VideoFileCreator.NewVideoFile(path, scanner.StripFileExtension)
			if err != nil {
				return err
			}
		}
		container := ffmpeg.MatchContainer(videoFile.Container, path)
		logger.Infof("Adding container %s to file %s", container, path)
		s.Format = models.NullString(string(container))
		changed = true
	}

	if changed {
		// we are operating on a checksum now, so grab a mutex on the checksum
		done := make(chan struct{})
		if scanned.New.OSHash != "" {
			scanner.MutexManager.Claim(mutexType, scanned.New.OSHash, done)
		}
		if scanned.New.Checksum != "" {
			scanner.MutexManager.Claim(mutexType, scanned.New.Checksum, done)
		}

		if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			defer close(done)
			qb := r.Scene()

			// ensure no clashes of hashes
			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
				dupe, _ := qb.FindByChecksum(s.Checksum.String)
				if dupe != nil {
					return fmt.Errorf("MD5 for file %s is the same as that of %s", path, dupe.Path)
				}
			}

			if scanned.New.OSHash != "" && scanned.Old.OSHash != scanned.New.OSHash {
				dupe, _ := qb.FindByOSHash(scanned.New.OSHash)
				if dupe != nil {
					return fmt.Errorf("OSHash for file %s is the same as that of %s", path, dupe.Path)
				}
			}

			s.Interactive = interactive
			s.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}

			_, err := qb.UpdateFull(*s)
			return err
		}); err != nil {
			return err
		}

		// Migrate any generated files if the hash has changed
		newHash := s.GetHash(scanner.FileNamingAlgorithm)
		if newHash != oldHash {
			MigrateHash(scanner.Paths, oldHash, newHash)
		}

		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
	}

	// We already have this item in the database
	// check for thumbnails, screenshots
	scanner.makeScreenshots(path, videoFile, s.GetHash(scanner.FileNamingAlgorithm))

	return nil
}

func (scanner *Scanner) ScanNew(file file.SourceFile) (retScene *models.Scene, err error) {
	scanned, err := scanner.Scanner.ScanNew(file)
	if err != nil {
		return nil, err
	}

	path := file.Path()
	checksum := scanned.Checksum
	oshash := scanned.OSHash

	// grab a mutex on the checksum and oshash
	done := make(chan struct{})
	if oshash != "" {
		scanner.MutexManager.Claim(mutexType, oshash, done)
	}
	if checksum != "" {
		scanner.MutexManager.Claim(mutexType, checksum, done)
	}

	defer close(done)

	// check for scene by checksum and oshash - MD5 should be
	// redundant, but check both
	var s *models.Scene
	if err := scanner.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		qb := r.Scene()
		if checksum != "" {
			s, _ = qb.FindByChecksum(checksum)
		}

		if s == nil {
			s, _ = qb.FindByOSHash(oshash)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	sceneHash := oshash

	if scanner.FileNamingAlgorithm == models.HashAlgorithmMd5 {
		sceneHash = checksum
	}

	interactive := getInteractive(file.Path())

	if s != nil {
		exists, _ := fsutil.FileExists(s.Path)
		if !scanner.CaseSensitiveFs {
			// #1426 - if file exists but is a case-insensitive match for the
			// original filename, then treat it as a move
			if exists && strings.EqualFold(path, s.Path) {
				exists = false
			}
		}

		if exists {
			logger.Infof("%s already exists. Duplicate of %s", path, s.Path)
		} else {
			logger.Infof("%s already exists. Updating path...", path)
			scenePartial := models.ScenePartial{
				ID:          s.ID,
				Path:        &path,
				Interactive: &interactive,
			}
			if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				_, err := r.Scene().Update(scenePartial)
				return err
			}); err != nil {
				return nil, err
			}

			scanner.makeScreenshots(path, nil, sceneHash)
			scanner.PluginCache.ExecutePostHooks(scanner.Ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
		}
	} else {
		logger.Infof("%s doesn't exist. Creating new item...", path)
		currentTime := time.Now()

		videoFile, err := scanner.VideoFileCreator.NewVideoFile(path, scanner.StripFileExtension)
		if err != nil {
			return nil, err
		}

		// Override title to be filename if UseFileMetadata is false
		if !scanner.UseFileMetadata {
			videoFile.SetTitleFromPath(scanner.StripFileExtension)
		}

		newScene := models.Scene{
			Checksum: sql.NullString{String: checksum, Valid: checksum != ""},
			OSHash:   sql.NullString{String: oshash, Valid: oshash != ""},
			Path:     path,
			FileModTime: models.NullSQLiteTimestamp{
				Timestamp: scanned.FileModTime,
				Valid:     true,
			},
			Title:       sql.NullString{String: videoFile.Title, Valid: true},
			CreatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			Interactive: interactive,
		}

		videoFileToScene(&newScene, videoFile)

		if scanner.UseFileMetadata {
			newScene.Details = sql.NullString{String: videoFile.Comment, Valid: true}
			_ = newScene.Date.Scan(videoFile.CreationTime)
		}

		if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			var err error
			retScene, err = r.Scene().Create(newScene)
			return err
		}); err != nil {
			return nil, err
		}

		scanner.makeScreenshots(path, videoFile, sceneHash)
		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, retScene.ID, plugin.SceneCreatePost, nil, nil)
	}

	return retScene, nil
}

func videoFileToScene(s *models.Scene, videoFile *ffmpeg.VideoFile) {
	container := ffmpeg.MatchContainer(videoFile.Container, s.Path)

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

func (scanner *Scanner) makeScreenshots(path string, probeResult *ffmpeg.VideoFile, checksum string) {
	thumbPath := scanner.Paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := scanner.Paths.Scene.GetScreenshotPath(checksum)

	thumbExists, _ := fsutil.FileExists(thumbPath)
	normalExists, _ := fsutil.FileExists(normalPath)

	if thumbExists && normalExists {
		return
	}

	if probeResult == nil {
		var err error
		probeResult, err = scanner.VideoFileCreator.NewVideoFile(path, scanner.StripFileExtension)

		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Infof("Regenerating images for %s", path)
	}

	at := float64(probeResult.Duration) * 0.2

	if !thumbExists {
		logger.Debugf("Creating thumbnail for %s", path)
		makeScreenshot(scanner.Screenshotter, *probeResult, thumbPath, 5, 320, at)
	}

	if !normalExists {
		logger.Debugf("Creating screenshot for %s", path)
		makeScreenshot(scanner.Screenshotter, *probeResult, normalPath, 2, probeResult.Width, at)
	}
}

func getInteractive(path string) bool {
	_, err := os.Stat(GetFunscriptPath(path))
	return err == nil
}
