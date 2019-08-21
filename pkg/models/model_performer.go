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

var DefaultPerformerImage string = "data:image/jpeg;base64,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"
