package api

import (
	"context"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/manager/config"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

func (r *galleryResolver) getPrimaryFile(ctx context.Context, obj *models.Gallery) (file.File, error) {
	if obj.PrimaryFileID != nil {
		f, err := loaders.From(ctx).FileByID.Load(*obj.PrimaryFileID)
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	return nil, nil
}

func (r *galleryResolver) getFiles(ctx context.Context, obj *models.Gallery) ([]file.File, error) {
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
		base := f.Base()
		ret[i] = &GalleryFile{
			ID:             strconv.Itoa(int(base.ID)),
			Path:           base.Path,
			Basename:       base.Basename,
			ParentFolderID: strconv.Itoa(int(base.ParentFolderID)),
			ModTime:        base.ModTime,
			Size:           base.Size,
			CreatedAt:      base.CreatedAt,
			UpdatedAt:      base.UpdatedAt,
			Fingerprints:   resolveFingerprints(base),
		}

		if base.ZipFileID != nil {
			zipFileID := strconv.Itoa(int(*base.ZipFileID))
			ret[i].ZipFileID = &zipFileID
		}
	}

	return ret, nil
}

func (r *galleryResolver) Folder(ctx context.Context, obj *models.Gallery) (*Folder, error) {
	if obj.FolderID == nil {
		return nil, nil
	}

	var ret *file.Folder

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

	rr := &Folder{
		ID:        ret.ID.String(),
		Path:      ret.Path,
		ModTime:   ret.ModTime,
		CreatedAt: ret.CreatedAt,
		UpdatedAt: ret.UpdatedAt,
	}

	if ret.ParentFolderID != nil {
		pfidStr := ret.ParentFolderID.String()
		rr.ParentFolderID = &pfidStr
	}

	if ret.ZipFileID != nil {
		zfidStr := ret.ZipFileID.String()
		rr.ZipFileID = &zfidStr
	}

	return rr, nil
}

func (r *galleryResolver) FileModTime(ctx context.Context, obj *models.Gallery) (*time.Time, error) {
	f, err := r.getPrimaryFile(ctx, obj)
	if err != nil {
		return nil, err
	}
	if f != nil {
		return &f.Base().ModTime, nil
	}

	return nil, nil
}

// Images is deprecated, slow and shouldn't be used
func (r *galleryResolver) Images(ctx context.Context, obj *models.Gallery) (ret []*models.Image, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error

		// #2376 - sort images by path
		// doing this via Query is really slow, so stick with FindByGalleryID
		ret, err = r.repository.Image.FindByGalleryID(ctx, obj.ID)
		if err != nil {
			return err
		}

		return err
	}); err != nil {
		return nil, err
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

func (r *galleryResolver) Checksum(ctx context.Context, obj *models.Gallery) (string, error) {
	if !obj.Files.PrimaryLoaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadPrimaryFile(ctx, r.repository.File)
		}); err != nil {
			return "", err
		}
	}

	return obj.PrimaryChecksum(), nil
}

func (r *galleryResolver) Rating(ctx context.Context, obj *models.Gallery) (*int, error) {
	if obj.Rating != nil {
		rating := models.Rating100To5(*obj.Rating)
		return &rating, nil
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
