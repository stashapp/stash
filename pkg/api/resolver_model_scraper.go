package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *scrapedSceneTagResolver) StoredID(ctx context.Context, obj *models.ScrapedSceneTag) (*string, error) {
	return obj.ID, nil
}

func (r *scrapedSceneMovieResolver) StoredID(ctx context.Context, obj *models.ScrapedSceneMovie) (*string, error) {
	return obj.ID, nil
}

func (r *scrapedScenePerformerResolver) StoredID(ctx context.Context, obj *models.ScrapedScenePerformer) (*string, error) {
	return obj.ID, nil
}

func (r *scrapedSceneStudioResolver) StoredID(ctx context.Context, obj *models.ScrapedSceneStudio) (*string, error) {
	return obj.ID, nil
}
