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

	"github.com/stashapp/stash/pkg/models"
)

const (
	sceneSegmentTable   = "scene_segments"
	sceneSegmentIDColumn = "scene_segment_id"
)

type sceneSegmentRow struct {
	ID           int       `db:"id" goqu:"skipinsert"`
	SceneID      int       `db:"scene_id"`
	Title        string    `db:"title"`
	StartSeconds float64   `db:"start_seconds"`
	EndSeconds   float64   `db:"end_seconds"`
	CreatedAt    Timestamp `db:"created_at"`
	UpdatedAt    Timestamp `db:"updated_at"`
}

func (r *sceneSegmentRow) fromSceneSegment(o models.SceneSegment) {
	r.ID = o.ID
	r.SceneID = o.SceneID
	r.Title = o.Title
	r.StartSeconds = o.StartSeconds
	r.EndSeconds = o.EndSeconds
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

func (r *sceneSegmentRow) resolve() *models.SceneSegment {
	ret := &models.SceneSegment{
		ID:           r.ID,
		SceneID:      r.SceneID,
		Title:        r.Title,
		StartSeconds: r.StartSeconds,
		EndSeconds:   r.EndSeconds,
		CreatedAt:    r.CreatedAt.Timestamp,
		UpdatedAt:    r.UpdatedAt.Timestamp,
	}

	return ret
}

type sceneSegmentRowRecord struct {
	updateRecord
}

func (r *sceneSegmentRowRecord) fromPartial(o models.SceneSegmentPartial) {
	r.setInt("scene_id", o.SceneID)
	r.setString("title", o.Title)
	r.setFloat64("start_seconds", o.StartSeconds)
	r.setFloat64("end_seconds", o.EndSeconds)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
}

type sceneSegmentRepositoryType struct {
	repository
	scenes repository
}

var (
	sceneSegmentRepository = sceneSegmentRepositoryType{
		repository: repository{
			tableName: sceneSegmentTable,
			idColumn:  idColumn,
		},
		scenes: repository{
			tableName: sceneTable,
			idColumn:  idColumn,
		},
	}
)

type SceneSegmentStore struct{}

func NewSceneSegmentStore() *SceneSegmentStore {
	return &SceneSegmentStore{}
}

func (qb *SceneSegmentStore) table() exp.IdentifierExpression {
	return sceneSegmentTableMgr.table
}

func (qb *SceneSegmentStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *SceneSegmentStore) Create(ctx context.Context, newObject *models.SceneSegment) error {
	var r sceneSegmentRow
	r.fromSceneSegment(*newObject)

	id, err := sceneSegmentTableMgr.insertID(ctx, r)
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

func (qb *SceneSegmentStore) UpdatePartial(ctx context.Context, id int, partial models.SceneSegmentPartial) (*models.SceneSegment, error) {
	r := sceneSegmentRowRecord{
		updateRecord{
			Record: make(exp.Record),
		},
	}

	r.fromPartial(partial)

	if len(r.Record) > 0 {
		if err := sceneSegmentTableMgr.updateByID(ctx, id, r.Record); err != nil {
			return nil, err
		}
	}

	return qb.find(ctx, id)
}

func (qb *SceneSegmentStore) Update(ctx context.Context, updatedObject *models.SceneSegment) error {
	var r sceneSegmentRow
	r.fromSceneSegment(*updatedObject)

	if err := sceneSegmentTableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *SceneSegmentStore) Destroy(ctx context.Context, id int) error {
	return sceneSegmentRepository.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *SceneSegmentStore) Find(ctx context.Context, id int) (*models.SceneSegment, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *SceneSegmentStore) FindMany(ctx context.Context, ids []int) ([]*models.SceneSegment, error) {
	ret := make([]*models.SceneSegment, len(ids))

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
			return nil, fmt.Errorf("scene segment with id %d not found", ids[i])
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *SceneSegmentStore) find(ctx context.Context, id int) (*models.SceneSegment, error) {
	q := qb.selectDataset().Where(sceneSegmentTableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *SceneSegmentStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.SceneSegment, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *SceneSegmentStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.SceneSegment, error) {
	const single = false
	var ret []*models.SceneSegment
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f sceneSegmentRow
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

func (qb *SceneSegmentStore) FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneSegment, error) {
	table := qb.table()
	q := qb.selectDataset().
		Prepared(true).
		Where(table.Col("scene_id").Eq(sceneID)).
		Order(table.Col("start_seconds").Asc())
	
	return qb.getMany(ctx, q)
}

func (qb *SceneSegmentStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *SceneSegmentStore) All(ctx context.Context) ([]*models.SceneSegment, error) {
	return qb.getMany(ctx, qb.selectDataset())
}
