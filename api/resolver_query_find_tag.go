package api

import (
	"context"
	"github.com/stashapp/stash/models"
	"strconv"
)

func (r *queryResolver) FindTag(ctx context.Context, id string) (*models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt, nil)
}

func (r *queryResolver) AllTags(ctx context.Context) ([]models.Tag, error) {
	qb := models.NewTagQueryBuilder()
	return qb.All()
}
