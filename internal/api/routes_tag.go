package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/tag"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type TagFinder interface {
	tag.Finder
	GetImage(ctx context.Context, tagID int) ([]byte, error)
}

type tagRoutes struct {
	txnManager txn.Manager
	tagFinder  TagFinder
}

func (rs tagRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{tagId}", func(r chi.Router) {
		r.Use(rs.TagCtx)
		r.Get("/image", rs.Image)
	})

	return r
}

func (rs tagRoutes) Image(w http.ResponseWriter, r *http.Request) {
	tag := r.Context().Value(tagKey).(*models.Tag)
	defaultParam := r.URL.Query().Get("default")

	var image []byte
	if defaultParam != "true" {
		err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			image, _ = rs.tagFinder.GetImage(ctx, tag.ID)
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

func (rs tagRoutes) TagCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(chi.URLParam(r, "tagId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var tag *models.Tag
		if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			var err error
			tag, err = rs.tagFinder.Find(ctx, tagID)
			return err
		}); err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), tagKey, tag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
