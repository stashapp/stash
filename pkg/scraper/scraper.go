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

type sceneScraper interface {
	scrapeByName(name string) ([]*models.ScrapedScene, error)
	scrapeByScene(scene *models.Scene) (*models.ScrapedScene, error)
	scrapeByFragment(scene models.ScrapedSceneInput) (*models.ScrapedScene, error)
	scrapeByURL(url string) (*models.ScrapedScene, error)
}

type galleryScraper interface {
	scrapeByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error)
	scrapeByFragment(gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error)
	scrapeByURL(url string) (*models.ScrapedGallery, error)
}

type movieScraper interface {
	scrapeByURL(url string) (*models.ScrapedMovie, error)
}

type scraper struct {
	ID   string
	Spec *models.Scraper

	Performer performerScraper
	Scene     sceneScraper
	Gallery   galleryScraper
	Movie     movieScraper
}

func matchesURL(maybeURLMatcher interface{}, url string) bool {
	if maybeURLMatcher != nil {
		matcher, ok := maybeURLMatcher.(urlMatcher)
		if ok {
			return matcher.matchesURL(url)
		}
	}

	return false
}
