package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindPerformer(ctx context.Context, id string) (ret *models.Performer, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.performer.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindPerformers(ctx context.Context, performerFilter *models.PerformerFilterType, filter *models.FindFilterType) (ret *FindPerformersResultType, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		performers, total, err := r.performer.Query(ctx, performerFilter, filter)
		if err != nil {
			return err
		}

		ret = &FindPerformersResultType{
			Count:      total,
			Performers: performers,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllPerformers(ctx context.Context) (ret []*models.Performer, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.performer.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
