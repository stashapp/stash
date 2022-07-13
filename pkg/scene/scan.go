package scene

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
)

var (
	ErrNotVideoFile = errors.New("not a video file")
)

// const mutexType = "scene"

type CreatorUpdater interface {
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Scene, error)
	FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Scene, error)
	Create(ctx context.Context, newScene *models.Scene, fileIDs []file.ID) error
	Update(ctx context.Context, updatedScene *models.Scene) error
	UpdatePartial(ctx context.Context, id int, updatedScene models.ScenePartial) (*models.Scene, error)
}

type ScanGenerator interface {
	Generate(ctx context.Context, s *models.Scene, f *file.VideoFile) error
}

type ScanHandler struct {
	CreatorUpdater CreatorUpdater

	CoverGenerator CoverGenerator
	ScanGenerator  ScanGenerator
	PluginCache    *plugin.Cache
}

func (h *ScanHandler) validate() error {
	if h.CreatorUpdater == nil {
		return errors.New("CreatorUpdater is required")
	}
	if h.CoverGenerator == nil {
		return errors.New("CoverGenerator is required")
	}
	if h.ScanGenerator == nil {
		return errors.New("ScanGenerator is required")
	}

	return nil
}

func (h *ScanHandler) Handle(ctx context.Context, f file.File) error {
	if err := h.validate(); err != nil {
		return err
	}

	videoFile, ok := f.(*file.VideoFile)
	if !ok {
		return ErrNotVideoFile
	}

	// try to match the file to a scene
	existing, err := h.CreatorUpdater.FindByFileID(ctx, f.Base().ID)
	if err != nil {
		return fmt.Errorf("finding existing scene: %w", err)
	}

	if len(existing) == 0 {
		// try also to match file by fingerprints
		existing, err = h.CreatorUpdater.FindByFingerprints(ctx, videoFile.Fingerprints)
		if err != nil {
			return fmt.Errorf("finding existing scene by fingerprints: %w", err)
		}
	}

	if len(existing) > 0 {
		if err := h.associateExisting(ctx, existing, videoFile); err != nil {
			return err
		}
	} else {
		// create a new scene
		now := time.Now()
		newScene := &models.Scene{
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := h.CreatorUpdater.Create(ctx, newScene, []file.ID{videoFile.ID}); err != nil {
			return fmt.Errorf("creating new scene: %w", err)
		}

		h.PluginCache.ExecutePostHooks(ctx, newScene.ID, plugin.SceneCreatePost, nil, nil)

		existing = []*models.Scene{newScene}
	}

	for _, s := range existing {
		if err := h.CoverGenerator.GenerateCover(ctx, s, videoFile); err != nil {
			// just log if cover generation fails. We can try again on rescan
			logger.Errorf("Error generating cover for %s: %v", videoFile.Path, err)
		}

		if err := h.ScanGenerator.Generate(ctx, s, videoFile); err != nil {
			// just log if cover generation fails. We can try again on rescan
			logger.Errorf("Error generating content for %s: %v", videoFile.Path, err)
		}
	}

	return nil
}

func (h *ScanHandler) associateExisting(ctx context.Context, existing []*models.Scene, f *file.VideoFile) error {
	for _, s := range existing {
		found := false
		for _, sf := range s.Files {
			if sf.ID == f.Base().ID {
				found = true
				break
			}
		}

		if !found {
			logger.Infof("Adding %s to scene %s", f.Path, s.GetTitle())
			s.Files = append(s.Files, f)
		}

		if err := h.CreatorUpdater.Update(ctx, s); err != nil {
			return fmt.Errorf("updating scene: %w", err)
		}
	}

	return nil
}

// type videoFileCreator interface {
// 	NewVideoFile(path string) (*ffmpeg.VideoFile, error)
// }

// type Scanner struct {
// 	file.Scanner

// 	StripFileExtension  bool
// 	UseFileMetadata     bool
// 	FileNamingAlgorithm models.HashAlgorithm

// 	CaseSensitiveFs  bool
// 	TxnManager       txn.Manager
// 	CreatorUpdater   CreatorUpdater
// 	Paths            *paths.Paths
// 	Screenshotter    screenshotter
// 	VideoFileCreator videoFileCreator
// 	PluginCache      *plugin.Cache
// 	MutexManager     *utils.MutexManager
// }

// func FileScanner(hasher file.Hasher, fileNamingAlgorithm models.HashAlgorithm, calculateMD5 bool) file.Scanner {
// 	return file.Scanner{
// 		Hasher:          hasher,
// 		CalculateOSHash: true,
// 		CalculateMD5:    fileNamingAlgorithm == models.HashAlgorithmMd5 || calculateMD5,
// 	}
// }

// func (scanner *Scanner) ScanExisting(ctx context.Context, existing file.FileBased, file file.SourceFile) (err error) {
// 	scanned, err := scanner.Scanner.ScanExisting(existing, file)
// 	if err != nil {
// 		return err
// 	}

// 	s := existing.(*models.Scene)

// 	path := scanned.New.Path
// 	interactive := getInteractive(path)

// 	oldHash := s.GetHash(scanner.FileNamingAlgorithm)
// 	changed := false

// 	var videoFile *ffmpeg.VideoFile

// 	if scanned.ContentsChanged() {
// 		logger.Infof("%s has been updated: rescanning", path)

// 		s.SetFile(*scanned.New)

// 		videoFile, err = scanner.VideoFileCreator.NewVideoFile(path)
// 		if err != nil {
// 			return err
// 		}

// 		if err := videoFileToScene(s, videoFile); err != nil {
// 			return err
// 		}
// 		changed = true
// 	} else if scanned.FileUpdated() || s.Interactive != interactive {
// 		logger.Infof("Updated scene file %s", path)

// 		// update fields as needed
// 		s.SetFile(*scanned.New)
// 		changed = true
// 	}

// 	// check for container
// 	if s.Format == nil {
// 		if videoFile == nil {
// 			videoFile, err = scanner.VideoFileCreator.NewVideoFile(path)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		container, err := ffmpeg.MatchContainer(videoFile.Container, path)
// 		if err != nil {
// 			return fmt.Errorf("getting container for %s: %w", path, err)
// 		}
// 		logger.Infof("Adding container %s to file %s", container, path)
// 		containerStr := string(container)
// 		s.Format = &containerStr
// 		changed = true
// 	}

// 	qb := scanner.CreatorUpdater

// 	if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 		var err error

// 		captions, er := qb.GetCaptions(ctx, s.ID)
// 		if er == nil {
// 			if len(captions) > 0 {
// 				clean, altered := CleanCaptions(s.Path, captions)
// 				if altered {
// 					er = qb.UpdateCaptions(ctx, s.ID, clean)
// 					if er == nil {
// 						logger.Debugf("Captions for %s cleaned: %s -> %s", path, captions, clean)
// 					}
// 				}
// 			}
// 		}
// 		return err
// 	}); err != nil {
// 		logger.Error(err.Error())
// 	}

// 	if changed {
// 		// we are operating on a checksum now, so grab a mutex on the checksum
// 		done := make(chan struct{})
// 		if scanned.New.OSHash != "" {
// 			scanner.MutexManager.Claim(mutexType, scanned.New.OSHash, done)
// 		}
// 		if scanned.New.Checksum != "" {
// 			scanner.MutexManager.Claim(mutexType, scanned.New.Checksum, done)
// 		}

// 		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 			defer close(done)
// 			qb := scanner.CreatorUpdater

// 			// ensure no clashes of hashes
// 			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
// 				dupe, _ := qb.FindByChecksum(ctx, *s.Checksum)
// 				if dupe != nil {
// 					return fmt.Errorf("MD5 for file %s is the same as that of %s", path, dupe.Path)
// 				}
// 			}

// 			if scanned.New.OSHash != "" && scanned.Old.OSHash != scanned.New.OSHash {
// 				dupe, _ := qb.FindByOSHash(ctx, scanned.New.OSHash)
// 				if dupe != nil {
// 					return fmt.Errorf("OSHash for file %s is the same as that of %s", path, dupe.Path)
// 				}
// 			}

// 			s.Interactive = interactive
// 			s.UpdatedAt = time.Now()

// 			return qb.Update(ctx, s)
// 		}); err != nil {
// 			return err
// 		}

// 		// Migrate any generated files if the hash has changed
// 		newHash := s.GetHash(scanner.FileNamingAlgorithm)
// 		if newHash != oldHash {
// 			MigrateHash(scanner.Paths, oldHash, newHash)
// 		}

// 		scanner.PluginCache.ExecutePostHooks(ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
// 	}

// 	// We already have this item in the database
// 	// check for thumbnails, screenshots
// 	scanner.makeScreenshots(path, videoFile, s.GetHash(scanner.FileNamingAlgorithm))

// 	return nil
// }

// func (scanner *Scanner) ScanNew(ctx context.Context, file file.SourceFile) (retScene *models.Scene, err error) {
// 	scanned, err := scanner.Scanner.ScanNew(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	path := file.Path()
// 	checksum := scanned.Checksum
// 	oshash := scanned.OSHash

// 	// grab a mutex on the checksum and oshash
// 	done := make(chan struct{})
// 	if oshash != "" {
// 		scanner.MutexManager.Claim(mutexType, oshash, done)
// 	}
// 	if checksum != "" {
// 		scanner.MutexManager.Claim(mutexType, checksum, done)
// 	}

// 	defer close(done)

// 	// check for scene by checksum and oshash - MD5 should be
// 	// redundant, but check both
// 	var s *models.Scene
// 	if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 		qb := scanner.CreatorUpdater
// 		if checksum != "" {
// 			s, _ = qb.FindByChecksum(ctx, checksum)
// 		}

// 		if s == nil {
// 			s, _ = qb.FindByOSHash(ctx, oshash)
// 		}

// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}

// 	sceneHash := oshash

// 	if scanner.FileNamingAlgorithm == models.HashAlgorithmMd5 {
// 		sceneHash = checksum
// 	}

// 	interactive := getInteractive(file.Path())

// 	if s != nil {
// 		exists, _ := fsutil.FileExists(s.Path)
// 		if !scanner.CaseSensitiveFs {
// 			// #1426 - if file exists but is a case-insensitive match for the
// 			// original filename, then treat it as a move
// 			if exists && strings.EqualFold(path, s.Path) {
// 				exists = false
// 			}
// 		}

// 		if exists {
// 			logger.Infof("%s already exists. Duplicate of %s", path, s.Path)
// 		} else {
// 			logger.Infof("%s already exists. Updating path...", path)
// 			scenePartial := models.ScenePartial{
// 				Path:        models.NewOptionalString(path),
// 				Interactive: models.NewOptionalBool(interactive),
// 			}
// 			if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 				_, err := scanner.CreatorUpdater.UpdatePartial(ctx, s.ID, scenePartial)
// 				return err
// 			}); err != nil {
// 				return nil, err
// 			}

// 			scanner.makeScreenshots(path, nil, sceneHash)
// 			scanner.PluginCache.ExecutePostHooks(ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
// 		}
// 	} else {
// 		logger.Infof("%s doesn't exist. Creating new item...", path)
// 		currentTime := time.Now()

// 		videoFile, err := scanner.VideoFileCreator.NewVideoFile(path)
// 		if err != nil {
// 			return nil, err
// 		}

// 		title := filepath.Base(path)
// 		if scanner.StripFileExtension {
// 			title = stripExtension(title)
// 		}

// 		if scanner.UseFileMetadata && videoFile.Title != "" {
// 			title = videoFile.Title
// 		}

// 		newScene := models.Scene{
// 			Path:        path,
// 			FileModTime: &scanned.FileModTime,
// 			Title:       title,
// 			CreatedAt:   currentTime,
// 			UpdatedAt:   currentTime,
// 			Interactive: interactive,
// 		}

// 		if checksum != "" {
// 			newScene.Checksum = &checksum
// 		}
// 		if oshash != "" {
// 			newScene.OSHash = &oshash
// 		}

// 		if err := videoFileToScene(&newScene, videoFile); err != nil {
// 			return nil, err
// 		}

// 		if scanner.UseFileMetadata {
// 			newScene.Details = videoFile.Comment
// 			d := models.SQLiteDate{}
// 			_ = d.Scan(videoFile.CreationTime)
// 			newScene.Date = d.DatePtr()
// 		}

// 		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
// 			return scanner.CreatorUpdater.Create(ctx, &newScene)
// 		}); err != nil {
// 			return nil, err
// 		}

// 		retScene = &newScene

// 		scanner.makeScreenshots(path, videoFile, sceneHash)
// 		scanner.PluginCache.ExecutePostHooks(ctx, retScene.ID, plugin.SceneCreatePost, nil, nil)
// 	}

// 	return retScene, nil
// }

// func stripExtension(path string) string {
// 	ext := filepath.Ext(path)
// 	return strings.TrimSuffix(path, ext)
// }

// func videoFileToScene(s *models.Scene, videoFile *ffmpeg.VideoFile) error {
// 	container, err := ffmpeg.MatchContainer(videoFile.Container, s.Path)
// 	if err != nil {
// 		return fmt.Errorf("matching container: %w", err)
// 	}

// 	s.Duration = &videoFile.Duration
// 	s.VideoCodec = &videoFile.VideoCodec
// 	s.AudioCodec = &videoFile.AudioCodec
// 	containerStr := string(container)
// 	s.Format = &containerStr
// 	s.Width = &videoFile.Width
// 	s.Height = &videoFile.Height
// 	s.Framerate = &videoFile.FrameRate
// 	s.Bitrate = &videoFile.Bitrate
// 	size := strconv.FormatInt(videoFile.Size, 10)
// 	s.Size = &size

// 	return nil
// }

// func (h *ScanHandler) makeScreenshots(ctx context.Context, scene *models.Scene, f *file.VideoFile) {
// 	checksum := scene.GetHash()
// 	thumbPath := h.Paths.Scene.GetThumbnailScreenshotPath(checksum)
// 	normalPath := h.Paths.Scene.GetScreenshotPath(checksum)

// 	thumbExists, _ := fsutil.FileExists(thumbPath)
// 	normalExists, _ := fsutil.FileExists(normalPath)

// 	if thumbExists && normalExists {
// 		return
// 	}

// 	if !thumbExists {
// 		logger.Debugf("Creating thumbnail for %s", f.Path)
// 		if err := h.Screenshotter.GenerateThumbnail(ctx, probeResult, checksum); err != nil {
// 			logger.Errorf("Error creating thumbnail for %s: %v", err)
// 		}
// 	}

// 	if !normalExists {
// 		logger.Debugf("Creating screenshot for %s", f.Path)
// 		if err := h.Screenshotter.GenerateScreenshot(ctx, probeResult, checksum); err != nil {
// 			logger.Errorf("Error creating screenshot for %s: %v", err)
// 		}
// 	}
// }

// func getInteractive(path string) bool {
// 	_, err := os.Stat(GetFunscriptPath(path))
// 	return err == nil
// }
