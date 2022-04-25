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

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Movie().Find(idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindMovies(ctx context.Context, movieFilter *models.MovieFilterType, filter *models.FindFilterType) (ret *FindMoviesResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		movies, total, err := repo.Movie().Query(movieFilter, filter)
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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Movie().All()
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
