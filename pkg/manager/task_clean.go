package manager

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

type CleanTask struct {
	ctx                 context.Context
	TxnManager          models.TransactionManager
	Scene               *models.Scene
	Gallery             *models.Gallery
	Image               *models.Image
	fileNamingAlgorithm models.HashAlgorithm
}

func (t *CleanTask) Start(wg *sync.WaitGroup, dryRun bool) {
	defer wg.Done()

	if t.Scene != nil && t.shouldCleanScene(t.Scene) && !dryRun {
		t.deleteScene(t.Scene.ID)
	}

	if t.Gallery != nil && t.shouldCleanGallery(t.Gallery) && !dryRun {
		t.deleteGallery(t.Gallery.ID)
	}

	if t.Image != nil && t.shouldCleanImage(t.Image) && !dryRun {
		t.deleteImage(t.Image.ID)
	}
}

func (t *CleanTask) shouldClean(path string) bool {
	// use image.FileExists for zip file checking
	fileExists := image.FileExists(path)

	// #1102 - clean anything in generated path
	generatedPath := config.GetInstance().GetGeneratedPath()
	if !fileExists || getStashFromPath(path) == nil || utils.IsPathInDir(generatedPath, path) {
		logger.Infof("File not found. Cleaning: \"%s\"", path)
		return true
	}

	return false
}

func (t *CleanTask) shouldCleanScene(s *models.Scene) bool {
	if t.shouldClean(s.Path) {
		return true
	}

	stash := getStashFromPath(s.Path)
	if stash.ExcludeVideo {
		logger.Infof("File in stash library that excludes video. Cleaning: \"%s\"", s.Path)
		return true
	}

	config := config.GetInstance()
	if !utils.MatchExtension(s.Path, config.GetVideoExtensions()) {
		logger.Infof("File extension does not match video extensions. Cleaning: \"%s\"", s.Path)
		return true
	}

	if matchFile(s.Path, config.GetExcludes()) {
		logger.Infof("File matched regex. Cleaning: \"%s\"", s.Path)
		return true
	}

	return false
}

func (t *CleanTask) shouldCleanGallery(g *models.Gallery) bool {
	// never clean manually created galleries
	if !g.Zip {
		return false
	}

	path := g.Path.String
	if t.shouldClean(path) {
		return true
	}

	stash := getStashFromPath(path)
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Cleaning: \"%s\"", path)
		return true
	}

	config := config.GetInstance()
	if !utils.MatchExtension(path, config.GetGalleryExtensions()) {
		logger.Infof("File extension does not match gallery extensions. Cleaning: \"%s\"", path)
		return true
	}

	if matchFile(path, config.GetImageExcludes()) {
		logger.Infof("File matched regex. Cleaning: \"%s\"", path)
		return true
	}

	if countImagesInZip(path) == 0 {
		logger.Infof("Gallery has 0 images. Cleaning: \"%s\"", path)
		return true
	}

	return false
}

func (t *CleanTask) shouldCleanImage(s *models.Image) bool {
	if t.shouldClean(s.Path) {
		return true
	}

	stash := getStashFromPath(s.Path)
	if stash.ExcludeImage {
		logger.Infof("File in stash library that excludes images. Cleaning: \"%s\"", s.Path)
		return true
	}

	config := config.GetInstance()
	if !utils.MatchExtension(s.Path, config.GetImageExtensions()) {
		logger.Infof("File extension does not match image extensions. Cleaning: \"%s\"", s.Path)
		return true
	}

	if matchFile(s.Path, config.GetImageExcludes()) {
		logger.Infof("File matched regex. Cleaning: \"%s\"", s.Path)
		return true
	}

	return false
}

func (t *CleanTask) deleteScene(sceneID int) {
	var postCommitFunc func()
	var scene *models.Scene
	if err := t.TxnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		qb := repo.Scene()

		var err error
		scene, err = qb.Find(sceneID)
		if err != nil {
			return err
		}
		postCommitFunc, err = DestroyScene(scene, repo)
		return err
	}); err != nil {
		logger.Errorf("Error deleting scene from database: %s", err.Error())
		return
	}

	postCommitFunc()

	DeleteGeneratedSceneFiles(scene, t.fileNamingAlgorithm)

	GetInstance().PluginCache.ExecutePostHooks(t.ctx, sceneID, plugin.SceneDestroyPost, nil, nil)
}

func (t *CleanTask) deleteGallery(galleryID int) {
	if err := t.TxnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		qb := repo.Gallery()
		return qb.Destroy(galleryID)
	}); err != nil {
		logger.Errorf("Error deleting gallery from database: %s", err.Error())
		return
	}

	GetInstance().PluginCache.ExecutePostHooks(t.ctx, galleryID, plugin.GalleryDestroyPost, nil, nil)
}

func (t *CleanTask) deleteImage(imageID int) {

	if err := t.TxnManager.WithTxn(context.TODO(), func(repo models.Repository) error {
		qb := repo.Image()

		return qb.Destroy(imageID)
	}); err != nil {
		logger.Errorf("Error deleting image from database: %s", err.Error())
		return
	}

	pathErr := os.Remove(GetInstance().Paths.Generated.GetThumbnailPath(t.Image.Checksum, models.DefaultGthumbWidth)) // remove cache dir of gallery
	if pathErr != nil {
		logger.Errorf("Error deleting thumbnail image from cache: %s", pathErr)
	}

	GetInstance().PluginCache.ExecutePostHooks(t.ctx, imageID, plugin.ImageDestroyPost, nil, nil)
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
