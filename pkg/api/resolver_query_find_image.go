package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func (r *queryResolver) FindImage(ctx context.Context, id *string, checksum *string) (*models.Image, error) {
	qb := sqlite.NewImageQueryBuilder()
	var image *models.Image
	var err error
	if id != nil {
		idInt, _ := strconv.Atoi(*id)
		image, err = qb.Find(idInt)
	} else if checksum != nil {
		image, err = qb.FindByChecksum(*checksum)
	}
	return image, err
}

func (r *queryResolver) FindImages(ctx context.Context, imageFilter *models.ImageFilterType, imageIds []int, filter *models.FindFilterType) (*models.FindImagesResultType, error) {
	qb := sqlite.NewImageQueryBuilder()
	images, total := qb.Query(imageFilter, filter)
	return &models.FindImagesResultType{
		Count:  total,
		Images: images,
	}, nil
}
