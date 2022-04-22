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

func handleFloat64Ptr(v *float64) *float64 {
	if v == nil {
		return nil
	}
	return handleFloat64(*v)
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

func translateSceneMovieIDs(input BulkUpdateIds) (*models.UpdateMovieIDs, error) {
	ids, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids [%v]: %w", input.Ids, err)
	}

	ret := &models.UpdateMovieIDs{
		Mode: input.Mode,
	}

	for _, id := range ids {
		ret.Movies = append(ret.Movies, models.MoviesScenes{
			MovieID: id,
		})
	}

	return ret, nil
}

func stringPtrToString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func float64PtrToFloat64(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
