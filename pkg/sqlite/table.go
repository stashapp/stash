package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
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

func (t *table) byIDInts(ids ...int) exp.Expression {
	ii := make([]interface{}, len(ids))
	for i, id := range ids {
		ii[i] = id
	}
	return t.idColumn.In(ii...)
}

func (t *table) idExists(ctx context.Context, id interface{}) (bool, error) {
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

func (t *table) join(j joiner, as string, parentIDCol string) {
	tableName := t.table.GetTable()
	tt := tableName
	if as != "" {
		tt = as
	}
	j.addLeftJoin(tableName, as, fmt.Sprintf("%s.%s = %s", tt, t.idColumn.GetCol(), parentIDCol))
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

	// required for ordering
	foreignTable *table
	orderBy      exp.OrderedExpression
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

	if t.orderBy != nil {
		if t.foreignTable != nil {
			q = q.InnerJoin(t.foreignTable.table, goqu.On(t.foreignTable.idColumn.Eq(t.fkColumn)))
		}
		q = q.Order(t.orderBy)
	}

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

func (t *joinTable) insertJoins(ctx context.Context, id int, foreignIDs []int) error {
	// manually create SQL so that we can prepare once
	// ignore duplicates
	q := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?) ON CONFLICT (%[2]s, %s) DO NOTHING", t.table.table.GetTable(), t.idColumn.GetCol(), t.fkColumn.GetCol())

	stmt, err := dbWrapper.Prepare(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// eliminate duplicates
	foreignIDs = sliceutil.AppendUniques(nil, foreignIDs)

	for _, fk := range foreignIDs {
		if _, err := dbWrapper.ExecStmt(ctx, stmt, id, fk); err != nil {
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
	foreignIDs = sliceutil.Exclude(foreignIDs, fks)
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

func (r *stashIDRow) resolve() models.StashID {
	return models.StashID{
		StashID:  r.StashID.String,
		Endpoint: r.Endpoint.String,
	}
}

func (t *stashIDTable) get(ctx context.Context, id int) ([]models.StashID, error) {
	q := dialect.Select("endpoint", "stash_id").From(t.table.table).Where(t.idColumn.Eq(id))

	const single = false
	var ret []models.StashID
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

type stringTable struct {
	table
	stringColumn exp.IdentifierExpression
}

func (t *stringTable) get(ctx context.Context, id int) ([]string, error) {
	q := dialect.Select(t.stringColumn).From(t.table.table).Where(t.idColumn.Eq(id))

	const single = false
	var ret []string
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var v string
		if err := rows.Scan(&v); err != nil {
			return err
		}

		ret = append(ret, v)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("getting stash ids from %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *stringTable) insertJoin(ctx context.Context, id int, v string) (sql.Result, error) {
	q := dialect.Insert(t.table.table).Cols(t.idColumn.GetCol(), t.stringColumn.GetCol()).Vals(
		goqu.Vals{id, v},
	)
	ret, err := exec(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("inserting into %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *stringTable) insertJoins(ctx context.Context, id int, v []string) error {
	for _, fk := range v {
		if _, err := t.insertJoin(ctx, id, fk); err != nil {
			return err
		}
	}

	return nil
}

func (t *stringTable) replaceJoins(ctx context.Context, id int, v []string) error {
	if err := t.destroy(ctx, []int{id}); err != nil {
		return err
	}

	return t.insertJoins(ctx, id, v)
}

func (t *stringTable) addJoins(ctx context.Context, id int, v []string) error {
	// get existing foreign keys
	existing, err := t.get(ctx, id)
	if err != nil {
		return err
	}

	// only add values that are not already present
	filtered := sliceutil.Exclude(v, existing)
	return t.insertJoins(ctx, id, filtered)
}

func (t *stringTable) destroyJoins(ctx context.Context, id int, v []string) error {
	for _, vv := range v {
		q := dialect.Delete(t.table.table).Where(
			t.idColumn.Eq(id),
			t.stringColumn.Eq(vv),
		)

		if _, err := exec(ctx, q); err != nil {
			return fmt.Errorf("destroying %s: %w", t.table.table.GetTable(), err)
		}
	}

	return nil
}

func (t *stringTable) modifyJoins(ctx context.Context, id int, v []string, mode models.RelationshipUpdateMode) error {
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

type orderedValueTable[T comparable] struct {
	table
	valueColumn exp.IdentifierExpression
}

func (t *orderedValueTable[T]) positionColumn() exp.IdentifierExpression {
	const positionColumn = "position"
	return t.table.table.Col(positionColumn)
}

func (t *orderedValueTable[T]) get(ctx context.Context, id int) ([]T, error) {
	q := dialect.Select(t.valueColumn).From(t.table.table).Where(t.idColumn.Eq(id)).Order(t.positionColumn().Asc())

	const single = false
	var ret []T
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var v T
		if err := rows.Scan(&v); err != nil {
			return err
		}

		ret = append(ret, v)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("getting stash ids from %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *orderedValueTable[T]) insertJoin(ctx context.Context, id int, position int, v T) (sql.Result, error) {
	q := dialect.Insert(t.table.table).Cols(t.idColumn.GetCol(), t.positionColumn().GetCol(), t.valueColumn.GetCol()).Vals(
		goqu.Vals{id, position, v},
	)
	ret, err := exec(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("inserting into %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *orderedValueTable[T]) insertJoins(ctx context.Context, id int, startPos int, v []T) error {
	for i, fk := range v {
		if _, err := t.insertJoin(ctx, id, i+startPos, fk); err != nil {
			return err
		}
	}

	return nil
}

func (t *orderedValueTable[T]) replaceJoins(ctx context.Context, id int, v []T) error {
	if err := t.destroy(ctx, []int{id}); err != nil {
		return err
	}

	const startPos = 0
	return t.insertJoins(ctx, id, startPos, v)
}

func (t *orderedValueTable[T]) addJoins(ctx context.Context, id int, v []T) error {
	// get existing foreign keys
	existing, err := t.get(ctx, id)
	if err != nil {
		return err
	}

	// only add values that are not already present
	filtered := sliceutil.Exclude(v, existing)

	if len(filtered) == 0 {
		return nil
	}

	startPos := len(existing)
	return t.insertJoins(ctx, id, startPos, filtered)
}

func (t *orderedValueTable[T]) destroyJoins(ctx context.Context, id int, v []T) error {
	existing, err := t.get(ctx, id)
	if err != nil {
		return fmt.Errorf("getting existing %s: %w", t.table.table.GetTable(), err)
	}

	newValue := sliceutil.Exclude(existing, v)
	if len(newValue) == len(existing) {
		return nil
	}

	return t.replaceJoins(ctx, id, newValue)
}

func (t *orderedValueTable[T]) modifyJoins(ctx context.Context, id int, v []T, mode models.RelationshipUpdateMode) error {
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

type scenesGroupsTable struct {
	table
}

type groupsScenesRow struct {
	SceneID    null.Int `db:"scene_id"`
	GroupID    null.Int `db:"group_id"`
	SceneIndex null.Int `db:"scene_index"`
}

func (r groupsScenesRow) resolve(sceneID int) models.GroupsScenes {
	return models.GroupsScenes{
		GroupID:    int(r.GroupID.Int64),
		SceneIndex: nullIntPtr(r.SceneIndex),
	}
}

func (t *scenesGroupsTable) get(ctx context.Context, id int) ([]models.GroupsScenes, error) {
	q := dialect.Select("group_id", "scene_index").From(t.table.table).Where(t.idColumn.Eq(id))

	const single = false
	var ret []models.GroupsScenes
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var v groupsScenesRow
		if err := rows.StructScan(&v); err != nil {
			return err
		}

		ret = append(ret, v.resolve(id))

		return nil
	}); err != nil {
		return nil, fmt.Errorf("getting scene groups from %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *scenesGroupsTable) insertJoin(ctx context.Context, id int, v models.GroupsScenes) (sql.Result, error) {
	q := dialect.Insert(t.table.table).Cols(t.idColumn.GetCol(), "group_id", "scene_index").Vals(
		goqu.Vals{id, v.GroupID, intFromPtr(v.SceneIndex)},
	)
	ret, err := exec(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("inserting into %s: %w", t.table.table.GetTable(), err)
	}

	return ret, nil
}

func (t *scenesGroupsTable) insertJoins(ctx context.Context, id int, v []models.GroupsScenes) error {
	for _, fk := range v {
		if _, err := t.insertJoin(ctx, id, fk); err != nil {
			return err
		}
	}

	return nil
}

func (t *scenesGroupsTable) replaceJoins(ctx context.Context, id int, v []models.GroupsScenes) error {
	if err := t.destroy(ctx, []int{id}); err != nil {
		return err
	}

	return t.insertJoins(ctx, id, v)
}

func (t *scenesGroupsTable) addJoins(ctx context.Context, id int, v []models.GroupsScenes) error {
	// get existing foreign keys
	fks, err := t.get(ctx, id)
	if err != nil {
		return err
	}

	// only add values that are not already present
	var filtered []models.GroupsScenes
	for _, vv := range v {
		found := false

		for _, e := range fks {
			if vv.GroupID == e.GroupID {
				found = true
				break
			}
		}

		if !found {
			filtered = append(filtered, vv)
		}
	}
	return t.insertJoins(ctx, id, filtered)
}

func (t *scenesGroupsTable) destroyJoins(ctx context.Context, id int, v []models.GroupsScenes) error {
	for _, vv := range v {
		q := dialect.Delete(t.table.table).Where(
			t.idColumn.Eq(id),
			t.table.table.Col("group_id").Eq(vv.GroupID),
		)

		if _, err := exec(ctx, q); err != nil {
			return fmt.Errorf("destroying %s: %w", t.table.table.GetTable(), err)
		}
	}

	return nil
}

func (t *scenesGroupsTable) modifyJoins(ctx context.Context, id int, v []models.GroupsScenes, mode models.RelationshipUpdateMode) error {
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

type relatedFilesTable struct {
	table
}

// type scenesFilesRow struct {
// 	SceneID int           `db:"scene_id"`
// 	Primary bool          `db:"primary"`
// 	FileID  models.FileID `db:"file_id"`
// }

func (t *relatedFilesTable) insertJoin(ctx context.Context, id int, primary bool, fileID models.FileID) error {
	q := dialect.Insert(t.table.table).Cols(t.idColumn.GetCol(), "primary", "file_id").Vals(
		goqu.Vals{id, primary, fileID},
	)
	_, err := exec(ctx, q)
	if err != nil {
		return fmt.Errorf("inserting into %s: %w", t.table.table.GetTable(), err)
	}

	return nil
}

func (t *relatedFilesTable) insertJoins(ctx context.Context, id int, firstPrimary bool, fileIDs []models.FileID) error {
	for i, fk := range fileIDs {
		if err := t.insertJoin(ctx, id, firstPrimary && i == 0, fk); err != nil {
			return err
		}
	}

	return nil
}

func (t *relatedFilesTable) replaceJoins(ctx context.Context, id int, fileIDs []models.FileID) error {
	if err := t.destroy(ctx, []int{id}); err != nil {
		return err
	}

	const firstPrimary = true
	return t.insertJoins(ctx, id, firstPrimary, fileIDs)
}

// destroyJoins destroys all entries in the table with the provided fileIDs
func (t *relatedFilesTable) destroyJoins(ctx context.Context, fileIDs []models.FileID) error {
	q := dialect.Delete(t.table.table).Where(t.table.table.Col("file_id").In(fileIDs))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("destroying file joins in %s: %w", t.table.table.GetTable(), err)
	}

	return nil
}

func (t *relatedFilesTable) setPrimary(ctx context.Context, id int, fileID models.FileID) error {
	table := t.table.table

	q := dialect.Update(table).Prepared(true).Set(goqu.Record{
		"primary": 0,
	}).Where(t.idColumn.Eq(id), table.Col(fileIDColumn).Neq(fileID))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("unsetting primary flags in %s: %w", t.table.table.GetTable(), err)
	}

	q = dialect.Update(table).Prepared(true).Set(goqu.Record{
		"primary": 1,
	}).Where(t.idColumn.Eq(id), table.Col(fileIDColumn).Eq(fileID))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("setting primary flag in %s: %w", t.table.table.GetTable(), err)
	}

	return nil
}

type viewHistoryTable struct {
	table
	dateColumn exp.IdentifierExpression
}

func (t *viewHistoryTable) getDates(ctx context.Context, id int) ([]time.Time, error) {
	table := t.table.table

	q := dialect.Select(
		t.dateColumn,
	).From(table).Where(
		t.idColumn.Eq(id),
	).Order(t.dateColumn.Desc())

	const single = false
	var ret []time.Time
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var date Timestamp
		if err := rows.Scan(&date); err != nil {
			return err
		}
		ret = append(ret, date.Timestamp)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *viewHistoryTable) getManyDates(ctx context.Context, ids []int) ([][]time.Time, error) {
	table := t.table.table

	q := dialect.Select(
		t.idColumn,
		t.dateColumn,
	).From(table).Where(
		t.idColumn.In(ids),
	).Order(t.dateColumn.Desc())

	ret := make([][]time.Time, len(ids))
	idToIndex := idToIndexMap(ids)

	const single = false
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var id int
		var date Timestamp
		if err := rows.Scan(&id, &date); err != nil {
			return err
		}

		idx := idToIndex[id]
		ret[idx] = append(ret[idx], date.Timestamp)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *viewHistoryTable) getLastDate(ctx context.Context, id int) (*time.Time, error) {
	table := t.table.table
	q := dialect.Select(t.dateColumn).From(table).Where(
		t.idColumn.Eq(id),
	).Order(t.dateColumn.Desc()).Limit(1)

	var date NullTimestamp
	if err := querySimple(ctx, q, &date); err != nil {
		return nil, err
	}

	return date.TimePtr(), nil
}

func (t *viewHistoryTable) getManyLastDate(ctx context.Context, ids []int) ([]*time.Time, error) {
	table := t.table.table

	q := dialect.Select(
		t.idColumn,
		goqu.MAX(t.dateColumn),
	).From(table).Where(
		t.idColumn.In(ids),
	).GroupBy(t.idColumn)

	ret := make([]*time.Time, len(ids))
	idToIndex := idToIndexMap(ids)

	const single = false
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var id int

		// MAX appears to return a string, so handle it manually
		var dateString string

		if err := rows.Scan(&id, &dateString); err != nil {
			return err
		}

		t, err := time.Parse(TimestampFormat, dateString)
		if err != nil {
			return fmt.Errorf("parsing date %v: %w", dateString, err)
		}

		idx := idToIndex[id]
		ret[idx] = &t
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *viewHistoryTable) getCount(ctx context.Context, id int) (int, error) {
	table := t.table.table
	q := dialect.Select(goqu.COUNT("*")).From(table).Where(t.idColumn.Eq(id))

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

func (t *viewHistoryTable) getManyCount(ctx context.Context, ids []int) ([]int, error) {
	table := t.table.table

	q := dialect.Select(
		t.idColumn,
		goqu.COUNT(t.dateColumn),
	).From(table).Where(
		t.idColumn.In(ids),
	).GroupBy(t.idColumn)

	ret := make([]int, len(ids))
	idToIndex := idToIndexMap(ids)

	const single = false
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		var id int
		var count int
		if err := rows.Scan(&id, &count); err != nil {
			return err
		}

		idx := idToIndex[id]
		ret[idx] = count
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *viewHistoryTable) getAllCount(ctx context.Context) (int, error) {
	table := t.table.table
	q := dialect.Select(goqu.COUNT("*")).From(table)

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

func (t *viewHistoryTable) getUniqueCount(ctx context.Context) (int, error) {
	table := t.table.table
	q := dialect.Select(goqu.COUNT(goqu.DISTINCT(t.idColumn))).From(table)

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

func (t *viewHistoryTable) addDates(ctx context.Context, id int, dates []time.Time) ([]time.Time, error) {
	table := t.table.table

	if len(dates) == 0 {
		dates = []time.Time{time.Now()}
	}

	for _, d := range dates {
		q := dialect.Insert(table).Cols(t.idColumn.GetCol(), t.dateColumn.GetCol()).Vals(
			// convert all dates to UTC
			goqu.Vals{id, UTCTimestamp{Timestamp{d}}},
		)

		if _, err := exec(ctx, q); err != nil {
			return nil, fmt.Errorf("inserting into %s: %w", table.GetTable(), err)
		}
	}

	return t.getDates(ctx, id)
}

func (t *viewHistoryTable) deleteDates(ctx context.Context, id int, dates []time.Time) ([]time.Time, error) {
	table := t.table.table

	mostRecent := false
	if len(dates) == 0 {
		mostRecent = true
		dates = []time.Time{time.Now()}
	}

	for _, date := range dates {
		var subquery *goqu.SelectDataset
		if mostRecent {
			// delete the most recent
			subquery = dialect.Select("rowid").From(table).Where(
				t.idColumn.Eq(id),
			).Order(t.dateColumn.Desc()).Limit(1)
		} else {
			subquery = dialect.Select("rowid").From(table).Where(
				t.idColumn.Eq(id),
				t.dateColumn.Eq(UTCTimestamp{Timestamp{date}}),
			).Limit(1)
		}

		q := dialect.Delete(table).Where(goqu.I("rowid").Eq(subquery))

		if _, err := exec(ctx, q); err != nil {
			return nil, fmt.Errorf("deleting from %s: %w", table.GetTable(), err)
		}
	}

	return t.getDates(ctx, id)
}

func (t *viewHistoryTable) deleteAllDates(ctx context.Context, id int) (int, error) {
	table := t.table.table
	q := dialect.Delete(table).Where(t.idColumn.Eq(id))

	if _, err := exec(ctx, q); err != nil {
		return 0, fmt.Errorf("resetting dates for id %v: %w", id, err)
	}

	return t.getCount(ctx, id)
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

	rows, err := dbWrapper.QueryxContext(ctx, q, args...)

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

	rows, err := dbWrapper.QueryxContext(ctx, q, args...)
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

// func cols(table exp.IdentifierExpression, cols []string) []interface{} {
// 	var ret []interface{}
// 	for _, c := range cols {
// 		ret = append(ret, table.Col(c))
// 	}
// 	return ret
// }
