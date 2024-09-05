package urlbuilders

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type GroupURLBuilder struct {
	BaseURL   string
	GroupID   string
	UpdatedAt string
}

func NewGroupURLBuilder(baseURL string, group *models.Group) GroupURLBuilder {
	return GroupURLBuilder{
		BaseURL:   baseURL,
		GroupID:   strconv.Itoa(group.ID),
		UpdatedAt: strconv.FormatInt(group.UpdatedAt.Unix(), 10),
	}
}

func (b GroupURLBuilder) GetGroupFrontImageURL(hasImage bool) string {
	url := b.BaseURL + "/group/" + b.GroupID + "/frontimage?t=" + b.UpdatedAt
	if !hasImage {
		url += "&default=true"
	}
	return url
}

func (b GroupURLBuilder) GetGroupBackImageURL() string {
	return b.BaseURL + "/group/" + b.GroupID + "/backimage?t=" + b.UpdatedAt
}
