package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type MovieFinder interface {
	GetFrontImage(ctx context.Context, movieID int) ([]byte, error)
	GetBackImage(ctx context.Context, movieID int) ([]byte, error)
	Find(ctx context.Context, id int) (*models.Movie, error)
}

type movieRoutes struct {
	txnManager  txn.Manager
	movieFinder MovieFinder
}

func (rs movieRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{movieId}", func(r chi.Router) {
		r.Use(rs.MovieCtx)
		r.Get("/frontimage", rs.FrontImage)
		r.Get("/backimage", rs.BackImage)
	})

	return r
}

func (rs movieRoutes) FrontImage(w http.ResponseWriter, r *http.Request) {
	movie := r.Context().Value(movieKey).(*models.Movie)
	defaultParam := r.URL.Query().Get("default")
	var image []byte
	if defaultParam != "true" {
		err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			image, _ = rs.movieFinder.GetFrontImage(ctx, movie.ID)
			return nil
		})
		if err != nil {
			logger.Warnf("read transaction error while getting front image: %v", err)
		}
	}

	if len(image) == 0 {
		image, _ = utils.ProcessBase64Image(models.DefaultMovieImage)
	}

	if err := utils.ServeImage(image, w, r); err != nil {
		logger.Warnf("error serving front image: %v", err)
	}
}

func (rs movieRoutes) BackImage(w http.ResponseWriter, r *http.Request) {
	movie := r.Context().Value(movieKey).(*models.Movie)
	defaultParam := r.URL.Query().Get("default")
	var image []byte
	if defaultParam != "true" {
		err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			image, _ = rs.movieFinder.GetBackImage(ctx, movie.ID)
			return nil
		})
		if err != nil {
			logger.Warnf("read transaction error on fetch back image: %v", err)
		}
	}

	if len(image) == 0 {
		image, _ = utils.ProcessBase64Image(models.DefaultMovieImage)
	}

	if err := utils.ServeImage(image, w, r); err != nil {
		logger.Warnf("error while serving image: %v", err)
	}
}

func (rs movieRoutes) MovieCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		movieID, err := strconv.Atoi(chi.URLParam(r, "movieId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var movie *models.Movie
		if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			var err error
			movie, err = rs.movieFinder.Find(ctx, movieID)
			return err
		}); err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), movieKey, movie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
