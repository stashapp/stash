// +build integration

package sqlite_test

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func TestPerformerFindBySceneID(t *testing.T) {
	withTxn(func(r models.Repository) error {
		pqb := r.Performer()
		sceneID := sceneIDs[sceneIdxWithPerformer]

		performers, err := pqb.FindBySceneID(sceneID)

		if err != nil {
			t.Errorf("Error finding performer: %s", err.Error())
		}

		assert.Equal(t, 1, len(performers))
		performer := performers[0]

		assert.Equal(t, getPerformerStringValue(performerIdxWithScene, "Name"), performer.Name.String)

		performers, err = pqb.FindBySceneID(0)

		if err != nil {
			t.Errorf("Error finding performer: %s", err.Error())
		}

		assert.Equal(t, 0, len(performers))

		return nil
	})
}

func TestPerformerFindByNames(t *testing.T) {
	getNames := func(p []*models.Performer) []string {
		var ret []string
		for _, pp := range p {
			ret = append(ret, pp.Name.String)
		}
		return ret
	}

	withTxn(func(r models.Repository) error {
		var names []string

		pqb := r.Performer()

		names = append(names, performerNames[performerIdxWithScene]) // find performers by names

		performers, err := pqb.FindByNames(names, false)
		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}
		assert.Len(t, performers, 1)
		assert.Equal(t, performerNames[performerIdxWithScene], performers[0].Name.String)

		performers, err = pqb.FindByNames(names, true) // find performers by names nocase
		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}
		assert.Len(t, performers, 2) // performerIdxWithScene and performerIdxWithDupName
		assert.Equal(t, strings.ToLower(performerNames[performerIdxWithScene]), strings.ToLower(performers[0].Name.String))
		assert.Equal(t, strings.ToLower(performerNames[performerIdxWithScene]), strings.ToLower(performers[1].Name.String))

		names = append(names, performerNames[performerIdx1WithScene]) // find performers by names ( 2 names )

		performers, err = pqb.FindByNames(names, false)
		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}
		retNames := getNames(performers)
		assert.Equal(t, names, retNames)

		performers, err = pqb.FindByNames(names, true) // find performers by names ( 2 names nocase)
		if err != nil {
			t.Errorf("Error finding performers: %s", err.Error())
		}
		retNames = getNames(performers)
		assert.Equal(t, []string{
			performerNames[performerIdxWithScene],
			performerNames[performerIdx1WithScene],
			performerNames[performerIdx1WithDupName],
			performerNames[performerIdxWithDupName],
		}, retNames)

		return nil
	})
}

func TestPerformerUpdatePerformerImage(t *testing.T) {
	if err := withTxn(func(r models.Repository) error {
		qb := r.Performer()

		// create performer to test against
		const name = "TestPerformerUpdatePerformerImage"
		performer := models.Performer{
			Name:     sql.NullString{String: name, Valid: true},
			Checksum: utils.MD5FromString(name),
			Favorite: sql.NullBool{Bool: false, Valid: true},
		}
		created, err := qb.Create(performer)
		if err != nil {
			return fmt.Errorf("Error creating performer: %s", err.Error())
		}

		image := []byte("image")
		err = qb.UpdateImage(created.ID, image)
		if err != nil {
			return fmt.Errorf("Error updating performer image: %s", err.Error())
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

func TestPerformerDestroyPerformerImage(t *testing.T) {
	if err := withTxn(func(r models.Repository) error {
		qb := r.Performer()

		// create performer to test against
		const name = "TestPerformerDestroyPerformerImage"
		performer := models.Performer{
			Name:     sql.NullString{String: name, Valid: true},
			Checksum: utils.MD5FromString(name),
			Favorite: sql.NullBool{Bool: false, Valid: true},
		}
		created, err := qb.Create(performer)
		if err != nil {
			return fmt.Errorf("Error creating performer: %s", err.Error())
		}

		image := []byte("image")
		err = qb.UpdateImage(created.ID, image)
		if err != nil {
			return fmt.Errorf("Error updating performer image: %s", err.Error())
		}

		err = qb.DestroyImage(created.ID)
		if err != nil {
			return fmt.Errorf("Error destroying performer image: %s", err.Error())
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

func TestPerformerQueryAge(t *testing.T) {
	const age = 19
	ageCriterion := models.IntCriterionInput{
		Value:    age,
		Modifier: models.CriterionModifierEquals,
	}

	verifyPerformerAge(t, ageCriterion)

	ageCriterion.Modifier = models.CriterionModifierNotEquals
	verifyPerformerAge(t, ageCriterion)

	ageCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyPerformerAge(t, ageCriterion)

	ageCriterion.Modifier = models.CriterionModifierLessThan
	verifyPerformerAge(t, ageCriterion)
}

func verifyPerformerAge(t *testing.T, ageCriterion models.IntCriterionInput) {
	withTxn(func(r models.Repository) error {
		qb := r.Performer()
		performerFilter := models.PerformerFilterType{
			Age: &ageCriterion,
		}

		performers, _, err := qb.Query(&performerFilter, nil)
		if err != nil {
			t.Errorf("Error querying performer: %s", err.Error())
		}

		now := time.Now()
		for _, performer := range performers {
			bd := performer.Birthdate.String
			d, _ := time.Parse("2006-01-02", bd)
			age := now.Year() - d.Year()
			if now.YearDay() < d.YearDay() {
				age = age - 1
			}

			verifyInt(t, age, ageCriterion)
		}

		return nil
	})
}

func queryPerformers(t *testing.T, qb models.PerformerReader, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) []*models.Performer {
	performers, _, err := qb.Query(performerFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying performers: %s", err.Error())
	}

	return performers
}

func TestPerformerQueryTags(t *testing.T) {
	withTxn(func(r models.Repository) error {
		sqb := r.Performer()
		tagCriterion := models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdxWithPerformer]),
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludes,
		}

		performerFilter := models.PerformerFilterType{
			Tags: &tagCriterion,
		}

		// ensure ids are correct
		performers := queryPerformers(t, sqb, &performerFilter, nil)
		assert.Len(t, performers, 2)
		for _, performer := range performers {
			assert.True(t, performer.ID == performerIDs[performerIdxWithTag] || performer.ID == performerIDs[performerIdxWithTwoTags])
		}

		tagCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
				strconv.Itoa(tagIDs[tagIdx2WithPerformer]),
			},
			Modifier: models.CriterionModifierIncludesAll,
		}

		performers = queryPerformers(t, sqb, &performerFilter, nil)

		assert.Len(t, performers, 1)
		assert.Equal(t, sceneIDs[performerIdxWithTwoTags], performers[0].ID)

		tagCriterion = models.MultiCriterionInput{
			Value: []string{
				strconv.Itoa(tagIDs[tagIdx1WithPerformer]),
			},
			Modifier: models.CriterionModifierExcludes,
		}

		q := getSceneStringValue(performerIdxWithTwoTags, titleField)
		findFilter := models.FindFilterType{
			Q: &q,
		}

		performers = queryPerformers(t, sqb, &performerFilter, &findFilter)
		assert.Len(t, performers, 0)

		return nil
	})
}

func TestPerformerStashIDs(t *testing.T) {
	if err := withTxn(func(r models.Repository) error {
		qb := r.Performer()

		// create performer to test against
		const name = "TestStashIDs"
		performer := models.Performer{
			Name:     sql.NullString{String: name, Valid: true},
			Checksum: utils.MD5FromString(name),
			Favorite: sql.NullBool{Bool: false, Valid: true},
		}
		created, err := qb.Create(performer)
		if err != nil {
			return fmt.Errorf("Error creating performer: %s", err.Error())
		}

		testStashIDReaderWriter(t, qb, created.ID)
		return nil
	}); err != nil {
		t.Error(err.Error())
	}
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
