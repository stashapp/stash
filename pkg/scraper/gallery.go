package scraper

import "github.com/stashapp/stash/pkg/models"

type ScrapedGallery struct {
	Title        *string                    `json:"title"`
	Code         *string                    `json:"code"`
	Details      *string                    `json:"details"`
	Photographer *string                    `json:"photographer"`
	URLs         []string                   `json:"urls"`
	Date         *string                    `json:"date"`
	Studio       *models.ScrapedStudio      `json:"studio"`
	Tags         []*models.ScrapedTag       `json:"tags"`
	Performers   []*models.ScrapedPerformer `json:"performers"`

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
