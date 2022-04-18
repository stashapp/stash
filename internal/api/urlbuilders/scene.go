package urlbuilders

import (
	"fmt"
	"strconv"
	"time"
)

type SceneURLBuilder struct {
	BaseURL string
	SceneID string
	APIKey  string
}

func NewSceneURLBuilder(baseURL string, sceneID int) SceneURLBuilder {
	return SceneURLBuilder{
		BaseURL: baseURL,
		SceneID: strconv.Itoa(sceneID),
	}
}

func (b SceneURLBuilder) GetStreamURL() string {
	var apiKeyParam string
	if b.APIKey != "" {
		apiKeyParam = fmt.Sprintf("?apikey=%s", b.APIKey)
	}
	return fmt.Sprintf("%s/scene/%s/stream%s", b.BaseURL, b.SceneID, apiKeyParam)
}

func (b SceneURLBuilder) GetStreamPreviewURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/preview"
}

func (b SceneURLBuilder) GetStreamPreviewImageURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/webp"
}

func (b SceneURLBuilder) GetSpriteVTTURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "_thumbs.vtt"
}

func (b SceneURLBuilder) GetSpriteURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "_sprite.jpg"
}

func (b SceneURLBuilder) GetScreenshotURL(updateTime time.Time) string {
	return b.BaseURL + "/scene/" + b.SceneID + "/screenshot?" + strconv.FormatInt(updateTime.Unix(), 10)
}

func (b SceneURLBuilder) GetChaptersVTTURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/vtt/chapter"
}

func (b SceneURLBuilder) GetSceneMarkerStreamURL(sceneMarkerID int) string {
	return b.BaseURL + "/scene/" + b.SceneID + "/scene_marker/" + strconv.Itoa(sceneMarkerID) + "/stream"
}

func (b SceneURLBuilder) GetSceneMarkerStreamPreviewURL(sceneMarkerID int) string {
	return b.BaseURL + "/scene/" + b.SceneID + "/scene_marker/" + strconv.Itoa(sceneMarkerID) + "/preview"
}

func (b SceneURLBuilder) GetSceneMarkerStreamScreenshotURL(sceneMarkerID int) string {
	return b.BaseURL + "/scene/" + b.SceneID + "/scene_marker/" + strconv.Itoa(sceneMarkerID) + "/screenshot"
}

func (b SceneURLBuilder) GetFunscriptURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/funscript"
}

func (b SceneURLBuilder) GetCaptionURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/caption"
}

func (b SceneURLBuilder) GetInteractiveHeatmapURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/interactive_heatmap"
}
