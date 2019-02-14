package models

import (
	"database/sql"
)

type Scene struct {
	ID         int             `db:"id" json:"id"`
	Checksum   string          `db:"checksum" json:"checksum"`
	Path       string          `db:"path" json:"path"`
	Title      sql.NullString  `db:"title" json:"title"`
	Details    sql.NullString  `db:"details" json:"details"`
	URL        sql.NullString  `db:"url" json:"url"`
	Date       sql.NullString  `db:"date" json:"date"` // TODO dates?
	Rating     sql.NullInt64   `db:"rating" json:"rating"`
	Size       sql.NullString  `db:"size" json:"size"`
	Duration   sql.NullFloat64 `db:"duration" json:"duration"`
	VideoCodec sql.NullString  `db:"video_codec" json:"video_codec"`
	AudioCodec sql.NullString  `db:"audio_codec" json:"audio_codec"`
	Width      sql.NullInt64   `db:"width" json:"width"`
	Height     sql.NullInt64   `db:"height" json:"height"`
	Framerate  sql.NullFloat64 `db:"framerate" json:"framerate"`
	Bitrate    sql.NullInt64   `db:"bitrate" json:"bitrate"`
	StudioID   sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	CreatedAt  SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt  SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}
