package scraper

type ScrapedMovieInput struct {
	Name     *string  `json:"name"`
	Aliases  *string  `json:"aliases"`
	Duration *string  `json:"duration"`
	Date     *string  `json:"date"`
	Rating   *string  `json:"rating"`
	Director *string  `json:"director"`
	URLs     []string `json:"urls"`
	Synopsis *string  `json:"synopsis"`

	// deprecated
	URL *string `json:"url"`
}
