package scraper

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stashapp/stash/pkg/models"
)

// definedScraper implements the scraper interface using a Definition object.
type definedScraper struct {
	config Definition

	globalConf GlobalConfig
}

func scraperFromDefinition(c Definition, globalConfig GlobalConfig) definedScraper {
	return definedScraper{
		config:     c,
		globalConf: globalConfig,
	}
}

func (g definedScraper) spec() Scraper {
	return g.config.spec()
}

// fragmentScraper finds an appropriate fragment scraper based on input.
func (g definedScraper) fragmentScraper(input Input) *ByFragmentDefinition {
	switch {
	case input.Performer != nil:
		return g.config.PerformerByFragment
	case input.Gallery != nil:
		// TODO - this should be galleryByQueryFragment
		return g.config.GalleryByFragment
	case input.Image != nil:
		// TODO - this should be imageByImageFragment
		return g.config.ImageByFragment
	case input.Scene != nil:
		return g.config.SceneByQueryFragment
	}

	return nil
}

func (g definedScraper) viaFragment(ctx context.Context, client *http.Client, input Input) (ScrapedContent, error) {
	stc := g.fragmentScraper(input)
	if stc == nil {
		// If there's no performer fragment scraper in the group, we try to use
		// the URL scraper. Check if there's an URL in the input, and then shift
		// to an URL scrape if it's present.
		if input.Performer != nil && input.Performer.URL != nil && *input.Performer.URL != "" {
			return g.viaURL(ctx, client, *input.Performer.URL, ScrapeContentTypePerformer)
		}

		return nil, ErrNotSupported
	}

	s := g.config.getFragmentScraper(*stc, client, g.globalConf)
	return s.scrapeByFragment(ctx, input)
}

func (g definedScraper) viaScene(ctx context.Context, client *http.Client, scene *models.Scene) (*models.ScrapedScene, error) {
	if g.config.SceneByFragment == nil {
		return nil, ErrNotSupported
	}

	s := g.config.getFragmentScraper(*g.config.SceneByFragment, client, g.globalConf)
	return s.scrapeSceneByScene(ctx, scene)
}

func (g definedScraper) viaGallery(ctx context.Context, client *http.Client, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	if g.config.GalleryByFragment == nil {
		return nil, ErrNotSupported
	}

	s := g.config.getFragmentScraper(*g.config.GalleryByFragment, client, g.globalConf)
	return s.scrapeGalleryByGallery(ctx, gallery)
}

func (g definedScraper) viaImage(ctx context.Context, client *http.Client, gallery *models.Image) (*models.ScrapedImage, error) {
	if g.config.ImageByFragment == nil {
		return nil, ErrNotSupported
	}

	s := g.config.getFragmentScraper(*g.config.ImageByFragment, client, g.globalConf)
	return s.scrapeImageByImage(ctx, gallery)
}

func loadUrlCandidates(c Definition, ty ScrapeContentType) []*ByURLDefinition {
	switch ty {
	case ScrapeContentTypePerformer:
		return c.PerformerByURL
	case ScrapeContentTypeScene:
		return c.SceneByURL
	case ScrapeContentTypeMovie, ScrapeContentTypeGroup:
		return append(c.MovieByURL, c.GroupByURL...)
	case ScrapeContentTypeGallery:
		return c.GalleryByURL
	case ScrapeContentTypeImage:
		return c.ImageByURL
	}

	panic("loadUrlCandidates: unreachable")
}

func (g definedScraper) viaURL(ctx context.Context, client *http.Client, url string, ty ScrapeContentType) (ScrapedContent, error) {
	candidates := loadUrlCandidates(g.config, ty)
	for _, scraper := range candidates {
		if scraper.matchesURL(url) {
			u := replaceURL(url, *scraper) // allow a URL Replace for url-queries
			s := g.config.getURLScraper(*scraper, client, g.globalConf)
			ret, err := s.scrapeByURL(ctx, u, ty)
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

func (g definedScraper) viaName(ctx context.Context, client *http.Client, name string, ty ScrapeContentType) ([]ScrapedContent, error) {
	switch ty {
	case ScrapeContentTypePerformer:
		if g.config.PerformerByName == nil {
			break
		}

		s := g.config.getNameScraper(*g.config.PerformerByName, client, g.globalConf)
		return s.scrapeByName(ctx, name, ty)
	case ScrapeContentTypeScene:
		if g.config.SceneByName == nil {
			break
		}

		s := g.config.getNameScraper(*g.config.SceneByName, client, g.globalConf)
		return s.scrapeByName(ctx, name, ty)
	}

	return nil, fmt.Errorf("%w: cannot load %v by name", ErrNotSupported, ty)
}

func (g definedScraper) supports(ty ScrapeContentType) bool {
	return g.config.supports(ty)
}

func (g definedScraper) supportsURL(url string, ty ScrapeContentType) bool {
	return g.config.matchesURL(url, ty)
}
