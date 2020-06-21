package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *galleryResolver) Title(ctx context.Context, obj *models.Gallery) (*string, error) {
	return nil, nil // TODO remove this from schema
}

func (r *galleryResolver) Files(ctx context.Context, obj *models.Gallery) ([]*models.GalleryFilesType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	return obj.GetFiles(baseURL), nil
}

func (r *galleryResolver) Scene(ctx context.Context, obj *models.Gallery) (*models.Scene, error) {
	if !obj.SceneID.Valid {
		return nil, nil
	}

	qb := models.NewSceneQueryBuilder()
	return qb.Find(int(obj.SceneID.Int64))
}
