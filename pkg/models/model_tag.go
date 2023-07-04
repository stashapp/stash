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

type TagPartial struct {
	Name          OptionalString
	Description   OptionalString
	IgnoreAutoTag OptionalBool
	CreatedAt     OptionalTime
	UpdatedAt     OptionalTime
}

type TagPath struct {
	Tag
	Path string `json:"path"`
}

func NewTag(name string) *Tag {
	currentTime := time.Now()
	return &Tag{
		Name:      name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

func NewTagPartial() TagPartial {
	updatedTime := time.Now()
	return TagPartial{
		UpdatedAt: NewOptionalTime(updatedTime),
	}
}

type Tags []*Tag

func (t *Tags) Append(o interface{}) {
	*t = append(*t, o.(*Tag))
}

func (t *Tags) New() interface{} {
	return &Tag{}
}

type TagPaths []*TagPath

func (t *TagPaths) Append(o interface{}) {
	*t = append(*t, o.(*TagPath))
}

func (t *TagPaths) New() interface{} {
	return &TagPath{}
}
