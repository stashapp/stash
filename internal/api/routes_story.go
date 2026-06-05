package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type StoryFinder interface {
	models.StoryGetter
	GetFrontImage(ctx context.Context, storyID int) ([]byte, error)
	GetBackImage(ctx context.Context, storyID int) ([]byte, error)
}

type storyRoutes struct {
	routes
	storyFinder StoryFinder
	sfwConfig   sfwConfig
}

func (rs storyRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{storyId}", func(r chi.Router) {
		r.Use(rs.StoryCtx)
		r.Get("/front_cover", rs.FrontCover)
		r.Get("/back_cover", rs.BackCover)
	})

	return r
}

func (rs storyRoutes) FrontCover(w http.ResponseWriter, r *http.Request) {
	story := r.Context().Value(storyKey).(*models.Story)
	var image []byte
	readTxnErr := rs.withReadTxn(r, func(ctx context.Context) error {
		var err error
		image, err = rs.storyFinder.GetFrontImage(ctx, story.ID)
		return err
	})
	if errors.Is(readTxnErr, context.Canceled) {
		return
	}
	if readTxnErr != nil {
		logger.Warnf("read transaction error on fetch story front cover: %v", readTxnErr)
	}
	if len(image) == 0 {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	utils.ServeImage(w, r, image)
}

func (rs storyRoutes) BackCover(w http.ResponseWriter, r *http.Request) {
	story := r.Context().Value(storyKey).(*models.Story)
	var image []byte
	readTxnErr := rs.withReadTxn(r, func(ctx context.Context) error {
		var err error
		image, err = rs.storyFinder.GetBackImage(ctx, story.ID)
		return err
	})
	if errors.Is(readTxnErr, context.Canceled) {
		return
	}
	if readTxnErr != nil {
		logger.Warnf("read transaction error on fetch story back cover: %v", readTxnErr)
	}
	if len(image) == 0 {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	utils.ServeImage(w, r, image)
}

func (rs storyRoutes) StoryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storyID, err := strconv.Atoi(chi.URLParam(r, "storyId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var story *models.Story
		_ = rs.withReadTxn(r, func(ctx context.Context) error {
			var err error
			story, err = rs.storyFinder.Find(ctx, storyID)
			return err
		})
		if story == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), storyKey, story)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
