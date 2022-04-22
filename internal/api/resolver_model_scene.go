package api

import (
	"context"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *sceneResolver) Date(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.Date != nil {
		result := obj.Date.String()
		return &result, nil
	}
	return nil, nil
}

func (r *sceneResolver) File(ctx context.Context, obj *models.Scene) (*models.SceneFileType, error) {
	var bitrate *int
	if obj.Bitrate != nil {
		v := int(*obj.Bitrate)
		bitrate = &v
	}

	return &models.SceneFileType{
		Size:       obj.Size,
		Duration:   handleFloat64Ptr(obj.Duration),
		VideoCodec: obj.VideoCodec,
		AudioCodec: obj.AudioCodec,
		Width:      obj.Width,
		Height:     obj.Height,
		Framerate:  handleFloat64Ptr(obj.Framerate),
		Bitrate:    bitrate,
	}, nil
}

func (r *sceneResolver) Paths(ctx context.Context, obj *models.Scene) (*ScenePathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	config := manager.GetInstance().Config
	builder := urlbuilders.NewSceneURLBuilder(baseURL, obj.ID)
	builder.APIKey = config.GetAPIKey()
	screenshotPath := builder.GetScreenshotURL(obj.UpdatedAt)
	previewPath := builder.GetStreamPreviewURL()
	streamPath := builder.GetStreamURL()
	webpPath := builder.GetStreamPreviewImageURL()
	vttPath := builder.GetSpriteVTTURL()
	spritePath := builder.GetSpriteURL()
	chaptersVttPath := builder.GetChaptersVTTURL()
	funscriptPath := builder.GetFunscriptURL()
	captionBasePath := builder.GetCaptionURL()
	interactiveHeatmap := builder.GetInteractiveHeatmapURL()

	return &ScenePathsType{
		Screenshot:         &screenshotPath,
		Preview:            &previewPath,
		Stream:             &streamPath,
		Webp:               &webpPath,
		Vtt:                &vttPath,
		ChaptersVtt:        &chaptersVttPath,
		Sprite:             &spritePath,
		Funscript:          &funscriptPath,
		InteractiveHeatmap: &interactiveHeatmap,
		Caption:            &captionBasePath,
	}, nil
}

func (r *sceneResolver) SceneMarkers(ctx context.Context, obj *models.Scene) (ret []*models.SceneMarker, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.SceneMarker.FindBySceneID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Captions(ctx context.Context, obj *models.Scene) (ret []*models.SceneCaption, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.GetCaptions(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *sceneResolver) Galleries(ctx context.Context, obj *models.Scene) (ret []*models.Gallery, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.FindBySceneID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Studio(ctx context.Context, obj *models.Scene) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.Find(ctx, *obj.StudioID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Movies(ctx context.Context, obj *models.Scene) (ret []*SceneMovie, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		mqb := r.repository.Movie

		for _, sm := range obj.Movies {
			movie, err := mqb.Find(ctx, sm.MovieID)
			if err != nil {
				return err
			}

			sceneIdx := sm.SceneIndex
			sceneMovie := &SceneMovie{
				Movie:      movie,
				SceneIndex: sceneIdx,
			}

			ret = append(ret, sceneMovie)
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *sceneResolver) Tags(ctx context.Context, obj *models.Scene) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.FindBySceneID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Performers(ctx context.Context, obj *models.Scene) (ret []*models.Performer, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.FindBySceneID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Phash(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.Phash != nil {
		hexval := utils.PhashToString(*obj.Phash)
		return &hexval, nil
	}
	return nil, nil
}

func (r *sceneResolver) SceneStreams(ctx context.Context, obj *models.Scene) ([]*manager.SceneStreamEndpoint, error) {
	config := manager.GetInstance().Config

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewSceneURLBuilder(baseURL, obj.ID)

	return manager.GetSceneStreamPaths(obj, builder.GetStreamURL(), config.GetMaxStreamingTranscodeSize())
}
