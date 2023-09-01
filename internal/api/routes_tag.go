package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type TagFinder interface {
	models.TagGetter
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
		readTxnErr := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			var err error
			image, err = rs.tagFinder.GetImage(ctx, tag.ID)
			return err
		})
		if errors.Is(readTxnErr, context.Canceled) {
			return
		}
		if readTxnErr != nil {
			logger.Warnf("read transaction error on fetch tag image: %v", readTxnErr)
		}
	}

	if len(image) == 0 {
		const defaultTagImage = "tag/tag.svg"

		// fall back to static image
		f, _ := static.Tag.Open(defaultTagImage)
		defer f.Close()
		stat, _ := f.Stat()
		http.ServeContent(w, r, "tag.svg", stat.ModTime(), f.(io.ReadSeeker))
		return
	}

	utils.ServeImage(w, r, image)
}

func (rs tagRoutes) TagCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(chi.URLParam(r, "tagId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var tag *models.Tag
		_ = txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			var err error
			tag, err = rs.tagFinder.Find(ctx, tagID)
			return err
		})
		if tag == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), tagKey, tag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
