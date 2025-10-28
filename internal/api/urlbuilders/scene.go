package urlbuilders

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

type SceneURLBuilder struct {
	BaseURL   string
	SceneID   string
	UpdatedAt string
}

func NewSceneURLBuilder(baseURL string, scene *models.Scene) SceneURLBuilder {
	return SceneURLBuilder{
		BaseURL:   baseURL,
		SceneID:   strconv.Itoa(scene.ID),
		UpdatedAt: strconv.FormatInt(scene.UpdatedAt.Unix(), 10),
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

func (b SceneURLBuilder) GetScreenshotURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/screenshot?t=" + b.UpdatedAt
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
