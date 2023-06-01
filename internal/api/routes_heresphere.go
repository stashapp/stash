package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/manager"
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
	Access int `json:"access"`
	//Banner  HeresphereBanner       `json:"banner"`
	Library []HeresphereIndexEntry `json:"library"`
}
type HeresphereVideoScript struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	//Rating float32 `json:"rating,omitempty"`
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
	Access         int    `json:"access"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	ThumbnailImage string `json:"thumbnailImage"`
	ThumbnailVideo string `json:"thumbnailVideo,omitempty"`
	DateReleased   string `json:"dateReleased"`
	DateAdded      string `json:"dateAdded"`
	Duration       uint   `json:"duration,omitempty"`
	/*Rating         float32                   `json:"rating,omitempty"`
	Favorites      float32                   `json:"favorites,omitempty"`
	Comments       int                       `json:"comments"`*/
	IsFavorite bool                 `json:"isFavorite"`
	Projection HeresphereProjection `json:"projection"`
	Stereo     HeresphereStereo     `json:"stereo"`
	//IsEyeSwapped   bool                      `json:"isEyeSwapped"`
	Fov  float32        `json:"fov,omitempty"`
	Lens HeresphereLens `json:"lens"`
	//CameraIPD      float32                   `json:"cameraIPD,omitempty"`
	/*Hsp            string                    `json:"hsp,omitempty"`
	EventServer    string                    `json:"eventServer,omitempty"`*/
	Scripts       []HeresphereVideoScript   `json:"scripts,omitempty"`
	Subtitles     []HeresphereVideoSubtitle `json:"subtitles,omitempty"`
	Tags          []HeresphereVideoTag      `json:"tags,omitempty"`
	Media         []HeresphereVideoMedia    `json:"media,omitempty"`
	WriteFavorite bool                      `json:"writeFavorite"`
	WriteRating   bool                      `json:"writeRating"`
	WriteTags     bool                      `json:"writeTags"`
	WriteHSP      bool                      `json:"writeHSP"`
}
type HeresphereVideoEntryShort struct {
	Link         string `json:"link"`
	Title        string `json:"title"`
	DateReleased string `json:"dateReleased"`
	DateAdded    string `json:"dateAdded"`
	Duration     uint   `json:"duration,omitempty"`
	/*Rating       float32              `json:"rating,omitempty"`
	Favorites    int                  `json:"favorites"`
	Comments     int                  `json:"comments"`*/
	IsFavorite bool                 `json:"isFavorite"`
	Tags       []HeresphereVideoTag `json:"tags"`
}
type HeresphereVideoEntryUpdate struct {
	Username   string               `json:"username"`
	Password   string               `json:"password"`
	IsFavorite bool                 `json:"isFavorite"`
	Rating     float32              `json:"rating,omitempty"`
	Tags       []HeresphereVideoTag `json:"tags"`
	//In base64
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
			//Ours
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

// TODO: This probably isnt necessary
func relUrlToAbs(r *http.Request, rel string) string {
	// Get the scheme (http or https) from the request
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Get the host from the request
	host := r.Host

	// Combine the scheme, host, and relative path to form the absolute URL
	return fmt.Sprintf("%s://%s%s", scheme, host, rel)
}

func (rs heresphereRoutes) getVideoTags(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoTag {
	processedTags := []HeresphereVideoTag{}

	testTag := HeresphereVideoTag{
		Name: "Test:a",
		/*Start:  0.0,
		End:    1.0,
		Track:  0,
		Rating: 0,*/
	}
	processedTags = append(processedTags, testTag)

	if scene.LoadRelationships(ctx, rs.repository.Scene) == nil {
		if scene.GalleryIDs.Loaded() {
			for _, sceneTags := range scene.GalleryIDs.List() {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Gallery:%v", sceneTags),
					/*Start:  0.0,
					End:    1.0,
					Track:  0,
					Rating: 0,*/
				}
				processedTags = append(processedTags, genTag)
			}
		}

		if scene.TagIDs.Loaded() {
			for _, sceneTags := range scene.TagIDs.List() {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Tag:%v", sceneTags),
					/*Start:  0.0,
					End:    1.0,
					Track:  0,
					Rating: 0,*/
				}
				processedTags = append(processedTags, genTag)
			}
		}

		if scene.PerformerIDs.Loaded() {
			for _, sceneTags := range scene.PerformerIDs.List() {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Talent:%s", sceneTags),
					/*Start:  0.0,
					End:    1.0,
					Track:  0,
					Rating: 0,*/
				}
				processedTags = append(processedTags, genTag)
			}
		}

		if scene.Movies.Loaded() {
			for _, sceneTags := range scene.Movies.List() {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Movie:%v", sceneTags),
					/*Start:  0.0,
					End:    1.0,
					Track:  0,
					Rating: 0,*/
				}
				processedTags = append(processedTags, genTag)
			}
		}

		if scene.StashIDs.Loaded() {
			//TODO: Markers have timestamps?
			for _, sceneTags := range scene.StashIDs.List() {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Scene:%v", sceneTags),
					/*Start:  0.0,
					End:    1.0,
					Track:  0,
					Rating: 0,*/
				}
				processedTags = append(processedTags, genTag)
			}
		}
	}

	if scene.StudioID != nil {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Studio:%v", scene.StudioID),
			/*Start:  0.0,
			End:    1.0,
			Track:  0,
			Rating: 0,*/
		}
		processedTags = append(processedTags, genTag)
	}

	return processedTags
}
func (rs heresphereRoutes) getVideoScripts(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoScript {
	//TODO: Check if exists
	processedScript := HeresphereVideoScript{
		Name: "Default script",
		Url:  relUrlToAbs(r, fmt.Sprintf("/scene/%v/funscript", scene.ID)),
		//Rating: 2.5,
	}
	processedScripts := []HeresphereVideoScript{processedScript}
	return processedScripts
}
func (rs heresphereRoutes) getVideoSubtitles(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoSubtitle {
	processedSubtitles := []HeresphereVideoSubtitle{}

	//TODO: /scene/123/caption?lang=00&type=srt

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
func (rs heresphereRoutes) getVideoMedia(ctx context.Context, r *http.Request, scene *models.Scene) []HeresphereVideoMedia {
	processedMedia := []HeresphereVideoMedia{}

	//TODO: Gather also secondary files, by Format, use a map/dict
	mediaFile := scene.Files.Primary()
	if mediaFile != nil {
		processedEntry := HeresphereVideoMediaSource{
			Resolution: mediaFile.Height,
			Height:     mediaFile.Height,
			Width:      mediaFile.Width,
			Size:       mediaFile.Size,
			Url:        relUrlToAbs(r, fmt.Sprintf("/scene/%v/stream", scene.ID)),
		}
		processedSources := []HeresphereVideoMediaSource{processedEntry}
		processedMedia = append(processedMedia, HeresphereVideoMedia{
			Name:    mediaFile.Format,
			Sources: processedSources,
		})
	}
	//TODO: Transcode etc. /scene/%v/stream.mp4?resolution=ORIGINAL

	return processedMedia
}

func (rs heresphereRoutes) HeresphereIndex(w http.ResponseWriter, r *http.Request) {
	/*banner := HeresphereBanner{
		Image: relUrlToAbs(r, "/apple-touch-icon.png"),
		Link:  relUrlToAbs(r, "/"),
	}*/

	var scenes []*models.Scene
	if err := rs.repository.WithTxn(r.Context(), func(ctx context.Context) error {
		var err error
		scenes, err = rs.repository.Scene.All(ctx)
		return err
	}); err != nil {
		http.Error(w, "Failed to fetch scenes!", http.StatusInternalServerError)
		return
	}

	var sceneUrls []string
	for _, scene := range scenes {
		sceneUrls = append(sceneUrls, relUrlToAbs(r, fmt.Sprintf("/heresphere/%v", scene.ID)))
	}

	library := HeresphereIndexEntry{
		Name: "All",
		List: sceneUrls,
	}
	idx := HeresphereIndex{
		Access: HeresphereMember,
		//Banner:  banner,
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
	//TODO: Auth
	var scenes []*models.Scene
	if err := rs.repository.WithTxn(r.Context(), func(ctx context.Context) error {
		var err error
		scenes, err = rs.repository.Scene.All(ctx)
		return err
	}); err != nil {
		http.Error(w, "Failed to fetch scenes!", http.StatusInternalServerError)
		return
	}

	// Processing each scene and creating a new list
	var processedScenes []HeresphereVideoEntryShort
	for _, scene := range scenes {
		// Perform the necessary processing on each scene
		processedScene := HeresphereVideoEntryShort{
			Link:         relUrlToAbs(r, fmt.Sprintf("/heresphere/%v", scene.ID)),
			Title:        scene.GetTitle(),
			DateReleased: scene.CreatedAt.Format("2006-01-02"),
			DateAdded:    scene.CreatedAt.Format("2006-01-02"),
			Duration:     60,
			/*Rating:       2.5,
			Favorites:    0,
			Comments:     scene.OCounter,*/
			IsFavorite: false,
			Tags:       rs.getVideoTags(r.Context(), r, scene),
		}
		if scene.Date != nil {
			processedScene.DateReleased = scene.Date.Format("2006-01-02")
		}
		if scene.Rating != nil {
			isFavorite := *scene.Rating > 85
			//processedScene.Rating = float32(*scene.Rating) * 0.05 // 0-5
			processedScene.IsFavorite = isFavorite
		}
		//TODO: panic="relationship has not been loaded"
		/*primaryFile := scene.Files.Primary()
		if primaryFile != nil {
			processedScene.Duration = handleFloat64Value(primaryFile.Duration)
		}*/
		processedScenes = append(processedScenes, processedScene)
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

func (rs heresphereRoutes) HeresphereVideoHsp(w http.ResponseWriter, r *http.Request) {
	//TODO: Auth
	w.WriteHeader(http.StatusNotImplemented)
}

func (rs heresphereRoutes) HeresphereVideoDataUpdate(w http.ResponseWriter, r *http.Request) {
	//TODO: This
	//TODO: Auth

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

	//TODO: Auth

	processedScene := HeresphereVideoEntry{
		Access:         HeresphereMember,
		Title:          scene.GetTitle(),
		Description:    scene.Details,
		ThumbnailImage: relUrlToAbs(r, fmt.Sprintf("/scene/%v/screenshot", scene.ID)),
		ThumbnailVideo: relUrlToAbs(r, fmt.Sprintf("/scene/%v/preview", scene.ID)),
		DateReleased:   scene.CreatedAt.Format("2006-01-02"),
		DateAdded:      scene.CreatedAt.Format("2006-01-02"),
		Duration:       60.0,
		/*Rating:         2.5,
		Favorites:      0,
		Comments:       scene.OCounter,*/
		IsFavorite: false,
		Projection: HeresphereProjectionPerspective, // Default to flat cause i have no idea
		Stereo:     HeresphereStereoMono,            // Default to flat cause i have no idea
		//IsEyeSwapped: false,
		Fov:  180,
		Lens: HeresphereLensLinear,
		/*CameraIPD:      6.5,
		Hsp:            relUrlToAbs(r, fmt.Sprintf("/heresphere/%v/hsp", scene.ID)),
		EventServer:    relUrlToAbs(r, fmt.Sprintf("/heresphere/%v/event", scene.ID)),*/
		Scripts:       rs.getVideoScripts(r.Context(), r, scene),
		Subtitles:     rs.getVideoSubtitles(r.Context(), r, scene),
		Tags:          rs.getVideoTags(r.Context(), r, scene),
		Media:         []HeresphereVideoMedia{},
		WriteFavorite: false,
		WriteRating:   false,
		WriteTags:     false,
		WriteHSP:      false,
	}

	if user.NeedsMediaSource {
		processedScene.Media = rs.getVideoMedia(r.Context(), r, scene)
	}
	if scene.Date != nil {
		processedScene.DateReleased = scene.Date.Format("2006-01-02")
	}
	if scene.Rating != nil {
		isFavorite := *scene.Rating > 85
		//processedScene.Rating = float32(*scene.Rating) * 0.05 // 0-5
		processedScene.IsFavorite = isFavorite
	}
	primaryFile := scene.Files.Primary()
	if primaryFile != nil {
		processedScene.Duration = uint(handleFloat64Value(primaryFile.Duration))
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
	if r.Body == nil {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}

	var user HeresphereAuthReq
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Auth
	rs.HeresphereIndex(w, r)
}

func (rs heresphereRoutes) HeresphereLoginToken(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}

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
