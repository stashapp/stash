package scraper

type ScrapedStudio struct {
	// Set if studio matched
	StoredID     *string   `json:"stored_id"`
	Name         string    `json:"name"`
	URLS         []*string `json:"urls"`
	Image        []*string `json:"images"`
	RemoteSiteID *string   `json:"remote_site_id"`
}

func (ScrapedStudio) IsScrapedContent() {}

type ScrapedStudioInput struct {
	Name  *string `json:"name"`
	URL   *string `json:"url"`
	Image *string `json:"image"`
}
