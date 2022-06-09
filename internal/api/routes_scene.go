package api

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
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
	FindByChecksum(ctx context.Context, checksum string) (*models.Scene, error)
	FindByOSHash(ctx context.Context, oshash string) (*models.Scene, error)
	GetCaptions(ctx context.Context, sceneID int) ([]*models.SceneCaption, error)
}

type SceneMarkerFinder interface {
	Find(ctx context.Context, id int) (*models.SceneMarker, error)
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneMarker, error)
}

type sceneRoutes struct {
	txnManager        txn.Manager
	sceneFinder       SceneFinder
	sceneMarkerFinder SceneMarkerFinder
	tagFinder         scene.MarkerTagFinder
}

func (rs sceneRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{sceneId}", func(r chi.Router) {
		r.Use(rs.SceneCtx)

		// streaming endpoints
		r.Get("/stream", rs.StreamDirect)
		r.Get("/stream.mkv", rs.StreamMKV)
		r.Get("/stream.webm", rs.StreamWebM)
		r.Get("/stream.m3u8", rs.StreamHLS)
		r.Get("/stream.ts", rs.StreamTS)
		r.Get("/stream.mp4", rs.StreamMp4)

		r.Get("/screenshot", rs.Screenshot)
		r.Get("/preview", rs.Preview)
		r.Get("/webp", rs.Webp)
		r.Get("/vtt/chapter", rs.ChapterVtt)
		r.Get("/funscript", rs.Funscript)
		r.Get("/interactive_heatmap", rs.InteractiveHeatmap)
		r.Get("/caption", rs.CaptionLang)

		r.Get("/scene_marker/{sceneMarkerId}/stream", rs.SceneMarkerStream)
		r.Get("/scene_marker/{sceneMarkerId}/preview", rs.SceneMarkerPreview)
		r.Get("/scene_marker/{sceneMarkerId}/screenshot", rs.SceneMarkerScreenshot)
	})
	r.With(rs.SceneCtx).Get("/{sceneId}_thumbs.vtt", rs.VttThumbs)
	r.With(rs.SceneCtx).Get("/{sceneId}_sprite.jpg", rs.VttSprite)

	return r
}

// region Handlers

func (rs sceneRoutes) StreamDirect(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)

	ss := manager.SceneServer{
		TxnManager:       rs.txnManager,
		SceneCoverGetter: rs.sceneFinder,
	}
	ss.StreamSceneDirect(scene, w, r)
}

func (rs sceneRoutes) StreamMKV(w http.ResponseWriter, r *http.Request) {
	// only allow mkv streaming if the scene container is an mkv already
	scene := r.Context().Value(sceneKey).(*models.Scene)

	container, err := manager.GetSceneFileContainer(scene)
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

	rs.streamTranscode(w, r, ffmpeg.StreamFormatMKVAudio)
}

func (rs sceneRoutes) StreamWebM(w http.ResponseWriter, r *http.Request) {
	rs.streamTranscode(w, r, ffmpeg.StreamFormatVP9)
}

func (rs sceneRoutes) StreamMp4(w http.ResponseWriter, r *http.Request) {
	rs.streamTranscode(w, r, ffmpeg.StreamFormatH264)
}

func (rs sceneRoutes) StreamHLS(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)

	ffprobe := manager.GetInstance().FFProbe
	videoFile, err := ffprobe.NewVideoFile(scene.Path)
	if err != nil {
		logger.Errorf("[stream] error reading video file: %v", err)
		return
	}

	logger.Debug("Returning HLS playlist")

	// getting the playlist manifest only
	w.Header().Set("Content-Type", ffmpeg.MimeHLS)
	var str strings.Builder

	ffmpeg.WriteHLSPlaylist(videoFile.Duration, r.URL.String(), &str)

	requestByteRange := createByteRange(r.Header.Get("Range"))
	if requestByteRange.RawString != "" {
		logger.Debugf("Requested range: %s", requestByteRange.RawString)
	}

	ret := requestByteRange.apply([]byte(str.String()))
	rangeStr := requestByteRange.toHeaderValue(int64(str.Len()))
	w.Header().Set("Content-Range", rangeStr)

	if n, err := w.Write(ret); err != nil {
		logger.Warnf("[stream] error writing stream (wrote %v bytes): %v", n, err)
	}
}

func (rs sceneRoutes) StreamTS(w http.ResponseWriter, r *http.Request) {
	rs.streamTranscode(w, r, ffmpeg.StreamFormatHLS)
}

func (rs sceneRoutes) streamTranscode(w http.ResponseWriter, r *http.Request, streamFormat ffmpeg.StreamFormat) {
	logger.Debugf("Streaming as %s", streamFormat.MimeType)
	scene := r.Context().Value(sceneKey).(*models.Scene)

	// start stream based on query param, if provided
	if err := r.ParseForm(); err != nil {
		logger.Warnf("[stream] error parsing query form: %v", err)
	}

	startTime := r.Form.Get("start")
	ss, _ := strconv.ParseFloat(startTime, 64)
	requestedSize := r.Form.Get("resolution")

	audioCodec := ffmpeg.MissingUnsupported
	if scene.AudioCodec != nil {
		audioCodec = ffmpeg.ProbeAudioCodec(*scene.AudioCodec)
	}

	var (
		width  int
		height int
	)

	if scene.Width != nil {
		width = *scene.Width
	}
	if scene.Height != nil {
		height = *scene.Height
	}

	options := ffmpeg.TranscodeStreamOptions{
		Input:     scene.Path,
		Codec:     streamFormat,
		VideoOnly: audioCodec == ffmpeg.MissingUnsupported,

		VideoWidth:  width,
		VideoHeight: height,

		StartTime:        ss,
		MaxTranscodeSize: config.GetInstance().GetMaxStreamingTranscodeSize().GetMaxResolution(),
	}

	if requestedSize != "" {
		options.MaxTranscodeSize = models.StreamingResolutionEnum(requestedSize).GetMaxResolution()
	}

	encoder := manager.GetInstance().FFMPEG

	lm := manager.GetInstance().ReadLockManager
	streamRequestCtx := manager.NewStreamRequestContext(w, r)
	lockCtx := lm.ReadLock(streamRequestCtx, scene.Path)
	defer lockCtx.Cancel()

	stream, err := encoder.GetTranscodeStream(lockCtx, options)

	if err != nil {
		logger.Errorf("[stream] error transcoding video file: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			logger.Warnf("[stream] error writing response: %v", err)
		}
		return
	}

	lockCtx.AttachCommand(stream.Cmd)

	stream.Serve(w, r)
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
	filepath := manager.GetInstance().Paths.Scene.GetVideoPreviewPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))
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
	filepath := manager.GetInstance().Paths.Scene.GetWebpPreviewPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) getChapterVttTitle(ctx context.Context, marker *models.SceneMarker) string {
	if marker.Title != "" {
		return marker.Title
	}

	var ret string
	if err := txn.WithTxn(ctx, rs.txnManager, func(ctx context.Context) error {
		qb := rs.tagFinder
		primaryTag, err := qb.Find(ctx, marker.PrimaryTagID)
		if err != nil {
			return err
		}

		ret = primaryTag.Name

		tags, err := qb.FindBySceneMarkerID(ctx, marker.ID)
		if err != nil {
			return err
		}

		for _, t := range tags {
			ret += ", " + t.Name
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return ret
}

func (rs sceneRoutes) ChapterVtt(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	var sceneMarkers []*models.SceneMarker
	if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		sceneMarkers, err = rs.sceneMarkerFinder.FindBySceneID(ctx, scene.ID)
		return err
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vttLines := []string{"WEBVTT", ""}
	for i, marker := range sceneMarkers {
		vttLines = append(vttLines, strconv.Itoa(i+1))
		time := utils.GetVTTTime(marker.Seconds)
		vttLines = append(vttLines, time+" --> "+time)
		vttLines = append(vttLines, rs.getChapterVttTitle(r.Context(), marker))
		vttLines = append(vttLines, "")
	}
	vtt := strings.Join(vttLines, "\n")

	w.Header().Set("Content-Type", "text/vtt")
	_, _ = w.Write([]byte(vtt))
}

func (rs sceneRoutes) Funscript(w http.ResponseWriter, r *http.Request) {
	s := r.Context().Value(sceneKey).(*models.Scene)
	funscript := scene.GetFunscriptPath(s.Path)
	serveFileNoCache(w, r, funscript)
}

func (rs sceneRoutes) InteractiveHeatmap(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	w.Header().Set("Content-Type", "image/png")
	filepath := manager.GetInstance().Paths.Scene.GetInteractiveHeatmapPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) Caption(w http.ResponseWriter, r *http.Request, lang string, ext string) {
	s := r.Context().Value(sceneKey).(*models.Scene)

	if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		captions, err := rs.sceneFinder.GetCaptions(ctx, s.ID)
		for _, caption := range captions {
			if lang == caption.LanguageCode && ext == caption.CaptionType {
				sub, err := scene.ReadSubs(caption.Path(s.Path))
				if err == nil {
					var b bytes.Buffer
					err = sub.WriteToWebVTT(&b)
					if err == nil {
						w.Header().Set("Content-Type", "text/vtt")
						w.Header().Add("Cache-Control", "no-cache")
						_, _ = b.WriteTo(w)
					}
					return err
				}
				logger.Debugf("Error while reading subs: %v", err)
			}
		}
		return err
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func (rs sceneRoutes) VttThumbs(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	w.Header().Set("Content-Type", "text/vtt")
	filepath := manager.GetInstance().Paths.Scene.GetSpriteVttFilePath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) VttSprite(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	w.Header().Set("Content-Type", "image/jpeg")
	filepath := manager.GetInstance().Paths.Scene.GetSpriteImageFilePath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()))
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) SceneMarkerStream(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneMarkerID, _ := strconv.Atoi(chi.URLParam(r, "sceneMarkerId"))
	var sceneMarker *models.SceneMarker
	if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		sceneMarker, err = rs.sceneMarkerFinder.Find(ctx, sceneMarkerID)
		return err
	}); err != nil {
		logger.Warnf("Error when getting scene marker for stream: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if sceneMarker == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	filepath := manager.GetInstance().Paths.SceneMarkers.GetVideoPreviewPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()), int(sceneMarker.Seconds))
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) SceneMarkerPreview(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneMarkerID, _ := strconv.Atoi(chi.URLParam(r, "sceneMarkerId"))
	var sceneMarker *models.SceneMarker
	if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		sceneMarker, err = rs.sceneMarkerFinder.Find(ctx, sceneMarkerID)
		return err
	}); err != nil {
		logger.Warnf("Error when getting scene marker for stream: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if sceneMarker == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	filepath := manager.GetInstance().Paths.SceneMarkers.GetWebpPreviewPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()), int(sceneMarker.Seconds))

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
	sceneMarkerID, _ := strconv.Atoi(chi.URLParam(r, "sceneMarkerId"))
	var sceneMarker *models.SceneMarker
	if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		sceneMarker, err = rs.sceneMarkerFinder.Find(ctx, sceneMarkerID)
		return err
	}); err != nil {
		logger.Warnf("Error when getting scene marker for stream: %s", err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if sceneMarker == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	filepath := manager.GetInstance().Paths.SceneMarkers.GetScreenshotPath(scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm()), int(sceneMarker.Seconds))

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
		sceneIdentifierQueryParam := chi.URLParam(r, "sceneId")
		sceneID, _ := strconv.Atoi(sceneIdentifierQueryParam)

		var scene *models.Scene
		readTxnErr := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			qb := rs.sceneFinder
			if sceneID == 0 {
				// determine checksum/os by the length of the query param
				if len(sceneIdentifierQueryParam) == 32 {
					scene, _ = qb.FindByChecksum(ctx, sceneIdentifierQueryParam)
				} else {
					scene, _ = qb.FindByOSHash(ctx, sceneIdentifierQueryParam)
				}
			} else {
				scene, _ = qb.Find(ctx, sceneID)
			}

			return nil
		})
		if readTxnErr != nil {
			logger.Warnf("error executing SceneCtx transaction: %v", readTxnErr)
		}

		if scene == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), sceneKey, scene)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
