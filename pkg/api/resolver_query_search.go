package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) Search(ctx context.Context, query string, ty models.SearchType) (*models.SearchResult, error) {
	s, err := r.searchEngine.Search(ctx, query, ty)
	if err != nil {
		return nil, err
	}

	var x []models.SearchItem
	for _, c := range s.Content {
		x = append(x, models.SearchSceneItem{
			ID: c,
		})
	}
	res := models.SearchResult{
		Content: x,
		Took:    s.Took.Seconds(),
	}

	return &res, nil
}
