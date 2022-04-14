package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type PerformerFinder interface {
	Find(ctx context.Context, id int) (*models.Performer, error)
	GetImage(ctx context.Context, performerID int) ([]byte, error)
}

type performerRoutes struct {
	txnManager      txn.Manager
	performerFinder PerformerFinder
}

func (rs performerRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{performerId}", func(r chi.Router) {
		r.Use(rs.PerformerCtx)
		r.Get("/image", rs.Image)
	})

	return r
}

func (rs performerRoutes) Image(w http.ResponseWriter, r *http.Request) {
	performer := r.Context().Value(performerKey).(*models.Performer)
	defaultParam := r.URL.Query().Get("default")

	var image []byte
	if defaultParam != "true" {
		readTxnErr := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			image, _ = rs.performerFinder.GetImage(ctx, performer.ID)
			return nil
		})
		if readTxnErr != nil {
			logger.Warnf("couldn't execute getting a performer image from read transaction: %v", readTxnErr)
		}
	}

	if len(image) == 0 || defaultParam == "true" {
		image, _ = getRandomPerformerImageUsingName(performer.Name.String, performer.Gender.String, config.GetInstance().GetCustomPerformerImageLocation())
	}

	if err := utils.ServeImage(image, w, r); err != nil {
		logger.Warnf("error serving image: %v", err)
	}
}

func (rs performerRoutes) PerformerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		performerID, err := strconv.Atoi(chi.URLParam(r, "performerId"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		var performer *models.Performer
		if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			var err error
			performer, err = rs.performerFinder.Find(ctx, performerID)
			return err
		}); err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), performerKey, performer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
