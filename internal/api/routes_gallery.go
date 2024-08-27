package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GalleryFinder interface {
	models.GalleryGetter
	FindByChecksum(ctx context.Context, checksum string) ([]*models.Gallery, error)
}

type ImageByIndexer interface {
	FindByGalleryIDIndex(ctx context.Context, galleryID int, index uint) (*models.Image, error)
}

type galleryRoutes struct {
	routes
	imageRoutes   imageRoutes
	galleryFinder GalleryFinder
	imageFinder   ImageByIndexer
	fileGetter    models.FileGetter
}

func (rs galleryRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{galleryId}", func(r chi.Router) {
		r.Use(rs.GalleryCtx)

		r.Get("/preview/{imageIndex}", rs.Preview)
	})

	return r
}

func (rs galleryRoutes) Preview(w http.ResponseWriter, r *http.Request) {
	g := r.Context().Value(galleryKey).(*models.Gallery)
	indexQueryParam := chi.URLParam(r, "imageIndex")
	var i *models.Image

	index, err := strconv.Atoi(indexQueryParam)
	if err != nil || index < 0 {
		http.Error(w, "bad index", 400)
		return
	}

	_ = rs.withReadTxn(r, func(ctx context.Context) error {
		qb := rs.imageFinder
		i, _ = qb.FindByGalleryIDIndex(ctx, g.ID, uint(index))
		// TODO - handle errors?

		// serveThumbnail needs files populated
		if err := i.LoadPrimaryFile(ctx, rs.fileGetter); err != nil {
			if !errors.Is(err, context.Canceled) {
				logger.Errorf("error loading primary file for image %d: %v", i.ID, err)
			}
			// set image to nil so that it doesn't try to use the primary file
			i = nil
		}

		return nil
	})
	if i == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	rs.imageRoutes.serveThumbnail(w, r, i)
}

func (rs galleryRoutes) GalleryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		galleryIdentifierQueryParam := chi.URLParam(r, "galleryId")
		galleryID, _ := strconv.Atoi(galleryIdentifierQueryParam)

		var gallery *models.Gallery
		_ = rs.withReadTxn(r, func(ctx context.Context) error {
			qb := rs.galleryFinder
			if galleryID == 0 {
				galleries, _ := qb.FindByChecksum(ctx, galleryIdentifierQueryParam)
				if len(galleries) > 0 {
					gallery = galleries[0]
				}
			} else {
				gallery, _ = qb.Find(ctx, galleryID)
			}

			if gallery != nil {
				if err := gallery.LoadPrimaryFile(ctx, rs.fileGetter); err != nil {
					if !errors.Is(err, context.Canceled) {
						logger.Errorf("error loading primary file for gallery %d: %v", galleryID, err)
					}
					// set image to nil so that it doesn't try to use the primary file
					gallery = nil
				}
			}

			return nil
		})
		if gallery == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), galleryKey, gallery)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
