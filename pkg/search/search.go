package search

import (
	"context"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/stashapp/stash/pkg/models"
)

type Result struct {
	Content []string
	Took    time.Duration
}

func (e *Engine) Search(ctx context.Context, in string, ty models.SearchType) (*Result, error) {
	query := bleve.NewQueryStringQuery(in)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := e.sceneIdx.SearchInContext(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	var content []string
	for _, match := range searchResult.Hits {
		id := match.ID
		content = append(content, id)
	}
	res := Result{
		Content: content,
		Took:    searchResult.Took,
	}

	return &res, nil
}
