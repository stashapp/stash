package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

const scenePlayDateTable = "scenes_playdates"

type scenePlayDateRow struct {
	ID       int       `db:"id" goqu:"skipinsert"`
	SceneID  zero.Int  `db:"scene_id,omitempty"` // TODO: make schema non-nullable
	PlayDate Timestamp `db:"playdate"`
}

func (r *scenePlayDateRow) fromScenePlayDate(o models.ScenePlayDate) {
	r.ID = o.ID
	r.SceneID = zero.IntFrom(int64(o.SceneID))
	r.PlayDate = Timestamp{Timestamp: o.PlayDate}
}

func (r *scenePlayDateRow) resolve() *models.ScenePlayDate {
	ret := &models.ScenePlayDate{
		ID:       r.ID,
		SceneID:  int(r.SceneID.Int64),
		PlayDate: r.PlayDate.Timestamp,
	}

	return ret
}

type ScenePlayDateStore struct {
	repository

	tableMgr *table
	playDateManager
}

func NewScenePlayDateStore() *ScenePlayDateStore {
	return &ScenePlayDateStore{
		repository: repository{
			tableName: scenePlayDateTable,
			idColumn:  idColumn,
		},
		tableMgr:        scenePlayDateTableMgr,
		playDateManager: playDateManager{scenePlayDateTableMgr},
	}
}

func (qb *ScenePlayDateStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *ScenePlayDateStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *ScenePlayDateStore) Create(ctx context.Context, newObject *models.ScenePlayDate) error {
	var r scenePlayDateRow
	r.fromScenePlayDate(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	updated, err := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject = *updated

	return nil
}

func (qb *ScenePlayDateStore) Update(ctx context.Context, updatedObject *models.ScenePlayDate) error {
	var r scenePlayDateRow
	r.fromScenePlayDate(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *ScenePlayDateStore) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *ScenePlayDateStore) Find(ctx context.Context, id int) (*models.ScenePlayDate, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *ScenePlayDateStore) FindMany(ctx context.Context, ids []int) ([]*models.ScenePlayDate, error) {
	ret := make([]*models.ScenePlayDate, len(ids))

	table := qb.table()
	q := qb.selectDataset().Prepared(true).Where(table.Col(idColumn).In(ids))
	unsorted, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	for _, s := range unsorted {
		i := intslice.IntIndex(ids, s.ID)
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
func (qb *ScenePlayDateStore) find(ctx context.Context, id int) (*models.ScenePlayDate, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *ScenePlayDateStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.ScenePlayDate, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *ScenePlayDateStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.ScenePlayDate, error) {
	const single = false
	var ret []*models.ScenePlayDate
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f scenePlayDateRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *ScenePlayDateStore) FindBySceneID(ctx context.Context, sceneID int) ([]*models.ScenePlayDate, error) {
	query := `
		SELECT scenes_playdates.* FROM scenes_playdates
		WHERE scenes_playdates.scene_id = ?
		GROUP BY scenes_playdates.id
	`
	args := []interface{}{sceneID}
	return qb.queryscenePlayDates(ctx, query, args)
}

func scenePlayDatePerformersCriterionHandler(qb *ScenePlayDateStore, performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    performersScenesTable,
		joinAs:       "performers_join",
		primaryFK:    sceneIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			f.addLeftJoin(performersScenesTable, "performers_join", "performers_join.scene_id = scenes_playdates.scene_id")
		},
	}

	handler := h.handler(performers)
	return func(ctx context.Context, f *filterBuilder) {
		// Make sure scenes is included, otherwise excludes filter fails
		f.addLeftJoin(sceneTable, "", "scenes.id = scenes_playdates.scene_id")
		handler(ctx, f)
	}
}

func (qb *ScenePlayDateStore) queryscenePlayDates(ctx context.Context, query string, args []interface{}) ([]*models.ScenePlayDate, error) {
	const single = false
	var ret []*models.ScenePlayDate
	if err := qb.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f scenePlayDateRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *ScenePlayDateStore) queryMarkerStringsResultType(ctx context.Context, query string, args []interface{}) ([]*models.MarkerStringsResultType, error) {
	rows, err := qb.tx.Queryx(ctx, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()

	markerStrings := make([]*models.MarkerStringsResultType, 0)
	for rows.Next() {
		markerString := models.MarkerStringsResultType{}
		if err := rows.StructScan(&markerString); err != nil {
			return nil, err
		}
		markerStrings = append(markerStrings, &markerString)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return markerStrings, nil
}

func (qb *ScenePlayDateStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *ScenePlayDateStore) All(ctx context.Context) ([]*models.ScenePlayDate, error) {
	return qb.getMany(ctx, qb.selectDataset())
}
