package urlbuilders

import (
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

type PerformerURLBuilder struct {
	BaseURL     string
	PerformerID string
	UpdatedAt   string
}

func NewPerformerURLBuilder(baseURL string, performer *models.Performer) PerformerURLBuilder {
	return PerformerURLBuilder{
		BaseURL:     baseURL,
		PerformerID: strconv.Itoa(performer.ID),
		UpdatedAt:   strconv.FormatInt(performer.UpdatedAt.Timestamp.Unix(), 10),
	}
}

func (b PerformerURLBuilder) GetPerformerImageURL() string {
	return b.BaseURL + "/performer/" + b.PerformerID + "/image?" + b.UpdatedAt
}
