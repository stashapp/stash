package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindTag(ctx context.Context, id string) (*models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt, nil)
}

func (r *queryResolver) FindTags(ctx context.Context, tagFilter *models.TagFilterType, filter *models.FindFilterType) (*models.FindTagsResultType, error) {
	qb := models.NewTagQueryBuilder()
	tags, total := qb.Query(tagFilter, filter)
	return &models.FindTagsResultType{
		Count: total,
		Tags:  tags,
	}, nil
}

func (r *queryResolver) AllTags(ctx context.Context) ([]*models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	return qb.All()
}

func (r *queryResolver) AllTagsSlim(ctx context.Context) ([]*models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	return qb.AllSlim()
}
