package urlbuilders

import (
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

type TagURLBuilder struct {
	BaseURL   string
	TagID     string
	UpdatedAt string
}

func NewTagURLBuilder(baseURL string, tag *models.Tag) TagURLBuilder {
	return TagURLBuilder{
		BaseURL:   baseURL,
		TagID:     strconv.Itoa(tag.ID),
		UpdatedAt: strconv.FormatInt(tag.UpdatedAt.Timestamp.Unix(), 10),
	}
}

func (b TagURLBuilder) GetTagImageURL() string {
	return b.BaseURL + "/tag/" + b.TagID + "/image?" + b.UpdatedAt
}
