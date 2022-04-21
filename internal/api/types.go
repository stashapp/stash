package api

import (
	"fmt"
	"math"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

// An enum https://golang.org/ref/spec#Iota
const (
	create = iota // 0
	update = iota // 1
)

// #1572 - Inf and NaN values cause the JSON marshaller to fail
// Return nil for these values
func handleFloat64(v float64) *float64 {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return nil
	}

	return &v
}

func translateUpdateIDs(strIDs []string, mode models.RelationshipUpdateMode) (*models.UpdateIDs, error) {
	ids, err := stringslice.StringSliceToIntSlice(strIDs)
	if err != nil {
		return nil, fmt.Errorf("converting ids [%v]: %w", strIDs, err)
	}
	return &models.UpdateIDs{
		IDs:  ids,
		Mode: mode,
	}, nil
}
