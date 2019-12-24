package manager

import (
	"context"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CleanTask struct {
	Scene models.Scene
}

func (t *CleanTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.fileExists(t.Scene.Path) && t.pathInStash() {
		logger.Debugf("File Found: %s", t.Scene.Path)
		if matchFile(t.Scene.Path, config.GetExcludes()) {
			logger.Infof("File matched regex. Cleaning: \"%s\"", t.Scene.Path)
			t.deleteScene(t.Scene.ID)
		}
	} else {
		logger.Infof("File not found. Cleaning: \"%s\"", t.Scene.Path)
		t.deleteScene(t.Scene.ID)
	}
}

func (t *CleanTask) deleteScene(sceneID int) {
	ctx := context.TODO()
	qb := models.NewSceneQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)

	scene, err := qb.Find(sceneID)
	err = DestroyScene(sceneID, tx)

	if err != nil {
		logger.Infof("Error deleting scene from database: %s", err.Error())
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Infof("Error deleting scene from database: %s", err.Error())
		return
	}

	DeleteGeneratedSceneFiles(scene)
}

func (t *CleanTask) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (t *CleanTask) pathInStash() bool {
	for _, path := range config.GetStashPaths() {

		rel, error := filepath.Rel(path, filepath.Dir(t.Scene.Path))

		if error == nil {
			if !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				logger.Debugf("File %s belongs to stash path %s", t.Scene.Path, path)
				return true
			}
		}

	}
	logger.Debugf("File %s is out from stash path", t.Scene.Path)
	return false
}
