package sqlite

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

	"github.com/stashapp/stash/pkg/models"
)

// currentUserID returns the ID of the user set in the context.
// Returns 0 if there is no stored user, such as for the guest user or background tasks.
// 0 never matches a user, so joins using it return no user data.
func currentUserID(ctx context.Context) int {
	id, _ := models.UserIDFromContext(ctx)
	return id
}

// userDataTable is a table storing per-user data for an object type,
// keyed by the object ID and user ID.
type userDataTable struct {
	table       string
	idColumn    string
	objectTable string
}

func (t userDataTable) goquTable() exp.IdentifierExpression {
	return goqu.T(t.table)
}

// onClause returns the clause joining the user data table to the object table for the given user.
func (t userDataTable) onClause(userID int) string {
	return fmt.Sprintf("%[1]s.%[2]s = %[3]s.id AND %[1]s.%[4]s = %[5]d", t.table, t.idColumn, t.objectTable, userIDColumn, userID)
}

// join adds a left join on the user data table for the current user to the filter.
func (t userDataTable) join(ctx context.Context, f *filterBuilder) {
	f.addLeftJoin(t.table, "", t.onClause(currentUserID(ctx)))
}

// joinSort adds a left join on the user data table for the current user to the query, for sorting.
func (t userDataTable) joinSort(ctx context.Context, q *queryBuilder) {
	q.joinSort(t.table, "", t.onClause(currentUserID(ctx)))
}

// column returns the fully qualified name of a column in the user data table.
func (t userDataTable) column(col string) string {
	return t.table + "." + col
}

// columnOrDefault returns an expression for a column in the user data table,
// using def when the user has no row for the object.
func (t userDataTable) columnOrDefault(col string, def string) string {
	return "COALESCE(" + t.column(col) + ", " + def + ")"
}

// orderBy returns an ORDER BY clause on a column of the current user's data for the object.
// It uses a subquery, so it can be used where the sort cannot add a join.
func (t userDataTable) orderBy(ctx context.Context, col string, direction string) string {
	return fmt.Sprintf(" ORDER BY (SELECT %[1]s.%[2]s FROM %[1]s WHERE %[1]s.%[3]s = %[4]s.id AND %[1]s.%[5]s = %[6]d) %[7]s",
		t.table, col, t.idColumn, t.objectTable, userIDColumn, currentUserID(ctx), getSortDirection(direction))
}

// upsert sets the provided values for the object and the current user.
// Existing values not in values are left unchanged.
// Does nothing if there is no current user.
func (t userDataTable) upsert(ctx context.Context, id int, values goqu.Record) error {
	updates := goqu.Record{}
	for col := range values {
		updates[col] = goqu.I("excluded." + col)
	}

	return t.upsertWith(ctx, id, values, updates)
}

// upsertWith inserts the provided values for the object and the current user,
// or applies updates if a row already exists.
// Does nothing if there is no current user.
func (t userDataTable) upsertWith(ctx context.Context, id int, values goqu.Record, updates goqu.Record) error {
	userID, ok := models.UserIDFromContext(ctx)
	if !ok || len(values) == 0 {
		return nil
	}

	row := goqu.Record{
		t.idColumn:   id,
		userIDColumn: userID,
	}
	for col, v := range values {
		row[col] = v
	}

	conflictKey := t.idColumn + ", " + userIDColumn
	q := dialect.Insert(t.goquTable()).Prepared(true).Rows(row).
		OnConflict(goqu.DoUpdate(conflictKey, updates))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("upserting %s: %w", t.table, err)
	}

	return nil
}

// selectForUser returns a query selecting the current user's rows for the provided object IDs.
// Returns nil if there is no current user.
func (t userDataTable) selectForUser(ctx context.Context, ids []int) *goqu.SelectDataset {
	userID, ok := models.UserIDFromContext(ctx)
	if !ok || len(ids) == 0 {
		return nil
	}

	table := t.goquTable()
	return dialect.From(table).Select(table.All()).Where(
		table.Col(userIDColumn).Eq(userID),
		table.Col(t.idColumn).In(ids),
	)
}

// joinFn returns a join function for criterion handlers that joins the user data table for the current user.
// The join is always a left join, since objects without user data must still match null and default values.
func (t userDataTable) joinFn(ctx context.Context) func(f *filterBuilder, _ joinType) {
	return func(f *filterBuilder, _ joinType) {
		t.join(ctx, f)
	}
}

// intCriterionHandler returns a criterion handler filtering on an integer expression of the user data table.
func (t userDataTable) intCriterionHandler(c *models.IntCriterionInput, expr string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		intCriterionHandler(c, expr, t.joinFn(ctx))(ctx, f)
	}
}

// boolCriterionHandler returns a criterion handler filtering on a boolean expression of the user data table.
func (t userDataTable) boolCriterionHandler(c *bool, expr string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		boolCriterionHandler(c, expr, t.joinFn(ctx))(ctx, f)
	}
}

// floatIntCriterionHandler returns a criterion handler filtering on a float expression of the user data table,
// using an integer criterion.
func (t userDataTable) floatIntCriterionHandler(c *models.IntCriterionInput, expr string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		floatIntCriterionHandler(c, expr, t.joinFn(ctx))(ctx, f)
	}
}
