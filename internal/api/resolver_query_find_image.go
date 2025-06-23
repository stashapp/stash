package api

import (
	"context"
	"slices"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
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

func (r *queryResolver) FindImages(
	ctx context.Context,
	imageFilter *models.ImageFilterType,
	imageIds []int,
	ids []string,
	filter *models.FindFilterType,
) (ret *FindImagesResultType, err error) {
	if len(ids) > 0 {
		imageIds, err = stringslice.StringSliceToIntSlice(ids)
		if err != nil {
			return nil, err
		}
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image

		var images []*models.Image
		fields := graphql.CollectAllFields(ctx)
		result := &models.ImageQueryResult{}

		if len(imageIds) > 0 {
			images, err = r.repository.Image.FindMany(ctx, imageIds)
			if err == nil {
				result.Count = len(images)
				for _, s := range images {
					if err = s.LoadPrimaryFile(ctx, r.repository.File); err != nil {
						break
					}

					f := s.Files.Primary()
					if f == nil {
						continue
					}

					imageFile, ok := f.(*models.ImageFile)
					if !ok {
						continue
					}

					result.Megapixels += float64(imageFile.Width*imageFile.Height) / float64(1000000)
					result.TotalSize += float64(f.Base().Size)
				}
			}
		} else {
			result, err = qb.Query(ctx, models.ImageQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: filter,
					Count:      slices.Contains(fields, "count"),
				},
				ImageFilter: imageFilter,
				Megapixels:  slices.Contains(fields, "megapixels"),
				TotalSize:   slices.Contains(fields, "filesize"),
			})
			if err == nil {
				images, err = result.Resolve(ctx)
			}
		}

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
