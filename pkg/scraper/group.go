package scraper

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stashapp/stash/pkg/models"
)

type group struct {
	config config

	txnManager models.TransactionManager
	globalConf GlobalConfig
}

func newGroupScraper(c config, txnManager models.TransactionManager, globalConfig GlobalConfig) scraper {
	return group{
		config:     c,
		txnManager: txnManager,
		globalConf: globalConfig,
	}
}

func (g group) spec() Scraper {
	return g.config.spec()
}

// fragmentScraper finds an appropriate fragment scraper based on input.
func (g group) fragmentScraper(input Input) *scraperTypeConfig {
	switch {
	case input.Performer != nil:
		return g.config.PerformerByFragment
	case input.Gallery != nil:
		// TODO - this should be galleryByQueryFragment
		return g.config.GalleryByFragment
	case input.Scene != nil:
		return g.config.SceneByQueryFragment
	}

	return nil
}

func (g group) viaFragment(ctx context.Context, client *http.Client, input Input) (ScrapedContent, error) {
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

	s := g.config.getScraper(*stc, client, g.txnManager, g.globalConf)
	return s.scrapeByFragment(ctx, input)
}

func (g group) viaScene(ctx context.Context, client *http.Client, scene *models.Scene) (*ScrapedScene, error) {
	if g.config.SceneByFragment == nil {
		return nil, ErrNotSupported
	}

	s := g.config.getScraper(*g.config.SceneByFragment, client, g.txnManager, g.globalConf)
	return s.scrapeSceneByScene(ctx, scene)
}

func (g group) viaGallery(ctx context.Context, client *http.Client, gallery *models.Gallery) (*ScrapedGallery, error) {
	if g.config.GalleryByFragment == nil {
		return nil, ErrNotSupported
	}

	s := g.config.getScraper(*g.config.GalleryByFragment, client, g.txnManager, g.globalConf)
	return s.scrapeGalleryByGallery(ctx, gallery)
}

func loadUrlCandidates(c config, ty ScrapeContentType) []*scrapeByURLConfig {
	switch ty {
	case ScrapeContentTypePerformer:
		return c.PerformerByURL
	case ScrapeContentTypeScene:
		return c.SceneByURL
	case ScrapeContentTypeMovie:
		return c.MovieByURL
	case ScrapeContentTypeGallery:
		return c.GalleryByURL
	}

	panic("loadUrlCandidates: unreachable")
}

func (g group) viaURL(ctx context.Context, client *http.Client, url string, ty ScrapeContentType) (ScrapedContent, error) {
	candidates := loadUrlCandidates(g.config, ty)
	for _, scraper := range candidates {
		if scraper.matchesURL(url) {
			s := g.config.getScraper(scraper.scraperTypeConfig, client, g.txnManager, g.globalConf)
			ret, err := s.scrapeByURL(ctx, url, ty)
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

func (g group) viaName(ctx context.Context, client *http.Client, name string, ty ScrapeContentType) ([]ScrapedContent, error) {
	switch ty {
	case ScrapeContentTypePerformer:
		if g.config.PerformerByName == nil {
			break
		}

		s := g.config.getScraper(*g.config.PerformerByName, client, g.txnManager, g.globalConf)
		return s.scrapeByName(ctx, name, ty)
	case ScrapeContentTypeScene:
		if g.config.SceneByName == nil {
			break
		}

		s := g.config.getScraper(*g.config.SceneByName, client, g.txnManager, g.globalConf)
		return s.scrapeByName(ctx, name, ty)
	}

	return nil, fmt.Errorf("%w: cannot load %v by name", ErrNotSupported, ty)
}

func (g group) supports(ty ScrapeContentType) bool {
	return g.config.supports(ty)
}

func (g group) supportsURL(url string, ty ScrapeContentType) bool {
	return g.config.matchesURL(url, ty)
}
