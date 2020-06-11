// +build integration

package models_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
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

// TODO Update
// TODO Destroy
// TODO Find
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
