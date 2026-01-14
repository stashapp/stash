package manager

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
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

type cleaner interface {
	Clean(ctx context.Context, options file.CleanOptions, progress *job.Progress)
}

type cleanJob struct {
	cleaner      cleaner
	repository   models.Repository
	input        CleanMetadataInput
	sceneService SceneService
	imageService ImageService
	scanSubs     *subscriptionManager
}

func (j *cleanJob) Execute(ctx context.Context, progress *job.Progress) error {
	logger.Infof("Starting cleaning of tracked files")
	start := time.Now()
	if j.input.DryRun {
		logger.Infof("Running in Dry Mode")
	}

	j.cleaner.Clean(ctx, file.CleanOptions{
		Paths:      j.input.Paths,
		DryRun:     j.input.DryRun,
		PathFilter: newCleanFilter(instance.Config),
	}, progress)

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return nil
	}

	j.cleanEmptyGalleries(ctx)

	j.scanSubs.notify()
	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Finished Cleaning (%s)", elapsed))
	return nil
}

func (j *cleanJob) cleanEmptyGalleries(ctx context.Context) {
	const batchSize = 1000
	var toClean []int
	findFilter := models.BatchFindFilter(batchSize)
	r := j.repository
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		found := true
		for found {
			emptyGalleries, _, err := r.Gallery.Query(ctx, &models.GalleryFilterType{
				ImageCount: &models.IntCriterionInput{
					Value:    0,
					Modifier: models.CriterionModifierEquals,
				},
			}, findFilter)

			if err != nil {
				return err
			}

			found = len(emptyGalleries) > 0

			for _, g := range emptyGalleries {
				if g.Path == "" {
					continue
				}

				if len(j.input.Paths) > 0 && !fsutil.IsPathInDirs(j.input.Paths, g.Path) {
					continue
				}

				logger.Infof("Gallery has 0 images. Marking to clean: %s", g.DisplayName())
				toClean = append(toClean, g.ID)
			}

			*findFilter.Page++
		}

		return nil
	}); err != nil {
		logger.Errorf("Error finding empty galleries: %v", err)
		return
	}

	if !j.input.DryRun {
		for _, id := range toClean {
			j.deleteGallery(ctx, id)
		}
	}
}

func (j *cleanJob) deleteGallery(ctx context.Context, id int) {
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

type cleanFilter struct {
	scanFilter
}

func newCleanFilter(c *config.Config) *cleanFilter {
	return &cleanFilter{
		scanFilter: scanFilter{
			extensionConfig:   newExtensionConfig(c),
			stashPaths:        c.GetStashPaths(),
			generatedPath:     c.GetGeneratedPath(),
			videoExcludeRegex: generateRegexps(c.GetExcludes()),
			imageExcludeRegex: generateRegexps(c.GetImageExcludes()),
		},
	}
}

func (f *cleanFilter) Accept(ctx context.Context, path string, info fs.FileInfo) bool {
	//  #1102 - clean anything in generated path
	generatedPath := f.generatedPath

	var stash *config.StashConfig
	fileOrFolder := "File"

	if info.IsDir() {
		fileOrFolder = "Folder"
		stash = f.stashPaths.GetStashFromDirPath(path)
	} else {
		stash = f.stashPaths.GetStashFromPath(path)
	}

	if stash == nil {
		logger.Infof("%s not in any stash library directories. Marking to clean: \"%s\"", fileOrFolder, path)
		return false
	}

	if fsutil.IsPathInDir(generatedPath, path) {
		logger.Infof("%s is in generated path. Marking to clean: \"%s\"", fileOrFolder, path)
		return false
	}

	if info.IsDir() {
		return !f.shouldCleanFolder(path, stash)
	}

	return !f.shouldCleanFile(path, info, stash)
}

func (f *cleanFilter) shouldCleanFolder(path string, s *config.StashConfig) bool {
	// only delete folders where it is excluded from everything
	pathExcludeTest := path + string(filepath.Separator)
	if (s.ExcludeVideo || matchFileRegex(pathExcludeTest, f.videoExcludeRegex)) && (s.ExcludeImage || matchFileRegex(pathExcludeTest, f.imageExcludeRegex)) {
		logger.Infof("Folder is excluded from both video and image. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (f *cleanFilter) shouldCleanFile(path string, info fs.FileInfo, stash *config.StashConfig) bool {
	switch {
	case info.IsDir() || fsutil.MatchExtension(path, f.zipExt):
		return f.shouldCleanGallery(path, stash)
	case useAsVideo(path):
		return f.shouldCleanVideoFile(path, stash)
	case useAsImage(path):
		return f.shouldCleanImage(path, stash)
	default:
		logger.Infof("File extension does not match any media extensions. Marking to clean: \"%s\"", path)
		return true
	}
}

func (f *cleanFilter) shouldCleanVideoFile(path string, stash *config.StashConfig) bool {
	if stash.ExcludeVideo {
		logger.Infof("File in stash library that excludes video. Marking to clean: \"%s\"", path)
		return true
	}

	if matchFileRegex(path, f.videoExcludeRegex) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (f *cleanFilter) shouldCleanGallery(path string, stash *config.StashConfig) bool {
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Marking to clean: \"%s\"", path)
		return true
	}

	if matchFileRegex(path, f.imageExcludeRegex) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

func (f *cleanFilter) shouldCleanImage(path string, stash *config.StashConfig) bool {
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Marking to clean: \"%s\"", path)
		return true
	}

	if matchFileRegex(path, f.imageExcludeRegex) {
		logger.Infof("File matched regex. Marking to clean: \"%s\"", path)
		return true
	}

	return false
}

type cleanHandler struct{}

func (h *cleanHandler) HandleFile(ctx context.Context, fileDeleter *file.Deleter, fileID models.FileID) error {
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

func (h *cleanHandler) HandleFolder(ctx context.Context, fileDeleter *file.Deleter, folderID models.FolderID) error {
	return h.deleteRelatedFolderGalleries(ctx, folderID)
}

func (h *cleanHandler) handleRelatedScenes(ctx context.Context, fileDeleter *file.Deleter, fileID models.FileID) error {
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

func (h *cleanHandler) handleRelatedGalleries(ctx context.Context, fileID models.FileID) error {
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

func (h *cleanHandler) deleteRelatedFolderGalleries(ctx context.Context, folderID models.FolderID) error {
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

func (h *cleanHandler) handleRelatedImages(ctx context.Context, fileDeleter *file.Deleter, fileID models.FileID) error {
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
