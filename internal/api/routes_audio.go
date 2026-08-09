package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/static"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type AudioFinder interface {
	models.AudioGetter

	FindByChecksum(ctx context.Context, checksum string) ([]*models.Audio, error)
	GetCover(ctx context.Context, audioID int) ([]byte, error)
}

type audioRoutes struct {
	routes
	audioFinder AudioFinder
	fileGetter  models.FileGetter
}

func (rs audioRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{audioId}", func(r chi.Router) {
		r.Use(rs.AudioCtx)

		r.Get("/screenshot", rs.Screenshot)

		// streaming endpoints
		r.Get("/stream", rs.StreamDirect)
	})

	return r
}

func (rs audioRoutes) Screenshot(w http.ResponseWriter, r *http.Request) {
	// if default flag is set, return the default image
	if r.URL.Query().Get("default") == "true" {
		utils.ServeImage(w, r, static.ReadAll(static.DefaultSceneImage))
		return
	}

	audio := r.Context().Value(audioKey).(*models.Audio)

	ss := manager.AudioServer{
		TxnManager:       rs.txnManager,
		AudioCoverGetter: rs.audioFinder,
	}
	ss.ServeScreenshot(audio, w, r)
}

func (rs audioRoutes) StreamDirect(w http.ResponseWriter, r *http.Request) {
	audio := r.Context().Value(audioKey).(*models.Audio)
	ss := manager.AudioServer{
		TxnManager: rs.txnManager,
	}
	ss.StreamAudioDirect(audio, w, r)
}

func (rs audioRoutes) AudioCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		audioID, err := strconv.Atoi(chi.URLParam(r, "audioId"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var audio *models.Audio
		_ = rs.withReadTxn(r, func(ctx context.Context) error {
			qb := rs.audioFinder
			audio, _ = qb.Find(ctx, audioID)

			if audio != nil {
				if err := audio.LoadPrimaryFile(ctx, rs.fileGetter); err != nil {
					if !errors.Is(err, context.Canceled) {
						logger.Errorf("error loading primary file for audio %d: %v", audioID, err)
					}
					// set audio to nil so that it doesn't try to use the primary file
					audio = nil
				}
			}

			return nil
		})
		if audio == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), audioKey, audio)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
