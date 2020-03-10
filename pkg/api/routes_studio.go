package api

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/models"
	"net/http"
	"strconv"
	"strings"
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
	etag := fmt.Sprintf("%x", md5.Sum(studio.Image))
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	contentType := http.DetectContentType(studio.Image)
	if contentType == "text/xml; charset=utf-8" {
		contentType = "image/svg+xml"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Add("Etag", etag)
	_, _ = w.Write(studio.Image)
}

func StudioCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		studioID, err := strconv.Atoi(chi.URLParam(r, "studioId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		qb := models.NewStudioQueryBuilder()
		studio, err := qb.Find(studioID, nil)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), studioKey, studio)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
