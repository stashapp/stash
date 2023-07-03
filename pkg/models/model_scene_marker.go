package models

import (
	"time"
)

type SceneMarker struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Seconds      float64   `json:"seconds"`
	PrimaryTagID int       `json:"primary_tag_id"`
	SceneID      int       `json:"scene_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SceneMarkers []*SceneMarker

func (m *SceneMarkers) Append(o interface{}) {
	*m = append(*m, o.(*SceneMarker))
}

func (m *SceneMarkers) New() interface{} {
	return &SceneMarker{}
}
