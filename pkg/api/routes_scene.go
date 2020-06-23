package api

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type sceneRoutes struct{}

func (rs sceneRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{sceneId}", func(r chi.Router) {
		r.Use(SceneCtx)
		r.Get("/stream", rs.Stream)
		r.Get("/stream.mp4", rs.Stream)
		r.Get("/screenshot", rs.Screenshot)
		r.Get("/preview", rs.Preview)
		r.Get("/webp", rs.Webp)
		r.Get("/vtt/chapter", rs.ChapterVtt)

		r.Get("/scene_marker/{sceneMarkerId}/stream", rs.SceneMarkerStream)
		r.Get("/scene_marker/{sceneMarkerId}/preview", rs.SceneMarkerPreview)
	})
	r.With(SceneCtx).Get("/{sceneId}_thumbs.vtt", rs.VttThumbs)
	r.With(SceneCtx).Get("/{sceneId}_sprite.jpg", rs.VttSprite)

	return r
}

// region Handlers

func (rs sceneRoutes) Stream(w http.ResponseWriter, r *http.Request) {

	scene := r.Context().Value(sceneKey).(*models.Scene)

	container := ""
	if scene.Format.Valid {
		container = scene.Format.String
	} else { // container isn't in the DB
		// shouldn't happen, fallback to ffprobe
		tmpVideoFile, err := ffmpeg.NewVideoFile(manager.GetInstance().FFProbePath, scene.Path)
		if err != nil {
			logger.Errorf("[transcode] error reading video file: %s", err.Error())
			return
		}

		container = string(ffmpeg.MatchContainer(tmpVideoFile.Container, scene.Path))
	}

	// detect if not a streamable file and try to transcode it instead
	filepath := manager.GetInstance().Paths.Scene.GetStreamPath(scene.Path, scene.Checksum)
	videoCodec := scene.VideoCodec.String
	audioCodec := ffmpeg.MissingUnsupported
	if scene.AudioCodec.Valid {
		audioCodec = ffmpeg.AudioCodec(scene.AudioCodec.String)
	}
	hasTranscode, _ := manager.HasTranscode(scene)
	if ffmpeg.IsValidCodec(videoCodec) && ffmpeg.IsValidCombo(videoCodec, ffmpeg.Container(container)) && ffmpeg.IsValidAudioForContainer(audioCodec, ffmpeg.Container(container)) || hasTranscode {
		manager.RegisterStream(filepath, &w)
		http.ServeFile(w, r, filepath)
		manager.WaitAndDeregisterStream(filepath, &w, r)

		return
	}

	// needs to be transcoded
	videoFile, err := ffmpeg.NewVideoFile(manager.GetInstance().FFProbePath, scene.Path)
	if err != nil {
		logger.Errorf("[stream] error reading video file: %s", err.Error())
		return
	}

	// start stream based on query param, if provided
	r.ParseForm()
	startTime := r.Form.Get("start")

	encoder := ffmpeg.NewEncoder(manager.GetInstance().FFMPEGPath)

	var stream io.ReadCloser
	var process *os.Process
	mimeType := ffmpeg.MimeWebm

	if audioCodec == ffmpeg.MissingUnsupported {
		//ffmpeg fails if it trys to transcode a non supported audio codec
		stream, process, err = encoder.StreamTranscodeVideo(*videoFile, startTime, config.GetMaxStreamingTranscodeSize())
	} else {
		copyVideo := false // try to be smart if the video to be transcoded is in a Matroska container
		//  mp4 has always supported audio so it doesn't need to be checked
		//  while mpeg_ts has seeking issues if we don't reencode the video

		if config.GetForceMKV() { // If MKV is forced as supported and video codec is also supported then only transcode audio
			if ffmpeg.Container(container) == ffmpeg.Matroska {
				switch videoCodec {
				case ffmpeg.H264, ffmpeg.Vp9, ffmpeg.Vp8:
					copyVideo = true
				case ffmpeg.Hevc:
					if config.GetForceHEVC() {
						copyVideo = true
					}

				}
			}
		}

		if copyVideo { // copy video stream instead of transcoding it
			stream, process, err = encoder.StreamMkvTranscodeAudio(*videoFile, startTime, config.GetMaxStreamingTranscodeSize())
			mimeType = ffmpeg.MimeMkv

		} else {
			stream, process, err = encoder.StreamTranscode(*videoFile, startTime, config.GetMaxStreamingTranscodeSize())
		}
	}

	if err != nil {
		logger.Errorf("[stream] error transcoding video file: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", mimeType)

	logger.Infof("[stream] transcoding video file to %s", mimeType)

	// handle if client closes the connection
	notify := r.Context().Done()
	go func() {
		<-notify
		logger.Info("[stream] client closed the connection. Killing stream process.")
		process.Kill()
	}()

	_, err = io.Copy(w, stream)
	if err != nil {
		logger.Errorf("[stream] error serving transcoded video file: %s", err.Error())
	}
}

func (rs sceneRoutes) Screenshot(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	filepath := manager.GetInstance().Paths.Scene.GetScreenshotPath(scene.Checksum)

	// fall back to the scene image blob if the file isn't present
	screenshotExists, _ := utils.FileExists(filepath)
	if screenshotExists {
		http.ServeFile(w, r, filepath)
	} else {
		qb := models.NewSceneQueryBuilder()
		cover, _ := qb.GetSceneCover(scene.ID, nil)
		utils.ServeImage(cover, w, r)
	}
}

func (rs sceneRoutes) Preview(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	filepath := manager.GetInstance().Paths.Scene.GetStreamPreviewPath(scene.Checksum)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) Webp(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	filepath := manager.GetInstance().Paths.Scene.GetStreamPreviewImagePath(scene.Checksum)
	http.ServeFile(w, r, filepath)
}

func getChapterVttTitle(marker *models.SceneMarker) string {
	if marker.Title != "" {
		return marker.Title
	}

	qb := models.NewTagQueryBuilder()
	primaryTag, err := qb.Find(marker.PrimaryTagID, nil)
	if err != nil {
		// should not happen
		panic(err)
	}

	ret := primaryTag.Name

	tags, err := qb.FindBySceneMarkerID(marker.ID, nil)
	if err != nil {
		// should not happen
		panic(err)
	}

	for _, t := range tags {
		ret += ", " + t.Name
	}

	return ret
}

func (rs sceneRoutes) ChapterVtt(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, err := qb.FindBySceneID(scene.ID, nil)
	if err != nil {
		panic("invalid scene markers for chapter vtt")
	}

	vttLines := []string{"WEBVTT", ""}
	for i, marker := range sceneMarkers {
		vttLines = append(vttLines, strconv.Itoa(i+1))
		time := utils.GetVTTTime(marker.Seconds)
		vttLines = append(vttLines, time+" --> "+time)
		vttLines = append(vttLines, getChapterVttTitle(marker))
		vttLines = append(vttLines, "")
	}
	vtt := strings.Join(vttLines, "\n")

	w.Header().Set("Content-Type", "text/vtt")
	_, _ = w.Write([]byte(vtt))
}

func (rs sceneRoutes) VttThumbs(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	w.Header().Set("Content-Type", "text/vtt")
	filepath := manager.GetInstance().Paths.Scene.GetSpriteVttFilePath(scene.Checksum)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) VttSprite(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	w.Header().Set("Content-Type", "image/jpeg")
	filepath := manager.GetInstance().Paths.Scene.GetSpriteImageFilePath(scene.Checksum)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) SceneMarkerStream(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneMarkerID, _ := strconv.Atoi(chi.URLParam(r, "sceneMarkerId"))
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarker, err := qb.Find(sceneMarkerID)
	if err != nil {
		logger.Warn("Error when getting scene marker for stream")
		http.Error(w, http.StatusText(404), 404)
		return
	}
	filepath := manager.GetInstance().Paths.SceneMarkers.GetStreamPath(scene.Checksum, int(sceneMarker.Seconds))
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) SceneMarkerPreview(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	sceneMarkerID, _ := strconv.Atoi(chi.URLParam(r, "sceneMarkerId"))
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarker, err := qb.Find(sceneMarkerID)
	if err != nil {
		logger.Warn("Error when getting scene marker for stream")
		http.Error(w, http.StatusText(404), 404)
		return
	}
	filepath := manager.GetInstance().Paths.SceneMarkers.GetStreamPreviewImagePath(scene.Checksum, int(sceneMarker.Seconds))

	// If the image doesn't exist, send the placeholder
	exists, _ := utils.FileExists(filepath)
	if !exists {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(utils.PendingGenerateResource)
		return
	}

	http.ServeFile(w, r, filepath)
}

// endregion

func SceneCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sceneIdentifierQueryParam := chi.URLParam(r, "sceneId")
		sceneID, _ := strconv.Atoi(sceneIdentifierQueryParam)

		var scene *models.Scene
		var err error
		qb := models.NewSceneQueryBuilder()
		if sceneID == 0 {
			scene, err = qb.FindByChecksum(sceneIdentifierQueryParam)
		} else {
			scene, err = qb.Find(sceneID)
		}

		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), sceneKey, scene)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
