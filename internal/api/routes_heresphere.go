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

func (rs heresphereRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(rs.HeresphereCtx)
		r.Post("/", rs.HeresphereIndex)
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

func (rs heresphereRoutes) HeresphereVideoHsp(w http.ResponseWriter, r *http.Request) {
	// TODO: This
	w.WriteHeader(http.StatusNotImplemented)
}

func (rs heresphereRoutes) HeresphereVideoDataUpdate(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(heresphereKey).(*models.Scene)
	user := r.Context().Value(heresphereUserKey).(HeresphereAuthReq)

	if err := txn.WithTxn(r.Context(), rs.repository.TxnManager, func(ctx context.Context) error {
		qb := rs.repository.Scene

		rating := int((user.Rating / 5.0) * 100.0)
		scene.Rating = &rating
		// TODO: user.Hsp
		// TODO: user.DeleteFile

		err := qb.Update(ctx, scene)
		return err
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
				Name:  fmt.Sprintf("Marker:%v", mark.Title),
				Start: mark.Seconds * 1000,
				End:   (mark.Seconds + 60) * 1000,
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

	// stash_ids, err := rs.resolver.Scene().StashIds(ctx, scene)

	studio_id, err := rs.resolver.Scene().Studio(ctx, scene)
	if err == nil && studio_id != nil {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Studio:%v", studio_id.Name.String),
		}
		processedTags = append(processedTags, genTag)
	}

	exists, err := rs.resolver.Scene().Interactive(ctx, scene)
	if err == nil && exists {
		shouldAdd := true
		for _, tag := range processedTags {
			if strings.Contains(tag.Name, "Tag:Interactive") {
				shouldAdd = false
				break
			}
		}

		if shouldAdd {
			genTag := HeresphereVideoTag{
				Name: "Tag:Interactive",
			}
			processedTags = append(processedTags, genTag)
		}
	}

	return processedTags
}

func (rs heresphereRoutes) getVideoScripts(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoScript {
	processedScripts := []HeresphereVideoScript{}

	// TODO: Use urlbuilders
	exists, err := rs.resolver.Scene().Interactive(ctx, scene)
	if err == nil && exists {
		processedScript := HeresphereVideoScript{
			Name:   "Default script",
			Url:    relUrlToAbs(r, fmt.Sprintf("/scene/%v/funscript?apikey=%v", scene.ID, config.GetInstance().GetAPIKey())),
			Rating: 4.2,
		}
		processedScripts = append(processedScripts, processedScript)
	}

	return processedScripts
}
func (rs heresphereRoutes) getVideoSubtitles(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoSubtitle {
	processedSubtitles := []HeresphereVideoSubtitle{}

	captions_id, err := rs.resolver.Scene().Captions(ctx, scene)
	if err == nil {
		for _, caption := range captions_id {
			processedCaption := HeresphereVideoSubtitle{
				Name:     caption.Filename,
				Language: caption.LanguageCode,
				Url:      relUrlToAbs(r, fmt.Sprintf("/scene/%v/caption?lang=%v&type=%v&apikey=%v", scene.ID, caption.LanguageCode, caption.CaptionType, config.GetInstance().GetAPIKey())),
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
				Url:        relUrlToAbs(r, fmt.Sprintf("/scene/%v/stream?apikey=%v", scene.ID, config.GetInstance().GetAPIKey())),
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
	// TODO: A very helpful feature would be to create multiple libraries from the saved filters: Every saved filter scene should become a separate library, the "All"-Library should remain as well.

	filters, err := rs.repository.SavedFilter.All(r.Context())
	if err == nil {
		for _, filter := range filters {
			fmt.Printf("Filter: %v -> %v\n", filter.Name, filter.Filter)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(idx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TODO: Consider removing, manual scan is faster
func (rs heresphereRoutes) HeresphereScan(w http.ResponseWriter, r *http.Request) {
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
			Link:  relUrlToAbs(r, fmt.Sprintf("/heresphere/%v", scene.ID)),
			Title: scene.GetTitle(),
			// DateReleased: scene.CreatedAt.Format("2006-01-02"),
			DateAdded:  scene.CreatedAt.Format("2006-01-02"),
			Duration:   60000.0,
			Rating:     0.0,
			Favorites:  scene.OCounter,
			Comments:   0,
			IsFavorite: false,
			Tags:       rs.getVideoTags(r.Context(), r, scene),
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
		// fmt.Printf("Done scene: %v\n", idx)
		processedScenes[idx] = processedScene
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(processedScenes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Check against filename for vr modes
func FindProjection(path string) (proj HeresphereProjection, stereo HeresphereStereo, eyeswap bool, fov float32, lens HeresphereLens, ipd float32) {
	proj, stereo, eyeswap, fov, lens, ipd = HeresphereProjectionPerspective, HeresphereStereoMono, false, 180, HeresphereLensLinear, 6.5

	path = strings.ToUpper(path)
	// TODO: Detect video tags

	if strings.Contains(path, "_LR") || strings.Contains(path, "_3DH") {
		stereo = HeresphereStereoSbs
	}
	if strings.Contains(path, "_RL") {
		stereo = HeresphereStereoSbs
		eyeswap = true
	}
	if strings.Contains(path, "_TB") || strings.Contains(path, "_3DV") {
		stereo = HeresphereStereoTB
	}
	if strings.Contains(path, "_BT") {
		stereo = HeresphereStereoTB
		eyeswap = true
	}

	if strings.Contains(path, "_EAC360") || strings.Contains(path, "_360EAC") {
		proj = HeresphereProjectionEquirectangularCubemap
	}
	if strings.Contains(path, "_360") {
		proj = HeresphereProjectionEquirectangular360
	}
	if strings.Contains(path, "_F180") || strings.Contains(path, "_180F") || strings.Contains(path, "_VR180") {
		proj = HeresphereProjectionFisheye
	} else if strings.Contains(path, "_180") {
		proj = HeresphereProjectionEquirectangular
	}
	if strings.Contains(path, "_MKX200") {
		proj = HeresphereProjectionFisheye
		fov = 200
		lens = HeresphereLensMKX200
	}
	if strings.Contains(path, "_MKX220") {
		proj = HeresphereProjectionFisheye
		fov = 220
		lens = HeresphereLensMKX220
	}
	if strings.Contains(path, "_RF52") {
		proj = HeresphereProjectionFisheye
		fov = 190
	}
	if strings.Contains(path, "_VRCA220") {
		proj = HeresphereProjectionFisheye
		fov = 220
		lens = HeresphereLensVRCA220
	}

	return
}

// Check against stashdb tags for vr modes
func FindProjectionTags(scene *HeresphereVideoEntry) {
	for _, tag := range scene.Tags {
		if strings.Contains(tag.Name, "°") {
			deg := strings.ReplaceAll(tag.Name, "°", "")
			if s, err := strconv.ParseFloat(deg, 32); err == nil {
				scene.Fov = float32(s)
			}
		}
		if strings.Contains(tag.Name, "Virtual Reality") || strings.Contains(tag.Name, "JAVR") {
			if scene.Projection == HeresphereProjectionPerspective {
				scene.Projection = HeresphereProjectionEquirectangular
			}
			if scene.Stereo == HeresphereStereoMono {
				scene.Stereo = HeresphereStereoSbs
			}
		}
		if strings.Contains(tag.Name, "Fisheye") {
			scene.Projection = HeresphereProjectionFisheye
			if scene.Stereo == HeresphereStereoMono {
				scene.Stereo = HeresphereStereoSbs
			}
		}
	}
}

func (rs heresphereRoutes) HeresphereVideoData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(heresphereUserKey).(HeresphereAuthReq)
	if user.Tags != nil {
		rs.HeresphereVideoDataUpdate(w, r)
		return
	}

	scene := r.Context().Value(heresphereKey).(*models.Scene)
	proj, stereo, isEyeSwapped, fov, lens, ipd := FindProjection(scene.Files.Primary().Basename)
	processedScene := HeresphereVideoEntry{
		Access:         HeresphereMember,
		Title:          scene.GetTitle(),
		Description:    scene.Details,
		ThumbnailImage: relUrlToAbs(r, fmt.Sprintf("/scene/%v/screenshot?apikey=%v", scene.ID, config.GetInstance().GetAPIKey())),
		ThumbnailVideo: relUrlToAbs(r, fmt.Sprintf("/scene/%v/preview?apikey=%v", scene.ID, config.GetInstance().GetAPIKey())),
		// DateReleased:   scene.CreatedAt.Format("2006-01-02"),
		DateAdded:    scene.CreatedAt.Format("2006-01-02"),
		Duration:     60000.0,
		Rating:       0.0,
		Favorites:    scene.OCounter,
		Comments:     0,
		IsFavorite:   false,
		Projection:   proj,
		Stereo:       stereo,
		IsEyeSwapped: isEyeSwapped,
		Fov:          fov,
		Lens:         lens,
		CameraIPD:    ipd,
		// Hsp:            relUrlToAbs(r, fmt.Sprintf("/heresphere/%v/hsp?apikey=%v", scene.ID, config.GetInstance().GetAPIKey())),
		EventServer:   relUrlToAbs(r, fmt.Sprintf("/heresphere/%v/event?apikey=%v", scene.ID, config.GetInstance().GetAPIKey())),
		Scripts:       rs.getVideoScripts(r.Context(), r, scene),
		Subtitles:     rs.getVideoSubtitles(r.Context(), r, scene),
		Tags:          rs.getVideoTags(r.Context(), r, scene),
		Media:         []HeresphereVideoMedia{},
		WriteFavorite: false,
		WriteRating:   false,
		WriteTags:     false,
		WriteHSP:      false,
	}
	FindProjectionTags(&processedScene)

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

func basicLogin(username string, password string) bool {
	if config.GetInstance().HasCredentials() {
		err := manager.GetInstance().SessionStore.LoginPlain(username, password)
		return err != nil
	}
	return false
}

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

// TODO: This is a copy of the Ctx from routes_scene
// Create a general version
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

func HeresphereHasValidToken(r *http.Request) bool {
	apiKey := r.Header.Get("auth-token")

	if apiKey == "" {
		apiKey = r.URL.Query().Get(session.ApiKeyParameter)
	}

	return len(apiKey) > 0 && apiKey == config.GetInstance().GetAPIKey()
}

func writeNotAuthorized(w http.ResponseWriter, r *http.Request, msg string) {
	banner := HeresphereBanner{
		Image: relUrlToAbs(r, "/apple-touch-icon.png"),
		Link:  relUrlToAbs(r, "/"),
	}
	library := HeresphereIndexEntry{
		Name: msg,
		List: []string{relUrlToAbs(r, "/heresphere/doesnt-exist")},
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

func (rs heresphereRoutes) HeresphereCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header()["HereSphere-JSON-Version"] = []string{strconv.Itoa(HeresphereJsonVersion)}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		isAuth := config.GetInstance().HasCredentials() && !HeresphereHasValidToken(r)
		user := HeresphereAuthReq{NeedsMediaSource: true}
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
