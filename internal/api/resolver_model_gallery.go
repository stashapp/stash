package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *galleryResolver) Path(ctx context.Context, obj *models.Gallery) (*string, error) {
	if obj.Path.Valid {
		return &obj.Path.String, nil
	}
	return nil, nil
}

func (r *galleryResolver) Title(ctx context.Context, obj *models.Gallery) (*string, error) {
	if obj.Title.Valid {
		return &obj.Title.String, nil
	}
	return nil, nil
}

func (r *galleryResolver) Images(ctx context.Context, obj *models.Gallery) (ret []*models.Image, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var err error

		// #2376 - sort images by path
		// doing this via Query is really slow, so stick with FindByGalleryID
		ret, err = repo.Image().FindByGalleryID(obj.ID)
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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		// doing this via Query is really slow, so stick with FindByGalleryID
		imgs, err := repo.Image().FindByGalleryID(obj.ID)
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
	if obj.Date.Valid {
		result := utils.GetYMDFromDatabaseDate(obj.Date.String)
		return &result, nil
	}
	return nil, nil
}

func (r *galleryResolver) URL(ctx context.Context, obj *models.Gallery) (*string, error) {
	if obj.URL.Valid {
		return &obj.URL.String, nil
	}
	return nil, nil
}

func (r *galleryResolver) Details(ctx context.Context, obj *models.Gallery) (*string, error) {
	if obj.Details.Valid {
		return &obj.Details.String, nil
	}
	return nil, nil
}

func (r *galleryResolver) Rating(ctx context.Context, obj *models.Gallery) (*int, error) {
	if obj.Rating.Valid {
		rating := int(obj.Rating.Int64)
		return &rating, nil
	}
	return nil, nil
}

func (r *galleryResolver) Scenes(ctx context.Context, obj *models.Gallery) (ret []*models.Scene, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var err error
		ret, err = repo.Scene().FindByGalleryID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) Studio(ctx context.Context, obj *models.Gallery) (ret *models.Studio, err error) {
	if !obj.StudioID.Valid {
		return nil, nil
	}

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var err error
		ret, err = repo.Studio().Find(int(obj.StudioID.Int64))
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) Tags(ctx context.Context, obj *models.Gallery) (ret []*models.Tag, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var err error
		ret, err = repo.Tag().FindByGalleryID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) Performers(ctx context.Context, obj *models.Gallery) (ret []*models.Performer, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var err error
		ret, err = repo.Performer().FindByGalleryID(obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *galleryResolver) ImageCount(ctx context.Context, obj *models.Gallery) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		var err error
		ret, err = repo.Image().CountByGalleryID(obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *galleryResolver) CreatedAt(ctx context.Context, obj *models.Gallery) (*time.Time, error) {
	return &obj.CreatedAt.Timestamp, nil
}

func (r *galleryResolver) UpdatedAt(ctx context.Context, obj *models.Gallery) (*time.Time, error) {
	return &obj.UpdatedAt.Timestamp, nil
}

func (r *galleryResolver) FileModTime(ctx context.Context, obj *models.Gallery) (*time.Time, error) {
	return &obj.FileModTime.Timestamp, nil
}
