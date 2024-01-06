package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindGallery(ctx context.Context, id string) (ret *models.Gallery, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindGalleries(ctx context.Context, galleryFilter *models.GalleryFilterType, filter *models.FindFilterType) (ret *FindGalleriesResultType, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var galleries []*models.Gallery
		var err error

		var result *models.GalleryQueryResult

		result, err = r.repository.Gallery.Query(ctx, models.GalleryQueryOptions{
			QueryOptions: models.QueryOptions{
				FindFilter: filter,
				Count:      true,
			},
			GalleryFilter: galleryFilter,
		})
		if err == nil {
			galleries, err = result.Resolve(ctx)
		}

		if err != nil {
			return err
		}

		ret = &FindGalleriesResultType{
			Count:     result.Count,
			Galleries: galleries,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllGalleries(ctx context.Context) (ret []*models.Gallery, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
