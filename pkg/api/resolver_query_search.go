package api

import (
	"context"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *queryResolver) Search(ctx context.Context, query string, ty models.SearchType) (*models.SearchResultItemConnection, error) {
	s, err := r.searchEngine.Search(ctx, query, ty)
	if err != nil {
		return nil, err
	}

	var edges []*models.SearchItemEdge
	for _, c := range s.Content {
		h := r.hydrate(ctx, c)
		if h != nil {

			edges = append(edges, &models.SearchItemEdge{
				Node: h,
			})
		}
	}
	res := models.SearchResultItemConnection{
		Edges: edges,
		Took:  s.Took.Seconds(),
	}

	return &res, nil
}

func (r *queryResolver) hydrate(ctx context.Context, nodeID string) models.SearchResultItem {
	ty, id, ok := utils.Cut(nodeID, ":")
	if !ok {
		logger.Warnf("hydration of search engine id failed: id=%s", nodeID)
		return nil
	}

	switch ty {
	case "Scene":
		scene, err := r.FindScene(ctx, &id, nil)
		if err != nil {
			return nil
		}

		return scene
	}

	return nil
}
