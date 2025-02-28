//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"slices"
	"strconv"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stretchr/testify/assert"
)

func TestMarkerFindBySceneID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.SceneMarker

		sceneID := sceneIDs[sceneIdxWithMarkers]
		markers, err := mqb.FindBySceneID(ctx, sceneID)

		if err != nil {
			t.Errorf("Error finding markers: %s", err.Error())
		}

		assert.Greater(t, len(markers), 0)
		for _, marker := range markers {
			assert.Equal(t, sceneIDs[sceneIdxWithMarkers], marker.SceneID)
		}

		markers, err = mqb.FindBySceneID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding marker: %s", err.Error())
		}

		assert.Len(t, markers, 0)

		return nil
	})
}

func TestMarkerCountByTagID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.SceneMarker

		markerCount, err := mqb.CountByTagID(ctx, tagIDs[tagIdxWithPrimaryMarkers])

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 6, markerCount)

		markerCount, err = mqb.CountByTagID(ctx, tagIDs[tagIdxWithMarkers])

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 2, markerCount)

		markerCount, err = mqb.CountByTagID(ctx, 0)

		if err != nil {
			t.Errorf("error calling CountByTagID: %s", err.Error())
		}

		assert.Equal(t, 0, markerCount)

		return nil
	})
}

func TestMarkerQueryQ(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		q := getSceneTitle(sceneIdxWithMarkers)
		m, _, err := db.SceneMarker.Query(ctx, nil, &models.FindFilterType{
			Q: &q,
		})

		if err != nil {
			t.Errorf("Error querying scene markers: %s", err.Error())
		}

		if !assert.Greater(t, len(m), 0) {
			return nil
		}

		assert.Equal(t, sceneIDs[sceneIdxWithMarkers], m[0].SceneID)

		return nil
	})
}

func TestMarkerQuerySortBySceneUpdated(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sort := "scenes_updated_at"
		_, _, err := db.SceneMarker.Query(ctx, nil, &models.FindFilterType{
			Sort: &sort,
		})

		if err != nil {
			t.Errorf("Error querying scene markers: %s", err.Error())
		}

		return nil
	})
}

func verifyIDs(t *testing.T, modifier models.CriterionModifier, values []int, results []int) {
	t.Helper()
	switch modifier {
	case models.CriterionModifierIsNull:
		assert.Len(t, results, 0)
	case models.CriterionModifierNotNull:
		assert.NotEqual(t, 0, len(results))
	case models.CriterionModifierIncludes:
		for _, v := range values {
			assert.Contains(t, results, v)
		}
	case models.CriterionModifierExcludes:
		for _, v := range values {
			assert.NotContains(t, results, v)
		}
	case models.CriterionModifierEquals:
		for _, v := range values {
			assert.Contains(t, results, v)
		}
		assert.Len(t, results, len(values))
	case models.CriterionModifierNotEquals:
		foundAll := true
		for _, v := range values {
			if !slices.Contains(results, v) {
				foundAll = false
				break
			}
		}
		if foundAll && len(results) == len(values) {
			t.Errorf("expected ids not equal to %v - found %v", values, results)
		}
	}
}

func TestMarkerQueryTags(t *testing.T) {
	type test struct {
		name         string
		markerFilter *models.SceneMarkerFilterType
		findFilter   *models.FindFilterType
	}

	withTxn(func(ctx context.Context) error {
		testTags := func(t *testing.T, m *models.SceneMarker, markerFilter *models.SceneMarkerFilterType) {
			tagIDs, err := db.SceneMarker.GetTagIDs(ctx, m.ID)
			if err != nil {
				t.Errorf("error getting marker tag ids: %v", err)
			}

			// HACK - if modifier isn't null/not null, then add the primary tag id
			if markerFilter.Tags.Modifier != models.CriterionModifierIsNull && markerFilter.Tags.Modifier != models.CriterionModifierNotNull {
				tagIDs = append(tagIDs, m.PrimaryTagID)
			}

			values, _ := stringslice.StringSliceToIntSlice(markerFilter.Tags.Value)
			verifyIDs(t, markerFilter.Tags.Modifier, values, tagIDs)
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
			{
				"includes",
				&models.SceneMarkerFilterType{
					Tags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierIncludes,
						Value: []string{
							strconv.Itoa(tagIDs[tagIdxWithMarkers]),
						},
					},
				},
				nil,
			},
			{
				"includes all",
				&models.SceneMarkerFilterType{
					Tags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierIncludesAll,
						Value: []string{
							strconv.Itoa(tagIDs[tagIdxWithMarkers]),
							strconv.Itoa(tagIDs[tagIdx2WithMarkers]),
						},
					},
				},
				nil,
			},
			{
				"equals",
				&models.SceneMarkerFilterType{
					Tags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierEquals,
						Value: []string{
							strconv.Itoa(tagIDs[tagIdxWithPrimaryMarkers]),
							strconv.Itoa(tagIDs[tagIdxWithMarkers]),
							strconv.Itoa(tagIDs[tagIdx2WithMarkers]),
						},
					},
				},
				nil,
			},
			// not equals not supported
			// {
			// 	"not equals",
			// 	&models.SceneMarkerFilterType{
			// 		Tags: &models.HierarchicalMultiCriterionInput{
			// 			Modifier: models.CriterionModifierNotEquals,
			// 			Value: []string{
			// 				strconv.Itoa(tagIDs[tagIdx2WithScene]),
			// 				strconv.Itoa(tagIDs[tagIdx3WithScene]),
			// 			},
			// 		},
			// 	},
			// 	nil,
			// },
			{
				"excludes",
				&models.SceneMarkerFilterType{
					Tags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierIncludes,
						Value: []string{
							strconv.Itoa(tagIDs[tagIdx2WithMarkers]),
						},
					},
				},
				nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				markers := queryMarkers(ctx, t, db.SceneMarker, tc.markerFilter, tc.findFilter)
				assert.Greater(t, len(markers), 0)
				for _, m := range markers {
					testTags(t, m, tc.markerFilter)
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

	withTxn(func(ctx context.Context) error {
		testTags := func(t *testing.T, m *models.SceneMarker, markerFilter *models.SceneMarkerFilterType) {
			s, err := db.Scene.Find(ctx, m.SceneID)
			if err != nil {
				t.Errorf("error getting marker tag ids: %v", err)
				return
			}

			if err := s.LoadTagIDs(ctx, db.Scene); err != nil {
				t.Errorf("error getting marker tag ids: %v", err)
				return
			}

			tagIDs := s.TagIDs.List()
			values, _ := stringslice.StringSliceToIntSlice(markerFilter.SceneTags.Value)
			verifyIDs(t, markerFilter.SceneTags.Modifier, values, tagIDs)
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
			{
				"includes",
				&models.SceneMarkerFilterType{
					SceneTags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierIncludes,
						Value: []string{
							strconv.Itoa(tagIDs[tagIdx3WithScene]),
						},
					},
				},
				nil,
			},
			{
				"includes all",
				&models.SceneMarkerFilterType{
					SceneTags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierIncludesAll,
						Value: []string{
							strconv.Itoa(tagIDs[tagIdx2WithScene]),
							strconv.Itoa(tagIDs[tagIdx3WithScene]),
						},
					},
				},
				nil,
			},
			{
				"equals",
				&models.SceneMarkerFilterType{
					SceneTags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierEquals,
						Value: []string{
							strconv.Itoa(tagIDs[tagIdx2WithScene]),
							strconv.Itoa(tagIDs[tagIdx3WithScene]),
						},
					},
				},
				nil,
			},
			// not equals not supported
			// {
			// 	"not equals",
			// 	&models.SceneMarkerFilterType{
			// 		SceneTags: &models.HierarchicalMultiCriterionInput{
			// 			Modifier: models.CriterionModifierNotEquals,
			// 			Value: []string{
			// 				strconv.Itoa(tagIDs[tagIdx2WithScene]),
			// 				strconv.Itoa(tagIDs[tagIdx3WithScene]),
			// 			},
			// 		},
			// 	},
			// 	nil,
			// },
			{
				"excludes",
				&models.SceneMarkerFilterType{
					SceneTags: &models.HierarchicalMultiCriterionInput{
						Modifier: models.CriterionModifierIncludes,
						Value: []string{
							strconv.Itoa(tagIDs[tagIdx2WithScene]),
						},
					},
				},
				nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				markers := queryMarkers(ctx, t, db.SceneMarker, tc.markerFilter, tc.findFilter)
				assert.Greater(t, len(markers), 0)
				for _, m := range markers {
					testTags(t, m, tc.markerFilter)
				}
			})
		}

		return nil
	})
}

func markersToIDs(i []*models.SceneMarker) []int {
	ret := make([]int, len(i))
	for i, v := range i {
		ret[i] = v.ID
	}

	return ret
}

func TestMarkerQueryDuration(t *testing.T) {
	type test struct {
		name         string
		markerFilter *models.SceneMarkerFilterType
		include      []int
		exclude      []int
	}

	cases := []test{
		{
			"is null",
			&models.SceneMarkerFilterType{
				Duration: &models.FloatCriterionInput{
					Modifier: models.CriterionModifierIsNull,
				},
			},
			[]int{markerIdxWithScene},
			[]int{markerIdxWithDuration},
		},
		{
			"not null",
			&models.SceneMarkerFilterType{
				Duration: &models.FloatCriterionInput{
					Modifier: models.CriterionModifierNotNull,
				},
			},
			[]int{markerIdxWithDuration},
			[]int{markerIdxWithScene},
		},
		{
			"equals",
			&models.SceneMarkerFilterType{
				Duration: &models.FloatCriterionInput{
					Modifier: models.CriterionModifierEquals,
					Value:    markerIdxWithDuration,
				},
			},
			[]int{markerIdxWithDuration},
			[]int{markerIdx2WithDuration, markerIdxWithScene},
		},
		{
			"not equals",
			&models.SceneMarkerFilterType{
				Duration: &models.FloatCriterionInput{
					Modifier: models.CriterionModifierNotEquals,
					Value:    markerIdx2WithDuration,
				},
			},
			[]int{markerIdxWithDuration},
			[]int{markerIdx2WithDuration, markerIdxWithScene},
		},
		{
			"greater than",
			&models.SceneMarkerFilterType{
				Duration: &models.FloatCriterionInput{
					Modifier: models.CriterionModifierGreaterThan,
					Value:    markerIdxWithDuration,
				},
			},
			[]int{markerIdx2WithDuration},
			[]int{markerIdxWithDuration, markerIdxWithScene},
		},
		{
			"less than",
			&models.SceneMarkerFilterType{
				Duration: &models.FloatCriterionInput{
					Modifier: models.CriterionModifierLessThan,
					Value:    markerIdx2WithDuration,
				},
			},
			[]int{markerIdxWithDuration},
			[]int{markerIdx2WithDuration, markerIdxWithScene},
		},
	}

	qb := db.SceneMarker

	for _, tt := range cases {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, _, err := qb.Query(ctx, tt.markerFilter, nil)
			if err != nil {
				t.Errorf("SceneMarkerStore.Query() error = %v", err)
				return
			}

			ids := markersToIDs(got)
			include := indexesToIDs(markerIDs, tt.include)
			exclude := indexesToIDs(markerIDs, tt.exclude)

			for _, i := range include {
				assert.Contains(ids, i)
			}
			for _, e := range exclude {
				assert.NotContains(ids, e)
			}
		})
	}

}

func queryMarkers(ctx context.Context, t *testing.T, sqb models.SceneMarkerReader, markerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) []*models.SceneMarker {
	t.Helper()
	result, _, err := sqb.Query(ctx, markerFilter, findFilter)
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
// TODO Count
// TODO All
// TODO Query
