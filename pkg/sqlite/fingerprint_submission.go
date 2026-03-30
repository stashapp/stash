package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/models"
)

const (
	fingerprintSubmissionsTable = "fingerprint_submissions"
)

var (
	fingerprintSubmissionsTableMgr = &table{
		table:    goqu.T(fingerprintSubmissionsTable),
		idColumn: goqu.T(fingerprintSubmissionsTable).Col("endpoint"), // not a real ID column, but needed for table struct
	}
)

type fingerprintSubmissionRow struct {
	Endpoint  string    `db:"endpoint"`
	StashID   string    `db:"stash_id"`
	SceneID   int       `db:"scene_id"`
	Vote      string    `db:"vote"`
	CreatedAt Timestamp `db:"created_at"`
}

func (r *fingerprintSubmissionRow) fromFingerprintSubmission(o models.FingerprintSubmission) {
	r.Endpoint = o.Endpoint
	r.StashID = o.StashID
	r.SceneID = o.SceneID
	r.Vote = string(o.Vote)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
}

func (r *fingerprintSubmissionRow) resolve() *models.FingerprintSubmission {
	return &models.FingerprintSubmission{
		Endpoint:  r.Endpoint,
		StashID:   r.StashID,
		SceneID:   r.SceneID,
		Vote:      models.FingerprintVote(r.Vote),
		CreatedAt: r.CreatedAt.Timestamp,
	}
}

type FingerprintSubmissionStore struct{}

func NewFingerprintSubmissionStore() *FingerprintSubmissionStore {
	return &FingerprintSubmissionStore{}
}

func (qb *FingerprintSubmissionStore) table() exp.IdentifierExpression {
	return fingerprintSubmissionsTableMgr.table
}

func (qb *FingerprintSubmissionStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *FingerprintSubmissionStore) Create(ctx context.Context, newObject *models.FingerprintSubmission) error {
	var r fingerprintSubmissionRow
	r.fromFingerprintSubmission(*newObject)

	q := dialect.Insert(qb.table()).Prepared(true).Rows(r).OnConflict(goqu.DoNothing())
	if _, err := exec(ctx, q); err != nil {
		return err
	}

	return nil
}

func (qb *FingerprintSubmissionStore) Delete(ctx context.Context, endpoint string, stashID string) error {
	q := dialect.Delete(qb.table()).Where(
		qb.table().Col("endpoint").Eq(endpoint),
		qb.table().Col("stash_id").Eq(stashID),
	)

	if _, err := exec(ctx, q); err != nil {
		return err
	}

	return nil
}

func (qb *FingerprintSubmissionStore) DeleteByEndpoint(ctx context.Context, endpoint string) error {
	q := dialect.Delete(qb.table()).Where(
		qb.table().Col("endpoint").Eq(endpoint),
	)

	if _, err := exec(ctx, q); err != nil {
		return err
	}

	return nil
}

func (qb *FingerprintSubmissionStore) Find(ctx context.Context, endpoint string, stashID string) (*models.FingerprintSubmission, error) {
	q := qb.selectDataset().Where(
		qb.table().Col("endpoint").Eq(endpoint),
		qb.table().Col("stash_id").Eq(stashID),
	)

	ret, err := qb.get(ctx, q)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *FingerprintSubmissionStore) FindByEndpoint(ctx context.Context, endpoint string) ([]*models.FingerprintSubmission, error) {
	q := qb.selectDataset().Where(
		qb.table().Col("endpoint").Eq(endpoint),
	).Order(qb.table().Col("created_at").Asc())

	return qb.getMany(ctx, q)
}

func (qb *FingerprintSubmissionStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.FingerprintSubmission, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *FingerprintSubmissionStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.FingerprintSubmission, error) {
	const single = false
	var ret []*models.FingerprintSubmission
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f fingerprintSubmissionRow
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
