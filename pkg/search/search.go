package search

import (
	"context"
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/search/documents"
	"github.com/stashapp/stash/pkg/utils"
)

type Item struct {
	ID    string
	Type  documents.DocType
	Score float64
}

func newItem(nodeID string, score float64) *Item {
	ty, id, ok := utils.Cut(nodeID, ":")
	if !ok {
		return nil
	}

	return &Item{
		ID:    id,
		Type:  documents.NewDocType(ty),
		Score: score,
	}
}

type Result struct {
	Items  []Item
	Facets search.FacetResults

	Total    uint64
	MaxScore float64
	Status   *models.SearchResultStatus
	Took     time.Duration
}

func (e *Engine) Search(ctx context.Context, in string, ty *models.SearchType, facets []*models.SearchFacet) (*Result, error) {
	queryString := bleve.NewQueryStringQuery(in)

	var q query.Query
	if ty == nil {
		q = queryString
	} else {
		var filter *query.MatchQuery

		switch *ty {
		case models.SearchTypeSearchPerformer:
			filter = bleve.NewMatchQuery(string(documents.TypePerformer))
		case models.SearchTypeSearchScene:
			filter = bleve.NewMatchQuery(string(documents.TypeScene))
		case models.SearchTypeSearchStudio:
			filter = bleve.NewMatchQuery(string(documents.TypeStudio))
		case models.SearchTypeSearchTag:
			filter = bleve.NewMatchQuery(string(documents.TypeTag))
		}

		filter.SetField("stash_type")
		q = bleve.NewConjunctionQuery(queryString, filter)
	}

	searchRequest := bleve.NewSearchRequest(q)

	for _, f := range facets {
		if f == nil {
			continue
		}

		switch *f {
		case models.SearchFacetDateRange:
			var cutOffDate = time.Now().Add(-30 * 24 * time.Hour)
			dateFacet := bleve.NewFacetRequest("date", 2)
			dateFacet.AddDateTimeRange("old", time.Unix(0, 0), cutOffDate)
			dateFacet.AddDateTimeRange("new", cutOffDate, time.Unix(9999999999999, 999999999))
			searchRequest.AddFacet("released", dateFacet)
		}
	}

	// Hold e.mu for as short as possible
	e.mu.RLock()
	searchResult, err := e.idx.SearchInContext(ctx, searchRequest)
	e.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	var items []Item
	for _, match := range searchResult.Hits {
		i := newItem(match.ID, match.Score)
		items = append(items, *i)
	}

	var status *models.SearchResultStatus
	if searchResult.Status != nil {
		st := searchResult.Status

		status = &models.SearchResultStatus{
			Successful: st.Successful,
			Failed:     st.Failed,
			Total:      st.Total,
		}
	}

	res := Result{
		Items:  items,
		Facets: searchResult.Facets,

		Took:     searchResult.Took,
		Total:    searchResult.Total,
		MaxScore: searchResult.MaxScore,

		Status: status,
	}

	return &res, nil
}

func tagID(id int) string {
	return fmt.Sprintf("tag:%d", id)
}

func sceneID(id int) string {
	return fmt.Sprintf("scene:%d", id)
}

func performerID(id int) string {
	return fmt.Sprintf("performer:%d", id)
}

func studioID(id int) string {
	return fmt.Sprintf("studio:%d", id)
}
