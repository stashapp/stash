package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type StudioFinder interface {
	models.StudioGetter
	GetImage(ctx context.Context, studioID int) ([]byte, error)
}

type studioRoutes struct {
	routes
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
		readTxnErr := rs.withReadTxn(r, func(ctx context.Context) error {
			var err error
			image, err = rs.studioFinder.GetImage(ctx, studio.ID)
			return err
		})
		if errors.Is(readTxnErr, context.Canceled) {
			return
		}
		if readTxnErr != nil {
			logger.Warnf("read transaction error on fetch studio image: %v", readTxnErr)
		}
	}

	// fallback to default image
	if len(image) == 0 {
		image = static.ReadAll(static.DefaultStudioImage)
	}

	utils.ServeImage(w, r, image)
}

func (rs studioRoutes) StudioCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		studioID, err := strconv.Atoi(chi.URLParam(r, "studioId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var studio *models.Studio
		_ = rs.withReadTxn(r, func(ctx context.Context) error {
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
