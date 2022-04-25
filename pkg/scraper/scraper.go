package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type Source struct {
	// Index of the configured stash-box instance to use. Should be unset if scraper_id is set
	StashBoxIndex *int `json:"stash_box_index"`
	// Stash-box endpoint
	StashBoxEndpoint *string `json:"stash_box_endpoint"`
	// Scraper ID to scrape with. Should be unset if stash_box_index is set
	ScraperID *string `json:"scraper_id"`
}

// Scraped Content is the forming union over the different scrapers
type ScrapedContent interface {
	IsScrapedContent()
}

// Type of the content a scraper generates
type ScrapeContentType string

const (
	ScrapeContentTypeGallery   ScrapeContentType = "GALLERY"
	ScrapeContentTypeMovie     ScrapeContentType = "MOVIE"
	ScrapeContentTypePerformer ScrapeContentType = "PERFORMER"
	ScrapeContentTypeScene     ScrapeContentType = "SCENE"
)

var AllScrapeContentType = []ScrapeContentType{
	ScrapeContentTypeGallery,
	ScrapeContentTypeMovie,
	ScrapeContentTypePerformer,
	ScrapeContentTypeScene,
}

func (e ScrapeContentType) IsValid() bool {
	switch e {
	case ScrapeContentTypeGallery, ScrapeContentTypeMovie, ScrapeContentTypePerformer, ScrapeContentTypeScene:
		return true
	}
	return false
}

func (e ScrapeContentType) String() string {
	return string(e)
}

func (e *ScrapeContentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ScrapeContentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ScrapeContentType", str)
	}
	return nil
}

func (e ScrapeContentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Scraper struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Details for performer scraper
	Performer *ScraperSpec `json:"performer"`
	// Details for scene scraper
	Scene *ScraperSpec `json:"scene"`
	// Details for gallery scraper
	Gallery *ScraperSpec `json:"gallery"`
	// Details for movie scraper
	Movie *ScraperSpec `json:"movie"`
}

type ScraperSpec struct {
	// URLs matching these can be scraped with
	Urls             []string     `json:"urls"`
	SupportedScrapes []ScrapeType `json:"supported_scrapes"`
}

type ScrapeType string

const (
	// From text query
	ScrapeTypeName ScrapeType = "NAME"
	// From existing object
	ScrapeTypeFragment ScrapeType = "FRAGMENT"
	// From URL
	ScrapeTypeURL ScrapeType = "URL"
)

var AllScrapeType = []ScrapeType{
	ScrapeTypeName,
	ScrapeTypeFragment,
	ScrapeTypeURL,
}

func (e ScrapeType) IsValid() bool {
	switch e {
	case ScrapeTypeName, ScrapeTypeFragment, ScrapeTypeURL:
		return true
	}
	return false
}

func (e ScrapeType) String() string {
	return string(e)
}

func (e *ScrapeType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ScrapeType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ScrapeType", str)
	}
	return nil
}

func (e ScrapeType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

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
	Performer *ScrapedPerformerInput
	Scene     *ScrapedSceneInput
	Gallery   *ScrapedGalleryInput
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
	spec() Scraper
	// supports tests if the scraper supports a given content type
	supports(ScrapeContentType) bool
	// supportsURL tests if the scraper supports scrapes of a given url, producing a given content type
	supportsURL(url string, ty ScrapeContentType) bool
}

// urlScraper is the interface of scrapers supporting url loads
type urlScraper interface {
	scraper

	viaURL(ctx context.Context, client *http.Client, url string, ty ScrapeContentType) (ScrapedContent, error)
}

// nameScraper is the interface of scrapers supporting name loads
type nameScraper interface {
	scraper

	viaName(ctx context.Context, client *http.Client, name string, ty ScrapeContentType) ([]ScrapedContent, error)
}

// fragmentScraper is the interface of scrapers supporting fragment loads
type fragmentScraper interface {
	scraper

	viaFragment(ctx context.Context, client *http.Client, input Input) (ScrapedContent, error)
}

// sceneScraper is a scraper which supports scene scrapes with
// scene data as the input.
type sceneScraper interface {
	scraper

	viaScene(ctx context.Context, client *http.Client, scene *models.Scene) (*ScrapedScene, error)
}

// galleryScraper is a scraper which supports gallery scrapes with
// gallery data as the input.
type galleryScraper interface {
	scraper

	viaGallery(ctx context.Context, client *http.Client, gallery *models.Gallery) (*ScrapedGallery, error)
}
