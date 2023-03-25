package urlbuilders

import (
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

type TagURLBuilder struct {
	BaseURL string
	TagID   string
}

func NewTagURLBuilder(baseURL string, tag *models.Tag) TagURLBuilder {
	return TagURLBuilder{
		BaseURL: baseURL,
		TagID:   strconv.Itoa(tag.ID),
	}
}

func (b TagURLBuilder) GetTagImageURL() string {
	return b.BaseURL + "/tag/" + b.TagID + "/image"
}
