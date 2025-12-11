package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindStudio(ctx context.Context, id string) (ret *models.Studio, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Studio.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindStudios(ctx context.Context, studioFilter *models.StudioFilterType, filter *models.FindFilterType, ids []string) (ret *FindStudiosResultType, err error) {
	idInts, err := handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var studios []*models.Studio
		var err error
		var total int

		if len(idInts) > 0 {
			studios, err = r.repository.Studio.FindMany(ctx, idInts)
			total = len(studios)
		} else {
			studios, total, err = r.repository.Studio.Query(ctx, studioFilter, filter)
		}
		if err != nil {
			return err
		}

		ret = &FindStudiosResultType{
			Count:   total,
			Studios: studios,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllStudios(ctx context.Context) (ret []*models.Studio, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
