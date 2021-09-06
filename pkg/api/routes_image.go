package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
	image := r.Context().Value(imageKey).(*models.Image)
	filepath := manager.GetInstance().Paths.Generated.GetThumbnailPath(image.Checksum, models.DefaultGthumbWidth)

	// if the thumbnail doesn't exist, fall back to the original file
	exists, _ := utils.FileExists(filepath)
	if exists {
		http.ServeFile(w, r, filepath)
	} else {
		rs.Image(w, r)
	}
}

func (rs imageRoutes) Image(w http.ResponseWriter, r *http.Request) {
	i := r.Context().Value(imageKey).(*models.Image)

	// if image is in a zip file, we need to serve it specifically
	image.Serve(w, r, i.Path)
}

// endregion

func ImageCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageIdentifierQueryParam := chi.URLParam(r, "imageId")
		imageID, _ := strconv.Atoi(imageIdentifierQueryParam)

		var image *models.Image
		readTxnErr := manager.GetInstance().TxnManager.WithReadTxn(r.Context(), func(repo models.ReaderRepository) error {
			qb := repo.Image()
			if imageID == 0 {
				image, _ = qb.FindByChecksum(imageIdentifierQueryParam)
			} else {
				image, _ = qb.Find(imageID)
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
