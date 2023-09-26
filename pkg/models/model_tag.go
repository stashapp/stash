package models

import (
	"time"
)

type Tag struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	IgnoreAutoTag bool      `json:"ignore_auto_tag"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewTag() Tag {
	currentTime := time.Now()
	return Tag{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

type TagPartial struct {
	Name          OptionalString
	Description   OptionalString
	IgnoreAutoTag OptionalBool
	CreatedAt     OptionalTime
	UpdatedAt     OptionalTime
}

func NewTagPartial() TagPartial {
	currentTime := time.Now()
	return TagPartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}

type TagPath struct {
	Tag
	Path string `json:"path"`
}
