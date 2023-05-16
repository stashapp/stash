package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

func convertImageFile(f *file.ImageFile) *ImageFile {
	ret := &ImageFile{
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
		ret.ZipFileID = &zipFileID
	}

	return ret
}

func (r *imageResolver) getPrimaryFile(ctx context.Context, obj *models.Image) (file.VisualFile, error) {
	if obj.PrimaryFileID != nil {
		f, err := loaders.From(ctx).FileByID.Load(*obj.PrimaryFileID)
		if err != nil {
			return nil, err
		}

		asFrame, ok := f.(file.VisualFile)
		if !ok {
			return nil, fmt.Errorf("file %T is not an frame", f)
		}

		return asFrame, nil
	}

	return nil, nil
}

func (r *imageResolver) getFiles(ctx context.Context, obj *models.Image) ([]file.File, error) {
	fileIDs, err := loaders.From(ctx).ImageFiles.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	files, errs := loaders.From(ctx).FileByID.LoadAll(fileIDs)
	return files, firstError(errs)
}

func (r *imageResolver) Title(ctx context.Context, obj *models.Image) (*string, error) {
	ret := obj.GetTitle()
	return &ret, nil
}

func (r *imageResolver) File(ctx context.Context, obj *models.Image) (*ImageFileType, error) {
	f, err := r.getPrimaryFile(ctx, obj)
	if err != nil {
		return nil, err
	}

	if f == nil {
		return nil, nil
	}

	width := f.GetWidth()
	height := f.GetHeight()
	size := f.Base().Size
	return &ImageFileType{
		Size:   int(size),
		Width:  width,
		Height: height,
	}, nil
}

func convertVisualFile(f file.File) VisualFile {
	switch f := f.(type) {
	case *file.ImageFile:
		return convertImageFile(f)
	case *file.VideoFile:
		return convertVideoFile(f)
	default:
		panic(fmt.Sprintf("unknown file type %T", f))
	}
}

func (r *imageResolver) VisualFiles(ctx context.Context, obj *models.Image) ([]VisualFile, error) {
	fileIDs, err := loaders.From(ctx).ImageFiles.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	files, errs := loaders.From(ctx).FileByID.LoadAll(fileIDs)
	ret := make([]VisualFile, len(files))
	for i, f := range files {
		ret[i] = convertVisualFile(f)
	}

	return ret, firstError(errs)
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
		imageFile, ok := f.(*file.ImageFile)
		if !ok {
			continue
		}

		thisFile := convertImageFile(imageFile)

		ret = append(ret, thisFile)
	}

	return ret, nil
}

func (r *imageResolver) FileModTime(ctx context.Context, obj *models.Image) (*time.Time, error) {
	f, err := r.getPrimaryFile(ctx, obj)
	if err != nil {
		return nil, err
	}
	if f != nil {
		return &f.Base().ModTime, nil
	}

	return nil, nil
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

func (r *imageResolver) Rating(ctx context.Context, obj *models.Image) (*int, error) {
	if obj.Rating != nil {
		rating := models.Rating100To5(*obj.Rating)
		return &rating, nil
	}
	return nil, nil
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
