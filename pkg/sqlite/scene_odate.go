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

const sceneODateTable = "scenes_odates"

type sceneODateRow struct {
	ID      int       `db:"id" goqu:"skipinsert"`
	SceneID zero.Int  `db:"scene_id,omitempty"` // TODO: make schema non-nullable
	ODate   Timestamp `db:"ODate"`
}

func (r *sceneODateRow) fromSceneODate(o models.SceneODate) {
	r.ID = o.ID
	r.SceneID = zero.IntFrom(int64(o.SceneID))
	r.ODate = Timestamp{Timestamp: o.ODate}
}

func (r *sceneODateRow) resolve() *models.SceneODate {
	ret := &models.SceneODate{
		ID:      r.ID,
		SceneID: int(r.SceneID.Int64),
		ODate:   r.ODate.Timestamp,
	}

	return ret
}

type SceneODateStore struct {
	repository

	tableMgr *table
	oDateManager
}

func NewSceneODateStore() *SceneODateStore {
	return &SceneODateStore{
		repository: repository{
			tableName: sceneODateTable,
			idColumn:  idColumn,
		},
		tableMgr:     sceneODateTableMgr,
		oDateManager: oDateManager{sceneODateTableMgr},
	}
}

func (qb *SceneODateStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *SceneODateStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *SceneODateStore) Create(ctx context.Context, newObject *models.SceneODate) error {
	var r sceneODateRow
	r.fromSceneODate(*newObject)

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

func (qb *SceneODateStore) Update(ctx context.Context, updatedObject *models.SceneODate) error {
	var r sceneODateRow
	r.fromSceneODate(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *SceneODateStore) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *SceneODateStore) Find(ctx context.Context, id int) (*models.SceneODate, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *SceneODateStore) FindMany(ctx context.Context, ids []int) ([]*models.SceneODate, error) {
	ret := make([]*models.SceneODate, len(ids))

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
func (qb *SceneODateStore) find(ctx context.Context, id int) (*models.SceneODate, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *SceneODateStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.SceneODate, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *SceneODateStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.SceneODate, error) {
	const single = false
	var ret []*models.SceneODate
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f sceneODateRow
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

func (qb *SceneODateStore) FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneODate, error) {
	query := `
		SELECT scenes_odates.* FROM scenes_odates
		WHERE scenes_odates.scene_id = ?
		GROUP BY scenes_odates.id
	`
	args := []interface{}{sceneID}
	return qb.querysceneODates(ctx, query, args)
}

func sceneODatePerformersCriterionHandler(qb *SceneODateStore, performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    performersScenesTable,
		joinAs:       "performers_join",
		primaryFK:    sceneIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			f.addLeftJoin(performersScenesTable, "performers_join", "performers_join.scene_id = scenes_odates.scene_id")
		},
	}

	handler := h.handler(performers)
	return func(ctx context.Context, f *filterBuilder) {
		// Make sure scenes is included, otherwise excludes filter fails
		f.addLeftJoin(sceneTable, "", "scenes.id = scenes_odates.scene_id")
		handler(ctx, f)
	}
}

func (qb *SceneODateStore) querysceneODates(ctx context.Context, query string, args []interface{}) ([]*models.SceneODate, error) {
	const single = false
	var ret []*models.SceneODate
	if err := qb.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f sceneODateRow
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

func (qb *SceneODateStore) queryMarkerStringsResultType(ctx context.Context, query string, args []interface{}) ([]*models.MarkerStringsResultType, error) {
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

func (qb *SceneODateStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *SceneODateStore) All(ctx context.Context) ([]*models.SceneODate, error) {
	return qb.getMany(ctx, qb.selectDataset())
}
