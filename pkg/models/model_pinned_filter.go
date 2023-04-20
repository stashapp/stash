package models

type PinnedFilter struct {
	ID   int        `db:"id" json:"id"`
	Mode FilterMode `db:"mode" json:"mode"`
	Name string     `db:"name" json:"name"`
}

type PinnedFilters []*PinnedFilter

func (m *PinnedFilters) Append(o interface{}) {
	*m = append(*m, o.(*PinnedFilter))
}

func (m *PinnedFilters) New() interface{} {
	return &PinnedFilter{}
}
