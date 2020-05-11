// +build integration

package models_test

import (
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

// TODO Update
// TODO Destroy
// TODO Find
// TODO FindByNames
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
