package stashbox

import "github.com/stashapp/stash/pkg/models"

type StashBoxStudioQueryResult struct {
	Query   string                  `json:"query"`
	Results []*models.ScrapedStudio `json:"results"`
}

type StashBoxPerformerQueryResult struct {
	Query   string                     `json:"query"`
	Results []*models.ScrapedPerformer `json:"results"`
}
