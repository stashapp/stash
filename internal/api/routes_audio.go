// TODO(audio): update this file
package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type AudioFinder interface {
	models.AudioGetter

	FindByChecksum(ctx context.Context, checksum string) ([]*models.Audio, error)
	FindByOSHash(ctx context.Context, oshash string) ([]*models.Audio, error)
}

type audioRoutes struct {
	routes
	audioFinder   AudioFinder
	fileGetter    models.FileGetter
	captionFinder CaptionFinder
}

func (rs audioRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{audioId}", func(r chi.Router) {
		r.Use(rs.AudioCtx)

		// streaming endpoints
		r.Get("/stream", rs.StreamDirect)
		// TODO(audio): slightly difficult to support StreamHLS/StreamDASH...do last
		// r.Get("/stream.m3u8", rs.StreamHLS)
		// r.Get("/stream.m3u8/{segment}.ts", rs.StreamHLSSegment)
		// r.Get("/stream.mpd", rs.StreamDASH)
		// r.Get("/stream.mpd/{segment}_a.webm", rs.StreamDASHAudioSegment)

		r.Get("/funscript", rs.Funscript)
		r.Get("/caption", rs.CaptionLang)
	})

	return r
}

func (rs audioRoutes) StreamDirect(w http.ResponseWriter, r *http.Request) {
	audio := r.Context().Value(audioKey).(*models.Audio)
	ss := manager.AudioServer{
		TxnManager: rs.txnManager,
	}
	ss.StreamAudioDirect(audio, w, r)
}

// func (rs audioRoutes) StreamHLS(w http.ResponseWriter, r *http.Request) {
// 	rs.streamManifest(w, r, ffmpeg.StreamTypeHLS, "HLS")
// }

// func (rs audioRoutes) StreamDASH(w http.ResponseWriter, r *http.Request) {
// 	rs.streamManifest(w, r, ffmpeg.StreamTypeDASHAudio, "DASH")
// }

// func (rs audioRoutes) streamManifest(w http.ResponseWriter, r *http.Request, streamType *ffmpeg.StreamType, logName string) {
// 	audio := r.Context().Value(audioKey).(*models.Audio)

// 	streamManager := manager.GetInstance().StreamManager
// 	if streamManager == nil {
// 		http.Error(w, "Live transcoding disabled", http.StatusServiceUnavailable)
// 		return
// 	}

// 	f := audio.Files.Primary()
// 	if f == nil {
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		logger.Warnf("[transcode] error parsing query form: %v", err)
// 	}

// 	resolution := r.Form.Get("resolution")

// 	logger.Debugf("[transcode] returning %s manifest for audio %d", logName, audio.ID)
// 	streamManager.ServeManifest(w, r, streamType, f, resolution)
// }

// func (rs audioRoutes) StreamHLSSegment(w http.ResponseWriter, r *http.Request) {
// 	rs.streamSegment(w, r, ffmpeg.StreamTypeHLS)
// }

// func (rs audioRoutes) StreamDASHAudioSegment(w http.ResponseWriter, r *http.Request) {
// 	rs.streamSegment(w, r, ffmpeg.StreamTypeDASHAudio)
// }

// func (rs audioRoutes) streamSegment(w http.ResponseWriter, r *http.Request, streamType *ffmpeg.StreamType) {
// 	audio := r.Context().Value(audioKey).(*models.Audio)

// 	streamManager := manager.GetInstance().StreamManager
// 	if streamManager == nil {
// 		http.Error(w, "Live transcoding disabled", http.StatusServiceUnavailable)
// 		return
// 	}

// 	f := audio.Files.Primary()
// 	if f == nil {
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		logger.Warnf("[transcode] error parsing query form: %v", err)
// 	}

// 	audioHash := audio.GetHash(config.GetInstance().GetAudioFileNamingAlgorithm())

// 	segment := chi.URLParam(r, "segment")
// 	resolution := r.Form.Get("resolution")

// 	options := ffmpeg.StreamOptions{
// 		StreamType: streamType,
// 		AudioFile:  f,
// 		Resolution: resolution,
// 		Hash:       audioHash,
// 		Segment:    segment,
// 	}

// 	streamManager.ServeSegment(w, r, options)
// }

func (rs audioRoutes) Funscript(w http.ResponseWriter, r *http.Request) {
	s := r.Context().Value(audioKey).(*models.Audio)
	filepath := video.GetFunscriptPath(s.Path)

	utils.ServeStaticFile(w, r, filepath)
}

func (rs audioRoutes) Caption(w http.ResponseWriter, r *http.Request, lang string, ext string) {
	s := r.Context().Value(audioKey).(*models.Audio)

	var captions []*models.VideoCaption
	readTxnErr := rs.withReadTxn(r, func(ctx context.Context) error {
		var err error
		primaryFile := s.Files.Primary()
		if primaryFile == nil {
			return nil
		}

		captions, err = rs.captionFinder.GetCaptions(ctx, primaryFile.Base().ID)

		return err
	})
	if errors.Is(readTxnErr, context.Canceled) {
		return
	}
	if readTxnErr != nil {
		logger.Warnf("read transaction error on fetch audio captions: %v", readTxnErr)
		http.Error(w, readTxnErr.Error(), http.StatusInternalServerError)
		return
	}

	for _, caption := range captions {
		if lang != caption.LanguageCode || ext != caption.CaptionType {
			continue
		}

		sub, err := video.ReadSubs(caption.Path(s.Path))
		if err != nil {
			logger.Warnf("error while reading subs: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer

		err = sub.WriteToWebVTT(&buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/vtt")
		utils.ServeStaticContent(w, r, buf.Bytes())
		return
	}
}

func (rs audioRoutes) CaptionLang(w http.ResponseWriter, r *http.Request) {
	// serve caption based on lang query param, if provided
	if err := r.ParseForm(); err != nil {
		logger.Warnf("[caption] error parsing query form: %v", err)
	}

	l := r.Form.Get("lang")
	ext := r.Form.Get("type")
	rs.Caption(w, r, l, ext)
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
