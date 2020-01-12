package models

import (
	"database/sql"
	"path/filepath"
)

type Scene struct {
	ID         int             `db:"id" json:"id"`
	Checksum   string          `db:"checksum" json:"checksum"`
	Path       string          `db:"path" json:"path"`
	Cover      []byte          `db:"cover" json:"cover"`
	Title      sql.NullString  `db:"title" json:"title"`
	Details    sql.NullString  `db:"details" json:"details"`
	URL        sql.NullString  `db:"url" json:"url"`
	Date       SQLiteDate      `db:"date" json:"date"`
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
	DvdID      sql.NullInt64   `db:"dvd_id,omitempty" json:"dvd_id"`
	CreatedAt  SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt  SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type ScenePartial struct {
	ID         int              `db:"id" json:"id"`
	Checksum   *string          `db:"checksum" json:"checksum"`
	Path       *string          `db:"path" json:"path"`
	Cover      *[]byte          `db:"cover" json:"cover"`
	Title      *sql.NullString  `db:"title" json:"title"`
	Details    *sql.NullString  `db:"details" json:"details"`
	URL        *sql.NullString  `db:"url" json:"url"`
	Date       *SQLiteDate      `db:"date" json:"date"`
	Rating     *sql.NullInt64   `db:"rating" json:"rating"`
	Size       *sql.NullString  `db:"size" json:"size"`
	Duration   *sql.NullFloat64 `db:"duration" json:"duration"`
	VideoCodec *sql.NullString  `db:"video_codec" json:"video_codec"`
	AudioCodec *sql.NullString  `db:"audio_codec" json:"audio_codec"`
	Width      *sql.NullInt64   `db:"width" json:"width"`
	Height     *sql.NullInt64   `db:"height" json:"height"`
	Framerate  *sql.NullFloat64 `db:"framerate" json:"framerate"`
	Bitrate    *sql.NullInt64   `db:"bitrate" json:"bitrate"`
	StudioID   *sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	DvdID      *sql.NullInt64   `db:"dvd_id,omitempty" json:"dvd_id"`
	CreatedAt  *SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt  *SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

func (s Scene) GetTitle() string {
	if s.Title.String != "" {
		return s.Title.String
	}

	return filepath.Base(s.Path)
}

type SceneFileType struct {
	Size       *string  `graphql:"size" json:"size"`
	Duration   *float64 `graphql:"duration" json:"duration"`
	VideoCodec *string  `graphql:"video_codec" json:"video_codec"`
	AudioCodec *string  `graphql:"audio_codec" json:"audio_codec"`
	Width      *int     `graphql:"width" json:"width"`
	Height     *int     `graphql:"height" json:"height"`
	Framerate  *float64 `graphql:"framerate" json:"framerate"`
	Bitrate    *int     `graphql:"bitrate" json:"bitrate"`
}
