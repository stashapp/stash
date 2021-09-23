//go:build integration
// +build integration

package sqlite_test

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSavedFilterFind(t *testing.T) {
	withTxn(func(r models.Repository) error {
		savedFilter, err := r.SavedFilter().Find(savedFilterIDs[savedFilterIdxImage])

		if err != nil {
			t.Errorf("Error finding saved filter: %s", err.Error())
		}

		assert.Equal(t, savedFilterIDs[savedFilterIdxImage], savedFilter.ID)

		return nil
	})
}

func TestSavedFilterFindByMode(t *testing.T) {
	withTxn(func(r models.Repository) error {
		savedFilters, err := r.SavedFilter().FindByMode(models.FilterModeScenes)

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
	withTxn(func(r models.Repository) error {
		created, err := r.SavedFilter().Create(models.SavedFilter{
			Name:   filterName,
			Mode:   models.FilterModeScenes,
			Filter: testFilter,
		})

		if err == nil {
			id = created.ID
		}

		return err
	})

	withTxn(func(r models.Repository) error {
		qb := r.SavedFilter()

		return qb.Destroy(id)
	})

	// now try to find it
	withTxn(func(r models.Repository) error {
		found, err := r.SavedFilter().Find(id)
		if err == nil {
			assert.Nil(t, found)
		}

		return err
	})
}

func TestSavedFilterFindDefault(t *testing.T) {
	withTxn(func(r models.Repository) error {
		def, err := r.SavedFilter().FindDefault(models.FilterModeScenes)
		if err == nil {
			assert.Equal(t, savedFilterIDs[savedFilterIdxDefaultScene], def.ID)
		}

		return err
	})
}

func TestSavedFilterSetDefault(t *testing.T) {
	const newFilter = "foo"

	withTxn(func(r models.Repository) error {
		_, err := r.SavedFilter().SetDefault(models.SavedFilter{
			Mode:   models.FilterModeMovies,
			Filter: newFilter,
		})

		return err
	})

	var defID int
	withTxn(func(r models.Repository) error {
		def, err := r.SavedFilter().FindDefault(models.FilterModeMovies)
		if err == nil {
			defID = def.ID
			assert.Equal(t, newFilter, def.Filter)
		}

		return err
	})

	// destroy it again
	withTxn(func(r models.Repository) error {
		return r.SavedFilter().Destroy(defID)
	})
}

// TODO Update
// TODO Destroy
// TODO Find
// TODO GetMarkerStrings
// TODO Wall
// TODO Query
