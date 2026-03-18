package api

import (
	"errors"

	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindMovie(ctx context.Context, id string) (ret *models.Group, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Group.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindMovies(
	ctx context.Context,
	movieFilter *models.GroupFilterType,
	savedFilterID *string,
	filter *models.FindFilterType,
	ids []string,
) (ret *FindMoviesResultType, err error) {
	if movieFilter != nil && savedFilterID != nil {
		return nil, errors.New("cannot provide both movieFilter and saved_filter_id")
	}

	var finalFilter *models.GroupFilterType
	if savedFilterID != nil {
		finalFilter = &models.GroupFilterType{}
		mode := models.FilterModeMovies

		mergedFindFilter, err := r.resolveSavedFilter(ctx, *savedFilterID, mode, finalFilter, filter)
		if err != nil {
			return nil, err
		}
		filter = mergedFindFilter
	} else {
		finalFilter = movieFilter
	}
	idInts, err := handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var groups []*models.Group
		var err error
		var total int

		if len(idInts) > 0 {
			groups, err = r.repository.Group.FindMany(ctx, idInts)
			total = len(groups)
		} else {
			groups, total, err = r.repository.Group.Query(ctx, finalFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &FindMoviesResultType{
			Count:  total,
			Movies: groups,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllMovies(ctx context.Context) (ret []*models.Group, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Group.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
