package scraper

import (
	"github.com/stashapp/stash/pkg/models"
)

type ScrapedScene struct {
	Title   *string `json:"title"`
	Details *string `json:"details"`
	URL     *string `json:"url"`
	Date    *string `json:"date"`
	// This should be a base64 encoded data URL
	Image        *string                       `json:"image"`
	File         *models.SceneFileType         `json:"file"`
	Studio       *models.ScrapedStudio         `json:"studio"`
	Tags         []*models.ScrapedTag          `json:"tags"`
	Performers   []*models.ScrapedPerformer    `json:"performers"`
	Movies       []*models.ScrapedMovie        `json:"movies"`
	RemoteSiteID *string                       `json:"remote_site_id"`
	Duration     *int                          `json:"duration"`
	Fingerprints []*models.StashBoxFingerprint `json:"fingerprints"`
}

func (ScrapedScene) IsScrapedContent() {}

type ScrapedSceneInput struct {
	Title        *string `json:"title"`
	Details      *string `json:"details"`
	URL          *string `json:"url"`
	Date         *string `json:"date"`
	RemoteSiteID *string `json:"remote_site_id"`
}
