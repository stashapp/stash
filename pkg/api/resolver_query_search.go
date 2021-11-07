package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/search"
)

func (r *queryResolver) Search(ctx context.Context, query string, ty models.SearchType, facets []*models.SearchFacet) (*models.SearchResultItemConnection, error) {
	s, err := r.searchEngine.Search(ctx, query, ty, facets)
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

	var facetResults []*models.SearchFacetResult
	for k, f := range s.Facets {
		var dateRanges []*models.SearchDateRangeFacetResult
		for _, dr := range f.DateRanges {
			drRes := &models.SearchDateRangeFacetResult{
				Name:  dr.Name,
				Count: dr.Count,
				Start: dr.Start,
				End:   dr.End,
			}

			dateRanges = append(dateRanges, drRes)
		}

		facetResults = append(facetResults, &models.SearchFacetResult{
			Name:    k,
			Total:   f.Total,
			Missing: f.Missing,
			Other:   f.Other,

			DateRanges: dateRanges,
		})
	}

	res := models.SearchResultItemConnection{
		Edges:  edges,
		Facets: facetResults,
		Took:   s.Took.Seconds(),
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
