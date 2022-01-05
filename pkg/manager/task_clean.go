package manager

import (
	"context"
	"fmt"

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

		if err := j.processScenes(ctx, progress, r.Scene(), r.File()); err != nil {
			return fmt.Errorf("error cleaning scenes: %w", err)
		}
		// doing galleries first since it might result in image files being deleted
		if err := j.processGalleries(ctx, progress, r.Gallery(), r.Image(), r.File()); err != nil {
			return fmt.Errorf("error cleaning galleries: %w", err)
		}
		if err := j.processImages(ctx, progress, r.Image(), r.File()); err != nil {
			return fmt.Errorf("error cleaning images: %w", err)
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
	sceneFilter := scene.PathsFilter(j.input.Paths)
	sceneResult, err := r.Scene().Query(models.SceneQueryOptions{
		QueryOptions: models.QueryOptions{
			Count: true,
		},
		SceneFilter: sceneFilter,
	})
	if err != nil {
		return 0, err
	}

	imageCount, err := r.Image().QueryCount(image.PathsFilter(j.input.Paths), nil)
	if err != nil {
		return 0, err
	}

	galleryCount, err := r.Gallery().QueryCount(gallery.PathsFilter(j.input.Paths), nil)
	if err != nil {
		return 0, err
	}

	return sceneResult.Count + imageCount + galleryCount, nil
}

type deleteSet struct {
	// 0 means to not delete the object
	objectID int
	fileIDs  []int
}

func (s deleteSet) deleteFiles(w models.FileWriter) error {
	// delete files
	for _, fileID := range s.fileIDs {
		if err := w.Destroy(fileID); err != nil {
			return fmt.Errorf("deleting file with id %d: %w", fileID, err)
		}
	}

	return nil
}

func (j *cleanJob) processScenes(ctx context.Context, progress *job.Progress, qb models.SceneReader, fileReader models.FileReader) error {
	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	sceneFilter := scene.PathsFilter(j.input.Paths)
	sort := "path"
	findFilter.Sort = &sort

	var toDelete []deleteSet

	cleaner := file.Cleaner{
		FileReader: fileReader,
		// DryRun:           j.input.DryRun,
		Config: config.GetInstance(),
	}

	more := true
	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		scenes, err := scene.Query(qb, sceneFilter, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for scenes: %w", err)
		}

		for _, scene := range scenes {
			progress.ExecuteTask(fmt.Sprintf("Assessing scene %s for clean", scene.Path), func() {
				var fileIDs []int
				fileIDs, err = qb.GetFileIDs(scene.ID)
				if err != nil {
					err = fmt.Errorf("error getting file ids for scene %q: %w", scene.Path, err)
					return
				}

				var toClean []int
				toClean, err = cleaner.GetFilesToClean(fileIDs, scene)
				if err != nil {
					err = fmt.Errorf("error cleaning files for scene %q: %w", scene.Path, err)
					return
				}

				if len(toClean) > 0 || len(fileIDs) == 0 {
					set := deleteSet{
						fileIDs: toClean,
					}

					if len(toClean) == len(fileIDs) {
						// all files for this object are removed, remove this object as well
						logger.Infof("All scene files removed. Marking to clean: %s", scene.Path)
						set.objectID = scene.ID
					}

					toDelete = append(toDelete, set)
				} else {
					// increment progress, no further processing
					progress.Increment()
				}
			})

			if err != nil {
				return err
			}
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
			for _, set := range toDelete {
				if job.IsCancelled(ctx) {
					return
				}

				j.deleteScene(ctx, fileNamingAlgorithm, set)

				progress.Increment()
			}
		})
	}

	return nil
}

func (j *cleanJob) processGalleries(ctx context.Context, progress *job.Progress, qb models.GalleryReader, iqb models.ImageReader, fileReader models.FileReader) error {
	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	galleryFilter := gallery.PathsFilter(j.input.Paths)

	if galleryFilter == nil {
		galleryFilter = &models.GalleryFilterType{
			Path: &models.StringCriterionInput{
				Modifier: models.CriterionModifierNotNull,
			},
		}
	}

	sort := "path"
	findFilter.Sort = &sort

	var toDelete []deleteSet

	cleaner := file.Cleaner{
		FileReader: fileReader,
		Config:     config.GetInstance(),
	}

	more := true
	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		galleries, _, err := qb.Query(galleryFilter, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for galleries: %w", err)
		}

		for _, gallery := range galleries {
			progress.ExecuteTask(fmt.Sprintf("Assessing gallery %s for clean", gallery.GetTitle()), func() {
				var fileIDs []int
				fileIDs, err = qb.GetFileIDs(gallery.ID)
				if err != nil {
					err = fmt.Errorf("error getting file ids for gallery %q: %w", gallery.Path.String, err)
					return
				}

				if gallery.Zip {
					var toClean []int
					toClean, err = cleaner.GetFilesToClean(fileIDs, gallery)
					if err != nil {
						err = fmt.Errorf("error cleaning files for gallery %q: %w", gallery.Path.String, err)
						return
					}

					if len(toClean) > 0 || len(fileIDs) == 0 {
						set := deleteSet{
							fileIDs: toClean,
						}

						if len(toClean) == len(fileIDs) {
							// all files for this object are removed, remove this object as well
							logger.Infof("All gallery files removed. Marking to clean: %s", gallery.Path.String)
							set.objectID = gallery.ID
						}

						toDelete = append(toDelete, set)
					} else {
						// increment progress, no further processing
						progress.Increment()
					}
				} else {
					// folder-based - delete if it has no images
					var count int
					count, err = iqb.CountByGalleryID(gallery.ID)
					if err != nil {
						err = fmt.Errorf("counting gallery images for %q: %v", gallery.Path.String, err)
						return
					}

					if count == 0 {
						toDelete = append(toDelete, deleteSet{
							objectID: gallery.ID,
						})
					}
				}
			})

			if err != nil {
				return err
			}
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
			for _, set := range toDelete {
				if job.IsCancelled(ctx) {
					return
				}

				j.deleteGallery(ctx, set)

				progress.Increment()
			}
		})
	}

	return nil
}

func (j *cleanJob) processImages(ctx context.Context, progress *job.Progress, qb models.ImageReader, fileReader models.FileReader) error {
	batchSize := 1000

	findFilter := models.BatchFindFilter(batchSize)
	imageFilter := image.PathsFilter(j.input.Paths)

	// performance consideration: order by path since default ordering by
	// title is slow
	sortBy := "path"
	findFilter.Sort = &sortBy

	var toDelete []deleteSet
	cleaner := file.Cleaner{
		FileReader: fileReader,
		Config:     config.GetInstance(),
	}

	more := true
	for more {
		if job.IsCancelled(ctx) {
			return nil
		}

		images, err := image.Query(qb, imageFilter, findFilter)
		if err != nil {
			return fmt.Errorf("error querying for images: %w", err)
		}

		for _, image := range images {
			progress.ExecuteTask(fmt.Sprintf("Assessing image %s for clean", image.Path), func() {
				var fileIDs []int
				fileIDs, err = qb.GetFileIDs(image.ID)
				if err != nil {
					err = fmt.Errorf("error getting file ids for image %q: %w", image.Path, err)
					return
				}

				var toClean []int
				toClean, err = cleaner.GetFilesToClean(fileIDs, image)
				if err != nil {
					err = fmt.Errorf("error cleaning files for image %q: %w", image.Path, err)
					return
				}

				if len(toClean) > 0 || len(fileIDs) == 0 {
					set := deleteSet{
						fileIDs: toClean,
					}

					if len(toClean) == len(fileIDs) {
						// all files for this object are removed, remove this object as well
						logger.Infof("All image files removed. Marking to clean: %s", image.Path)
						set.objectID = image.ID
					}

					toDelete = append(toDelete, set)
				} else {
					// increment progress, no further processing
					progress.Increment()
				}
			})

			if err != nil {
				return err
			}
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

func (j *cleanJob) deleteScene(ctx context.Context, fileNamingAlgorithm models.HashAlgorithm, toDelete deleteSet) {
	fileNamingAlgo := GetInstance().Config.GetVideoFileNamingAlgorithm()

	fileDeleter := &scene.FileDeleter{
		Deleter:        *file.NewDeleter(),
		FileNamingAlgo: fileNamingAlgo,
		Paths:          GetInstance().Paths,
	}

	var s *models.Scene
	if err := j.txnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		var err error

		if err := toDelete.deleteFiles(repo.File()); err != nil {
			return err
		}

		if toDelete.objectID != 0 {
			qb := repo.Scene()

			s, err = qb.Find(toDelete.objectID)
			if err != nil {
				return err
			}

			return scene.Destroy(s, repo, fileDeleter, true, false)
		}

		return nil
	}); err != nil {
		fileDeleter.Rollback()

		logger.Errorf("Error deleting scene from database: %s", err.Error())
		return
	}

	if toDelete.objectID != 0 {
		// perform the post-commit actions
		fileDeleter.Commit()

		GetInstance().PluginCache.ExecutePostHooks(ctx, toDelete.objectID, plugin.SceneDestroyPost, plugin.SceneDestroyInput{
			Checksum: s.Checksum.String,
			OSHash:   s.OSHash.String,
			Path:     s.Path,
		}, nil)
	}
}

func (j *cleanJob) deleteGallery(ctx context.Context, toDelete deleteSet) {
	var g *models.Gallery

	if err := j.txnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		if err := toDelete.deleteFiles(repo.File()); err != nil {
			return err
		}

		if toDelete.objectID != 0 {
			qb := repo.Gallery()

			var err error
			g, err = qb.Find(toDelete.objectID)
			if err != nil {
				return err
			}

			return qb.Destroy(toDelete.objectID)
		}

		return nil
	}); err != nil {
		logger.Errorf("Error deleting gallery from database: %s", err.Error())
		return
	}

	if toDelete.objectID != 0 {
		GetInstance().PluginCache.ExecutePostHooks(ctx, toDelete.objectID, plugin.GalleryDestroyPost, plugin.GalleryDestroyInput{
			Checksum: g.Checksum,
			Path:     g.Path.String,
		}, nil)
	}
}

func (j *cleanJob) deleteImage(ctx context.Context, toDelete deleteSet) {
	fileDeleter := &image.FileDeleter{
		Deleter: *file.NewDeleter(),
		Paths:   GetInstance().Paths,
	}

	var i *models.Image
	if err := j.txnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		if err := toDelete.deleteFiles(repo.File()); err != nil {
			return err
		}

		if toDelete.objectID != 0 {
			qb := repo.Image()

			var err error
			i, err = qb.Find(toDelete.objectID)
			if err != nil {
				return err
			}

			if i == nil {
				return fmt.Errorf("image not found: %d", toDelete.objectID)
			}

			return image.Destroy(i, repo, fileDeleter, true, false)
		}

		return nil
	}); err != nil {
		fileDeleter.Rollback()

		logger.Errorf("Error deleting image from database: %s", err.Error())
		return
	}

	if toDelete.objectID != 0 {
		// perform the post-commit actions
		fileDeleter.Commit()
		GetInstance().PluginCache.ExecutePostHooks(ctx, toDelete.objectID, plugin.ImageDestroyPost, plugin.ImageDestroyInput{
			Checksum: i.Checksum,
			Path:     i.Path,
		}, nil)
	}
}

func getStashFromDirPath(pathToCheck string) *models.StashConfig {
	for _, s := range config.GetInstance().GetStashPaths() {
		if utils.IsPathInDir(s.Path, pathToCheck) {
			return s
		}
	}
	return nil
}
