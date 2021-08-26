package models

import (
	"database/sql"
)

type ScrapedItem struct {
	ID              int             `db:"id" json:"id"`
	Title           sql.NullString  `db:"title" json:"title"`
	Description     sql.NullString  `db:"description" json:"description"`
	URL             sql.NullString  `db:"url" json:"url"`
	Date            SQLiteDate      `db:"date" json:"date"`
	Rating          sql.NullString  `db:"rating" json:"rating"`
	Tags            sql.NullString  `db:"tags" json:"tags"`
	Models          sql.NullString  `db:"models" json:"models"`
	Episode         sql.NullInt64   `db:"episode" json:"episode"`
	GalleryFilename sql.NullString  `db:"gallery_filename" json:"gallery_filename"`
	GalleryURL      sql.NullString  `db:"gallery_url" json:"gallery_url"`
	VideoFilename   sql.NullString  `db:"video_filename" json:"video_filename"`
	VideoURL        sql.NullString  `db:"video_url" json:"video_url"`
	StudioID        sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	CreatedAt       SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt       SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type ScrapedItems []*ScrapedItem

func (s *ScrapedItems) Append(o interface{}) {
	*s = append(*s, o.(*ScrapedItem))
}

func (s *ScrapedItems) New() interface{} {
	return &ScrapedItem{}
}
