package models

import (
	"time"
)

type Story struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Author    string     `json:"author"`
	URLs      []string   `json:"urls"`
	Date      *time.Time `json:"date"`
	Language  string     `json:"language"`
	TagLine   string     `json:"tag_line"`
	Details   string     `json:"details"`
	StudioID  *int       `json:"studio_id"`
	Rating    *int       `json:"rating"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	TagIDs       []int `json:"tag_ids"`
	PerformerIDs []int `json:"performer_ids"`
}

func NewStory() Story {
	currentTime := time.Now()
	return Story{
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
		TagIDs:       []int{},
		PerformerIDs: []int{},
		URLs:         []string{},
	}
}

type StoryPartial struct {
	Title       OptionalString
	Author      OptionalString
	Date        OptionalDate
	Language    OptionalString
	TagLine     OptionalString
	Details     OptionalString
	StudioID    OptionalInt
	Rating      OptionalInt
	TagIDs      *UpdateIDs
	PerformerIDs *UpdateIDs
	URLs        *UpdateStrings
}

func NewStoryPartial() StoryPartial {
	return StoryPartial{}
}
