package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindMovie(ctx context.Context, id string) (*models.Movie, error) {
	qb := models.NewMovieQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt, nil)
}

func (r *queryResolver) FindMovies(ctx context.Context, movieFilter *models.MovieFilterType, filter *models.FindFilterType) (*models.FindMoviesResultType, error) {
	qb := models.NewMovieQueryBuilder()
	movies, total := qb.Query(movieFilter, filter)
	return &models.FindMoviesResultType{
		Count:  total,
		Movies: movies,
	}, nil
}

func (r *queryResolver) AllMovies(ctx context.Context) ([]*models.Movie, error) {
	qb := models.NewMovieQueryBuilder()
	return qb.All()
}

func (r *queryResolver) AllMoviesSlim(ctx context.Context) ([]*models.Movie, error) {
	qb := models.NewMovieQueryBuilder()
	return qb.AllSlim()
}
