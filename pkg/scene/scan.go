package scene

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

type videoFileCreator interface {
	NewVideoFile(path string, stripFileExtension bool) (*ffmpeg.VideoFile, error)
}

type Scanner struct {
	file.Scanner

	StripFileExtension  bool
	UseFileMetadata     bool
	FileNamingAlgorithm models.HashAlgorithm

	Ctx              context.Context
	TxnManager       models.TransactionManager
	Paths            *paths.Paths
	Screenshotter    screenshotter
	VideoFileCreator videoFileCreator
	PluginCache      *plugin.Cache

	NewScene *models.Scene
}

func FileScanner(hasher file.Hasher, statter file.Statter, fileNamingAlgorithm models.HashAlgorithm, calculateMD5 bool) file.Scanner {
	return file.Scanner{
		Hasher:          hasher,
		Statter:         statter,
		CalculateOSHash: true,
		CalculateMD5:    fileNamingAlgorithm == models.HashAlgorithmMd5 || calculateMD5,
		Done:            make(chan struct{}),
	}
}

func (scanner *Scanner) PostScan(scanned file.Scanned) error {
	if scanned.Old != nil {
		// should be an existing scene
		var scene *models.Scene
		if err := scanner.TxnManager.WithReadTxn(scanner.Ctx, func(r models.ReaderRepository) error {
			scenes, err := r.Scene().FindByFileID(scanned.Old.ID)
			if err != nil {
				return err
			}

			// assume only one scene for now
			if len(scenes) > 0 {
				scene = scenes[0]
			}
			return err
		}); err != nil {
			logger.Error(err.Error())
			return nil
		}

		if scene != nil {
			return scanner.ScanExisting(scene, scanned)
		}

		// we shouldn't be able to have an existing file without a scene, but
		// assuming that it's happened, treat it as a new scene
	}

	// assume a new file/scene
	return scanner.ScanNew(scanned.New)
}

func (scanner *Scanner) GenerateMetadata(dest *models.File, src file.SourceFile) error {
	videoFile, err := scanner.VideoFileCreator.NewVideoFile(src.Path(), scanner.StripFileExtension)
	if err != nil {
		return err
	}

	videoFileToFile(dest, videoFile)

	return nil
}

func (scanner *Scanner) ScanExisting(s *models.Scene, scanned file.Scanned) (err error) {
	path := scanned.New.Path
	interactive := getInteractive(path)

	config := config.GetInstance()
	oldHash := s.GetHash(scanner.FileNamingAlgorithm)
	changed := false

	if scanned.ContentsChanged() {
		logger.Infof("%s has been updated: rescanning", path)

		s.SetFile(*scanned.New)

		changed = true
	} else if scanned.FileUpdated() || s.Interactive != interactive {
		logger.Infof("Updated scene file %s", path)

		// update fields as needed
		s.SetFile(*scanned.New)
		changed = true
	}

	// check for container
	if !s.Format.Valid {
		s.Format = scanned.New.Format
		container := s.Format.String
		logger.Infof("Adding container %s to file %s", container, path)
		changed = true
	}

	if changed {
		if err := scanner.TxnManager.WithTxn(scanner.Ctx, func(r models.Repository) error {
			if err := scanner.Scanner.ApplyChanges(r.File(), &scanned); err != nil {
				return err
			}

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
		newHash := s.GetHash(config.GetVideoFileNamingAlgorithm())
		if newHash != oldHash {
			MigrateHash(scanner.Paths, oldHash, newHash)
		}

		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
	}

	// We already have this item in the database
	// check for thumbnails, screenshots
	scanner.makeScreenshots(path, nil, s.GetHash(scanner.FileNamingAlgorithm))

	return nil
}

func (scanner *Scanner) ScanNew(f *models.File) error {
	path := f.Path
	checksum := f.Checksum
	oshash := f.OSHash

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
		return err
	}

	sceneHash := oshash

	if scanner.FileNamingAlgorithm == models.HashAlgorithmMd5 {
		sceneHash = checksum
	}

	interactive := getInteractive(f.Path)

	if s != nil {
		// if exists {
		logger.Infof("%s already exists. Duplicate of %s", path, s.Path)

		if err := scanner.TxnManager.WithTxn(scanner.Ctx, func(r models.Repository) error {
			if err := scanner.Scanner.ApplyChanges(r.File(), &file.Scanned{
				New: f,
			}); err != nil {
				return err
			}

			// link scene to file
			return addFile(r.Scene(), s, f)
		}); err != nil {
			return err
		}

		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
	} else {
		logger.Infof("%s doesn't exist. Creating new item...", path)
		currentTime := time.Now()

		videoFile, err := scanner.VideoFileCreator.NewVideoFile(path, scanner.StripFileExtension)
		if err != nil {
			return err
		}

		// Override title to be filename if UseFileMetadata is false
		if !scanner.UseFileMetadata {
			videoFile.SetTitleFromPath(scanner.StripFileExtension)
		}

		newScene := models.Scene{
			Title:       sql.NullString{String: videoFile.Title, Valid: true},
			CreatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			Interactive: interactive,
		}

		newScene.SetFile(*f)

		if scanner.UseFileMetadata {
			newScene.Details = sql.NullString{String: videoFile.Comment, Valid: true}
			_ = newScene.Date.Scan(videoFile.CreationTime)
		}

		var retScene *models.Scene
		if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			if err := scanner.Scanner.ApplyChanges(r.File(), &file.Scanned{
				New: f,
			}); err != nil {
				return err
			}

			var err error
			retScene, err = r.Scene().Create(newScene)
			if err != nil {
				return err
			}

			// link scene to file
			return addFile(r.Scene(), retScene, f)
		}); err != nil {
			return err
		}

		scanner.makeScreenshots(path, videoFile, sceneHash)
		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, retScene.ID, plugin.SceneCreatePost, nil, nil)

		scanner.NewScene = retScene
	}

	return nil
}

func addFile(rw models.SceneReaderWriter, s *models.Scene, f *models.File) error {
	ids, err := rw.GetFileIDs(s.ID)
	if err != nil {
		return err
	}

	ids = utils.IntAppendUnique(ids, f.ID)
	return rw.UpdateFiles(s.ID, ids)
}

func videoFileToFile(s *models.File, videoFile *ffmpeg.VideoFile) {
	container := ffmpeg.MatchContainer(videoFile.Container, s.Path)

	s.Duration = sql.NullFloat64{Float64: videoFile.Duration, Valid: true}
	s.VideoCodec = sql.NullString{String: videoFile.VideoCodec, Valid: true}
	s.AudioCodec = sql.NullString{String: videoFile.AudioCodec, Valid: true}
	s.Format = sql.NullString{String: string(container), Valid: true}
	s.Width = sql.NullInt64{Int64: int64(videoFile.Width), Valid: true}
	s.Height = sql.NullInt64{Int64: int64(videoFile.Height), Valid: true}
	s.Framerate = sql.NullFloat64{Float64: videoFile.FrameRate, Valid: true}
	s.Bitrate = sql.NullInt64{Int64: videoFile.Bitrate, Valid: true}
}

func (scanner *Scanner) makeScreenshots(path string, probeResult *ffmpeg.VideoFile, checksum string) {
	thumbPath := scanner.Paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := scanner.Paths.Scene.GetScreenshotPath(checksum)

	thumbExists, _ := utils.FileExists(thumbPath)
	normalExists, _ := utils.FileExists(normalPath)

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
	_, err := os.Stat(utils.GetFunscriptPath(path))
	return err == nil
}
