package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type movieRoutes struct {
	txnManager models.TransactionManager
}

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
	defaultParam := r.URL.Query().Get("default")
	var image []byte
	if defaultParam != "true" {
		err := rs.txnManager.WithReadTxn(r.Context(), func(repo models.ReaderRepository) error {
			image, _ = repo.Movie().GetFrontImage(movie.ID)
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
		err := rs.txnManager.WithReadTxn(r.Context(), func(repo models.ReaderRepository) error {
			image, _ = repo.Movie().GetBackImage(movie.ID)
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

func MovieCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		movieID, err := strconv.Atoi(chi.URLParam(r, "movieId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var movie *models.Movie
		if err := manager.GetInstance().TxnManager.WithReadTxn(r.Context(), func(repo models.ReaderRepository) error {
			var err error
			movie, err = repo.Movie().Find(movieID)
			return err
		}); err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), movieKey, movie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
