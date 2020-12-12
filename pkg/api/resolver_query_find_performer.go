package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func (r *queryResolver) FindPerformer(ctx context.Context, id string) (*models.Performer, error) {
	qb := sqlite.NewPerformerQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt)
}

func (r *queryResolver) FindPerformers(ctx context.Context, performerFilter *models.PerformerFilterType, filter *models.FindFilterType) (*models.FindPerformersResultType, error) {
	qb := sqlite.NewPerformerQueryBuilder()
	performers, total := qb.Query(performerFilter, filter)
	return &models.FindPerformersResultType{
		Count:      total,
		Performers: performers,
	}, nil
}

func (r *queryResolver) AllPerformers(ctx context.Context) ([]*models.Performer, error) {
	qb := sqlite.NewPerformerQueryBuilder()
	return qb.All()
}

func (r *queryResolver) AllPerformersSlim(ctx context.Context) ([]*models.Performer, error) {
	qb := sqlite.NewPerformerQueryBuilder()
	return qb.AllSlim()
}
