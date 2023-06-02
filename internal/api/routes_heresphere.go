package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
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

type HeresphereAuthReq struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	NeedsMediaSource bool   `json:"needsMediaSource,omitempty"`
}
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
	// TODO: Floats become null, same with lists, should have default value instead
	// This is technically an api violation
}
type HeresphereVideoSubtitle struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Url      string `json:"url"`
}
type HeresphereVideoTag struct {
	// Should start with any of the following: "Scene:", "Category:", "Talent:", "Studio:", "Position:"
	Name string `json:"name"`
	/*Start  float64 `json:"start,omitempty"`
	End    float64 `json:"end,omitempty"`
	Track  int     `json:"track,omitempty"`
	Rating float32 `json:"rating,omitempty"`*/
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
	Favorites      float32                   `json:"favorites,omitempty"`
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
	scene          *models.Scene
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
	scene        *models.Scene
}
type HeresphereVideoEntryUpdate struct {
	Username   string               `json:"username"`
	Password   string               `json:"password"`
	IsFavorite bool                 `json:"isFavorite"`
	Rating     float32              `json:"rating,omitempty"`
	Tags       []HeresphereVideoTag `json:"tags"`
	// In base64
	Hsp        string `json:"hsp"`
	DeleteFile bool   `json:"deleteFile"`
}
type HeresphereVideoEvent struct {
	Username      string  `json:"username"`
	Id            string  `json:"id"`
	Title         string  `json:"title"`
	Event         int     `json:"event"`
	Time          float64 `json:"time,omitempty"`
	Speed         float32 `json:"speed,omitempty"`
	Utc           float64 `json:"utc,omitempty"`
	ConnectionKey string  `json:"connectionKey"`
}

type heresphereRoutes struct {
	txnManager  txn.Manager
	sceneFinder SceneFinder
	fileFinder  file.Finder
	repository  manager.Repository
	resolver    ResolverRoot
}

func (rs heresphereRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(rs.HeresphereCtx)
		r.Post("/", rs.HeresphereLogin)
		r.Get("/", rs.HeresphereIndex)
		r.Head("/", rs.HeresphereIndex)

		r.Post("/auth", rs.HeresphereLoginToken)
		r.Post("/scan", rs.HeresphereScan)
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

func (rs heresphereRoutes) HeresphereVideoEvent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	// TODO: Auth
}

func relUrlToAbs(r *http.Request, rel string) string {
	// Get the scheme (http or https) from the request
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Get the host from the request
	host := r.Host

	// TODO: Support Forwarded standard

	// Combine the scheme, host, and relative path to form the absolute URL
	return fmt.Sprintf("%s://%s%s", scheme, host, rel)
}

func (rs heresphereRoutes) getVideosTags(ctx context.Context, vids []HeresphereVideoEntryShort) {
	gallery_ids := make(map[int]*models.Gallery)
	tag_ids := make(map[int]*models.Tag)
	perf_ids := make(map[int]*models.Performer)
	//stash_ids

	for _, vid := range vids {
		vid.Tags = []HeresphereVideoTag{}

		if err := txn.WithReadTxn(ctx, rs.txnManager, func(ctx context.Context) error {
			return vid.scene.LoadRelationships(ctx, rs.repository.Scene)
		}); err != nil {
			continue
		}

		if vid.scene.GalleryIDs.Loaded() {
			for _, id := range vid.scene.GalleryIDs.List() {
				gallery_ids[id] = nil
			}
		}
		if vid.scene.TagIDs.Loaded() {
			for _, id := range vid.scene.TagIDs.List() {
				tag_ids[id] = nil
			}
		}
		if vid.scene.PerformerIDs.Loaded() {
			for _, id := range vid.scene.PerformerIDs.List() {
				perf_ids[id] = nil
			}
		}

		mark_ids, err := rs.resolver.Scene().SceneMarkers(ctx, vid.scene)
		if err == nil {
			for _, mark := range mark_ids {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Marker:%v", mark.Title),
				}
				vid.Tags = append(vid.Tags, genTag)
			}
		}

		movie_ids, err := rs.resolver.Scene().Movies(ctx, vid.scene)
		if err == nil {
			for _, movie := range movie_ids {
				if movie.Movie != nil {
					genTag := HeresphereVideoTag{
						Name: fmt.Sprintf("Movie:%v", movie.Movie.Name),
					}
					vid.Tags = append(vid.Tags, genTag)
				}
			}
		}

		studio_id, err := rs.resolver.Scene().Studio(ctx, vid.scene)
		if err == nil && studio_id != nil {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Studio:%v", studio_id.Name.String),
			}
			vid.Tags = append(vid.Tags, genTag)
		}
	}

	/*
	 * TODO: DAMN this is ugly
	 * When you make the code uglier for the sake of performance
	 * I cri
	 * 😭
	 */
	_gallery_ids := []int{}
	_tag_ids := []int{}
	_perf_ids := []int{}
	for id, _ := range gallery_ids {
		_gallery_ids = append(_gallery_ids, id)
	}
	for id, _ := range tag_ids {
		_tag_ids = append(_tag_ids, id)
	}
	for id, _ := range perf_ids {
		_perf_ids = append(_perf_ids, id)
	}

	r_gallery_ids, errs := loaders.From(ctx).GalleryByID.LoadAll(_gallery_ids)
	if firstError(errs) != nil {
		return
	}
	r_tag_ids, errs := loaders.From(ctx).TagByID.LoadAll(_tag_ids)
	if firstError(errs) != nil {
		return
	}
	r_perf_ids, errs := loaders.From(ctx).PerformerByID.LoadAll(_perf_ids)
	if firstError(errs) != nil {
		return
	}

	for idx, obj := range r_gallery_ids {
		gallery_ids[idx] = obj
	}
	for idx, obj := range r_tag_ids {
		tag_ids[idx] = obj
	}
	for idx, obj := range r_perf_ids {
		perf_ids[idx] = obj
	}

	for _, vid := range vids {
		for _, id := range vid.scene.GalleryIDs.List() {
			if gallery_ids[id] != nil {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Gallery:%v", gallery_ids[id].GetTitle()),
				}
				vid.Tags = append(vid.Tags, genTag)
			}
		}
		for _, id := range vid.scene.TagIDs.List() {
			if tag_ids[id] != nil {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Tag:%v", tag_ids[id].Name),
				}
				vid.Tags = append(vid.Tags, genTag)
			}
		}
		for _, id := range vid.scene.GalleryIDs.List() {
			if perf_ids[id] != nil {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Talent:%v", perf_ids[id].Name),
				}
				vid.Tags = append(vid.Tags, genTag)
			}
		}
	}
}

// TODO: Consolidate into one
func (rs heresphereRoutes) getVideoTags(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoTag {
	processedTags := []HeresphereVideoTag{}

	if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.repository.Scene)
	}); err != nil {
		return processedTags
	}

	mark_ids, err := rs.resolver.Scene().SceneMarkers(ctx, scene)
	if err == nil {
		for _, mark := range mark_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Marker:%v", mark.Title),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	gallery_ids, err := rs.resolver.Scene().Galleries(ctx, scene)
	if err == nil {
		for _, gal := range gallery_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Gallery:%v", gal.GetTitle()),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	tag_ids, err := rs.resolver.Scene().Tags(ctx, scene)
	if err == nil {
		for _, tag := range tag_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Tag:%v", tag.Name),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	perf_ids, err := rs.resolver.Scene().Performers(ctx, scene)
	if err == nil {
		for _, perf := range perf_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Talent:%s", perf.Name),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	movie_ids, err := rs.resolver.Scene().Movies(ctx, scene)
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

	//stash_ids, err := rs.resolver.Scene().StashIds(ctx, scene)

	studio_id, err := rs.resolver.Scene().Studio(ctx, scene)
	if err == nil && studio_id != nil {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Studio:%v", studio_id.Name.String),
		}
		processedTags = append(processedTags, genTag)
	}

	return processedTags
}

func (rs heresphereRoutes) getVideoScripts(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoScript {
	// TODO: Check if exists
	processedScripts := []HeresphereVideoScript{}

	exists, err := rs.resolver.Scene().Interactive(ctx, scene)
	if err == nil && exists {
		processedScript := HeresphereVideoScript{
			Name:   "Default script",
			Url:    relUrlToAbs(r, fmt.Sprintf("/scene/%v/funscript", scene.ID)),
			Rating: 4.2,
		}
		processedScripts = append(processedScripts, processedScript)
	}
	return processedScripts
}
func (rs heresphereRoutes) getVideoSubtitles(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoSubtitle {
	processedSubtitles := []HeresphereVideoSubtitle{}

	// TODO: Use sceneResolver Paths => rs.resolver.Scene().Paths()

	captions_id, err := rs.resolver.Scene().Captions(ctx, scene)
	if err == nil {
		for _, caption := range captions_id {
			processedCaption := HeresphereVideoSubtitle{
				Name:     caption.Filename,
				Language: caption.LanguageCode,
				Url:      relUrlToAbs(r, fmt.Sprintf("/scene/%v/caption?lang=%v&type=%v", scene.ID, caption.LanguageCode, caption.CaptionType)),
			}
			processedSubtitles = append(processedSubtitles, processedCaption)
		}
	}

	return processedSubtitles
}
func (rs heresphereRoutes) getVideoMedia(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoMedia {
	processedMedia := []HeresphereVideoMedia{}

	mediaTypes := make(map[string][]HeresphereVideoMediaSource)

	file_ids, err := rs.resolver.Scene().Files(ctx, scene)
	if err == nil {
		for _, mediaFile := range file_ids {
			processedEntry := HeresphereVideoMediaSource{
				Resolution: mediaFile.Height,
				Height:     mediaFile.Height,
				Width:      mediaFile.Width,
				Size:       mediaFile.Size,
				Url:        relUrlToAbs(r, fmt.Sprintf("/scene/%v/stream", scene.ID)),
			}
			mediaTypes[mediaFile.Format] = append(mediaTypes[mediaFile.Format], processedEntry)
		}
	}

	for codec, sources := range mediaTypes {
		processedMedia = append(processedMedia, HeresphereVideoMedia{
			Name:    codec,
			Sources: sources,
		})
	}
	// TODO: Transcode etc. /scene/%v/stream.mp4?resolution=ORIGINAL

	return processedMedia
}

func (rs heresphereRoutes) HeresphereIndex(w http.ResponseWriter, r *http.Request) {
	banner := HeresphereBanner{
		Image: relUrlToAbs(r, "/apple-touch-icon.png"),
		Link:  relUrlToAbs(r, "/"),
	}

	var scenes []*models.Scene
	if err := rs.repository.WithTxn(r.Context(), func(ctx context.Context) error {
		var err error
		scenes, err = rs.repository.Scene.All(ctx)
		return err
	}); err != nil {
		http.Error(w, "Failed to fetch scenes!", http.StatusInternalServerError)
		return
	}

	sceneUrls := make([]string, len(scenes))
	for idx, scene := range scenes {
		sceneUrls[idx] = relUrlToAbs(r, fmt.Sprintf("/heresphere/%v", scene.ID))
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
func (rs heresphereRoutes) HeresphereScan(w http.ResponseWriter, r *http.Request) {
	// TODO: Auth
	var scenes []*models.Scene
	if err := rs.repository.WithTxn(r.Context(), func(ctx context.Context) error {
		var err error
		scenes, err = rs.repository.Scene.All(ctx)
		return err
	}); err != nil {
		http.Error(w, "Failed to fetch scenes!", http.StatusInternalServerError)
		return
	}

	// TODO: /scan is damn slow
	// Processing each scene and creating a new list
	processedScenes := make([]HeresphereVideoEntryShort, len(scenes))
	for idx, scene := range scenes {
		// Perform the necessary processing on each scene
		processedScene := HeresphereVideoEntryShort{
			Link:         relUrlToAbs(r, fmt.Sprintf("/heresphere/%v", scene.ID)),
			Title:        scene.GetTitle(),
			DateReleased: scene.CreatedAt.Format("2006-01-02"),
			DateAdded:    scene.CreatedAt.Format("2006-01-02"),
			Duration:     60000.0,
			Rating:       0.0,
			Favorites:    0,
			Comments:     scene.OCounter,
			IsFavorite:   false,
			//Tags:         rs.getVideoTags(r.Context(), r, scene),
			scene: scene,
		}
		if scene.Date != nil {
			processedScene.DateReleased = scene.Date.Format("2006-01-02")
		}
		if scene.Rating != nil {
			isFavorite := *scene.Rating > 85
			processedScene.Rating = float32(*scene.Rating) * 0.05 // 0-5
			processedScene.IsFavorite = isFavorite
		}
		file_ids, err := rs.resolver.Scene().Files(r.Context(), scene)
		if err == nil && len(file_ids) > 0 {
			processedScene.Duration = handleFloat64Value(file_ids[0].Duration * 1000.0)
		}
		processedScenes[idx] = processedScene
	}
	rs.getVideosTags(r.Context(), processedScenes)

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(processedScenes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rs heresphereRoutes) HeresphereVideoHsp(w http.ResponseWriter, r *http.Request) {
	// TODO: Auth
	w.WriteHeader(http.StatusNotImplemented)
}

func (rs heresphereRoutes) HeresphereVideoDataUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: This
	// TODO: Auth
	w.WriteHeader(http.StatusNotImplemented)
}

func (rs heresphereRoutes) HeresphereVideoData(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(heresphereKey).(*models.Scene)

	// This endpoint can receive 2 types of requests
	// One is a video request (HeresphereAuthReq)
	// Other is an update (HeresphereVideoEntryUpdate)

	user := HeresphereAuthReq{
		NeedsMediaSource: true,
	}
	userupd := HeresphereVideoEntryUpdate{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		err = json.NewDecoder(r.Body).Decode(&userupd)
		if err == nil {
			rs.HeresphereVideoDataUpdate(w, r)
			return
		}

		/*http.Error(w, err.Error(), http.StatusBadRequest)
		return*/
	}

	// TODO: Auth

	processedScene := HeresphereVideoEntry{
		Access:         HeresphereMember,
		Title:          scene.GetTitle(),
		Description:    scene.Details,
		ThumbnailImage: relUrlToAbs(r, fmt.Sprintf("/scene/%v/screenshot", scene.ID)),
		ThumbnailVideo: relUrlToAbs(r, fmt.Sprintf("/scene/%v/preview", scene.ID)),
		DateReleased:   scene.CreatedAt.Format("2006-01-02"),
		DateAdded:      scene.CreatedAt.Format("2006-01-02"),
		Duration:       60000.0,
		Rating:         0.0,
		Favorites:      0,
		Comments:       scene.OCounter,
		IsFavorite:     false,
		Projection:     HeresphereProjectionPerspective, // Default to flat cause i have no idea
		Stereo:         HeresphereStereoMono,            // Default to flat cause i have no idea
		IsEyeSwapped:   false,
		Fov:            180,
		Lens:           HeresphereLensLinear,
		CameraIPD:      6.5,
		/*Hsp:            relUrlToAbs(r, fmt.Sprintf("/heresphere/%v/hsp", scene.ID)),
		EventServer:    relUrlToAbs(r, fmt.Sprintf("/heresphere/%v/event", scene.ID)),*/
		Scripts:       rs.getVideoScripts(r.Context(), r, scene),
		Subtitles:     rs.getVideoSubtitles(r.Context(), r, scene),
		Tags:          rs.getVideoTags(r.Context(), r, scene),
		Media:         []HeresphereVideoMedia{},
		WriteFavorite: false,
		WriteRating:   false,
		WriteTags:     false,
		WriteHSP:      false,
		scene:         scene,
	}
	//rs.getVideosTags(r.Context(), []HeresphereVideoEntry{processedScene})

	if user.NeedsMediaSource {
		processedScene.Media = rs.getVideoMedia(r.Context(), r, scene)
	}
	if scene.Date != nil {
		processedScene.DateReleased = scene.Date.Format("2006-01-02")
	}
	if scene.Rating != nil {
		isFavorite := *scene.Rating > 85
		processedScene.Rating = float32(*scene.Rating) * 0.05 // 0-5
		processedScene.IsFavorite = isFavorite
	}
	file_ids, err := rs.resolver.Scene().Files(r.Context(), scene)
	if err == nil && len(file_ids) > 0 {
		processedScene.Duration = handleFloat64Value(file_ids[0].Duration * 1000.0)
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(processedScene)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rs heresphereRoutes) HeresphereLogin(w http.ResponseWriter, r *http.Request) {
	var user HeresphereAuthReq
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if config.GetInstance().HasCredentials() {
		/*err := manager.GetInstance().SessionStore.Login(w, r)
		if err != nil {
			// always log the error
			logger.Errorf("Error logging in: %v", err)
		}*/

	}

	// TODO: Auth
	rs.HeresphereIndex(w, r)
}

func (rs heresphereRoutes) HeresphereLoginToken(w http.ResponseWriter, r *http.Request) {
	var user HeresphereAuthReq
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Auth

	//TODO: Will supply header auth-token in future requests, check it in other functions
	auth := &HeresphereAuthResp{
		AuthToken: "yes",
		Access:    HeresphereMember,
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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

func (rs heresphereRoutes) HeresphereCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header()["HereSphere-JSON-Version"] = []string{strconv.Itoa(HeresphereJsonVersion)}
		next.ServeHTTP(w, r)
	})
}
