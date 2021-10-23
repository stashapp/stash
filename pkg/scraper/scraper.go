package scraper

import (
	"fmt"
	"net/http"

	"github.com/stashapp/stash/pkg/models"
)

// Input coalesces inputs of diffrent types into a single structure.
// The system expects one of these to be set, and the remaining to be
// set to nil.
type Input struct {
	Performer *models.ScrapedPerformerInput
	Scene     *models.ScrapedSceneInput
	Gallery   *models.ScrapedGalleryInput
}

type performerScraper interface {
	scrapeByName(name string) ([]*models.ScrapedPerformer, error)
	scrapeByFragment(scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error)
}

type sceneScraper interface {
	scrapeByName(name string) ([]*models.ScrapedScene, error)
	scrapeByScene(scene *models.Scene) (*models.ScrapedScene, error)
	scrapeByFragment(scene models.ScrapedSceneInput) (*models.ScrapedScene, error)
}

type galleryScraper interface {
	scrapeByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error)
	scrapeByFragment(gallery models.ScrapedGalleryInput) (*models.ScrapedGallery, error)
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

	loadByURL(client *http.Client, url string, ty models.ScrapeContentType) (models.ScrapedContent, error)
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

// sceneLoader is a scraper which supports scene scrapes with
// scene data as the input.
type sceneLoader interface {
	scraper

	loadByScene(scene *models.Scene) (*models.ScrapedScene, error)
}

// galleryLoader is a sraper which supports gallery scrapes with
// gallery data as the input.
type galleryLoader interface {
	scraper

	loadByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error)
}

type group struct {
	config        config
	specification *models.Scraper

	txnManager models.TransactionManager
	globalConf GlobalConfig

	performer performerScraper
	scene     sceneScraper
	gallery   galleryScraper
}

func (g group) spec() models.Scraper {
	return *g.specification
}

func (g group) loadByFragment(input Input) (models.ScrapedContent, error) {
	switch {
	case input.Performer != nil:
		g.performer.scrapeByFragment(*input.Performer)
	case input.Gallery != nil:
		g.gallery.scrapeByFragment(*input.Gallery)
	case input.Scene != nil:
		g.scene.scrapeByFragment(*input.Scene)
	}

	return nil, ErrNotSupported
}

func (g group) loadByScene(scene *models.Scene) (*models.ScrapedScene, error) {
	if g.scene == nil {
		return nil, ErrNotSupported
	}

	return g.scene.scrapeByScene(scene)
}

func (g group) loadByGallery(gallery *models.Gallery) (*models.ScrapedGallery, error) {
	if g.gallery == nil {
		return nil, ErrNotSupported
	}

	return g.gallery.scrapeByGallery(gallery)
}

func loadUrlCandidates(c config, ty models.ScrapeContentType) []*scrapeByURLConfig {
	switch ty {
	case models.ScrapeContentTypePerformer:
		return c.PerformerByURL
	case models.ScrapeContentTypeScene:
		return c.SceneByURL
	case models.ScrapeContentTypeMovie:
		return c.MovieByURL
	case models.ScrapeContentTypeGallery:
		return c.GalleryByURL
	}

	panic("loadUrlCandidates: unreachable")
}

func scrapeByUrl(url string, s scraperActionImpl, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	switch ty {
	case models.ScrapeContentTypePerformer:
		return s.scrapePerformerByURL(url)
	case models.ScrapeContentTypeScene:
		return s.scrapeSceneByURL(url)
	case models.ScrapeContentTypeMovie:
		return s.scrapeMovieByURL(url)
	case models.ScrapeContentTypeGallery:
		return s.scrapeGalleryByURL(url)
	}

	panic("scrapeByUrl: unreachable")
}

func (g group) loadByURL(client *http.Client, url string, ty models.ScrapeContentType) (models.ScrapedContent, error) {
	candidates := loadUrlCandidates(g.config, ty)
	for _, scraper := range candidates {
		if scraper.matchesURL(url) {
			s := g.config.getScraper(scraper.scraperTypeConfig, client, g.txnManager, g.globalConf)
			ret, err := scrapeByUrl(url, s, ty)
			if err != nil {
				return nil, err
			}

			if ret != nil {
				return ret, nil
			}
		}
	}

	return nil, nil
}

func (s group) loadByName(name string, ty models.ScrapeContentType) ([]models.ScrapedContent, error) {
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
		return nil, fmt.Errorf("loading %v by name: %w", ty, ErrNotSupported)
	}
}

func (s group) supports(ty models.ScrapeContentType) bool {
	return s.config.supports(ty)
}

func (s group) supportsURL(url string, ty models.ScrapeContentType) bool {
	return s.config.matchesURL(url, ty)
}
