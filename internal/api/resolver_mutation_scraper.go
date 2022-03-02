package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
)

func (r *mutationResolver) ReloadScrapers(ctx context.Context) (bool, error) {
	err := manager.GetInstance().ScraperCache.ReloadScrapers()

	if err != nil {
		return false, err
	}

	return true, nil
}
