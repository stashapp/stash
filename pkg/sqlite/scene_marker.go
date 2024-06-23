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

const sceneMarkerTable = "scene_markers"

const countSceneMarkersForTagQuery = `
SELECT scene_markers.id FROM scene_markers
LEFT JOIN scene_markers_tags as tags_join on tags_join.scene_marker_id = scene_markers.id
WHERE tags_join.tag_id = ? OR scene_markers.primary_tag_id = ?
GROUP BY scene_markers.id
`

type sceneMarkerRow struct {
	ID           int       `db:"id" goqu:"skipinsert"`
	Title        string    `db:"title"` // TODO: make db schema (and gql schema) nullable
	Seconds      float64   `db:"seconds"`
	PrimaryTagID int       `db:"primary_tag_id"`
	SceneID      int       `db:"scene_id"`
	CreatedAt    Timestamp `db:"created_at"`
	UpdatedAt    Timestamp `db:"updated_at"`
}

func (r *sceneMarkerRow) fromSceneMarker(o models.SceneMarker) {
	r.ID = o.ID
	r.Title = o.Title
	r.Seconds = o.Seconds
	r.PrimaryTagID = o.PrimaryTagID
	r.SceneID = o.SceneID
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

func (r *sceneMarkerRow) resolve() *models.SceneMarker {
	ret := &models.SceneMarker{
		ID:           r.ID,
		Title:        r.Title,
		Seconds:      r.Seconds,
		PrimaryTagID: r.PrimaryTagID,
		SceneID:      r.SceneID,
		CreatedAt:    r.CreatedAt.Timestamp,
		UpdatedAt:    r.UpdatedAt.Timestamp,
	}

	return ret
}

type sceneMarkerRowRecord struct {
	updateRecord
}

func (r *sceneMarkerRowRecord) fromPartial(o models.SceneMarkerPartial) {
	// TODO: replace with setNullString after schema is made nullable
	// r.setNullString("title", o.Title)
	// saves a null input as the empty string
	if o.Title.Set {
		r.set("title", o.Title.Value)
	}
	r.setFloat64("seconds", o.Seconds)
	r.setInt("primary_tag_id", o.PrimaryTagID)
	r.setInt("scene_id", o.SceneID)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
}

type sceneMarkerRepositoryType struct {
	repository

	scenes repository
	tags   joinRepository
}

var (
	sceneMarkerRepository = sceneMarkerRepositoryType{
		repository: repository{
			tableName: sceneMarkerTable,
			idColumn:  idColumn,
		},
		scenes: repository{
			tableName: sceneTable,
			idColumn:  idColumn,
		},
		tags: joinRepository{
			repository: repository{
				tableName: "scene_markers_tags",
				idColumn:  "scene_marker_id",
			},
			fkColumn: tagIDColumn,
		},
	}
)

type SceneMarkerStore struct{}

func NewSceneMarkerStore() *SceneMarkerStore {
	return &SceneMarkerStore{}
}

func (qb *SceneMarkerStore) table() exp.IdentifierExpression {
	return sceneMarkerTableMgr.table
}

func (qb *SceneMarkerStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *SceneMarkerStore) Create(ctx context.Context, newObject *models.SceneMarker) error {
	var r sceneMarkerRow
	r.fromSceneMarker(*newObject)

	id, err := sceneMarkerTableMgr.insertID(ctx, r)
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

func (qb *SceneMarkerStore) UpdatePartial(ctx context.Context, id int, partial models.SceneMarkerPartial) (*models.SceneMarker, error) {
	r := sceneMarkerRowRecord{
		updateRecord{
			Record: make(exp.Record),
		},
	}

	r.fromPartial(partial)

	if len(r.Record) > 0 {
		if err := sceneMarkerTableMgr.updateByID(ctx, id, r.Record); err != nil {
			return nil, err
		}
	}

	return qb.find(ctx, id)
}

func (qb *SceneMarkerStore) Update(ctx context.Context, updatedObject *models.SceneMarker) error {
	var r sceneMarkerRow
	r.fromSceneMarker(*updatedObject)

	if err := sceneMarkerTableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *SceneMarkerStore) Destroy(ctx context.Context, id int) error {
	return sceneMarkerRepository.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *SceneMarkerStore) Find(ctx context.Context, id int) (*models.SceneMarker, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *SceneMarkerStore) FindMany(ctx context.Context, ids []int) ([]*models.SceneMarker, error) {
	ret := make([]*models.SceneMarker, len(ids))

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
			return nil, fmt.Errorf("scene marker with id %d not found", ids[i])
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *SceneMarkerStore) find(ctx context.Context, id int) (*models.SceneMarker, error) {
	q := qb.selectDataset().Where(sceneMarkerTableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *SceneMarkerStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.SceneMarker, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *SceneMarkerStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.SceneMarker, error) {
	const single = false
	var ret []*models.SceneMarker
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f sceneMarkerRow
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

func (qb *SceneMarkerStore) FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneMarker, error) {
	query := `
		SELECT scene_markers.* FROM scene_markers
		WHERE scene_markers.scene_id = ?
		GROUP BY scene_markers.id
		ORDER BY scene_markers.seconds ASC
	`
	args := []interface{}{sceneID}
	return qb.querySceneMarkers(ctx, query, args)
}

func (qb *SceneMarkerStore) CountByTagID(ctx context.Context, tagID int) (int, error) {
	args := []interface{}{tagID, tagID}
	return sceneMarkerRepository.runCountQuery(ctx, sceneMarkerRepository.buildCountQuery(countSceneMarkersForTagQuery), args)
}

func (qb *SceneMarkerStore) GetMarkerStrings(ctx context.Context, q *string, sort *string) ([]*models.MarkerStringsResultType, error) {
	query := "SELECT count(*) as `count`, scene_markers.id as id, scene_markers.title as title FROM scene_markers"
	if q != nil {
		query += " WHERE title LIKE '%" + *q + "%'"
	}
	query += " GROUP BY title"
	if sort != nil && *sort == "count" {
		query += " ORDER BY `count` DESC"
	} else {
		query += " ORDER BY title ASC"
	}
	var args []interface{}
	return qb.queryMarkerStringsResultType(ctx, query, args)
}

func (qb *SceneMarkerStore) Wall(ctx context.Context, q *string) ([]*models.SceneMarker, error) {
	s := ""
	if q != nil {
		s = *q
	}

	table := qb.table()
	qq := qb.selectDataset().Prepared(true).Where(table.Col("title").Like("%" + s + "%")).Order(goqu.L("RANDOM()").Asc()).Limit(80)
	return qb.getMany(ctx, qq)
}

func (qb *SceneMarkerStore) makeQuery(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if sceneMarkerFilter == nil {
		sceneMarkerFilter = &models.SceneMarkerFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := sceneMarkerRepository.newQuery()
	distinctIDs(&query, sceneMarkerTable)

	if q := findFilter.Q; q != nil && *q != "" {
		query.join(tagTable, "", "scene_markers.primary_tag_id = tags.id")
		searchColumns := []string{"scene_markers.title", "scenes.title", "tags.name"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := filterBuilderFromHandler(ctx, &sceneMarkerFilterHandler{
		sceneMarkerFilter: sceneMarkerFilter,
	})

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	if err := qb.setSceneMarkerSort(&query, findFilter); err != nil {
		return nil, err
	}
	query.sortAndPagination += getPagination(findFilter)

	return &query, nil
}

func (qb *SceneMarkerStore) Query(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) ([]*models.SceneMarker, int, error) {
	query, err := qb.makeQuery(ctx, sceneMarkerFilter, findFilter)
	if err != nil {
		return nil, 0, err
	}

	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	sceneMarkers, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return sceneMarkers, countResult, nil
}

func (qb *SceneMarkerStore) QueryCount(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, sceneMarkerFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

var sceneMarkerSortOptions = sortOptions{
	"created_at",
	"id",
	"title",
	"random",
	"scene_id",
	"scenes_updated_at",
	"seconds",
	"updated_at",
}

func (qb *SceneMarkerStore) setSceneMarkerSort(query *queryBuilder, findFilter *models.FindFilterType) error {
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()

	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
	if err := sceneMarkerSortOptions.validateSort(sort); err != nil {
		return err
	}

	switch sort {
	case "scenes_updated_at":
		sort = "updated_at"
		query.join(sceneTable, "", "scenes.id = scene_markers.scene_id")
		query.sortAndPagination += getSort(sort, direction, sceneTable)
	case "title":
		query.join(tagTable, "", "scene_markers.primary_tag_id = tags.id")
		query.sortAndPagination += " ORDER BY COALESCE(NULLIF(scene_markers.title,''), tags.name) COLLATE NATURAL_CI " + direction
	default:
		query.sortAndPagination += getSort(sort, direction, sceneMarkerTable)
	}

	query.sortAndPagination += ", scene_markers.scene_id ASC, scene_markers.seconds ASC"
	return nil
}

func (qb *SceneMarkerStore) querySceneMarkers(ctx context.Context, query string, args []interface{}) ([]*models.SceneMarker, error) {
	const single = false
	var ret []*models.SceneMarker
	if err := sceneMarkerRepository.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f sceneMarkerRow
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

func (qb *SceneMarkerStore) queryMarkerStringsResultType(ctx context.Context, query string, args []interface{}) ([]*models.MarkerStringsResultType, error) {
	rows, err := dbWrapper.Queryx(ctx, query, args...)
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

func (qb *SceneMarkerStore) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return sceneMarkerRepository.tags.getIDs(ctx, id)
}

func (qb *SceneMarkerStore) UpdateTags(ctx context.Context, id int, tagIDs []int) error {
	// Delete the existing joins and then create new ones
	return sceneMarkerRepository.tags.replace(ctx, id, tagIDs)
}

func (qb *SceneMarkerStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *SceneMarkerStore) All(ctx context.Context) ([]*models.SceneMarker, error) {
	return qb.getMany(ctx, qb.selectDataset())
}
