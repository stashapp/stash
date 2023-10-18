package scraper

import "github.com/stashapp/stash/pkg/models"

type ScrapedGallery struct {
	Title      *string                    `json:"title"`
	Details    *string                    `json:"details"`
	URLs       []string                   `json:"urls"`
	Date       *string                    `json:"date"`
	Studio     *models.ScrapedStudio      `json:"studio"`
	Tags       []*models.ScrapedTag       `json:"tags"`
	Performers []*models.ScrapedPerformer `json:"performers"`

	// deprecated
	URL *string `json:"url"`
}

func (ScrapedGallery) IsScrapedContent() {}

type ScrapedGalleryInput struct {
	Title   *string  `json:"title"`
	Details *string  `json:"details"`
	URLs    []string `json:"urls"`
	Date    *string  `json:"date"`

	// deprecated
	URL *string `json:"url"`
}
