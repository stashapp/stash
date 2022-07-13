package api

import (
	"context"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *imageResolver) Title(ctx context.Context, obj *models.Image) (*string, error) {
	ret := obj.GetTitle()
	return &ret, nil
}

func (r *imageResolver) File(ctx context.Context, obj *models.Image) (*ImageFileType, error) {
	f := obj.PrimaryFile()
	width := f.Width
	height := f.Height
	size := f.Size
	return &ImageFileType{
		Size:   int(size),
		Width:  width,
		Height: height,
	}, nil
}

func (r *imageResolver) Files(ctx context.Context, obj *models.Image) ([]*ImageFile, error) {
	ret := make([]*ImageFile, len(obj.Files))

	for i, f := range obj.Files {
		ret[i] = &ImageFile{
			ID:             strconv.Itoa(int(f.ID)),
			Path:           f.Path,
			Basename:       f.Basename,
			ParentFolderID: strconv.Itoa(int(f.ParentFolderID)),
			ModTime:        f.ModTime,
			Size:           f.Size,
			Width:          f.Width,
			Height:         f.Height,
			CreatedAt:      f.CreatedAt,
			UpdatedAt:      f.UpdatedAt,
			Fingerprints:   resolveFingerprints(f.Base()),
		}

		if f.ZipFileID != nil {
			zipFileID := strconv.Itoa(int(*f.ZipFileID))
			ret[i].ZipFileID = &zipFileID
		}
	}

	return ret, nil
}

func (r *imageResolver) FileModTime(ctx context.Context, obj *models.Image) (*time.Time, error) {
	f := obj.PrimaryFile()
	if f != nil {
		return &f.ModTime, nil
	}

	return nil, nil
}

func (r *imageResolver) Paths(ctx context.Context, obj *models.Image) (*ImagePathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewImageURLBuilder(baseURL, obj)
	thumbnailPath := builder.GetThumbnailURL()
	imagePath := builder.GetImageURL()
	return &ImagePathsType{
		Image:     &imagePath,
		Thumbnail: &thumbnailPath,
	}, nil
}

func (r *imageResolver) Galleries(ctx context.Context, obj *models.Image) (ret []*models.Gallery, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Gallery.FindMany(ctx, obj.GalleryIDs)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *imageResolver) Studio(ctx context.Context, obj *models.Image) (ret *models.Studio, err error) {
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

func (r *imageResolver) Tags(ctx context.Context, obj *models.Image) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.FindMany(ctx, obj.TagIDs)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *imageResolver) Performers(ctx context.Context, obj *models.Image) (ret []*models.Performer, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.FindMany(ctx, obj.PerformerIDs)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
