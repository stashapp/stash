package models

import (
	"database/sql"
)

type Gallery struct {
	ID          int                 `db:"id" json:"id"`
	Path        sql.NullString      `db:"path" json:"path"`
	Checksum    string              `db:"checksum" json:"checksum"`
	Zip         bool                `db:"zip" json:"zip"`
	Title       sql.NullString      `db:"title" json:"title"`
	URL         sql.NullString      `db:"url" json:"url"`
	Date        SQLiteDate          `db:"date" json:"date"`
	Details     sql.NullString      `db:"details" json:"details"`
	Rating      sql.NullInt64       `db:"rating" json:"rating"`
	Organized   bool                `db:"organized" json:"organized"`
	StudioID    sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	SceneID     sql.NullInt64       `db:"scene_id,omitempty" json:"scene_id"`
	FileModTime NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

// GalleryPartial represents part of a Gallery object. It is used to update
// the database entry. Only non-nil fields will be updated.
type GalleryPartial struct {
	ID          int                  `db:"id" json:"id"`
	Path        *sql.NullString      `db:"path" json:"path"`
	Checksum    *string              `db:"checksum" json:"checksum"`
	Title       *sql.NullString      `db:"title" json:"title"`
	URL         *sql.NullString      `db:"url" json:"url"`
	Date        *SQLiteDate          `db:"date" json:"date"`
	Details     *sql.NullString      `db:"details" json:"details"`
	Rating      *sql.NullInt64       `db:"rating" json:"rating"`
	Organized   *bool                `db:"organized" json:"organized"`
	StudioID    *sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	SceneID     *sql.NullInt64       `db:"scene_id,omitempty" json:"scene_id"`
	FileModTime *NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   *SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   *SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

const DefaultGthumbWidth int = 640
