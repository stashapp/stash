package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
)

func (r *queryResolver) ScrapeFreeones(ctx context.Context, performer_name string) (*models.ScrapedPerformer, error) {
	return scraper.GetPerformer(performer_name)
}

func (r *queryResolver) ScrapeFreeonesPerformerList(ctx context.Context, query string) ([]string, error) {
	return scraper.GetPerformerNames(query)
}

func (r *queryResolver) ListScrapers(ctx context.Context, scraperType models.ScraperType) ([]*models.Scraper, error) {
	return scraper.ListScrapers(scraperType)
}

func (r *queryResolver) ScrapePerformerList(ctx context.Context, scraperID string, query string) ([]string, error) {
	return scraper.ScrapePerformerList(scraperID, query)
}

func (r *queryResolver) ScrapePerformer(ctx context.Context, scraperID string, performerName string) (*models.ScrapedPerformer, error) {
	return scraper.ScrapePerformer(scraperID, performerName)
}
