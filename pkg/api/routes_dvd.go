package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/models"
)

type dvdRoutes struct{}

func (rs dvdRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{dvdId}", func(r chi.Router) {
		r.Use(DvdCtx)
		r.Get("/frontimage", rs.FrontImage)
		r.Get("/backimage", rs.BackImage)
	})

	return r
}

func (rs dvdRoutes) FrontImage(w http.ResponseWriter, r *http.Request) {
	dvd := r.Context().Value(dvdKey).(*models.Dvd)
	_, _ = w.Write(dvd.FrontImage)
}

func (rs dvdRoutes) BackImage(w http.ResponseWriter, r *http.Request) {
	dvd := r.Context().Value(dvdKey).(*models.Dvd)
	_, _ = w.Write(dvd.BackImage)
}

func DvdCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dvdID, err := strconv.Atoi(chi.URLParam(r, "dvdId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		qb := models.NewDvdQueryBuilder()
		dvd, err := qb.Find(dvdID, nil)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), dvdKey, dvd)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
