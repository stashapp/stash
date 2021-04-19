package manager

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

const autoTagSeparatorChars = `.\-_ `

type AutoTagFileTask struct {
	Scene      *models.Scene
	txnManager models.TransactionManager

	Performers bool
	Tags       bool
	Studios    bool
}

func (t *AutoTagFileTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTag()
}

func (t *AutoTagFileTask) getQueryString() string {
	ret := t.Scene.Path
	ret = strings.ReplaceAll(ret, string(filepath.Separator), " ")

	// handle path separators
	const separator = `[` + autoTagSeparatorChars + `]+`
	re := regexp.MustCompile(separator)
	ret = re.ReplaceAllString(ret, " ")

	return ret
}

func (t *AutoTagFileTask) nameMatchesPath(name string) bool {
	// handle path separators
	const separator = `[` + autoTagSeparatorChars + `]`

	reStr := strings.Replace(name, " ", separator+"*", -1)
	reStr = `(?:^|_|[^\w\d])` + reStr + `(?:$|_|[^\w\d])`

	re := regexp.MustCompile(reStr)
	return re.MatchString(t.Scene.Path)
}

func (t *AutoTagFileTask) autoTag() {
	queryStr := t.getQueryString()
	perPage := -1
	filter := &models.FindFilterType{
		Q:       &queryStr,
		PerPage: &perPage,
	}

	if err := t.txnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		if t.Performers {
			if err := t.autoTagPerformers(r, filter); err != nil {
				return err
			}
		}
		if t.Tags {
			if err := t.autoTagTags(r, filter); err != nil {
				return err
			}
		}
		if t.Studios {
			if err := t.autoTagStudios(r, filter); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		logger.Error(err.Error())
	}
}

func (t *AutoTagFileTask) autoTagPerformers(r models.Repository, filter *models.FindFilterType) error {
	performers, _, err := r.Performer().Query(nil, filter)

	if err != nil {
		return fmt.Errorf("error querying performers for auto-tag: %s", err.Error())
	}

	for _, p := range performers {
		if t.nameMatchesPath(p.Name.String) {
			added, err := scene.AddPerformer(r.Scene(), t.Scene.ID, p.ID)

			if err != nil {
				return fmt.Errorf("error adding performer '%s' to scene '%s': %s", p.Name.String, t.Scene.GetTitle(), err.Error())
			}

			if added {
				logger.Infof("Added performer '%s' to scene '%s'", p.Name.String, t.Scene.GetTitle())
			}
		}
	}

	return nil
}

func (t *AutoTagFileTask) autoTagStudios(r models.Repository, filter *models.FindFilterType) error {
	// #306 - don't overwrite studio if already present
	if t.Scene.StudioID.Valid {
		// don't modify
		return nil
	}

	studios, _, err := r.Studio().Query(nil, filter)

	if err != nil {
		return fmt.Errorf("error querying studios for auto-tag: %s", err.Error())
	}

	for _, s := range studios {
		if t.nameMatchesPath(s.Name.String) {
			logger.Infof("Adding studio '%s' to scene '%s'", s.Name.String, t.Scene.GetTitle())

			// set the studio id
			studioID := sql.NullInt64{Int64: int64(s.ID), Valid: true}
			scenePartial := models.ScenePartial{
				ID:       t.Scene.ID,
				StudioID: &studioID,
			}

			if _, err := r.Scene().Update(scenePartial); err != nil {
				return fmt.Errorf("error adding studio to scene: %s", err.Error())
			}

			// only set the first one that matches
			return nil
		}
	}

	return nil
}

func (t *AutoTagFileTask) autoTagTags(r models.Repository, filter *models.FindFilterType) error {
	tags, _, err := r.Tag().Query(nil, filter)

	if err != nil {
		return fmt.Errorf("error querying tags for auto-tag: %s", err.Error())
	}

	for _, tag := range tags {
		if t.nameMatchesPath(tag.Name) {
			added, err := scene.AddTag(r.Scene(), t.Scene.ID, tag.ID)

			if err != nil {
				return fmt.Errorf("error adding tag '%s' to scene '%s': %s", tag.Name, t.Scene.GetTitle(), err.Error())
			}

			if added {
				logger.Infof("Added tag '%s' to scene '%s'", tag.Name, t.Scene.GetTitle())
			}
		}
	}

	return nil
}
