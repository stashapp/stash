//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stretchr/testify/assert"
)

func TestFindBySceneID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.SceneFilter

		sceneID := sceneIDs[sceneIdxWithFilters]
		filters, err := mqb.FindBySceneID(ctx, sceneID)

		if err != nil {
			t.Errorf("Error finding filters: %s", err.Error())
		}

		assert.Greater(t, len(filters), 0)
		for _, filter := range filters {
			assert.Equal(t, sceneIDs[sceneIdxWithFilters], filter.SceneID)
		}

		filters, err = mqb.FindBySceneID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding filter: %s", err.Error())
		}

		assert.Len(t, filters, 0)

		return nil
	})
}

func TestFilterFindBySceneID(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		mqb := db.SceneFilter

		sceneID := sceneIDs[sceneIdxWithFilters]
		filters, err := mqb.FindBySceneID(ctx, sceneID)

		if err != nil {
			t.Errorf("Error finding filters: %s", err.Error())
		}

		assert.Greater(t, len(filters), 0)
		for _, filter := range filters {
			assert.Equal(t, sceneIDs[sceneIdxWithFilters], filter.SceneID)
		}

		filters, err = mqb.FindBySceneID(ctx, 0)

		if err != nil {
			t.Errorf("Error finding filter: %s", err.Error())
		}

		assert.Len(t, filters, 0)

		return nil
	})
}

func TestFilterQuerySortBySceneUpdated(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		sort := "scenes_updated_at"
		_, _, err := db.SceneFilter.Query(ctx, nil, &models.FindFilterType{
			Sort: &sort,
		})

		if err != nil {
			t.Errorf("Error querying scene filters: %s", err.Error())
		}

		return nil
	})
}

func verifyFilterIDs(t *testing.T, modifier models.CriterionModifier, values []int, results []int) {
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
			if !intslice.IntInclude(results, v) {
				foundAll = false
				break
			}
		}
		if foundAll && len(results) == len(values) {
			t.Errorf("expected ids not equal to %v - found %v", values, results)
		}
	}
}

func queryFilters(ctx context.Context, t *testing.T, sqb models.SceneFilterReader, filterFilter *models.SceneFilterFilterType, findFilter *models.FindFilterType) []*models.SceneFilter {
	t.Helper()
	result, _, err := sqb.Query(ctx, filterFilter, findFilter)
	if err != nil {
		t.Errorf("Error querying filters: %v", err)
	}

	return result
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO Wall
// TODO Count
// TODO All
// TODO Query
