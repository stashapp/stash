package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type SceneFinder interface {
	manager.SceneCoverGetter

	scene.IDFinder
	FindByChecksum(ctx context.Context, checksum string) ([]*models.Scene, error)
	FindByOSHash(ctx context.Context, oshash string) ([]*models.Scene, error)
}

type SceneMarkerFinder interface {
	Find(ctx context.Context, id int) (*models.SceneMarker, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneMarker, error)
}

type CaptionFinder interface {
	GetCaptions(ctx context.Context, fileID file.ID) ([]*models.VideoCaption, error)
}

type sceneRoutes struct {
	txnManager        txn.Manager
	sceneFinder       SceneFinder
	fileFinder        file.Finder
	captionFinder     CaptionFinder
	sceneMarkerFinder SceneMarkerFinder
	tagFinder         scene.MarkerTagFinder
}

func (rs sceneRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{sceneId}", func(r chi.Router) {
		r.Use(rs.SceneCtx)

		// streaming endpoints
		r.Get("/stream", rs.StreamDirect)
		r.Get("/stream.mp4", rs.StreamMp4)
		r.Get("/stream.webm", rs.StreamWebM)
		r.Get("/stream.mkv", rs.StreamMKV)
		r.Get("/stream.m3u8", rs.StreamHLS)
		r.Get("/stream.m3u8/{segment}.ts", rs.StreamHLSSegment)
		r.Get("/stream.mpd", rs.StreamDASH)
		r.Get("/stream.mpd/{segment}_v.webm", rs.StreamDASHVideoSegment)
		r.Get("/stream.mpd/{segment}_a.webm", rs.StreamDASHAudioSegment)

		r.Get("/screenshot", rs.Screenshot)
		r.Get("/preview", rs.Preview)
		r.Get("/webp", rs.Webp)
		r.Get("/vtt/chapter", rs.VttChapter)
		r.Get("/vtt/thumbs", rs.VttThumbs)
		r.Get("/vtt/sprite", rs.VttSprite)
		r.Get("/funscript", rs.Funscript)
		r.Get("/interactive_heatmap", rs.InteractiveHeatmap)
		r.Get("/caption", rs.CaptionLang)

		r.Get("/scene_marker/{sceneMarkerId}/stream", rs.SceneMarkerStream)
		r.Get("/scene_marker/{sceneMarkerId}/preview", rs.SceneMarkerPreview)
		r.Get("/scene_marker/{sceneMarkerId}/screenshot", rs.SceneMarkerScreenshot)
	})
	r.Get("/{sceneHash}_thumbs.vtt", rs.VttThumbs)
	r.Get("/{sceneHash}_sprite.jpg", rs.VttSprite)

	return r
}

// region Handlers

func (rs sceneRoutes) StreamDirect(w http.ResponseWriter, r *http.Request) {

	scene := r.Context().Value(sceneKey).(*models.Scene)
	// #3526 - return 404 if the scene does not have any files
	if scene.Path == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())

	filepath := manager.GetInstance().Paths.Scene.GetStreamPath(scene.Path, sceneHash)
	streamRequestCtx := ffmpeg.NewStreamRequestContext(w, r)

	// #2579 - hijacking and closing the connection here causes video playback to fail in Safari
	// We trust that the request context will be closed, so we don't need to call Cancel on the
	// returned context here.
	_ = manager.GetInstance().ReadLockManager.ReadLock(streamRequestCtx, filepath)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) StreamMp4(w http.ResponseWriter, r *http.Request) {
	rs.streamTranscode(w, r, ffmpeg.StreamTypeMP4)
}

func (rs sceneRoutes) StreamWebM(w http.ResponseWriter, r *http.Request) {
	rs.streamTranscode(w, r, ffmpeg.StreamTypeWEBM)
}

func (rs sceneRoutes) StreamMKV(w http.ResponseWriter, r *http.Request) {
	// only allow mkv streaming if the scene container is an mkv already
	scene := r.Context().Value(sceneKey).(*models.Scene)

	pf := scene.Files.Primary()
	if pf == nil {
		return
	}

	container, err := manager.GetVideoFileContainer(pf)
	if err != nil {
		logger.Errorf("[transcode] error getting container: %v", err)
	}

	if container != ffmpeg.Matroska {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("not an mkv file")); err != nil {
			logger.Warnf("[stream] error writing to stream: %v", err)
		}
		return
	}

	rs.streamTranscode(w, r, ffmpeg.StreamTypeMKV)
}

func (rs sceneRoutes) streamTranscode(w http.ResponseWriter, r *http.Request, streamType ffmpeg.StreamFormat) {
	scene := r.Context().Value(sceneKey).(*models.Scene)

	streamManager := manager.GetInstance().StreamManager
	if streamManager == nil {
		http.Error(w, "Live transcoding disabled", http.StatusServiceUnavailable)
		return
	}

	f := scene.Files.Primary()
	if f == nil {
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Warnf("[transcode] error parsing query form: %v", err)
	}

	startTime := r.Form.Get("start")
	ss, _ := strconv.ParseFloat(startTime, 64)
	resolution := r.Form.Get("resolution")

	options := ffmpeg.TranscodeOptions{
		StreamType: streamType,
		VideoFile:  f,
		Resolution: resolution,
		StartTime:  ss,
	}

	logger.Debugf("[transcode] streaming scene %d as %s", scene.ID, streamType.MimeType)
	streamManager.ServeTranscode(w, r, options)
}

func (rs sceneRoutes) StreamHLS(w http.ResponseWriter, r *http.Request) {
	rs.streamManifest(w, r, ffmpeg.StreamTypeHLS, "HLS")
}

func (rs sceneRoutes) StreamDASH(w http.ResponseWriter, r *http.Request) {
	rs.streamManifest(w, r, ffmpeg.StreamTypeDASHVideo, "DASH")
}

func (rs sceneRoutes) streamManifest(w http.ResponseWriter, r *http.Request, streamType *ffmpeg.StreamType, logName string) {
	scene := r.Context().Value(sceneKey).(*models.Scene)

	streamManager := manager.GetInstance().StreamManager
	if streamManager == nil {
		http.Error(w, "Live transcoding disabled", http.StatusServiceUnavailable)
		return
	}

	f := scene.Files.Primary()
	if f == nil {
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Warnf("[transcode] error parsing query form: %v", err)
	}

	resolution := r.Form.Get("resolution")

	logger.Debugf("[transcode] returning %s manifest for scene %d", logName, scene.ID)
	streamManager.ServeManifest(w, r, streamType, f, resolution)
}

func (rs sceneRoutes) StreamHLSSegment(w http.ResponseWriter, r *http.Request) {
	rs.streamSegment(w, r, ffmpeg.StreamTypeHLS)
}

func (rs sceneRoutes) StreamDASHVideoSegment(w http.ResponseWriter, r *http.Request) {
	rs.streamSegment(w, r, ffmpeg.StreamTypeDASHVideo)
}

func (rs sceneRoutes) StreamDASHAudioSegment(w http.ResponseWriter, r *http.Request) {
	rs.streamSegment(w, r, ffmpeg.StreamTypeDASHAudio)
}

func (rs sceneRoutes) streamSegment(w http.ResponseWriter, r *http.Request, streamType *ffmpeg.StreamType) {
	scene := r.Context().Value(sceneKey).(*models.Scene)

	streamManager := manager.GetInstance().StreamManager
	if streamManager == nil {
		http.Error(w, "Live transcoding disabled", http.StatusServiceUnavailable)
		return
	}

	f := scene.Files.Primary()
	if f == nil {
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Warnf("[transcode] error parsing query form: %v", err)
	}

	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())

	segment := chi.URLParam(r, "segment")
	resolution := r.Form.Get("resolution")

	options := ffmpeg.StreamOptions{
		StreamType: streamType,
		VideoFile:  f,
		Resolution: resolution,
		Hash:       sceneHash,
		Segment:    segment,
	}

	streamManager.ServeSegment(w, r, options)
}

func (rs sceneRoutes) Screenshot(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)

	ss := manager.SceneServer{
		TxnManager:       rs.txnManager,
		SceneCoverGetter: rs.sceneFinder,
	}
	ss.ServeScreenshot(scene, w, r)
}

func (rs sceneRoutes) Preview(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	filepath := manager.GetInstance().Paths.Scene.GetVideoPreviewPath(sceneHash)
	serveFileNoCache(w, r, filepath)
}

// serveFileNoCache serves the provided file, ensuring that the response
// contains headers to prevent caching.
func serveFileNoCache(w http.ResponseWriter, r *http.Request, filepath string) {
	w.Header().Add("Cache-Control", "no-cache")

	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) Webp(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	filepath := manager.GetInstance().Paths.Scene.GetWebpPreviewPath(sceneHash)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) getChapterVttTitle(ctx context.Context, marker *models.SceneMarker) (*string, error) {
	if marker.Title != "" {
		return &marker.Title, nil
	}

	var title string
	if err := txn.WithReadTxn(ctx, rs.txnManager, func(ctx context.Context) error {
		qb := rs.tagFinder
		primaryTag, err := qb.Find(ctx, marker.PrimaryTagID)
		if err != nil {
			return err
		}

		title = primaryTag.Name

		tags, err := qb.FindBySceneMarkerID(ctx, marker.ID)
		if err != nil {
			return err
		}

		for _, t := range tags {
			title += ", " + t.Name
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &title, nil
}

func (rs sceneRoutes) VttChapter(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	var sceneMarkers []*models.SceneMarker
	readTxnErr := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		sceneMarkers, err = rs.sceneMarkerFinder.FindBySceneID(ctx, scene.ID)
		return err
	})
	if errors.Is(readTxnErr, context.Canceled) {
		return
	}
	if readTxnErr != nil {
		logger.Warnf("read transaction error on fetch scene markers: %v", readTxnErr)
		http.Error(w, readTxnErr.Error(), http.StatusInternalServerError)
		return
	}

	vttLines := []string{"WEBVTT", ""}
	for i, marker := range sceneMarkers {
		vttLines = append(vttLines, strconv.Itoa(i+1))
		time := utils.GetVTTTime(marker.Seconds)
		vttLines = append(vttLines, time+" --> "+time)

		vttTitle, err := rs.getChapterVttTitle(r.Context(), marker)
		if errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			logger.Warnf("read transaction error on fetch scene marker title: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vttLines = append(vttLines, *vttTitle)
		vttLines = append(vttLines, "")
	}
	vtt := strings.Join(vttLines, "\n")

	w.Header().Set("Content-Type", "text/vtt")
	_, _ = w.Write([]byte(vtt))
}

func (rs sceneRoutes) VttThumbs(w http.ResponseWriter, r *http.Request) {
	scene, ok := r.Context().Value(sceneKey).(*models.Scene)
	var sceneHash string
	if ok && scene != nil {
		sceneHash = scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	} else {
		sceneHash = chi.URLParam(r, "sceneHash")
	}
	w.Header().Set("Content-Type", "text/vtt")
	filepath := manager.GetInstance().Paths.Scene.GetSpriteVttFilePath(sceneHash)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) VttSprite(w http.ResponseWriter, r *http.Request) {
	scene, ok := r.Context().Value(sceneKey).(*models.Scene)
	var sceneHash string
	if ok && scene != nil {
		sceneHash = scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	} else {
		sceneHash = chi.URLParam(r, "sceneHash")
	}
	w.Header().Set("Content-Type", "image/jpeg")
	filepath := manager.GetInstance().Paths.Scene.GetSpriteImageFilePath(sceneHash)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) Funscript(w http.ResponseWriter, r *http.Request) {
	s := r.Context().Value(sceneKey).(*models.Scene)
	funscript := video.GetFunscriptPath(s.Path)
	serveFileNoCache(w, r, funscript)
}

func (rs sceneRoutes) InteractiveHeatmap(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	w.Header().Set("Content-Type", "image/png")
	filepath := manager.GetInstance().Paths.Scene.GetInteractiveHeatmapPath(sceneHash)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) Caption(w http.ResponseWriter, r *http.Request, lang string, ext string) {
	s := r.Context().Value(sceneKey).(*models.Scene)

	var captions []*models.VideoCaption
	readTxnErr := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
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
		logger.Warnf("read transaction error on fetch scene captions: %v", readTxnErr)
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

		var b bytes.Buffer
		err = sub.WriteToWebVTT(&b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/vtt")
		w.Header().Add("Cache-Control", "no-cache")
		_, _ = b.WriteTo(w)
	}
}

func (rs sceneRoutes) CaptionLang(w http.ResponseWriter, r *http.Request) {
	// serve caption based on lang query param, if provided
	if err := r.ParseForm(); err != nil {
		logger.Warnf("[caption] error parsing query form: %v", err)
	}

	l := r.Form.Get("lang")
	ext := r.Form.Get("type")
	rs.Caption(w, r, l, ext)
}

func (rs sceneRoutes) SceneMarkerStream(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	sceneMarkerID, _ := strconv.Atoi(chi.URLParam(r, "sceneMarkerId"))
	var sceneMarker *models.SceneMarker
	readTxnErr := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		sceneMarker, err = rs.sceneMarkerFinder.Find(ctx, sceneMarkerID)
		return err
	})
	if errors.Is(readTxnErr, context.Canceled) {
		return
	}
	if readTxnErr != nil {
		logger.Warnf("read transaction error on fetch scene marker: %v", readTxnErr)
		http.Error(w, readTxnErr.Error(), http.StatusInternalServerError)
		return
	}

	if sceneMarker == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	filepath := manager.GetInstance().Paths.SceneMarkers.GetVideoPreviewPath(sceneHash, int(sceneMarker.Seconds))
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) SceneMarkerPreview(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	sceneMarkerID, _ := strconv.Atoi(chi.URLParam(r, "sceneMarkerId"))
	var sceneMarker *models.SceneMarker
	readTxnErr := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		sceneMarker, err = rs.sceneMarkerFinder.Find(ctx, sceneMarkerID)
		return err
	})
	if errors.Is(readTxnErr, context.Canceled) {
		return
	}
	if readTxnErr != nil {
		logger.Warnf("read transaction error on fetch scene marker preview: %v", readTxnErr)
		http.Error(w, readTxnErr.Error(), http.StatusInternalServerError)
		return
	}

	if sceneMarker == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	filepath := manager.GetInstance().Paths.SceneMarkers.GetWebpPreviewPath(sceneHash, int(sceneMarker.Seconds))

	// If the image doesn't exist, send the placeholder
	exists, _ := fsutil.FileExists(filepath)
	if !exists {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(utils.PendingGenerateResource)
		return
	}

	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) SceneMarkerScreenshot(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	sceneMarkerID, _ := strconv.Atoi(chi.URLParam(r, "sceneMarkerId"))
	var sceneMarker *models.SceneMarker
	readTxnErr := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		sceneMarker, err = rs.sceneMarkerFinder.Find(ctx, sceneMarkerID)
		return err
	})
	if errors.Is(readTxnErr, context.Canceled) {
		return
	}
	if readTxnErr != nil {
		logger.Warnf("read transaction error on fetch scene marker screenshot: %v", readTxnErr)
		http.Error(w, readTxnErr.Error(), http.StatusInternalServerError)
		return
	}

	if sceneMarker == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	filepath := manager.GetInstance().Paths.SceneMarkers.GetScreenshotPath(sceneHash, int(sceneMarker.Seconds))

	// If the image doesn't exist, send the placeholder
	exists, _ := fsutil.FileExists(filepath)
	if !exists {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(utils.PendingGenerateResource)
		return
	}

	http.ServeFile(w, r, filepath)
}

// endregion

func (rs sceneRoutes) SceneCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sceneID, err := strconv.Atoi(chi.URLParam(r, "sceneId"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var scene *models.Scene
		_ = txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			qb := rs.sceneFinder
			scene, _ = qb.Find(ctx, sceneID)

			if scene != nil {
				if err := scene.LoadPrimaryFile(ctx, rs.fileFinder); err != nil {
					if !errors.Is(err, context.Canceled) {
						logger.Errorf("error loading primary file for scene %d: %v", sceneID, err)
					}
					// set scene to nil so that it doesn't try to use the primary file
					scene = nil
				}
			}

			return nil
		})
		if scene == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), sceneKey, scene)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
