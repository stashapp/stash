// +build integration

package models_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func TestPerformerFindBySceneID(t *testing.T) {
	pqb := models.NewPerformerQueryBuilder()
	sceneID := sceneIDs[sceneIdxWithPerformer]

	performers, err := pqb.FindBySceneID(sceneID, nil)

	if err != nil {
		t.Fatalf("Error finding performer: %s", err.Error())
	}

	assert.Equal(t, 1, len(performers))
	performer := performers[0]

	assert.Equal(t, getPerformerStringValue(performerIdxWithScene, "Name"), performer.Name.String)

	performers, err = pqb.FindBySceneID(0, nil)

	if err != nil {
		t.Fatalf("Error finding performer: %s", err.Error())
	}

	assert.Equal(t, 0, len(performers))
}

func TestPerformerFindNameBySceneID(t *testing.T) {
	pqb := models.NewPerformerQueryBuilder()
	sceneID := sceneIDs[sceneIdxWithPerformer]

	performers, err := pqb.FindNameBySceneID(sceneID, nil)

	if err != nil {
		t.Fatalf("Error finding performer: %s", err.Error())
	}

	assert.Equal(t, 1, len(performers))
	performer := performers[0]

	assert.Equal(t, getPerformerStringValue(performerIdxWithScene, "Name"), performer.Name.String)

	performers, err = pqb.FindBySceneID(0, nil)

	if err != nil {
		t.Fatalf("Error finding performer: %s", err.Error())
	}

	assert.Equal(t, 0, len(performers))
}

func TestPerformerFindByNames(t *testing.T) {
	var names []string

	pqb := models.NewPerformerQueryBuilder()

	names = append(names, performerNames[performerIdxWithScene]) // find performers by names

	performers, err := pqb.FindByNames(names, nil, false)
	if err != nil {
		t.Fatalf("Error finding performers: %s", err.Error())
	}
	assert.Len(t, performers, 1)
	assert.Equal(t, performerNames[performerIdxWithScene], performers[0].Name.String)

	performers, err = pqb.FindByNames(names, nil, true) // find performers by names nocase
	if err != nil {
		t.Fatalf("Error finding performers: %s", err.Error())
	}
	assert.Len(t, performers, 2) // performerIdxWithScene and performerIdxWithDupName
	assert.Equal(t, strings.ToLower(performerNames[performerIdxWithScene]), strings.ToLower(performers[0].Name.String))
	assert.Equal(t, strings.ToLower(performerNames[performerIdxWithScene]), strings.ToLower(performers[1].Name.String))

	names = append(names, performerNames[performerIdx1WithScene]) // find performers by names ( 2 names )

	performers, err = pqb.FindByNames(names, nil, false)
	if err != nil {
		t.Fatalf("Error finding performers: %s", err.Error())
	}
	assert.Len(t, performers, 2) // performerIdxWithScene and performerIdx1WithScene
	assert.Equal(t, performerNames[performerIdxWithScene], performers[0].Name.String)
	assert.Equal(t, performerNames[performerIdx1WithScene], performers[1].Name.String)

	performers, err = pqb.FindByNames(names, nil, true) // find performers by names ( 2 names nocase)
	if err != nil {
		t.Fatalf("Error finding performers: %s", err.Error())
	}
	assert.Len(t, performers, 4) // performerIdxWithScene and performerIdxWithDupName , performerIdx1WithScene and performerIdx1WithDupName
	assert.Equal(t, performerNames[performerIdxWithScene], performers[0].Name.String)
	assert.Equal(t, performerNames[performerIdx1WithScene], performers[1].Name.String)
	assert.Equal(t, performerNames[performerIdx1WithDupName], performers[2].Name.String)
	assert.Equal(t, performerNames[performerIdxWithDupName], performers[3].Name.String)

}

func TestPerformerUpdatePerformerImage(t *testing.T) {
	qb := models.NewPerformerQueryBuilder()

	// create performer to test against
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	const name = "TestPerformerUpdatePerformerImage"
	performer := models.Performer{
		Name:     sql.NullString{String: name, Valid: true},
		Checksum: utils.MD5FromString(name),
		Favorite: sql.NullBool{Bool: false, Valid: true},
	}
	created, err := qb.Create(performer, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating performer: %s", err.Error())
	}

	image := []byte("image")
	err = qb.UpdatePerformerImage(created.ID, image, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error updating performer image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	// ensure image set
	storedImage, err := qb.GetPerformerImage(created.ID, nil)
	if err != nil {
		t.Fatalf("Error getting image: %s", err.Error())
	}
	assert.Equal(t, storedImage, image)

	// set nil image
	tx = database.DB.MustBeginTx(ctx, nil)
	err = qb.UpdatePerformerImage(created.ID, nil, tx)
	if err == nil {
		t.Fatalf("Expected error setting nil image")
	}

	tx.Rollback()
}

func TestPerformerDestroyPerformerImage(t *testing.T) {
	qb := models.NewPerformerQueryBuilder()

	// create performer to test against
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	const name = "TestPerformerDestroyPerformerImage"
	performer := models.Performer{
		Name:     sql.NullString{String: name, Valid: true},
		Checksum: utils.MD5FromString(name),
		Favorite: sql.NullBool{Bool: false, Valid: true},
	}
	created, err := qb.Create(performer, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating performer: %s", err.Error())
	}

	image := []byte("image")
	err = qb.UpdatePerformerImage(created.ID, image, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error updating performer image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	tx = database.DB.MustBeginTx(ctx, nil)

	err = qb.DestroyPerformerImage(created.ID, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error destroying performer image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	// image should be nil
	storedImage, err := qb.GetPerformerImage(created.ID, nil)
	if err != nil {
		t.Fatalf("Error getting image: %s", err.Error())
	}
	assert.Nil(t, storedImage)
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
