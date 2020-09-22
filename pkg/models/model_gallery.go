package models

import (
	"database/sql"
)

type Gallery struct {
	ID        int             `db:"id" json:"id"`
	Path      string          `db:"path" json:"path"`
	Checksum  string          `db:"checksum" json:"checksum"`
	Title     sql.NullString  `db:"title" json:"title"`
	URL       sql.NullString  `db:"url" json:"url"`
	Details   sql.NullString  `db:"details" json:"details"`
	Rating    sql.NullInt64   `db:"rating" json:"rating"`
	StudioID  sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	SceneID   sql.NullInt64   `db:"scene_id,omitempty" json:"scene_id"`
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

const DefaultGthumbWidth int = 640
