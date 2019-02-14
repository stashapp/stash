package models

import (
	"database/sql"
)

type Studio struct {
	ID        int             `db:"id" json:"id"`
	Image     []byte          `db:"image" json:"image"`
	Checksum  string          `db:"checksum" json:"checksum"`
	Name      sql.NullString  `db:"name" json:"name"`
	URL       sql.NullString  `db:"url" json:"url"`
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}
