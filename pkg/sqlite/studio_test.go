// +build integration

package sqlite_test

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestStudioFindByName(t *testing.T) {
	withTxn(func(r models.Repository) error {
		sqb := r.Studio()

		name := studioNames[studioIdxWithScene] // find a studio by name

		studio, err := sqb.FindByName(name, false)

		if err != nil {
			t.Errorf("Error finding studios: %s", err.Error())
		}

		assert.Equal(t, studioNames[studioIdxWithScene], studio.Name.String)

		name = studioNames[studioIdxWithDupName] // find a studio by name nocase

		studio, err = sqb.FindByName(name, true)

		if err != nil {
			t.Errorf("Error finding studios: %s", err.Error())
		}
		// studioIdxWithDupName and studioIdxWithScene should have similar names ( only diff should be Name vs NaMe)
		//studio.Name should match with studioIdxWithScene since its ID is before studioIdxWithDupName
		assert.Equal(t, studioNames[studioIdxWithScene], studio.Name.String)
		//studio.Name should match with studioIdxWithDupName if the check is not case sensitive
		assert.Equal(t, strings.ToLower(studioNames[studioIdxWithDupName]), strings.ToLower(studio.Name.String))

		return nil
	})
}

func TestStudioQueryForAutoTag(t *testing.T) {
	withTxn(func(r models.Repository) error {
		tqb := r.Studio()

		name := studioNames[studioIdxWithScene] // find a studio by name

		studios, err := tqb.QueryForAutoTag([]string{name})

		if err != nil {
			t.Errorf("Error finding studios: %s", err.Error())
		}

		assert.Len(t, studios, 2)
		assert.Equal(t, strings.ToLower(studioNames[studioIdxWithScene]), strings.ToLower(studios[0].Name.String))
		assert.Equal(t, strings.ToLower(studioNames[studioIdxWithScene]), strings.ToLower(studios[1].Name.String))

		return nil
	})
}

func TestStudioQueryParent(t *testing.T) {
	withTxn(func(r models.Repository) error {
		sqb := r.Studio()
		studioCriterion := models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(studioIDs[studioIdxWithChildStudio]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		studioFilter := models.StudioFilterType{
			Parents: &studioCriterion,
		}

		studios, _, err := sqb.Query(&studioFilter, nil)
		if err != nil {
			t.Errorf("Error querying studio: %s", err.Error())
		}

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

		studios, _, err = sqb.Query(&studioFilter, &findFilter)
		if err != nil {
			t.Errorf("Error querying studio: %s", err.Error())
		}
		assert.Len(t, studios, 0)

		return nil
	})
}

func TestStudioDestroyParent(t *testing.T) {
	const parentName = "parent"
	const childName = "child"

	// create parent and child studios
	if err := withTxn(func(r models.Repository) error {
		createdParent, err := createStudio(r.Studio(), parentName, nil)
		if err != nil {
			return fmt.Errorf("Error creating parent studio: %s", err.Error())
		}

		parentID := int64(createdParent.ID)
		createdChild, err := createStudio(r.Studio(), childName, &parentID)
		if err != nil {
			return fmt.Errorf("Error creating child studio: %s", err.Error())
		}

		sqb := r.Studio()

		// destroy the parent
		err = sqb.Destroy(createdParent.ID)
		if err != nil {
			return fmt.Errorf("Error destroying parent studio: %s", err.Error())
		}

		// destroy the child
		err = sqb.Destroy(createdChild.ID)
		if err != nil {
			return fmt.Errorf("Error destroying child studio: %s", err.Error())
		}

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestStudioFindChildren(t *testing.T) {
	withTxn(func(r models.Repository) error {
		sqb := r.Studio()

		studios, err := sqb.FindChildren(studioIDs[studioIdxWithChildStudio])

		if err != nil {
			t.Errorf("error calling FindChildren: %s", err.Error())
		}

		assert.Len(t, studios, 1)
		assert.Equal(t, studioIDs[studioIdxWithParentStudio], studios[0].ID)

		studios, err = sqb.FindChildren(0)

		if err != nil {
			t.Errorf("error calling FindChildren: %s", err.Error())
		}

		assert.Len(t, studios, 0)

		return nil
	})
}

func TestStudioUpdateClearParent(t *testing.T) {
	const parentName = "clearParent_parent"
	const childName = "clearParent_child"

	// create parent and child studios
	if err := withTxn(func(r models.Repository) error {
		createdParent, err := createStudio(r.Studio(), parentName, nil)
		if err != nil {
			return fmt.Errorf("Error creating parent studio: %s", err.Error())
		}

		parentID := int64(createdParent.ID)
		createdChild, err := createStudio(r.Studio(), childName, &parentID)
		if err != nil {
			return fmt.Errorf("Error creating child studio: %s", err.Error())
		}

		sqb := r.Studio()

		// clear the parent id from the child
		updatePartial := models.StudioPartial{
			ID:       createdChild.ID,
			ParentID: &sql.NullInt64{Valid: false},
		}

		updatedStudio, err := sqb.Update(updatePartial)

		if err != nil {
			return fmt.Errorf("Error updated studio: %s", err.Error())
		}

		if updatedStudio.ParentID.Valid {
			return errors.New("updated studio has parent ID set")
		}

		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestStudioUpdateStudioImage(t *testing.T) {
	if err := withTxn(func(r models.Repository) error {
		qb := r.Studio()

		// create performer to test against
		const name = "TestStudioUpdateStudioImage"
		created, err := createStudio(r.Studio(), name, nil)
		if err != nil {
			return fmt.Errorf("Error creating studio: %s", err.Error())
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

func TestStudioDestroyStudioImage(t *testing.T) {
	if err := withTxn(func(r models.Repository) error {
		qb := r.Studio()

		// create performer to test against
		const name = "TestStudioDestroyStudioImage"
		created, err := createStudio(r.Studio(), name, nil)
		if err != nil {
			return fmt.Errorf("Error creating studio: %s", err.Error())
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

func TestStudioStashIDs(t *testing.T) {
	if err := withTxn(func(r models.Repository) error {
		qb := r.Studio()

		// create studio to test against
		const name = "TestStudioStashIDs"
		created, err := createStudio(r.Studio(), name, nil)
		if err != nil {
			return fmt.Errorf("Error creating studio: %s", err.Error())
		}

		testStashIDReaderWriter(t, qb, created.ID)
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

func TestStudioQueryURL(t *testing.T) {
	const sceneIdx = 1
	studioURL := getStudioStringValue(sceneIdx, urlField)

	urlCriterion := models.StringCriterionInput{
		Value:    studioURL,
		Modifier: models.CriterionModifierEquals,
	}

	filter := models.StudioFilterType{
		URL: &urlCriterion,
	}

	verifyFn := func(g *models.Studio) {
		t.Helper()
		verifyNullString(t, g.URL, urlCriterion)
	}

	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotEquals
	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierMatchesRegex
	urlCriterion.Value = "studio_.*1_URL"
	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotMatchesRegex
	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierIsNull
	urlCriterion.Value = ""
	verifyStudioQuery(t, filter, verifyFn)

	urlCriterion.Modifier = models.CriterionModifierNotNull
	verifyStudioQuery(t, filter, verifyFn)
}

func verifyStudioQuery(t *testing.T, filter models.StudioFilterType, verifyFn func(s *models.Studio)) {
	withTxn(func(r models.Repository) error {
		t.Helper()
		sqb := r.Studio()

		galleries := queryStudio(t, sqb, &filter, nil)

		// assume it should find at least one
		assert.Greater(t, len(galleries), 0)

		for _, studio := range galleries {
			verifyFn(studio)
		}

		return nil
	})
}

func queryStudio(t *testing.T, sqb models.StudioReader, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) []*models.Studio {
	studios, _, err := sqb.Query(studioFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying studio: %s", err.Error())
	}

	return studios
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
