package models

import (
	"database/sql"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
)

type Studio struct {
	ID            int             `db:"id" json:"id"`
	Checksum      string          `db:"checksum" json:"checksum"`
	Name          sql.NullString  `db:"name" json:"name"`
	URL           sql.NullString  `db:"url" json:"url"`
	ParentID      sql.NullInt64   `db:"parent_id,omitempty" json:"parent_id"`
	CreatedAt     SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt     SQLiteTimestamp `db:"updated_at" json:"updated_at"`
	Rating        sql.NullInt64   `db:"rating" json:"rating"`
	Details       sql.NullString  `db:"details" json:"details"`
	IgnoreAutoTag bool            `db:"ignore_auto_tag" json:"ignore_auto_tag"`
}

type StudioPartial struct {
	ID            int              `db:"id" json:"id"`
	Checksum      *string          `db:"checksum" json:"checksum"`
	Name          *sql.NullString  `db:"name" json:"name"`
	URL           *sql.NullString  `db:"url" json:"url"`
	ParentID      *sql.NullInt64   `db:"parent_id,omitempty" json:"parent_id"`
	CreatedAt     *SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt     *SQLiteTimestamp `db:"updated_at" json:"updated_at"`
	Rating        *sql.NullInt64   `db:"rating" json:"rating"`
	Details       *sql.NullString  `db:"details" json:"details"`
	IgnoreAutoTag *bool            `db:"ignore_auto_tag" json:"ignore_auto_tag"`
}

var DefaultStudioImage = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAYAAABw4pVUAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAA3XAAAN1wFCKJt4AAAAB3RJTUUH4wgVBQsJl1CMZAAAASJJREFUeNrt3N0JwyAYhlEj3cj9R3Cm5rbkqtAP+qrnGaCYHPwJpLlaa++mmLpbAERAgAgIEAEBIiBABERAgAgIEAEBIiBABERAgAgIEAHZuVflj40x4i94zhk9vqsVvEq6AsQqMP1EjORx20OACAgQRRx7T+zzcFBxcjNDfoB4ntQqTm5Awo7MlqywZxcgYQ+RlqywJ3ozJAQCSBiEJSsQA0gYBpDAgAARECACAkRAgAgIEAERECACAmSjUv6eAOSB8m8YIGGzBUjYbAESBgMkbBkDEjZbgITBAClcxiqQvEoatreYIWEBASIgJ4Gkf11ntXH3nS9uxfGWfJ5J9hAgAgJEQAQEiIAAERAgAgJEQAQEiIAAERAgAgJEQAQEiL7qBuc6RKLHxr0CAAAAAElFTkSuQmCC"

func NewStudio(name string) *Studio {
	currentTime := time.Now()
	return &Studio{
		Checksum:  md5.FromString(name),
		Name:      sql.NullString{String: name, Valid: true},
		CreatedAt: SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: SQLiteTimestamp{Timestamp: currentTime},
	}
}

type Studios []*Studio

func (s *Studios) Append(o interface{}) {
	*s = append(*s, o.(*Studio))
}

func (s *Studios) New() interface{} {
	return &Studio{}
}
