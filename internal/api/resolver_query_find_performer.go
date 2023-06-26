package api

import (
	"context"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *queryResolver) FindPerformer(ctx context.Context, id string) (ret *models.Performer, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindPerformers(ctx context.Context, performerFilter *models.PerformerFilterType, filter *models.FindFilterType) (ret *FindPerformersResultType, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {

		var performers []*models.Performer
		var filteredCounts []*models.FilteredCounts

		fields := graphql.CollectAllFields(ctx)
		result := &models.PerformerQueryResult{}

		result, err = r.repository.Performer.QueryWithOptions(ctx, models.PerformerQueryOptions{
			QueryOptions: models.QueryOptions{
				FindFilter: filter,
			},
			PerformerFilter: performerFilter,
			FilteredCounts:  performerFilter != nil && performerFilter.Performers != nil && stringslice.StrInclude(fields, "filteredCounts"),
		})
		if err == nil {
			performers, filteredCounts, err = result.Resolve(ctx)
		}

		ret = &FindPerformersResultType{
			Count:          result.Count,
			Performers:     performers,
			FilteredCounts: filteredCounts,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllPerformers(ctx context.Context) (ret []*models.Performer, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
