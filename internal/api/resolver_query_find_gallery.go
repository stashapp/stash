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

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Gallery().Find(idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindGalleries(ctx context.Context, galleryFilter *models.GalleryFilterType, filter *models.FindFilterType) (ret *FindGalleriesResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		galleries, total, err := repo.Gallery().Query(galleryFilter, filter)
		if err != nil {
			return err
		}

		ret = &FindGalleriesResultType{
			Count:     total,
			Galleries: galleries,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
