package manager

import (
	"archive/zip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/remeh/sizedwaitgroup"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/utils"
)

type ScanJob struct {
	txnManager    models.TransactionManager
	input         models.ScanMetadataInput
	subscriptions *subscriptionManager
}

func (j *ScanJob) Execute(ctx context.Context, progress *job.Progress) {
	input := j.input
	paths := getScanPaths(input.Paths)

	var total *int
	var newFiles *int
	progress.ExecuteTask("Counting files to scan...", func() {
		total, newFiles = j.neededScan(ctx, paths)
	})

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	if total == nil || newFiles == nil {
		logger.Infof("Taking too long to count content. Skipping...")
		logger.Infof("Starting scan")
	} else {
		logger.Infof("Starting scan of %d files. %d New files found", *total, *newFiles)
	}

	start := time.Now()
	config := config.GetInstance()
	parallelTasks := config.GetParallelTasksWithAutoDetection()
	logger.Infof("Scan started with %d parallel tasks", parallelTasks)
	wg := sizedwaitgroup.New(parallelTasks)

	if total != nil {
		progress.SetTotal(*total)
	}

	fileNamingAlgo := config.GetVideoFileNamingAlgorithm()
	calculateMD5 := config.IsCalculateMD5()

	stoppingErr := errors.New("stopping")
	var err error

	var galleries []string

	for _, sp := range paths {
		csFs, er := utils.IsFsPathCaseSensitive(sp.Path)
		if er != nil {
			logger.Warnf("Cannot determine fs case sensitivity: %s", er.Error())
		}

		err = walkFilesToScan(sp, func(path string, info os.FileInfo, err error) error {
			if job.IsCancelled(ctx) {
				return stoppingErr
			}

			if isGallery(path) {
				galleries = append(galleries, path)
			}

			instance.Paths.Generated.EnsureTmpDir()

			wg.Add()
			task := ScanTask{
				TxnManager:           j.txnManager,
				FilePath:             path,
				FileInfo:             info,
				UseFileMetadata:      utils.IsTrue(input.UseFileMetadata),
				StripFileExtension:   utils.IsTrue(input.StripFileExtension),
				fileNamingAlgorithm:  fileNamingAlgo,
				calculateMD5:         calculateMD5,
				GeneratePreview:      utils.IsTrue(input.ScanGeneratePreviews),
				GenerateImagePreview: utils.IsTrue(input.ScanGenerateImagePreviews),
				GenerateSprite:       utils.IsTrue(input.ScanGenerateSprites),
				GeneratePhash:        utils.IsTrue(input.ScanGeneratePhashes),
				progress:             progress,
				CaseSensitiveFs:      csFs,
				ctx:                  ctx,
			}

			go func() {
				task.Start(&wg)
				progress.Increment()
			}()

			return nil
		})

		if err == stoppingErr {
			logger.Info("Stopping due to user request")
			break
		}

		if err != nil {
			logger.Errorf("Error encountered scanning files: %s", err.Error())
			break
		}
	}

	wg.Wait()
	instance.Paths.Generated.EmptyTmpDir()
	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Scan finished (%s)", elapsed))

	if job.IsCancelled(ctx) || err != nil {
		return
	}

	progress.ExecuteTask("Associating galleries", func() {
		for _, path := range galleries {
			wg.Add()
			task := ScanTask{
				TxnManager:      j.txnManager,
				FilePath:        path,
				UseFileMetadata: false,
			}

			go task.associateGallery(&wg)
			wg.Wait()
		}
		logger.Info("Finished gallery association")
	})

	j.subscriptions.notify()
}

func (j *ScanJob) neededScan(ctx context.Context, paths []*models.StashConfig) (total *int, newFiles *int) {
	const timeout = 90 * time.Second

	// create a control channel through which to signal the counting loop when the timeout is reached
	chTimeout := time.After(timeout)

	logger.Infof("Counting files to scan...")

	t := 0
	n := 0

	timeoutErr := errors.New("timed out")

	for _, sp := range paths {
		err := walkFilesToScan(sp, func(path string, info os.FileInfo, err error) error {
			t++
			task := ScanTask{FilePath: path, TxnManager: j.txnManager}
			if !task.doesPathExist() {
				n++
			}

			//check for timeout
			select {
			case <-chTimeout:
				return timeoutErr
			default:
			}

			// check stop
			if job.IsCancelled(ctx) {
				return timeoutErr
			}

			return nil
		})

		if err == timeoutErr {
			// timeout should return nil counts
			return nil, nil
		}

		if err != nil {
			logger.Errorf("Error encountered counting files to scan: %s", err.Error())
			return nil, nil
		}
	}

	return &t, &n
}

type ScanTask struct {
	ctx                  context.Context
	TxnManager           models.TransactionManager
	FilePath             string
	FileInfo             os.FileInfo
	UseFileMetadata      bool
	StripFileExtension   bool
	calculateMD5         bool
	fileNamingAlgorithm  models.HashAlgorithm
	GenerateSprite       bool
	GeneratePhash        bool
	GeneratePreview      bool
	GenerateImagePreview bool
	zipGallery           *models.Gallery
	progress             *job.Progress
	CaseSensitiveFs      bool
}

func (t *ScanTask) Start(wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()

	var s *models.Scene

	t.progress.ExecuteTask("Scanning "+t.FilePath, func() {
		if isGallery(t.FilePath) {
			t.scanGallery()
		} else if isVideo(t.FilePath) {
			s = t.scanScene()
		} else if isImage(t.FilePath) {
			t.scanImage()
		}
	})

	if s != nil {
		iwg := sizedwaitgroup.New(2)

		if t.GenerateSprite {
			iwg.Add()

			go t.progress.ExecuteTask(fmt.Sprintf("Generating sprites for %s", t.FilePath), func() {
				taskSprite := GenerateSpriteTask{
					Scene:               *s,
					Overwrite:           false,
					fileNamingAlgorithm: t.fileNamingAlgorithm,
				}
				taskSprite.Start(&iwg)
			})
		}

		if t.GeneratePhash {
			iwg.Add()

			go t.progress.ExecuteTask(fmt.Sprintf("Generating phash for %s", t.FilePath), func() {
				taskPhash := GeneratePhashTask{
					Scene:               *s,
					fileNamingAlgorithm: t.fileNamingAlgorithm,
					txnManager:          t.TxnManager,
				}
				taskPhash.Start(&iwg)
			})
		}

		if t.GeneratePreview {
			iwg.Add()

			go t.progress.ExecuteTask(fmt.Sprintf("Generating preview for %s", t.FilePath), func() {
				config := config.GetInstance()
				var previewSegmentDuration = config.GetPreviewSegmentDuration()
				var previewSegments = config.GetPreviewSegments()
				var previewExcludeStart = config.GetPreviewExcludeStart()
				var previewExcludeEnd = config.GetPreviewExcludeEnd()
				var previewPresent = config.GetPreviewPreset()

				// NOTE: the reuse of this model like this is painful.
				previewOptions := models.GeneratePreviewOptionsInput{
					PreviewSegments:        &previewSegments,
					PreviewSegmentDuration: &previewSegmentDuration,
					PreviewExcludeStart:    &previewExcludeStart,
					PreviewExcludeEnd:      &previewExcludeEnd,
					PreviewPreset:          &previewPresent,
				}

				taskPreview := GeneratePreviewTask{
					Scene:               *s,
					ImagePreview:        t.GenerateImagePreview,
					Options:             previewOptions,
					Overwrite:           false,
					fileNamingAlgorithm: t.fileNamingAlgorithm,
				}
				taskPreview.Start(wg)
			})
		}

		iwg.Wait()
	}
}

func (t *ScanTask) scanGallery() {
	var g *models.Gallery
	images := 0
	scanImages := false

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		g, err = r.Gallery().FindByPath(t.FilePath)

		if g != nil && err != nil {
			images, err = r.Image().CountByGalleryID(g.ID)
			if err != nil {
				return fmt.Errorf("error getting images for zip gallery %s: %s", t.FilePath, err.Error())
			}
		}

		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	fileModTime, err := t.getFileModTime()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if g != nil {
		// We already have this item in the database, keep going

		// if file mod time is not set, set it now
		if !g.FileModTime.Valid {
			// we will also need to rescan the zip contents
			scanImages = true
			logger.Infof("setting file modification time on %s", t.FilePath)

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Gallery()
				if _, err := gallery.UpdateFileModTime(qb, g.ID, models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				}); err != nil {
					return err
				}

				// update our copy of the gallery
				var err error
				g, err = qb.Find(g.ID)
				return err
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}

		// if the mod time of the zip file is different than that of the associated
		// gallery, then recalculate the checksum
		modified := t.isFileModified(fileModTime, g.FileModTime)
		if modified {
			scanImages = true
			logger.Infof("%s has been updated: rescanning", t.FilePath)

			// update the checksum and the modification time
			checksum, err := t.calculateChecksum()
			if err != nil {
				logger.Error(err.Error())
				return
			}

			currentTime := time.Now()
			galleryPartial := models.GalleryPartial{
				ID:       g.ID,
				Checksum: &checksum,
				FileModTime: &models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				},
				UpdatedAt: &models.SQLiteTimestamp{Timestamp: currentTime},
			}

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				_, err := r.Gallery().UpdatePartial(galleryPartial)
				return err
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}

		// scan the zip files if the gallery has no images
		scanImages = scanImages || images == 0
	} else {
		// Ignore directories.
		if isDir, _ := utils.DirExists(t.FilePath); isDir {
			return
		}

		checksum, err := t.calculateChecksum()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			qb := r.Gallery()
			g, _ = qb.FindByChecksum(checksum)
			if g != nil {
				exists, _ := utils.FileExists(g.Path.String)
				if !t.CaseSensitiveFs {
					// #1426 - if file exists but is a case-insensitive match for the
					// original filename, then treat it as a move
					if exists && strings.EqualFold(t.FilePath, g.Path.String) {
						exists = false
					}
				}

				if exists {
					logger.Infof("%s already exists.  Duplicate of %s ", t.FilePath, g.Path.String)
				} else {
					logger.Infof("%s already exists.  Updating path...", t.FilePath)
					g.Path = sql.NullString{
						String: t.FilePath,
						Valid:  true,
					}
					g, err = qb.Update(*g)
					if err != nil {
						return err
					}

					GetInstance().PluginCache.ExecutePostHooks(t.ctx, g.ID, plugin.GalleryUpdatePost, nil, nil)
				}
			} else {
				currentTime := time.Now()

				newGallery := models.Gallery{
					Checksum: checksum,
					Zip:      true,
					Path: sql.NullString{
						String: t.FilePath,
						Valid:  true,
					},
					FileModTime: models.NullSQLiteTimestamp{
						Timestamp: fileModTime,
						Valid:     true,
					},
					Title: sql.NullString{
						String: utils.GetNameFromPath(t.FilePath, t.StripFileExtension),
						Valid:  true,
					},
					CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
					UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
				}

				// don't create gallery if it has no images
				if countImagesInZip(t.FilePath) > 0 {
					// only warn when creating the gallery
					ok, err := utils.IsZipFileUncompressed(t.FilePath)
					if err == nil && !ok {
						logger.Warnf("%s is using above store (0) level compression.", t.FilePath)
					}

					logger.Infof("%s doesn't exist.  Creating new item...", t.FilePath)
					g, err = qb.Create(newGallery)
					if err != nil {
						return err
					}
					scanImages = true

					GetInstance().PluginCache.ExecutePostHooks(t.ctx, g.ID, plugin.GalleryCreatePost, nil, nil)
				}
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			return
		}
	}

	if g != nil {
		if scanImages {
			t.scanZipImages(g)
		} else {
			// in case thumbnails have been deleted, regenerate them
			t.regenerateZipImages(g)
		}
	}
}

func (t *ScanTask) getFileModTime() (time.Time, error) {
	fi, err := os.Stat(t.FilePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("error performing stat on %s: %s", t.FilePath, err.Error())
	}

	ret := fi.ModTime()
	// truncate to seconds, since we don't store beyond that in the database
	ret = ret.Truncate(time.Second)

	return ret, nil
}

func (t *ScanTask) getInteractive() bool {
	_, err := os.Stat(utils.GetFunscriptPath(t.FilePath))
	return err == nil

}

func (t *ScanTask) isFileModified(fileModTime time.Time, modTime models.NullSQLiteTimestamp) bool {
	return !modTime.Timestamp.Equal(fileModTime)
}

// associates a gallery to a scene with the same basename
func (t *ScanTask) associateGallery(wg *sizedwaitgroup.SizedWaitGroup) {
	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Gallery()
		sqb := r.Scene()
		g, err := qb.FindByPath(t.FilePath)
		if err != nil {
			return err
		}

		if g == nil {
			// associate is run after scan is finished
			// should only happen if gallery is a directory or an io error occurs during hashing
			logger.Warnf("associate: gallery %s not found in DB", t.FilePath)
			return nil
		}

		basename := strings.TrimSuffix(t.FilePath, filepath.Ext(t.FilePath))
		var relatedFiles []string
		vExt := config.GetInstance().GetVideoExtensions()
		// make a list of media files that can be related to the gallery
		for _, ext := range vExt {
			related := basename + "." + ext
			// exclude gallery extensions from the related files
			if !isGallery(related) {
				relatedFiles = append(relatedFiles, related)
			}
		}
		for _, scenePath := range relatedFiles {
			scene, _ := sqb.FindByPath(scenePath)
			// found related Scene
			if scene != nil {
				sceneGalleries, _ := sqb.FindByGalleryID(g.ID) // check if gallery is already associated to the scene
				isAssoc := false
				for _, sg := range sceneGalleries {
					if scene.ID == sg.ID {
						isAssoc = true
						break
					}
				}
				if !isAssoc {
					logger.Infof("associate: Gallery %s is related to scene: %d", t.FilePath, scene.ID)
					if err := sqb.UpdateGalleries(scene.ID, []int{g.ID}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
	wg.Done()
}

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

	interactive := t.getInteractive()

	if s != nil {
		scanned, err := scanner.ScanExisting(s, t.FilePath, t.FileInfo)
		if err != nil {
			return logError(err)
		}

		if scanned.ContentsChanged() {
			config := config.GetInstance()
			oldHash := s.GetHash(config.GetVideoFileNamingAlgorithm())
			s, err = t.rescanScene(s, *scanned.New.FileModTime)
			if err != nil {
				return logError(err)
			}

			// Migrate any generated files if the hash has changed
			newHash := s.GetHash(config.GetVideoFileNamingAlgorithm())
			if newHash != oldHash {
				MigrateHash(oldHash, newHash)
			}
		} else if scanned.FileUpdated() || s.Interactive != interactive {
			// update fields as needed
			scenePartial := models.ScenePartial{
				ID:          s.ID,
				Interactive: &interactive,
			}

			scenePartial.SetFile(*scanned.New)

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Scene()
				// ensure no clashes of hashes
				if scenePartial.Checksum != nil {
					dupe, _ := qb.FindByChecksum(*scanned.New.Checksum)
					if dupe != nil {
						return fmt.Errorf("MD5 for file %s is the same as that of %s", t.FilePath, dupe.Path)
					}
				}

				if scenePartial.OSHash != nil {
					dupe, _ := qb.FindByOSHash(*scanned.New.OSHash)
					if dupe != nil {
						return fmt.Errorf("OSHash for file %s is the same as that of %s", t.FilePath, dupe.Path)
					}
				}

				_, err := qb.Update(scenePartial)
				return err
			}); err != nil {
				return logError(err)
			}
		}

		// We already have this item in the database
		// check for thumbnails,screenshots
		t.makeScreenshots(nil, s.GetHash(t.fileNamingAlgorithm))

		// check for container
		if !s.Format.Valid {
			videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
			if err != nil {
				return logError(err)
			}
			container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)
			logger.Infof("Adding container %s to file %s", container, t.FilePath)

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				_, err := scene.UpdateFormat(r.Scene(), s.ID, string(container))
				return err
			}); err != nil {
				return logError(err)
			}
		}

		return nil
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)

	// Override title to be filename if UseFileMetadata is false
	if !t.UseFileMetadata {
		videoFile.SetTitleFromPath(t.StripFileExtension)
	}

	file, err := scanner.ScanNew(t.FilePath, t.FileInfo)
	if err != nil {
		return logError(err)
	}

	var checksum string
	if file.Checksum != nil {
		checksum = *file.Checksum
	}
	oshash := *file.OSHash

	// check for scene by checksum and oshash - MD5 should be
	// redundant, but check both
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
				return logError(err)
			}

			GetInstance().PluginCache.ExecutePostHooks(t.ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
		}
	} else {
		logger.Infof("%s doesn't exist. Creating new item...", t.FilePath)
		currentTime := time.Now()
		newScene := models.Scene{
			Checksum:   sql.NullString{String: checksum, Valid: checksum != ""},
			OSHash:     sql.NullString{String: oshash, Valid: oshash != ""},
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
			Size:       sql.NullString{String: strconv.FormatInt(videoFile.Size, 10), Valid: true},
			FileModTime: models.NullSQLiteTimestamp{
				Timestamp: *file.FileModTime,
				Valid:     true,
			},
			CreatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			Interactive: interactive,
		}

		if t.UseFileMetadata {
			newScene.Details = sql.NullString{String: videoFile.Comment, Valid: true}
			newScene.Date = models.SQLiteDate{String: videoFile.CreationTime.Format("2006-01-02")}
		}

		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			var err error
			retScene, err = r.Scene().Create(newScene)
			return err
		}); err != nil {
			return logError(err)
		}

		GetInstance().PluginCache.ExecutePostHooks(t.ctx, retScene.ID, plugin.SceneCreatePost, nil, nil)
	}

	return retScene
}

func (t *ScanTask) rescanScene(s *models.Scene, fileModTime time.Time) (*models.Scene, error) {
	logger.Infof("%s has been updated: rescanning", t.FilePath)

	// update the oshash/checksum and the modification time
	logger.Infof("Calculating oshash for existing file %s ...", t.FilePath)
	oshash, err := utils.OSHashFromFilePath(t.FilePath)
	if err != nil {
		return nil, err
	}

	var checksum *sql.NullString
	if t.calculateMD5 {
		cs, err := t.calculateChecksum()
		if err != nil {
			return nil, err
		}

		checksum = &sql.NullString{
			String: cs,
			Valid:  true,
		}
	}

	// regenerate the file details as well
	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
	if err != nil {
		return nil, err
	}
	container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)

	currentTime := time.Now()
	scenePartial := models.ScenePartial{
		ID:       s.ID,
		Checksum: checksum,
		OSHash: &sql.NullString{
			String: oshash,
			Valid:  true,
		},
		Duration:   &sql.NullFloat64{Float64: videoFile.Duration, Valid: true},
		VideoCodec: &sql.NullString{String: videoFile.VideoCodec, Valid: true},
		AudioCodec: &sql.NullString{String: videoFile.AudioCodec, Valid: true},
		Format:     &sql.NullString{String: string(container), Valid: true},
		Width:      &sql.NullInt64{Int64: int64(videoFile.Width), Valid: true},
		Height:     &sql.NullInt64{Int64: int64(videoFile.Height), Valid: true},
		Framerate:  &sql.NullFloat64{Float64: videoFile.FrameRate, Valid: true},
		Bitrate:    &sql.NullInt64{Int64: videoFile.Bitrate, Valid: true},
		Size:       &sql.NullString{String: strconv.FormatInt(videoFile.Size, 10), Valid: true},
		FileModTime: &models.NullSQLiteTimestamp{
			Timestamp: fileModTime,
			Valid:     true,
		},
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: currentTime},
	}

	var ret *models.Scene
	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		var err error
		ret, err = r.Scene().Update(scenePartial)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	GetInstance().PluginCache.ExecutePostHooks(t.ctx, ret.ID, plugin.SceneUpdatePost, nil, nil)

	// leave the generated files as is - the scene file may have been moved
	// elsewhere

	return ret, nil
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

func (t *ScanTask) scanZipImages(zipGallery *models.Gallery) {
	err := walkGalleryZip(zipGallery.Path.String, func(file *zip.File) error {
		// copy this task and change the filename
		subTask := *t

		// filepath is the zip file and the internal file name, separated by a null byte
		subTask.FilePath = image.ZipFilename(zipGallery.Path.String, file.Name)
		subTask.zipGallery = zipGallery

		// run the subtask and wait for it to complete
		iwg := sizedwaitgroup.New(1)
		iwg.Add()
		subTask.Start(&iwg)
		return nil
	})
	if err != nil {
		logger.Warnf("failed to scan zip file images for %s: %s", zipGallery.Path.String, err.Error())
	}
}

func (t *ScanTask) regenerateZipImages(zipGallery *models.Gallery) {
	var images []*models.Image
	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		iqb := r.Image()

		var err error
		images, err = iqb.FindByGalleryID(zipGallery.ID)
		return err
	}); err != nil {
		logger.Warnf("failed to find gallery images: %s", err.Error())
		return
	}

	for _, img := range images {
		t.generateThumbnail(img)
	}
}

func (t *ScanTask) scanImage() {
	var i *models.Image

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		i, err = r.Image().FindByPath(t.FilePath)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	fileModTime, err := image.GetFileModTime(t.FilePath)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if i != nil {
		// if file mod time is not set, set it now
		if !i.FileModTime.Valid {
			logger.Infof("setting file modification time on %s", t.FilePath)

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Image()
				if _, err := image.UpdateFileModTime(qb, i.ID, models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				}); err != nil {
					return err
				}

				// update our copy of the gallery
				var err error
				i, err = qb.Find(i.ID)
				return err
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}

		// if the mod time of the file is different than that of the associated
		// image, then recalculate the checksum and regenerate the thumbnail
		modified := t.isFileModified(fileModTime, i.FileModTime)
		if modified {
			i, err = t.rescanImage(i, fileModTime)
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}

		// We already have this item in the database
		// check for thumbnails
		t.generateThumbnail(i)
	} else {
		// Ignore directories.
		if isDir, _ := utils.DirExists(t.FilePath); isDir {
			return
		}

		var checksum string

		logger.Infof("%s not found.  Calculating checksum...", t.FilePath)
		checksum, err = t.calculateImageChecksum()
		if err != nil {
			logger.Errorf("error calculating checksum for %s: %s", t.FilePath, err.Error())
			return
		}

		// check for scene by checksum and oshash - MD5 should be
		// redundant, but check both
		if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			var err error
			i, err = r.Image().FindByChecksum(checksum)
			return err
		}); err != nil {
			logger.Error(err.Error())
			return
		}

		if i != nil {
			exists := image.FileExists(i.Path)
			if !t.CaseSensitiveFs {
				// #1426 - if file exists but is a case-insensitive match for the
				// original filename, then treat it as a move
				if exists && strings.EqualFold(t.FilePath, i.Path) {
					exists = false
				}
			}

			if exists {
				logger.Infof("%s already exists.  Duplicate of %s ", image.PathDisplayName(t.FilePath), image.PathDisplayName(i.Path))
			} else {
				logger.Infof("%s already exists.  Updating path...", image.PathDisplayName(t.FilePath))
				imagePartial := models.ImagePartial{
					ID:   i.ID,
					Path: &t.FilePath,
				}

				if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
					_, err := r.Image().Update(imagePartial)
					return err
				}); err != nil {
					logger.Error(err.Error())
					return
				}

				GetInstance().PluginCache.ExecutePostHooks(t.ctx, i.ID, plugin.ImageUpdatePost, nil, nil)
			}
		} else {
			logger.Infof("%s doesn't exist.  Creating new item...", image.PathDisplayName(t.FilePath))
			currentTime := time.Now()
			newImage := models.Image{
				Checksum: checksum,
				Path:     t.FilePath,
				FileModTime: models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				},
				CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
				UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			}
			newImage.Title.String = image.GetFilename(&newImage, t.StripFileExtension)
			newImage.Title.Valid = true

			if err := image.SetFileDetails(&newImage); err != nil {
				logger.Error(err.Error())
				return
			}

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				var err error
				i, err = r.Image().Create(newImage)
				return err
			}); err != nil {
				logger.Error(err.Error())
				return
			}

			GetInstance().PluginCache.ExecutePostHooks(t.ctx, i.ID, plugin.ImageCreatePost, nil, nil)
		}

		if t.zipGallery != nil {
			// associate with gallery
			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				return gallery.AddImage(r.Gallery(), t.zipGallery.ID, i.ID)
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		} else if config.GetInstance().GetCreateGalleriesFromFolders() {
			// create gallery from folder or associate with existing gallery
			logger.Infof("Associating image %s with folder gallery", i.Path)
			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				return t.associateImageWithFolderGallery(i.ID, r.Gallery())
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}
	}

	if i != nil {
		t.generateThumbnail(i)
	}
}

func (t *ScanTask) rescanImage(i *models.Image, fileModTime time.Time) (*models.Image, error) {
	logger.Infof("%s has been updated: rescanning", t.FilePath)

	oldChecksum := i.Checksum

	// update the checksum and the modification time
	checksum, err := t.calculateImageChecksum()
	if err != nil {
		return nil, err
	}

	// regenerate the file details as well
	fileDetails, err := image.GetFileDetails(t.FilePath)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	imagePartial := models.ImagePartial{
		ID:       i.ID,
		Checksum: &checksum,
		Width:    &fileDetails.Width,
		Height:   &fileDetails.Height,
		Size:     &fileDetails.Size,
		FileModTime: &models.NullSQLiteTimestamp{
			Timestamp: fileModTime,
			Valid:     true,
		},
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: currentTime},
	}

	var ret *models.Image
	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		var err error
		ret, err = r.Image().Update(imagePartial)
		return err
	}); err != nil {
		return nil, err
	}

	// remove the old thumbnail if the checksum changed - we'll regenerate it
	if oldChecksum != checksum {
		err = os.Remove(GetInstance().Paths.Generated.GetThumbnailPath(oldChecksum, models.DefaultGthumbWidth)) // remove cache dir of gallery
		if err != nil {
			logger.Errorf("Error deleting thumbnail image: %s", err)
		}
	}

	GetInstance().PluginCache.ExecutePostHooks(t.ctx, ret.ID, plugin.ImageUpdatePost, nil, nil)

	return ret, nil
}

func (t *ScanTask) associateImageWithFolderGallery(imageID int, qb models.GalleryReaderWriter) error {
	// find a gallery with the path specified
	path := filepath.Dir(t.FilePath)
	g, err := qb.FindByPath(path)
	if err != nil {
		return err
	}

	if g == nil {
		checksum := utils.MD5FromString(path)

		// create the gallery
		currentTime := time.Now()

		newGallery := models.Gallery{
			Checksum: checksum,
			Path: sql.NullString{
				String: path,
				Valid:  true,
			},
			CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			Title: sql.NullString{
				String: utils.GetNameFromPath(path, false),
				Valid:  true,
			},
		}

		logger.Infof("Creating gallery for folder %s", path)
		g, err = qb.Create(newGallery)
		if err != nil {
			return err
		}
	}

	// associate image with gallery
	err = gallery.AddImage(qb, g.ID, imageID)
	return err
}

func (t *ScanTask) generateThumbnail(i *models.Image) {
	thumbPath := GetInstance().Paths.Generated.GetThumbnailPath(i.Checksum, models.DefaultGthumbWidth)
	exists, _ := utils.FileExists(thumbPath)
	if exists {
		return
	}

	srcImage, err := image.GetSourceImage(i)
	if err != nil {
		logger.Errorf("error reading image %s: %s", i.Path, err.Error())
		return
	}

	if image.ThumbnailNeeded(srcImage, models.DefaultGthumbWidth) {
		data, err := image.GetThumbnail(srcImage, models.DefaultGthumbWidth)
		if err != nil {
			logger.Errorf("error getting thumbnail for image %s: %s", i.Path, err.Error())
			return
		}

		err = utils.WriteFile(thumbPath, data)
		if err != nil {
			logger.Errorf("error writing thumbnail for image %s: %s", i.Path, err)
		}
	}
}

func (t *ScanTask) calculateChecksum() (string, error) {
	logger.Infof("Calculating checksum for %s...", t.FilePath)
	checksum, err := utils.MD5FromFilePath(t.FilePath)
	if err != nil {
		return "", err
	}
	logger.Debugf("Checksum calculated: %s", checksum)
	return checksum, nil
}

func (t *ScanTask) calculateImageChecksum() (string, error) {
	logger.Infof("Calculating checksum for %s...", image.PathDisplayName(t.FilePath))
	// uses image.CalculateMD5 to read files in zips
	checksum, err := image.CalculateMD5(t.FilePath)
	if err != nil {
		return "", err
	}
	logger.Debugf("Checksum calculated: %s", checksum)
	return checksum, nil
}

func (t *ScanTask) doesPathExist() bool {
	config := config.GetInstance()
	vidExt := config.GetVideoExtensions()
	imgExt := config.GetImageExtensions()
	gExt := config.GetGalleryExtensions()

	ret := false
	t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		if matchExtension(t.FilePath, gExt) {
			gallery, _ := r.Gallery().FindByPath(t.FilePath)
			if gallery != nil {
				ret = true
			}
		} else if matchExtension(t.FilePath, vidExt) {
			s, _ := r.Scene().FindByPath(t.FilePath)
			if s != nil {
				ret = true
			}
		} else if matchExtension(t.FilePath, imgExt) {
			i, _ := r.Image().FindByPath(t.FilePath)
			if i != nil {
				ret = true
			}
		}

		return nil
	})

	return ret
}

func walkFilesToScan(s *models.StashConfig, f filepath.WalkFunc) error {
	config := config.GetInstance()
	vidExt := config.GetVideoExtensions()
	imgExt := config.GetImageExtensions()
	gExt := config.GetGalleryExtensions()
	excludeVidRegex := generateRegexps(config.GetExcludes())
	excludeImgRegex := generateRegexps(config.GetImageExcludes())

	// don't scan zip images directly
	if image.IsZipPath(s.Path) {
		logger.Warnf("Cannot rescan zip image %s. Rescan zip gallery instead.", s.Path)
		return nil
	}

	generatedPath := config.GetGeneratedPath()

	return utils.SymWalk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warnf("error scanning %s: %s", path, err.Error())
			return nil
		}

		if info.IsDir() {
			// #1102 - ignore files in generated path
			if utils.IsPathInDir(generatedPath, path) {
				return filepath.SkipDir
			}

			// shortcut: skip the directory entirely if it matches both exclusion patterns
			// add a trailing separator so that it correctly matches against patterns like path/.*
			pathExcludeTest := path + string(filepath.Separator)
			if (s.ExcludeVideo || matchFileRegex(pathExcludeTest, excludeVidRegex)) && (s.ExcludeImage || matchFileRegex(pathExcludeTest, excludeImgRegex)) {
				return filepath.SkipDir
			}

			return nil
		}

		if !s.ExcludeVideo && matchExtension(path, vidExt) && !matchFileRegex(path, excludeVidRegex) {
			return f(path, info, err)
		}

		if !s.ExcludeImage {
			if (matchExtension(path, imgExt) || matchExtension(path, gExt)) && !matchFileRegex(path, excludeImgRegex) {
				return f(path, info, err)
			}
		}

		return nil
	})
}
