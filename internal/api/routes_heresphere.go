package api

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
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/txn"
)

// Based on HereSphere_JSON_API_Version_1.txt

const HeresphereJsonVersion = 1

const (
	HeresphereGuest    = 0
	HeresphereMember   = 1
	HeresphereBadLogin = -1
)

type HeresphereProjection string

const (
	HeresphereProjectionEquirectangular        HeresphereProjection = "equirectangular"
	HeresphereProjectionPerspective            HeresphereProjection = "perspective"
	HeresphereProjectionEquirectangular360     HeresphereProjection = "equirectangular360"
	HeresphereProjectionFisheye                HeresphereProjection = "fisheye"
	HeresphereProjectionCubemap                HeresphereProjection = "cubemap"
	HeresphereProjectionEquirectangularCubemap HeresphereProjection = "equiangularCubemap"
)

type HeresphereStereo string

const (
	HeresphereStereoMono HeresphereStereo = "mono"
	HeresphereStereoSbs  HeresphereStereo = "sbs"
	HeresphereStereoTB   HeresphereStereo = "tb"
)

type HeresphereLens string

const (
	HeresphereLensLinear  HeresphereLens = "Linear"
	HeresphereLensMKX220  HeresphereLens = "MKX220"
	HeresphereLensMKX200  HeresphereLens = "MKX200"
	HeresphereLensVRCA220 HeresphereLens = "VRCA220"
)

const HeresphereAuthHeader = "auth-token"

type HeresphereAuthResp struct {
	AuthToken string `json:"auth-token"`
	Access    int    `json:"access"`
}

type HeresphereBanner struct {
	Image string `json:"image"`
	Link  string `json:"link"`
}
type HeresphereIndexEntry struct {
	Name string   `json:"name"`
	List []string `json:"list"`
}
type HeresphereIndex struct {
	Access  int                    `json:"access"`
	Banner  HeresphereBanner       `json:"banner"`
	Library []HeresphereIndexEntry `json:"library"`
}
type HeresphereVideoScript struct {
	Name   string  `json:"name"`
	Url    string  `json:"url"`
	Rating float32 `json:"rating,omitempty"`
}
type HeresphereVideoSubtitle struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Url      string `json:"url"`
}
type HeresphereVideoTag struct {
	Name   string  `json:"name"`
	Start  float64 `json:"start,omitempty"`
	End    float64 `json:"end,omitempty"`
	Track  int     `json:"track,omitempty"`
	Rating float32 `json:"rating,omitempty"`
}
type HeresphereVideoMediaSource struct {
	Resolution int `json:"resolution"`
	Height     int `json:"height"`
	Width      int `json:"width"`
	// In bytes
	Size int64  `json:"size"`
	Url  string `json:"url"`
}
type HeresphereVideoMedia struct {
	// Media type (h265 etc.)
	Name    string                       `json:"name"`
	Sources []HeresphereVideoMediaSource `json:"sources"`
}
type HeresphereVideoEntry struct {
	Access         int                       `json:"access"`
	Title          string                    `json:"title"`
	Description    string                    `json:"description"`
	ThumbnailImage string                    `json:"thumbnailImage"`
	ThumbnailVideo string                    `json:"thumbnailVideo,omitempty"`
	DateReleased   string                    `json:"dateReleased"`
	DateAdded      string                    `json:"dateAdded"`
	Duration       float64                   `json:"duration,omitempty"`
	Rating         float32                   `json:"rating,omitempty"`
	Favorites      int                       `json:"favorites"`
	Comments       int                       `json:"comments"`
	IsFavorite     bool                      `json:"isFavorite"`
	Projection     HeresphereProjection      `json:"projection"`
	Stereo         HeresphereStereo          `json:"stereo"`
	IsEyeSwapped   bool                      `json:"isEyeSwapped"`
	Fov            float32                   `json:"fov,omitempty"`
	Lens           HeresphereLens            `json:"lens"`
	CameraIPD      float32                   `json:"cameraIPD"`
	Hsp            string                    `json:"hsp,omitempty"`
	EventServer    string                    `json:"eventServer,omitempty"`
	Scripts        []HeresphereVideoScript   `json:"scripts,omitempty"`
	Subtitles      []HeresphereVideoSubtitle `json:"subtitles,omitempty"`
	Tags           []HeresphereVideoTag      `json:"tags,omitempty"`
	Media          []HeresphereVideoMedia    `json:"media,omitempty"`
	WriteFavorite  bool                      `json:"writeFavorite"`
	WriteRating    bool                      `json:"writeRating"`
	WriteTags      bool                      `json:"writeTags"`
	WriteHSP       bool                      `json:"writeHSP"`
}
type HeresphereVideoEntryShort struct {
	Link         string               `json:"link"`
	Title        string               `json:"title"`
	DateReleased string               `json:"dateReleased"`
	DateAdded    string               `json:"dateAdded"`
	Duration     float64              `json:"duration,omitempty"`
	Rating       float32              `json:"rating,omitempty"`
	Favorites    int                  `json:"favorites"`
	Comments     int                  `json:"comments"`
	IsFavorite   bool                 `json:"isFavorite"`
	Tags         []HeresphereVideoTag `json:"tags"`
}
type HeresphereAuthReq struct {
	Username         string               `json:"username"`
	Password         string               `json:"password"`
	NeedsMediaSource bool                 `json:"needsMediaSource,omitempty"`
	IsFavorite       bool                 `json:"isFavorite,omitempty"`
	Rating           float32              `json:"rating,omitempty"`
	Tags             []HeresphereVideoTag `json:"tags,omitempty"`
	// In base64
	Hsp        string `json:"hsp,omitempty"`
	DeleteFile bool   `json:"deleteFile,omitempty"`
}
type HeresphereVideoEvent struct {
	Username      string  `json:"username"`
	Id            string  `json:"id"`
	Title         string  `json:"title"`
	Event         int     `json:"event"`
	Time          float64 `json:"time"`
	Speed         float32 `json:"speed"`
	Utc           float64 `json:"utc"`
	ConnectionKey string  `json:"connectionKey"`
}

type heresphereRoutes struct {
	txnManager  txn.Manager
	sceneFinder SceneFinder
	fileFinder  file.Finder
	repository  manager.Repository
	resolver    ResolverRoot
}

/*
 * This function provides the possible routes for this api.
 */
func (rs heresphereRoutes) Routes() chi.Router {
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

			r.Get("/hsp", rs.HeresphereVideoHsp)
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
func (rs heresphereRoutes) HeresphereVideoEvent(w http.ResponseWriter, r *http.Request) {
	event := HeresphereVideoEvent{}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	scene := r.Context().Value(heresphereKey).(*models.Scene)

	newTime := event.Time / 1000
	newDuration := scene.PlayDuration
	if newTime > scene.ResumeTime {
		newDuration += (newTime - scene.ResumeTime)
	}

	if _, err := rs.resolver.Mutation().SceneSaveActivity(r.Context(), strconv.Itoa(scene.ID), &newTime, &newDuration); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/*
 * HSP is a HereSphere config file
 * It stores the players local config such as projection or color settings etc.
 */
func (rs heresphereRoutes) HeresphereVideoHsp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

/*
 * This endpoint is for letting the user update scene data
 */
func (rs heresphereRoutes) HeresphereVideoDataUpdate(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(heresphereKey).(*models.Scene)
	user := r.Context().Value(heresphereUserKey).(HeresphereAuthReq)
	fileDeleter := file.NewDeleter()

	if err := txn.WithTxn(r.Context(), rs.repository.TxnManager, func(ctx context.Context) error {
		qb := rs.repository.Scene

		rating := models.Rating5To100F(user.Rating)
		scene.Rating = &rating

		if user.DeleteFile {
			qe := rs.repository.File
			if err := scene.LoadPrimaryFile(r.Context(), qe); err != nil {
				ff := scene.Files.Primary()
				if ff != nil {
					if err := file.Destroy(ctx, qe, ff, fileDeleter, true); err != nil {
						return fmt.Errorf("destroying file %s: %w", ff.Base().Path, err)
					}
				}
			}
		}

		err := qb.Update(ctx, scene)
		return err
	}); err != nil {
		fileDeleter.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileDeleter.Commit()
	w.WriteHeader(http.StatusOK)
}

/*
 * This auxillary function gathers various tags from the scene to feed the api.
 */
func (rs heresphereRoutes) getVideoTags(r *http.Request, scene *models.Scene) []HeresphereVideoTag {
	processedTags := []HeresphereVideoTag{}

	if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.repository.Scene)
	}); err != nil {
		return processedTags
	}

	mark_ids, err := rs.resolver.Scene().SceneMarkers(r.Context(), scene)
	if err == nil {
		for _, mark := range mark_ids {
			genTag := HeresphereVideoTag{
				Name:  fmt.Sprintf("Marker:%v", mark.Title),
				Start: mark.Seconds * 1000,
				End:   (mark.Seconds + 60) * 1000,
			}
			processedTags = append(processedTags, genTag)
		}
	}

	gallery_ids, err := rs.resolver.Scene().Galleries(r.Context(), scene)
	if err == nil {
		for _, gal := range gallery_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Gallery:%v", gal.GetTitle()),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	tag_ids, err := rs.resolver.Scene().Tags(r.Context(), scene)
	if err == nil {
		for _, tag := range tag_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Tag:%v", tag.Name),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	perf_ids, err := rs.resolver.Scene().Performers(r.Context(), scene)
	if err == nil {
		for _, perf := range perf_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Talent:%s", perf.Name),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	movie_ids, err := rs.resolver.Scene().Movies(r.Context(), scene)
	if err == nil {
		for _, movie := range movie_ids {
			if movie.Movie != nil {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Movie:%v", movie.Movie.Name),
				}
				processedTags = append(processedTags, genTag)
			}
		}
	}

	// stash_ids, err := rs.resolver.Scene().StashIds(r.Context(), scene)

	studio_id, err := rs.resolver.Scene().Studio(r.Context(), scene)
	if err == nil && studio_id != nil {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Studio:%v", studio_id.Name.String),
		}
		processedTags = append(processedTags, genTag)
	}

	return processedTags
}

/*
 * This auxillary function gathers a script if applicable
 */
func (rs heresphereRoutes) getVideoScripts(r *http.Request, scene *models.Scene) []HeresphereVideoScript {
	processedScripts := []HeresphereVideoScript{}

	if interactive, err := rs.resolver.Scene().Interactive(r.Context(), scene); err == nil && interactive {
		processedScript := HeresphereVideoScript{
			Name: "Default script",
			Url: fmt.Sprintf("%s?apikey=%v",
				urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetFunscriptURL(),
				config.GetInstance().GetAPIKey(),
			),
			Rating: 5,
		}
		processedScripts = append(processedScripts, processedScript)
	}

	return processedScripts
}

/*
 * This auxillary function gathers subtitles if applicable
 */
func (rs heresphereRoutes) getVideoSubtitles(r *http.Request, scene *models.Scene) []HeresphereVideoSubtitle {
	processedSubtitles := []HeresphereVideoSubtitle{}

	if captions_id, err := rs.resolver.Scene().Captions(r.Context(), scene); err == nil {
		for _, caption := range captions_id {
			processedCaption := HeresphereVideoSubtitle{
				Name:     caption.Filename,
				Language: caption.LanguageCode,
				Url: fmt.Sprintf("%s?lang=%v&type=%v&apikey=%v",
					urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetCaptionURL(),
					caption.LanguageCode,
					caption.CaptionType,
					config.GetInstance().GetAPIKey(),
				),
			}
			processedSubtitles = append(processedSubtitles, processedCaption)
		}
	}

	return processedSubtitles
}

/*
 * This auxillary function gathers media information + transcoding options.
 */
func (rs heresphereRoutes) getVideoMedia(r *http.Request, scene *models.Scene) []HeresphereVideoMedia {
	processedMedia := []HeresphereVideoMedia{}

	mediaTypes := make(map[string][]HeresphereVideoMediaSource)

	if file_ids, err := rs.resolver.Scene().Files(r.Context(), scene); err == nil {
		for _, mediaFile := range file_ids {
			if mediaFile.ID == scene.PrimaryFileID.String() {
				sourceUrl := urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetStreamURL("").String()
				processedEntry := HeresphereVideoMediaSource{
					Resolution: mediaFile.Height,
					Height:     mediaFile.Height,
					Width:      mediaFile.Width,
					Size:       mediaFile.Size,
					Url:        fmt.Sprintf("%s?apikey=%s", sourceUrl, config.GetInstance().GetAPIKey()),
				}
				processedMedia = append(processedMedia, HeresphereVideoMedia{
					Name:    "direct stream",
					Sources: []HeresphereVideoMediaSource{processedEntry},
				})

				resRatio := mediaFile.Width / mediaFile.Height
				transcodeSize := config.GetInstance().GetMaxStreamingTranscodeSize()
				transNames := []string{"HLS", "DASH"}
				for i, trans := range []string{".m3u8", ".mpd"} {
					for _, res := range models.AllStreamingResolutionEnum {
						maxTrans := transcodeSize.GetMaxResolution()
						if height := res.GetMaxResolution(); (maxTrans == 0 || maxTrans >= height) && height <= mediaFile.Height {
							processedEntry.Resolution = height
							processedEntry.Height = height
							processedEntry.Width = resRatio * height
							processedEntry.Size = 0
							if height == 0 {
								processedEntry.Resolution = mediaFile.Height
								processedEntry.Height = mediaFile.Height
								processedEntry.Width = mediaFile.Width
								processedEntry.Size = mediaFile.Size
							}
							processedEntry.Url = fmt.Sprintf("%s%s?resolution=%s&apikey=%s", sourceUrl, trans, res.String(), config.GetInstance().GetAPIKey())

							typeName := transNames[i]
							mediaTypes[typeName] = append(mediaTypes[typeName], processedEntry)
						}
					}
				}
			}
		}
	}

	for codec, sources := range mediaTypes {
		processedMedia = append(processedMedia, HeresphereVideoMedia{
			Name:    codec,
			Sources: sources,
		})
	}

	return processedMedia
}

/*
 * This endpoint provides the main libraries that are available to browse.
 */
func (rs heresphereRoutes) HeresphereIndex(w http.ResponseWriter, r *http.Request) {
	banner := HeresphereBanner{
		Image: fmt.Sprintf("%s%s",
			GetBaseURL(r),
			"/apple-touch-icon.png",
		),
		Link: fmt.Sprintf("%s%s",
			GetBaseURL(r),
			"/",
		),
	}

	var scenes []*models.Scene
	if err := txn.WithReadTxn(r.Context(), rs.repository.TxnManager, func(ctx context.Context) error {
		var err error
		scenes, err = rs.repository.Scene.All(ctx)
		return err
	}); err != nil {
		http.Error(w, "Failed to fetch scenes!", http.StatusInternalServerError)
		return
	}

	sceneUrls := make([]string, len(scenes))
	for idx, scene := range scenes {
		sceneUrls[idx] = fmt.Sprintf("%s/heresphere/%v",
			GetBaseURL(r),
			scene.ID,
		)
	}

	library := HeresphereIndexEntry{
		Name: "All",
		List: sceneUrls,
	}
	idx := HeresphereIndex{
		Access:  HeresphereMember,
		Banner:  banner,
		Library: []HeresphereIndexEntry{library},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(idx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This auxillary function finds vr projection modes from tags and the filename.
 */
func FindProjectionTags(scene *models.Scene, processedScene *HeresphereVideoEntry) {
	// Detect VR modes from tags
	for _, tag := range processedScene.Tags {
		if strings.Contains(tag.Name, "°") {
			deg := strings.ReplaceAll(tag.Name, "°", "")
			if s, err := strconv.ParseFloat(deg, 32); err == nil {
				processedScene.Fov = float32(s)
			}
		}
		if strings.Contains(tag.Name, "Virtual Reality") || strings.Contains(tag.Name, "JAVR") {
			if processedScene.Projection == HeresphereProjectionPerspective {
				processedScene.Projection = HeresphereProjectionEquirectangular
			}
			if processedScene.Stereo == HeresphereStereoMono {
				processedScene.Stereo = HeresphereStereoSbs
			}
		}
		if strings.Contains(tag.Name, "Fisheye") {
			processedScene.Projection = HeresphereProjectionFisheye
			if processedScene.Stereo == HeresphereStereoMono {
				processedScene.Stereo = HeresphereStereoSbs
			}
		}
	}

	// Detect VR modes from filename
	file := scene.Files.Primary()
	if file != nil {
		path := strings.ToUpper(file.Basename)

		if strings.Contains(path, "_LR") || strings.Contains(path, "_3DH") {
			processedScene.Stereo = HeresphereStereoSbs
		}
		if strings.Contains(path, "_RL") {
			processedScene.Stereo = HeresphereStereoSbs
			processedScene.IsEyeSwapped = true
		}
		if strings.Contains(path, "_TB") || strings.Contains(path, "_3DV") {
			processedScene.Stereo = HeresphereStereoTB
		}
		if strings.Contains(path, "_BT") {
			processedScene.Stereo = HeresphereStereoTB
			processedScene.IsEyeSwapped = true
		}

		if strings.Contains(path, "_EAC360") || strings.Contains(path, "_360EAC") {
			processedScene.Projection = HeresphereProjectionEquirectangularCubemap
		}
		if strings.Contains(path, "_360") {
			processedScene.Projection = HeresphereProjectionEquirectangular360
		}
		if strings.Contains(path, "_F180") || strings.Contains(path, "_180F") || strings.Contains(path, "_VR180") {
			processedScene.Projection = HeresphereProjectionFisheye
		} else if strings.Contains(path, "_180") {
			processedScene.Projection = HeresphereProjectionEquirectangular
		}
		if strings.Contains(path, "_MKX200") {
			processedScene.Projection = HeresphereProjectionFisheye
			processedScene.Fov = 200.0
			processedScene.Lens = HeresphereLensMKX200
		}
		if strings.Contains(path, "_MKX220") {
			processedScene.Projection = HeresphereProjectionFisheye
			processedScene.Fov = 220.0
			processedScene.Lens = HeresphereLensMKX220
		}
		if strings.Contains(path, "_RF52") {
			processedScene.Projection = HeresphereProjectionFisheye
			processedScene.Fov = 190.0
		}
		if strings.Contains(path, "_VRCA220") {
			processedScene.Projection = HeresphereProjectionFisheye
			processedScene.Fov = 220.0
			processedScene.Lens = HeresphereLensVRCA220
		}
	}
}

/*
 * This endpoint provides a single scenes full information.
 */
func (rs heresphereRoutes) HeresphereVideoData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(heresphereUserKey).(HeresphereAuthReq)
	if user.Tags != nil {
		rs.HeresphereVideoDataUpdate(w, r)
		return
	}

	scene := r.Context().Value(heresphereKey).(*models.Scene)

	processedScene := HeresphereVideoEntry{}
	if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.repository.Scene)
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	processedScene = HeresphereVideoEntry{
		Access:      HeresphereMember,
		Title:       scene.GetTitle(),
		Description: scene.Details,
		ThumbnailImage: fmt.Sprintf("%s?apikey=%v",
			urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetScreenshotURL(),
			config.GetInstance().GetAPIKey(),
		),
		ThumbnailVideo: fmt.Sprintf("%s?apikey=%v",
			urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetStreamPreviewURL(),
			config.GetInstance().GetAPIKey(),
		),
		DateAdded:    scene.CreatedAt.Format("2006-01-02"),
		Duration:     60000.0,
		Rating:       0,
		Favorites:    scene.OCounter,
		Comments:     0,
		IsFavorite:   false,
		Projection:   HeresphereProjectionPerspective,
		Stereo:       HeresphereStereoMono,
		IsEyeSwapped: false,
		Fov:          180.0,
		Lens:         HeresphereLensLinear,
		CameraIPD:    6.5,
		Hsp: fmt.Sprintf("%s/heresphere/%v/hsp?apikey=%v",
			GetBaseURL(r),
			scene.ID,
			config.GetInstance().GetAPIKey(),
		),
		EventServer: fmt.Sprintf("%s/heresphere/%v/event?apikey=%v",
			GetBaseURL(r),
			scene.ID,
			config.GetInstance().GetAPIKey(),
		),
		Scripts:       rs.getVideoScripts(r, scene),
		Subtitles:     rs.getVideoSubtitles(r, scene),
		Tags:          rs.getVideoTags(r, scene),
		Media:         []HeresphereVideoMedia{},
		WriteFavorite: false,
		WriteRating:   true,
		WriteTags:     false,
		WriteHSP:      false,
	}
	FindProjectionTags(scene, &processedScene)

	if user.NeedsMediaSource {
		processedScene.Media = rs.getVideoMedia(r, scene)
	}
	if scene.Date != nil {
		processedScene.DateReleased = scene.Date.Format("2006-01-02")
	}
	if scene.Rating != nil {
		fiveScale := models.Rating100To5F(*scene.Rating)
		processedScene.Rating = fiveScale
		processedScene.IsFavorite = fiveScale >= 4
	}

	file_ids := scene.Files.Primary()
	if file_ids != nil {
		processedScene.Duration = handleFloat64Value(file_ids.Duration * 1000.0)
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(processedScene)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This auxillary function finds if a login is needed, and auth is correct.
 */
func basicLogin(username string, password string) bool {
	if config.GetInstance().HasCredentials() {
		err := manager.GetInstance().SessionStore.LoginPlain(username, password)
		return err != nil
	}
	return false
}

/*
 * This endpoint function allows the user to login and receive a token if successful.
 */
func (rs heresphereRoutes) HeresphereLoginToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(heresphereUserKey).(HeresphereAuthReq)

	if basicLogin(user.Username, user.Password) {
		writeNotAuthorized(w, r, "Invalid credentials")
		return
	}

	key := config.GetInstance().GetAPIKey()
	if len(key) == 0 {
		writeNotAuthorized(w, r, "Missing auth key!")
		return
	}

	auth := &HeresphereAuthResp{
		AuthToken: key,
		Access:    HeresphereMember,
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This auxillary function finds if the request has a valid auth token.
 */
func HeresphereHasValidToken(r *http.Request) bool {
	apiKey := r.Header.Get(HeresphereAuthHeader)

	if apiKey == "" {
		apiKey = r.URL.Query().Get(session.ApiKeyParameter)
	}

	return len(apiKey) > 0 && apiKey == config.GetInstance().GetAPIKey()
}

/*
 * This auxillary writes a library with a fake name upon auth failure
 */
func writeNotAuthorized(w http.ResponseWriter, r *http.Request, msg string) {
	banner := HeresphereBanner{
		Image: fmt.Sprintf("%s%s",
			GetBaseURL(r),
			"/apple-touch-icon.png",
		),
		Link: fmt.Sprintf("%s%s",
			GetBaseURL(r),
			"/",
		),
	}
	library := HeresphereIndexEntry{
		Name: msg,
		List: []string{fmt.Sprintf("%s/heresphere/doesnt-exist", GetBaseURL(r))},
	}
	idx := HeresphereIndex{
		Access:  HeresphereBadLogin,
		Banner:  banner,
		Library: []HeresphereIndexEntry{library},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(idx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This http handler redirects HereSphere if enabled
 */
func heresphereHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := config.GetInstance()

			if strings.Contains(r.UserAgent(), "HereSphere") && c.GetRedirectHeresphere() && (r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/login")) {
				http.Redirect(w, r, "/heresphere", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/*
 * This context function finds the applicable scene from the request and stores it.
 */
func (rs heresphereRoutes) HeresphereSceneCtx(next http.Handler) http.Handler {
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

		ctx := context.WithValue(r.Context(), heresphereKey, scene)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*
 * This context function finds if the authentication is correct, otherwise rejects the request.
 */
func (rs heresphereRoutes) HeresphereCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header()["HereSphere-JSON-Version"] = []string{strconv.Itoa(HeresphereJsonVersion)}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		isAuth := config.GetInstance().HasCredentials() && !HeresphereHasValidToken(r)
		user := HeresphereAuthReq{NeedsMediaSource: true, DeleteFile: false}
		if err := json.Unmarshal(body, &user); err != nil && isAuth {
			writeNotAuthorized(w, r, "Not logged in!")
			return
		}

		if isAuth && !strings.HasPrefix(r.URL.Path, "/heresphere/auth") {
			writeNotAuthorized(w, r, "Unauthorized!")
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		ctx := context.WithValue(r.Context(), heresphereUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
