package urlbuilders

import (
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type SceneMarkerURLBuilder struct {
	BaseURL  string
	SceneID  string
	MarkerID string
}

func NewSceneMarkerURLBuilder(baseURL string, sceneMarker *models.SceneMarker) SceneMarkerURLBuilder {
	return SceneMarkerURLBuilder{
		BaseURL:  baseURL,
		SceneID:  strconv.Itoa(sceneMarker.SceneID),
		MarkerID: strconv.Itoa(sceneMarker.ID),
	}
}

func (b SceneMarkerURLBuilder) GetStreamURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/scene_marker/" + b.MarkerID + "/stream"
}

func (b SceneMarkerURLBuilder) GetPreviewURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/scene_marker/" + b.MarkerID + "/preview"
}

func (b SceneMarkerURLBuilder) GetScreenshotURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/scene_marker/" + b.MarkerID + "/screenshot"
}
