package models

// PerPageAll is the value used for perPage to indicate all results should be
// returned.
const PerPageAll = -1

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
	const maxPerPage = 1000

	if ff.PerPage == nil {
		return defaultPerPage
	}

	if *ff.PerPage > maxPerPage {
		return maxPerPage
	} else if *ff.PerPage < minPerPage {
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
