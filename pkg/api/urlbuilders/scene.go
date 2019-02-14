package urlbuilders

import "strconv"

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

func (b SceneURLBuilder) GetStreamURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/stream.mp4"
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

func (b SceneURLBuilder) GetScreenshotURL() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/screenshot"
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
