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

func TestPinnedFilterFind(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		pinnedFilter, err := sqlite.PinnedFilterReaderWriter.Find(ctx, pinnedFilterIDs[pinnedFilterIdxImage])

		if err != nil {
			t.Errorf("Error finding pinned filter: %s", err.Error())
		}

		assert.Equal(t, pinnedFilterIDs[pinnedFilterIdxImage], pinnedFilter.ID)

		return nil
	})
}

func TestPinnedFilterFindByMode(t *testing.T) {
	withTxn(func(ctx context.Context) error {
		pinnedFilters, err := sqlite.PinnedFilterReaderWriter.FindByMode(ctx, models.FilterModeScenes)

		if err != nil {
			t.Errorf("Error finding pinned filters: %s", err.Error())
		}

		assert.Len(t, pinnedFilters, 1)
		assert.Equal(t, pinnedFilterIDs[pinnedFilterIdxScene], pinnedFilters[0].ID)

		return nil
	})
}

func TestPinnedFilterDestroy(t *testing.T) {
	const filterName = "filterToDestroy"
	var id int

	// create the pinned filter to destroy
	withTxn(func(ctx context.Context) error {
		created, err := sqlite.PinnedFilterReaderWriter.Create(ctx, models.PinnedFilter{
			Name: filterName,
			Mode: models.FilterModeScenes,
		})

		if err == nil {
			id = created.ID
		}

		return err
	})

	withTxn(func(ctx context.Context) error {
		qb := sqlite.PinnedFilterReaderWriter

		return qb.Destroy(ctx, id)
	})

	// now try to find it
	withTxn(func(ctx context.Context) error {
		found, err := sqlite.PinnedFilterReaderWriter.Find(ctx, id)
		if err == nil {
			assert.Nil(t, found)
		}

		return err
	})
}
