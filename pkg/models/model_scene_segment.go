package models

import (
	"fmt"
	"time"
)

type SceneSegment struct {
	ID           int       `json:"id"`
	SceneID      int       `json:"scene_id"`
	Title        string    `json:"title"`
	StartSeconds float64   `json:"start_seconds"`
	EndSeconds   float64   `json:"end_seconds"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SceneSegmentPartial struct {
	ID           int
	SceneID      OptionalInt
	Title        OptionalString
	StartSeconds OptionalFloat64
	EndSeconds   OptionalFloat64
	CreatedAt    OptionalTime
	UpdatedAt    OptionalTime
}

func NewSceneSegment() SceneSegment {
	currentTime := time.Now()
	return SceneSegment{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

func (s *SceneSegment) LoadRelationships(r SceneSegmentReader) error {
	return nil
}

func (s *SceneSegment) Validate() error {
	if s.Title == "" {
		return fmt.Errorf("title is required")
	}
	if s.StartSeconds < 0 {
		return fmt.Errorf("start_seconds must be >= 0")
	}
	if s.EndSeconds <= s.StartSeconds {
		return fmt.Errorf("end_seconds must be > start_seconds")
	}
	return nil
}
