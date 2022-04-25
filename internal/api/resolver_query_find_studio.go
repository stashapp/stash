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

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var err error
		ret, err = repo.Studio().Find(idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindStudios(ctx context.Context, studioFilter *models.StudioFilterType, filter *models.FindFilterType) (ret *FindStudiosResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		studios, total, err := repo.Studio().Query(studioFilter, filter)
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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Studio().All()
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
