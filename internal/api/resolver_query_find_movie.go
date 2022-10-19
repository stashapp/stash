package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindMovie(ctx context.Context, id string) (ret *models.Movie, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindMovies(ctx context.Context, movieFilter *models.MovieFilterType, filter *models.FindFilterType) (ret *FindMoviesResultType, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		movies, total, err := r.repository.Movie.Query(ctx, movieFilter, filter)
		if err != nil {
			return err
		}

		ret = &FindMoviesResultType{
			Count:  total,
			Movies: movies,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllMovies(ctx context.Context) (ret []*models.Movie, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
