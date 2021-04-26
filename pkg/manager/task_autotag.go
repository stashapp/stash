package manager

import (
	"context"
	"sync"

	"github.com/stashapp/stash/pkg/autotag"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type autoTagSceneTask struct {
	txnManager models.TransactionManager
	scene      *models.Scene

	performers bool
	studios    bool
	tags       bool
}

func (t *autoTagSceneTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		if t.performers {
			if err := autotag.ScenePerformers(t.scene, r.Scene(), r.Performer()); err != nil {
				return err
			}
		}
		if t.studios {
			if err := autotag.SceneStudios(t.scene, r.Scene(), r.Studio()); err != nil {
				return err
			}
		}
		if t.tags {
			if err := autotag.SceneTags(t.scene, r.Scene(), r.Tag()); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}
