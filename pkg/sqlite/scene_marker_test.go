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

		sceneID := sceneIDs[sceneIdxWithMarkers]
		markers, err := mqb.FindBySceneID(sceneID)

		if err != nil {
			t.Errorf("Error finding markers: %s", err.Error())
		}

		assert.Greater(t, len(markers), 0)
		for _, marker := range markers {
			assert.Equal(t, sceneIDs[sceneIdxWithMarkers], int(marker.SceneID.Int64))
		}

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

		markerCount, err := mqb.CountByTagID(tagIDs[tagIdxWithPrimaryMarkers])

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 3, markerCount)

		markerCount, err = mqb.CountByTagID(tagIDs[tagIdxWithMarkers])

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

func TestMarkerQueryTags(t *testing.T) {
	type test struct {
		name         string
		markerFilter *models.SceneMarkerFilterType
		findFilter   *models.FindFilterType
	}

	withTxn(func(r models.Repository) error {
		testTags := func(m *models.SceneMarker, markerFilter *models.SceneMarkerFilterType) {
			tagIDs, err := r.SceneMarker().GetTagIDs(m.ID)
			if err != nil {
				t.Errorf("error getting marker tag ids: %v", err)
			}
			if markerFilter.Tags.Modifier == models.CriterionModifierIsNull && len(tagIDs) > 0 {
				t.Errorf("expected marker %d to have no tags - found %d", m.ID, len(tagIDs))
			}
			if markerFilter.Tags.Modifier == models.CriterionModifierNotNull && len(tagIDs) == 0 {
				t.Errorf("expected marker %d to have tags - found 0", m.ID)
			}
		}

		cases := []test{
			{
				"is null",
				&models.SceneMarkerFilterType{
					Tags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierIsNull,
					},
				},
				nil,
			},
			{
				"not null",
				&models.SceneMarkerFilterType{
					Tags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierNotNull,
					},
				},
				nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				markers := queryMarkers(t, r.SceneMarker(), tc.markerFilter, tc.findFilter)
				assert.Greater(t, len(markers), 0)
				for _, m := range markers {
					testTags(m, tc.markerFilter)
				}
			})
		}

		return nil
	})
}

func TestMarkerQuerySceneTags(t *testing.T) {
	type test struct {
		name         string
		markerFilter *models.SceneMarkerFilterType
		findFilter   *models.FindFilterType
	}

	withTxn(func(r models.Repository) error {
		testTags := func(m *models.SceneMarker, markerFilter *models.SceneMarkerFilterType) {
			tagIDs, err := r.Scene().GetTagIDs(int(m.SceneID.Int64))
			if err != nil {
				t.Errorf("error getting marker tag ids: %v", err)
			}
			if markerFilter.SceneTags.Modifier == models.CriterionModifierIsNull && len(tagIDs) > 0 {
				t.Errorf("expected marker %d to have no scene tags - found %d", m.ID, len(tagIDs))
			}
			if markerFilter.SceneTags.Modifier == models.CriterionModifierNotNull && len(tagIDs) == 0 {
				t.Errorf("expected marker %d to have scene tags - found 0", m.ID)
			}
		}

		cases := []test{
			{
				"is null",
				&models.SceneMarkerFilterType{
					SceneTags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierIsNull,
					},
				},
				nil,
			},
			{
				"not null",
				&models.SceneMarkerFilterType{
					SceneTags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierNotNull,
					},
				},
				nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				markers := queryMarkers(t, r.SceneMarker(), tc.markerFilter, tc.findFilter)
				assert.Greater(t, len(markers), 0)
				for _, m := range markers {
					testTags(m, tc.markerFilter)
				}
			})
		}

		return nil
	})
}

func queryMarkers(t *testing.T, sqb models.SceneMarkerReader, markerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) []*models.SceneMarker {
	t.Helper()
	result, _, err := sqb.Query(markerFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying markers: %v", err)
	}

	return result
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO GetMarkerStrings
// TODO Wall
// TODO Query
