package stashbox

import "github.com/stashapp/stash/pkg/models"

type StashBoxPerformerQueryResult struct {
	Query   string                     `json:"query"`
	Results []*models.ScrapedPerformer `json:"results"`
}

type StashBoxStudioQueryResult struct {
	Query   string                  `json:"query"`
	Results []*models.ScrapedStudio `json:"results"`
}
