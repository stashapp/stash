package api

import (
	"context"

	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
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

func (r *imageResolver) Galleries(ctx context.Context, obj *models.Image) (ret []*models.Gallery, err error) {
	if err := r.withTxn(ctx, func(r models.Repository) error {
		var err error
		ret, err = r.Gallery().FindByImageID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *imageResolver) Studio(ctx context.Context, obj *models.Image) (ret *models.Studio, err error) {
	if !obj.StudioID.Valid {
		return nil, nil
	}

	if err := r.withTxn(ctx, func(r models.Repository) error {
		ret, err = r.Studio().Find(int(obj.StudioID.Int64))
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *imageResolver) Tags(ctx context.Context, obj *models.Image) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(r models.Repository) error {
		ret, err = r.Tag().FindByImageID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *imageResolver) Performers(ctx context.Context, obj *models.Image) (ret []*models.Performer, err error) {
	if err := r.withTxn(ctx, func(r models.Repository) error {
		ret, err = r.Performer().FindByImageID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
