package api

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type ImageFinder interface {
	Find(ctx context.Context, id int) (*models.Image, error)
	FindByChecksum(ctx context.Context, checksum string) ([]*models.Image, error)
}

type imageRoutes struct {
	txnManager  txn.Manager
	imageFinder ImageFinder
	fileFinder  file.Finder
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
		const useDefault = true

		f := img.Files.Primary()
		if f == nil {
			rs.serveImage(w, r, img, useDefault)
			return
		}

		encoder := image.NewThumbnailEncoder(manager.GetInstance().FFMPEG)
		data, err := encoder.GetThumbnail(f, models.DefaultGthumbWidth)
		if err != nil {
			// don't log for unsupported image format
			// don't log for file not found - can optionally be logged in serveImage
			if !errors.Is(err, image.ErrNotSupportedForThumbnail) && !errors.Is(err, fs.ErrNotExist) {
				logger.Errorf("error generating thumbnail for %s: %v", f.Path, err)

				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					logger.Errorf("stderr: %s", string(exitErr.Stderr))
				}
			}

			// backwards compatibility - fallback to original image instead
			rs.serveImage(w, r, img, useDefault)
			return
		}

		// write the generated thumbnail to disk if enabled
		if manager.GetInstance().Config.IsWriteImageThumbnails() {
			logger.Debugf("writing thumbnail to disk: %s", img.Path)
			if err := fsutil.WriteFile(filepath, data); err != nil {
				logger.Errorf("error writing thumbnail for image %s: %v", img.Path, err)
			}
		}
		if n, err := w.Write(data); err != nil && !errors.Is(err, syscall.EPIPE) {
			logger.Errorf("error serving thumbnail (wrote %v bytes out of %v): %v", n, len(data), err)
		}
	}
}

func (rs imageRoutes) Image(w http.ResponseWriter, r *http.Request) {
	i := r.Context().Value(imageKey).(*models.Image)

	const useDefault = false
	rs.serveImage(w, r, i, useDefault)
}

func (rs imageRoutes) serveImage(w http.ResponseWriter, r *http.Request, i *models.Image, useDefault bool) {
	const defaultImageImage = "image/image.svg"

	if i.Files.Primary() != nil {
		err := i.Files.Primary().Serve(&file.OsFS{}, w, r)
		if err == nil {
			return
		}

		if !useDefault {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// only log in debug since it can get noisy
		logger.Debugf("Error serving %s: %v", i.DisplayName(), err)
	}

	if !useDefault {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// fall back to static image
	f, _ := static.Image.Open(defaultImageImage)
	defer f.Close()
	stat, _ := f.Stat()
	http.ServeContent(w, r, "image.svg", stat.ModTime(), f.(io.ReadSeeker))
}

// endregion

func (rs imageRoutes) ImageCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageIdentifierQueryParam := chi.URLParam(r, "imageId")
		imageID, _ := strconv.Atoi(imageIdentifierQueryParam)

		var image *models.Image
		_ = txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			qb := rs.imageFinder
			if imageID == 0 {
				images, _ := qb.FindByChecksum(ctx, imageIdentifierQueryParam)
				if len(images) > 0 {
					image = images[0]
				}
			} else {
				image, _ = qb.Find(ctx, imageID)
			}

			if image != nil {
				if err := image.LoadPrimaryFile(ctx, rs.fileFinder); err != nil {
					if !errors.Is(err, context.Canceled) {
						logger.Errorf("error loading primary file for image %d: %v", imageID, err)
					}
					// set image to nil so that it doesn't try to use the primary file
					image = nil
				}
			}

			return nil
		})
		if image == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), imageKey, image)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
