package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
)

func (r *mutationResolver) ReloadScrapers(ctx context.Context) (bool, error) {
	manager.GetInstance().RefreshScraperCache()
	return true, nil
}
