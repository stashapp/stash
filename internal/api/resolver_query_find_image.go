package api

import (
	"context"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
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

func (r *queryResolver) FindImages(ctx context.Context, imageFilter *models.ImageFilterType, imageIds []int, filter *models.FindFilterType) (ret *FindImagesResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		qb := repo.Image()

		fields := graphql.CollectAllFields(ctx)

		result, err := qb.Query(models.ImageQueryOptions{
			QueryOptions: models.QueryOptions{
				FindFilter: filter,
				Count:      stringslice.StrInclude(fields, "count"),
			},
			ImageFilter: imageFilter,
			Megapixels:  stringslice.StrInclude(fields, "megapixels"),
			TotalSize:   stringslice.StrInclude(fields, "filesize"),
		})
		if err != nil {
			return err
		}

		images, err := result.Resolve()
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
