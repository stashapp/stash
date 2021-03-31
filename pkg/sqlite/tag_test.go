// +build integration

package sqlite_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMarkerFindBySceneMarkerID(t *testing.T) {
	withTxn(func(r models.Repository) error {
		tqb := r.Tag()

		markerID := markerIDs[markerIdxWithScene]

		tags, err := tqb.FindBySceneMarkerID(markerID)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Len(t, tags, 1)
		assert.Equal(t, tagIDs[tagIdxWithMarker], tags[0].ID)

		tags, err = tqb.FindBySceneMarkerID(0)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Len(t, tags, 0)

		return nil
	})
}

func TestTagFindByName(t *testing.T) {
	withTxn(func(r models.Repository) error {
		tqb := r.Tag()

		name := tagNames[tagIdxWithScene] // find a tag by name

		tag, err := tqb.FindByName(name, false)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}

		assert.Equal(t, tagNames[tagIdxWithScene], tag.Name)

		name = tagNames[tagIdxWithDupName] // find a tag by name nocase

		tag, err = tqb.FindByName(name, true)

		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		// tagIdxWithDupName and tagIdxWithScene should have similar names ( only diff should be Name vs NaMe)
		//tag.Name should match with tagIdxWithScene since its ID is before tagIdxWithDupName
		assert.Equal(t, tagNames[tagIdxWithScene], tag.Name)
		//tag.Name should match with tagIdxWithDupName if the check is not case sensitive
		assert.Equal(t, strings.ToLower(tagNames[tagIdxWithDupName]), strings.ToLower(tag.Name))

		return nil
	})
}

func TestTagFindByNames(t *testing.T) {
	var names []string

	withTxn(func(r models.Repository) error {
		tqb := r.Tag()

		names = append(names, tagNames[tagIdxWithScene]) // find tags by names

		tags, err := tqb.FindByNames(names, false)
		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		assert.Len(t, tags, 1)
		assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)

		tags, err = tqb.FindByNames(names, true) // find tags by names nocase
		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		assert.Len(t, tags, 2) // tagIdxWithScene and tagIdxWithDupName
		assert.Equal(t, strings.ToLower(tagNames[tagIdxWithScene]), strings.ToLower(tags[0].Name))
		assert.Equal(t, strings.ToLower(tagNames[tagIdxWithScene]), strings.ToLower(tags[1].Name))

		names = append(names, tagNames[tagIdx1WithScene]) // find tags by names ( 2 names )

		tags, err = tqb.FindByNames(names, false)
		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		assert.Len(t, tags, 2) // tagIdxWithScene and tagIdx1WithScene
		assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)
		assert.Equal(t, tagNames[tagIdx1WithScene], tags[1].Name)

		tags, err = tqb.FindByNames(names, true) // find tags by names ( 2 names nocase)
		if err != nil {
			t.Errorf("Error finding tags: %s", err.Error())
		}
		assert.Len(t, tags, 4) // tagIdxWithScene and tagIdxWithDupName , tagIdx1WithScene and tagIdx1WithDupName
		assert.Equal(t, tagNames[tagIdxWithScene], tags[0].Name)
		assert.Equal(t, tagNames[tagIdx1WithScene], tags[1].Name)
		assert.Equal(t, tagNames[tagIdx1WithDupName], tags[2].Name)
		assert.Equal(t, tagNames[tagIdxWithDupName], tags[3].Name)

		return nil
	})
}

func TestTagQueryIsMissingImage(t *testing.T) {
	withTxn(func(r models.Repository) error {
		qb := r.Tag()
		isMissing := "image"
		tagFilter := models.TagFilterType{
			IsMissing: &isMissing,
		}

		q := getTagStringValue(tagIdxWithCoverImage, "name")
		findFilter := models.FindFilterType{
			Q: &q,
		}

		tags, _, err := qb.Query(&tagFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		assert.Len(t, tags, 0)

		findFilter.Q = nil
		tags, _, err = qb.Query(&tagFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		// ensure non of the ids equal the one with image
		for _, tag := range tags {
			assert.NotEqual(t, tagIDs[tagIdxWithCoverImage], tag.ID)
		}

		return nil
	})
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
	withTxn(func(r models.Repository) error {
		qb := r.Tag()
		tagFilter := models.TagFilterType{
			SceneCount: &sceneCountCriterion,
		}

		tags, _, err := qb.Query(&tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagSceneCount(tag.ID)),
				Valid: true,
			}, sceneCountCriterion)
		}

		return nil
	})
}

// disabled due to performance issues

// func TestTagQueryMarkerCount(t *testing.T) {
// 	countCriterion := models.IntCriterionInput{
// 		Value:    1,
// 		Modifier: models.CriterionModifierEquals,
// 	}

// 	verifyTagMarkerCount(t, countCriterion)

// 	countCriterion.Modifier = models.CriterionModifierNotEquals
// 	verifyTagMarkerCount(t, countCriterion)

// 	countCriterion.Modifier = models.CriterionModifierLessThan
// 	verifyTagMarkerCount(t, countCriterion)

// 	countCriterion.Value = 0
// 	countCriterion.Modifier = models.CriterionModifierGreaterThan
// 	verifyTagMarkerCount(t, countCriterion)
// }

func verifyTagMarkerCount(t *testing.T, markerCountCriterion models.IntCriterionInput) {
	withTxn(func(r models.Repository) error {
		qb := r.Tag()
		tagFilter := models.TagFilterType{
			MarkerCount: &markerCountCriterion,
		}

		tags, _, err := qb.Query(&tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagMarkerCount(tag.ID)),
				Valid: true,
			}, markerCountCriterion)
		}

		return nil
	})
}

func TestTagQueryImageCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagImageCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagImageCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagImageCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagImageCount(t, countCriterion)
}

func verifyTagImageCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(r models.Repository) error {
		qb := r.Tag()
		tagFilter := models.TagFilterType{
			ImageCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(&tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagImageCount(tag.ID)),
				Valid: true,
			}, imageCountCriterion)
		}

		return nil
	})
}

func TestTagQueryGalleryCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagGalleryCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagGalleryCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagGalleryCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagGalleryCount(t, countCriterion)
}

func verifyTagGalleryCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(r models.Repository) error {
		qb := r.Tag()
		tagFilter := models.TagFilterType{
			GalleryCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(&tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagGalleryCount(tag.ID)),
				Valid: true,
			}, imageCountCriterion)
		}

		return nil
	})
}

func TestTagQueryPerformerCount(t *testing.T) {
	countCriterion := models.IntCriterionInput{
		Value:    1,
		Modifier: models.CriterionModifierEquals,
	}

	verifyTagPerformerCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierNotEquals
	verifyTagPerformerCount(t, countCriterion)

	countCriterion.Modifier = models.CriterionModifierLessThan
	verifyTagPerformerCount(t, countCriterion)

	countCriterion.Value = 0
	countCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyTagPerformerCount(t, countCriterion)
}

func verifyTagPerformerCount(t *testing.T, imageCountCriterion models.IntCriterionInput) {
	withTxn(func(r models.Repository) error {
		qb := r.Tag()
		tagFilter := models.TagFilterType{
			PerformerCount: &imageCountCriterion,
		}

		tags, _, err := qb.Query(&tagFilter, nil)
		if err != nil {
			t.Errorf("Error querying tag: %s", err.Error())
		}

		for _, tag := range tags {
			verifyInt64(t, sql.NullInt64{
				Int64: int64(getTagPerformerCount(tag.ID)),
				Valid: true,
			}, imageCountCriterion)
		}

		return nil
	})
}

func TestTagUpdateTagImage(t *testing.T) {
	if err := withTxn(func(r models.Repository) error {
		qb := r.Tag()

		// create tag to test against
		const name = "TestTagUpdateTagImage"
		tag := models.Tag{
			Name: name,
		}
		created, err := qb.Create(tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		image := []byte("image")
		err = qb.UpdateImage(created.ID, image)
		if err != nil {
			return fmt.Errorf("Error updating studio image: %s", err.Error())
		}

		// ensure image set
		storedImage, err := qb.GetImage(created.ID)
		if err != nil {
			return fmt.Errorf("Error getting image: %s", err.Error())
		}
		assert.Equal(t, storedImage, image)

		// set nil image
		err = qb.UpdateImage(created.ID, nil)
		if err == nil {
			return fmt.Errorf("Expected error setting nil image")
		}

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestTagDestroyTagImage(t *testing.T) {
	if err := withTxn(func(r models.Repository) error {
		qb := r.Tag()

		// create performer to test against
		const name = "TestTagDestroyTagImage"
		tag := models.Tag{
			Name: name,
		}
		created, err := qb.Create(tag)
		if err != nil {
			return fmt.Errorf("Error creating tag: %s", err.Error())
		}

		image := []byte("image")
		err = qb.UpdateImage(created.ID, image)
		if err != nil {
			return fmt.Errorf("Error updating studio image: %s", err.Error())
		}

		err = qb.DestroyImage(created.ID)
		if err != nil {
			return fmt.Errorf("Error destroying studio image: %s", err.Error())
		}

		// image should be nil
		storedImage, err := qb.GetImage(created.ID)
		if err != nil {
			return fmt.Errorf("Error getting image: %s", err.Error())
		}
		assert.Nil(t, storedImage)

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
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
