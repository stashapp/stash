package models

import (
	"time"
)

type SceneODate struct {
	ID      int       `json:"id"`
	SceneID int       `json:"scene_id"`
	ODate   time.Time `json:"odate"`
}

type SceneODates []*SceneODate

func (m *SceneODates) Append(o interface{}) {
	*m = append(*m, o.(*SceneODate))
}

func (m *SceneODates) New() interface{} {
	return &SceneODate{}
}
