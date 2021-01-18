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

type AutoTagPerformerTask struct {
	performer  *models.Performer
	txnManager models.TransactionManager
}

func (t *AutoTagPerformerTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagPerformer()
}

func getQueryRegex(name string) string {
	const separatorChars = `.\-_ `
	// handle path separators
	const separator = `[` + separatorChars + `]`

	ret := strings.Replace(name, " ", separator+"*", -1)
	ret = `(?:^|_|[^\w\d])` + ret + `(?:$|_|[^\w\d])`
	return ret
}

func (t *AutoTagPerformerTask) autoTagPerformer() {
	regex := getQueryRegex(t.performer.Name.String)

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()
		const ignoreOrganized = true
		scenes, err := qb.QueryAllByPathRegex(regex, ignoreOrganized)

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
	studio     *models.Studio
	txnManager models.TransactionManager
}

func (t *AutoTagStudioTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagStudio()
}

func (t *AutoTagStudioTask) autoTagStudio() {
	regex := getQueryRegex(t.studio.Name.String)

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()
		const ignoreOrganized = true
		scenes, err := qb.QueryAllByPathRegex(regex, ignoreOrganized)

		if err != nil {
			return fmt.Errorf("Error querying scenes with regex '%s': %s", regex, err.Error())
		}

		for _, scene := range scenes {
			// #306 - don't overwrite studio if already present
			if scene.StudioID.Valid {
				// don't modify
				continue
			}

			logger.Infof("Adding studio '%s' to scene '%s'", t.studio.Name.String, scene.GetTitle())

			// set the studio id
			studioID := sql.NullInt64{Int64: int64(t.studio.ID), Valid: true}
			scenePartial := models.ScenePartial{
				ID:       scene.ID,
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
	tag        *models.Tag
	txnManager models.TransactionManager
}

func (t *AutoTagTagTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagTag()
}

func (t *AutoTagTagTask) autoTagTag() {
	regex := getQueryRegex(t.tag.Name)

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()
		const ignoreOrganized = true
		scenes, err := qb.QueryAllByPathRegex(regex, ignoreOrganized)

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
