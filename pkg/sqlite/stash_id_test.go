//go:build integration
// +build integration

package sqlite_test

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

type stashIDReaderWriter interface {
	GetStashIDs(performerID int) ([]*models.StashID, error)
	UpdateStashIDs(performerID int, stashIDs []models.StashID) error
}

func testStashIDReaderWriter(t *testing.T, r stashIDReaderWriter, id int) {
	// ensure no stash IDs to begin with
	testNoStashIDs(t, r, id)

	// ensure GetStashIDs with non-existing also returns none
	testNoStashIDs(t, r, -1)

	// add stash ids
	const stashIDStr = "stashID"
	const endpoint = "endpoint"
	stashID := models.StashID{
		StashID:  stashIDStr,
		Endpoint: endpoint,
	}

	// update stash ids and ensure was updated
	if err := r.UpdateStashIDs(id, []models.StashID{stashID}); err != nil {
		t.Error(err.Error())
	}

	testStashIDs(t, r, id, []*models.StashID{&stashID})

	// update non-existing id - should return error
	if err := r.UpdateStashIDs(-1, []models.StashID{stashID}); err == nil {
		t.Error("expected error when updating non-existing id")
	}

	// remove stash ids and ensure was updated
	if err := r.UpdateStashIDs(id, []models.StashID{}); err != nil {
		t.Error(err.Error())
	}

	testNoStashIDs(t, r, id)
}

func testNoStashIDs(t *testing.T, r stashIDReaderWriter, id int) {
	t.Helper()
	stashIDs, err := r.GetStashIDs(id)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Len(t, stashIDs, 0)
}

func testStashIDs(t *testing.T, r stashIDReaderWriter, id int, expected []*models.StashID) {
	t.Helper()
	stashIDs, err := r.GetStashIDs(id)
	if err != nil {
		t.Error(err.Error())
		return
	}

	assert.Equal(t, stashIDs, expected)
}
