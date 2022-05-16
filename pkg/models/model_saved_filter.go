package models

type SavedFilter struct {
	ID   int        `db:"id" json:"id"`
	Mode FilterMode `db:"mode" json:"mode"`
	Name string     `db:"name" json:"name"`
	// JSON-encoded filter string
	Filter              string `db:"filter" json:"filter"`
	RecommendationIndex int    `db:"recommendation_index" json:"recommendation_index"`
}

type SavedFilters []*SavedFilter

func (m *SavedFilters) Append(o interface{}) {
	*m = append(*m, o.(*SavedFilter))
}

func (m *SavedFilters) New() interface{} {
	return &SavedFilter{}
}
