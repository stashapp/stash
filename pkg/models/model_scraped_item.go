package models

import (
	"context"
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

func (s *ScrapedStudio) ToStudio(endpoint string, excluded map[string]bool) *Studio {
	now := time.Now()

	// Populate a new studio from the input
	newStudio := Studio{
		Name: s.Name,
		StashIDs: NewRelatedStashIDs([]StashID{
			{
				Endpoint: endpoint,
				StashID:  *s.RemoteSiteID,
			},
		}),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if s.URL != nil && !excluded["url"] {
		newStudio.URL = *s.URL
	}

	if s.Parent != nil && s.Parent.StoredID != nil && !excluded["parent"] {
		parentId, _ := strconv.Atoi(*s.Parent.StoredID)
		newStudio.ParentID = &parentId
	}

	return &newStudio
}

func (s *ScrapedStudio) GetImage(ctx context.Context, excluded map[string]bool) ([]byte, error) {
	// Process the base 64 encoded image string
	if len(s.Images) > 0 && !excluded["image"] {
		var err error
		img, err := utils.ProcessImageInput(ctx, *s.Image)
		if err != nil {
			return nil, err
		}

		return img, nil
	}

	return nil, nil
}

func (s *ScrapedStudio) ToPartial(id *string, endpoint string, excluded map[string]bool, existingStashIDs []StashID) *StudioPartial {
	partial := StudioPartial{
		UpdatedAt: NewOptionalTime(time.Now()),
	}
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

	partial.StashIDs = &UpdateStashIDs{
		StashIDs: existingStashIDs,
		Mode:     RelationshipUpdateModeSet,
	}

	partial.StashIDs.Set(StashID{
		Endpoint: endpoint,
		StashID:  *s.RemoteSiteID,
	})

	return &partial
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
