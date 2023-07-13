package models

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/utils"
)

type ScrapedStudio struct {
	// Set if studio matched
	StoredID     *string        `json:"stored_id"`
	Name         string         `json:"name"`
	URL          *string        `json:"url"`
	Parent       *ScrapedStudio `json:"parent"`
	Image        *string        `json:"image"`
	Images       []string       `json:"images"`
	RemoteSiteID *string        `json:"remote_site_id"`
}

func (ScrapedStudio) IsScrapedContent() {}

func (s *ScrapedStudio) ToStudio(ctx context.Context, endpoint string, excluded map[string]bool) (*Studio, error) {
	// Populate a new studio from the input
	newStudio := Studio{
		Name: s.Name,
		StashIDs: NewRelatedStashIDs([]StashID{
			{
				Endpoint: endpoint,
				StashID:  *s.RemoteSiteID,
			},
		}),
	}

	if s.URL != nil && !excluded["url"] {
		newStudio.URL = *s.URL
	}

	if s.Parent != nil && s.Parent.StoredID != nil && !excluded["parent"] {
		parentId, _ := strconv.Atoi(*s.Parent.StoredID)
		newStudio.ParentID = &parentId
	}

	// Process the base 64 encoded image string
	if s.Image != nil && !excluded["image"] {
		var err error
		newStudio.ImageBytes, err = utils.ProcessImageInput(ctx, *s.Image)
		if err != nil {
			return nil, err
		}
	}

	return &newStudio, nil
}

func (s *ScrapedStudio) ToPartial(ctx context.Context, id *string, endpoint string, excluded map[string]bool, existingStashIDs []StashID) (*StudioPartial, error) {
	partial := StudioPartial{}
	partial.ID, _ = strconv.Atoi(*id)

	if s.Name != "" && !excluded["name"] {
		partial.Name = NewOptionalString(s.Name)

	}

	if s.URL != nil && !excluded["url"] {
		partial.URL = NewOptionalString(*s.URL)
	}

	if s.Parent != nil && !excluded["parent"] {
		if s.Parent.StoredID != nil {
			parentID, _ := strconv.Atoi(*s.Parent.StoredID)
			if parentID > 0 {
				// This is to be set directly as we know it has a value and the translator won't have the field
				partial.ParentID = NewOptionalInt(parentID)
			}
		}
	} else {
		partial.ParentID = NewOptionalIntPtr(nil)
	}

	// Process the base 64 encoded image string
	if len(s.Images) > 0 && !excluded["image"] {
		partial.Image = OptionalBytes{
			Set: true,
		}

		var err error
		partial.Image.Value, err = utils.ProcessImageInput(ctx, s.Images[0])
		if err != nil {
			return nil, err
		}
	}

	partial.StashIDs = &UpdateStashIDs{
		StashIDs: existingStashIDs,
		Mode:     RelationshipUpdateModeSet,
	}

	partial.StashIDs.Set(StashID{
		Endpoint: endpoint,
		StashID:  *s.RemoteSiteID,
	})

	return &partial, nil
}

// A performer from a scraping operation...
type ScrapedPerformer struct {
	// Set if performer matched
	StoredID       *string       `json:"stored_id"`
	Name           *string       `json:"name"`
	Disambiguation *string       `json:"disambiguation"`
	Gender         *string       `json:"gender"`
	URL            *string       `json:"url"`
	Twitter        *string       `json:"twitter"`
	Instagram      *string       `json:"instagram"`
	Birthdate      *string       `json:"birthdate"`
	Ethnicity      *string       `json:"ethnicity"`
	Country        *string       `json:"country"`
	EyeColor       *string       `json:"eye_color"`
	Height         *string       `json:"height"`
	Measurements   *string       `json:"measurements"`
	FakeTits       *string       `json:"fake_tits"`
	PenisLength    *string       `json:"penis_length"`
	Circumcised    *string       `json:"circumcised"`
	CareerLength   *string       `json:"career_length"`
	Tattoos        *string       `json:"tattoos"`
	Piercings      *string       `json:"piercings"`
	Aliases        *string       `json:"aliases"`
	Tags           []*ScrapedTag `json:"tags"`
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
	ID              int            `db:"id" json:"id"`
	Title           sql.NullString `db:"title" json:"title"`
	Code            sql.NullString `db:"code" json:"code"`
	Description     sql.NullString `db:"description" json:"description"`
	Director        sql.NullString `db:"director" json:"director"`
	URL             sql.NullString `db:"url" json:"url"`
	Date            *Date          `db:"date" json:"date"`
	Rating          sql.NullString `db:"rating" json:"rating"`
	Tags            sql.NullString `db:"tags" json:"tags"`
	Models          sql.NullString `db:"models" json:"models"`
	Episode         sql.NullInt64  `db:"episode" json:"episode"`
	GalleryFilename sql.NullString `db:"gallery_filename" json:"gallery_filename"`
	GalleryURL      sql.NullString `db:"gallery_url" json:"gallery_url"`
	VideoFilename   sql.NullString `db:"video_filename" json:"video_filename"`
	VideoURL        sql.NullString `db:"video_url" json:"video_url"`
	StudioID        sql.NullInt64  `db:"studio_id,omitempty" json:"studio_id"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
}

type ScrapedItems []*ScrapedItem

func (s *ScrapedItems) Append(o interface{}) {
	*s = append(*s, o.(*ScrapedItem))
}

func (s *ScrapedItems) New() interface{} {
	return &ScrapedItem{}
}
