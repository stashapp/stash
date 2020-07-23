package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

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

func (r *sceneResolver) File(ctx context.Context, obj *models.Scene) (*models.SceneFileType, error) {
	width := int(obj.Width.Int64)
	height := int(obj.Height.Int64)
	bitrate := int(obj.Bitrate.Int64)
	return &models.SceneFileType{
		Size:       &obj.Size.String,
		Duration:   &obj.Duration.Float64,
		VideoCodec: &obj.VideoCodec.String,
		AudioCodec: &obj.AudioCodec.String,
		Width:      &width,
		Height:     &height,
		Framerate:  &obj.Framerate.Float64,
		Bitrate:    &bitrate,
	}, nil
}

func (r *sceneResolver) Paths(ctx context.Context, obj *models.Scene) (*models.ScenePathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewSceneURLBuilder(baseURL, obj.ID)
	screenshotPath := builder.GetScreenshotURL(obj.UpdatedAt.Timestamp)
	previewPath := builder.GetStreamPreviewURL()
	streamPath := builder.GetStreamURL()
	webpPath := builder.GetStreamPreviewImageURL()
	vttPath := builder.GetSpriteVTTURL()
	chaptersVttPath := builder.GetChaptersVTTURL()
	return &models.ScenePathsType{
		Screenshot:  &screenshotPath,
		Preview:     &previewPath,
		Stream:      &streamPath,
		Webp:        &webpPath,
		Vtt:         &vttPath,
		ChaptersVtt: &chaptersVttPath,
	}, nil
}

func (r *sceneResolver) SceneMarkers(ctx context.Context, obj *models.Scene) ([]*models.SceneMarker, error) {
	qb := models.NewSceneMarkerQueryBuilder()
	return qb.FindBySceneID(obj.ID, nil)
}

func (r *sceneResolver) Gallery(ctx context.Context, obj *models.Scene) (*models.Gallery, error) {
	qb := models.NewGalleryQueryBuilder()
	return qb.FindBySceneID(obj.ID, nil)
}

func (r *sceneResolver) Studio(ctx context.Context, obj *models.Scene) (*models.Studio, error) {
	qb := models.NewStudioQueryBuilder()
	return qb.FindBySceneID(obj.ID)
}

func (r *sceneResolver) Movies(ctx context.Context, obj *models.Scene) ([]*models.SceneMovie, error) {
	joinQB := models.NewJoinsQueryBuilder()
	qb := models.NewMovieQueryBuilder()

	sceneMovies, err := joinQB.GetSceneMovies(obj.ID, nil)
	if err != nil {
		return nil, err
	}

	var ret []*models.SceneMovie
	for _, sm := range sceneMovies {
		movie, err := qb.Find(sm.MovieID, nil)
		if err != nil {
			return nil, err
		}

		sceneIdx := sm.SceneIndex
		sceneMovie := &models.SceneMovie{
			Movie: movie,
		}

		if sceneIdx.Valid {
			var idx int
			idx = int(sceneIdx.Int64)
			sceneMovie.SceneIndex = &idx
		}

		ret = append(ret, sceneMovie)
	}
	return ret, nil
}

func (r *sceneResolver) Tags(ctx context.Context, obj *models.Scene) ([]*models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	return qb.FindBySceneID(obj.ID, nil)
}

func (r *sceneResolver) Performers(ctx context.Context, obj *models.Scene) ([]*models.Performer, error) {
	qb := models.NewPerformerQueryBuilder()
	return qb.FindBySceneID(obj.ID, nil)
}
