package sqlite

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type storyFilterHandler struct {
	storyFilter *models.StoryFilterType
}

func (s *storyFilterHandler) validate() error {
	return nil
}

func (s *storyFilterHandler) handler() *filterHandler {
	return &filterHandler{}
}
