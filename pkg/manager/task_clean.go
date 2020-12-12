package manager

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

type CleanTask struct {
	Scene               *models.Scene
	Gallery             *models.Gallery
	Image               *models.Image
	fileNamingAlgorithm models.HashAlgorithm
}

func (t *CleanTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.Scene != nil && t.shouldCleanScene(t.Scene) {
		t.deleteScene(t.Scene.ID)
	}

	if t.Gallery != nil && t.shouldCleanGallery(t.Gallery) {
		t.deleteGallery(t.Gallery.ID)
	}

	if t.Image != nil && t.shouldCleanImage(t.Image) {
		t.deleteImage(t.Image.ID)
	}
}

func (t *CleanTask) shouldClean(path string) bool {
	// use image.FileExists for zip file checking
	fileExists := image.FileExists(path)

	if !fileExists || getStashFromPath(path) == nil {
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

	if !matchExtension(s.Path, config.GetVideoExtensions()) {
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

	if !matchExtension(path, config.GetGalleryExtensions()) {
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

	if !matchExtension(s.Path, config.GetImageExtensions()) {
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
	ctx := context.TODO()
	qb := sqlite.NewSceneQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	scene, err := qb.Find(sceneID)
	err = DestroyScene(sceneID, tx)

	if err != nil {
		logger.Errorf("Error deleting scene from database: %s", err.Error())
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Error deleting scene from database: %s", err.Error())
		return
	}

	DeleteGeneratedSceneFiles(scene, t.fileNamingAlgorithm)
}

func (t *CleanTask) deleteGallery(galleryID int) {
	ctx := context.TODO()
	qb := sqlite.NewGalleryQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	err := qb.Destroy(galleryID, tx)

	if err != nil {
		logger.Errorf("Error deleting gallery from database: %s", err.Error())
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Error deleting gallery from database: %s", err.Error())
		return
	}
}

func (t *CleanTask) deleteImage(imageID int) {
	ctx := context.TODO()
	qb := sqlite.NewImageQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	err := qb.Destroy(imageID, tx)

	if err != nil {
		logger.Errorf("Error deleting image from database: %s", err.Error())
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Error deleting image from database: %s", err.Error())
		return
	}

	pathErr := os.Remove(GetInstance().Paths.Generated.GetThumbnailPath(t.Image.Checksum, models.DefaultGthumbWidth)) // remove cache dir of gallery
	if pathErr != nil {
		logger.Errorf("Error deleting thumbnail image from cache: %s", pathErr)
	}
}

func (t *CleanTask) fileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}

	// handle if error is something else
	if err != nil {
		return false, err
	}

	return !info.IsDir(), nil
}

func getStashFromPath(pathToCheck string) *models.StashConfig {
	for _, s := range config.GetStashPaths() {
		rel, error := filepath.Rel(s.Path, filepath.Dir(pathToCheck))

		if error == nil {
			if !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return s
			}
		}

	}
	return nil
}

func getStashFromDirPath(pathToCheck string) *models.StashConfig {
	for _, s := range config.GetStashPaths() {
		rel, error := filepath.Rel(s.Path, pathToCheck)

		if error == nil {
			if !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return s
			}
		}

	}
	return nil
}
