// +build integration

package models_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMarkerFindBySceneMarkerID(t *testing.T) {
	tqb := models.NewTagQueryBuilder()

	markerID := markerIDs[markerIdxWithScene]

	tags, err := tqb.FindBySceneMarkerID(markerID, nil)

	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}

	assert.Len(t, tags, 1)
	assert.Equal(t, tagIDs[tagIdxWithMarker], tags[0].ID)

	tags, err = tqb.FindBySceneMarkerID(0, nil)

	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}

	assert.Len(t, tags, 0)
}

func TestTagFindByName(t *testing.T) {

	tqb := models.NewTagQueryBuilder()

	name := tagNames[tagIdxWithScene] // find a tag by name

	tag, err := tqb.FindByName(name, nil, false)

	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}

	assert.Equal(t, tagNames[tagIdxWithScene], tag.Name)

	name = tagNames[tagIdxWithDupName] // find a tag by name nocase

	tag, err = tqb.FindByName(name, nil, true)

	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	// tagIdxWithDupName and tagIdxWithScene should have similar names ( only diff should be Name vs NaMe)
	//tag.Name should match with tagIdxWithScene since its ID is before tagIdxWithDupName
	assert.Equal(t, tagNames[tagIdxWithScene], tag.Name)
	//tag.Name should match with tagIdxWithDupName if the check is not case sensitive
	assert.Equal(t, strings.ToLower(tagNames[tagIdxWithDupName]), strings.ToLower(tag.Name))

}

func TestTagFindByNames(t *testing.T) {
	var names []string

	tqb := models.NewTagQueryBuilder()

	names = append(names, tagNames[tagIdxWithScene]) // find tags by names

	tags, err := tqb.FindByNames(names, nil, false)
	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	assert.Len(t, tags, 1)
	assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)

	tags, err = tqb.FindByNames(names, nil, true) // find tags by names nocase
	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	assert.Len(t, tags, 2) // tagIdxWithScene and tagIdxWithDupName
	assert.Equal(t, strings.ToLower(tagNames[tagIdxWithScene]), strings.ToLower(tags[0].Name))
	assert.Equal(t, strings.ToLower(tagNames[tagIdxWithScene]), strings.ToLower(tags[1].Name))

	names = append(names, tagNames[tagIdx1WithScene]) // find tags by names ( 2 names )

	tags, err = tqb.FindByNames(names, nil, false)
	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	assert.Len(t, tags, 2) // tagIdxWithScene and tagIdx1WithScene
	assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)
	assert.Equal(t, tagNames[tagIdx1WithScene], tags[1].Name)

	tags, err = tqb.FindByNames(names, nil, true) // find tags by names ( 2 names nocase)
	if err != nil {
		t.Fatalf("Error finding tags: %s", err.Error())
	}
	assert.Len(t, tags, 4) // tagIdxWithScene and tagIdxWithDupName , tagIdx1WithScene and tagIdx1WithDupName
	assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)
	assert.Equal(t, tagNames[tagIdx1WithScene], tags[1].Name)
	assert.Equal(t, tagNames[tagIdx1WithDupName], tags[2].Name)
	assert.Equal(t, tagNames[tagIdxWithDupName], tags[3].Name)

}

func TestTagQueryIsMissingImage(t *testing.T) {
	qb := models.NewTagQueryBuilder()
	isMissing := "image"
	tagFilter := models.TagFilterType{
		IsMissing: &isMissing,
	}

	q := getTagStringValue(tagIdxWithImage, "name")
	findFilter := models.FindFilterType{
		Q: &q,
	}

	tags, _ := qb.Query(&tagFilter, &findFilter)

	assert.Len(t, tags, 0)

	findFilter.Q = nil
	tags, _ = qb.Query(&tagFilter, &findFilter)

	// ensure non of the ids equal the one with image
	for _, tag := range tags {
		assert.NotEqual(t, tagIDs[tagIdxWithImage], tag.ID)
	}
}

func TestTagQuerySceneCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagSceneCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagSceneCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagSceneCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagSceneCount(t, countCriterion)
}

func verifyTagSceneCount(t *testing.T, sceneCountCriterion models.IntCriterionInput) {
	qb := models.NewTagQueryBuilder()
	tagFilter := models.TagFilterType{
		SceneCount: &sceneCountCriterion,
	}

	tags, _ := qb.Query(&tagFilter, nil)

	for _, tag := range tags {
		verifyInt64(t, sql.NullInt64{
			Int64: int64(getTagSceneCount(tag.ID)),
			Valid: true,
		}, sceneCountCriterion)
	}
}

func TestTagQueryMarkerCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagMarkerCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagMarkerCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagMarkerCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagMarkerCount(t, countCriterion)
}

func verifyTagMarkerCount(t *testing.T, markerCountCriterion models.IntCriterionInput) {
	qb := models.NewTagQueryBuilder()
	tagFilter := models.TagFilterType{
		MarkerCount: &markerCountCriterion,
	}

	tags, _ := qb.Query(&tagFilter, nil)

	for _, tag := range tags {
		verifyInt64(t, sql.NullInt64{
			Int64: int64(getTagMarkerCount(tag.ID)),
			Valid: true,
		}, markerCountCriterion)
	}
}

func TestTagUpdateTagImage(t *testing.T) {
	qb := models.NewTagQueryBuilder()

	// create tag to test against
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	const name = "TestTagUpdateTagImage"
	tag := models.Tag{
		Name: name,
	}
	created, err := qb.Create(tag, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating tag: %s", err.Error())
	}

	image := []byte("image")
	err = qb.UpdateTagImage(created.ID, image, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error updating studio image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	// ensure image set
	storedImage, err := qb.GetTagImage(created.ID, nil)
	if err != nil {
		t.Fatalf("Error getting image: %s", err.Error())
	}
	assert.Equal(t, storedImage, image)

	// set nil image
	tx = database.DB.MustBeginTx(ctx, nil)
	err = qb.UpdateTagImage(created.ID, nil, tx)
	if err == nil {
		t.Fatalf("Expected error setting nil image")
	}

	tx.Rollback()
}

func TestTagDestroyTagImage(t *testing.T) {
	qb := models.NewTagQueryBuilder()

	// create performer to test against
	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	const name = "TestTagDestroyTagImage"
	tag := models.Tag{
		Name: name,
	}
	created, err := qb.Create(tag, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error creating tag: %s", err.Error())
	}

	image := []byte("image")
	err = qb.UpdateTagImage(created.ID, image, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error updating studio image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	tx = database.DB.MustBeginTx(ctx, nil)

	err = qb.DestroyTagImage(created.ID, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Error destroying studio image: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		t.Fatalf("Error committing: %s", err.Error())
	}

	// image should be nil
	storedImage, err := qb.GetTagImage(created.ID, nil)
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
// TODO FindBySceneMarkerID
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
