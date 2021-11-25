package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/search"
	"github.com/stashapp/stash/pkg/search/documents"
)

var ErrUnknownType = errors.New("unknown item type")

func (r *queryResolver) Search(ctx context.Context, query string, ty *models.SearchType) (*models.SearchResultItemConnection, error) {
	s, err := r.searchEngine.Search(ctx, query, ty)
	if err != nil {
		return nil, err
	}

	var edges []*models.SearchItemEdge
	for _, item := range s.Items {
		h, err := r.hydrate(ctx, item)
		if err != nil {
			edges = append(edges, nil)
			continue
		}

		edges = append(edges, &models.SearchItemEdge{
			Score: item.Score,
			Node:  h,
		})
	}

	res := models.SearchResultItemConnection{
		Edges:    edges,
		Took:     s.Took.Seconds(),
		MaxScore: s.MaxScore,
		Total:    int(s.Total),

		Status: s.Status,
	}

	return &res, nil
}

func (r *queryResolver) hydrate(ctx context.Context, item search.Item) (models.SearchResultItem, error) {
	switch item.Type {
	case documents.TypeScene:
		return r.FindScene(ctx, &item.ID, nil)
	case documents.TypePerformer:
		return r.FindPerformer(ctx, item.ID)
	case documents.TypeTag:
		return r.FindTag(ctx, item.ID)
	case documents.TypeStudio:
		return r.FindStudio(ctx, item.ID)
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownType, item.Type)
	}
}
