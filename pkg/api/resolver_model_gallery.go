package api

import (
	"context"

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

func (r *galleryResolver) Images(ctx context.Context, obj *models.Gallery) ([]*models.Image, error) {
	qb := models.NewImageQueryBuilder()

	return qb.FindByGalleryID(obj.ID)
}

func (r *galleryResolver) Cover(ctx context.Context, obj *models.Gallery) (*models.Image, error) {
	qb := models.NewImageQueryBuilder()

	imgs, err := qb.FindByGalleryID(obj.ID)
	if err != nil {
		return nil, err
	}

	var ret *models.Image
	if len(imgs) > 0 {
		ret = imgs[0]
	}

	for _, img := range imgs {
		if image.IsCover(img) {
			ret = img
			break
		}
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

func (r *galleryResolver) Scene(ctx context.Context, obj *models.Gallery) (*models.Scene, error) {
	if !obj.SceneID.Valid {
		return nil, nil
	}

	qb := models.NewSceneQueryBuilder()
	return qb.Find(int(obj.SceneID.Int64))
}

func (r *galleryResolver) Studio(ctx context.Context, obj *models.Gallery) (*models.Studio, error) {
	if !obj.StudioID.Valid {
		return nil, nil
	}

	qb := models.NewStudioQueryBuilder()
	return qb.Find(int(obj.StudioID.Int64), nil)
}

func (r *galleryResolver) Tags(ctx context.Context, obj *models.Gallery) ([]*models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	return qb.FindByGalleryID(obj.ID, nil)
}

func (r *galleryResolver) Performers(ctx context.Context, obj *models.Gallery) ([]*models.Performer, error) {
	qb := models.NewPerformerQueryBuilder()
	return qb.FindByGalleryID(obj.ID, nil)
}

func (r *galleryResolver) ImageCount(ctx context.Context, obj *models.Gallery) (int, error) {
	qb := models.NewImageQueryBuilder()
	return qb.CountByGalleryID(obj.ID)
}
