package api

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/internal/manager/config"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

func (r *galleryResolver) getFiles(ctx context.Context, obj *models.Gallery) ([]models.File, error) {
	fileIDs, err := loaders.From(ctx).GalleryFiles.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	files, errs := loaders.From(ctx).FileByID.LoadAll(fileIDs)
	return files, firstError(errs)
}

func (r *galleryResolver) Files(ctx context.Context, obj *models.Gallery) ([]*GalleryFile, error) {
	files, err := r.getFiles(ctx, obj)
	if err != nil {
		return nil, err
	}

	ret := make([]*GalleryFile, len(files))

	for i, f := range files {
		ret[i] = &GalleryFile{
			BaseFile: f.Base(),
		}
	}

	return ret, nil
}

func (r *galleryResolver) Folder(ctx context.Context, obj *models.Gallery) (*models.Folder, error) {
	if obj.FolderID == nil {
		return nil, nil
	}

	var ret *models.Folder

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error

		ret, err = r.repository.Folder.Find(ctx, *obj.FolderID)
		if err != nil {
			return err
		}

		return err
	}); err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, nil
	}

	return ret, nil
}

func (r *galleryResolver) Cover(ctx context.Context, obj *models.Gallery) (ret *models.Image, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		// Find cover image first
		ret, err = image.FindGalleryCover(ctx, r.repository.Image, obj.ID, config.GetInstance().GetGalleryCoverRegex())
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) Date(ctx context.Context, obj *models.Gallery) (*string, error) {
	if obj.Date != nil {
		result := obj.Date.String()
		return &result, nil
	}
	return nil, nil
}

func (r *galleryResolver) Rating100(ctx context.Context, obj *models.Gallery) (*int, error) {
	return obj.Rating, nil
}

func (r *galleryResolver) Scenes(ctx context.Context, obj *models.Gallery) (ret []*models.Scene, err error) {
	if !obj.SceneIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadSceneIDs(ctx, r.repository.Gallery)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).SceneByID.LoadAll(obj.SceneIDs.List())
	return ret, firstError(errs)
}

func (r *galleryResolver) Studio(ctx context.Context, obj *models.Gallery) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	return loaders.From(ctx).StudioByID.Load(*obj.StudioID)
}

func (r *galleryResolver) Tags(ctx context.Context, obj *models.Gallery) (ret []*models.Tag, err error) {
	if !obj.TagIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadTagIDs(ctx, r.repository.Gallery)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.TagIDs.List())
	return ret, firstError(errs)
}

func (r *galleryResolver) Performers(ctx context.Context, obj *models.Gallery) (ret []*models.Performer, err error) {
	if !obj.PerformerIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadPerformerIDs(ctx, r.repository.Gallery)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).PerformerByID.LoadAll(obj.PerformerIDs.List())
	return ret, firstError(errs)
}

func (r *galleryResolver) ImageCount(ctx context.Context, obj *models.Gallery) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Image.CountByGalleryID(ctx, obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *galleryResolver) Chapters(ctx context.Context, obj *models.Gallery) (ret []*models.GalleryChapter, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.GalleryChapter.FindByGalleryID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) URL(ctx context.Context, obj *models.Gallery) (*string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Gallery)
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

func (r *galleryResolver) Urls(ctx context.Context, obj *models.Gallery) ([]string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Gallery)
		}); err != nil {
			return nil, err
		}
	}

	return obj.URLs.List(), nil
}

func (r *galleryResolver) Paths(ctx context.Context, obj *models.Gallery) (*GalleryPathsType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewGalleryURLBuilder(baseURL, obj)

	return &GalleryPathsType{
		Cover:   builder.GetCoverURL(),
		Preview: builder.GetPreviewURL(),
	}, nil
}

func (r *galleryResolver) Image(ctx context.Context, obj *models.Gallery, index int) (ret *models.Image, err error) {
	if index < 0 {
		return nil, fmt.Errorf("index must >= 0")
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Image.FindByGalleryIDIndex(ctx, obj.ID, uint(index))
		return err
	}); err != nil {
		return nil, err
	}

	return
}
