package heresphere

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/txn"
)

type HeresphereCustomTag string

const (
	HeresphereCustomTagInteractive HeresphereCustomTag = "Interactive"

	HeresphereCustomTagPlayCount HeresphereCustomTag = "PlayCount"
	HeresphereCustomTagWatched   HeresphereCustomTag = "Watched"

	HeresphereCustomTagOrganized HeresphereCustomTag = "Organized"

	HeresphereCustomTagOCounter HeresphereCustomTag = "OCounter"
	HeresphereCustomTagOrgasmed HeresphereCustomTag = "Orgasmed"

	HeresphereCustomTagRated HeresphereCustomTag = "Rated"
)

type Routes struct {
	TxnManager  txn.Manager
	SceneFinder sceneFinder
	FileFinder  models.FileFinder
	Repository  manager.Repository
}

/*
 * This function provides the possible routes for this api.
 */
func (rs Routes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(rs.HeresphereCtx)

		r.Post("/", rs.HeresphereIndex)
		r.Get("/", rs.HeresphereIndex)
		r.Head("/", rs.HeresphereIndex)

		r.Post("/auth", rs.HeresphereLoginToken)
		r.Route("/{sceneId}", func(r chi.Router) {
			r.Use(rs.HeresphereSceneCtx)

			r.Post("/", rs.HeresphereVideoData)
			r.Get("/", rs.HeresphereVideoData)

			r.Post("/event", rs.HeresphereVideoEvent)
		})
	})

	return r
}

/*
 * This is a video playback event
 * Intended for server-sided script playback.
 * But since we dont need that, we just use it for timestamps.
 */
func (rs Routes) HeresphereVideoEvent(w http.ResponseWriter, r *http.Request) {
	scn := r.Context().Value(sceneKey).(*models.Scene)

	var event HeresphereVideoEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		logger.Errorf("Heresphere HeresphereVideoEvent decode error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if event.Event == HeresphereEventClose {
		newTime := event.Time / 1000
		newDuration := 0.0
		if newTime > scn.ResumeTime {
			newDuration += (newTime - scn.ResumeTime)
		}

		if err := updatePlayCount(r.Context(), scn, event, rs.TxnManager, rs.Repository.Scene); err != nil {
			logger.Errorf("Heresphere HeresphereVideoEvent updatePlayCount error: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := txn.WithReadTxn(r.Context(), rs.TxnManager, func(ctx context.Context) error {
			_, err := rs.Repository.Scene.SaveActivity(ctx, scn.ID, &newTime, &newDuration)
			return err
		}); err != nil {
			logger.Errorf("Heresphere HeresphereVideoEvent SaveActivity error: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

/*
 * This endpoint is for letting the user update scene data
 */
func (rs Routes) HeresphereVideoDataUpdate(w http.ResponseWriter, r *http.Request) error {
	scn := r.Context().Value(sceneKey).(*models.Scene)
	user := r.Context().Value(authKey).(HeresphereAuthReq)
	fileDeleter := file.NewDeleter()
	c := config.GetInstance()
	shouldUpdate := false

	ret := &scene.UpdateSet{
		ID:      scn.ID,
		Partial: models.NewScenePartial(),
	}

	var b bool
	var err error
	if user.Rating != nil && c.GetHSPWriteRatings() {
		if b, err = updateRating(user, ret); err != nil {
			fileDeleter.Rollback()
			return err
		}
		shouldUpdate = b || shouldUpdate
	}

	if user.DeleteFile != nil && *user.DeleteFile && c.GetHSPWriteDeletes() {
		if b, err = handleDeletePrimaryFile(r.Context(), rs.TxnManager, scn, rs.Repository.File, fileDeleter); err != nil {
			fileDeleter.Rollback()
			return err
		}
		shouldUpdate = b || shouldUpdate
	}

	if user.IsFavorite != nil && c.GetHSPWriteFavorites() {
		if b, err = handleFavoriteTag(r.Context(), rs, scn, &user, rs.TxnManager, ret); err != nil {
			return err
		}
		shouldUpdate = b || shouldUpdate
	}

	if user.Tags != nil && c.GetHSPWriteTags() {
		if b, err = handleTags(r.Context(), scn, &user, rs, ret); err != nil {
			return err
		}
		shouldUpdate = b || shouldUpdate
	}

	if shouldUpdate {
		if err := txn.WithTxn(r.Context(), rs.TxnManager, func(ctx context.Context) error {
			_, err := ret.Update(ctx, rs.Repository.Scene)
			return err
		}); err != nil {
			return err
		}

		fileDeleter.Commit()
		return nil
	}
	return nil
}

/*
 * This endpoint provides the main libraries that are available to browse.
 */
func (rs Routes) HeresphereIndex(w http.ResponseWriter, r *http.Request) {
	// Banner
	banner := HeresphereBanner{
		Image: fmt.Sprintf("%s%s", manager.GetBaseURL(r), "/apple-touch-icon.png"),
		Link:  fmt.Sprintf("%s%s", manager.GetBaseURL(r), "/"),
	}

	// Index
	libraryObj := HeresphereIndex{
		Access:  HeresphereMember,
		Banner:  banner,
		Library: []HeresphereIndexEntry{},
	}

	// Add filters
	parsedFilters, err := getAllFilters(r.Context(), rs.Repository)
	if err == nil {
		for key, value := range parsedFilters {
			sceneUrls := make([]string, len(value))

			for idx, sceneID := range value {
				sceneUrls[idx] = addApiKey(fmt.Sprintf("%s/heresphere/%d", manager.GetBaseURL(r), sceneID))
			}

			libraryObj.Library = append(libraryObj.Library, HeresphereIndexEntry{
				Name: key,
				List: sceneUrls,
			})
		}
	} else {
		logger.Warnf("Heresphere HeresphereIndex getAllFilters error: %s\n", err.Error())
	}

	// Set response headers and encode JSON
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(libraryObj); err != nil {
		logger.Errorf("Heresphere HeresphereIndex encode error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This endpoint provides a single scenes full information.
 */
func (rs Routes) HeresphereVideoData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authKey).(HeresphereAuthReq)
	c := config.GetInstance()

	// Update request
	if err := rs.HeresphereVideoDataUpdate(w, r); err != nil {
		logger.Errorf("Heresphere HeresphereVideoData HeresphereVideoDataUpdate error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch scene
	scene := r.Context().Value(sceneKey).(*models.Scene)

	// Load relationships
	processedScene := HeresphereVideoEntry{}
	if err := txn.WithReadTxn(r.Context(), rs.TxnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.Repository.Scene)
	}); err != nil {
		logger.Errorf("Heresphere HeresphereVideoData LoadRelationships error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create scene
	processedScene = HeresphereVideoEntry{
		Access:         HeresphereMember,
		Title:          scene.GetTitle(),
		Description:    scene.Details,
		ThumbnailImage: addApiKey(urlbuilders.NewSceneURLBuilder(manager.GetBaseURL(r), scene).GetScreenshotURL()),
		ThumbnailVideo: addApiKey(urlbuilders.NewSceneURLBuilder(manager.GetBaseURL(r), scene).GetStreamPreviewURL()),
		DateAdded:      scene.CreatedAt.Format("2006-01-02"),
		Duration:       60000.0,
		Rating:         0,
		Favorites:      0,
		Comments:       scene.OCounter,
		IsFavorite:     getVideoFavorite(rs, r, scene),
		Projection:     HeresphereProjectionPerspective,
		Stereo:         HeresphereStereoMono,
		IsEyeSwapped:   false,
		Fov:            180.0,
		Lens:           HeresphereLensLinear,
		CameraIPD:      6.5,
		EventServer: addApiKey(fmt.Sprintf("%s/heresphere/%d/event",
			manager.GetBaseURL(r),
			scene.ID,
		)),
		Scripts:       getVideoScripts(rs, r, scene),
		Subtitles:     getVideoSubtitles(rs, r, scene),
		Tags:          getVideoTags(r.Context(), rs, scene),
		Media:         []HeresphereVideoMedia{},
		WriteFavorite: c.GetHSPWriteFavorites(),
		WriteRating:   c.GetHSPWriteRatings(),
		WriteTags:     c.GetHSPWriteTags(),
		WriteHSP:      false,
	}

	// Find projection options
	FindProjectionTags(scene, &processedScene)

	// Additional info
	if user.NeedsMediaSource != nil && *user.NeedsMediaSource {
		processedScene.Media = getVideoMedia(rs, r, scene)
	}
	if scene.Date != nil {
		processedScene.DateReleased = scene.Date.Format("2006-01-02")
	}
	if scene.Rating != nil {
		fiveScale := models.Rating100To5F(*scene.Rating)
		processedScene.Rating = fiveScale
	}
	if processedScene.IsFavorite {
		processedScene.Favorites++
	}
	if scene.Files.PrimaryLoaded() {
		file_ids := scene.Files.Primary()
		if file_ids != nil {
			if val := manager.HandleFloat64(file_ids.Duration * 1000.0); val != nil {
				processedScene.Duration = *val
			}
		}
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(processedScene); err != nil {
		logger.Errorf("Heresphere HeresphereVideoData encode error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This endpoint function allows the user to login and receive a token if successful.
 */
func (rs Routes) HeresphereLoginToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authKey).(HeresphereAuthReq)

	// Try login
	if basicLogin(user.Username, user.Password) {
		writeNotAuthorized(w, r, "Invalid credentials")
		return
	}

	// Fetch key
	key := config.GetInstance().GetAPIKey()
	if len(key) == 0 {
		writeNotAuthorized(w, r, "Missing auth key!")
		return
	}

	// Generate auth response
	auth := &HeresphereAuthResp{
		AuthToken: key,
		Access:    HeresphereMember,
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(auth); err != nil {
		logger.Errorf("Heresphere HeresphereLoginToken encode error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This context function finds the applicable scene from the request and stores it.
 */
func (rs Routes) HeresphereSceneCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get sceneId
		sceneID, err := strconv.Atoi(chi.URLParam(r, "sceneId"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Resolve scene
		var scene *models.Scene
		_ = txn.WithReadTxn(r.Context(), rs.TxnManager, func(ctx context.Context) error {
			qb := rs.SceneFinder
			scene, _ = qb.Find(ctx, sceneID)

			if scene != nil {
				// A valid scene should have a attached video
				if err := scene.LoadPrimaryFile(ctx, rs.FileFinder); err != nil {
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

/*
 * This context function finds if the authentication is correct, otherwise rejects the request.
 */
func (rs Routes) HeresphereCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add JSON Header (using Add uses camel case and makes it invalid because "Json")
		w.Header()["HereSphere-JSON-Version"] = []string{strconv.Itoa(HeresphereJsonVersion)}

		// Only if enabled
		if !config.GetInstance().GetHSPDefaultEnabled() {
			writeNotAuthorized(w, r, "HereSphere API not enabled!")
			return
		}

		// Read HTTP Body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		// Make the Body re-readable (afaik only /event uses this)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Auth enabled and not has valid credentials (TLDR: needs to be blocked)
		isAuth := config.GetInstance().HasCredentials() && !HeresphereHasValidToken(r)

		// Default request
		user := HeresphereAuthReq{}

		// Attempt decode, and if err and invalid auth, fail
		if err := json.Unmarshal(body, &user); err != nil && isAuth {
			writeNotAuthorized(w, r, "Not logged in!")
			return
		}

		// If empty, fill as true
		if user.NeedsMediaSource == nil {
			needsMedia := true
			user.NeedsMediaSource = &needsMedia
		}

		// If invalid creds, only allow auth endpoint
		if isAuth && !strings.HasPrefix(r.URL.Path, "/heresphere/auth") {
			writeNotAuthorized(w, r, "Unauthorized!")
			return
		}

		ctx := context.WithValue(r.Context(), authKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
