package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func (r *queryResolver) FindGallery(ctx context.Context, id string) (*models.Gallery, error) {
	qb := sqlite.NewGalleryQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt, nil)
}

func (r *queryResolver) FindGalleries(ctx context.Context, galleryFilter *models.GalleryFilterType, filter *models.FindFilterType) (*models.FindGalleriesResultType, error) {
	qb := sqlite.NewGalleryQueryBuilder()
	galleries, total := qb.Query(galleryFilter, filter)
	return &models.FindGalleriesResultType{
		Count:     total,
		Galleries: galleries,
	}, nil
}
