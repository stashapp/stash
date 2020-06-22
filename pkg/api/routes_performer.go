package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type performerRoutes struct{}

func (rs performerRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{performerId}", func(r chi.Router) {
		r.Use(PerformerCtx)
		r.Get("/image", rs.Image)
	})

	return r
}

func (rs performerRoutes) Image(w http.ResponseWriter, r *http.Request) {
	performer := r.Context().Value(performerKey).(*models.Performer)
	qb := models.NewPerformerQueryBuilder()
	image, _ := qb.GetPerformerImage(performer.ID, nil)
	utils.ServeImage(image, w, r)
}

func PerformerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		performerID, err := strconv.Atoi(chi.URLParam(r, "performerId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		qb := models.NewPerformerQueryBuilder()
		performer, err := qb.Find(performerID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), performerKey, performer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
