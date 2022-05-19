package scene

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
	"github.com/stashapp/stash/pkg/utils"
)

type CoverStashIDGetter interface {
	GetCover(ctx context.Context, sceneID int) ([]byte, error)
	GetStashIDs(ctx context.Context, sceneID int) ([]*models.StashID, error)
}

type MovieGetter interface {
	GetMovies(ctx context.Context, sceneID int) ([]models.MoviesScenes, error)
}

type MarkerTagFinder interface {
	tag.Finder
	TagFinder
	FindBySceneMarkerID(ctx context.Context, sceneMarkerID int) ([]*models.Tag, error)
}

type MarkerFinder interface {
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneMarker, error)
}

type TagFinder interface {
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.Tag, error)
}

// ToBasicJSON converts a scene object into its JSON object equivalent. It
// does not convert the relationships to other objects, with the exception
// of cover image.
func ToBasicJSON(ctx context.Context, reader CoverStashIDGetter, scene *models.Scene) (*jsonschema.Scene, error) {
	newSceneJSON := jsonschema.Scene{
		CreatedAt: json.JSONTime{Time: scene.CreatedAt.Timestamp},
		UpdatedAt: json.JSONTime{Time: scene.UpdatedAt.Timestamp},
	}

	if scene.Checksum.Valid {
		newSceneJSON.Checksum = scene.Checksum.String
	}

	if scene.OSHash.Valid {
		newSceneJSON.OSHash = scene.OSHash.String
	}

	if scene.Phash.Valid {
		newSceneJSON.Phash = utils.PhashToString(scene.Phash.Int64)
	}

	if scene.Title.Valid {
		newSceneJSON.Title = scene.Title.String
	}

	if scene.URL.Valid {
		newSceneJSON.URL = scene.URL.String
	}

	if scene.Date.Valid {
		newSceneJSON.Date = utils.GetYMDFromDatabaseDate(scene.Date.String)
	}

	if scene.Rating.Valid {
		newSceneJSON.Rating = int(scene.Rating.Int64)
	}

	newSceneJSON.Organized = scene.Organized
	newSceneJSON.OCounter = scene.OCounter

	if scene.Details.Valid {
		newSceneJSON.Details = scene.Details.String
	}

	newSceneJSON.File = getSceneFileJSON(scene)

	cover, err := reader.GetCover(ctx, scene.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene cover: %v", err)
	}

	if len(cover) > 0 {
		newSceneJSON.Cover = utils.GetBase64StringFromData(cover)
	}

	stashIDs, _ := reader.GetStashIDs(ctx, scene.ID)
	var ret []models.StashID
	for _, stashID := range stashIDs {
		newJoin := models.StashID{
			StashID:  stashID.StashID,
			Endpoint: stashID.Endpoint,
		}
		ret = append(ret, newJoin)
	}

	newSceneJSON.StashIDs = ret

	return &newSceneJSON, nil
}

func getSceneFileJSON(scene *models.Scene) *jsonschema.SceneFile {
	ret := &jsonschema.SceneFile{}

	if scene.FileModTime.Valid {
		ret.ModTime = json.JSONTime{Time: scene.FileModTime.Timestamp}
	}

	if scene.Size.Valid {
		ret.Size = scene.Size.String
	}

	if scene.Duration.Valid {
		ret.Duration = getDecimalString(scene.Duration.Float64)
	}

	if scene.VideoCodec.Valid {
		ret.VideoCodec = scene.VideoCodec.String
	}

	if scene.AudioCodec.Valid {
		ret.AudioCodec = scene.AudioCodec.String
	}

	if scene.Format.Valid {
		ret.Format = scene.Format.String
	}

	if scene.Width.Valid {
		ret.Width = int(scene.Width.Int64)
	}

	if scene.Height.Valid {
		ret.Height = int(scene.Height.Int64)
	}

	if scene.Framerate.Valid {
		ret.Framerate = getDecimalString(scene.Framerate.Float64)
	}

	if scene.Bitrate.Valid {
		ret.Bitrate = int(scene.Bitrate.Int64)
	}

	return ret
}

// GetStudioName returns the name of the provided scene's studio. It returns an
// empty string if there is no studio assigned to the scene.
func GetStudioName(ctx context.Context, reader studio.Finder, scene *models.Scene) (string, error) {
	if scene.StudioID.Valid {
		studio, err := reader.Find(ctx, int(scene.StudioID.Int64))
		if err != nil {
			return "", err
		}

		if studio != nil {
			return studio.Name.String, nil
		}
	}

	return "", nil
}

// GetTagNames returns a slice of tag names corresponding to the provided
// scene's tags.
func GetTagNames(ctx context.Context, reader TagFinder, scene *models.Scene) ([]string, error) {
	tags, err := reader.FindBySceneID(ctx, scene.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene tags: %v", err)
	}

	return getTagNames(tags), nil
}

func getTagNames(tags []*models.Tag) []string {
	var results []string
	for _, tag := range tags {
		if tag.Name != "" {
			results = append(results, tag.Name)
		}
	}

	return results
}

// GetDependentTagIDs returns a slice of unique tag IDs that this scene references.
func GetDependentTagIDs(ctx context.Context, tags MarkerTagFinder, markerReader MarkerFinder, scene *models.Scene) ([]int, error) {
	var ret []int

	t, err := tags.FindBySceneID(ctx, scene.ID)
	if err != nil {
		return nil, err
	}

	for _, tt := range t {
		ret = intslice.IntAppendUnique(ret, tt.ID)
	}

	sm, err := markerReader.FindBySceneID(ctx, scene.ID)
	if err != nil {
		return nil, err
	}

	for _, smm := range sm {
		ret = intslice.IntAppendUnique(ret, smm.PrimaryTagID)
		smmt, err := tags.FindBySceneMarkerID(ctx, smm.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid tags for scene marker: %v", err)
		}

		for _, smmtt := range smmt {
			ret = intslice.IntAppendUnique(ret, smmtt.ID)
		}
	}

	return ret, nil
}

type MovieFinder interface {
	Find(ctx context.Context, id int) (*models.Movie, error)
}

// GetSceneMoviesJSON returns a slice of SceneMovie JSON representation objects
// corresponding to the provided scene's scene movie relationships.
func GetSceneMoviesJSON(ctx context.Context, movieReader MovieFinder, sceneReader MovieGetter, scene *models.Scene) ([]jsonschema.SceneMovie, error) {
	sceneMovies, err := sceneReader.GetMovies(ctx, scene.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene movies: %v", err)
	}

	var results []jsonschema.SceneMovie
	for _, sceneMovie := range sceneMovies {
		movie, err := movieReader.Find(ctx, sceneMovie.MovieID)
		if err != nil {
			return nil, fmt.Errorf("error getting movie: %v", err)
		}

		if movie.Name.Valid {
			sceneMovieJSON := jsonschema.SceneMovie{
				MovieName:  movie.Name.String,
				SceneIndex: int(sceneMovie.SceneIndex.Int64),
			}
			results = append(results, sceneMovieJSON)
		}
	}

	return results, nil
}

// GetDependentMovieIDs returns a slice of movie IDs that this scene references.
func GetDependentMovieIDs(ctx context.Context, sceneReader MovieGetter, scene *models.Scene) ([]int, error) {
	var ret []int

	m, err := sceneReader.GetMovies(ctx, scene.ID)
	if err != nil {
		return nil, err
	}

	for _, mm := range m {
		ret = append(ret, mm.MovieID)
	}

	return ret, nil
}

// GetSceneMarkersJSON returns a slice of SceneMarker JSON representation
// objects corresponding to the provided scene's markers.
func GetSceneMarkersJSON(ctx context.Context, markerReader MarkerFinder, tagReader MarkerTagFinder, scene *models.Scene) ([]jsonschema.SceneMarker, error) {
	sceneMarkers, err := markerReader.FindBySceneID(ctx, scene.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene markers: %v", err)
	}

	var results []jsonschema.SceneMarker

	for _, sceneMarker := range sceneMarkers {
		primaryTag, err := tagReader.Find(ctx, sceneMarker.PrimaryTagID)
		if err != nil {
			return nil, fmt.Errorf("invalid primary tag for scene marker: %v", err)
		}

		sceneMarkerTags, err := tagReader.FindBySceneMarkerID(ctx, sceneMarker.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid tags for scene marker: %v", err)
		}

		sceneMarkerJSON := jsonschema.SceneMarker{
			Title:      sceneMarker.Title,
			Seconds:    getDecimalString(sceneMarker.Seconds),
			PrimaryTag: primaryTag.Name,
			Tags:       getTagNames(sceneMarkerTags),
			CreatedAt:  json.JSONTime{Time: sceneMarker.CreatedAt.Timestamp},
			UpdatedAt:  json.JSONTime{Time: sceneMarker.UpdatedAt.Timestamp},
		}

		results = append(results, sceneMarkerJSON)
	}

	return results, nil
}

func getDecimalString(num float64) string {
	if num == 0 {
		return ""
	}

	precision := getPrecision(num)
	if precision == 0 {
		precision = 1
	}
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", num)
}

func getPrecision(num float64) int {
	if num == 0 {
		return 0
	}

	e := 1.0
	p := 0
	for (math.Round(num*e) / e) != num {
		e *= 10
		p++
	}
	return p
}
