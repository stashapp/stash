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
