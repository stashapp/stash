package models

import (
	"database/sql"
)

// Image stores the metadata for a single image.
type Image struct {
	ID          int                 `db:"id" json:"id"`
	Checksum    string              `db:"checksum" json:"checksum"`
	Path        string              `db:"path" json:"path"`
	Title       sql.NullString      `db:"title" json:"title"`
	Rating      sql.NullInt64       `db:"rating" json:"rating"`
	Organized   sql.NullBool        `db:"organized" json:"organized"`
	OCounter    int                 `db:"o_counter" json:"o_counter"`
	Size        sql.NullInt64       `db:"size" json:"size"`
	Width       sql.NullInt64       `db:"width" json:"width"`
	Height      sql.NullInt64       `db:"height" json:"height"`
	StudioID    sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	FileModTime NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

// ImagePartial represents part of a Image object. It is used to update
// the database entry. Only non-nil fields will be updated.
type ImagePartial struct {
	ID          int                  `db:"id" json:"id"`
	Checksum    *string              `db:"checksum" json:"checksum"`
	Path        *string              `db:"path" json:"path"`
	Title       *sql.NullString      `db:"title" json:"title"`
	Rating      *sql.NullInt64       `db:"rating" json:"rating"`
	Organized   *sql.NullBool        `db:"organized" json:"organized"`
	Size        *sql.NullInt64       `db:"size" json:"size"`
	Width       *sql.NullInt64       `db:"width" json:"width"`
	Height      *sql.NullInt64       `db:"height" json:"height"`
	StudioID    *sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	FileModTime *NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   *SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   *SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

// ImageFileType represents the file metadata for an image.
type ImageFileType struct {
	Size   *int `graphql:"size" json:"size"`
	Width  *int `graphql:"width" json:"width"`
	Height *int `graphql:"height" json:"height"`
}
