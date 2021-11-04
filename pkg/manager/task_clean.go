package manager

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/utils"
)

type cleanJob struct {
	txnManager models.TransactionManager
	input      models.CleanMetadataInput
	scanSubs   *subscriptionManager
}

func (j *cleanJob) Execute(ctx context.Context, progress *job.Progress) {
	logger.Infof("Starting cleaning of tracked files")
	if j.input.DryRun {
		logger.Infof("Running in Dry Mode")
	}

	if err := j.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		total, err := j.getCount(r)
		if err != nil {
			return fmt.Errorf("error getting count: %w", err)
		}

		progress.SetTotal(total)

		if job.IsCancelled(ctx) {
			return nil
		}

		if err := j.processScenes(ctx, progress, r.Scene()); err != nil {
			return fmt.Errorf("error cleaning scenes: %w", err)
		}
		if err := j.processImages(ctx, progress, r.Image()); err != nil {
			return fmt.Errorf("error cleaning images: %w", err)
		}
		if err := j.processGalleries(ctx, progress, r.Gallery(), r.Image()); err != nil {
			return fmt.Errorf("error cleaning galleries: %w", err)
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
		return
	}

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	j.scanSubs.notify()
	logger.Info("Finished Cleaning")
}

func (j *cleanJob) getCount(r models.ReaderRepository) (int, error) {
	sceneCount, err := r.Scene().Count()
	if err != nil {
		return 0, err
	}

	imageCount, err := r.Image().Count()
	if err != nil {
		return 0, err
	}

	galleryCount, err := r.Gallery().Count()
	if err != nil {
		return 0, err
	}

	return sceneCount + imageCount + galleryCount, nil
}

func (j *cleanJob) processScenes(ctx context.Context, progress *job.Progress, qb models.SceneReader) error {
	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	sort := "path"
	findFilter.Sort = &sort

	var toDelete []int

	more := true
	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		scenes, err := scene.Query(qb, nil, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for scenes: %w", err)
		}

		for _, scene := range scenes {
			progress.ExecuteTask(fmt.Sprintf("Assessing scene %s for clean", scene.Path), func() {
				if j.shouldCleanScene(scene) {
					toDelete = append(toDelete, scene.ID)
				} else {
					// increment progress, no further processing
					progress.Increment()
				}
			})
		}

		if len(scenes) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	if j.input.DryRun && len(toDelete) > 0 {
		// add progress for scenes that would've been deleted
		progress.AddProcessed(len(toDelete))
	}

	fileNamingAlgorithm := instance.Config.GetVideoFileNamingAlgorithm()

	if !j.input.DryRun && len(toDelete) > 0 {
		progress.ExecuteTask(fmt.Sprintf("Cleaning %d scenes", len(toDelete)), func() {
			for _, sceneID := range toDelete {
				if job.IsCancelled(ctx) {
					return
				}

				j.deleteScene(ctx, fileNamingAlgorithm, sceneID)

				progress.Increment()
			}
		})
	}

	return nil
}

func (j *cleanJob) processGalleries(ctx context.Context, progress *job.Progress, qb models.GalleryReader, iqb models.ImageReader) error {
	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	sort := "path"
	findFilter.Sort = &sort

	var toDelete []int

	more := true
	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		galleries, _, err := qb.Query(nil, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for galleries: %w", err)
		}

		for _, gallery := range galleries {
			progress.ExecuteTask(fmt.Sprintf("Assessing gallery %s for clean", gallery.GetTitle()), func() {
				if j.shouldCleanGallery(gallery, iqb) {
					toDelete = append(toDelete, gallery.ID)
				} else {
					// increment progress, no further processing
					progress.Increment()
				}
			})
		}

		if len(galleries) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	if j.input.DryRun && len(toDelete) > 0 {
		// add progress for galleries that would've been deleted
		progress.AddProcessed(len(toDelete))
	}

	if !j.input.DryRun && len(toDelete) > 0 {
		progress.ExecuteTask(fmt.Sprintf("Cleaning %d galleries", len(toDelete)), func() {
			for _, galleryID := range toDelete {
				if job.IsCancelled(ctx) {
					return
				}

				j.deleteGallery(ctx, galleryID)

				progress.Increment()
			}
		})
	}

	return nil
}

func (j *cleanJob) processImages(ctx context.Context, progress *job.Progress, qb models.ImageReader) error {
	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)

	// performance consideration: order by path since default ordering by
	// title is slow
	sortBy := "path"
	findFilter.Sort = &sortBy

	var toDelete []int

	more := true
	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		images, err := image.Query(qb, nil, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for images: %w", err)
		}

		for _, image := range images {
			progress.ExecuteTask(fmt.Sprintf("Assessing image %s for clean", image.Path), func() {
				if j.shouldCleanImage(image) {
					toDelete = append(toDelete, image.ID)
				} else {
					// increment progress, no further processing
					progress.Increment()
				}
			})
		}

		if len(images) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}

	if j.input.DryRun && len(toDelete) > 0 {
		// add progress for images that would've been deleted
		progress.AddProcessed(len(toDelete))
	}

	if !j.input.DryRun && len(toDelete) > 0 {
		progress.ExecuteTask(fmt.Sprintf("Cleaning %d images", len(toDelete)), func() {
			for _, imageID := range toDelete {
				if job.IsCancelled(ctx) {
					return
				}

				j.deleteImage(ctx, imageID)

				progress.Increment()
			}
		})
	}

	return nil
}

func (j *cleanJob) shouldClean(path string) bool {
	// use image.FileExists for zip file checking
	fileExists := image.FileExists(path)

	// #1102 - clean anything in generated path
	generatedPath := config.GetInstance().GetGeneratedPath()
	if !fileExists || getStashFromPath(path) == nil || utils.IsPathInDir(generatedPath, path) {
		logger.Infof("File not found. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (j *cleanJob) shouldCleanScene(s *models.Scene) bool {
	if j.shouldClean(s.Path) {
		return true
	}

	stash := getStashFromPath(s.Path)
	if stash.ExcludeVideo {
		logger.Infof("File in stash library that excludes video. Marking to clean: \"%s\"", s.Path)
		return true
	}

	config := config.GetInstance()
	if !utils.MatchExtension(s.Path, config.GetVideoExtensions()) {
		logger.Infof("File extension does not match video extensions. Marking to clean: \"%s\"", s.Path)
		return true
	}

	if matchFile(s.Path, config.GetExcludes()) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", s.Path)
		return true
	}

	return false
}

func (j *cleanJob) shouldCleanGallery(g *models.Gallery, qb models.ImageReader) bool {
	// never clean manually created galleries
	if !g.Path.Valid {
		return false
	}

	path := g.Path.String
	if j.shouldClean(path) {
		return true
	}

	stash := getStashFromPath(path)
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Marking to clean: \"%s\"", path)
		return true
	}

	config := config.GetInstance()
	if g.Zip {
		if !utils.MatchExtension(path, config.GetGalleryExtensions()) {
			logger.Infof("File extension does not match gallery extensions. Marking to clean: \"%s\"", path)
			return true
		}

		if countImagesInZip(path) == 0 {
			logger.Infof("Gallery has 0 images. Marking to clean: \"%s\"", path)
			return true
		}
	} else {
		// folder-based - delete if it has no images
		count, err := qb.CountByGalleryID(g.ID)
		if err != nil {
			logger.Warnf("Error trying to count gallery images for %q: %v", path, err)
			return false
		}

		if count == 0 {
			return true
		}
	}

	if matchFile(path, config.GetImageExcludes()) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (j *cleanJob) shouldCleanImage(s *models.Image) bool {
	if j.shouldClean(s.Path) {
		return true
	}

	stash := getStashFromPath(s.Path)
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Marking to clean: \"%s\"", s.Path)
		return true
	}

	config := config.GetInstance()
	if !utils.MatchExtension(s.Path, config.GetImageExtensions()) {
		logger.Infof("File extension does not match image extensions. Marking to clean: \"%s\"", s.Path)
		return true
	}

	if matchFile(s.Path, config.GetImageExcludes()) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", s.Path)
		return true
	}

	return false
}

func (j *cleanJob) deleteScene(ctx context.Context, fileNamingAlgorithm models.HashAlgorithm, sceneID int) {
	fileNamingAlgo := GetInstance().Config.GetVideoFileNamingAlgorithm()

	fileDeleter := &scene.FileDeleter{
		Deleter:        *file.NewDeleter(),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          GetInstance().Paths,
	}
	var s *models.Scene
	if err := j.txnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		qb := repo.Scene()

		var err error
		s, err = qb.Find(sceneID)
		if err != nil {
			return err
		}

		if err := fileDeleter.MarkGeneratedFiles(s); err != nil {
			return err
		}

		return scene.Destroy(s, repo, fileDeleter)
	}); err != nil {
		errs := fileDeleter.Abort()
		for _, delErr := range errs {
			logger.Warn(delErr)
		}

		logger.Errorf("Error deleting scene from database: %s", err.Error())
		return
	}

	// perform the post-commit actions
	errs := fileDeleter.Complete()
	for _, delErr := range errs {
		logger.Warn(delErr)
	}

	GetInstance().PluginCache.ExecutePostHooks(ctx, sceneID, plugin.SceneDestroyPost, nil, nil)
}

func (j *cleanJob) deleteGallery(ctx context.Context, galleryID int) {
	if err := j.txnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		qb := repo.Gallery()
		return qb.Destroy(galleryID)
	}); err != nil {
		logger.Errorf("Error deleting gallery from database: %s", err.Error())
		return
	}

	GetInstance().PluginCache.ExecutePostHooks(ctx, galleryID, plugin.GalleryDestroyPost, nil, nil)
}

func (j *cleanJob) deleteImage(ctx context.Context, imageID int) {
	fileDeleter := &image.FileDeleter{
		Deleter: *file.NewDeleter(),
		Paths:   GetInstance().Paths,
	}

	if err := j.txnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		qb := repo.Image()

		i, err := qb.Find(imageID)
		if err != nil {
			return err
		}

		if i == nil {
			return fmt.Errorf("image not found: %d", imageID)
		}

		return image.Destroy(i, qb, fileDeleter, true, false)
	}); err != nil {
		errs := fileDeleter.Abort()
		for _, delErr := range errs {
			logger.Warn(delErr)
		}

		logger.Errorf("Error deleting image from database: %s", err.Error())
		return
	}

	// perform the post-commit actions
	errs := fileDeleter.Complete()
	for _, delErr := range errs {
		logger.Warn(delErr)
	}

	GetInstance().PluginCache.ExecutePostHooks(ctx, imageID, plugin.ImageDestroyPost, nil, nil)
}

func getStashFromPath(pathToCheck string) *models.StashConfig {
	for _, s := range config.GetInstance().GetStashPaths() {
		if utils.IsPathInDir(s.Path, filepath.Dir(pathToCheck)) {
			return s
		}
	}
	return nil
}

func getStashFromDirPath(pathToCheck string) *models.StashConfig {
	for _, s := range config.GetInstance().GetStashPaths() {
		if utils.IsPathInDir(s.Path, pathToCheck) {
			return s
		}
	}
	return nil
}
