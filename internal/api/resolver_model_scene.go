package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *sceneResolver) Checksum(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.Checksum.Valid {
		return &obj.Checksum.String, nil
	}
	return nil, nil
}

func (r *sceneResolver) Oshash(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.OSHash.Valid {
		return &obj.OSHash.String, nil
	}
	return nil, nil
}

func (r *sceneResolver) Title(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.Title.Valid {
		return &obj.Title.String, nil
	}
	return nil, nil
}

func (r *sceneResolver) Details(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.Details.Valid {
		return &obj.Details.String, nil
	}
	return nil, nil
}

func (r *sceneResolver) URL(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.URL.Valid {
		return &obj.URL.String, nil
	}
	return nil, nil
}

func (r *sceneResolver) Date(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.Date.Valid {
		result := utils.GetYMDFromDatabaseDate(obj.Date.String)
		return &result, nil
	}
	return nil, nil
}

func (r *sceneResolver) Rating(ctx context.Context, obj *models.Scene) (*int, error) {
	if obj.Rating.Valid {
		rating := int(obj.Rating.Int64)
		return &rating, nil
	}
	return nil, nil
}

func (r *sceneResolver) InteractiveSpeed(ctx context.Context, obj *models.Scene) (*int, error) {
	if obj.InteractiveSpeed.Valid {
		interactive_speed := int(obj.InteractiveSpeed.Int64)
		return &interactive_speed, nil
	}
	return nil, nil
}

func (r *sceneResolver) File(ctx context.Context, obj *models.Scene) (*models.SceneFileType, error) {
	width := int(obj.Width.Int64)
	height := int(obj.Height.Int64)
	bitrate := int(obj.Bitrate.Int64)
	return &models.SceneFileType{
		Size:       &obj.Size.String,
		Duration:   handleFloat64(obj.Duration.Float64),
		VideoCodec: &obj.VideoCodec.String,
		AudioCodec: &obj.AudioCodec.String,
		Width:      &width,
		Height:     &height,
		Framerate:  handleFloat64(obj.Framerate.Float64),
		Bitrate:    &bitrate,
	}, nil
}

func (r *sceneResolver) Paths(ctx context.Context, obj *models.Scene) (*ScenePathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	config := manager.GetInstance().Config
	builder := urlbuilders.NewSceneURLBuilder(baseURL, obj.ID)
	builder.APIKey = config.GetAPIKey()
	screenshotPath := builder.GetScreenshotURL(obj.UpdatedAt.Timestamp)
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
	if !obj.StudioID.Valid {
		return nil, nil
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.Find(ctx, int(obj.StudioID.Int64))
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Movies(ctx context.Context, obj *models.Scene) (ret []*SceneMovie, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Scene
		mqb := r.repository.Movie

		sceneMovies, err := qb.GetMovies(ctx, obj.ID)
		if err != nil {
			return err
		}

		for _, sm := range sceneMovies {
			movie, err := mqb.Find(ctx, sm.MovieID)
			if err != nil {
				return err
			}

			sceneIdx := sm.SceneIndex
			sceneMovie := &SceneMovie{
				Movie: movie,
			}

			if sceneIdx.Valid {
				idx := int(sceneIdx.Int64)
				sceneMovie.SceneIndex = &idx
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

func (r *sceneResolver) StashIds(ctx context.Context, obj *models.Scene) (ret []*models.StashID, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.GetStashIDs(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *sceneResolver) Phash(ctx context.Context, obj *models.Scene) (*string, error) {
	if obj.Phash.Valid {
		hexval := utils.PhashToString(obj.Phash.Int64)
		return &hexval, nil
	}
	return nil, nil
}

func (r *sceneResolver) CreatedAt(ctx context.Context, obj *models.Scene) (*time.Time, error) {
	return &obj.CreatedAt.Timestamp, nil
}

func (r *sceneResolver) UpdatedAt(ctx context.Context, obj *models.Scene) (*time.Time, error) {
	return &obj.UpdatedAt.Timestamp, nil
}

func (r *sceneResolver) FileModTime(ctx context.Context, obj *models.Scene) (*time.Time, error) {
	return &obj.FileModTime.Timestamp, nil
}

func (r *sceneResolver) SceneStreams(ctx context.Context, obj *models.Scene) ([]*manager.SceneStreamEndpoint, error) {
	config := manager.GetInstance().Config

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewSceneURLBuilder(baseURL, obj.ID)

	return manager.GetSceneStreamPaths(obj, builder.GetStreamURL(), config.GetMaxStreamingTranscodeSize())
}
