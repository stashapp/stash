package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *galleryResolver) Title(ctx context.Context, obj *models.Gallery) (*string, error) {
	if obj.Title.Valid {
		return &obj.Title.String, nil
	}
	return nil, nil
}

func (r *galleryResolver) Files(ctx context.Context, obj *models.Gallery) ([]*models.GalleryFilesType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	return obj.GetFiles(baseURL), nil
}

func (r *galleryResolver) URL(ctx context.Context, obj *models.Gallery) (*string, error) {
	if obj.URL.Valid {
		return &obj.URL.String, nil
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
