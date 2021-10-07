package scraper

import "github.com/stashapp/stash/pkg/models"

type urlMatcher interface {
	matchesURL(url string) bool
}

type performerScraper interface {
	scrapeByName(name string) ([]*models.ScrapedPerformer, error)
	scrapeByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error)
	scrapeByURL(url string) (*models.ScrapedPerformer, error)
}

type performerScraperMatcher interface {
	performerScraper
	urlMatcher
}

type sceneScraper interface {
	scrapeByName(name string) ([]*models.ScrapedScene, error)
	scrapeByScene(scene *models.Scene) (*models.ScrapedScene, error)
	scrapeByFragment(scene models.ScrapedSceneInput) (*models.ScrapedScene, error)
	scrapeByURL(url string) (*models.ScrapedScene, error)
}

type sceneScraperMatcher interface {
	sceneScraper
	urlMatcher
}

type galleryScraper interface {
	scrapeByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error)
	scrapeByFragment(gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error)
	scrapeByURL(url string) (*models.ScrapedGallery, error)
}

type galleryScraperMatcher interface {
	galleryScraper
	urlMatcher
}

type movieScraper interface {
	scrapeByURL(url string) (*models.ScrapedMovie, error)
}

type movieScraperMatcher interface {
	movieScraper
	urlMatcher
}

type scraperImpl interface {
	ID() string
	ScraperSpec() *models.Scraper

	Performer() performerScraperMatcher
	Scene() sceneScraperMatcher
	Gallery() galleryScraperMatcher
	Movie() movieScraperMatcher
}
