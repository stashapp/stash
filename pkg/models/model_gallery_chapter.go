package models

import (
	"database/sql"
)

type GalleryChapter struct {
	ID           int             `db:"id" json:"id"`
	Title        string          `db:"title" json:"title"`
	PageNumber   int             `db:"page_number" json:"page_number"`
	GalleryID    sql.NullInt64   `db:"gallery_id,omitempty" json:"gallery_id"`
	CreatedAt    SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt    SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type GalleryChapters []*GalleryChapter

func (m *GalleryChapters) Append(o interface{}) {
	*m = append(*m, o.(*GalleryChapter))
}

func (m *GalleryChapters) New() interface{} {
	return &GalleryChapter{}
}
