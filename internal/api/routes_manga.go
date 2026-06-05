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

type MangaFinder interface {
	models.MangaGetter
}

type mangaRoutes struct {
	routes
	mangaFinder MangaFinder
}

func (rs mangaRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{mangaId}", func(r chi.Router) {
		r.Use(rs.MangaCtx)

		r.Get("/cover", rs.Cover)
	})

	return r
}

func (rs mangaRoutes) Cover(w http.ResponseWriter, r *http.Request) {
	m := r.Context().Value(mangaKey).(*models.Manga)

	// serve the cover image blob if available
	if m.CoverImageBlob != "" {
		// For now, serve a default since we don't have blob serving set up
		image := static.ReadAll(static.DefaultGalleryImage)
		utils.ServeImage(w, r, image)
		return
	}

	// fallback to default image
	image := static.ReadAll(static.DefaultGalleryImage)
	utils.ServeImage(w, r, image)
}

func (rs mangaRoutes) MangaCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mangaIdentifierQueryParam := chi.URLParam(r, "mangaId")
		mangaID, _ := strconv.Atoi(mangaIdentifierQueryParam)

		var manga *models.Manga
		_ = rs.withReadTxn(r, func(ctx context.Context) error {
			qb := rs.mangaFinder
			if mangaID == 0 {
				http.Error(w, http.StatusText(404), 404)
				return nil
			} else {
				manga, _ = qb.Find(ctx, mangaID)
			}

			if manga == nil {
				return errors.New("manga not found")
			}

			return nil
		})
		if manga == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), mangaKey, manga)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
