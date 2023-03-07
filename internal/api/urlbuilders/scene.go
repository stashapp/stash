package urlbuilders

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type SceneURLBuilder struct {
	BaseURL string
	SceneID string
}

func NewSceneURLBuilder(baseURL string, sceneID int) SceneURLBuilder {
	return SceneURLBuilder{
		BaseURL: baseURL,
		SceneID: strconv.Itoa(sceneID),
	}
}

func (b SceneURLBuilder) GetStreamURL(apiKey string) *url.URL {
	u, err := url.Parse(fmt.Sprintf("%s/scene/%s/stream", b.BaseURL, b.SceneID))
	if err != nil {
		// shouldn't happen
		panic(err)
	}

	if apiKey != "" {
		v := u.Query()
		v.Set("apikey", apiKey)
		u.RawQuery = v.Encode()
	}
	return u
}

func (b SceneURLBuilder) GetStreamPreviewURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/preview"
}

func (b SceneURLBuilder) GetStreamPreviewImageURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/webp"
}

func (b SceneURLBuilder) GetSpriteVTTURL(checksum string) string {
	return b.BaseURL + "/scene/" + checksum + "_thumbs.vtt"
}

func (b SceneURLBuilder) GetSpriteURL(checksum string) string {
	return b.BaseURL + "/scene/" + checksum + "_sprite.jpg"
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
