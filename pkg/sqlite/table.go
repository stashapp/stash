package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

type table struct {
	table    exp.IdentifierExpression
	idColumn exp.IdentifierExpression
}

type NotFoundError struct {
	ID    int
	Table string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("id %d does not exist in %s", e.ID, e.Table)
}

func (t *table) insert(ctx context.Context, o interface{}) (sql.Result, error) {
	q := dialect.Insert(t.table).Prepared(true).Rows(o)
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
	q := dialect.Update(t.table).Prepared(true).Set(o).Where(t.byID(id))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("updating %s: %w", t.table.GetTable(), err)
	}

	return nil
}

func (t *table) byID(id interface{}) exp.Expression {
	return t.idColumn.Eq(id)
}

func (t *table) idExists(ctx context.Context, id int) (bool, error) {
	q := dialect.Select(goqu.COUNT("*")).From(t.table).Where(t.byID(id))

	var count int
	if err := querySimple(ctx, q, &count); err != nil {
		return false, err
	}

	return count == 1, nil
}

func (t *table) checkIDExists(ctx context.Context, id int) error {
	exists, err := t.idExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return &NotFoundError{ID: id, Table: t.table.GetTable()}
	}

	return nil
}

func (t *table) destroyExisting(ctx context.Context, ids []int) error {
	for _, id := range ids {
		exists, err := t.idExists(ctx, id)
		if err != nil {
			return err
		}

		if !exists {
			return &NotFoundError{
				ID:    id,
				Table: t.table.GetTable(),
			}
		}
	}

	return t.destroy(ctx, ids)
}

func (t *table) destroy(ctx context.Context, ids []int) error {
	q := dialect.Delete(t.table).Where(t.idColumn.In(ids))

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

func (t *joinTable) invert() *joinTable {
	return &joinTable{
		table: table{
			table:    t.table.table,
			idColumn: t.fkColumn,
		},
		fkColumn: t.table.idColumn,
	}
}

func (t *joinTable) get(ctx context.Context, id int) ([]int, error) {
	q := dialect.Select(t.fkColumn).From(t.table.table).Where(t.idColumn.Eq(id))

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
	q := dialect.Insert(t.table.table).Cols(t.idColumn.GetCol(), t.fkColumn.GetCol()).Vals(
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
	q := dialect.Delete(t.table.table).Where(
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

type stashIDTable struct {
	table
}

type stashIDRow struct {
	StashID  null.String `db:"stash_id"`
	Endpoint null.String `db:"endpoint"`
}

func (r *stashIDRow) resolve() *models.StashID {
	return &models.StashID{
		StashID:  r.StashID.String,
		Endpoint: r.Endpoint.String,
	}
}

func (t *stashIDTable) get(ctx context.Context, id int) ([]*models.StashID, error) {
	q := dialect.Select("endpoint", "stash_id").From(t.table.table).Where(t.idColumn.Eq(id))

	const single = false
	var ret []*models.StashID
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var v stashIDRow
		if err := rows.StructScan(&v); err != nil {
			return err
		}

		ret = append(ret, v.resolve())

		return nil
	}); err != nil {
		return nil, fmt.Errorf("getting stash ids from %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *stashIDTable) insertJoin(ctx context.Context, id int, v models.StashID) (sql.Result, error) {
	q := dialect.Insert(t.table.table).Cols(t.idColumn.GetCol(), "endpoint", "stash_id").Vals(
		goqu.Vals{id, v.Endpoint, v.StashID},
	)
	ret, err := exec(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("inserting into %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *stashIDTable) insertJoins(ctx context.Context, id int, v []models.StashID) error {
	for _, fk := range v {
		if _, err := t.insertJoin(ctx, id, fk); err != nil {
			return err
		}
	}

	return nil
}

func (t *stashIDTable) replaceJoins(ctx context.Context, id int, v []models.StashID) error {
	if err := t.destroy(ctx, []int{id}); err != nil {
		return err
	}

	return t.insertJoins(ctx, id, v)
}

func (t *stashIDTable) addJoins(ctx context.Context, id int, v []models.StashID) error {
	// get existing foreign keys
	fks, err := t.get(ctx, id)
	if err != nil {
		return err
	}

	// only add values that are not already present
	var filtered []models.StashID
	for _, vv := range v {
		for _, e := range fks {
			if vv.Endpoint == e.Endpoint {
				continue
			}

			filtered = append(filtered, vv)
		}
	}
	return t.insertJoins(ctx, id, filtered)
}

func (t *stashIDTable) destroyJoins(ctx context.Context, id int, v []models.StashID) error {
	for _, vv := range v {
		q := dialect.Delete(t.table.table).Where(
			t.idColumn.Eq(id),
			t.table.table.Col("endpoint").Eq(vv.Endpoint),
			t.table.table.Col("stash_id").Eq(vv.StashID),
		)

		if _, err := exec(ctx, q); err != nil {
			return fmt.Errorf("destroying %s: %w", t.table.table.GetTable(), err)
		}
	}

	return nil
}

func (t *stashIDTable) modifyJoins(ctx context.Context, id int, v []models.StashID, mode models.RelationshipUpdateMode) error {
	switch mode {
	case models.RelationshipUpdateModeSet:
		return t.replaceJoins(ctx, id, v)
	case models.RelationshipUpdateModeAdd:
		return t.addJoins(ctx, id, v)
	case models.RelationshipUpdateModeRemove:
		return t.destroyJoins(ctx, id, v)
	}

	return nil
}

type scenesMoviesTable struct {
	table
}

type moviesScenesRow struct {
	MovieID    null.Int `db:"movie_id"`
	SceneIndex null.Int `db:"scene_index"`
}

func (r moviesScenesRow) resolve(sceneID int) models.MoviesScenes {
	return models.MoviesScenes{
		MovieID:    int(r.MovieID.Int64),
		SceneIndex: nullIntPtr(r.SceneIndex),
	}
}

func (t *scenesMoviesTable) get(ctx context.Context, id int) ([]models.MoviesScenes, error) {
	q := dialect.Select("movie_id", "scene_index").From(t.table.table).Where(t.idColumn.Eq(id))

	const single = false
	var ret []models.MoviesScenes
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var v moviesScenesRow
		if err := rows.StructScan(&v); err != nil {
			return err
		}

		ret = append(ret, v.resolve(id))

		return nil
	}); err != nil {
		return nil, fmt.Errorf("getting scene movies from %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *scenesMoviesTable) insertJoin(ctx context.Context, id int, v models.MoviesScenes) (sql.Result, error) {
	q := dialect.Insert(t.table.table).Cols(t.idColumn.GetCol(), "movie_id", "scene_index").Vals(
		goqu.Vals{id, v.MovieID, intFromPtr(v.SceneIndex)},
	)
	ret, err := exec(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("inserting into %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *scenesMoviesTable) insertJoins(ctx context.Context, id int, v []models.MoviesScenes) error {
	for _, fk := range v {
		if _, err := t.insertJoin(ctx, id, fk); err != nil {
			return err
		}
	}

	return nil
}

func (t *scenesMoviesTable) replaceJoins(ctx context.Context, id int, v []models.MoviesScenes) error {
	if err := t.destroy(ctx, []int{id}); err != nil {
		return err
	}

	return t.insertJoins(ctx, id, v)
}

func (t *scenesMoviesTable) addJoins(ctx context.Context, id int, v []models.MoviesScenes) error {
	// get existing foreign keys
	fks, err := t.get(ctx, id)
	if err != nil {
		return err
	}

	// only add values that are not already present
	var filtered []models.MoviesScenes
	for _, vv := range v {
		for _, e := range fks {
			if vv.MovieID == e.MovieID {
				continue
			}

			filtered = append(filtered, vv)
		}
	}
	return t.insertJoins(ctx, id, filtered)
}

func (t *scenesMoviesTable) destroyJoins(ctx context.Context, id int, v []models.MoviesScenes) error {
	for _, vv := range v {
		q := dialect.Delete(t.table.table).Where(
			t.idColumn.Eq(id),
			t.table.table.Col("movie_id").Eq(vv.MovieID),
		)

		if _, err := exec(ctx, q); err != nil {
			return fmt.Errorf("destroying %s: %w", t.table.table.GetTable(), err)
		}
	}

	return nil
}

func (t *scenesMoviesTable) modifyJoins(ctx context.Context, id int, v []models.MoviesScenes, mode models.RelationshipUpdateMode) error {
	switch mode {
	case models.RelationshipUpdateModeSet:
		return t.replaceJoins(ctx, id, v)
	case models.RelationshipUpdateModeAdd:
		return t.addJoins(ctx, id, v)
	case models.RelationshipUpdateModeRemove:
		return t.destroyJoins(ctx, id, v)
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

	logger.Tracef("SQL: %s [%v]", sql, args)
	ret, err := tx.ExecContext(ctx, sql, args...)
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

	logger.Tracef("SQL: %s [%v]", q, args)
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

	logger.Tracef("SQL: %s [%v]", q, args)
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
