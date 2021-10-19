package scraper

import "github.com/stashapp/stash/pkg/models"

// Kind is the categorization of scrapers
type Kind int

const (
	// Unknown scraper
	Unknown Kind = iota
	// The scraper can handle Performer scrapes
	Performer
	// The scraper can handle Scene scrapes
	Scene
	// The scraper can handle Gallery scrapes
	Gallery
	// The scraper can handle Movie scrapes
	Movie
)

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

func (k Kind) String() string {
	switch k {
	case Unknown:
		return "Unknown"
	case Performer:
		return "Performer"
	case Scene:
		return "Scene"
	case Gallery:
		return "Gallery"
	case Movie:
		return "Movie"
	}

	panic("missing implementation of Kind.String()")
}

func (s scraper) matchKind(k Kind) bool {
	switch k {
	case Unknown:
		return false
	case Performer:
		return s.Performer != nil
	case Scene:
		return s.Scene != nil
	case Gallery:
		return s.Gallery != nil
	case Movie:
		return s.Movie != nil
	default:
		return false
	}
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
