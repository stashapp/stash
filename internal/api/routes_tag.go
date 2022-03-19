package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type tagRoutes struct {
	txnManager models.TransactionManager
}

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
	defaultParam := r.URL.Query().Get("default")

	var image []byte
	if defaultParam != "true" {
		err := rs.txnManager.WithReadTxn(r.Context(), func(repo models.ReaderRepository) error {
			image, _ = repo.Tag().GetImage(tag.ID)
			return nil
		})
		if err != nil {
			logger.Warnf("read transaction error while getting tag image: %v", err)
		}
	}

	if len(image) == 0 {
		image = models.DefaultTagImage
	}

	if err := utils.ServeImage(image, w, r); err != nil {
		logger.Warnf("error serving tag image: %v", err)
	}
}

func TagCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(chi.URLParam(r, "tagId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var tag *models.Tag
		if err := manager.GetInstance().TxnManager.WithReadTxn(r.Context(), func(repo models.ReaderRepository) error {
			var err error
			tag, err = repo.Tag().Find(tagID)
			return err
		}); err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), tagKey, tag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
