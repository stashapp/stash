package manager

import (
	"context"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"os"
	"strconv"
	"sync"
)

type CleanTask struct {
	Scene models.Scene
}

func (t *CleanTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if t.fileExists(t.Scene.Path) {
		logger.Debugf("Found: %s", t.Scene.Path)
	} else {
		logger.Debugf("Deleting missing file: %s", t.Scene.Path)
		t.deleteScene(strconv.Itoa(t.Scene.ID))
	}
}

func (t *CleanTask) deleteScene(id string) {
	ctx := context.TODO()
	qb := models.NewSceneQueryBuilder()
	tx := database.DB.MustBeginTx(ctx, nil)
	defer tx.Commit()
	qb.Destroy(strconv.Itoa(t.Scene.ID), tx)
}

func (t *CleanTask) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
