package scraper

type ScrapedMovieInput struct {
	Name     *string `json:"name"`
	Aliases  *string `json:"aliases"`
	Duration *string `json:"duration"`
	Date     *string `json:"date"`
	Rating   *string `json:"rating"`
	Director *string `json:"director"`
	URL      *string `json:"url"`
	Synopsis *string `json:"synopsis"`
}
