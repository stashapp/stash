package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *tagResolver) SceneCount(ctx context.Context, obj *models.Tag) (ret *int, err error) {
	var count int
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		count, err = repo.Scene().CountByTagID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &count, err
}

func (r *tagResolver) SceneMarkerCount(ctx context.Context, obj *models.Tag) (ret *int, err error) {
	var count int
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		count, err = repo.SceneMarker().CountByTagID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &count, err
}

func (r *tagResolver) PerformerCount(ctx context.Context, obj *models.Tag) (ret *int, err error) {
	var count int
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		count, err = repo.Performer().CountByTagID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return &count, err
}

func (r *tagResolver) ImagePath(ctx context.Context, obj *models.Tag) (*string, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewTagURLBuilder(baseURL, obj.ID).GetTagImageURL()
	return &imagePath, nil
}
