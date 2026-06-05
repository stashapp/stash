package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindStory(ctx context.Context, id string) (*models.Story, error) {
	qb := r.repository.Story
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return qb.Find(ctx, idNum)
}

func (r *queryResolver) FindStories(ctx context.Context, storyFilter *models.StoryFilterType, filter *models.FindFilterType, ids []string) (*FindStoriesResultType, error) {
	qb := r.repository.Story
	stories, count, err := qb.Query(ctx, storyFilter, filter)
	if err != nil {
		return nil, err
	}
	return &FindStoriesResultType{
		Count:   count,
		Stories: stories,
	}, nil
}

func (r *queryResolver) AllStories(ctx context.Context) ([]*models.Story, error) {
	return r.repository.Story.All(ctx)
}
