package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindTag(ctx context.Context, id string) (ret *models.Tag, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Tag().Find(idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindTags(ctx context.Context, tagFilter *models.TagFilterType, filter *models.FindFilterType) (ret *FindTagsResultType, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		tags, total, err := repo.Tag().Query(tagFilter, filter)
		if err != nil {
			return err
		}

		ret = &FindTagsResultType{
			Count: total,
			Tags:  tags,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllTags(ctx context.Context) (ret []*models.Tag, err error) {
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Tag().All()
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
