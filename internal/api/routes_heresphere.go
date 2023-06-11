package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
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
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"

	"golang.org/x/image/draw"
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

type HeresphereEventType int

const (
	HeresphereEventOpen  HeresphereEventType = 0
	HeresphereEventPlay  HeresphereEventType = 1
	HeresphereEventPause HeresphereEventType = 2
	HeresphereEventClose HeresphereEventType = 3
)

type HeresphereCustomTag string

const (
	HeresphereCustomTagInteractive HeresphereCustomTag = "Interactive"
	HeresphereCustomTagUnwatched   HeresphereCustomTag = "Unwatched"
	HeresphereCustomTagWatched     HeresphereCustomTag = "Watched"
	HeresphereCustomTagUnorganized HeresphereCustomTag = "Unorganized"
	HeresphereCustomTagOrganized   HeresphereCustomTag = "Organized"
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
	DateReleased   string                    `json:"dateReleased,omitempty"`
	DateAdded      string                    `json:"dateAdded,omitempty"`
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
	DateReleased string               `json:"dateReleased,omitempty"`
	DateAdded    string               `json:"dateAdded,omitempty"`
	Duration     float64              `json:"duration,omitempty"`
	Rating       float32              `json:"rating,omitempty"`
	Favorites    int                  `json:"favorites"`
	Comments     int                  `json:"comments"`
	IsFavorite   bool                 `json:"isFavorite"`
	Tags         []HeresphereVideoTag `json:"tags"`
}
type HeresphereScanIndex struct {
	ScanData []HeresphereVideoEntryShort `json:"scanData"`
}
type HeresphereAuthReq struct {
	Username         string                `json:"username"`
	Password         string                `json:"password"`
	NeedsMediaSource *bool                 `json:"needsMediaSource,omitempty"`
	IsFavorite       *bool                 `json:"isFavorite,omitempty"`
	Rating           *float32              `json:"rating,omitempty"`
	Tags             *[]HeresphereVideoTag `json:"tags,omitempty"`
	HspBase64        *string               `json:"hsp,omitempty"`
	DeleteFile       *bool                 `json:"deleteFile,omitempty"`
}
type HeresphereVideoEvent struct {
	Username      string              `json:"username"`
	Id            string              `json:"id"`
	Title         string              `json:"title"`
	Event         HeresphereEventType `json:"event"`
	Time          float64             `json:"time"`
	Speed         float32             `json:"speed"`
	Utc           float64             `json:"utc"`
	ConnectionKey string              `json:"connectionKey"`
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

		r.Post("/scan", rs.HeresphereScan)
		r.Post("/auth", rs.HeresphereLoginToken)
		r.Route("/{sceneId}", func(r chi.Router) {
			r.Use(rs.HeresphereSceneCtx)

			r.Post("/", rs.HeresphereVideoData)
			r.Get("/", rs.HeresphereVideoData)

			// r.Get("/hsp", rs.HeresphereVideoHsp)
			r.Post("/event", rs.HeresphereVideoEvent)
			r.Get("/thumbnail", rs.HeresphereThumbnail)
		})
	})

	return r
}

// TODO: Move these to be more generic functions
func getVrTag() string {
	varTag := "Virtual Reality"
	cfgMap := config.GetInstance().GetUIConfiguration()
	if val, ok := cfgMap["vrTag"]; ok {
		rval := val.(string)
		if len(rval) > 0 {
			varTag = rval
		}
	}
	return varTag
}
func getFavoriteTag() string {
	varTag := "Favorite"
	// TODO: .
	return varTag
}
func getHeatmapOverlayEnabled() bool {
	// TODO: .
	return true
}

/*
 * This is a video playback event
 * Intended for server-sided script playback.
 * But since we dont need that, we just use it for timestamps.
 */
func (rs heresphereRoutes) HeresphereVideoEvent(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(heresphereKey).(*models.Scene)

	var event HeresphereVideoEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTime := event.Time / 1000
	newDuration := scene.PlayDuration
	if newTime > scene.ResumeTime {
		newDuration += (newTime - scene.ResumeTime)
	}

	// TODO: Unless we track playing, we cant increment playcount (also check minimumPlayPercent)

	if _, err := rs.resolver.Mutation().SceneSaveActivity(r.Context(), strconv.Itoa(scene.ID), &newTime, &newDuration); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rs heresphereRoutes) HeresphereThumbnail(w http.ResponseWriter, r *http.Request) {
	scene := r.Context().Value(heresphereKey).(*models.Scene)
	defaultUrl := addApiKey(urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetScreenshotURL())

	if !getHeatmapOverlayEnabled() {
		http.Redirect(w, r, defaultUrl, http.StatusSeeOther)
		return
	}

	sceneHash := scene.GetHash(config.GetInstance().GetVideoFileNamingAlgorithm())
	filePath := manager.GetInstance().Paths.Scene.GetInteractiveHeatmapPath(sceneHash)

	// Get cover image
	var cover []byte
	if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		sceneCoverGetter := (rs.sceneFinder).(manager.SceneCoverGetter)
		var err error
		cover, err = sceneCoverGetter.GetCover(ctx, scene.ID)
		return err
	}); err != nil {
		logger.Debugf("Scene %v failed to get cover\n", scene.ID)
		http.Redirect(w, r, defaultUrl, http.StatusSeeOther)
		return
	}

	// Heatmap open
	f, err := os.Open(filePath)
	if err != nil {
		logger.Debugf("Scene %v failed to open heatmap\n", scene.ID)
		utils.ServeImage(w, r, cover)
		return
	}
	defer f.Close()

	// Cover decode
	coverDec, err := jpeg.Decode(bytes.NewBuffer(cover))
	if err != nil {
		logger.Debugf("Scene %v failed to decode cover\n", scene.ID)
		utils.ServeImage(w, r, cover)
		return
	}

	// Heatmap decode
	heatMapDec, err := png.Decode(f)
	if err != nil {
		logger.Debugf("Scene %v failed to decode heatmap\n", scene.ID)
		utils.ServeImage(w, r, cover)
		return
	}

	// Calculate the new width and height based on the desired ratio
	newWidth := models.DefaultGthumbWidth
	newHeight := int(float64(newWidth) / 16.0 * 9.0)

	// Calculate the height for pasting heatMapDec onto coverDec
	pasteHeight := int(float64(newHeight) * 0.10)

	// Calculate the scaled dimensions of the coverDec image
	// TODO: The ignores segment the image
	// Consider adding back
	scale := math.Min(float64(newWidth)/float64(coverDec.Bounds().Dx()), float64(newHeight /*-pasteHeight*/)/float64(coverDec.Bounds().Dy()))
	scaledWidth := int(float64(coverDec.Bounds().Dx()) * scale)
	scaledHeight := int(float64(coverDec.Bounds().Dy()) * scale)

	// Calculate the position to center the scaled coverDec image
	x := (newWidth - scaledWidth) / 2
	y := (newHeight /*- pasteHeight*/ - scaledHeight) / 2

	// Create a new image with the specified width and height
	newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize and draw coverDec onto the new image
	coverResized := image.NewRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))
	draw.CatmullRom.Scale(coverResized, coverResized.Bounds(), coverDec, coverDec.Bounds(), draw.Src, nil)

	// Paste the resized coverDec onto the new image
	draw.Draw(newImage, image.Rect(x, y, x+scaledWidth, y+scaledHeight), coverResized, coverResized.Bounds().Min, draw.Src)

	// Calculate the destination rectangle for pasting heatMapDec
	pasteRect := image.Rect(0, newHeight-pasteHeight, newWidth, newHeight)

	// Resize heatMapDec to fit the paste rectangle
	resizedHeatMap := image.NewRGBA(pasteRect)
	draw.CatmullRom.Scale(resizedHeatMap, resizedHeatMap.Bounds(), heatMapDec, heatMapDec.Bounds(), draw.Over, nil)

	// Paste heatMapDec onto the new image
	draw.Draw(newImage, pasteRect, resizedHeatMap, resizedHeatMap.Bounds().Min, draw.Over)

	// Encode overlayed image
	var b bytes.Buffer
	iw := bufio.NewWriter(&b)
	if err := jpeg.Encode(iw, newImage, &jpeg.Options{Quality: 90}); err != nil {
		logger.Debugf("Scene %v failed to encode image with heatmap overlaid\n", scene.ID)
		utils.ServeImage(w, r, cover)
		return
	}

	// Flush buffer
	if err := iw.Flush(); err != nil {
		logger.Debugf("Scene %v failed to flush encoded heatmap overlay\n", scene.ID)
		utils.ServeImage(w, r, cover)
		return
	}

	utils.ServeImage(w, r, b.Bytes())
}

/*
 * HSP is a HereSphere config file
 * It stores the players local config such as projection or color settings etc.
 */
func (rs heresphereRoutes) HeresphereVideoHsp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	// TODO: Need an SQL entry, either link to file (annoying path, needs cleanup etc.), or just put binary in SQL
}

/*
 * This endpoint provides a list of all videos in a short format
 */
func (rs heresphereRoutes) HeresphereScan(w http.ResponseWriter, r *http.Request) {
	processedScenes := []HeresphereVideoEntryShort{}
	var scenes []*models.Scene

	if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		var err error
		scenes, err = rs.repository.Scene.All(ctx)
		return err
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Processing each scene and creating a new list
	for _, scene := range scenes {
		// Perform the necessary processing on each scene
		processedScene := HeresphereVideoEntryShort{
			Link: fmt.Sprintf("%s/heresphere/%v",
				GetBaseURL(r),
				scene.ID,
			),
			Title:      scene.GetTitle(),
			DateAdded:  scene.CreatedAt.Format("2006-01-02"),
			Duration:   60000.0,
			Rating:     0,
			Favorites:  0,
			Comments:   scene.OCounter,
			IsFavorite: rs.getVideoFavorite(r, scene),
			Tags:       rs.getVideoTags(r, scene),
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
				processedScene.Duration = handleFloat64Value(file_ids.Duration * 1000.0)
			}
		}

		processedScenes = append(processedScenes, processedScene)
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(HeresphereScanIndex{ScanData: processedScenes})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This endpoint is for letting the user update scene data
 */
func (rs heresphereRoutes) HeresphereVideoDataUpdate(w http.ResponseWriter, r *http.Request) {
	scn := r.Context().Value(heresphereKey).(*models.Scene)
	user := r.Context().Value(heresphereUserKey).(HeresphereAuthReq)
	fileDeleter := file.NewDeleter()

	ret := &scene.UpdateSet{
		ID: scn.ID,
	}
	ret.Partial = models.NewScenePartial()

	if user.Rating != nil {
		rating := models.Rating5To100F(*user.Rating)
		ret.Partial.Rating = models.NewOptionalInt(rating)
	}

	if user.DeleteFile != nil && *user.DeleteFile {
		if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			fqb := rs.repository.File

			if err := scn.LoadPrimaryFile(ctx, fqb); err != nil {
				return err
			}

			ff := scn.Files.Primary()
			if ff != nil {
				if err := file.Destroy(ctx, fqb, ff, fileDeleter, true); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			fileDeleter.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	favName := getFavoriteTag()
	if user.IsFavorite != nil {
		favtag := HeresphereVideoTag{Name: fmt.Sprintf("Tag:%v", favName)}
		if *user.IsFavorite {
			if user.Tags == nil {
				user.Tags = &[]HeresphereVideoTag{favtag}
			} else {
				*user.Tags = append(*user.Tags, favtag)
			}
		} else if user.Tags != nil {
			for i, tag := range *user.Tags {
				if tag.Name == favtag.Name {
					*user.Tags = append((*user.Tags)[:i], (*user.Tags)[i+1:]...)
					break
				}
			}
		}
	}

	if user.Tags != nil {
		// Search input tags and add/create any new ones
		var tagIDs []int
		var perfIDs []int

		for _, tagI := range *user.Tags {
			fmt.Printf("Tag name: %v\n", tagI.Name)

			if len(tagI.Name) == 0 {
				continue
			}

			if strings.HasPrefix(tagI.Name, "Tag:") {
				tagName := strings.TrimPrefix(tagI.Name, "Tag:")

				// TODO: How to increment
				if tagName == string(HeresphereCustomTagWatched) && scn.PlayCount == 0 {
					scn.PlayCount++
					continue
				}
				if tagName == string(HeresphereCustomTagUnwatched) {
					scn.PlayCount = 0
					continue
				}
				if tagName == string(HeresphereCustomTagOrganized) {
					scn.Organized = true
					continue
				}
				if tagName == string(HeresphereCustomTagUnorganized) {
					scn.Organized = false
					continue
				}
				if tagName == favName {
					continue
				}

				var err error
				var tagMod *models.Tag
				if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
					tagMod, err = rs.repository.Tag.FindByName(ctx, tagName, true)
					return err
				}); err != nil || tagMod == nil {
					newTag := TagCreateInput{
						Name: tagName,
					}
					if tagMod, err = rs.resolver.Mutation().TagCreate(r.Context(), newTag); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}

				tagIDs = append(tagIDs, tagMod.ID)
			}
			if strings.HasPrefix(tagI.Name, "Performer:") {
				tagName := strings.TrimPrefix(tagI.Name, "Performer:")

				var err error
				var tagMod *models.Performer
				if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
					var tagMods []*models.Performer
					if tagMods, err = rs.repository.Performer.FindByNames(ctx, []string{tagName}, true); err == nil && len(tagMods) > 0 {
						tagMod = tagMods[0]
					}
					return err
				}); err != nil || tagMod == nil {
					newTag := PerformerCreateInput{
						Name: tagName,
					}
					if tagMod, err = rs.resolver.Mutation().PerformerCreate(r.Context(), newTag); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}

				perfIDs = append(perfIDs, tagMod.ID)
			}
		}

		ret.Partial.TagIDs = &models.UpdateIDs{
			IDs:  tagIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
		ret.Partial.PerformerIDs = &models.UpdateIDs{
			IDs:  perfIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
	}

	if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		_, err := ret.Update(ctx, rs.repository.Scene)
		return err
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileDeleter.Commit()
	w.WriteHeader(http.StatusOK)
}

/*
 * This auxiliary function gathers various tags from the scene to feed the api.
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
			if tag.Name == getFavoriteTag() {
				continue
			}

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
				Name: fmt.Sprintf("Performer:%s", perf.Name),
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

	studio_id, err := rs.resolver.Scene().Studio(r.Context(), scene)
	if err == nil && studio_id != nil {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Studio:%v", studio_id.Name.String),
		}
		processedTags = append(processedTags, genTag)
	}

	interactive, err := rs.resolver.Scene().Interactive(r.Context(), scene)
	if err == nil && interactive {
		genTag := HeresphereVideoTag{
			Name: string(HeresphereCustomTagInteractive),
		}
		processedTags = append(processedTags, genTag)
	}

	if len(scene.Director) > 0 {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Director:%v", scene.Director),
		}
		processedTags = append(processedTags, genTag)
	}

	if scene.PlayCount > 0 {
		genTag := HeresphereVideoTag{
			Name: string(HeresphereCustomTagWatched),
		}
		processedTags = append(processedTags, genTag)
	} else {
		genTag := HeresphereVideoTag{
			Name: string(HeresphereCustomTagUnwatched),
		}
		processedTags = append(processedTags, genTag)
	}

	if scene.Organized {
		genTag := HeresphereVideoTag{
			Name: string(HeresphereCustomTagOrganized),
		}
		processedTags = append(processedTags, genTag)
	} else {
		genTag := HeresphereVideoTag{
			Name: string(HeresphereCustomTagUnorganized),
		}
		processedTags = append(processedTags, genTag)
	}

	// TODO: PlayCount tag (replace watch/unwatched?, just set to number like PlayCount:5 and let user edit)
	// TODO: OCount tag
	// TODO: More?

	return processedTags
}

/*
 * This auxiliary function gathers a script if applicable
 */
func (rs heresphereRoutes) getVideoScripts(r *http.Request, scene *models.Scene) []HeresphereVideoScript {
	processedScripts := []HeresphereVideoScript{}

	if interactive, err := rs.resolver.Scene().Interactive(r.Context(), scene); err == nil && interactive {
		processedScript := HeresphereVideoScript{
			Name:   "Default script",
			Url:    addApiKey(urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetFunscriptURL()),
			Rating: 5,
		}
		processedScripts = append(processedScripts, processedScript)
	}

	return processedScripts
}

/*
 * This auxiliary function gathers subtitles if applicable
 */
func (rs heresphereRoutes) getVideoSubtitles(r *http.Request, scene *models.Scene) []HeresphereVideoSubtitle {
	processedSubtitles := []HeresphereVideoSubtitle{}

	// TODO: Seemingly broken in heresphere
	if captions_id, err := rs.resolver.Scene().Captions(r.Context(), scene); err == nil {
		for _, caption := range captions_id {
			processedCaption := HeresphereVideoSubtitle{
				Name:     caption.Filename,
				Language: caption.LanguageCode,
				Url: addApiKey(fmt.Sprintf("%s?lang=%v&type=%v",
					urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetCaptionURL(),
					caption.LanguageCode,
					caption.CaptionType,
				)),
			}
			processedSubtitles = append(processedSubtitles, processedCaption)
		}
	}

	return processedSubtitles
}

/*
 * This auxiliary function searches for the "favorite" tag
 */
func (rs heresphereRoutes) getVideoFavorite(r *http.Request, scene *models.Scene) bool {
	tag_ids, err := rs.resolver.Scene().Tags(r.Context(), scene)
	if err == nil {
		favTag := getFavoriteTag()
		for _, tag := range tag_ids {
			if tag.Name == favTag {
				return true
			}
		}
	}

	return false
}

/*
 * This auxiliary function gathers media information + transcoding options.
 */
func (rs heresphereRoutes) getVideoMedia(r *http.Request, scene *models.Scene) []HeresphereVideoMedia {
	processedMedia := []HeresphereVideoMedia{}

	mediaTypes := make(map[string][]HeresphereVideoMediaSource)

	if err := txn.WithTxn(r.Context(), rs.repository.TxnManager, func(ctx context.Context) error {
		return scene.LoadPrimaryFile(ctx, rs.repository.File)
	}); err != nil {
		return processedMedia
	}

	if mediaFile := scene.Files.Primary(); mediaFile != nil {
		sourceUrl := urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetStreamURL("").String()
		processedEntry := HeresphereVideoMediaSource{
			Resolution: mediaFile.Height,
			Height:     mediaFile.Height,
			Width:      mediaFile.Width,
			Size:       mediaFile.Size,
			Url:        addApiKey(sourceUrl),
		}
		processedMedia = append(processedMedia, HeresphereVideoMedia{
			Name:    "direct stream",
			Sources: []HeresphereVideoMediaSource{processedEntry},
		})

		resRatio := float32(mediaFile.Width) / float32(mediaFile.Height)
		transcodeSize := config.GetInstance().GetMaxStreamingTranscodeSize()
		transNames := []string{"HLS", "DASH"}
		for i, trans := range []string{".m3u8", ".mpd"} {
			for _, res := range models.AllStreamingResolutionEnum {
				maxTrans := transcodeSize.GetMaxResolution()
				if height := res.GetMaxResolution(); (maxTrans == 0 || maxTrans >= height) && height <= mediaFile.Height {
					processedEntry.Resolution = height
					processedEntry.Height = height
					processedEntry.Width = int(resRatio * float32(height))
					processedEntry.Size = 0
					if height == 0 {
						processedEntry.Resolution = mediaFile.Height
						processedEntry.Height = mediaFile.Height
						processedEntry.Width = mediaFile.Width
					}
					processedEntry.Url = addApiKey(fmt.Sprintf("%s%s?resolution=%s", sourceUrl, trans, res.String()))

					typeName := transNames[i]
					mediaTypes[typeName] = append(mediaTypes[typeName], processedEntry)
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
	err := json.NewEncoder(w).Encode(idx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This auxiliary function finds vr projection modes from tags and the filename.
 */
func FindProjectionTags(scene *models.Scene, processedScene *HeresphereVideoEntry) {
	// Detect VR modes from tags
	for _, tag := range processedScene.Tags {
		if strings.Contains(tag.Name, "°") {
			deg := strings.TrimSuffix(tag.Name, "°")
			deg = strings.TrimPrefix(deg, "Tag:")
			if s, err := strconv.ParseFloat(deg, 32); err == nil {
				processedScene.Fov = float32(s)
			}
		}
		if strings.Contains(tag.Name, getVrTag()) {
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
		if strings.Contains(path, "_FISHEYE") {
			processedScene.Projection = HeresphereProjectionFisheye
		}
		if strings.Contains(path, "_RF52") || strings.Contains(path, "_FISHEYE190") {
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

	if user.DeleteFile != nil || user.Rating != nil || user.Tags != nil {
		rs.HeresphereVideoDataUpdate(w, r)
	}

	scene := r.Context().Value(heresphereKey).(*models.Scene)

	processedScene := HeresphereVideoEntry{}
	if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.repository.Scene)
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Use image lib to overlay heatmap?
	processedScene = HeresphereVideoEntry{
		Access:      HeresphereMember,
		Title:       scene.GetTitle(),
		Description: scene.Details,
		ThumbnailImage: addApiKey(fmt.Sprintf("%s/heresphere/%v/thumbnail",
			GetBaseURL(r),
			scene.ID,
		)),
		ThumbnailVideo: addApiKey(urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetStreamPreviewURL()),
		DateAdded:      scene.CreatedAt.Format("2006-01-02"),
		Duration:       60000.0,
		Rating:         0,
		Favorites:      0,
		Comments:       scene.OCounter,
		IsFavorite:     rs.getVideoFavorite(r, scene),
		Projection:     HeresphereProjectionPerspective,
		Stereo:         HeresphereStereoMono,
		IsEyeSwapped:   false,
		Fov:            180.0,
		Lens:           HeresphereLensLinear,
		CameraIPD:      6.5,
		EventServer: addApiKey(fmt.Sprintf("%s/heresphere/%v/event",
			GetBaseURL(r),
			scene.ID,
		)),
		Scripts:       rs.getVideoScripts(r, scene),
		Subtitles:     rs.getVideoSubtitles(r, scene),
		Tags:          rs.getVideoTags(r, scene),
		Media:         []HeresphereVideoMedia{},
		WriteFavorite: true,
		WriteRating:   true,
		WriteTags:     true,
		WriteHSP:      false,
	}
	FindProjectionTags(scene, &processedScene)

	if user.NeedsMediaSource != nil && *user.NeedsMediaSource {
		processedScene.Media = rs.getVideoMedia(r, scene)
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
			processedScene.Duration = handleFloat64Value(file_ids.Duration * 1000.0)
		}
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(processedScene)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This auxiliary function finds if a login is needed, and auth is correct.
 */
// TODO: Move to utils?
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
	err := json.NewEncoder(w).Encode(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This auxiliary function finds if the request has a valid auth token.
 */
func HeresphereHasValidToken(r *http.Request) bool {
	apiKey := r.Header.Get(HeresphereAuthHeader)

	if apiKey == "" {
		apiKey = r.URL.Query().Get(session.ApiKeyParameter)
	}

	return len(apiKey) > 0 && apiKey == config.GetInstance().GetAPIKey()
}

/*
 * This auxiliary function adds an auth token to a url
 */
// TODO: Move this to utils
func addApiKey(urlS string) string {
	u, err := url.Parse(urlS)
	if err != nil {
		// shouldn't happen
		panic(err)
	}

	if config.GetInstance().GetAPIKey() != "" {
		v := u.Query()
		if !v.Has("apikey") {
			v.Set("apikey", config.GetInstance().GetAPIKey())
		}
		u.RawQuery = v.Encode()
	}

	return u.String()
}

/*
 * This auxiliary writes a library with a fake name upon auth failure
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
	err := json.NewEncoder(w).Encode(idx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		isAuth := config.GetInstance().HasCredentials() && !HeresphereHasValidToken(r)
		needsMedia := true
		user := HeresphereAuthReq{NeedsMediaSource: &needsMedia}
		if err := json.Unmarshal(body, &user); err != nil && isAuth {
			writeNotAuthorized(w, r, "Not logged in!")
			return
		}

		if isAuth && !strings.HasPrefix(r.URL.Path, "/heresphere/auth") {
			writeNotAuthorized(w, r, "Unauthorized!")
			return
		}

		ctx := context.WithValue(r.Context(), heresphereUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
