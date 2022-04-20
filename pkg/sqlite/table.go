package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type table struct {
	table exp.IdentifierExpression
}

func (t *table) insert(ctx context.Context, o interface{}) (sql.Result, error) {
	q := goqu.Insert(t.table).Rows(o)
	ret, err := exec(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("inserting into %s: %w", t.table.GetTable(), err)
	}

	return ret, nil
}

func (t *table) insertID(ctx context.Context, o interface{}) (int, error) {
	result, err := t.insert(ctx, o)
	if err != nil {
		return 0, err
	}

	ret, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ret), nil
}

func (t *table) updateByID(ctx context.Context, id interface{}, o interface{}) error {
	q := goqu.Update(t.table).Set(o).Where(t.byID(id))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("updating %s: %w", t.table.GetTable(), err)
	}

	return nil
}

func (t *table) byID(id interface{}) exp.Expression {
	return t.table.Col(idColumn).Eq(id)
}

func (t *table) idExists(ctx context.Context, id int) (bool, error) {
	q := goqu.Select(goqu.COUNT("*")).From(t.table).Where(t.byID(id))

	var count int
	if err := querySimple(ctx, q, &count); err != nil {
		return false, err
	}

	return count == 1, nil
}

func (t *table) destroyExisting(ctx context.Context, ids []int) error {
	for _, id := range ids {
		exists, err := t.idExists(ctx, id)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("id %d does not exist in %s", id, t.table.GetTable())
		}
	}

	return t.destroy(ctx, ids)
}

func (t *table) destroy(ctx context.Context, ids []int) error {
	q := goqu.Delete(t.table).Where(t.table.Col(idColumn).In(ids))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("destroying %s: %w", t.table.GetTable(), err)
	}

	return nil
}

func (t *table) get(ctx context.Context, q *goqu.SelectDataset, dest interface{}) error {
	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	sql, args, err := q.ToSQL()
	if err != nil {
		return fmt.Errorf("generating sql: %w", err)
	}

	return tx.GetContext(ctx, dest, sql, args...)
}

type sqler interface {
	ToSQL() (sql string, params []interface{}, err error)
}

func exec(ctx context.Context, stmt sqler) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, err
	}

	sql, args, err := stmt.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("generating sql: %w", err)
	}

	ret, err := tx.ExecContext(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("executing `%s` [%v]: %w", sql, args, err)
	}

	return ret, nil
}

func count(ctx context.Context, q *goqu.SelectDataset) (int, error) {
	var count int
	if err := querySimple(ctx, q, &count); err != nil {
		return 0, err
	}

	return count, nil
}

func queryFunc(ctx context.Context, query *goqu.SelectDataset, single bool, f func(rows *sqlx.Rows) error) error {
	q, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	rows, err := tx.QueryxContext(ctx, q, args...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("querying `%s` [%v]: %w", q, args, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := f(rows); err != nil {
			return err
		}
		if single {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func querySimple(ctx context.Context, query *goqu.SelectDataset, out interface{}) error {
	q, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	tx, err := getTx(ctx)
	if err != nil {
		return err
	}

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("querying `%s` [%v]: %w", q, args, err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(out); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
