package models

import (
	"fmt"
	"io"
	"strconv"
)

// PerPageAll is the value used for perPage to indicate all results should be
// returned.
const PerPageAll = -1

type SortDirectionEnum string

const (
	SortDirectionEnumAsc  SortDirectionEnum = "ASC"
	SortDirectionEnumDesc SortDirectionEnum = "DESC"
)

var AllSortDirectionEnum = []SortDirectionEnum{
	SortDirectionEnumAsc,
	SortDirectionEnumDesc,
}

func (e SortDirectionEnum) IsValid() bool {
	switch e {
	case SortDirectionEnumAsc, SortDirectionEnumDesc:
		return true
	}
	return false
}

func (e SortDirectionEnum) String() string {
	return string(e)
}

func (e *SortDirectionEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortDirectionEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortDirectionEnum", str)
	}
	return nil
}

func (e SortDirectionEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FindFilterType struct {
	Q    *string `json:"q"`
	Page *int    `json:"page"`
	// use per_page = -1 to indicate all results. Defaults to 25.
	PerPage   *int               `json:"per_page"`
	Sort      *string            `json:"sort"`
	Direction *SortDirectionEnum `json:"direction"`
}

func (ff FindFilterType) GetSort(defaultSort string) string {
	var sort string
	if ff.Sort == nil {
		sort = defaultSort
	} else {
		sort = *ff.Sort
	}
	return sort
}

func (ff FindFilterType) GetDirection() string {
	var direction string
	if directionFilter := ff.Direction; directionFilter != nil {
		if dir := directionFilter.String(); directionFilter.IsValid() {
			direction = dir
		} else {
			direction = "ASC"
		}
	} else {
		direction = "ASC"
	}
	return direction
}

func (ff FindFilterType) GetPage() int {
	const defaultPage = 1
	if ff.Page == nil || *ff.Page < 1 {
		return defaultPage
	}

	return *ff.Page
}

func (ff FindFilterType) GetPageSize() int {
	const defaultPerPage = 25
	const minPerPage = 0

	if ff.PerPage == nil {
		return defaultPerPage
	}

	// removed the maxPerPage check. We already all -1 to indicate all results
	// so there is no conceivable reason we should limit the page size

	if *ff.PerPage < minPerPage {
		// negative page sizes should return all results
		// this is a sanity check in case GetPageSize is
		// called with a negative page size.
		return minPerPage
	}

	return *ff.PerPage
}

func (ff FindFilterType) IsGetAll() bool {
	return ff.PerPage != nil && *ff.PerPage < 0
}

// BatchFindFilter returns a FindFilterType suitable for batch finding
// using the provided batch size.
func BatchFindFilter(batchSize int) *FindFilterType {
	page := 1
	return &FindFilterType{
		PerPage: &batchSize,
		Page:    &page,
	}
}
