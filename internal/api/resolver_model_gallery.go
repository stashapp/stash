package api

import (
	"context"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

func (r *galleryResolver) Files(ctx context.Context, obj *models.Gallery) ([]*GalleryFile, error) {
	ret := make([]*GalleryFile, len(obj.Files))

	for i, f := range obj.Files {
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

	if err := r.withTxn(ctx, func(ctx context.Context) error {
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
	f := obj.PrimaryFile()
	if f != nil {
		return &f.Base().ModTime, nil
	}

	return nil, nil
}

func (r *galleryResolver) Images(ctx context.Context, obj *models.Gallery) (ret []*models.Image, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
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
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		// doing this via Query is really slow, so stick with FindByGalleryID
		imgs, err := r.repository.Image.FindByGalleryID(ctx, obj.ID)
		if err != nil {
			return err
		}

		if len(imgs) > 0 {
			ret = imgs[0]
		}

		for _, img := range imgs {
			if image.IsCover(img) {
				ret = img
				break
			}
		}

		return nil
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

func (r *galleryResolver) Scenes(ctx context.Context, obj *models.Gallery) (ret []*models.Scene, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Scene.FindByGalleryID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) Studio(ctx context.Context, obj *models.Gallery) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Studio.Find(ctx, *obj.StudioID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) Tags(ctx context.Context, obj *models.Gallery) (ret []*models.Tag, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Tag.FindByGalleryID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) Performers(ctx context.Context, obj *models.Gallery) (ret []*models.Performer, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Performer.FindByGalleryID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) ImageCount(ctx context.Context, obj *models.Gallery) (ret int, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Image.CountByGalleryID(ctx, obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}
