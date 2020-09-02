package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindGallery(ctx context.Context, id string) (*models.Gallery, error) {
	qb := models.NewGalleryQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt)
}

func (r *queryResolver) FindGalleries(ctx context.Context, galleryFilter *models.GalleryFilterType, filter *models.FindFilterType) (*models.FindGalleriesResultType, error) {
	qb := models.NewGalleryQueryBuilder()
	galleries, total := qb.Query(galleryFilter, filter)
	return &models.FindGalleriesResultType{
		Count:     total,
		Galleries: galleries,
	}, nil
}
