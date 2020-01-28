package api

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/models"
	"net/http"
	"strconv"
)

type galleryRoutes struct{}

func (rs galleryRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{galleryId}", func(r chi.Router) {
		r.Use(GalleryCtx)
		r.Get("/{fileIndex}", rs.File)
	})

	return r
}

func (rs galleryRoutes) File(w http.ResponseWriter, r *http.Request) {
	gallery := r.Context().Value(galleryKey).(*models.Gallery)
	fileIndex, _ := strconv.Atoi(chi.URLParam(r, "fileIndex"))
	thumb := r.URL.Query().Get("thumb")
	w.Header().Add("Cache-Control", "max-age=604800000") // 1 Week
	if thumb == "true" {
		_, _ = w.Write(gallery.GetThumbnail(fileIndex, 200))
	} else if thumb == "" {
		_, _ = w.Write(gallery.GetImage(fileIndex))
	} else {
		height, err := strconv.ParseInt(thumb, 0, 64)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		_, _ = w.Write(gallery.GetThumbnail(fileIndex, int(height)))
	}
}

func GalleryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		galleryID, err := strconv.Atoi(chi.URLParam(r, "galleryId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		qb := models.NewGalleryQueryBuilder()
		gallery, err := qb.Find(galleryID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), galleryKey, gallery)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
