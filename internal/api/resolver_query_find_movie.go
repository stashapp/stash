package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *queryResolver) FindMovie(ctx context.Context, id string) (ret *models.Movie, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindMovies(ctx context.Context, movieFilter *models.MovieFilterType, filter *models.FindFilterType, ids []string) (ret *FindMoviesResultType, err error) {
	idInts, err := stringslice.StringSliceToIntSlice(ids)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var movies []*models.Movie
		var err error
		var total int

		if len(idInts) > 0 {
			movies, err = r.repository.Movie.FindMany(ctx, idInts)
			total = len(movies)
		} else {
			movies, total, err = r.repository.Movie.Query(ctx, movieFilter, filter)
		}

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
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
