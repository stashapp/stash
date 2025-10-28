package api

import (
	"context"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *imageResolver) getFiles(ctx context.Context, obj *models.Image) ([]models.File, error) {
	fileIDs, err := loaders.From(ctx).ImageFiles.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	files, errs := loaders.From(ctx).FileByID.LoadAll(fileIDs)
	return files, firstError(errs)
}

func (r *imageResolver) VisualFiles(ctx context.Context, obj *models.Image) ([]VisualFile, error) {
	files, err := r.getFiles(ctx, obj)
	if err != nil {
		return nil, err
	}

	ret := make([]VisualFile, len(files))
	for i, f := range files {
		ret[i], err = convertVisualFile(f)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (r *imageResolver) Date(ctx context.Context, obj *models.Image) (*string, error) {
	if obj.Date != nil {
		result := obj.Date.String()
		return &result, nil
	}
	return nil, nil
}

func (r *imageResolver) Files(ctx context.Context, obj *models.Image) ([]*ImageFile, error) {
	files, err := r.getFiles(ctx, obj)
	if err != nil {
		return nil, err
	}

	var ret []*ImageFile

	for _, f := range files {
		// filter out non-image files
		imageFile, ok := f.(*models.ImageFile)
		if !ok {
			continue
		}

		ret = append(ret, &ImageFile{
			ImageFile: imageFile,
		})
	}

	return ret, nil
}

func (r *imageResolver) Paths(ctx context.Context, obj *models.Image) (*ImagePathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewImageURLBuilder(baseURL, obj)
	thumbnailPath := builder.GetThumbnailURL()
	previewPath := builder.GetPreviewURL()
	imagePath := builder.GetImageURL()
	return &ImagePathsType{
		Image:     &imagePath,
		Thumbnail: &thumbnailPath,
		Preview:   &previewPath,
	}, nil
}

func (r *imageResolver) Galleries(ctx context.Context, obj *models.Image) (ret []*models.Gallery, err error) {
	if !obj.GalleryIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadGalleryIDs(ctx, r.repository.Image)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).GalleryByID.LoadAll(obj.GalleryIDs.List())
	return ret, firstError(errs)
}

func (r *imageResolver) Rating100(ctx context.Context, obj *models.Image) (*int, error) {
	return obj.Rating, nil
}

func (r *imageResolver) Studio(ctx context.Context, obj *models.Image) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	return loaders.From(ctx).StudioByID.Load(*obj.StudioID)
}

func (r *imageResolver) Tags(ctx context.Context, obj *models.Image) (ret []*models.Tag, err error) {
	if !obj.TagIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadTagIDs(ctx, r.repository.Image)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.TagIDs.List())
	return ret, firstError(errs)
}

func (r *imageResolver) Performers(ctx context.Context, obj *models.Image) (ret []*models.Performer, err error) {
	if !obj.PerformerIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadPerformerIDs(ctx, r.repository.Image)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).PerformerByID.LoadAll(obj.PerformerIDs.List())
	return ret, firstError(errs)
}

func (r *imageResolver) URL(ctx context.Context, obj *models.Image) (*string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Image)
		}); err != nil {
			return nil, err
		}
	}

	urls := obj.URLs.List()
	if len(urls) == 0 {
		return nil, nil
	}

	return &urls[0], nil
}

func (r *imageResolver) Urls(ctx context.Context, obj *models.Image) ([]string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Image)
		}); err != nil {
			return nil, err
		}
	}

	return obj.URLs.List(), nil
}
