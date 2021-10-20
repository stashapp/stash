package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindImage(ctx context.Context, id *string, checksum *string) (*models.Image, error) {
	var image *models.Image

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		qb := repo.Image()
		var err error

		if id != nil {
			idInt, err := strconv.Atoi(*id)
			if err != nil {
				return err
			}

			image, err = qb.Find(idInt)
			if err != nil {
				return err
			}
		} else if checksum != nil {
			image, err = qb.FindByChecksum(*checksum)
		}

		return err
	}); err != nil {
		return nil, err
	}

	return image, nil
}

func (r *queryResolver) FindImages(ctx context.Context, imageFilter *models.ImageFilterType, imageIds []int, filter *models.FindFilterType) (ret *models.FindImagesResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		qb := repo.Image()
		images, total, megapixels, filesize, err := qb.Query(imageFilter, filter)
		if err != nil {
			return err
		}

		ret = &models.FindImagesResultType{
			Count:      total,
			Images:     images,
			Megapixels: megapixels,
			Filesize:   filesize,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
