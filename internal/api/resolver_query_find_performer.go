package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
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

func (r *queryResolver) FindPerformers(ctx context.Context, performerFilter *models.PerformerFilterType, filter *models.FindFilterType, performerIDs []int, ids []string) (ret *FindPerformersResultType, err error) {
	if len(ids) > 0 {
		performerIDs, err = handleIDList(ids, "ids")
		if err != nil {
			return nil, err
		}
	}

	// #5682 - convert JSON numbers to float64 or int64
	if performerFilter != nil {
		performerFilter.CustomFields = convertCustomFieldCriterionInputJSONNumbers(performerFilter.CustomFields)
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var performers []*models.Performer
		var err error
		var total int

		if len(performerIDs) > 0 {
			performers, err = r.repository.Performer.FindMany(ctx, performerIDs)
			total = len(performers)
		} else {
			performers, total, err = r.repository.Performer.Query(ctx, performerFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &FindPerformersResultType{
			Count:      total,
			Performers: performers,
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
