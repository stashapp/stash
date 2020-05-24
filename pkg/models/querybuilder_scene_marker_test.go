// +build integration

package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func TestMarkerFindBySceneID(t *testing.T) {
	mqb := models.NewSceneMarkerQueryBuilder()

	sceneID := sceneIDs[sceneIdxWithMarker]
	markers, err := mqb.FindBySceneID(sceneID, nil)

	if err != nil {
		t.Fatalf("Error finding markers: %s", err.Error())
	}

	assert.Len(t, markers, 1)
	assert.Equal(t, markerIDs[markerIdxWithScene], markers[0].ID)

	markers, err = mqb.FindBySceneID(0, nil)

	if err != nil {
		t.Fatalf("Error finding marker: %s", err.Error())
	}

	assert.Len(t, markers, 0)
}

func TestMarkerCountByTagID(t *testing.T) {
	mqb := models.NewSceneMarkerQueryBuilder()

	markerCount, err := mqb.CountByTagID(tagIDs[tagIdxWithPrimaryMarker])

	if err != nil {
		t.Fatalf("error calling CountByTagID: %s", err.Error())
	}

	assert.Equal(t, 1, markerCount)

	markerCount, err = mqb.CountByTagID(tagIDs[tagIdxWithMarker])

	if err != nil {
		t.Fatalf("error calling CountByTagID: %s", err.Error())
	}

	assert.Equal(t, 1, markerCount)

	markerCount, err = mqb.CountByTagID(0)

	if err != nil {
		t.Fatalf("error calling CountByTagID: %s", err.Error())
	}

	assert.Equal(t, 0, markerCount)
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO GetMarkerStrings
// TODO Wall
// TODO Query
