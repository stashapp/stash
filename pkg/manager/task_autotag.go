package manager

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type AutoTagPerformerTask struct {
	performer *models.Performer
}

func (t *AutoTagPerformerTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagPerformer()
}

func getQueryRegex(name string) string {
	const separatorChars = `.\-_ `
	// handle path separators
	const endSeparatorChars = separatorChars + `\\/`
	const separator = `[` + separatorChars + `]`
	const endSeparator = `[` + endSeparatorChars + `]`

	ret := strings.Replace(name, " ", separator+"*", -1)
	ret = "(?:^|" + endSeparator + "+)" + ret + "(?:$|" + endSeparator + "+)"
	return ret
}

func (t *AutoTagPerformerTask) autoTagPerformer() {
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	regex := getQueryRegex(t.performer.Name.String)

	scenes, err := qb.QueryAllByPathRegex(regex)

	if err != nil {
		logger.Infof("Error querying scenes with regex '%s': %s", regex, err.Error())
		return
	}

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	for _, scene := range scenes {
		added, err := jqb.AddPerformerScene(scene.ID, t.performer.ID, tx)

		if err != nil {
			logger.Infof("Error adding performer '%s' to scene '%s': %s", t.performer.Name.String, scene.GetTitle(), err.Error())
			tx.Rollback()
			return
		}

		if added {
			logger.Infof("Added performer '%s' to scene '%s'", t.performer.Name.String, scene.GetTitle())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Infof("Error adding performer to scene: %s", err.Error())
		return
	}
}

type AutoTagStudioTask struct {
	studio *models.Studio
}

func (t *AutoTagStudioTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagStudio()
}

func (t *AutoTagStudioTask) autoTagStudio() {
	qb := models.NewSceneQueryBuilder()

	regex := getQueryRegex(t.studio.Name.String)

	scenes, err := qb.QueryAllByPathRegex(regex)

	if err != nil {
		logger.Infof("Error querying scenes with regex '%s': %s", regex, err.Error())
		return
	}

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

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

		_, err := qb.Update(scenePartial, tx)

		if err != nil {
			logger.Infof("Error adding studio to scene: %s", err.Error())
			tx.Rollback()
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Infof("Error adding studio to scene: %s", err.Error())
		return
	}
}

type AutoTagDvdTask struct {
	dvd *models.Dvd
}

func (t *AutoTagDvdTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagDvd()
}

func (t *AutoTagDvdTask) autoTagDvd() {
	qb := models.NewSceneQueryBuilder()

	regex := getQueryRegex(t.dvd.Name.String)

	scenes, err := qb.QueryAllByPathRegex(regex)

	if err != nil {
		logger.Infof("Error querying scenes with regex '%s': %s", regex, err.Error())
		return
	}

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	for _, scene := range scenes {
		if scene.DvdID.Int64 == int64(t.dvd.ID) {
			// don't modify
			continue
		}

		logger.Infof("Adding dvd '%s' to scene '%s'", t.dvd.Name.String, scene.GetTitle())

		// set the dvd id
		dvdID := sql.NullInt64{Int64: int64(t.dvd.ID), Valid: true}
		scenePartial := models.ScenePartial{
			ID:    scene.ID,
			DvdID: &dvdID,
		}

		_, err := qb.Update(scenePartial, tx)

		if err != nil {
			logger.Infof("Error adding dvd to scene: %s", err.Error())
			tx.Rollback()
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Infof("Error adding dvd to scene: %s", err.Error())
		return
	}
}

type AutoTagTagTask struct {
	tag *models.Tag
}

func (t *AutoTagTagTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagTag()
}

func (t *AutoTagTagTask) autoTagTag() {
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	regex := getQueryRegex(t.tag.Name)

	scenes, err := qb.QueryAllByPathRegex(regex)

	if err != nil {
		logger.Infof("Error querying scenes with regex '%s': %s", regex, err.Error())
		return
	}

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	for _, scene := range scenes {
		added, err := jqb.AddSceneTag(scene.ID, t.tag.ID, tx)

		if err != nil {
			logger.Infof("Error adding tag '%s' to scene '%s': %s", t.tag.Name, scene.GetTitle(), err.Error())
			tx.Rollback()
			return
		}

		if added {
			logger.Infof("Added tag '%s' to scene '%s'", t.tag.Name, scene.GetTitle())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Infof("Error adding tag to scene: %s", err.Error())
		return
	}
}
