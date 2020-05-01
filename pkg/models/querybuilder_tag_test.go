// +build integration
package models_test

import (
	"testing"

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

// TODO Create
// TODO Update
// TODO Destroy
// TODO Find
// TODO FindBySceneID
// TODO FindBySceneMarkerID
// TODO FindByName
// TODO FindByNames
// TODO Count
// TODO All
// TODO AllSlim
// TODO Query
