package models

import (
	"time"
)

type ScenePlayDate struct {
	ID       int       `json:"id"`
	SceneID  int       `json:"scene_id"`
	PlayDate time.Time `json:"playdate"`
}

type ScenePlayDates []*ScenePlayDate

func (m *ScenePlayDates) Append(o interface{}) {
	*m = append(*m, o.(*ScenePlayDate))
}

func (m *ScenePlayDates) New() interface{} {
	return &ScenePlayDate{}
}
