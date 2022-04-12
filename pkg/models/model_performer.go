package models

import (
	"database/sql"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
)

type Performer struct {
	ID            int             `db:"id" json:"id"`
	Checksum      string          `db:"checksum" json:"checksum"`
	Name          sql.NullString  `db:"name" json:"name"`
	Gender        sql.NullString  `db:"gender" json:"gender"`
	URL           sql.NullString  `db:"url" json:"url"`
	Twitter       sql.NullString  `db:"twitter" json:"twitter"`
	Instagram     sql.NullString  `db:"instagram" json:"instagram"`
	Birthdate     SQLiteDate      `db:"birthdate" json:"birthdate"`
	Ethnicity     sql.NullString  `db:"ethnicity" json:"ethnicity"`
	Country       sql.NullString  `db:"country" json:"country"`
	EyeColor      sql.NullString  `db:"eye_color" json:"eye_color"`
	Height        sql.NullString  `db:"height" json:"height"`
	Measurements  sql.NullString  `db:"measurements" json:"measurements"`
	FakeTits      sql.NullString  `db:"fake_tits" json:"fake_tits"`
	CareerLength  sql.NullString  `db:"career_length" json:"career_length"`
	Tattoos       sql.NullString  `db:"tattoos" json:"tattoos"`
	Piercings     sql.NullString  `db:"piercings" json:"piercings"`
	Aliases       sql.NullString  `db:"aliases" json:"aliases"`
	Favorite      sql.NullBool    `db:"favorite" json:"favorite"`
	CreatedAt     SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt     SQLiteTimestamp `db:"updated_at" json:"updated_at"`
	Rating        sql.NullInt64   `db:"rating" json:"rating"`
	Details       sql.NullString  `db:"details" json:"details"`
	DeathDate     SQLiteDate      `db:"death_date" json:"death_date"`
	HairColor     sql.NullString  `db:"hair_color" json:"hair_color"`
	Weight        sql.NullInt64   `db:"weight" json:"weight"`
	IgnoreAutoTag bool            `db:"ignore_auto_tag" json:"ignore_auto_tag"`
}

type PerformerPartial struct {
	ID            int              `db:"id" json:"id"`
	Checksum      *string          `db:"checksum" json:"checksum"`
	Name          *sql.NullString  `db:"name" json:"name"`
	Gender        *sql.NullString  `db:"gender" json:"gender"`
	URL           *sql.NullString  `db:"url" json:"url"`
	Twitter       *sql.NullString  `db:"twitter" json:"twitter"`
	Instagram     *sql.NullString  `db:"instagram" json:"instagram"`
	Birthdate     *SQLiteDate      `db:"birthdate" json:"birthdate"`
	Ethnicity     *sql.NullString  `db:"ethnicity" json:"ethnicity"`
	Country       *sql.NullString  `db:"country" json:"country"`
	EyeColor      *sql.NullString  `db:"eye_color" json:"eye_color"`
	Height        *sql.NullString  `db:"height" json:"height"`
	Measurements  *sql.NullString  `db:"measurements" json:"measurements"`
	FakeTits      *sql.NullString  `db:"fake_tits" json:"fake_tits"`
	CareerLength  *sql.NullString  `db:"career_length" json:"career_length"`
	Tattoos       *sql.NullString  `db:"tattoos" json:"tattoos"`
	Piercings     *sql.NullString  `db:"piercings" json:"piercings"`
	Aliases       *sql.NullString  `db:"aliases" json:"aliases"`
	Favorite      *sql.NullBool    `db:"favorite" json:"favorite"`
	CreatedAt     *SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt     *SQLiteTimestamp `db:"updated_at" json:"updated_at"`
	Rating        *sql.NullInt64   `db:"rating" json:"rating"`
	Details       *sql.NullString  `db:"details" json:"details"`
	DeathDate     *SQLiteDate      `db:"death_date" json:"death_date"`
	HairColor     *sql.NullString  `db:"hair_color" json:"hair_color"`
	Weight        *sql.NullInt64   `db:"weight" json:"weight"`
	IgnoreAutoTag *bool            `db:"ignore_auto_tag" json:"ignore_auto_tag"`
}

func NewPerformer(name string) *Performer {
	currentTime := time.Now()
	return &Performer{
		Checksum:  md5.FromString(name),
		Name:      sql.NullString{String: name, Valid: true},
		Favorite:  sql.NullBool{Bool: false, Valid: true},
		CreatedAt: SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: SQLiteTimestamp{Timestamp: currentTime},
	}
}

type Performers []*Performer

func (p *Performers) Append(o interface{}) {
	*p = append(*p, o.(*Performer))
}

func (p *Performers) New() interface{} {
	return &Performer{}
}
