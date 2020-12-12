package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func (r *imageResolver) Title(ctx context.Context, obj *models.Image) (*string, error) {
	ret := image.GetTitle(obj)
	return &ret, nil
}

func (r *imageResolver) Rating(ctx context.Context, obj *models.Image) (*int, error) {
	if obj.Rating.Valid {
		rating := int(obj.Rating.Int64)
		return &rating, nil
	}
	return nil, nil
}

func (r *imageResolver) File(ctx context.Context, obj *models.Image) (*models.ImageFileType, error) {
	width := int(obj.Width.Int64)
	height := int(obj.Height.Int64)
	size := int(obj.Size.Int64)
	return &models.ImageFileType{
		Size:   &size,
		Width:  &width,
		Height: &height,
	}, nil
}

func (r *imageResolver) Paths(ctx context.Context, obj *models.Image) (*models.ImagePathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewImageURLBuilder(baseURL, obj.ID)
	thumbnailPath := builder.GetThumbnailURL()
	imagePath := builder.GetImageURL()
	return &models.ImagePathsType{
		Image:     &imagePath,
		Thumbnail: &thumbnailPath,
	}, nil
}

func (r *imageResolver) Galleries(ctx context.Context, obj *models.Image) ([]*models.Gallery, error) {
	qb := sqlite.NewGalleryQueryBuilder()
	return qb.FindByImageID(obj.ID, nil)
}

func (r *imageResolver) Studio(ctx context.Context, obj *models.Image) (*models.Studio, error) {
	if !obj.StudioID.Valid {
		return nil, nil
	}

	qb := sqlite.NewStudioQueryBuilder()
	return qb.Find(int(obj.StudioID.Int64), nil)
}

func (r *imageResolver) Tags(ctx context.Context, obj *models.Image) ([]*models.Tag, error) {
	qb := sqlite.NewTagQueryBuilder()
	return qb.FindByImageID(obj.ID, nil)
}

func (r *imageResolver) Performers(ctx context.Context, obj *models.Image) ([]*models.Performer, error) {
	qb := sqlite.NewPerformerQueryBuilder()
	return qb.FindByImageID(obj.ID, nil)
}
