//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSavedFilterFind(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		savedFilter, err := db.SavedFilter.Find(ctx, savedFilterIDs[savedFilterIdxImage])

		if err != nil {
			t.Errorf("Error finding saved filter: %s", err.Error())
		}

		assert.Equal(t, savedFilterIDs[savedFilterIdxImage], savedFilter.ID)

		return nil
	})
}

func TestSavedFilterFindByMode(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		savedFilters, err := db.SavedFilter.FindByMode(ctx, models.FilterModeScenes)

		if err != nil {
			t.Errorf("Error finding saved filters: %s", err.Error())
		}

		assert.Len(t, savedFilters, 1)
		assert.Equal(t, savedFilterIDs[savedFilterIdxScene], savedFilters[0].ID)

		return nil
	})
}

func TestSavedFilterDestroy(t *testing.T) {
	const filterName = "filterToDestroy"
	filterQ := ""
	filterPage := 1
	filterPerPage := 40
	filterSort := "date"
	filterDirection := models.SortDirectionEnumAsc
	findFilter := models.FindFilterType{
		Q:         &filterQ,
		Page:      &filterPage,
		PerPage:   &filterPerPage,
		Sort:      &filterSort,
		Direction: &filterDirection,
	}
	objectFilter := map[string]interface{}{
		"test": "foo",
	}
	uiOptions := map[string]interface{}{
		"display_mode": 1,
		"zoom_index":   1,
	}
	var id int

	// create the saved filter to destroy
	withTxn(func(ctx context.Context) error {
		newFilter := models.SavedFilter{
			Name:         filterName,
			Mode:         models.FilterModeScenes,
			FindFilter:   &findFilter,
			ObjectFilter: objectFilter,
			UIOptions:    uiOptions,
		}
		err := db.SavedFilter.Create(ctx, &newFilter)

		if err == nil {
			id = newFilter.ID
		}

		return err
	})

	withTxn(func(ctx context.Context) error {
		return db.SavedFilter.Destroy(ctx, id)
	})

	// now try to find it
	withTxn(func(ctx context.Context) error {
		found, err := db.SavedFilter.Find(ctx, id)
		if err == nil {
			assert.Nil(t, found)
		}

		return err
	})
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO GetMarkerStrings
// TODO Wall
// TODO Query

func TestSavedFilterPerUser(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		qb := db.SavedFilter

		other := models.User{
			Username: "other",
			Roles:    models.Roles{models.RoleEnumModify},
		}
		require.NoError(t, db.User.Create(ctx, &other, ""))
		otherCtx := models.WithUserID(ctx, other.ID)
		guestCtx := models.WithUserID(ctx, 0)

		const name = "per user filter"
		newFilter := models.SavedFilter{
			Name: name,
			Mode: models.FilterModeScenes,
		}
		require.NoError(t, qb.Create(ctx, &newFilter))

		// other users may use the same name
		otherFilter := models.SavedFilter{
			Name: name,
			Mode: models.FilterModeScenes,
		}
		require.NoError(t, qb.Create(otherCtx, &otherFilter))

		// the guest user cannot create filters
		assert.Error(t, qb.Create(guestCtx, &models.SavedFilter{Name: name, Mode: models.FilterModeScenes}))

		// users only see their own filters
		f, err := qb.Find(otherCtx, newFilter.ID)
		require.NoError(t, err)
		assert.Nil(t, f)

		f, err = qb.Find(guestCtx, newFilter.ID)
		require.NoError(t, err)
		assert.Nil(t, f)

		all, err := qb.All(otherCtx)
		require.NoError(t, err)
		require.Len(t, all, 1)
		assert.Equal(t, otherFilter.ID, all[0].ID)

		byMode, err := qb.FindByMode(guestCtx, models.FilterModeScenes)
		require.NoError(t, err)
		assert.Empty(t, byMode)

		// users cannot change or remove other users' filters
		assert.Error(t, qb.Update(otherCtx, &models.SavedFilter{ID: newFilter.ID, Name: "changed", Mode: models.FilterModeScenes}))
		assert.Error(t, qb.Destroy(otherCtx, newFilter.ID))

		f, err = qb.Find(ctx, newFilter.ID)
		require.NoError(t, err)
		require.NotNil(t, f)
		assert.Equal(t, name, f.Name)

		return nil
	})
}
