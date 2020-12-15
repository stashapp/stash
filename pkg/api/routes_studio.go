package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type studioRoutes struct{}

func (rs studioRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{studioId}", func(r chi.Router) {
		r.Use(StudioCtx)
		r.Get("/image", rs.Image)
	})

	return r
}

func (rs studioRoutes) Image(w http.ResponseWriter, r *http.Request) {
	studio := r.Context().Value(studioKey).(*models.Studio)
	defaultParam := r.URL.Query().Get("default")

	var image []byte
	if defaultParam != "true" {
		manager.GetInstance().WithTxn(r.Context(), func(repo models.Repository) error {
			image, _ = repo.Studio().GetImage(studio.ID)
			return nil
		})
	}

	if len(image) == 0 {
		_, image, _ = utils.ProcessBase64Image(models.DefaultStudioImage)
	}

	utils.ServeImage(image, w, r)
}

func StudioCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		studioID, err := strconv.Atoi(chi.URLParam(r, "studioId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var studio *models.Studio
		if err := manager.GetInstance().WithTxn(r.Context(), func(repo models.Repository) error {
			var err error
			studio, err = repo.Studio().Find(studioID)
			return err
		}); err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), studioKey, studio)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
