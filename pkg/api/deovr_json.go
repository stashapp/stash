package api

import (
	"context"
	"encoding/json"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

type MultipleVideoJsonResponse struct {
	Scenes []SceneLibrary `json:"scenes"`
}

type SceneLibrary struct {
	Name string         `json:"name"`
	List []SlimDeoScene `json:"list"`
}

type SlimDeoScene struct {
	Title        string `json:"title"`
	VideoLength  uint   `json:"videoLength"`
	ThumbnailURL string `json:"thumbnailUrl"`
	VideoJsonURL string `json:"video_url"`
	VideoPreview string `json:"videoPreview,omitempty"`
}

type FullDeoScene struct {
	Encodings    []DeoSceneEncoding `json:"encodings"`
	Title        string             `json:"title"`
	ID           uint               `json:"id"`
	VideoLength  uint               `json:"videoLength"`
	Is3D         bool               `json:"is3d"`
	ScreenType   string             `json:"screenType,omitempty"`
	StereoMode   string             `json:"stereoMode,omitempty"`
	VideoPreview string             `json:"videoPreview,omitempty"`
	ThumbnailURL string             `json:"thumbnailUrl"`
}

type DeoSceneEncoding struct {
	Name         string                `json:"name"` // This should be the video codec
	VideoSources []DeoSceneVideoSource `json:"videoSources"`
}

type DeoSceneVideoSource struct {
	Resolution uint   `json:"resolution"`
	URL        string `json:"url"`
}

func getEverySceneJSON(ctx context.Context) []byte {
	var err error
	txnManager := manager.GetInstance().TxnManager
	var scenes []*models.Scene
	err = txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		pageSize := -1
		scenes, err = scene.Query(r.Scene(), &models.SceneFilterType{}, &models.FindFilterType{
			PerPage: &pageSize,
		})
		if err != nil {
			logger.Errorf("Could not retrieve scene list: %s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return nil
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	var list []SlimDeoScene
	for _, sceneModel := range scenes {
		builder := urlbuilders.NewSceneURLBuilder(baseURL, sceneModel.ID)

		list = append(list, SlimDeoScene{
			Title:        sceneModel.GetTitle(),
			VideoLength:  uint(sceneModel.Duration.Float64),
			ThumbnailURL: builder.GetScreenshotURL(sceneModel.UpdatedAt.Timestamp),
			VideoPreview: builder.GetStreamPreviewURL(),
			VideoJsonURL: builder.GetDeoVRURL(false),
		})
	}

	library := SceneLibrary{
		Name: "Library",
		List: list,
	}
	libraries := []SceneLibrary{library}

	response := MultipleVideoJsonResponse{
		Scenes: libraries,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("Could not marshal JSON for deoVR all scenes: %s", err.Error())
	}
	return jsonBytes
}

func getSingleSceneJSON(ctx context.Context, sceneModel *models.Scene) []byte {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewSceneURLBuilder(baseURL, sceneModel.ID)

	videoSource := DeoSceneVideoSource{
		Resolution: uint(sceneModel.Height.Int64),
		URL:        builder.GetStreamURL(),
	}

	encoding := DeoSceneEncoding{
		Name: sceneModel.VideoCodec.String,
		VideoSources: []DeoSceneVideoSource{
			videoSource,
		},
	}

	sceneStruct := FullDeoScene{
		Encodings: []DeoSceneEncoding{
			encoding,
		},
		Title:        sceneModel.GetTitle(),
		ID:           uint(sceneModel.ID),
		VideoLength:  uint(sceneModel.Duration.Float64),
		Is3D:         true,
		VideoPreview: builder.GetStreamPreviewURL(),
		ThumbnailURL: builder.GetScreenshotURL(sceneModel.UpdatedAt.Timestamp),
	}

	jsonBytes, err := json.Marshal(sceneStruct)
	if err != nil {
		logger.Errorf("Could not marshal JSON for single deoVR scene: %s", err.Error())
	}
	return jsonBytes
}
