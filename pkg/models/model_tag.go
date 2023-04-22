package models

import (
	"database/sql"
	"time"
)

type Tag struct {
	ID            int            `db:"id" json:"id"`
	Name          string         `db:"name" json:"name"` // TODO make schema not null
	Description   sql.NullString `db:"description" json:"description"`
	IgnoreAutoTag bool           `db:"ignore_auto_tag" json:"ignore_auto_tag"`
	// TODO - this is only here because of database code in the models package
	ImageBlob sql.NullString  `db:"image_blob" json:"-"`
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type TagPartial struct {
	ID            int              `db:"id" json:"id"`
	Name          *string          `db:"name" json:"name"` // TODO make schema not null
	Description   *sql.NullString  `db:"description" json:"description"`
	IgnoreAutoTag *bool            `db:"ignore_auto_tag" json:"ignore_auto_tag"`
	CreatedAt     *SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt     *SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type TagPath struct {
	Tag
	Path string `db:"path" json:"path"`
}

func NewTag(name string) *Tag {
	currentTime := time.Now()
	return &Tag{
		Name:      name,
		CreatedAt: SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: SQLiteTimestamp{Timestamp: currentTime},
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
