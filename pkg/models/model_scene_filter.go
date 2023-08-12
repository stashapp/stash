package models

import (
	"time"
)

type SceneFilter struct {
	ID          int       `json:"id"`
	Contrast    int       `json:"contrast"`
	Brightness  int       `json:"brightness"`
	Gamma       int       `json:"gamma"`
	Saturate    int       `json:"saturate"`
	HueRotate   int       `json:"hue_rotate"`
	Warmth      int       `json:"warmth"`
	Red         int       `json:"red"`
	Green       int       `json:"green"`
	Blue        int       `json:"blue"`
	Blur        int       `json:"blur"`
	Rotate      float64   `json:"rotate"`
	Scale       int       `json:"scale"`
	AspectRatio int       `json:"aspect_ratio"`
	SceneID     int       `json:"scene_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SceneFilters []*SceneFilter

func (m *SceneFilters) Append(o interface{}) {
	*m = append(*m, o.(*SceneFilter))
}

func (m *SceneFilters) New() interface{} {
	return &SceneFilter{}
}
