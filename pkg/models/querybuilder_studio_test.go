// +build integration

package models_test

import (
	"strconv"
	"strings"
	"testing"

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

// TODO Create
// TODO Update
// TODO Destroy
// TODO Find
// TODO FindBySceneID
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
