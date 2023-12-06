package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/internal/manager"
)

type downloadsRoutes struct{}

func (rs downloadsRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{downloadHash}", func(r chi.Router) {
		r.Use(downloadCtx)
		r.Get("/{filename}", rs.file)
	})

	return r
}

func (rs downloadsRoutes) file(w http.ResponseWriter, r *http.Request) {
	hash := r.Context().Value(downloadKey).(string)
	if hash == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	manager.GetInstance().DownloadStore.Serve(hash, w, r)
}

func downloadCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downloadHash := chi.URLParam(r, "downloadHash")

		ctx := context.WithValue(r.Context(), downloadKey, downloadHash)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
