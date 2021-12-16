package scraper

import (
	"context"
	"errors"
	"net/http"

	"github.com/stashapp/stash/pkg/models"
)

var (
	// ErrMaxRedirects is returned if the max number of HTTP redirects are reached.
	ErrMaxRedirects = errors.New("maximum number of HTTP redirects reached")

	// ErrNotFound is returned when an entity isn't found
	ErrNotFound = errors.New("scraper not found")

	// ErrNotSupported is returned when a given invocation isn't supported, and there
	// is a guard function which should be able to guard against it.
	ErrNotSupported = errors.New("scraper operation not supported")
)

// Input coalesces inputs of different types into a single structure.
// The system expects one of these to be set, and the remaining to be
// set to nil.
type Input struct {
	Performer *models.ScrapedPerformerInput
	Scene     *models.ScrapedSceneInput
	Gallery   *models.ScrapedGalleryInput
}

// simple type definitions that can help customize
// actions per query
type QueryType int

const (
	// for now only SearchQuery is needed
	SearchQuery QueryType = iota + 1
)

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

	viaURL(ctx context.Context, client *http.Client, url string, ty models.ScrapeContentType) (models.ScrapedContent, error)
}

// nameScraper is the interface of scrapers supporting name loads
type nameScraper interface {
	scraper

	viaName(ctx context.Context, client *http.Client, name string, ty models.ScrapeContentType) ([]models.ScrapedContent, error)
}

// fragmentScraper is the interface of scrapers supporting fragment loads
type fragmentScraper interface {
	scraper

	viaFragment(ctx context.Context, client *http.Client, input Input) (models.ScrapedContent, error)
}

// sceneScraper is a scraper which supports scene scrapes with
// scene data as the input.
type sceneScraper interface {
	scraper

	viaScene(ctx context.Context, client *http.Client, scene *models.Scene) (*models.ScrapedScene, error)
}

// galleryScraper is a scraper which supports gallery scrapes with
// gallery data as the input.
type galleryScraper interface {
	scraper

	viaGallery(ctx context.Context, client *http.Client, gallery *models.Gallery) (*models.ScrapedGallery, error)
}
