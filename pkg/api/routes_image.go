package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type imageRoutes struct {
	txnManager models.TransactionManager
}

func (rs imageRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{imageId}", func(r chi.Router) {
		r.Use(ImageCtx)

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
	exists, _ := utils.FileExists(filepath)
	if exists {
		http.ServeFile(w, r, filepath)
	} else {
		encoder := image.NewThumbnailEncoder(manager.GetInstance().FFMPEG)

		// try to read the first associated file
		files := r.Context().Value(filesKey).([]*models.File)
		if len(files) == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		reader, err := file.Open(files[0])
		if err != nil {
			logger.Errorf("error generating thumbnail for %q: %v", files[0].Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusNotFound)
			return
		}
		defer reader.Close()

		data, err := encoder.GetThumbnail(reader, models.DefaultGthumbWidth)
		if err != nil {
			// don't log if too small
			if !errors.Is(err, image.ErrTooSmall) {
				logger.Errorf("error generating thumbnail for image: %s", err.Error())
			}

			// backwards compatibility - fallback to original image instead
			rs.Image(w, r)
			return
		}

		// write the generated thumbnail to disk if enabled
		if manager.GetInstance().Config.IsWriteImageThumbnails() {
			if err := utils.WriteFile(filepath, data); err != nil {
				logger.Errorf("error writing thumbnail for image %s: %s", img.Path, err)
			}
		}
		if n, err := w.Write(data); err != nil {
			logger.Errorf("error writing thumbnail response. Wrote %v bytes: %v", n, err)
		}
	}
}

func (rs imageRoutes) Image(w http.ResponseWriter, r *http.Request) {
	files := r.Context().Value(filesKey).([]*models.File)
	if len(files) == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// if image is in a zip file, we need to serve it specifically
	file.Serve(w, r, files[0])
}

// endregion

func ImageCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageIdentifierQueryParam := chi.URLParam(r, "imageId")
		imageID, _ := strconv.Atoi(imageIdentifierQueryParam)

		var image *models.Image
		var files []*models.File
		readTxnErr := manager.GetInstance().TxnManager.WithReadTxn(r.Context(), func(repo models.ReaderRepository) error {
			qb := repo.Image()
			if imageID == 0 {
				image, _ = qb.FindByChecksum(imageIdentifierQueryParam)
			} else {
				image, _ = qb.Find(imageID)
			}

			if image != nil {
				// get the file(s) as well
				fileIDs, _ := qb.GetFileIDs(imageID)
				files, _ = repo.File().Find(fileIDs)
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
		ctx = context.WithValue(ctx, filesKey, files)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
