package search

import (
	"context"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type Item struct {
	ID    string
	Type  string
	Score float64
}

type Result struct {
	Items []Item
	Took  time.Duration
}

func newItem(nodeID string, score float64) *Item {
	ty, id, ok := utils.Cut(nodeID, ":")
	if !ok {
		return nil
	}

	return &Item{
		ID:    id,
		Type:  ty,
		Score: score,
	}
}

func (e *Engine) Search(ctx context.Context, in string, ty models.SearchType) (*Result, error) {
	// Hold e.mu for as short as possible
	e.mu.RLock()
	query := bleve.NewQueryStringQuery(in)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := e.sceneIdx.SearchInContext(ctx, searchRequest)
	e.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, match := range searchResult.Hits {
		i := newItem(match.ID, match.Score)
		items = append(items, *i)
	}
	res := Result{
		Items: items,
		Took:  searchResult.Took,
	}

	return &res, nil
}
