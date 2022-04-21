package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

type table struct {
	table    exp.IdentifierExpression
	idColumn exp.IdentifierExpression
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
	return t.idColumn.Eq(id)
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
	q := goqu.Delete(t.table).Where(t.idColumn.In(ids))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("destroying %s: %w", t.table.GetTable(), err)
	}

	return nil
}

// func (t *table) get(ctx context.Context, q *goqu.SelectDataset, dest interface{}) error {
// 	tx, err := getTx(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	sql, args, err := q.ToSQL()
// 	if err != nil {
// 		return fmt.Errorf("generating sql: %w", err)
// 	}

// 	return tx.GetContext(ctx, dest, sql, args...)
// }

type joinTable struct {
	table
	fkColumn exp.IdentifierExpression
}

func (t *joinTable) get(ctx context.Context, id int) ([]int, error) {
	q := goqu.Select(t.fkColumn).From(t.table.table).Where(t.idColumn.Eq(id))

	const single = false
	var ret []int
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var fk int
		if err := rows.Scan(&fk); err != nil {
			return err
		}

		ret = append(ret, fk)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("getting foreign keys from %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *joinTable) insertJoin(ctx context.Context, id, foreignID int) (sql.Result, error) {
	q := goqu.Insert(t.table.table).Cols(t.idColumn.GetCol(), t.fkColumn.GetCol()).Vals(
		goqu.Vals{id, foreignID},
	)
	ret, err := exec(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("inserting into %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *joinTable) insertJoins(ctx context.Context, id int, foreignIDs []int) error {
	for _, fk := range foreignIDs {
		if _, err := t.insertJoin(ctx, id, fk); err != nil {
			return err
		}
	}

	return nil
}

func (t *joinTable) replaceJoins(ctx context.Context, id int, foreignIDs []int) error {
	if err := t.destroy(ctx, []int{id}); err != nil {
		return err
	}

	return t.insertJoins(ctx, id, foreignIDs)
}

func (t *joinTable) addJoins(ctx context.Context, id int, foreignIDs []int) error {
	// get existing foreign keys
	fks, err := t.get(ctx, id)
	if err != nil {
		return err
	}

	// only add foreign keys that are not already present
	foreignIDs = intslice.IntExclude(foreignIDs, fks)
	return t.insertJoins(ctx, id, foreignIDs)
}

func (t *joinTable) destroyJoins(ctx context.Context, id int, foreignIDs []int) error {
	q := goqu.Delete(t.table.table).Where(
		t.idColumn.Eq(id),
		t.fkColumn.In(foreignIDs),
	)

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("destroying %s: %w", t.table.table.GetTable(), err)
	}

	return nil
}

func (t *joinTable) modifyJoins(ctx context.Context, id int, foreignIDs []int, mode models.RelationshipUpdateMode) error {
	switch mode {
	case models.RelationshipUpdateModeSet:
		return t.replaceJoins(ctx, id, foreignIDs)
	case models.RelationshipUpdateModeAdd:
		return t.addJoins(ctx, id, foreignIDs)
	case models.RelationshipUpdateModeRemove:
		return t.destroyJoins(ctx, id, foreignIDs)
	}

	return nil
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
