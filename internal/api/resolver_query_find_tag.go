package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindTag(ctx context.Context, id string) (ret *models.Tag, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindTags(ctx context.Context, tagFilter *models.TagFilterType, filter *models.FindFilterType) (ret *FindTagsResultType, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		tags, total, err := r.repository.Tag.Query(ctx, tagFilter, filter)
		if err != nil {
			return err
		}

		ret = &FindTagsResultType{
			Count: total,
			Tags:  tags,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllTags(ctx context.Context) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
