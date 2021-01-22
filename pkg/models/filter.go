package models

type QueryFilter struct {
	FindFilterType
	All bool
}

func NewQueryFilter(f *FindFilterType) QueryFilter {
	if f == nil {
		f = &FindFilterType{}
	}

	ret := QueryFilter{
		FindFilterType: *f,
	}

	ret.SetDefaults()
	return ret
}

func (f *QueryFilter) SetDefaults() {
	// if we're getting all, we don't need to set page parameters
	if f.All {
		return
	}

	defaultPage := 1
	defaultPerPage := 25
	minPerPage := 1
	maxPerPage := 1000

	if f.Page == nil || *f.Page < 1 {
		f.Page = &defaultPage
	}

	if f.PerPage == nil {
		f.PerPage = &defaultPerPage
	}

	if *f.PerPage > 1000 {
		f.PerPage = &maxPerPage
	} else if *f.PerPage < 1 {
		f.PerPage = &minPerPage
	}
}
