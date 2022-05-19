package scene

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
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
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

const mutexType = "scene"

type CreatorUpdater interface {
	FindByChecksum(ctx context.Context, checksum string) (*models.Scene, error)
	FindByOSHash(ctx context.Context, oshash string) (*models.Scene, error)
	Create(ctx context.Context, newScene models.Scene) (*models.Scene, error)
	UpdateFull(ctx context.Context, updatedScene models.Scene) (*models.Scene, error)
	Update(ctx context.Context, updatedScene models.ScenePartial) (*models.Scene, error)

	GetCaptions(ctx context.Context, sceneID int) ([]*models.SceneCaption, error)
	UpdateCaptions(ctx context.Context, id int, captions []*models.SceneCaption) error
}

type videoFileCreator interface {
	NewVideoFile(path string) (*ffmpeg.VideoFile, error)
}

type Scanner struct {
	file.Scanner

	StripFileExtension  bool
	UseFileMetadata     bool
	FileNamingAlgorithm models.HashAlgorithm

	CaseSensitiveFs  bool
	TxnManager       txn.Manager
	CreatorUpdater   CreatorUpdater
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

func (scanner *Scanner) ScanExisting(ctx context.Context, existing file.FileBased, file file.SourceFile) (err error) {
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

		videoFile, err = scanner.VideoFileCreator.NewVideoFile(path)
		if err != nil {
			return err
		}

		if err := videoFileToScene(s, videoFile); err != nil {
			return err
		}
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
			videoFile, err = scanner.VideoFileCreator.NewVideoFile(path)
			if err != nil {
				return err
			}
		}
		container, err := ffmpeg.MatchContainer(videoFile.Container, path)
		if err != nil {
			return fmt.Errorf("getting container for %s: %w", path, err)
		}
		logger.Infof("Adding container %s to file %s", container, path)
		s.Format = models.NullString(string(container))
		changed = true
	}

	qb := scanner.CreatorUpdater

	if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
		var err error

		captions, er := qb.GetCaptions(ctx, s.ID)
		if er == nil {
			if len(captions) > 0 {
				clean, altered := CleanCaptions(s.Path, captions)
				if altered {
					er = qb.UpdateCaptions(ctx, s.ID, clean)
					if er == nil {
						logger.Debugf("Captions for %s cleaned: %s -> %s", path, captions, clean)
					}
				}
			}
		}
		return err
	}); err != nil {
		logger.Error(err.Error())
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

		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
			defer close(done)
			qb := scanner.CreatorUpdater

			// ensure no clashes of hashes
			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
				dupe, _ := qb.FindByChecksum(ctx, s.Checksum.String)
				if dupe != nil {
					return fmt.Errorf("MD5 for file %s is the same as that of %s", path, dupe.Path)
				}
			}

			if scanned.New.OSHash != "" && scanned.Old.OSHash != scanned.New.OSHash {
				dupe, _ := qb.FindByOSHash(ctx, scanned.New.OSHash)
				if dupe != nil {
					return fmt.Errorf("OSHash for file %s is the same as that of %s", path, dupe.Path)
				}
			}

			s.Interactive = interactive
			s.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}

			_, err := qb.UpdateFull(ctx, *s)
			return err
		}); err != nil {
			return err
		}

		// Migrate any generated files if the hash has changed
		newHash := s.GetHash(scanner.FileNamingAlgorithm)
		if newHash != oldHash {
			MigrateHash(scanner.Paths, oldHash, newHash)
		}

		scanner.PluginCache.ExecutePostHooks(ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
	}

	// We already have this item in the database
	// check for thumbnails, screenshots
	scanner.makeScreenshots(path, videoFile, s.GetHash(scanner.FileNamingAlgorithm))

	return nil
}

func (scanner *Scanner) ScanNew(ctx context.Context, file file.SourceFile) (retScene *models.Scene, err error) {
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
	if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
		qb := scanner.CreatorUpdater
		if checksum != "" {
			s, _ = qb.FindByChecksum(ctx, checksum)
		}

		if s == nil {
			s, _ = qb.FindByOSHash(ctx, oshash)
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
			if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
				_, err := scanner.CreatorUpdater.Update(ctx, scenePartial)
				return err
			}); err != nil {
				return nil, err
			}

			scanner.makeScreenshots(path, nil, sceneHash)
			scanner.PluginCache.ExecutePostHooks(ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
		}
	} else {
		logger.Infof("%s doesn't exist. Creating new item...", path)
		currentTime := time.Now()

		videoFile, err := scanner.VideoFileCreator.NewVideoFile(path)
		if err != nil {
			return nil, err
		}

		title := filepath.Base(path)
		if scanner.StripFileExtension {
			title = stripExtension(title)
		}

		if scanner.UseFileMetadata && videoFile.Title != "" {
			title = videoFile.Title
		}

		newScene := models.Scene{
			Checksum: sql.NullString{String: checksum, Valid: checksum != ""},
			OSHash:   sql.NullString{String: oshash, Valid: oshash != ""},
			Path:     path,
			FileModTime: models.NullSQLiteTimestamp{
				Timestamp: scanned.FileModTime,
				Valid:     true,
			},
			Title:       sql.NullString{String: title, Valid: true},
			CreatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			Interactive: interactive,
		}

		if err := videoFileToScene(&newScene, videoFile); err != nil {
			return nil, err
		}

		if scanner.UseFileMetadata {
			newScene.Details = sql.NullString{String: videoFile.Comment, Valid: true}
			_ = newScene.Date.Scan(videoFile.CreationTime)
		}

		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
			var err error
			retScene, err = scanner.CreatorUpdater.Create(ctx, newScene)
			return err
		}); err != nil {
			return nil, err
		}

		scanner.makeScreenshots(path, videoFile, sceneHash)
		scanner.PluginCache.ExecutePostHooks(ctx, retScene.ID, plugin.SceneCreatePost, nil, nil)
	}

	return retScene, nil
}

func stripExtension(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}

func videoFileToScene(s *models.Scene, videoFile *ffmpeg.VideoFile) error {
	container, err := ffmpeg.MatchContainer(videoFile.Container, s.Path)
	if err != nil {
		return fmt.Errorf("matching container: %w", err)
	}

	s.Duration = sql.NullFloat64{Float64: videoFile.Duration, Valid: true}
	s.VideoCodec = sql.NullString{String: videoFile.VideoCodec, Valid: true}
	s.AudioCodec = sql.NullString{String: videoFile.AudioCodec, Valid: true}
	s.Format = sql.NullString{String: string(container), Valid: true}
	s.Width = sql.NullInt64{Int64: int64(videoFile.Width), Valid: true}
	s.Height = sql.NullInt64{Int64: int64(videoFile.Height), Valid: true}
	s.Framerate = sql.NullFloat64{Float64: videoFile.FrameRate, Valid: true}
	s.Bitrate = sql.NullInt64{Int64: videoFile.Bitrate, Valid: true}
	s.Size = sql.NullString{String: strconv.FormatInt(videoFile.Size, 10), Valid: true}

	return nil
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
		probeResult, err = scanner.VideoFileCreator.NewVideoFile(path)

		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Infof("Regenerating images for %s", path)
	}

	if !thumbExists {
		logger.Debugf("Creating thumbnail for %s", path)
		if err := scanner.Screenshotter.GenerateThumbnail(context.TODO(), probeResult, checksum); err != nil {
			logger.Errorf("Error creating thumbnail for %s: %v", err)
		}
	}

	if !normalExists {
		logger.Debugf("Creating screenshot for %s", path)
		if err := scanner.Screenshotter.GenerateScreenshot(context.TODO(), probeResult, checksum); err != nil {
			logger.Errorf("Error creating screenshot for %s: %v", err)
		}
	}
}

func getInteractive(path string) bool {
	_, err := os.Stat(GetFunscriptPath(path))
	return err == nil
}
