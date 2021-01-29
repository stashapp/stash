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
	"github.com/stashapp/stash/pkg/utils"
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

func (t *AutoTagTask) getSceneFilter(regex string) *models.SceneFilterType {
	organized := false
	return &models.SceneFilterType{
		Path: &models.StringCriterionInput{
			Value:    "(?i)" + regex,
			Modifier: models.CriterionModifierMatchesRegex,
		},
		Organized: &organized,
	}
}

func (t *AutoTagTask) getQueryFilter() models.QueryFilter {
	return models.QueryFilter{
		All: true,
	}
}

// TODO - this should be replaced in the query filter once multiple conditions
// on the same field is possible
func (t *AutoTagTask) isSceneInPaths(s *models.Scene) bool {
	for _, p := range t.paths {
		if utils.IsPathInDir(p, s.Path) {
			return true
		}
	}

	return false
}

func (t *AutoTagPerformerTask) autoTagPerformer() {
	regex := t.getQueryRegex(t.performer.Name.String)

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Scene()

		sceneFilter := t.getSceneFilter(regex)
		findFilter := t.getQueryFilter()

		scenes, _, err := qb.Query(sceneFilter, findFilter)

		if err != nil {
			return fmt.Errorf("Error querying scenes with regex '%s': %s", regex, err.Error())
		}

		for _, s := range scenes {
			if t.paths == nil || t.isSceneInPaths(s) {
				added, err := scene.AddPerformer(qb, s.ID, t.performer.ID)

				if err != nil {
					return fmt.Errorf("Error adding performer '%s' to scene '%s': %s", t.performer.Name.String, s.GetTitle(), err.Error())
				}

				if added {
					logger.Infof("Added performer '%s' to scene '%s'", t.performer.Name.String, s.GetTitle())
				}
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
		sceneFilter := t.getSceneFilter(regex)
		findFilter := t.getQueryFilter()

		scenes, _, err := qb.Query(sceneFilter, findFilter)

		if err != nil {
			return fmt.Errorf("Error querying scenes with regex '%s': %s", regex, err.Error())
		}

		for _, s := range scenes {
			if t.paths == nil || t.isSceneInPaths(s) {
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
		sceneFilter := t.getSceneFilter(regex)
		findFilter := t.getQueryFilter()

		scenes, _, err := qb.Query(sceneFilter, findFilter)

		if err != nil {
			return fmt.Errorf("Error querying scenes with regex '%s': %s", regex, err.Error())
		}

		for _, s := range scenes {
			if t.paths == nil || t.isSceneInPaths(s) {
				added, err := scene.AddTag(qb, s.ID, t.tag.ID)

				if err != nil {
					return fmt.Errorf("Error adding tag '%s' to scene '%s': %s", t.tag.Name, s.GetTitle(), err.Error())
				}

				if added {
					logger.Infof("Added tag '%s' to scene '%s'", t.tag.Name, s.GetTitle())
				}
			}
		}

		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}
