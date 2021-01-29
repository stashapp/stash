package models

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
	const minPerPage = 1
	const maxPerPage = 1000

	if ff.PerPage == nil {
		return defaultPerPage
	}

	if *ff.PerPage > 1000 {
		return maxPerPage
	} else if *ff.PerPage < 0 {
		// PerPage == 0 -> no limit
		return minPerPage
	}

	return *ff.PerPage
}

func (ff FindFilterType) IsGetAll() bool {
	return ff.PerPage != nil && *ff.PerPage == 0
}
