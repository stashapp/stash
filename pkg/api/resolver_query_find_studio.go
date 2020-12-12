package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func (r *queryResolver) FindStudio(ctx context.Context, id string) (*models.Studio, error) {
	qb := sqlite.NewStudioQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt, nil)
}

func (r *queryResolver) FindStudios(ctx context.Context, studioFilter *models.StudioFilterType, filter *models.FindFilterType) (*models.FindStudiosResultType, error) {
	qb := sqlite.NewStudioQueryBuilder()
	studios, total := qb.Query(studioFilter, filter)
	return &models.FindStudiosResultType{
		Count:   total,
		Studios: studios,
	}, nil
}

func (r *queryResolver) AllStudios(ctx context.Context) ([]*models.Studio, error) {
	qb := sqlite.NewStudioQueryBuilder()
	return qb.All()
}

func (r *queryResolver) AllStudiosSlim(ctx context.Context) ([]*models.Studio, error) {
	qb := sqlite.NewStudioQueryBuilder()
	return qb.AllSlim()
}
