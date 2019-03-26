package api

import (
	"context"
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

func (r *queryResolver) FindGallery(ctx context.Context, id string) (*models.Gallery, error) {
	qb := models.NewGalleryQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt)
}

func (r *queryResolver) FindGalleries(ctx context.Context, filter *models.FindFilterType) (*models.FindGalleriesResultType, error) {
	qb := models.NewGalleryQueryBuilder()
	galleries, total := qb.Query(filter)
	return &models.FindGalleriesResultType{
		Count:     total,
		Galleries: galleries,
	}, nil
}
