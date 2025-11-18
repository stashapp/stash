package models

import (
	"time"
)

type SceneMarker struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Seconds      float64   `json:"seconds"`
	EndSeconds   *float64  `json:"end_seconds"`
	PrimaryTagID int       `json:"primary_tag_id"`
	SceneID      int       `json:"scene_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewSceneMarker() SceneMarker {
	currentTime := time.Now()
	return SceneMarker{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

// SceneMarkerPartial represents part of a SceneMarker object.
// It is used to update the database entry.
type SceneMarkerPartial struct {
	Title        OptionalString
	Seconds      OptionalFloat64
	EndSeconds   OptionalFloat64
	PrimaryTagID OptionalInt
	TagIDs       *UpdateIDs
	SceneID      OptionalInt
	CreatedAt    OptionalTime
	UpdatedAt    OptionalTime
}

func NewSceneMarkerPartial() SceneMarkerPartial {
	currentTime := time.Now()
	return SceneMarkerPartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}
