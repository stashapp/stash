package models

import (
	"database/sql"
)

type Performer struct {
	ID           int             `db:"id" json:"id"`
	Checksum     string          `db:"checksum" json:"checksum"`
	Name         sql.NullString  `db:"name" json:"name"`
	Gender       sql.NullString  `db:"gender" json:"gender"`
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
