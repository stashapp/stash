package models

import (
	"database/sql"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
)

type Movie struct {
	ID       int            `db:"id" json:"id"`
	Checksum string         `db:"checksum" json:"checksum"`
	Name     sql.NullString `db:"name" json:"name"`
	Aliases  sql.NullString `db:"aliases" json:"aliases"`
	Duration sql.NullInt64  `db:"duration" json:"duration"`
	Date     SQLiteDate     `db:"date" json:"date"`
	// Rating expressed in 1-100 scale
	Rating    sql.NullInt64   `db:"rating" json:"rating"`
	StudioID  sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	Director  sql.NullString  `db:"director" json:"director"`
	Synopsis  sql.NullString  `db:"synopsis" json:"synopsis"`
	URL       sql.NullString  `db:"url" json:"url"`
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type MoviePartial struct {
	ID       int             `db:"id" json:"id"`
	Checksum *string         `db:"checksum" json:"checksum"`
	Name     *sql.NullString `db:"name" json:"name"`
	Aliases  *sql.NullString `db:"aliases" json:"aliases"`
	Duration *sql.NullInt64  `db:"duration" json:"duration"`
	Date     *SQLiteDate     `db:"date" json:"date"`
	// Rating expressed in 1-100 scale
	Rating    *sql.NullInt64   `db:"rating" json:"rating"`
	StudioID  *sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	Director  *sql.NullString  `db:"director" json:"director"`
	Synopsis  *sql.NullString  `db:"synopsis" json:"synopsis"`
	URL       *sql.NullString  `db:"url" json:"url"`
	CreatedAt *SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt *SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

var DefaultMovieImage = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAA3XAAAN1wFCKJt4AAAAB3RJTUUH4wgVBQsJl1CMZAAAASJJREFUeNrt3N0JwyAYhlEj3cj9R3Cm5rbkqtAP+qrnGaCYHPwJpLlaa++mmLpbAERAgAgIEAEBIiBABERAgAgIEAEBIiBABERAgAgIEAHZuVflj40x4i94zhk9vqsVvEq6AsQqMP1EjORx20OACAgQRRx7T+zzcFBxcjNDfoB4ntQqTm5Awo7MlqywZxcgYQ+RlqywJ3ozJAQCSBiEJSsQA0gYBpDAgAARECACAkRAgAgIEAERECACAmSjUv6eAOSB8m8YIGGzBUjYbAESBgMkbBkDEjZbgITBAClcxiqQvEoatreYIWEBASIgJ4Gkf11ntXH3nS9uxfGWfJ5J9hAgAgJEQAQEiIAAERAgAgJEQAQEiIAAERAgAgJEQAQEiL7qBuc6RKLHxr0CAAAAAElFTkSuQmCC"

func NewMovie(name string) *Movie {
	currentTime := time.Now()
	return &Movie{
		Checksum:  md5.FromString(name),
		Name:      sql.NullString{String: name, Valid: true},
		CreatedAt: SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: SQLiteTimestamp{Timestamp: currentTime},
	}
}

type Movies []*Movie

func (m *Movies) Append(o interface{}) {
	*m = append(*m, o.(*Movie))
}

func (m *Movies) New() interface{} {
	return &Movie{}
}
