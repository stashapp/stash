package models

import (
	"database/sql"
)

type SceneMarker struct {
	ID           int             `db:"id" json:"id"`
	Title        string          `db:"title" json:"title"`
	Seconds      float64         `db:"seconds" json:"seconds"`
	PrimaryTagID int             `db:"primary_tag_id" json:"primary_tag_id"`
	SceneID      sql.NullInt64   `db:"scene_id,omitempty" json:"scene_id"`
	CreatedAt    SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt    SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type SceneMarkers []*SceneMarker

func (m *SceneMarkers) Append(o interface{}) {
	*m = append(*m, o.(*SceneMarker))
}

func (m *SceneMarkers) New() interface{} {
	return &SceneMarker{}
}
