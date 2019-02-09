package api

import (
	"context"
	"github.com/stashapp/stash/internal/models"
	"strconv"
)

func (r *galleryResolver) ID(ctx context.Context, obj *models.Gallery) (string, error) {
	return strconv.Itoa(obj.ID), nil
}

func (r *galleryResolver) Title(ctx context.Context, obj *models.Gallery) (*string, error) {
	return nil, nil // TODO remove this from schema
}

func (r *galleryResolver) Files(ctx context.Context, obj *models.Gallery) ([]models.GalleryFilesType, error) {
	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	return obj.GetFiles(baseURL), nil
}
