package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindDvd(ctx context.Context, id string) (*models.Dvd, error) {
	qb := models.NewDvdQueryBuilder()
	idInt, _ := strconv.Atoi(id)
	return qb.Find(idInt, nil)
}

func (r *queryResolver) FindDvds(ctx context.Context, filter *models.FindFilterType) (*models.FindDvdsResultType, error) {
	qb := models.NewDvdQueryBuilder()
	dvds, total := qb.Query(filter)
	return &models.FindDvdsResultType{
		Count: total,
		Dvds:  dvds,
	}, nil
}

func (r *queryResolver) AllDvds(ctx context.Context) ([]*models.Dvd, error) {
	qb := models.NewDvdQueryBuilder()
	return qb.All()
}
