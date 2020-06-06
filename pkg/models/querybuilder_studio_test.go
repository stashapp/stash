// +build integration

package models_test

import (
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

// TODO Create
// TODO Update
// TODO Destroy
// TODO Find
// TODO FindBySceneID
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
