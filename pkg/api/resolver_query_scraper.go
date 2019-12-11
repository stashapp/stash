package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
)

// deprecated
func (r *queryResolver) ScrapeFreeones(ctx context.Context, performer_name string) (*models.ScrapedPerformer, error) {
	scrapedPerformer := models.ScrapedPerformerInput{
		Name: &performer_name,
	}
	return scraper.GetFreeonesScraper().ScrapePerformer(scrapedPerformer)
}

// deprecated
func (r *queryResolver) ScrapeFreeonesPerformerList(ctx context.Context, query string) ([]string, error) {
	scrapedPerformers, err := scraper.GetFreeonesScraper().ScrapePerformerNames(query)

	if err != nil {
		return nil, err
	}

	var ret []string
	for _, v := range scrapedPerformers {
		name := v.Name
		ret = append(ret, *name)
	}

	return ret, nil
}

func (r *queryResolver) ListPerformerScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return scraper.ListPerformerScrapers()
}

func (r *queryResolver) ScrapePerformerList(ctx context.Context, scraperID string, query string) ([]*models.ScrapedPerformer, error) {
	if query == "" {
		return nil, nil
	}

	return scraper.ScrapePerformerList(scraperID, query)
}

func (r *queryResolver) ScrapePerformer(ctx context.Context, scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	return scraper.ScrapePerformer(scraperID, scrapedPerformer)
}

func (r *queryResolver) ScrapePerformerURL(ctx context.Context, url string) (*models.ScrapedPerformer, error) {
	return scraper.ScrapePerformerURL(url)
}
