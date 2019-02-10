package urlbuilders

import "strconv"

type sceneURLBuilder struct {
	BaseURL string
	SceneID string
}

func NewSceneURLBuilder(baseURL string, sceneID int) sceneURLBuilder {
	return sceneURLBuilder{
		BaseURL: baseURL,
		SceneID: strconv.Itoa(sceneID),
	}
}

func (b sceneURLBuilder) GetStreamUrl() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/stream.mp4"
}

func (b sceneURLBuilder) GetStreamPreviewUrl() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/preview"
}

func (b sceneURLBuilder) GetStreamPreviewImageUrl() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/webp"
}

func (b sceneURLBuilder) GetSpriteVttUrl() string {
	return b.BaseURL + "/scene/" + b.SceneID + "_thumbs.vtt"
}

func (b sceneURLBuilder) GetScreenshotUrl() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/screenshot"
}

func (b sceneURLBuilder) GetChaptersVttUrl() string {
	return b.BaseURL + "/scene/" + b.SceneID + "/vtt/chapter"
}

func (b sceneURLBuilder) GetSceneMarkerStreamUrl(sceneMarkerId int) string {
	return b.BaseURL + "/scene/" + b.SceneID + "/scene_marker/" + strconv.Itoa(sceneMarkerId) + "/stream"
}

func (b sceneURLBuilder) GetSceneMarkerStreamPreviewUrl(sceneMarkerId int) string {
	return b.BaseURL + "/scene/" + b.SceneID + "/scene_marker/" + strconv.Itoa(sceneMarkerId) + "/preview"
}
