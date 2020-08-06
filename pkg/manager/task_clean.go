package manager

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
)

type CleanTask struct {
	Scene               *models.Scene
	Gallery             *models.Gallery
	fileNamingAlgorithm models.HashAlgorithm
}

func (t *CleanTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.Scene != nil && t.shouldClean(t.Scene.Path) {
		t.deleteScene(t.Scene.ID)
	}

	if t.Gallery != nil && t.shouldCleanGallery(t.Gallery) {
		t.deleteGallery(t.Gallery.ID)
	}
}

func (t *CleanTask) shouldClean(path string) bool {
	fileExists, err := t.fileExists(path)
	if err != nil {
		logger.Errorf("Error checking existence of %s: %s", path, err.Error())
		return false
	}

	if fileExists && t.pathInStash(path) {
		logger.Debugf("File Found: %s", path)
		if matchFile(path, config.GetExcludes()) {
			logger.Infof("File matched regex. Cleaning: \"%s\"", path)
			return true
		}
	} else {
		logger.Infof("File not found. Cleaning: \"%s\"", path)
		return true
	}

	return false
}

func (t *CleanTask) shouldCleanGallery(g *models.Gallery) bool {
	if t.shouldClean(g.Path) {
		return true
	}

	if t.Gallery.CountFiles() == 0 {
		logger.Infof("Gallery has 0 images. Cleaning: \"%s\"", g.Path)
		return true
	}

	return false
}

func (t *CleanTask) deleteScene(sceneID int) {
	ctx := context.TODO()
	qb := models.NewSceneQueryBuilder()
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
	qb := models.NewGalleryQueryBuilder()
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

	pathErr := os.RemoveAll(paths.GetGthumbDir(t.Gallery.Checksum)) // remove cache dir of gallery
	if pathErr != nil {
		logger.Errorf("Error deleting gallery directory from cache: %s", pathErr)
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

func (t *CleanTask) pathInStash(pathToCheck string) bool {
	for _, path := range config.GetStashPaths() {

		rel, error := filepath.Rel(path, filepath.Dir(pathToCheck))

		if error == nil {
			if !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				logger.Debugf("File %s belongs to stash path %s", pathToCheck, path)
				return true
			}
		}

	}
	logger.Debugf("File %s is out from stash path", pathToCheck)
	return false
}
