package models

import (
	"context"
	"strconv"
	"time"
)

type Manga struct {
	ID int `json:"id"`

	Title   string `json:"title"`
	URL     string `json:"url"`
	Date    *Date  `json:"date"`
	Details string `json:"details"`
	// Rating expressed in 1-100 scale
	Rating    *int `json:"rating"`
	Organized bool `json:"organized"`
	StudioID  *int `json:"studio_id"`

	// cover image blob checksum
	CoverImageBlob string `json:"cover_image_blob"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	URLs         RelatedStrings `json:"urls"`
	TagIDs       RelatedIDs     `json:"tag_ids"`
	PerformerIDs RelatedIDs     `json:"performer_ids"`
}

func NewManga() Manga {
	currentTime := time.Now()
	return Manga{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}

type CreateMangaInput struct {
	*Manga

	CustomFields map[string]interface{} `json:"custom_fields"`
}

type UpdateMangaInput struct {
	*Manga

	CustomFields CustomFieldsInput `json:"custom_fields"`
}

// MangaPartial represents part of a Manga object. It is used to update
// the database entry. Only non-nil fields will be updated.
type MangaPartial struct {
	Title        OptionalString
	URL          OptionalString
	URLs         *UpdateStrings
	Date         OptionalDate
	Details      OptionalString
	// Rating expressed in 1-100 scale
	Rating    OptionalInt
	Organized OptionalBool
	StudioID  OptionalInt
	CreatedAt OptionalTime
	UpdatedAt OptionalTime

	TagIDs       *UpdateIDs
	PerformerIDs *UpdateIDs

	CoverImageBlob OptionalString

	CustomFields CustomFieldsInput
}

func NewMangaPartial() MangaPartial {
	currentTime := time.Now()
	return MangaPartial{
		UpdatedAt: NewOptionalTime(currentTime),
	}
}

func (m *Manga) LoadURLs(ctx context.Context, l URLLoader) error {
	return m.URLs.load(func() ([]string, error) {
		return l.GetURLs(ctx, m.ID)
	})
}

func (m *Manga) LoadPerformerIDs(ctx context.Context, l PerformerIDLoader) error {
	return m.PerformerIDs.load(func() ([]int, error) {
		return l.GetPerformerIDs(ctx, m.ID)
	})
}

func (m *Manga) LoadTagIDs(ctx context.Context, l TagIDLoader) error {
	return m.TagIDs.load(func() ([]int, error) {
		return l.GetTagIDs(ctx, m.ID)
	})
}

// GetTitle returns the title of the manga.
func (m Manga) GetTitle() string {
	if m.Title != "" {
		return m.Title
	}

	return strconv.Itoa(m.ID)
}

// DisplayName returns a display name for the manga for logging purposes.
func (m Manga) DisplayName() string {
	if m.Title != "" {
		return m.Title
	}

	return strconv.Itoa(m.ID)
}
