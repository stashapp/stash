package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type tagRoutes struct{}

func (rs tagRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{tagId}", func(r chi.Router) {
		r.Use(TagCtx)
		r.Get("/image", rs.Image)
	})

	return r
}

func (rs tagRoutes) Image(w http.ResponseWriter, r *http.Request) {
	tag := r.Context().Value(tagKey).(*models.Tag)
	qb := models.NewTagQueryBuilder()
	image, _ := qb.GetTagImage(tag.ID, nil)

	// use default image if not present
	if len(image) == 0 {
		image = models.DefaultTagImage
	}

	utils.ServeImage(image, w, r)
}

func TagCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(chi.URLParam(r, "tagId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		qb := models.NewTagQueryBuilder()
		tag, err := qb.Find(tagID, nil)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), tagKey, tag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
