package models

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

type ScrapedStudio struct {
	// Set if studio matched
	StoredID     *string        `json:"stored_id"`
	Name         string         `json:"name"`
	URL          *string        `json:"url"` // deprecated
	URLs         []string       `json:"urls"`
	Parent       *ScrapedStudio `json:"parent"`
	Image        *string        `json:"image"`
	Images       []string       `json:"images"`
	Details      *string        `json:"details"`
	Aliases      *string        `json:"aliases"`
	Tags         []*ScrapedTag  `json:"tags"`
	RemoteSiteID *string        `json:"remote_site_id"`
}

func (ScrapedStudio) IsScrapedContent() {}

func (s *ScrapedStudio) ToStudio(endpoint string, excluded map[string]bool) *Studio {
	// Populate a new studio from the input
	ret := NewStudio()
	ret.Name = strings.TrimSpace(s.Name)

	if s.RemoteSiteID != nil && endpoint != "" && *s.RemoteSiteID != "" {
		ret.StashIDs = NewRelatedStashIDs([]StashID{
			{
				Endpoint:  endpoint,
				StashID:   *s.RemoteSiteID,
				UpdatedAt: time.Now(),
			},
		})
	}

	// if URLs are provided, only use those
	if len(s.URLs) > 0 {
		if !excluded["urls"] {
			ret.URLs = NewRelatedStrings(s.URLs)
		}
	} else {
		urls := []string{}
		if s.URL != nil && !excluded["url"] {
			urls = append(urls, *s.URL)
		}

		if len(urls) > 0 {
			ret.URLs = NewRelatedStrings(urls)
		}
	}

	if s.Details != nil && !excluded["details"] {
		ret.Details = *s.Details
	}

	if s.Aliases != nil && *s.Aliases != "" && !excluded["aliases"] {
		ret.Aliases = NewRelatedStrings(stringslice.FromString(*s.Aliases, ","))
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

func (s *ScrapedStudio) ToPartial(id string, endpoint string, excluded map[string]bool, existingStashIDs []StashID) StudioPartial {
	ret := NewStudioPartial()
	ret.ID, _ = strconv.Atoi(id)
	currentTime := time.Now()

	if s.Name != "" && !excluded["name"] {
		ret.Name = NewOptionalString(strings.TrimSpace(s.Name))
	}

	if len(s.URLs) > 0 {
		if !excluded["urls"] {

			ret.URLs = &UpdateStrings{
				Values: stringslice.TrimSpace(s.URLs),
				Mode:   RelationshipUpdateModeSet,
			}
		}
	} else {
		urls := []string{}
		if s.URL != nil && !excluded["url"] {
			urls = append(urls, strings.TrimSpace(*s.URL))
		}

		if len(urls) > 0 {
			ret.URLs = &UpdateStrings{
				Values: stringslice.TrimSpace(urls),
				Mode:   RelationshipUpdateModeSet,
			}
		}
	}

	if s.Details != nil && !excluded["details"] {
		ret.Details = NewOptionalString(strings.TrimSpace(*s.Details))
	}

	if s.Aliases != nil && *s.Aliases != "" && !excluded["aliases"] {
		ret.Aliases = &UpdateStrings{
			Values: stringslice.TrimSpace(stringslice.FromString(*s.Aliases, ",")),
			Mode:   RelationshipUpdateModeSet,
		}
	}

	if s.Parent != nil && !excluded["parent"] {
		if s.Parent.StoredID != nil {
			parentID, _ := strconv.Atoi(*s.Parent.StoredID)
			if parentID > 0 {
				// This is to be set directly as we know it has a value and the translator won't have the field
				ret.ParentID = NewOptionalInt(parentID)
			}
		}
	}

	if s.RemoteSiteID != nil && endpoint != "" && *s.RemoteSiteID != "" {
		ret.StashIDs = &UpdateStashIDs{
			StashIDs: existingStashIDs,
			Mode:     RelationshipUpdateModeSet,
		}
		ret.StashIDs.Set(StashID{
			Endpoint:  endpoint,
			StashID:   *s.RemoteSiteID,
			UpdatedAt: currentTime,
		})
	}

	return ret
}

// A performer from a scraping operation...
type ScrapedPerformer struct {
	// Set if performer matched
	StoredID       *string       `json:"stored_id"`
	Name           *string       `json:"name"`
	Disambiguation *string       `json:"disambiguation"`
	Gender         *string       `json:"gender"`
	URLs           []string      `json:"urls"`
	URL            *string       `json:"url"`       // deprecated
	Twitter        *string       `json:"twitter"`   // deprecated
	Instagram      *string       `json:"instagram"` // deprecated
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
	Image              *string  `json:"image"` // deprecated: use Images
	Images             []string `json:"images"`
	Details            *string  `json:"details"`
	DeathDate          *string  `json:"death_date"`
	HairColor          *string  `json:"hair_color"`
	Weight             *string  `json:"weight"`
	RemoteSiteID       *string  `json:"remote_site_id"`
	RemoteDeleted      bool     `json:"remote_deleted"`
	RemoteMergedIntoId *string  `json:"remote_merged_into_id"`
}

func (ScrapedPerformer) IsScrapedContent() {}

func (p *ScrapedPerformer) ToPerformer(endpoint string, excluded map[string]bool) *Performer {
	ret := NewPerformer()
	currentTime := time.Now()
	ret.Name = strings.TrimSpace(*p.Name)

	if p.Aliases != nil && !excluded["aliases"] {
		aliases := stringslice.FromString(*p.Aliases, ",")
		for i, alias := range aliases {
			aliases[i] = strings.TrimSpace(alias)
		}
		ret.Aliases = NewRelatedStrings(aliases)
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

	// if URLs are provided, only use those
	if len(p.URLs) > 0 {
		if !excluded["urls"] {
			ret.URLs = NewRelatedStrings(p.URLs)
		}
	} else {
		urls := []string{}
		if p.URL != nil && !excluded["url"] {
			urls = append(urls, *p.URL)
		}
		if p.Twitter != nil && !excluded["twitter"] {
			urls = append(urls, *p.Twitter)
		}
		if p.Instagram != nil && !excluded["instagram"] {
			urls = append(urls, *p.Instagram)
		}

		if len(urls) > 0 {
			ret.URLs = NewRelatedStrings(urls)
		}
	}

	if p.RemoteSiteID != nil && endpoint != "" && *p.RemoteSiteID != "" {
		ret.StashIDs = NewRelatedStashIDs([]StashID{
			{
				Endpoint:  endpoint,
				StashID:   *p.RemoteSiteID,
				UpdatedAt: currentTime,
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

	// if URLs are provided, only use those
	if len(p.URLs) > 0 {
		if !excluded["urls"] {
			ret.URLs = &UpdateStrings{
				Values: p.URLs,
				Mode:   RelationshipUpdateModeSet,
			}
		}
	} else {
		urls := []string{}
		if p.URL != nil && !excluded["url"] {
			urls = append(urls, *p.URL)
		}
		if p.Twitter != nil && !excluded["twitter"] {
			urls = append(urls, *p.Twitter)
		}
		if p.Instagram != nil && !excluded["instagram"] {
			urls = append(urls, *p.Instagram)
		}

		if len(urls) > 0 {
			ret.URLs = &UpdateStrings{
				Values: urls,
				Mode:   RelationshipUpdateModeSet,
			}
		}
	}

	if p.RemoteSiteID != nil && endpoint != "" && *p.RemoteSiteID != "" {
		ret.StashIDs = &UpdateStashIDs{
			StashIDs: existingStashIDs,
			Mode:     RelationshipUpdateModeSet,
		}
		ret.StashIDs.Set(StashID{
			Endpoint:  endpoint,
			StashID:   *p.RemoteSiteID,
			UpdatedAt: time.Now(),
		})
	}

	return ret
}

type ScrapedTag struct {
	// Set if tag matched
	StoredID     *string `json:"stored_id"`
	Name         string  `json:"name"`
	RemoteSiteID *string `json:"remote_site_id"`
}

func (ScrapedTag) IsScrapedContent() {}

func (t *ScrapedTag) ToTag(endpoint string, excluded map[string]bool) *Tag {
	currentTime := time.Now()
	ret := NewTag()
	ret.Name = t.Name

	if t.RemoteSiteID != nil && endpoint != "" && *t.RemoteSiteID != "" {
		ret.StashIDs = NewRelatedStashIDs([]StashID{
			{
				Endpoint:  endpoint,
				StashID:   *t.RemoteSiteID,
				UpdatedAt: currentTime,
			},
		})
	}

	return &ret
}

func ScrapedTagSortFunction(a, b *ScrapedTag) int {
	return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
}

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
	Tags     []*ScrapedTag  `json:"tags"`
	// This should be a base64 encoded data URL
	FrontImage *string `json:"front_image"`
	// This should be a base64 encoded data URL
	BackImage *string `json:"back_image"`

	// deprecated
	URL *string `json:"url"`
}

func (ScrapedMovie) IsScrapedContent() {}

func (m ScrapedMovie) ScrapedGroup() ScrapedGroup {
	ret := ScrapedGroup{
		StoredID:   m.StoredID,
		Name:       m.Name,
		Aliases:    m.Aliases,
		Duration:   m.Duration,
		Date:       m.Date,
		Rating:     m.Rating,
		Director:   m.Director,
		URLs:       m.URLs,
		Synopsis:   m.Synopsis,
		Studio:     m.Studio,
		Tags:       m.Tags,
		FrontImage: m.FrontImage,
		BackImage:  m.BackImage,
	}

	if len(m.URLs) == 0 && m.URL != nil {
		ret.URLs = []string{*m.URL}
	}

	return ret
}

// ScrapedGroup is a group from a scraping operation
type ScrapedGroup struct {
	StoredID *string        `json:"stored_id"`
	Name     *string        `json:"name"`
	Aliases  *string        `json:"aliases"`
	Duration *string        `json:"duration"`
	Date     *string        `json:"date"`
	Rating   *string        `json:"rating"`
	Director *string        `json:"director"`
	URL      *string        `json:"url"` // included for backward compatibility
	URLs     []string       `json:"urls"`
	Synopsis *string        `json:"synopsis"`
	Studio   *ScrapedStudio `json:"studio"`
	Tags     []*ScrapedTag  `json:"tags"`
	// This should be a base64 encoded data URL
	FrontImage *string `json:"front_image"`
	// This should be a base64 encoded data URL
	BackImage *string `json:"back_image"`
}

func (ScrapedGroup) IsScrapedContent() {}

func (g ScrapedGroup) ScrapedMovie() ScrapedMovie {
	ret := ScrapedMovie{
		StoredID:   g.StoredID,
		Name:       g.Name,
		Aliases:    g.Aliases,
		Duration:   g.Duration,
		Date:       g.Date,
		Rating:     g.Rating,
		Director:   g.Director,
		URLs:       g.URLs,
		Synopsis:   g.Synopsis,
		Studio:     g.Studio,
		Tags:       g.Tags,
		FrontImage: g.FrontImage,
		BackImage:  g.BackImage,
	}

	if len(g.URLs) > 0 {
		ret.URL = &g.URLs[0]
	}

	return ret
}

type ScrapedScene struct {
	Title    *string  `json:"title"`
	Code     *string  `json:"code"`
	Details  *string  `json:"details"`
	Director *string  `json:"director"`
	URL      *string  `json:"url"`
	URLs     []string `json:"urls"`
	Date     *string  `json:"date"`
	// This should be a base64 encoded data URL
	Image        *string                `json:"image"`
	File         *SceneFileType         `json:"file"`
	Studio       *ScrapedStudio         `json:"studio"`
	Tags         []*ScrapedTag          `json:"tags"`
	Performers   []*ScrapedPerformer    `json:"performers"`
	Groups       []*ScrapedGroup        `json:"groups"`
	Movies       []*ScrapedMovie        `json:"movies"`
	RemoteSiteID *string                `json:"remote_site_id"`
	Duration     *int                   `json:"duration"`
	Fingerprints []*StashBoxFingerprint `json:"fingerprints"`
}

func (ScrapedScene) IsScrapedContent() {}

type ScrapedSceneInput struct {
	Title        *string  `json:"title"`
	Code         *string  `json:"code"`
	Details      *string  `json:"details"`
	Director     *string  `json:"director"`
	URL          *string  `json:"url"`
	URLs         []string `json:"urls"`
	Date         *string  `json:"date"`
	RemoteSiteID *string  `json:"remote_site_id"`
}

type ScrapedImage struct {
	Title        *string             `json:"title"`
	Code         *string             `json:"code"`
	Details      *string             `json:"details"`
	Photographer *string             `json:"photographer"`
	URLs         []string            `json:"urls"`
	Date         *string             `json:"date"`
	Studio       *ScrapedStudio      `json:"studio"`
	Tags         []*ScrapedTag       `json:"tags"`
	Performers   []*ScrapedPerformer `json:"performers"`
}

func (ScrapedImage) IsScrapedContent() {}

type ScrapedImageInput struct {
	Title   *string  `json:"title"`
	Code    *string  `json:"code"`
	Details *string  `json:"details"`
	URLs    []string `json:"urls"`
	Date    *string  `json:"date"`
}

type ScrapedGallery struct {
	Title        *string             `json:"title"`
	Code         *string             `json:"code"`
	Details      *string             `json:"details"`
	Photographer *string             `json:"photographer"`
	URLs         []string            `json:"urls"`
	Date         *string             `json:"date"`
	Studio       *ScrapedStudio      `json:"studio"`
	Tags         []*ScrapedTag       `json:"tags"`
	Performers   []*ScrapedPerformer `json:"performers"`

	// deprecated
	URL *string `json:"url"`
}

func (ScrapedGallery) IsScrapedContent() {}

type ScrapedGalleryInput struct {
	Title        *string  `json:"title"`
	Code         *string  `json:"code"`
	Details      *string  `json:"details"`
	Photographer *string  `json:"photographer"`
	URLs         []string `json:"urls"`
	Date         *string  `json:"date"`

	// deprecated
	URL *string `json:"url"`
}
