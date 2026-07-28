package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/scene"
)

type purgeMissingJob struct {
	purger *file.MissingPurger

	options      file.PurgeMissingOptions
	repository   models.Repository
	sceneService SceneService
	imageService ImageService
	scanSubs     *subscriptionManager
}

// galleryCleanCandidate avoids retaining full gallery models for large libraries.
type galleryCleanCandidate struct {
	id   int
	path string
}

func (j *purgeMissingJob) Execute(ctx context.Context, progress *job.Progress) error {
	logger.Infof("Starting purging of missing files and folders")
	start := time.Now()
	if j.options.DryRun {
		logger.Infof("Running in Dry Mode")
	}

	j.purger.PurgeMissing(ctx, j.options, progress)

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return nil
	}

	// only clean galleries if not in dry run mode
	if !j.options.DryRun {
		j.cleanGalleries(ctx)
	}

	j.scanSubs.notify()
	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Finished purging missing files and folders (%s)", elapsed))
	return nil
}

func (j *purgeMissingJob) cleanGalleries(ctx context.Context) {
	toClean, err := j.findGalleriesToClean(ctx)
	if err != nil {
		logger.Errorf("Error finding galleries to clean: %v", err)
		return
	}

	if !j.options.DryRun {
		for _, id := range toClean {
			j.deleteGallery(ctx, id)
		}
	}
}

func (j *purgeMissingJob) findGalleriesToClean(ctx context.Context) ([]int, error) {
	const batchSize = 1000
	var toClean []int
	var folderGalleries []galleryCleanCandidate
	toCleanSet := make(map[int]struct{})
	r := j.repository

	inCleanPaths := func(g *models.Gallery) bool {
		return g.Path != "" && (len(j.options.Paths) == 0 || fsutil.IsPathInDirs(j.options.Paths, g.Path))
	}

	markForClean := func(id int, displayName, reason string) {
		if _, exists := toCleanSet[id]; exists {
			return
		}

		logger.Infof("%s. Marking to clean: %s", reason, displayName)
		toCleanSet[id] = struct{}{}
		toClean = append(toClean, id)
	}

	queryInBatches := func(ctx context.Context, galleryFilter *models.GalleryFilterType, visit func(*models.Gallery)) error {
		findFilter := models.BatchFindFilter(batchSize)
		found := true
		for found {
			galleries, _, err := r.Gallery.Query(ctx, galleryFilter, findFilter)
			if err != nil {
				return err
			}

			found = len(galleries) > 0
			for _, g := range galleries {
				visit(g)
			}

			*findFilter.Page++
		}

		return nil
	}

	err := r.WithTxn(ctx, func(ctx context.Context) error {
		emptyGalleryFilter := &models.GalleryFilterType{
			ImageCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		if err := queryInBatches(ctx, emptyGalleryFilter, func(g *models.Gallery) {
			if inCleanPaths(g) {
				markForClean(g.ID, g.DisplayName(), "Gallery has 0 images")
			}
		}); err != nil {
			return err
		}

		folderGalleryFilter := &models.GalleryFilterType{
			FoldersFilter: &models.FolderFilterType{
				Path: &models.StringCriterionInput{
					Modifier: models.CriterionModifierNotNull,
				},
			},
		}
		return queryInBatches(ctx, folderGalleryFilter, func(g *models.Gallery) {
			if inCleanPaths(g) {
				folderGalleries = append(folderGalleries, galleryCleanCandidate{
					id:   g.ID,
					path: g.Path,
				})
			}
		})
	})
	if err != nil {
		return nil, err
	}

	// Gallery folders may be on slow network filesystems, so check them after
	// releasing the database transaction.
	for _, g := range folderGalleries {
		force, exempt, err := image.GalleryFolderMarkers(g.path)
		if err != nil {
			logger.Errorf("Error checking gallery marker files in %q: %v", g.path, err)
			continue
		}

		// .forcegallery wins, matching Scan
		if exempt && !force {
			markForClean(g.id, g.path, "Gallery folder contains .nogallery")
		}
	}

	return toClean, nil
}

func (j *purgeMissingJob) deleteGallery(ctx context.Context, id int) {
	pluginCache := GetInstance().PluginCache

	r := j.repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		qb := r.Gallery
		g, err := qb.Find(ctx, id)
		if err != nil {
			return err
		}

		if g == nil {
			return fmt.Errorf("gallery with id %d not found", id)
		}

		if err := g.LoadPrimaryFile(ctx, r.File); err != nil {
			return err
		}

		if err := qb.Destroy(ctx, id); err != nil {
			return err
		}

		pluginCache.RegisterPostHooks(ctx, id, hook.GalleryDestroyPost, plugin.GalleryDestroyInput{
			Checksum: g.PrimaryChecksum(),
			Path:     g.Path,
		}, nil)

		return nil
	}); err != nil {
		logger.Errorf("Error deleting gallery from database: %s", err.Error())
	}
}

type purgeHandler struct{}

func (h *purgeHandler) HandleFile(ctx context.Context, fileDeleter *file.Deleter, fileID models.FileID) error {
	if err := h.handleRelatedScenes(ctx, fileDeleter, fileID); err != nil {
		return err
	}
	if err := h.handleRelatedGalleries(ctx, fileID); err != nil {
		return err
	}
	if err := h.handleRelatedImages(ctx, fileDeleter, fileID); err != nil {
		return err
	}

	return nil
}

func (h *purgeHandler) HandleFolder(ctx context.Context, fileDeleter *file.Deleter, folderID models.FolderID) error {
	return h.deleteRelatedFolderGalleries(ctx, folderID)
}

func (h *purgeHandler) handleRelatedScenes(ctx context.Context, fileDeleter *file.Deleter, fileID models.FileID) error {
	mgr := GetInstance()
	sceneQB := mgr.Repository.Scene
	scenes, err := sceneQB.FindByFileID(ctx, fileID)
	if err != nil {
		return err
	}

	fileNamingAlgo := mgr.Config.GetVideoFileNamingAlgorithm()

	sceneFileDeleter := &scene.FileDeleter{
		Deleter:        fileDeleter,
		FileNamingAlgo: fileNamingAlgo,
		Paths:          mgr.Paths,
	}

	for _, scene := range scenes {
		if err := scene.LoadFiles(ctx, sceneQB); err != nil {
			return err
		}

		// only delete if the scene has no other files
		if len(scene.Files.List()) <= 1 {
			logger.Infof("Deleting scene %q since it has no other related files", scene.DisplayName())
			const deleteGenerated = true
			const deleteFile = false
			const destroyFileEntry = false
			if err := mgr.SceneService.Destroy(ctx, scene, sceneFileDeleter, deleteGenerated, deleteFile, destroyFileEntry); err != nil {
				return err
			}

			mgr.PluginCache.RegisterPostHooks(ctx, scene.ID, hook.SceneDestroyPost, plugin.SceneDestroyInput{
				Checksum: scene.Checksum,
				OSHash:   scene.OSHash,
				Path:     scene.Path,
			}, nil)
		} else {
			// set the primary file to a remaining file
			var newPrimaryID models.FileID
			for _, f := range scene.Files.List() {
				if f.ID != fileID {
					newPrimaryID = f.ID
					break
				}
			}

			scenePartial := models.NewScenePartial()
			scenePartial.PrimaryFileID = &newPrimaryID

			if _, err := mgr.Repository.Scene.UpdatePartial(ctx, scene.ID, scenePartial); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *purgeHandler) handleRelatedGalleries(ctx context.Context, fileID models.FileID) error {
	mgr := GetInstance()
	qb := mgr.Repository.Gallery
	galleries, err := qb.FindByFileID(ctx, fileID)
	if err != nil {
		return err
	}

	for _, g := range galleries {
		if err := g.LoadFiles(ctx, qb); err != nil {
			return err
		}

		// only delete if the gallery has no other files
		if len(g.Files.List()) <= 1 {
			logger.Infof("Deleting gallery %q since it has no other related files", g.DisplayName())
			if err := qb.Destroy(ctx, g.ID); err != nil {
				return err
			}

			mgr.PluginCache.RegisterPostHooks(ctx, g.ID, hook.GalleryDestroyPost, plugin.GalleryDestroyInput{
				Checksum: g.PrimaryChecksum(),
				Path:     g.Path,
			}, nil)
		} else {
			// set the primary file to a remaining file
			var newPrimaryID models.FileID
			for _, f := range g.Files.List() {
				if f.Base().ID != fileID {
					newPrimaryID = f.Base().ID
					break
				}
			}

			galleryPartial := models.NewGalleryPartial()
			galleryPartial.PrimaryFileID = &newPrimaryID

			if _, err := mgr.Repository.Gallery.UpdatePartial(ctx, g.ID, galleryPartial); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *purgeHandler) deleteRelatedFolderGalleries(ctx context.Context, folderID models.FolderID) error {
	mgr := GetInstance()
	qb := mgr.Repository.Gallery
	galleries, err := qb.FindByFolderID(ctx, folderID)
	if err != nil {
		return err
	}

	for _, g := range galleries {
		logger.Infof("Deleting folder-based gallery %q since the folder no longer exists", g.DisplayName())
		if err := qb.Destroy(ctx, g.ID); err != nil {
			return err
		}

		mgr.PluginCache.RegisterPostHooks(ctx, g.ID, hook.GalleryDestroyPost, plugin.GalleryDestroyInput{
			// No checksum for folders
			// Checksum: g.Checksum(),
			Path: g.Path,
		}, nil)
	}

	return nil
}

func (h *purgeHandler) handleRelatedImages(ctx context.Context, fileDeleter *file.Deleter, fileID models.FileID) error {
	mgr := GetInstance()
	imageQB := mgr.Repository.Image
	images, err := imageQB.FindByFileID(ctx, fileID)
	if err != nil {
		return err
	}

	imageFileDeleter := &image.FileDeleter{
		Deleter: fileDeleter,
		Paths:   mgr.Paths,
	}

	for _, i := range images {
		if err := i.LoadFiles(ctx, imageQB); err != nil {
			return err
		}

		if len(i.Files.List()) <= 1 {
			logger.Infof("Deleting image %q since it has no other related files", i.DisplayName())
			const deleteGenerated = true
			const deleteFile = false
			const destroyFileEntry = false
			if err := mgr.ImageService.Destroy(ctx, i, imageFileDeleter, deleteGenerated, deleteFile, destroyFileEntry); err != nil {
				return err
			}

			mgr.PluginCache.RegisterPostHooks(ctx, i.ID, hook.ImageDestroyPost, plugin.ImageDestroyInput{
				Checksum: i.Checksum,
				Path:     i.Path,
			}, nil)
		} else {
			// set the primary file to a remaining file
			var newPrimaryID models.FileID
			for _, f := range i.Files.List() {
				if f.Base().ID != fileID {
					newPrimaryID = f.Base().ID
					break
				}
			}

			imagePartial := models.NewImagePartial()
			imagePartial.PrimaryFileID = &newPrimaryID

			if _, err := mgr.Repository.Image.UpdatePartial(ctx, i.ID, imagePartial); err != nil {
				return err
			}
		}
	}

	return nil
}
