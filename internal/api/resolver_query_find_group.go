package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *queryResolver) FindGroup(ctx context.Context, id string) (ret *models.Group, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Group.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindGroups(ctx context.Context, groupFilter *models.GroupFilterType, filter *models.FindFilterType, ids []string) (ret *FindGroupsResultType, err error) {
	idInts, err := stringslice.StringSliceToIntSlice(ids)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var groups []*models.Group
		var err error
		var total int

		if len(idInts) > 0 {
			groups, err = r.repository.Group.FindMany(ctx, idInts)
			total = len(groups)
		} else {
			groups, total, err = r.repository.Group.Query(ctx, groupFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &FindGroupsResultType{
			Count:  total,
			Groups: groups,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
