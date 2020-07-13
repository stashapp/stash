package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *scrapedSceneTagResolver) StoredID(ctx context.Context, obj *models.ScrapedSceneTag) (*string, error) {
	return obj.ID, nil
}
