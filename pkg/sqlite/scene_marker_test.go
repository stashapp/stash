//go:build integration
// +build integration

package sqlite_test

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestMarkerFindBySceneID(t *testing.T) {
	withTxn(func(r models.Repository) error {
		mqb := r.SceneMarker()

		sceneID := sceneIDs[sceneIdxWithMarker]
		markers, err := mqb.FindBySceneID(sceneID)

		if err != nil {
			t.Errorf("Error finding markers: %s", err.Error())
		}

		assert.Len(t, markers, 1)
		assert.Equal(t, markerIDs[markerIdxWithScene], markers[0].ID)

		markers, err = mqb.FindBySceneID(0)

		if err != nil {
			t.Errorf("Error finding marker: %s", err.Error())
		}

		assert.Len(t, markers, 0)

		return nil
	})
}

func TestMarkerCountByTagID(t *testing.T) {
	withTxn(func(r models.Repository) error {
		mqb := r.SceneMarker()

		markerCount, err := mqb.CountByTagID(tagIDs[tagIdxWithPrimaryMarker])

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 1, markerCount)

		markerCount, err = mqb.CountByTagID(tagIDs[tagIdxWithMarker])

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 1, markerCount)

		markerCount, err = mqb.CountByTagID(0)

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 0, markerCount)

		return nil
	})
}

func TestMarkerQuerySortBySceneUpdated(t *testing.T) {
	withTxn(func(r models.Repository) error {
		sort := "scenes_updated_at"
		_, _, err := r.SceneMarker().Query(nil, &models.FindFilterType{
			Sort: &sort,
		})

		if err != nil {
			t.Errorf("Error querying scene markers: %s", err.Error())
		}

		return nil
	})
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO GetMarkerStrings
// TODO Wall
// TODO Query
