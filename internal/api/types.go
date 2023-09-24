package api

import (
	"math"

	"github.com/stashapp/stash/pkg/models"
)

// #1572 - Inf and NaN values cause the JSON marshaller to fail
// Return nil for these values
func handleFloat64(v float64) *float64 {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return nil
	}

	return &v
}

func stashIDsSliceToPtrSlice(v []models.StashID) []*models.StashID {
	ret := make([]*models.StashID, len(v))
	for i, vv := range v {
		c := vv
		ret[i] = &c
	}

	return ret
}
