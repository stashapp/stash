package api

import (
	"context"
	"errors"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type ImageFinder interface {
	Find(ctx context.Context, id int) (*models.Image, error)
	FindByChecksum(ctx context.Context, checksum string) (*models.Image, error)
}

type imageRoutes struct {
	txnManager  txn.Manager
	imageFinder ImageFinder
}

func (rs imageRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{imageId}", func(r chi.Router) {
		r.Use(rs.ImageCtx)

		r.Get("/image", rs.Image)
		r.Get("/thumbnail", rs.Thumbnail)
	})

	return r
}

// region Handlers

func (rs imageRoutes) Thumbnail(w http.ResponseWriter, r *http.Request) {
	img := r.Context().Value(imageKey).(*models.Image)
	filepath := manager.GetInstance().Paths.Generated.GetThumbnailPath(img.Checksum, models.DefaultGthumbWidth)

	w.Header().Add("Cache-Control", "max-age=604800000")

	// if the thumbnail doesn't exist, encode on the fly
	exists, _ := fsutil.FileExists(filepath)
	if exists {
		http.ServeFile(w, r, filepath)
	} else {
		encoder := image.NewThumbnailEncoder(manager.GetInstance().FFMPEG)
		data, err := encoder.GetThumbnail(img, models.DefaultGthumbWidth)
		if err != nil {
			// don't log for unsupported image format
			if !errors.Is(err, image.ErrNotSupportedForThumbnail) {
				logger.Errorf("error generating thumbnail for image: %s", err.Error())

				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					logger.Errorf("stderr: %s", string(exitErr.Stderr))
				}
			}

			// backwards compatibility - fallback to original image instead
			rs.Image(w, r)
			return
		}

		// write the generated thumbnail to disk if enabled
		if manager.GetInstance().Config.IsWriteImageThumbnails() {
			logger.Debugf("writing thumbnail to disk: %s", img.Path)
			if err := fsutil.WriteFile(filepath, data); err != nil {
				logger.Errorf("error writing thumbnail for image %s: %s", img.Path, err)
			}
		}
		if n, err := w.Write(data); err != nil {
			logger.Errorf("error writing thumbnail response. Wrote %v bytes: %v", n, err)
		}
	}
}

func (rs imageRoutes) Image(w http.ResponseWriter, r *http.Request) {
	i := r.Context().Value(imageKey).(*models.Image)

	// if image is in a zip file, we need to serve it specifically
	image.Serve(w, r, i.Path)
}

// endregion

func (rs imageRoutes) ImageCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageIdentifierQueryParam := chi.URLParam(r, "imageId")
		imageID, _ := strconv.Atoi(imageIdentifierQueryParam)

		var image *models.Image
		readTxnErr := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			qb := rs.imageFinder
			if imageID == 0 {
				image, _ = qb.FindByChecksum(ctx, imageIdentifierQueryParam)
			} else {
				image, _ = qb.Find(ctx, imageID)
			}

			return nil
		})
		if readTxnErr != nil {
			logger.Warnf("read transaction failure while trying to read image by id: %v", readTxnErr)
		}

		if image == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), imageKey, image)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
