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

type GroupFinder interface {
	models.GroupGetter
	GetFrontImage(ctx context.Context, groupID int) ([]byte, error)
	GetBackImage(ctx context.Context, groupID int) ([]byte, error)
}

type groupRoutes struct {
	routes
	groupFinder GroupFinder
}

func (rs groupRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{groupId}", func(r chi.Router) {
		r.Use(rs.GroupCtx)
		r.Get("/frontimage", rs.FrontImage)
		r.Get("/backimage", rs.BackImage)
	})

	return r
}

func (rs groupRoutes) FrontImage(w http.ResponseWriter, r *http.Request) {
	group := r.Context().Value(groupKey).(*models.Group)
	defaultParam := r.URL.Query().Get("default")
	var image []byte
	if defaultParam != "true" {
		readTxnErr := rs.withReadTxn(r, func(ctx context.Context) error {
			var err error
			image, err = rs.groupFinder.GetFrontImage(ctx, group.ID)
			return err
		})
		if errors.Is(readTxnErr, context.Canceled) {
			return
		}
		if readTxnErr != nil {
			logger.Warnf("read transaction error on fetch group front image: %v", readTxnErr)
		}
	}

	// fallback to default image
	if len(image) == 0 {
		image = static.ReadAll(static.DefaultGroupImage)
	}

	utils.ServeImage(w, r, image)
}

func (rs groupRoutes) BackImage(w http.ResponseWriter, r *http.Request) {
	group := r.Context().Value(groupKey).(*models.Group)
	defaultParam := r.URL.Query().Get("default")
	var image []byte
	if defaultParam != "true" {
		readTxnErr := rs.withReadTxn(r, func(ctx context.Context) error {
			var err error
			image, err = rs.groupFinder.GetBackImage(ctx, group.ID)
			return err
		})
		if errors.Is(readTxnErr, context.Canceled) {
			return
		}
		if readTxnErr != nil {
			logger.Warnf("read transaction error on fetch group back image: %v", readTxnErr)
		}
	}

	// fallback to default image
	if len(image) == 0 {
		image = static.ReadAll(static.DefaultGroupImage)
	}

	utils.ServeImage(w, r, image)
}

func (rs groupRoutes) GroupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupID, err := strconv.Atoi(chi.URLParam(r, "groupId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var group *models.Group
		_ = rs.withReadTxn(r, func(ctx context.Context) error {
			group, _ = rs.groupFinder.Find(ctx, groupID)
			return nil
		})
		if group == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), groupKey, group)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
