// +build integration

package models_test

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestStudioFindByName(t *testing.T) {

	sqb := models.NewStudioQueryBuilder()

	name := studioNames[studioIdxWithScene] // find a studio by name

	studio, err := sqb.FindByName(name, nil, false)

	if err != nil {
		t.Fatalf("Error finding studios: %s", err.Error())
	}

	assert.Equal(t, studioNames[studioIdxWithScene], studio.Name.String)

	name = studioNames[studioIdxWithDupName] // find a studio by name nocase

	studio, err = sqb.FindByName(name, nil, true)

	if err != nil {
		t.Fatalf("Error finding studios: %s", err.Error())
	}
	// studioIdxWithDupName and studioIdxWithScene should have similar names ( only diff should be Name vs NaMe)
	//studio.Name should match with studioIdxWithScene since its ID is before studioIdxWithDupName
	assert.Equal(t, studioNames[studioIdxWithScene], studio.Name.String)
	//studio.Name should match with studioIdxWithDupName if the check is not case sensitive
	assert.Equal(t, strings.ToLower(studioNames[studioIdxWithDupName]), strings.ToLower(studio.Name.String))

}

func TestStudioQueryParent(t *testing.T) {
	sqb := models.NewStudioQueryBuilder()
	studioCriterion := models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(studioIDs[studioIdxWithChildStudio]),
		},
		Modifier: models.CriterionModifierIncludes,
	}

	studioFilter := models.StudioFilterType{
		Parents: &studioCriterion,
	}

	studios, _ := sqb.Query(&studioFilter, nil)

	assert.Len(t, studios, 1)

	// ensure id is correct
	assert.Equal(t, sceneIDs[studioIdxWithParentStudio], studios[0].ID)

	studioCriterion = models.MultiCriterionInput{
		Value: []string{
			strconv.Itoa(studioIDs[studioIdxWithChildStudio]),
		},
		Modifier: models.CriterionModifierExcludes,
	}

	q := getStudioStringValue(studioIdxWithParentStudio, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	studios, _ = sqb.Query(&studioFilter, &findFilter)
	assert.Len(t, studios, 0)
}

func TestStudioDestroyParent(t *testing.T) {
	const parentName = "parent"
	const childName = "child"

	// create parent and child studios
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	createdParent, err := createStudio(tx, parentName, nil)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating parent studio: %s", err.Error())
	}

	parentID := int64(createdParent.ID)
	createdChild, err := createStudio(tx, childName, &parentID)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating child studio: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	sqb := models.NewStudioQueryBuilder()

	// destroy the parent
	tx = database.DB.MustBeginTx(ctx, nil)

	err = sqb.Destroy(strconv.Itoa(createdParent.ID), tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error destroying parent studio: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	// destroy the child
	tx = database.DB.MustBeginTx(ctx, nil)

	err = sqb.Destroy(strconv.Itoa(createdChild.ID), tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error destroying child studio: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}
}

func TestStudioFindChildren(t *testing.T) {
	sqb := models.NewStudioQueryBuilder()

	studios, err := sqb.FindChildren(studioIDs[studioIdxWithChildStudio], nil)

	if err != nil {
		t.Fatalf("error calling FindChildren: %s", err.Error())
	}

	assert.Len(t, studios, 1)
	assert.Equal(t, studioIDs[studioIdxWithParentStudio], studios[0].ID)

	studios, err = sqb.FindChildren(0, nil)

	if err != nil {
		t.Fatalf("error calling FindChildren: %s", err.Error())
	}

	assert.Len(t, studios, 0)
}

func TestStudioUpdateClearParent(t *testing.T) {
	const parentName = "clearParent_parent"
	const childName = "clearParent_child"

	// create parent and child studios
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	createdParent, err := createStudio(tx, parentName, nil)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating parent studio: %s", err.Error())
	}

	parentID := int64(createdParent.ID)
	createdChild, err := createStudio(tx, childName, &parentID)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating child studio: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	sqb := models.NewStudioQueryBuilder()

	// clear the parent id from the child
	tx = database.DB.MustBeginTx(ctx, nil)

	updatePartial := models.StudioPartial{
		ID:       createdChild.ID,
		ParentID: &sql.NullInt64{Valid: false},
	}

	updatedStudio, err := sqb.Update(updatePartial, tx)

	if err != nil {
		tx.Rollback()
		t.Fatalf("Error updated studio: %s", err.Error())
	}

	if updatedStudio.ParentID.Valid {
		t.Error("updated studio has parent ID set")
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}
}

func TestStudioUpdateStudioImage(t *testing.T) {
	qb := models.NewStudioQueryBuilder()

	// create performer to test against
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	const name = "TestStudioUpdateStudioImage"
	created, err := createStudio(tx, name, nil)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating studio: %s", err.Error())
	}

	image := []byte("image")
	err = qb.UpdateStudioImage(created.ID, image, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error updating studio image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	// ensure image set
	storedImage, err := qb.GetStudioImage(created.ID, nil)
	if err != nil {
		t.Fatalf("Error getting image: %s", err.Error())
	}
	assert.Equal(t, storedImage, image)

	// set nil image
	tx = database.DB.MustBeginTx(ctx, nil)
	err = qb.UpdateStudioImage(created.ID, nil, tx)
	if err == nil {
		t.Fatalf("Expected error setting nil image")
	}

	tx.Rollback()
}

func TestStudioDestroyStudioImage(t *testing.T) {
	qb := models.NewStudioQueryBuilder()

	// create performer to test against
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	const name = "TestStudioDestroyStudioImage"
	created, err := createStudio(tx, name, nil)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating studio: %s", err.Error())
	}

	image := []byte("image")
	err = qb.UpdateStudioImage(created.ID, image, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error updating studio image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	tx = database.DB.MustBeginTx(ctx, nil)

	err = qb.DestroyStudioImage(created.ID, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error destroying studio image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	// image should be nil
	storedImage, err := qb.GetStudioImage(created.ID, nil)
	if err != nil {
		t.Fatalf("Error getting image: %s", err.Error())
	}
	assert.Nil(t, storedImage)
}

// TODO Create
// TODO Update
// TODO Destroy
// TODO Find
// TODO FindBySceneID
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
