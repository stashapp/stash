package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

const (
	userTable      = "users"
	userRolesTable = "user_roles"
	userIDColumn   = "user_id"
	userRoleColumn = "role"
)

type userRow struct {
	ID           int         `db:"id" goqu:"skipinsert"`
	Username     string      `db:"username"`
	Locked       bool        `db:"locked"`
	ApiKey       zero.String `db:"api_key"`
	PasswordHash zero.String `db:"password_hash"`
	CreatedAt    Timestamp   `db:"created_at"`
	UpdatedAt    Timestamp   `db:"updated_at"`
}

func (r *userRow) fromUser(o models.User) {
	r.ID = o.ID
	r.Username = o.Username
	r.Locked = o.Locked
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

func (r *userRow) resolve() *models.User {
	ret := &models.User{
		ID:         r.ID,
		Username:   r.Username,
		Locked:     r.Locked,
		ApiKeyHash: r.ApiKey.String,
		CreatedAt:  r.CreatedAt.Timestamp,
		UpdatedAt:  r.UpdatedAt.Timestamp,
	}

	return ret
}

type userUpdateRow struct {
	ID        int         `db:"id" goqu:"skipinsert"`
	Username  string      `db:"username"`
	ApiKey    null.String `db:"api_key"`
	CreatedAt Timestamp   `db:"created_at"`
	UpdatedAt Timestamp   `db:"updated_at"`
}

func (r *userUpdateRow) fromUser(o models.User) {
	r.ID = o.ID
	r.Username = o.Username
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

type userQueryRow struct {
	userRow
	Role null.String `db:"role"`
}

func (r userQueryRow) appendRole(i *models.User) {
	if r.Role.Valid {
		i.Roles = sliceutil.AppendUnique(i.Roles, models.RoleEnum(r.Role.String))
	}
}

func (r *userQueryRow) resolve() *models.User {
	ret := r.userRow.resolve()
	r.appendRole(ret)

	return ret
}

type userQueryRows []userQueryRow

func (r userQueryRows) resolve() []*models.User {
	var ret []*models.User
	var last *models.User
	var lastID int

	for _, row := range r {
		if last == nil || lastID != row.ID {
			f := row.resolve()
			last = f
			lastID = row.ID
			ret = append(ret, last)
			continue
		}

		// must be merging with previous row
		row.appendRole(last)
	}

	return ret
}

// type userRowRecord struct {
// 	updateRecord
// }

type userRepositoryType struct {
	repository
}

var (
	userRepository = userRepositoryType{
		repository: repository{
			tableName: userTable,
			idColumn:  idColumn,
		},
	}
)

type UserStore struct{}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (qb *UserStore) table() exp.IdentifierExpression {
	return userTableMgr.table
}

func (qb *UserStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(
		qb.table().All(),
		userRolesJoinTable.Col(userRoleColumn).As("role"),
	).LeftJoin(
		userRolesJoinTable,
		goqu.On(
			userRolesJoinTable.Col(userIDColumn).Eq(qb.table().Col(idColumn)),
		),
	)
}

func (qb *UserStore) Create(ctx context.Context, newObject *models.User, passwordHash string) error {
	var r userRow
	r.fromUser(*newObject)
	r.PasswordHash = zero.StringFrom(passwordHash)

	id, err := userTableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if err := qb.SetRoles(ctx, id, newObject.Roles); err != nil {
		return fmt.Errorf("setting user roles: %w", err)
	}

	updated, err := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject = *updated

	return nil
}

func (qb *UserStore) Update(ctx context.Context, updatedObject *models.User) error {
	var r userUpdateRow
	r.fromUser(*updatedObject)

	if err := userTableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if err := qb.SetRoles(ctx, updatedObject.ID, updatedObject.Roles); err != nil {
		return fmt.Errorf("setting user roles: %w", err)
	}

	return nil
}

func (qb *UserStore) GetPasswordHash(ctx context.Context, id int) (string, error) {
	t := qb.table()
	q := dialect.Select(t.Col("password_hash")).From(t).Where(t.Col("id").Eq(id))

	var passwordHash string
	if err := queryFunc(ctx, q, true, func(r *sqlx.Rows) error {
		if err := r.Scan(&passwordHash); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", fmt.Errorf("getting password hash: %w", err)
	}

	return passwordHash, nil
}

func (qb *UserStore) SetUserPassword(ctx context.Context, id int, newPassword string) error {
	t := qb.table()
	q := dialect.Update(t).Prepared(true).Set(
		goqu.Record{"password_hash": null.StringFrom(newPassword)},
	).Where(t.Col("id").Eq(id))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("updating %s: %w", t.GetTable(), err)
	}

	return nil
}

func (qb *UserStore) SetUserAPIKey(ctx context.Context, id int, newAPIKey string) error {
	t := qb.table()
	q := dialect.Update(t).Prepared(true).Set(
		goqu.Record{"api_key": null.StringFrom(newAPIKey)},
	).Where(t.Col("id").Eq(id))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("updating %s: %w", t.GetTable(), err)
	}

	return nil
}

func (qb *UserStore) SetRoles(ctx context.Context, id int, roles models.Roles) error {
	// No need to update the user record itself, just ensure the roles are updated in the roles table
	return userRolesTableMgr.replaceJoins(ctx, id, roles.Strings())
}

func (qb *UserStore) SetLock(ctx context.Context, id int, locked bool) error {
	t := qb.table()
	q := dialect.Update(t).Prepared(true).Set(
		goqu.Record{"locked": locked},
	).Where(t.Col("id").Eq(id))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("updating %s: %w", t.GetTable(), err)
	}

	return nil
}

func (qb *UserStore) Destroy(ctx context.Context, id int) error {
	return userRepository.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *UserStore) Find(ctx context.Context, id int) (*models.User, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *UserStore) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	q := qb.selectDataset().Where(qb.table().Col("username").Eq(username))
	ret, err := qb.get(ctx, q)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *UserStore) FindAdminUsers(ctx context.Context) ([]*models.User, error) {
	q := qb.selectDataset().Where(qb.table().Col("id").In(
		dialect.Select(userIDColumn).From(userRolesJoinTable).Where(userRolesJoinTable.Col(userRoleColumn).Eq(models.RoleEnumAdmin.String())),
	))

	return qb.getMany(ctx, q)
}

func (qb *UserStore) FindMany(ctx context.Context, ids []int) ([]*models.User, error) {
	ret := make([]*models.User, len(ids))

	table := qb.table()
	q := qb.selectDataset().Prepared(true).Where(table.Col(idColumn).In(ids))
	unsorted, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	for _, s := range unsorted {
		i := slices.Index(ids, s.ID)
		ret[i] = s
	}

	for i := range ret {
		if ret[i] == nil {
			return nil, fmt.Errorf("scene marker with id %d not found", ids[i])
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *UserStore) find(ctx context.Context, id int) (*models.User, error) {
	q := qb.selectDataset().Where(userTableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *UserStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.User, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *UserStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.User, error) {
	const single = false
	var rows userQueryRows
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f userQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		rows = append(rows, f)
		return nil
	}); err != nil {
		return nil, err
	}

	return rows.resolve(), nil
}

// func (qb *UserStore) Query(ctx context.Context, userFilter *models.UserFilterType, findFilter *models.FindFilterType) ([]*models.User, int, error) {
// 	query, err := qb.makeQuery(ctx, userFilter, findFilter)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	idsResult, countResult, err := query.executeFind(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	users, err := qb.FindMany(ctx, idsResult)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return users, countResult, nil
// }

// func (qb *UserStore) QueryCount(ctx context.Context, userFilter *models.UserFilterType, findFilter *models.FindFilterType) (int, error) {
// 	query, err := qb.makeQuery(ctx, userFilter, findFilter)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return query.executeCount(ctx)
// }

// var userSortOptions = sortOptions{
// 	"created_at",
// 	"id",
// 	"username",
// 	"updated_at",
// }

// func (qb *UserStore) setUserSort(query *queryBuilder, findFilter *models.FindFilterType) error {
// 	sort := findFilter.GetSort("username")
// 	direction := findFilter.GetDirection()

// 	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
// 	if err := userSortOptions.validateSort(sort); err != nil {
// 		return err
// 	}

// 	query.sortAndPagination += getSort(sort, direction, userTable)

// 	return nil
// }

// func (qb *UserStore) queryUsers(ctx context.Context, query string, args []interface{}) ([]*models.User, error) {
// 	const single = false
// 	var ret []*models.User
// 	if err := userRepository.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
// 		var f userRow
// 		if err := r.StructScan(&f); err != nil {
// 			return err
// 		}

// 		s := f.resolve()

// 		ret = append(ret, s)
// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}

// 	return ret, nil
// }

func (qb *UserStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *UserStore) All(ctx context.Context) ([]*models.User, error) {
	return qb.getMany(ctx, qb.selectDataset())
}
