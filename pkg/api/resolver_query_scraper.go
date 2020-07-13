package api

import (
	"context"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
)

// deprecated
func (r *queryResolver) ScrapeFreeones(ctx context.Context, performer_name string) (*models.ScrapedPerformer, error) {
	scrapedPerformer := models.ScrapedPerformerInput{
		Name: &performer_name,
	}
	return manager.GetInstance().ScraperCache.ScrapePerformer(scraper.FreeonesScraperID, scrapedPerformer)
}

// deprecated
func (r *queryResolver) ScrapeFreeonesPerformerList(ctx context.Context, query string) ([]string, error) {
	scrapedPerformers, err := manager.GetInstance().ScraperCache.ScrapePerformerList(scraper.FreeonesScraperID, query)

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
	return manager.GetInstance().ScraperCache.ListPerformerScrapers(), nil
}

func (r *queryResolver) ListSceneScrapers(ctx context.Context) ([]*models.Scraper, error) {
	return manager.GetInstance().ScraperCache.ListSceneScrapers(), nil
}

func (r *queryResolver) ScrapePerformerList(ctx context.Context, scraperID string, query string) ([]*models.ScrapedPerformer, error) {
	if query == "" {
		return nil, nil
	}

	return manager.GetInstance().ScraperCache.ScrapePerformerList(scraperID, query)
}

func (r *queryResolver) ScrapePerformer(ctx context.Context, scraperID string, scrapedPerformer models.ScrapedPerformerInput) (*models.ScrapedPerformer, error) {
	return manager.GetInstance().ScraperCache.ScrapePerformer(scraperID, scrapedPerformer)
}

func (r *queryResolver) ScrapePerformerURL(ctx context.Context, url string) (*models.ScrapedPerformer, error) {
	return manager.GetInstance().ScraperCache.ScrapePerformerURL(url)
}

func (r *queryResolver) ScrapeScene(ctx context.Context, scraperID string, scene models.SceneUpdateInput) (*models.ScrapedScene, error) {
	return manager.GetInstance().ScraperCache.ScrapeScene(scraperID, scene)
}

func (r *queryResolver) ScrapeSceneURL(ctx context.Context, url string) (*models.ScrapedScene, error) {
	return manager.GetInstance().ScraperCache.ScrapeSceneURL(url)
}
