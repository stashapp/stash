package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

// Based on HereSphere_JSON_API_Version_1.txt

const heresphereJsonVersion = 1

const (
	heresphereGuest    = 0
	heresphereMember   = 1
	heresphereBadLogin = -1
)

type heresphereProjection string

const (
	heresphereProjectionEquirectangular        heresphereProjection = "equirectangular"
	heresphereProjectionPerspective            heresphereProjection = "perspective"
	heresphereProjectionEquirectangular360     heresphereProjection = "equirectangular360"
	heresphereProjectionFisheye                heresphereProjection = "fisheye"
	heresphereProjectionCubemap                heresphereProjection = "cubemap"
	heresphereProjectionEquirectangularCubemap heresphereProjection = "equiangularCubemap"
)

type heresphereStereo string

const (
	heresphereStereoMono heresphereStereo = "mono"
	heresphereStereoSbs  heresphereStereo = "sbs"
	heresphereStereoTB   heresphereStereo = "tb"
)

type heresphereLens string

const (
	heresphereLensLinear  heresphereLens = "Linear"
	heresphereLensMKX220  heresphereLens = "MKX220"
	heresphereLensMKX200  heresphereLens = "MKX200"
	heresphereLensVRCA220 heresphereLens = "VRCA220"
)

type heresphereAuthReq struct {
	username string `json:"username"`
	password string `json:"password"`
}
type heresphereAuthResp struct {
	auth_token string `json:"auth-token"`
	access     int    `json:"access"`
}

type heresphereBanner struct {
	image string `json:"image"`
	link  string `json:"link"`
}
type heresphereIndexEntry struct {
	name string   `json:"name"`
	list []string `json:"list"`
}
type heresphereIndex struct {
	access  int                    `json:"access"`
	banner  heresphereBanner       `json:"banner"`
	library []heresphereIndexEntry `json:"library"`
}
type heresphereVideoReq struct {
	username         string `json:"username"`
	password         string `json:"password"`
	needsMediaSource bool   `json:"needsMediaSource,omitempty"`
}
type heresphereVideoScript struct {
	name   string  `json:"name"`
	url    string  `json:"url"`
	rating float32 `json:"rating"`
}
type heresphereVideoSubtitle struct {
	name     string `json:"name"`
	language string `json:"language"`
	url      string `json:"url"`
}
type heresphereVideoTag struct {
	// Should start with any of the following: "Scene:", "Category:", "Talent:", "Studio:", "Position:"
	name   string  `json:"name"`
	start  float64 `json:"start"`
	end    float64 `json:"end"`
	track  int     `json:"track"`
	rating float32 `json:"rating"`
}
type heresphereVideoMediaSource struct {
	resolution int `json:"resolution"`
	height     int `json:"height"`
	width      int `json:"width"`
	// In bytes
	size int64  `json:"size"`
	url  string `json:"url"`
}
type heresphereVideoMedia struct {
	// Media type (h265 etc.)
	name    string                       `json:"name"`
	sources []heresphereVideoMediaSource `json:"sources"`
}
type heresphereVideoEntry struct {
	access         int                       `json:"access"`
	title          string                    `json:"title"`
	description    string                    `json:"description"`
	thumbnailImage string                    `json:"thumbnailImage"`
	thumbnailVideo string                    `json:"thumbnailVideo"`
	dateReleased   string                    `json:"dateReleased"`
	dateAdded      string                    `json:"dateAdded"`
	duration       float64                   `json:"duration"`
	rating         float32                   `json:"rating"`
	favorites      float32                   `json:"favorites"`
	comments       int                       `json:"comments"`
	isFavorite     bool                      `json:"isFavorite"`
	projection     heresphereProjection      `json:"projection"`
	stereo         heresphereStereo          `json:"stereo"`
	isEyeSwapped   bool                      `json:"isEyeSwapped"`
	fov            float32                   `json:"fov"`
	lens           heresphereLens            `json:"lens"`
	cameraIPD      float32                   `json:"cameraIPD"`
	hsp            string                    `json:"hsp,omitempty"`
	eventServer    string                    `json:"eventServer,omitempty"`
	scripts        []heresphereVideoScript   `json:"scripts"`
	subtitles      []heresphereVideoSubtitle `json:"subtitles"`
	tags           []heresphereVideoTag      `json:"tags"`
	media          []heresphereVideoMedia    `json:"media"`
	writeFavorite  bool                      `json:"writeFavorite"`
	writeRating    bool                      `json:"writeRating"`
	writeTags      bool                      `json:"writeTags"`
	writeHSP       bool                      `json:"writeHSP"`
}
type heresphereVideoEntryShort struct {
	link         string               `json:"link"`
	title        string               `json:"title"`
	dateReleased string               `json:"dateReleased"`
	dateAdded    string               `json:"dateAdded"`
	duration     float64              `json:"duration"`
	rating       float32              `json:"rating"`
	favorites    int                  `json:"favorites"`
	comments     int                  `json:"comments"`
	isFavorite   bool                 `json:"isFavorite"`
	tags         []heresphereVideoTag `json:"tags"`
}
type heresphereVideoEntryUpdate struct {
	username   string               `json:"username"`
	password   string               `json:"password"`
	isFavorite bool                 `json:"isFavorite"`
	rating     float32              `json:"rating"`
	tags       []heresphereVideoTag `json:"tags"`
	//In base64
	hsp        string `json:"hsp"`
	deleteFile bool   `json:"deleteFile"`
}
type heresphereVideoEvent struct {
	username      string  `json:"username"`
	id            string  `json:"id"`
	title         string  `json:"title"`
	event         int     `json:"event"`
	time          float64 `json:"time"`
	speed         float32 `json:"speed"`
	utc           float64 `json:"utc"`
	connectionKey string  `json:"connectionKey"`
}

type heresphereRoutes struct {
	txnManager txn.Manager
	repository manager.Repository
}

func (rs heresphereRoutes) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(rs.HeresphereCtx)
		r.Post("/", rs.heresphereLogin)

		r.Post("/auth", rs.heresphereLoginToken)
		r.Post("/scan", rs.heresphereScan)
		r.Route("/{sceneId}", func(r chi.Router) {
			r.Post("/", rs.heresphereVideoData)
			//Ours
			r.Get("/hsp", rs.heresphereVideoHsp)
			r.Post("/event", rs.heresphereVideoEvent)
		})
	})

	return r
}

func (rs heresphereRoutes) heresphereVideoEvent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func getVideoTags(scene *models.Scene) []heresphereVideoTag {
	var processedTags []heresphereVideoTag

	for _, sceneTags := range scene.GalleryIDs.List() {
		genTag := heresphereVideoTag{
			name:   fmt.Sprintf("Gallery:%s", sceneTags),
			start:  0.0,
			end:    1.0,
			track:  0,
			rating: 0,
		}
		processedTags = append(processedTags, genTag)
	}

	for _, sceneTags := range scene.TagIDs.List() {
		genTag := heresphereVideoTag{
			name:   fmt.Sprintf("Tag:%s", sceneTags),
			start:  0.0,
			end:    1.0,
			track:  0,
			rating: 0,
		}
		processedTags = append(processedTags, genTag)
	}

	for _, sceneTags := range scene.PerformerIDs.List() {
		genTag := heresphereVideoTag{
			name:   fmt.Sprintf("Talent:%s", sceneTags),
			start:  0.0,
			end:    1.0,
			track:  0,
			rating: 0,
		}
		processedTags = append(processedTags, genTag)
	}

	for _, sceneTags := range scene.Movies.List() {
		genTag := heresphereVideoTag{
			name:   fmt.Sprintf("Movie:%s", sceneTags),
			start:  0.0,
			end:    1.0,
			track:  0,
			rating: 0,
		}
		processedTags = append(processedTags, genTag)
	}

	for _, sceneTags := range scene.StashIDs.List() {
		genTag := heresphereVideoTag{
			name:   fmt.Sprintf("Scene:%s", sceneTags),
			start:  0.0,
			end:    1.0,
			track:  0,
			rating: 0,
		}
		processedTags = append(processedTags, genTag)
	}

	genTag := heresphereVideoTag{
		name:   fmt.Sprintf("Studio:%s", scene.StudioID),
		start:  0.0,
		end:    1.0,
		track:  0,
		rating: 0,
	}
	processedTags = append(processedTags, genTag)

	return processedTags
}
func getVideoScripts(scene *models.Scene) []heresphereVideoScript {
	processedScript := heresphereVideoScript{
		name:   "Default script",
		url:    fmt.Sprintf("/%s/funscript", scene.ID),
		rating: 2.5,
	}
	processedScripts := []heresphereVideoScript{processedScript}
	return processedScripts
}
func getVideoSubtitles(scene *models.Scene) []heresphereVideoSubtitle {
	var processedSubtitles []heresphereVideoSubtitle

	/*var captions []*models.VideoCaption
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

		var buf bytes.Buffer

		err = sub.WriteToWebVTT(&buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/vtt")
		utils.ServeStaticContent(w, r, buf.Bytes())
		return
	}*/

	return processedSubtitles
}
func getVideoMedia(scene *models.Scene) []heresphereVideoMedia {
	var processedMedia []heresphereVideoMedia

	//TODO: Gather all files by Format, use a map/dict
	mediaFile := scene.Files.Primary()
	processedEntry := heresphereVideoMediaSource{
		resolution: mediaFile.Height,
		height:     mediaFile.Height,
		width:      mediaFile.Width,
		size:       mediaFile.Size,
		url:        fmt.Sprintf("/%s/stream"),
	}
	processedSources := []heresphereVideoMediaSource{processedEntry}
	processedMedia = append(processedMedia, heresphereVideoMedia{
		name:    mediaFile.Format,
		sources: processedSources,
	})

	return processedMedia
}

func (rs heresphereRoutes) heresphereScan(w http.ResponseWriter, r *http.Request) {
	var scenes []*models.Scene
	_ = txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		scenes, err = rs.repository.Scene.All(ctx)
		return err
	})
	if scenes == nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Processing each scene and creating a new list
	var processedScenes []heresphereVideoEntryShort
	for _, scene := range scenes {
		isFavorite := *scene.Rating > 85

		// Perform the necessary processing on each scene
		processedScene := heresphereVideoEntryShort{
			link:         fmt.Sprintf("/heresphere/%s", scene.ID),
			title:        scene.Title,
			dateReleased: scene.Date.Format("2006-02-01"),
			duration:     scene.PlayDuration,
			rating:       float32(*scene.Rating) * 0.05, // 0-5
			favorites:    0,
			comments:     scene.OCounter,
			isFavorite:   isFavorite,
			tags:         getVideoTags(scene),
		}
		processedScenes = append(processedScenes, processedScene)
	}

	// Create a JSON encoder for the response writer
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)

	// Write the JSON response
	err := encoder.Encode(processedScenes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.ServeStaticContent(w, r, buf.Bytes())
}

func (rs heresphereRoutes) heresphereVideoHsp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (rs heresphereRoutes) heresphereVideoData(w http.ResponseWriter, r *http.Request) {
	sceneID, err := strconv.Atoi(chi.URLParam(r, "sceneId"))
	if err != nil {
		http.Error(w, "Missing sceneId", http.StatusBadRequest)
		return
	}

	var user heresphereVideoReq
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: This endpoint can receive 2 types of requests
	// One is a video request (heresphereVideoReq)
	// Other is an update (heresphereVideoEntryUpdate)

	var scene *models.Scene
	_ = txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		scene, err = rs.repository.Scene.Find(ctx, sceneID)
		return err
	})
	if scene == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	isFavorite := *scene.Rating > 85
	processedScene := heresphereVideoEntry{
		access:         1,
		title:          scene.Title,
		description:    scene.Details,
		thumbnailImage: "",
		thumbnailVideo: "",
		dateReleased:   scene.Date.Format("2006-02-01"),
		duration:       scene.PlayDuration,
		rating:         float32(*scene.Rating) * 0.05, // 0-5
		favorites:      0,
		comments:       scene.OCounter,
		isFavorite:     isFavorite,
		projection:     heresphereProjectionEquirectangular,
		stereo:         heresphereStereoMono,
		isEyeSwapped:   false,
		fov:            180,
		lens:           heresphereLensLinear,
		cameraIPD:      6.5,
		hsp:            fmt.Sprintf("/heresphere/%s/hsp", scene.ID),
		eventServer:    fmt.Sprintf("/heresphere/%s/event", scene.ID),
		scripts:        getVideoScripts(scene),
		subtitles:      getVideoSubtitles(scene),
		tags:           getVideoTags(scene),
		media:          getVideoMedia(scene),
		writeFavorite:  false,
		writeRating:    false,
		writeTags:      false,
		writeHSP:       false,
	}

	// Create a JSON encoder for the response writer
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)

	// Write the JSON response
	err = encoder.Encode(processedScene)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.ServeStaticContent(w, r, buf.Bytes())
}

func (rs heresphereRoutes) heresphereLogin(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}

	var user heresphereAuthReq
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Auth
	w.WriteHeader(http.StatusOK)
}

func (rs heresphereRoutes) heresphereLoginToken(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}

	var user heresphereAuthReq
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Auth

	auth := &heresphereAuthResp{
		auth_token: "yes",
		access:     heresphereMember,
	}

	// Create a JSON encoder for the response writer
	buf, err := json.Marshal(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.ServeStaticContent(w, r, buf)
}

func (rs heresphereRoutes) HeresphereCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("HereSphere-JSON-Version", strconv.Itoa(heresphereJsonVersion))
		next.ServeHTTP(w, r)
	})
}
