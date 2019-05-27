package api

import (
	"context"
	"github.com/stashapp/stash/pkg/models"
	"strconv"
)

func (r *queryResolver) FindPerformer(ctx context.Context, id string) (*models.Performer, error) {
	qb := models.NewPerformerQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt)
}

func (r *queryResolver) FindPerformers(ctx context.Context, performerFilter *models.PerformerFilterType, filter *models.FindFilterType) (*models.FindPerformersResultType, error) {
	qb := models.NewPerformerQueryBuilder()
	performers, total := qb.Query(performerFilter, filter)
	return &models.FindPerformersResultType{
		Count:      total,
		Performers: performers,
	}, nil
}

func (r *queryResolver) AllPerformers(ctx context.Context) ([]*models.Performer, error) {
	qb := models.NewPerformerQueryBuilder()
	return qb.All()
}
