package models

import (
	"database/sql"
)

type ScrapedStudio struct {
	// Set if studio matched
	StoredID     *string `json:"stored_id"`
	Name         string  `json:"name"`
	URL          *string `json:"url"`
	Image        *string `json:"image"`
	RemoteSiteID *string `json:"remote_site_id"`
}

func (ScrapedStudio) IsScrapedContent() {}

// A performer from a scraping operation...
type ScrapedPerformer struct {
	// Set if performer matched
	StoredID     *string       `json:"stored_id"`
	Name         *string       `json:"name"`
	Gender       *string       `json:"gender"`
	URL          *string       `json:"url"`
	Twitter      *string       `json:"twitter"`
	Instagram    *string       `json:"instagram"`
	Birthdate    *string       `json:"birthdate"`
	Ethnicity    *string       `json:"ethnicity"`
	Country      *string       `json:"country"`
	EyeColor     *string       `json:"eye_color"`
	Height       *string       `json:"height"`
	Measurements *string       `json:"measurements"`
	FakeTits     *string       `json:"fake_tits"`
	CareerLength *string       `json:"career_length"`
	Tattoos      *string       `json:"tattoos"`
	Piercings    *string       `json:"piercings"`
	Aliases      *string       `json:"aliases"`
	Tags         []*ScrapedTag `json:"tags"`
	// This should be a base64 encoded data URL
	Image        *string  `json:"image"`
	Images       []string `json:"images"`
	Details      *string  `json:"details"`
	DeathDate    *string  `json:"death_date"`
	HairColor    *string  `json:"hair_color"`
	Weight       *string  `json:"weight"`
	RemoteSiteID *string  `json:"remote_site_id"`
}

func (ScrapedPerformer) IsScrapedContent() {}

type ScrapedTag struct {
	// Set if tag matched
	StoredID *string `json:"stored_id"`
	Name     string  `json:"name"`
}

func (ScrapedTag) IsScrapedContent() {}

// A movie from a scraping operation...
type ScrapedMovie struct {
	StoredID *string        `json:"stored_id"`
	Name     *string        `json:"name"`
	Aliases  *string        `json:"aliases"`
	Duration *string        `json:"duration"`
	Date     *string        `json:"date"`
	Rating   *string        `json:"rating"`
	Director *string        `json:"director"`
	URL      *string        `json:"url"`
	Synopsis *string        `json:"synopsis"`
	Studio   *ScrapedStudio `json:"studio"`
	// This should be a base64 encoded data URL
	FrontImage *string `json:"front_image"`
	// This should be a base64 encoded data URL
	BackImage *string `json:"back_image"`
}

func (ScrapedMovie) IsScrapedContent() {}

type ScrapedItem struct {
	ID              int             `db:"id" json:"id"`
	Title           sql.NullString  `db:"title" json:"title"`
	Description     sql.NullString  `db:"description" json:"description"`
	URL             sql.NullString  `db:"url" json:"url"`
	Date            SQLiteDate      `db:"date" json:"date"`
	Rating          sql.NullString  `db:"rating" json:"rating"`
	Tags            sql.NullString  `db:"tags" json:"tags"`
	Models          sql.NullString  `db:"models" json:"models"`
	Episode         sql.NullInt64   `db:"episode" json:"episode"`
	GalleryFilename sql.NullString  `db:"gallery_filename" json:"gallery_filename"`
	GalleryURL      sql.NullString  `db:"gallery_url" json:"gallery_url"`
	VideoFilename   sql.NullString  `db:"video_filename" json:"video_filename"`
	VideoURL        sql.NullString  `db:"video_url" json:"video_url"`
	StudioID        sql.NullInt64   `db:"studio_id,omitempty" json:"studio_id"`
	CreatedAt       SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt       SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

type ScrapedItems []*ScrapedItem

func (s *ScrapedItems) Append(o interface{}) {
	*s = append(*s, o.(*ScrapedItem))
}

func (s *ScrapedItems) New() interface{} {
	return &ScrapedItem{}
}
