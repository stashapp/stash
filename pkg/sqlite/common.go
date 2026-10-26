package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/models"
)

var errNoCurrentUser = errors.New("no current user")

// oCounterManager manages the o-counter of the current user, stored in a user data table.
type oCounterManager struct {
	tableMgr *table
	userData userDataTable
}

func (qb *oCounterManager) getOCounter(ctx context.Context, id int) (int, error) {
	table := qb.userData.goquTable()
	q := dialect.From(table).Select("o_counter").Where(
		table.Col(qb.userData.idColumn).Eq(id),
		table.Col(userIDColumn).Eq(currentUserID(ctx)),
	)

	const single = true
	var ret int
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		if err := rows.Scan(&ret); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

// checkCanUpdate returns an error if the object does not exist or there is no current user.
func (qb *oCounterManager) checkCanUpdate(ctx context.Context, id int) error {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return err
	}

	if _, ok := models.UserIDFromContext(ctx); !ok {
		return fmt.Errorf("updating o-counter of %s %d: %w", qb.tableMgr.table.GetTable(), id, errNoCurrentUser)
	}

	return nil
}

func (qb *oCounterManager) IncrementOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.checkCanUpdate(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.userData.upsertWith(ctx, id, goqu.Record{
		"o_counter": 1,
	}, goqu.Record{
		"o_counter": goqu.L("o_counter + 1"),
	}); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}

func (qb *oCounterManager) DecrementOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.checkCanUpdate(ctx, id); err != nil {
		return 0, err
	}

	table := qb.userData.goquTable()
	q := dialect.Update(table).Set(goqu.Record{
		"o_counter": goqu.L("o_counter - 1"),
	}).Where(
		table.Col(qb.userData.idColumn).Eq(id),
		table.Col(userIDColumn).Eq(currentUserID(ctx)),
		goqu.L("o_counter > 0"),
	)

	if _, err := exec(ctx, q); err != nil {
		return 0, fmt.Errorf("updating %s: %w", qb.userData.table, err)
	}

	return qb.getOCounter(ctx, id)
}

func (qb *oCounterManager) ResetOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.checkCanUpdate(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.userData.upsert(ctx, id, goqu.Record{
		"o_counter": 0,
	}); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}
