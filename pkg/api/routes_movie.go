package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/models"
)

type movieRoutes struct{}

func (rs movieRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{movieId}", func(r chi.Router) {
		r.Use(MovieCtx)
		r.Get("/frontimage", rs.FrontImage)
		r.Get("/backimage", rs.BackImage)
	})

	return r
}

func (rs movieRoutes) FrontImage(w http.ResponseWriter, r *http.Request) {
	movie := r.Context().Value(movieKey).(*models.Movie)
	_, _ = w.Write(movie.Front_Image)
}

func (rs movieRoutes) BackImage(w http.ResponseWriter, r *http.Request) {
	movie := r.Context().Value(movieKey).(*models.Movie)
	_, _ = w.Write(movie.Back_Image)
}

func MovieCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		movieID, err := strconv.Atoi(chi.URLParam(r, "movieId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		qb := models.NewMovieQueryBuilder()
		movie, err := qb.Find(movieID, nil)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), movieKey, movie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
