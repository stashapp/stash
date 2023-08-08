package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
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

	HeresphereCustomTagPlayCount HeresphereCustomTag = "PlayCount"
	HeresphereCustomTagWatched   HeresphereCustomTag = "Watched"

	HeresphereCustomTagOrganized HeresphereCustomTag = "Organized"

	HeresphereCustomTagOCounter HeresphereCustomTag = "OCounter"
	HeresphereCustomTagOrgasmed HeresphereCustomTag = "Orgasmed"

	HeresphereCustomTagRated HeresphereCustomTag = "Rated"
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
	Rating float64 `json:"rating,omitempty"`
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
	Rating float64 `json:"rating,omitempty"`
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
	Rating         float64                   `json:"rating,omitempty"`
	Favorites      int                       `json:"favorites"`
	Comments       int                       `json:"comments"`
	IsFavorite     bool                      `json:"isFavorite"`
	Projection     HeresphereProjection      `json:"projection"`
	Stereo         HeresphereStereo          `json:"stereo"`
	IsEyeSwapped   bool                      `json:"isEyeSwapped"`
	Fov            float64                   `json:"fov,omitempty"`
	Lens           HeresphereLens            `json:"lens"`
	CameraIPD      float64                   `json:"cameraIPD"`
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
	Rating       float64              `json:"rating,omitempty"`
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
	Rating           *float64              `json:"rating,omitempty"`
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
	Speed         float64             `json:"speed"`
	Utc           float64             `json:"utc"`
	ConnectionKey string              `json:"connectionKey"`
}

type heresphereRoutes struct {
	txnManager    txn.Manager
	sceneFinder   SceneFinder
	fileFinder    file.Finder
	captionFinder CaptionFinder
	repository    manager.Repository
	resolver      ResolverRoot
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

			r.Post("/event", rs.HeresphereVideoEvent)
		})
	})

	return r
}

func getVrTag() (varTag string, err error) {
	// Find setting
	varTag = config.GetInstance().GetUIVRTag()
	if len(varTag) == 0 {
		err = fmt.Errorf("zero length vr tag")
	}
	return
}
func getMinPlayPercent() (per int, err error) {
	per = config.GetInstance().GetUIMinPlayPercent()
	if per == -1 {
		err = fmt.Errorf("unset minimum play percent")
	}
	return
}
func getFavoriteTag() (varTag string, err error) {
	varTag = config.GetInstance().GetUIFavoriteTag()
	if len(varTag) == 0 {
		// err = fmt.Errorf("zero length favorite tag")
		varTag = "Favorite"
		// TODO: This is for development, remove forced assign
	}
	return
}

/*
 * This auxiliary function searches for the "favorite" tag
 */
func (rs heresphereRoutes) getVideoFavorite(r *http.Request, scene *models.Scene) bool {
	tag_ids, err := rs.resolver.Scene().Tags(r.Context(), scene)
	if err == nil {
		if favTag, err := getFavoriteTag(); err == nil {
			for _, tag := range tag_ids {
				if tag.Name == favTag {
					return true
				}
			}
		}
	}

	return false
}

/*
 * This is a video playback event
 * Intended for server-sided script playback.
 * But since we dont need that, we just use it for timestamps.
 */
func (rs heresphereRoutes) HeresphereVideoEvent(w http.ResponseWriter, r *http.Request) {
	scn := r.Context().Value(heresphereKey).(*models.Scene)

	// Decode event
	var event HeresphereVideoEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add playDuration
	newTime := event.Time / 1000
	/*newDuration := scn.PlayDuration
	// TODO: Huge value bug
	if newTime > scene.ResumeTime {
		newDuration += (newTime - scene.ResumeTime)
	}*/

	// Update PlayCount
	if per, err := getMinPlayPercent(); err == nil {
		if file := scn.Files.Primary(); file != nil && newTime/file.Duration > float64(per)/100.0 {
			// Create update set
			ret := &scene.UpdateSet{
				ID: scn.ID,
			}
			ret.Partial = models.NewScenePartial()

			// Unless we track playback, we cant know if its a "played" or skip event
			if scn.PlayCount == 0 {
				ret.Partial.PlayCount.Set = true
				ret.Partial.PlayCount.Value = 1
			}

			// Update scene
			if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
				_, err := ret.Update(ctx, rs.repository.Scene)
				return err
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	// Write
	// TODO: Datebug still exists
	/*if _, err := rs.resolver.Mutation().SceneSaveActivity(r.Context(), strconv.Itoa(scn.ID), &newTime, &newDuration); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}*/

	w.WriteHeader(http.StatusOK)
}

/*
 * This endpoint is for letting the user update scene data
 */
func (rs heresphereRoutes) HeresphereVideoDataUpdate(w http.ResponseWriter, r *http.Request) error {
	scn := r.Context().Value(heresphereKey).(*models.Scene)
	user := r.Context().Value(heresphereUserKey).(HeresphereAuthReq)
	shouldUpdate := false
	fileDeleter := file.NewDeleter()

	// Create update set
	ret := &scene.UpdateSet{
		ID: scn.ID,
	}
	ret.Partial = models.NewScenePartial()

	// Update rating
	if user.Rating != nil {
		rating := models.Rating5To100F(*user.Rating)
		ret.Partial.Rating = models.NewOptionalInt(rating)
		shouldUpdate = true
	}

	// Delete primary file
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
			return err
		}
		shouldUpdate = true
	}

	// Favorites tag
	if favName, err := getFavoriteTag(); user.IsFavorite != nil && err == nil {
		favTag := HeresphereVideoTag{Name: fmt.Sprintf("Tag:%v", favName)}
		if *user.IsFavorite {
			if user.Tags == nil {
				user.Tags = &[]HeresphereVideoTag{favTag}
			} else {
				*user.Tags = append(*user.Tags, favTag)
			}
		} else if user.Tags != nil {
			for i, tag := range *user.Tags {
				if tag.Name == favTag.Name {
					*user.Tags = append((*user.Tags)[:i], (*user.Tags)[i+1:]...)
					break
				}
			}
		}
		shouldUpdate = true
	}

	// Tags
	if user.Tags != nil {
		// Search input tags and add/create any new ones
		var tagIDs []int
		var perfIDs []int

		for _, tagI := range *user.Tags {
			fmt.Printf("Tag name: %v\n", tagI.Name)

			// If missing
			if len(tagI.Name) == 0 {
				continue
			}

			// If add tag
			// TODO FUTURE: Switch to CutPrefix as it's nicer
			if strings.HasPrefix(tagI.Name, "Tag:") {
				after := strings.TrimPrefix(tagI.Name, "Tag:")
				var err error
				var tagMod *models.Tag
				if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
					// Search for tag
					tagMod, err = rs.repository.Tag.FindByName(ctx, after, true)
					return err
				}); err != nil {
					tagMod = nil
				}

				if tagMod != nil {
					tagIDs = append(tagIDs, tagMod.ID)
				}
			}

			// If add performer
			if strings.HasPrefix(tagI.Name, "Performer:") {
				after := strings.TrimPrefix(tagI.Name, "Performer:")
				var err error
				var tagMod *models.Performer
				if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
					var tagMods []*models.Performer

					// Search for performer
					if tagMods, err = rs.repository.Performer.FindByNames(ctx, []string{after}, true); err == nil && len(tagMods) > 0 {
						tagMod = tagMods[0]
					}

					return err
				}); err != nil {
					tagMod = nil
				}

				if tagMod != nil {
					perfIDs = append(perfIDs, tagMod.ID)
				}
			}

			// If add marker
			if strings.HasPrefix(tagI.Name, "Marker:") {
				after := strings.TrimPrefix(tagI.Name, "Marker:")
				var tagId *string
				if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
					var err error
					var tagMods []*models.MarkerStringsResultType
					searchType := "count"

					// Search for marker
					if tagMods, err = rs.repository.SceneMarker.GetMarkerStrings(ctx, &after, &searchType); err == nil && len(tagMods) > 0 {
						tagId = &tagMods[0].ID

						// Search for tag
						if markers, err := rs.repository.SceneMarker.FindBySceneID(r.Context(), scn.ID); err == nil {
							i, e := strconv.Atoi(*tagId)
							if e == nil {
								// Note: Currently we search if a marker exists.
								// If it doesn't, create it.
								// This also means that markers CANNOT be deleted using the api.
								for _, marker := range markers {
									if marker.Seconds == tagI.Start &&
										marker.SceneID == scn.ID &&
										marker.PrimaryTagID == i {
										tagId = nil
									}
								}
							}
						}
					}

					return err
				}); err != nil {
					// Create marker
					if tagId == nil {
						newTag := SceneMarkerCreateInput{
							Seconds:      tagI.Start,
							SceneID:      strconv.Itoa(scn.ID),
							PrimaryTagID: *tagId,
						}
						if _, err := rs.resolver.Mutation().SceneMarkerCreate(context.Background(), newTag); err != nil {
							return err
						}
					}
				}
			}

			if strings.HasPrefix(tagI.Name, "Movie:") {
				after := strings.TrimPrefix(tagI.Name, "Movie:")

				var err error
				var tagMod *models.Movie
				if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
					// Search for performer
					tagMod, err = rs.repository.Movie.FindByName(ctx, after, true)
					return err
				}); err == nil {
					ret.Partial.MovieIDs.Mode = models.RelationshipUpdateModeSet
					ret.Partial.MovieIDs.AddUnique(models.MoviesScenes{
						MovieID:    tagMod.ID,
						SceneIndex: &scn.ID,
					})
				}
			}
			if strings.HasPrefix(tagI.Name, "Studio:") {
				after := strings.TrimPrefix(tagI.Name, "Studio:")

				var err error
				var tagMod *models.Studio
				if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
					// Search for performer
					tagMod, err = rs.repository.Studio.FindByName(ctx, after, true)
					return err
				}); err == nil {
					ret.Partial.StudioID.Set = true
					ret.Partial.StudioID.Value = tagMod.ID
				}
			}
			if strings.HasPrefix(tagI.Name, "Director:") {
				after := strings.TrimPrefix(tagI.Name, "Director:")
				ret.Partial.Director.Set = true
				ret.Partial.Director.Value = after
			}

			// Custom
			{
				tagName := tagI.Name

				// Will be overwritten if PlayCount tag is updated
				prefix := string(HeresphereCustomTagWatched) + ":"
				if strings.HasPrefix(tagName, prefix) {
					after := strings.TrimPrefix(tagName, prefix)
					if b, err := strconv.ParseBool(after); err == nil {
						// Plays chicken
						if b && scn.PlayCount == 0 {
							ret.Partial.PlayCount.Set = true
							ret.Partial.PlayCount.Value = 1
						} else if !b {
							ret.Partial.PlayCount.Set = true
							ret.Partial.PlayCount.Value = 0
						}
					}
					continue
				}
				prefix = string(HeresphereCustomTagOrganized) + ":"
				if strings.HasPrefix(tagName, prefix) {
					after := strings.TrimPrefix(tagName, prefix)
					if b, err := strconv.ParseBool(after); err == nil {
						ret.Partial.Organized.Set = true
						ret.Partial.Organized.Value = b
					}
					continue
				}
				prefix = string(HeresphereCustomTagRated) + ":"
				if strings.HasPrefix(tagName, prefix) {
					after := strings.TrimPrefix(tagName, prefix)
					if b, err := strconv.ParseBool(after); err == nil && !b {
						ret.Partial.Rating.Set = true
						ret.Partial.Rating.Null = true
					}
					continue
				}

				// Set numbers
				prefix = string(HeresphereCustomTagPlayCount) + ":"
				if strings.HasPrefix(tagName, prefix) {
					after := strings.TrimPrefix(tagName, prefix)
					if numRes, err := strconv.Atoi(after); err != nil {
						ret.Partial.PlayCount.Set = true
						ret.Partial.PlayCount.Value = numRes
					}
					continue
				}
				prefix = string(HeresphereCustomTagOCounter) + ":"
				if strings.HasPrefix(tagName, prefix) {
					after := strings.TrimPrefix(tagName, prefix)
					if numRes, err := strconv.Atoi(after); err != nil {
						ret.Partial.OCounter.Set = true
						ret.Partial.OCounter.Value = numRes
					}
					continue
				}
			}
		}

		// Update tags
		ret.Partial.TagIDs = &models.UpdateIDs{
			IDs:  tagIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
		// Update performers
		ret.Partial.PerformerIDs = &models.UpdateIDs{
			IDs:  perfIDs,
			Mode: models.RelationshipUpdateModeSet,
		}
		shouldUpdate = true
	}

	if shouldUpdate {
		// Update scene
		if err := txn.WithTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			_, err := ret.Update(ctx, rs.repository.Scene)
			return err
		}); err != nil {
			return err
		}

		// Commit to delete file
		fileDeleter.Commit()
		w.WriteHeader(http.StatusOK)
	}
	return nil
}

/*
 * This auxiliary function gathers various tags from the scene to feed the api.
 */
func (rs heresphereRoutes) getVideoTags(r *http.Request, scene *models.Scene) []HeresphereVideoTag {
	processedTags := []HeresphereVideoTag{}

	// Load all relationships
	if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.repository.Scene)
	}); err != nil {
		return processedTags
	}

	if mark_ids, err := rs.resolver.Scene().SceneMarkers(r.Context(), scene); err == nil {
		for _, mark := range mark_ids {
			tagName := mark.Title

			// Add tag name
			if ret, err := loaders.From(r.Context()).TagByID.Load(mark.PrimaryTagID); err == nil {
				if len(tagName) == 0 {
					tagName = ret.Name
				} else {
					tagName = fmt.Sprintf("%s - %s", tagName, ret.Name)
				}
			}

			genTag := HeresphereVideoTag{
				Name:  fmt.Sprintf("Marker:%v", tagName),
				Start: mark.Seconds * 1000,
				End:   (mark.Seconds + 60) * 1000,
			}
			processedTags = append(processedTags, genTag)
		}
	}

	if gallery_ids, err := rs.resolver.Scene().Galleries(r.Context(), scene); err == nil {
		for _, gal := range gallery_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Gallery:%v", gal.GetTitle()),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	if tag_ids, err := rs.resolver.Scene().Tags(r.Context(), scene); err == nil {
		for _, tag := range tag_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Tag:%v", tag.Name),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	if perf_ids, err := rs.resolver.Scene().Performers(r.Context(), scene); err == nil {
		for _, perf := range perf_ids {
			genTag := HeresphereVideoTag{
				Name: fmt.Sprintf("Performer:%s", perf.Name),
			}
			processedTags = append(processedTags, genTag)
		}
	}

	if movie_ids, err := rs.resolver.Scene().Movies(r.Context(), scene); err == nil {
		for _, movie := range movie_ids {
			if movie.Movie != nil {
				genTag := HeresphereVideoTag{
					Name: fmt.Sprintf("Movie:%v", movie.Movie.Name),
				}
				processedTags = append(processedTags, genTag)
			}
		}
	}

	if studio_id, err := rs.resolver.Scene().Studio(r.Context(), scene); err == nil && studio_id != nil {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Studio:%v", studio_id.Name),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		interactive, _ := rs.resolver.Scene().Interactive(r.Context(), scene)
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagInteractive),
				strconv.FormatBool(interactive),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	if len(scene.Director) > 0 {
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("Director:%v", scene.Director),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagWatched),
				strconv.FormatBool(scene.PlayCount > 0),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagOrganized),
				strconv.FormatBool(scene.Organized),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagRated),
				strconv.FormatBool(scene.Rating != nil),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%s:%s",
				string(HeresphereCustomTagOrgasmed),
				strconv.FormatBool(scene.OCounter > 0),
			),
		}
		processedTags = append(processedTags, genTag)
	}

	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%v:%v", string(HeresphereCustomTagPlayCount), scene.PlayCount),
		}
		processedTags = append(processedTags, genTag)
	}
	{
		genTag := HeresphereVideoTag{
			Name: fmt.Sprintf("%v:%v", string(HeresphereCustomTagOCounter), scene.OCounter),
		}
		processedTags = append(processedTags, genTag)
	}

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

	if captions_id, err := rs.resolver.Scene().Captions(r.Context(), scene); err == nil {
		for _, caption := range captions_id {
			processedCaption := HeresphereVideoSubtitle{
				Name:     caption.Filename,
				Language: caption.LanguageCode,
				// & causes chi router bug with \u0026
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
 * This auxiliary function gathers media information + transcoding options.
 */
func (rs heresphereRoutes) getVideoMedia(r *http.Request, scene *models.Scene) []HeresphereVideoMedia {
	processedMedia := []HeresphereVideoMedia{}

	// Codec by source map
	mediaTypes := make(map[string][]HeresphereVideoMediaSource)

	// Load media file
	if err := txn.WithTxn(r.Context(), rs.repository.TxnManager, func(ctx context.Context) error {
		return scene.LoadPrimaryFile(ctx, rs.repository.File)
	}); err != nil {
		return processedMedia
	}

	// If valid primary file
	if mediaFile := scene.Files.Primary(); mediaFile != nil {
		// Get source URL
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

		// Add transcodes
		resRatio := float32(mediaFile.Width) / float32(mediaFile.Height)
		transcodeSize := config.GetInstance().GetMaxStreamingTranscodeSize()
		transNames := []string{"HLS", "DASH"}
		for i, trans := range []string{".m3u8", ".mpd"} {
			for _, res := range models.AllStreamingResolutionEnum {
				maxTrans := transcodeSize.GetMaxResolution()
				// If resolution is below or equal to allowed res (and original video res)
				if height := res.GetMaxResolution(); (maxTrans == 0 || maxTrans >= height) && height <= mediaFile.Height {
					processedEntry.Resolution = height
					processedEntry.Height = height
					processedEntry.Width = int(resRatio * float32(height))
					processedEntry.Size = 0
					// Resolution 0 means original
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

	// Reconstruct tables
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
	// Banner
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

	// Read scenes
	var scenes []*models.Scene
	if err := txn.WithReadTxn(r.Context(), rs.repository.TxnManager, func(ctx context.Context) error {
		var err error
		scenes, err = rs.repository.Scene.All(ctx)
		return err
	}); err != nil {
		http.Error(w, "Failed to fetch scenes!", http.StatusInternalServerError)
		return
	}

	// Create scene list
	sceneUrls := make([]string, len(scenes))
	for idx, scene := range scenes {
		sceneUrls[idx] = fmt.Sprintf("%s/heresphere/%v",
			GetBaseURL(r),
			scene.ID,
		)
	}

	// All library
	library := HeresphereIndexEntry{
		Name: "All",
		List: sceneUrls,
	}
	// Index
	idx := HeresphereIndex{
		Access:  HeresphereMember,
		Banner:  banner,
		Library: []HeresphereIndexEntry{library},
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(idx); err != nil {
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
		tagPre := strings.TrimPrefix(tag.Name, "Tag:")

		// Has degrees tag
		if strings.HasSuffix(tagPre, "°") {
			deg := strings.TrimSuffix(tagPre, "°")
			if s, err := strconv.ParseFloat(deg, 64); err == nil {
				processedScene.Fov = float64(s)
			}
		}
		// Has VR tag
		vrTag, err := getVrTag()
		if err == nil && tagPre == vrTag {
			if processedScene.Projection == HeresphereProjectionPerspective {
				processedScene.Projection = HeresphereProjectionEquirectangular
			}
			if processedScene.Stereo == HeresphereStereoMono {
				processedScene.Stereo = HeresphereStereoSbs
			}
		}
		// Has Fisheye tag
		if tagPre == "Fisheye" {
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

		// Stereo settings
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

		// Projection settings
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

	// Update request
	if err := rs.HeresphereVideoDataUpdate(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch scene
	scene := r.Context().Value(heresphereKey).(*models.Scene)

	// Load relationships
	processedScene := HeresphereVideoEntry{}
	if err := txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
		return scene.LoadRelationships(ctx, rs.repository.Scene)
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create scene
	processedScene = HeresphereVideoEntry{
		Access:         HeresphereMember,
		Title:          scene.GetTitle(),
		Description:    scene.Details,
		ThumbnailImage: addApiKey(urlbuilders.NewSceneURLBuilder(GetBaseURL(r), scene).GetScreenshotURL()),
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

	// Find projection options
	FindProjectionTags(scene, &processedScene)

	// Additional info
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
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(processedScene); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This auxiliary function finds if a login is needed, and auth is correct.
 */
func basicLogin(username string, password string) bool {
	// If needs creds, try login
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This auxiliary function finds if the request has a valid auth token.
 */
func HeresphereHasValidToken(r *http.Request) bool {
	// Check header auth
	apiKey := r.Header.Get(HeresphereAuthHeader)

	// Check url query auth
	if apiKey == "" {
		apiKey = r.URL.Query().Get(session.ApiKeyParameter)
	}

	return len(apiKey) > 0 && apiKey == config.GetInstance().GetAPIKey()
}

/*
 * This auxiliary function adds an auth token to a url
 */
func addApiKey(urlS string) string {
	// Parse URL
	u, err := url.Parse(urlS)
	if err != nil {
		// shouldn't happen
		panic(err)
	}

	// Add apikey if applicable
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
// TODO: Does this even work in HereSphere?
func writeNotAuthorized(w http.ResponseWriter, r *http.Request, msg string) {
	// Banner
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
	// Default video
	library := HeresphereIndexEntry{
		Name: msg,
		List: []string{fmt.Sprintf("%s/heresphere/doesnt-exist", GetBaseURL(r))},
	}
	// Index
	idx := HeresphereIndex{
		Access:  HeresphereBadLogin,
		Banner:  banner,
		Library: []HeresphereIndexEntry{library},
	}

	// Create a JSON encoder for the response writer
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(idx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
 * This context function finds the applicable scene from the request and stores it.
 */
func (rs heresphereRoutes) HeresphereSceneCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get sceneId
		sceneID, err := strconv.Atoi(chi.URLParam(r, "sceneId"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Resolve scene
		var scene *models.Scene
		_ = txn.WithReadTxn(r.Context(), rs.txnManager, func(ctx context.Context) error {
			qb := rs.sceneFinder
			scene, _ = qb.Find(ctx, sceneID)

			if scene != nil {
				// A valid scene should have a attached video
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
		// TODO: Enable and make the settings button
		/*if !config.GetInstance().GetHeresphereDefaultEnabled() {
			writeNotAuthorized(w, r, "HereSphere API not enabled!")
			return
		}*/

		// Add JSON Header (using Add uses camel case and makes it invalid because "Json")
		w.Header()["HereSphere-JSON-Version"] = []string{strconv.Itoa(HeresphereJsonVersion)}

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

		ctx := context.WithValue(r.Context(), heresphereUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
