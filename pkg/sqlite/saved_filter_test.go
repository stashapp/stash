//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestSavedFilterFind(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		savedFilter, err := sqlite.SavedFilterReaderWriter.Find(ctx, savedFilterIDs[savedFilterIdxImage])

		if err != nil {
			t.Errorf("Error finding saved filter: %s", err.Error())
		}

		assert.Equal(t, savedFilterIDs[savedFilterIdxImage], savedFilter.ID)

		return nil
	})
}

func TestSavedFilterFindByMode(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		savedFilters, err := sqlite.SavedFilterReaderWriter.FindByMode(ctx, models.FilterModeScenes)

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
	const testFilter = "{}"
	var id int

	// create the saved filter to destroy
	withTxn(func(ctx context.Context) error {
		created, err := sqlite.SavedFilterReaderWriter.Create(ctx, models.SavedFilter{
			Name:   filterName,
			Mode:   models.FilterModeScenes,
			Filter: testFilter,
		})

		if err == nil {
			id = created.ID
		}

		return err
	})

	withTxn(func(ctx context.Context) error {
		qb := sqlite.SavedFilterReaderWriter

		return qb.Destroy(ctx, id)
	})

	// now try to find it
	withTxn(func(ctx context.Context) error {
		found, err := sqlite.SavedFilterReaderWriter.Find(ctx, id)
		if err == nil {
			assert.Nil(t, found)
		}

		return err
	})
}

func TestSavedFilterFindDefault(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		def, err := sqlite.SavedFilterReaderWriter.FindDefault(ctx, models.FilterModeScenes)
		if err == nil {
			assert.Equal(t, savedFilterIDs[savedFilterIdxDefaultScene], def.ID)
		}

		return err
	})
}

func TestSavedFilterSetDefault(t *testing.T) {
	const newFilter = "foo"

	withTxn(func(ctx context.Context) error {
		_, err := sqlite.SavedFilterReaderWriter.SetDefault(ctx, models.SavedFilter{
			Mode:   models.FilterModeMovies,
			Filter: newFilter,
		})

		return err
	})

	var defID int
	withTxn(func(ctx context.Context) error {
		def, err := sqlite.SavedFilterReaderWriter.FindDefault(ctx, models.FilterModeMovies)
		if err == nil {
			defID = def.ID
			assert.Equal(t, newFilter, def.Filter)
		}

		return err
	})

	// destroy it again
	withTxn(func(ctx context.Context) error {
		return sqlite.SavedFilterReaderWriter.Destroy(ctx, defID)
	})
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO GetMarkerStrings
// TODO Wall
// TODO Query
