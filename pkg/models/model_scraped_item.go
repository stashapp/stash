package models

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
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
	// Populate a new studio from the input
	ret := NewStudio()
	ret.Name = s.Name

	if s.RemoteSiteID != nil && endpoint != "" {
		ret.StashIDs = NewRelatedStashIDs([]StashID{
			{
				Endpoint: endpoint,
				StashID:  *s.RemoteSiteID,
			},
		})
	}

	if s.URL != nil && !excluded["url"] {
		ret.URL = *s.URL
	}

	if s.Parent != nil && s.Parent.StoredID != nil && !excluded["parent"] && !excluded["parent_studio"] {
		parentId, _ := strconv.Atoi(*s.Parent.StoredID)
		ret.ParentID = &parentId
	}

	return &ret
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
	ret := NewStudioPartial()
	ret.ID, _ = strconv.Atoi(*id)

	if s.Name != "" && !excluded["name"] {
		ret.Name = NewOptionalString(s.Name)
	}

	if s.URL != nil && !excluded["url"] {
		ret.URL = NewOptionalString(*s.URL)
	}

	if s.Parent != nil && !excluded["parent"] {
		if s.Parent.StoredID != nil {
			parentID, _ := strconv.Atoi(*s.Parent.StoredID)
			if parentID > 0 {
				// This is to be set directly as we know it has a value and the translator won't have the field
				ret.ParentID = NewOptionalInt(parentID)
			}
		}
	} else {
		ret.ParentID = NewOptionalIntPtr(nil)
	}

	if s.RemoteSiteID != nil && endpoint != "" {
		ret.StashIDs = &UpdateStashIDs{
			StashIDs: existingStashIDs,
			Mode:     RelationshipUpdateModeSet,
		}
		ret.StashIDs.Set(StashID{
			Endpoint: endpoint,
			StashID:  *s.RemoteSiteID,
		})
	}

	return &ret
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
	Image        *string  `json:"image"` // deprecated: use Images
	Images       []string `json:"images"`
	Details      *string  `json:"details"`
	DeathDate    *string  `json:"death_date"`
	HairColor    *string  `json:"hair_color"`
	Weight       *string  `json:"weight"`
	RemoteSiteID *string  `json:"remote_site_id"`
}

func (ScrapedPerformer) IsScrapedContent() {}

func (p *ScrapedPerformer) ToPerformer(endpoint string, excluded map[string]bool) *Performer {
	ret := NewPerformer()
	ret.Name = *p.Name

	if p.Aliases != nil && !excluded["aliases"] {
		ret.Aliases = NewRelatedStrings(stringslice.FromString(*p.Aliases, ","))
	}
	if p.Birthdate != nil && !excluded["birthdate"] {
		date, err := ParseDate(*p.Birthdate)
		if err == nil {
			ret.Birthdate = &date
		}
	}
	if p.DeathDate != nil && !excluded["death_date"] {
		date, err := ParseDate(*p.DeathDate)
		if err == nil {
			ret.DeathDate = &date
		}
	}
	if p.CareerLength != nil && !excluded["career_length"] {
		ret.CareerLength = *p.CareerLength
	}
	if p.Country != nil && !excluded["country"] {
		ret.Country = *p.Country
	}
	if p.Ethnicity != nil && !excluded["ethnicity"] {
		ret.Ethnicity = *p.Ethnicity
	}
	if p.EyeColor != nil && !excluded["eye_color"] {
		ret.EyeColor = *p.EyeColor
	}
	if p.HairColor != nil && !excluded["hair_color"] {
		ret.HairColor = *p.HairColor
	}
	if p.FakeTits != nil && !excluded["fake_tits"] {
		ret.FakeTits = *p.FakeTits
	}
	if p.Gender != nil && !excluded["gender"] {
		v := GenderEnum(*p.Gender)
		if v.IsValid() {
			ret.Gender = &v
		}
	}
	if p.Height != nil && !excluded["height"] {
		h, err := strconv.Atoi(*p.Height)
		if err == nil {
			ret.Height = &h
		}
	}
	if p.Weight != nil && !excluded["weight"] {
		w, err := strconv.Atoi(*p.Weight)
		if err == nil {
			ret.Weight = &w
		}
	}
	if p.Instagram != nil && !excluded["instagram"] {
		ret.Instagram = *p.Instagram
	}
	if p.Measurements != nil && !excluded["measurements"] {
		ret.Measurements = *p.Measurements
	}
	if p.Disambiguation != nil && !excluded["disambiguation"] {
		ret.Disambiguation = *p.Disambiguation
	}
	if p.Details != nil && !excluded["details"] {
		ret.Details = *p.Details
	}
	if p.Piercings != nil && !excluded["piercings"] {
		ret.Piercings = *p.Piercings
	}
	if p.Tattoos != nil && !excluded["tattoos"] {
		ret.Tattoos = *p.Tattoos
	}
	if p.PenisLength != nil && !excluded["penis_length"] {
		l, err := strconv.ParseFloat(*p.PenisLength, 64)
		if err == nil {
			ret.PenisLength = &l
		}
	}
	if p.Circumcised != nil && !excluded["circumcised"] {
		v := CircumisedEnum(*p.Circumcised)
		if v.IsValid() {
			ret.Circumcised = &v
		}
	}
	if p.Twitter != nil && !excluded["twitter"] {
		ret.Twitter = *p.Twitter
	}
	if p.URL != nil && !excluded["url"] {
		ret.URL = *p.URL
	}

	if p.RemoteSiteID != nil && endpoint != "" {
		ret.StashIDs = NewRelatedStashIDs([]StashID{
			{
				Endpoint: endpoint,
				StashID:  *p.RemoteSiteID,
			},
		})
	}

	return &ret
}

func (p *ScrapedPerformer) GetImage(ctx context.Context, excluded map[string]bool) ([]byte, error) {
	// Process the base 64 encoded image string
	if len(p.Images) > 0 && !excluded["image"] {
		var err error
		img, err := utils.ProcessImageInput(ctx, p.Images[0])
		if err != nil {
			return nil, err
		}

		return img, nil
	}

	return nil, nil
}

func (p *ScrapedPerformer) ToPartial(endpoint string, excluded map[string]bool, existingStashIDs []StashID) PerformerPartial {
	ret := NewPerformerPartial()

	if p.Aliases != nil && !excluded["aliases"] {
		ret.Aliases = &UpdateStrings{
			Values: stringslice.FromString(*p.Aliases, ","),
			Mode:   RelationshipUpdateModeSet,
		}
	}
	if p.Birthdate != nil && !excluded["birthdate"] {
		date, err := ParseDate(*p.Birthdate)
		if err == nil {
			ret.Birthdate = NewOptionalDate(date)
		}
	}
	if p.DeathDate != nil && !excluded["death_date"] {
		date, err := ParseDate(*p.DeathDate)
		if err == nil {
			ret.DeathDate = NewOptionalDate(date)
		}
	}
	if p.CareerLength != nil && !excluded["career_length"] {
		ret.CareerLength = NewOptionalString(*p.CareerLength)
	}
	if p.Country != nil && !excluded["country"] {
		ret.Country = NewOptionalString(*p.Country)
	}
	if p.Ethnicity != nil && !excluded["ethnicity"] {
		ret.Ethnicity = NewOptionalString(*p.Ethnicity)
	}
	if p.EyeColor != nil && !excluded["eye_color"] {
		ret.EyeColor = NewOptionalString(*p.EyeColor)
	}
	if p.HairColor != nil && !excluded["hair_color"] {
		ret.HairColor = NewOptionalString(*p.HairColor)
	}
	if p.FakeTits != nil && !excluded["fake_tits"] {
		ret.FakeTits = NewOptionalString(*p.FakeTits)
	}
	if p.Gender != nil && !excluded["gender"] {
		ret.Gender = NewOptionalString(*p.Gender)
	}
	if p.Height != nil && !excluded["height"] {
		h, err := strconv.Atoi(*p.Height)
		if err == nil {
			ret.Height = NewOptionalInt(h)
		}
	}
	if p.Weight != nil && !excluded["weight"] {
		w, err := strconv.Atoi(*p.Weight)
		if err == nil {
			ret.Weight = NewOptionalInt(w)
		}
	}
	if p.Instagram != nil && !excluded["instagram"] {
		ret.Instagram = NewOptionalString(*p.Instagram)
	}
	if p.Measurements != nil && !excluded["measurements"] {
		ret.Measurements = NewOptionalString(*p.Measurements)
	}
	if p.Name != nil && !excluded["name"] {
		ret.Name = NewOptionalString(*p.Name)
	}
	if p.Disambiguation != nil && !excluded["disambiguation"] {
		ret.Disambiguation = NewOptionalString(*p.Disambiguation)
	}
	if p.Details != nil && !excluded["details"] {
		ret.Details = NewOptionalString(*p.Details)
	}
	if p.Piercings != nil && !excluded["piercings"] {
		ret.Piercings = NewOptionalString(*p.Piercings)
	}
	if p.Tattoos != nil && !excluded["tattoos"] {
		ret.Tattoos = NewOptionalString(*p.Tattoos)
	}
	if p.Twitter != nil && !excluded["twitter"] {
		ret.Twitter = NewOptionalString(*p.Twitter)
	}
	if p.URL != nil && !excluded["url"] {
		ret.URL = NewOptionalString(*p.URL)
	}

	if p.RemoteSiteID != nil && endpoint != "" {
		ret.StashIDs = &UpdateStashIDs{
			StashIDs: existingStashIDs,
			Mode:     RelationshipUpdateModeSet,
		}
		ret.StashIDs.Set(StashID{
			Endpoint: endpoint,
			StashID:  *p.RemoteSiteID,
		})
	}

	return ret
}

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
	URLs     []string       `json:"urls"`
	Synopsis *string        `json:"synopsis"`
	Studio   *ScrapedStudio `json:"studio"`
	// This should be a base64 encoded data URL
	FrontImage *string `json:"front_image"`
	// This should be a base64 encoded data URL
	BackImage *string `json:"back_image"`

	// deprecated
	URL *string `json:"url"`
}

func (ScrapedMovie) IsScrapedContent() {}
