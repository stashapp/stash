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
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scene"
)

type cleaner interface {
	Clean(ctx context.Context, options file.CleanOptions, progress *job.Progress)
}

type cleanJob struct {
	cleaner      cleaner
	txnManager   Repository
	input        CleanMetadataInput
	sceneService SceneService
	imageService ImageService
	scanSubs     *subscriptionManager
}

func (j *cleanJob) Execute(ctx context.Context, progress *job.Progress) {
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
		return
	}

	j.scanSubs.notify()
	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Finished Cleaning (%s)", elapsed))
}

type cleanFilter struct {
	scanFilter
}

func newCleanFilter(c *config.Instance) *cleanFilter {
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
		stash = getStashFromDirPath(f.stashPaths, path)
	} else {
		stash = getStashFromPath(f.stashPaths, path)
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
	case fsutil.MatchExtension(path, f.vidExt):
		return f.shouldCleanVideoFile(path, stash)
	case fsutil.MatchExtension(path, f.imgExt):
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

type cleanHandler struct {
	PluginCache *plugin.Cache
}

func (h *cleanHandler) HandleFile(ctx context.Context, fileDeleter *file.Deleter, fileID file.ID) error {
	if err := h.deleteRelatedScenes(ctx, fileDeleter, fileID); err != nil {
		return err
	}
	if err := h.deleteRelatedGalleries(ctx, fileID); err != nil {
		return err
	}
	if err := h.deleteRelatedImages(ctx, fileDeleter, fileID); err != nil {
		return err
	}

	return nil
}

func (h *cleanHandler) HandleFolder(ctx context.Context, fileDeleter *file.Deleter, folderID file.FolderID) error {
	return h.deleteRelatedFolderGalleries(ctx, folderID)
}

func (h *cleanHandler) deleteRelatedScenes(ctx context.Context, fileDeleter *file.Deleter, fileID file.ID) error {
	mgr := GetInstance()
	sceneQB := mgr.Database.Scene
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
		// only delete if the scene has no other files
		if len(scene.Files) <= 1 {
			logger.Infof("Deleting scene %q since it has no other related files", scene.GetTitle())
			if err := mgr.SceneService.Destroy(ctx, scene, sceneFileDeleter, true, false); err != nil {
				return err
			}

			checksum := scene.Checksum()
			oshash := scene.OSHash()

			mgr.PluginCache.RegisterPostHooks(ctx, mgr.Database, scene.ID, plugin.SceneDestroyPost, plugin.SceneDestroyInput{
				Checksum: checksum,
				OSHash:   oshash,
				Path:     scene.Path(),
			}, nil)
		}
	}

	return nil
}

func (h *cleanHandler) deleteRelatedGalleries(ctx context.Context, fileID file.ID) error {
	mgr := GetInstance()
	qb := mgr.Database.Gallery
	galleries, err := qb.FindByFileID(ctx, fileID)
	if err != nil {
		return err
	}

	for _, g := range galleries {
		// only delete if the gallery has no other files
		if len(g.Files) <= 1 {
			logger.Infof("Deleting gallery %q since it has no other related files", g.GetTitle())
			if err := qb.Destroy(ctx, g.ID); err != nil {
				return err
			}

			mgr.PluginCache.RegisterPostHooks(ctx, mgr.Database, g.ID, plugin.GalleryDestroyPost, plugin.GalleryDestroyInput{
				Checksum: g.Checksum(),
				Path:     g.Path(),
			}, nil)
		}
	}

	return nil
}

func (h *cleanHandler) deleteRelatedFolderGalleries(ctx context.Context, folderID file.FolderID) error {
	mgr := GetInstance()
	qb := mgr.Database.Gallery
	galleries, err := qb.FindByFolderID(ctx, folderID)
	if err != nil {
		return err
	}

	for _, g := range galleries {
		logger.Infof("Deleting folder-based gallery %q since the folder no longer exists", g.GetTitle())
		if err := qb.Destroy(ctx, g.ID); err != nil {
			return err
		}

		mgr.PluginCache.RegisterPostHooks(ctx, mgr.Database, g.ID, plugin.GalleryDestroyPost, plugin.GalleryDestroyInput{
			Checksum: g.Checksum(),
			Path:     g.Path(),
		}, nil)
	}

	return nil
}

func (h *cleanHandler) deleteRelatedImages(ctx context.Context, fileDeleter *file.Deleter, fileID file.ID) error {
	mgr := GetInstance()
	imageQB := mgr.Database.Image
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
			if err := mgr.ImageService.Destroy(ctx, i, imageFileDeleter, true, false); err != nil {
				return err
			}

			mgr.PluginCache.RegisterPostHooks(ctx, mgr.Database, i.ID, plugin.ImageDestroyPost, plugin.ImageDestroyInput{
				Checksum: i.Checksum(),
				Path:     i.Path(),
			}, nil)
		}
	}

	return nil
}

func getStashFromPath(stashes []*config.StashConfig, pathToCheck string) *config.StashConfig {
	for _, f := range stashes {
		if fsutil.IsPathInDir(f.Path, filepath.Dir(pathToCheck)) {
			return f
		}
	}
	return nil
}

func getStashFromDirPath(stashes []*config.StashConfig, pathToCheck string) *config.StashConfig {
	for _, f := range stashes {
		if fsutil.IsPathInDir(f.Path, pathToCheck) {
			return f
		}
	}
	return nil
}
