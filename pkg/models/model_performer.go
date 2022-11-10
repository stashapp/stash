package models

import (
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
)

type Performer struct {
	ID           int        `json:"id"`
	Checksum     string     `json:"checksum"`
	Name         string     `json:"name"`
	Gender       GenderEnum `json:"gender"`
	URL          string     `json:"url"`
	Twitter      string     `json:"twitter"`
	Instagram    string     `json:"instagram"`
	Birthdate    *Date      `json:"birthdate"`
	Ethnicity    string     `json:"ethnicity"`
	Country      string     `json:"country"`
	EyeColor     string     `json:"eye_color"`
	Height       *int       `json:"height"`
	Measurements string     `json:"measurements"`
	FakeTits     string     `json:"fake_tits"`
	CareerLength string     `json:"career_length"`
	Tattoos      string     `json:"tattoos"`
	Piercings    string     `json:"piercings"`
	Aliases      string     `json:"aliases"`
	Favorite     bool       `json:"favorite"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	// Rating expressed in 1-100 scale
	Rating        *int   `json:"rating"`
	Details       string `json:"details"`
	DeathDate     *Date  `json:"death_date"`
	HairColor     string `json:"hair_color"`
	Weight        *int   `json:"weight"`
	IgnoreAutoTag bool   `json:"ignore_auto_tag"`
}

// PerformerPartial represents part of a Performer object. It is used to update
// the database entry.
type PerformerPartial struct {
	ID           int
	Checksum     OptionalString
	Name         OptionalString
	Gender       OptionalString
	URL          OptionalString
	Twitter      OptionalString
	Instagram    OptionalString
	Birthdate    OptionalDate
	Ethnicity    OptionalString
	Country      OptionalString
	EyeColor     OptionalString
	Height       OptionalInt
	Measurements OptionalString
	FakeTits     OptionalString
	CareerLength OptionalString
	Tattoos      OptionalString
	Piercings    OptionalString
	Aliases      OptionalString
	Favorite     OptionalBool
	CreatedAt    OptionalTime
	UpdatedAt    OptionalTime
	// Rating expressed in 1-100 scale
	Rating        OptionalInt
	Details       OptionalString
	DeathDate     OptionalDate
	HairColor     OptionalString
	Weight        OptionalInt
	IgnoreAutoTag OptionalBool
}

func NewPerformer(name string) *Performer {
	currentTime := time.Now()
	return &Performer{
		Checksum:  md5.FromString(name),
		Name:      name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

func NewPerformerPartial() PerformerPartial {
	updatedTime := time.Now()
	return PerformerPartial{
		UpdatedAt: NewOptionalTime(updatedTime),
	}
}

type Performers []*Performer

func (p *Performers) Append(o interface{}) {
	*p = append(*p, o.(*Performer))
}

func (p *Performers) New() interface{} {
	return &Performer{}
}
