package api

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

func stashIDsSliceToPtrSlice(v []models.StashID) []*models.StashID {
	return sliceutil.ValuesToPtrs(v)
}
