package manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type AutoTagTask struct {
	paths      []string
	txnManager models.TransactionManager
}

type AutoTagPerformerTask struct {
	AutoTagTask
	performer *models.Performer
}

func (t *AutoTagPerformerTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagPerformer()
}

func (t *AutoTagTask) getQueryRegex(name string) string {
	const separatorChars = `.\-_ `
	// handle path separators
	const separator = `[` + separatorChars + `]`

	ret := strings.Replace(name, " ", separator+"*", -1)
	ret = `(?:^|_|[^\w\d])` + ret + `(?:$|_|[^\w\d])`
	return ret
}

func (t *AutoTagPerformerTask) autoTagPerformer() {
	regex := t.getQueryRegex(t.performer.Name.String)

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()

		scenes, err := qb.QueryForAutoTag(regex, t.paths)

		if err != nil {
			return fmt.Errorf("Error querying scenes with regex '%s': %s", regex, err.Error())
		}

		for _, s := range scenes {
			added, err := scene.AddPerformer(qb, s.ID, t.performer.ID)

			if err != nil {
				return fmt.Errorf("Error adding performer '%s' to scene '%s': %s", t.performer.Name.String, s.GetTitle(), err.Error())
			}

			if added {
				logger.Infof("Added performer '%s' to scene '%s'", t.performer.Name.String, s.GetTitle())
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}

type AutoTagStudioTask struct {
	AutoTagTask
	studio *models.Studio
}

func (t *AutoTagStudioTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagStudio()
}

func (t *AutoTagStudioTask) autoTagStudio() {
	regex := t.getQueryRegex(t.studio.Name.String)

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()
		scenes, err := qb.QueryForAutoTag(regex, t.paths)

		if err != nil {
			return fmt.Errorf("Error querying scenes with regex '%s': %s", regex, err.Error())
		}

		for _, s := range scenes {
			// #306 - don't overwrite studio if already present
			if s.StudioID.Valid {
				// don't modify
				continue
			}

			logger.Infof("Adding studio '%s' to scene '%s'", t.studio.Name.String, s.GetTitle())

			// set the studio id
			studioID := sql.NullInt64{Int64: int64(t.studio.ID), Valid: true}
			scenePartial := models.ScenePartial{
				ID:       s.ID,
				StudioID: &studioID,
			}

			if _, err := qb.Update(scenePartial); err != nil {
				return fmt.Errorf("Error adding studio to scene: %s", err.Error())
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}

type AutoTagTagTask struct {
	AutoTagTask
	tag *models.Tag
}

func (t *AutoTagTagTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagTag()
}

func (t *AutoTagTagTask) autoTagTag() {
	regex := t.getQueryRegex(t.tag.Name)

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()
		scenes, err := qb.QueryForAutoTag(regex, t.paths)

		if err != nil {
			return fmt.Errorf("Error querying scenes with regex '%s': %s", regex, err.Error())
		}

		for _, s := range scenes {
			added, err := scene.AddTag(qb, s.ID, t.tag.ID)

			if err != nil {
				return fmt.Errorf("Error adding tag '%s' to scene '%s': %s", t.tag.Name, s.GetTitle(), err.Error())
			}

			if added {
				logger.Infof("Added tag '%s' to scene '%s'", t.tag.Name, s.GetTitle())
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}
