package api

import (
	"context"

	"github.com/stashapp/stash/pkg/scraper"
)

func (r *mutationResolver) ReloadScrapers(ctx context.Context) (bool, error) {
	err := scraper.ReloadScrapers()

	if err != nil {
		return false, err
	}

	return true, nil
}
