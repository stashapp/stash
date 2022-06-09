package manager

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type cleanJob struct {
	txnManager   Repository
	input        CleanMetadataInput
	sceneService SceneService
	imageService ImageService
	scanSubs     *subscriptionManager
}

func (j *cleanJob) Execute(ctx context.Context, progress *job.Progress) {
	logger.Infof("Starting cleaning of tracked files")
	if j.input.DryRun {
		logger.Infof("Running in Dry Mode")
	}

	r := j.txnManager

	if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		total, err := j.getCount(ctx, r)
		if err != nil {
			return fmt.Errorf("error getting count: %w", err)
		}

		progress.SetTotal(total)

		return nil
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	if err := j.processFiles(ctx, progress); err != nil {
		logger.Errorf("error cleaning scenes: %w", err)
		return
	}

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	j.scanSubs.notify()
	logger.Info("Finished Cleaning")
}

func (j *cleanJob) getCount(ctx context.Context, r Repository) (int, error) {
	fileFilter := models.PathsFileFilter(j.input.Paths)
	fileResult, err := j.txnManager.File.Query(ctx, models.FileQueryOptions{
		QueryOptions: models.QueryOptions{
			Count: true,
		},
		FileFilter: fileFilter,
	})
	if err != nil {
		return 0, err
	}

	return fileResult.Count, nil
}

type deleteSet struct {
	list []file.ID
	set  map[file.ID]string
}

func newDeleteSet() deleteSet {
	return deleteSet{
		set: make(map[file.ID]string),
	}
}

func (s *deleteSet) add(id file.ID, path string) {
	if _, ok := s.set[id]; !ok {
		s.list = append(s.list, id)
		s.set[id] = path
	}
}

func (s *deleteSet) len() int {
	return len(s.list)
}

func (j *cleanJob) processFiles(ctx context.Context, progress *job.Progress) error {
	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	fileFilter := models.PathsFileFilter(j.input.Paths)
	sort := "path"
	findFilter.Sort = &sort

	toDelete := newDeleteSet()

	more := true
	if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		for more {
			if job.IsCancelled(ctx) {
				return nil
			}

			files, err := j.fileQuery(ctx, fileFilter, findFilter)
			if err != nil {
				return fmt.Errorf("error querying for file: %w", err)
			}

			for _, f := range files {
				path := f.Base().Path
				err = nil
				progress.ExecuteTask(fmt.Sprintf("Assessing file %s for clean", path), func() {
					if j.shouldCleanFile(f) {
						// add contained files first
						fileID := f.Base().ID
						var containedFiles []file.File
						containedFiles, err = j.txnManager.File.FindByZipFileID(ctx, fileID)
						if err != nil {
							err = fmt.Errorf("error finding contained files: %w", err)
							return
						}

						for _, cf := range containedFiles {
							toDelete.add(cf.Base().ID, cf.Base().Path)
						}

						toDelete.add(f.Base().ID, f.Base().Path)
					} else {
						// increment progress, no further processing
						progress.Increment()
					}
				})
				if err != nil {
					return err
				}
			}

			if len(files) != batchSize {
				more = false
			} else {
				*findFilter.Page++
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if j.input.DryRun && toDelete.len() > 0 {
		// add progress for scenes that would've been deleted
		progress.AddProcessed(toDelete.len())
	}

	if !j.input.DryRun && len(toDelete.list) > 0 {
		progress.ExecuteTask(fmt.Sprintf("Cleaning %d files", toDelete.len()), func() {
			for _, fileID := range toDelete.list {
				if job.IsCancelled(ctx) {
					return
				}

				j.deleteFile(ctx, fileID, toDelete.set[fileID])

				progress.Increment()
			}
		})
	}

	return nil
}

func (j *cleanJob) fileQuery(ctx context.Context, fileFilter *models.FileFilterType, findFilter *models.FindFilterType) ([]file.File, error) {
	result, err := j.txnManager.File.Query(ctx, models.FileQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      false,
		},
		FileFilter: fileFilter,
	})
	if err != nil {
		return nil, err
	}

	files, err := result.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (j *cleanJob) shouldClean(f file.File) bool {
	fileExists, _ := f.Base().Exists()
	path := f.Base().Path

	// // #1102 - clean anything in generated path
	generatedPath := config.GetInstance().GetGeneratedPath()
	if !fileExists || getStashFromPath(path) == nil || fsutil.IsPathInDir(generatedPath, path) {
		logger.Infof("File not found. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (j *cleanJob) shouldCleanFile(f file.File) bool {
	if j.shouldClean(f) {
		return true
	}

	path := f.Base().Path

	config := config.GetInstance()
	switch {
	case fsutil.MatchExtension(path, config.GetVideoExtensions()):
		return j.shouldCleanVideoFile(path)
	case fsutil.MatchExtension(path, config.GetImageExtensions()):
		return j.shouldCleanImage(path)
	case fsutil.MatchExtension(path, config.GetGalleryExtensions()):
		return j.shouldCleanGallery(path)
	default:
		logger.Infof("File extension does not any media extensions. Marking to clean: \"%s\"", path)
		return true
	}
}

func (j *cleanJob) shouldCleanVideoFile(path string) bool {
	stash := getStashFromPath(path)
	if stash.ExcludeVideo {
		logger.Infof("File in stash library that excludes video. Marking to clean: \"%s\"", path)
		return true
	}

	config := config.GetInstance()

	if matchFile(path, config.GetExcludes()) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (j *cleanJob) shouldCleanGallery(path string) bool {
	stash := getStashFromPath(path)
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Marking to clean: \"%s\"", path)
		return true
	}

	config := config.GetInstance()
	if matchFile(path, config.GetImageExcludes()) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (j *cleanJob) shouldCleanImage(path string) bool {
	stash := getStashFromPath(path)
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Marking to clean: \"%s\"", path)
		return true
	}

	config := config.GetInstance()

	if matchFile(path, config.GetImageExcludes()) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (j *cleanJob) deleteFile(ctx context.Context, fileID file.ID, fn string) {
	// delete associated objects
	fileDeleter := file.NewDeleter()
	if err := j.txnManager.WithTxn(ctx, func(ctx context.Context) error {
		repo := j.txnManager
		qb := repo.File

		if err := j.deleteRelatedObjects(ctx, fileDeleter, fileID); err != nil {
			return err
		}

		return qb.Destroy(ctx, fileID)
	}); err != nil {
		fileDeleter.Rollback()

		logger.Errorf("Error deleting file %q from database: %s", fn, err.Error())
		return
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	// FIXME
	// GetInstance().PluginCache.ExecutePostHooks(ctx, sceneID, plugin.SceneDestroyPost, plugin.SceneDestroyInput{
	// 	Checksum: checksum,
	// 	OSHash:   oshash,
	// 	Path:     s.Path(),
	// }, nil)
}

func (j *cleanJob) deleteRelatedObjects(ctx context.Context, fileDeleter *file.Deleter, fileID file.ID) error {
	if err := j.deleteRelatedScenes(ctx, fileDeleter, fileID); err != nil {
		return err
	}
	if err := j.deleteRelatedGalleries(ctx, fileID); err != nil {
		return err
	}
	if err := j.deleteRelatedImages(ctx, fileDeleter, fileID); err != nil {
		return err
	}

	return nil
}

func (j *cleanJob) deleteRelatedScenes(ctx context.Context, fileDeleter *file.Deleter, fileID file.ID) error {
	sceneQB := j.txnManager.Scene
	scenes, err := sceneQB.FindByFileID(ctx, fileID)
	if err != nil {
		return err
	}

	fileNamingAlgo := GetInstance().Config.GetVideoFileNamingAlgorithm()

	sceneFileDeleter := &scene.FileDeleter{
		Deleter:        fileDeleter,
		FileNamingAlgo: fileNamingAlgo,
		Paths:          GetInstance().Paths,
	}

	for _, scene := range scenes {
		// only delete if the scene has no other files
		if len(scene.Files) <= 1 {
			logger.Infof("Deleting scene %q since it has no other related files", scene.GetTitle())
			if err := j.sceneService.Destroy(ctx, scene, sceneFileDeleter, true, false); err != nil {
				return err
			}
			// FIXME
			// checksum := s.Checksum()
			// oshash := s.OSHash()

			// GetInstance().PluginCache.ExecutePostHooks(ctx, sceneID, plugin.SceneDestroyPost, plugin.SceneDestroyInput{
			// 	Checksum: checksum,
			// 	OSHash:   oshash,
			// 	Path:     s.Path(),
			// }, nil)
		}
	}

	return nil
}

func (j *cleanJob) deleteRelatedGalleries(ctx context.Context, fileID file.ID) error {
	qb := j.txnManager.Gallery
	galleries, err := qb.FindByFileID(ctx, fileID)
	if err != nil {
		return err
	}

	for _, g := range galleries {
		// only delete if the scene has no other files
		if len(g.Files) <= 1 {
			logger.Infof("Deleting gallery %q since it has no other related files", g.GetTitle())
			if err := qb.Destroy(ctx, g.ID); err != nil {
				return err
			}
			// FIXME
			// GetInstance().PluginCache.ExecutePostHooks(ctx, galleryID, plugin.GalleryDestroyPost, plugin.GalleryDestroyInput{
			// 	Checksum: g.Checksum(),
			// 	Path:     g.Path(),
			// }, nil)
		}
	}

	return nil
}

func (j *cleanJob) deleteRelatedImages(ctx context.Context, fileDeleter *file.Deleter, fileID file.ID) error {
	imageQB := j.txnManager.Image
	images, err := imageQB.FindByFileID(ctx, fileID)
	if err != nil {
		return err
	}

	imageFileDeleter := &image.FileDeleter{
		Deleter: fileDeleter,
		Paths:   GetInstance().Paths,
	}

	for _, i := range images {
		if len(i.Files) <= 1 {
			logger.Infof("Deleting image %q since it has no other related files", i.GetTitle())
			if err := j.imageService.Destroy(ctx, i, imageFileDeleter, true, false); err != nil {
				return err
			}
			// FIXME
			// GetInstance().PluginCache.ExecutePostHooks(ctx, imageID, plugin.ImageDestroyPost, plugin.ImageDestroyInput{
			// 	Checksum: i.Checksum(),
			// 	Path:     i.Path(),
			// }, nil)
		}
	}

	return nil
}

func getStashFromPath(pathToCheck string) *config.StashConfig {
	for _, f := range config.GetInstance().GetStashPaths() {
		if fsutil.IsPathInDir(f.Path, filepath.Dir(pathToCheck)) {
			return f
		}
	}
	return nil
}

func getStashFromDirPath(pathToCheck string) *config.StashConfig {
	for _, f := range config.GetInstance().GetStashPaths() {
		if fsutil.IsPathInDir(f.Path, pathToCheck) {
			return f
		}
	}
	return nil
}
