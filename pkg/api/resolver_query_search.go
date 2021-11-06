package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/search"
)

func (r *queryResolver) Search(ctx context.Context, query string, ty models.SearchType) (*models.SearchResultItemConnection, error) {
	s, err := r.searchEngine.Search(ctx, query, ty)
	if err != nil {
		return nil, err
	}

	var edges []*models.SearchItemEdge
	for _, item := range s.Items {
		h := r.hydrate(ctx, item)
		if h != nil {
			edges = append(edges, &models.SearchItemEdge{
				Score: item.Score,
				Node:  h,
			})
		}
	}
	res := models.SearchResultItemConnection{
		Edges: edges,
		Took:  s.Took.Seconds(),
	}

	return &res, nil
}

func (r *queryResolver) hydrate(ctx context.Context, item search.Item) models.SearchResultItem {
	switch item.Type {
	case "Scene":
		scene, err := r.FindScene(ctx, &item.ID, nil)
		if err != nil {
			return nil
		}

		return scene
	}

	return nil
}
