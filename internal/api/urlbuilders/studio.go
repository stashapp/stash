package urlbuilders

import (
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

type StudioURLBuilder struct {
	BaseURL  string
	StudioID string
}

func NewStudioURLBuilder(baseURL string, studio *models.Studio) StudioURLBuilder {
	return StudioURLBuilder{
		BaseURL:  baseURL,
		StudioID: strconv.Itoa(studio.ID),
	}
}

func (b StudioURLBuilder) GetStudioImageURL() string {
	return b.BaseURL + "/studio/" + b.StudioID + "/image"
}
