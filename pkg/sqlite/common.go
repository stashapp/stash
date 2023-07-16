package sqlite

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

type oCounterManager struct {
	tableMgr *table
}

func (qb *oCounterManager) getOCounter(ctx context.Context, id int) (int, error) {
	q := dialect.From(qb.tableMgr.table).Select("o_counter").Where(goqu.Ex{"id": id})

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

func (qb *oCounterManager) IncrementOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.tableMgr.updateByID(ctx, id, goqu.Record{
		"o_counter": goqu.L("o_counter + 1"),
	}); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}

func (qb *oCounterManager) DecrementOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	table := qb.tableMgr.table
	q := dialect.Update(table).Set(goqu.Record{
		"o_counter": goqu.L("o_counter - 1"),
	}).Where(qb.tableMgr.byID(id), goqu.L("o_counter > 0"))

	if _, err := exec(ctx, q); err != nil {
		return 0, fmt.Errorf("updating %s: %w", table.GetTable(), err)
	}

	return qb.getOCounter(ctx, id)
}

func (qb *oCounterManager) ResetOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.tableMgr.updateByID(ctx, id, goqu.Record{
		"o_counter": 0,
	}); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}
