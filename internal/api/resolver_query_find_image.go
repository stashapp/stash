package api

import (
	"context"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

func (r *queryResolver) FindImage(ctx context.Context, id *string, checksum *string) (*models.Image, error) {
	var image *models.Image

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image
		var err error

		if id != nil {
			idInt, err := strconv.Atoi(*id)
			if err != nil {
				return err
			}

			image, err = qb.Find(ctx, idInt)
			if err != nil {
				return err
			}
		} else if checksum != nil {
			var images []*models.Image
			images, err = qb.FindByChecksum(ctx, *checksum)
			if err != nil {
				return err
			}

			if len(images) > 0 {
				image = images[0]
			}
		}

		return err
	}); err != nil {
		return nil, err
	}

	return image, nil
}

func (r *queryResolver) FindImages(ctx context.Context, imageFilter *models.ImageFilterType, imageIds []int, filter *models.FindFilterType) (ret *FindImagesResultType, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image

		fields := graphql.CollectAllFields(ctx)

		result, err := qb.Query(ctx, models.ImageQueryOptions{
			QueryOptions: models.QueryOptions{
				FindFilter: filter,
				Count:      sliceutil.Contains(fields, "count"),
			},
			ImageFilter: imageFilter,
			Megapixels:  sliceutil.Contains(fields, "megapixels"),
			TotalSize:   sliceutil.Contains(fields, "filesize"),
		})
		if err != nil {
			return err
		}

		images, err := result.Resolve(ctx)
		if err != nil {
			return err
		}

		ret = &FindImagesResultType{
			Count:      result.Count,
			Images:     images,
			Megapixels: result.Megapixels,
			Filesize:   result.TotalSize,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllImages(ctx context.Context) (ret []*models.Image, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Image.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
