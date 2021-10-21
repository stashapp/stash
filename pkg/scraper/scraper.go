package scraper

import (
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

var ErrNotSupported = errors.New("not supported")

type urlMatcher interface {
	matchesURL(url string) bool
}

type Input struct {
	Performer *models.ScrapedPerformerInput
	Scene     *models.ScrapedSceneInput
	Gallery   *models.ScrapedGalleryInput
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

// scraper is the generic interface to the scraper subsystems
type scraper interface {
	// spec returns the scraper specification, suitable for graphql
	spec() models.Scraper
	// supports tests if the scraper supports a given content type
	supports(models.ScrapeContentType) bool
	// supportsURL tests if the scraper supports scrapes of a given url, producing a given content type
	supportsURL(url string, ty models.ScrapeContentType) bool
}

// urlScraper is the interface of scrapers supporting url loads
type urlScraper interface {
	scraper

	loadByURL(url string, ty models.ScrapeContentType) (models.ScrapedContent, error)
}

// nameScraper is the interface of scrapers supporting name loads
type nameScraper interface {
	scraper

	loadByName(name string, ty models.ScrapeContentType) ([]models.ScrapedContent, error)
}

// fragmentScraper is the interface of scrapers supporting fragment loads
type fragmentScraper interface {
	scraper

	loadByFragment(input Input) (models.ScrapedContent, error)
}

type sceneLoader interface {
	scraper

	loadByScene(scene *models.Scene) (*models.ScrapedScene, error)
}

type galleryLoader interface {
	scraper

	loadByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error)
}

type scraper_s struct {
	Spec *models.Scraper

	performer performerScraper
	scene     sceneScraper
	gallery   galleryScraper
	movie     movieScraper
}

func (s scraper_s) spec() models.Scraper {
	return *s.Spec
}

func (s scraper_s) loadByURL(url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	switch ty {
	case models.ScrapeContentTypePerformer:
		return s.performer.scrapeByURL(url)
	case models.ScrapeContentTypeScene:
		return s.scene.scrapeByURL(url)
	case models.ScrapeContentTypeGallery:
		return s.gallery.scrapeByURL(url)
	case models.ScrapeContentTypeMovie:
		return s.movie.scrapeByURL(url)
	default:
		panic("Unimplemented scraper type")
	}
}

func (s scraper_s) loadByName(name string, ty models.ScrapeContentType) ([]models.ScrapedContent, error) {
	switch ty {
	case models.ScrapeContentTypePerformer:
		performers, err := s.performer.scrapeByName(name)
		if err != nil {
			return nil, err
		}
		content := make([]models.ScrapedContent, len(performers))
		for i := range performers {
			content[i] = performers[i]
		}
		return content, nil
	case models.ScrapeContentTypeScene:
		scenes, err := s.scene.scrapeByName(name)
		if err != nil {
			return nil, err
		}
		content := make([]models.ScrapedContent, len(scenes))
		for i := range scenes {
			content[i] = scenes[i]
		}
		return content, nil
	default:
		return nil, fmt.Errorf("loading %v by name: %w", ty, ErrUnsupported)
	}
}

func (s scraper_s) supports(ty models.ScrapeContentType) bool {
	return s.matchesContentType(ty)
}

func (s scraper_s) matchesContentType(k models.ScrapeContentType) bool {
	switch k {
	case models.ScrapeContentTypePerformer:
		return s.performer != nil
	case models.ScrapeContentTypeScene:
		return s.scene != nil
	case models.ScrapeContentTypeGallery:
		return s.gallery != nil
	case models.ScrapeContentTypeMovie:
		return s.movie != nil
	default:
		return false
	}
}

func (s scraper_s) supportsURL(url string, ty models.ScrapeContentType) bool {
	return s.matchesURLContent(url, ty)
}

func (s scraper_s) matchesURLContent(url string, k models.ScrapeContentType) bool {
	var matcher urlMatcher
	var ok bool
	switch k {
	case models.ScrapeContentTypePerformer:
		matcher, ok = s.performer.(urlMatcher)
	case models.ScrapeContentTypeScene:
		matcher, ok = s.scene.(urlMatcher)
	case models.ScrapeContentTypeGallery:
		matcher, ok = s.gallery.(urlMatcher)
	case models.ScrapeContentTypeMovie:
		matcher, ok = s.movie.(urlMatcher)
	}

	if !ok {
		return false
	}

	return matcher.matchesURL(url)
}
