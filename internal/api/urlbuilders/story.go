package urlbuilders

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type StoryURLBuilder struct {
	BaseURL   string
	StoryID   string
	UpdatedAt string
}

func NewStoryURLBuilder(baseURL string, story *models.Story) StoryURLBuilder {
	return StoryURLBuilder{
		BaseURL:   baseURL,
		StoryID:   strconv.Itoa(story.ID),
		UpdatedAt: strconv.FormatInt(story.UpdatedAt.Unix(), 10),
	}
}

func (b StoryURLBuilder) GetFrontImageURL() string {
	return b.BaseURL + "/story/" + b.StoryID + "/front_cover?t=" + b.UpdatedAt
}

func (b StoryURLBuilder) GetBackImageURL() string {
	return b.BaseURL + "/story/" + b.StoryID + "/back_cover?t=" + b.UpdatedAt
}
