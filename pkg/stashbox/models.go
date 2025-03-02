package stashbox

import "github.com/stashapp/stash/pkg/models"

type StudioQueryResult struct {
	Query   string                  `json:"query"`
	Results []*models.ScrapedStudio `json:"results"`
}

type PerformerQueryResult struct {
	Query   string                     `json:"query"`
	Results []*models.ScrapedPerformer `json:"results"`
}
