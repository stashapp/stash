package manager

import (
	"context"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"os"
	"sync"
)

type CleanTask struct {
	Scene models.Scene
}

func (t *CleanTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.fileExists(t.Scene.Path) {
		logger.Debugf("File Found: %s", t.Scene.Path)
	} else {
		logger.Infof("File not found. Cleaning: %s", t.Scene.Path)
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
