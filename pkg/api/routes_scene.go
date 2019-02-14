package api

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"net/http"
	"strconv"
	"strings"
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
	filepath := manager.GetInstance().Paths.Scene.GetStreamPath(scene.Path, scene.Checksum)
	http.ServeFile(w, r, filepath)
}

func (rs sceneRoutes) Screenshot(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	filepath := manager.GetInstance().Paths.Scene.GetScreenshotPath(scene.Checksum)
	http.ServeFile(w, r, filepath)
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

func (rs sceneRoutes) ChapterVtt(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(sceneKey).(*models.Scene)
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, err := qb.FindBySceneID(scene.ID, nil)
	if err != nil {
		panic("invalid scene markers for chapter vtt")
	}

	vttLines := []string{"WEBVTT", ""}
	for _, marker := range sceneMarkers {
		time := utils.GetVTTTime(marker.Seconds)
		vttLines = append(vttLines, time+" --> "+time)
		vttLines = append(vttLines, marker.Title)
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
