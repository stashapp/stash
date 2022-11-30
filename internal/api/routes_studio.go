package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type StudioFinder interface {
	studio.Finder
	GetImage(ctx context.Context, studioID int) ([]byte, error)
}

type studioRoutes struct {
	txnManager   txn.Manager
	studioFinder StudioFinder
}

func (rs studioRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{studioId}", func(r chi.Router) {
		r.Use(rs.StudioCtx)
		r.Get("/image", rs.Image)
	})

	return r
}

func (rs studioRoutes) Image(w http.ResponseWriter, r *http.Request) {
	studio := r.Context().Value(studioKey).(*models.Studio)
	defaultParam := r.URL.Query().Get("default")

	var image []byte
	if defaultParam != "true" {
		readTxnErr := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			image, _ = rs.studioFinder.GetImage(ctx, studio.ID)
			return nil
		})
		if errors.Is(readTxnErr, context.Canceled) {
			return
		}
		if readTxnErr != nil {
			logger.Warnf("read transaction error on fetch studio image: %v", readTxnErr)
		}
	}

	if len(image) == 0 {
		image, _ = utils.ProcessBase64Image(models.DefaultStudioImage)
	}

	if err := utils.ServeImage(image, w, r); err != nil {
		logger.Warnf("error serving studio image: %v", err)
	}
}

func (rs studioRoutes) StudioCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		studioID, err := strconv.Atoi(chi.URLParam(r, "studioId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var studio *models.Studio
		_ = txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			var err error
			studio, err = rs.studioFinder.Find(ctx, studioID)
			return err
		})
		if studio == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), studioKey, studio)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
