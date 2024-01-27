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
	"github.com/stashapp/stash/pkg/sliceutil"
)

const sceneFilterTable = "scene_filters"

type sceneFilterRow struct {
	ID          int       `db:"id" goqu:"skipinsert"`
	Contrast    int       `db:"contrast"`
	Brightness  int       `db:"brightness"`
	Gamma       int       `db:"gamma"`
	Saturate    int       `db:"saturate"`
	HueRotate   int       `db:"hue_rotate"`
	Warmth      int       `db:"warmth"`
	Red         int       `db:"red"`
	Green       int       `db:"green"`
	Blue        int       `db:"blue"`
	Blur        int       `db:"blur"`
	Rotate      float64   `db:"rotate"`
	Scale       int       `db:"scale"`
	AspectRatio int       `db:"aspect_ratio"`
	SceneID     int       `db:"scene_id"`
	CreatedAt   Timestamp `db:"created_at"`
	UpdatedAt   Timestamp `db:"updated_at"`
}

func (r *sceneFilterRow) fromSceneFilter(o models.SceneFilter) {
	r.ID = o.ID
	r.SceneID = o.SceneID
	r.Contrast = o.Contrast
	r.Brightness = o.Brightness
	r.Gamma = o.Gamma
	r.Saturate = o.Saturate
	r.HueRotate = o.HueRotate
	r.Warmth = o.Warmth
	r.Red = o.Red
	r.Green = o.Green
	r.Blue = o.Blue
	r.Blur = o.Blur
	r.Rotate = o.Rotate
	r.Scale = o.Scale
	r.AspectRatio = o.AspectRatio
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

func (r *sceneFilterRow) resolve() *models.SceneFilter {
	ret := &models.SceneFilter{
		ID:          r.ID,
		SceneID:     r.SceneID,
		Contrast:    r.Contrast,
		Brightness:  r.Brightness,
		Gamma:       r.Gamma,
		Saturate:    r.Saturate,
		HueRotate:   r.HueRotate,
		Warmth:      r.Warmth,
		Red:         r.Red,
		Green:       r.Green,
		Blue:        r.Blue,
		Blur:        r.Blur,
		Rotate:      r.Rotate,
		Scale:       r.Scale,
		AspectRatio: r.AspectRatio,
		CreatedAt:   r.CreatedAt.Timestamp,
		UpdatedAt:   r.UpdatedAt.Timestamp,
	}

	return ret
}

type SceneFilterStore struct {
	repository

	tableMgr *table
}

func NewSceneFilterStore() *SceneFilterStore {
	return &SceneFilterStore{
		repository: repository{
			tableName: sceneFilterTable,
			idColumn:  idColumn,
		},
		tableMgr: sceneFilterTableMgr,
	}
}

func (qb *SceneFilterStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *SceneFilterStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *SceneFilterStore) Create(ctx context.Context, newObject *models.SceneFilter) error {
	var r sceneFilterRow
	r.fromSceneFilter(*newObject)

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

func (qb *SceneFilterStore) Update(ctx context.Context, updatedObject *models.SceneFilter) error {
	var r sceneFilterRow
	r.fromSceneFilter(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *SceneFilterStore) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

func (qb *SceneFilterStore) Find(ctx context.Context, id int) (*models.SceneFilter, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *SceneFilterStore) FindMany(ctx context.Context, ids []int) ([]*models.SceneFilter, error) {
	ret := make([]*models.SceneFilter, len(ids))

	table := qb.table()
	q := qb.selectDataset().Prepared(true).Where(table.Col(idColumn).In(ids))
	unsorted, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	for _, s := range unsorted {
		i := sliceutil.Index(ids, s.ID)
		ret[i] = s
	}

	for i := range ret {
		if ret[i] == nil {
			return nil, fmt.Errorf("scene filter with id %d not found", ids[i])
		}
	}

	return ret, nil
}

func (qb *SceneFilterStore) find(ctx context.Context, id int) (*models.SceneFilter, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *SceneFilterStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.SceneFilter, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *SceneFilterStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.SceneFilter, error) {
	const single = false
	var ret []*models.SceneFilter
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f sceneFilterRow
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

func (qb *SceneFilterStore) FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneFilter, error) {
	query := `
		SELECT scene_filters.* FROM scene_filters
		WHERE scene_filters.scene_id = ?
		GROUP BY scene_filters.id
		ORDER BY scene_filters.id ASC
	`
	args := []interface{}{sceneID}
	return qb.querySceneFilters(ctx, query, args)
}

func (qb *SceneFilterStore) makeFilter(ctx context.Context, sceneFilterFilter *models.SceneFilterFilterType) *filterBuilder {
	query := &filterBuilder{}

	query.handleCriterion(ctx, timestampCriterionHandler(sceneFilterFilter.CreatedAt, "scene_filters.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(sceneFilterFilter.UpdatedAt, "scene_filters.updated_at"))
	query.handleCriterion(ctx, dateCriterionHandler(sceneFilterFilter.SceneDate, "scenes.date"))
	query.handleCriterion(ctx, timestampCriterionHandler(sceneFilterFilter.SceneCreatedAt, "scenes.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(sceneFilterFilter.SceneUpdatedAt, "scenes.updated_at"))

	return query
}

func (qb *SceneFilterStore) makeQuery(ctx context.Context, sceneFilterFilter *models.SceneFilterFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if sceneFilterFilter == nil {
		sceneFilterFilter = &models.SceneFilterFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, sceneFilterTable)

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"scenes.title"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := qb.makeFilter(ctx, sceneFilterFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	query.sortAndPagination = qb.getSceneFilterSort(&query, findFilter) + getPagination(findFilter)

	return &query, nil
}

func (qb *SceneFilterStore) Query(ctx context.Context, sceneFilterFilter *models.SceneFilterFilterType, findFilter *models.FindFilterType) ([]*models.SceneFilter, int, error) {
	query, err := qb.makeQuery(ctx, sceneFilterFilter, findFilter)
	if err != nil {
		return nil, 0, err
	}

	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	sceneFilters, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return sceneFilters, countResult, nil
}

func (qb *SceneFilterStore) QueryCount(ctx context.Context, sceneFilterFilter *models.SceneFilterFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, sceneFilterFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

func (qb *SceneFilterStore) getSceneFilterSort(query *queryBuilder, findFilter *models.FindFilterType) string {
	sort := findFilter.GetSort("id")
	direction := findFilter.GetDirection()
	tableName := "scene_filters"
	if sort == "scenes_updated_at" {
		// ensure scene table is joined
		query.join(sceneTable, "", "scenes.id = scene_filters.scene_id")
		sort = "updated_at"
		tableName = "scenes"
	}

	additional := ", scene_filters.scene_id ASC"
	return getSort(sort, direction, tableName) + additional
}

func (qb *SceneFilterStore) querySceneFilters(ctx context.Context, query string, args []interface{}) ([]*models.SceneFilter, error) {
	const single = false
	var ret []*models.SceneFilter
	if err := qb.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f sceneFilterRow
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

func (qb *SceneFilterStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *SceneFilterStore) All(ctx context.Context) ([]*models.SceneFilter, error) {
	return qb.getMany(ctx, qb.selectDataset())
}
