package models

import (
	"database/sql"
)

type Performer struct {
	ID           int             `db:"id" json:"id"`
	Image        []byte          `db:"image" json:"image"`
	Checksum     string          `db:"checksum" json:"checksum"`
	Name         sql.NullString  `db:"name" json:"name"`
	URL          sql.NullString  `db:"url" json:"url"`
	Twitter      sql.NullString  `db:"twitter" json:"twitter"`
	Instagram    sql.NullString  `db:"instagram" json:"instagram"`
	Birthdate    SQLiteDate      `db:"birthdate" json:"birthdate"`
	Ethnicity    sql.NullString  `db:"ethnicity" json:"ethnicity"`
	Country      sql.NullString  `db:"country" json:"country"`
	EyeColor     sql.NullString  `db:"eye_color" json:"eye_color"`
	Height       sql.NullString  `db:"height" json:"height"`
	Measurements sql.NullString  `db:"measurements" json:"measurements"`
	FakeTits     sql.NullString  `db:"fake_tits" json:"fake_tits"`
	CareerLength sql.NullString  `db:"career_length" json:"career_length"`
	Tattoos      sql.NullString  `db:"tattoos" json:"tattoos"`
	Piercings    sql.NullString  `db:"piercings" json:"piercings"`
	Aliases      sql.NullString  `db:"aliases" json:"aliases"`
	Favorite     sql.NullBool    `db:"favorite" json:"favorite"`
	CreatedAt    SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt    SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

var DefaultPerformerImage string = "data:image/png;base64,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"
