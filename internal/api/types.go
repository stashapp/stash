package api

import (
	"github.com/stashapp/stash/pkg/models"
)

func stashIDsSliceToPtrSlice(v []models.StashID) []*models.StashID {
	ret := make([]*models.StashID, len(v))
	for i, vv := range v {
		c := vv
		ret[i] = &c
	}

	return ret
}
