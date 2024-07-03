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
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

const (
	groupTable    = "movies"
	groupIDColumn = "movie_id"

	groupFrontImageBlobColumn = "front_image_blob"
	groupBackImageBlobColumn  = "back_image_blob"

	groupsTagsTable = "movies_tags"

	groupURLsTable = "movie_urls"
	groupURLColumn = "url"
)

type groupRow struct {
	ID       int         `db:"id" goqu:"skipinsert"`
	Name     zero.String `db:"name"`
	Aliases  zero.String `db:"aliases"`
	Duration null.Int    `db:"duration"`
	Date     NullDate    `db:"date"`
	// expressed as 1-100
	Rating    null.Int    `db:"rating"`
	StudioID  null.Int    `db:"studio_id,omitempty"`
	Director  zero.String `db:"director"`
	Synopsis  zero.String `db:"synopsis"`
	CreatedAt Timestamp   `db:"created_at"`
	UpdatedAt Timestamp   `db:"updated_at"`

	// not used in resolutions or updates
	FrontImageBlob zero.String `db:"front_image_blob"`
	BackImageBlob  zero.String `db:"back_image_blob"`
}

func (r *groupRow) fromGroup(o models.Group) {
	r.ID = o.ID
	r.Name = zero.StringFrom(o.Name)
	r.Aliases = zero.StringFrom(o.Aliases)
	r.Duration = intFromPtr(o.Duration)
	r.Date = NullDateFromDatePtr(o.Date)
	r.Rating = intFromPtr(o.Rating)
	r.StudioID = intFromPtr(o.StudioID)
	r.Director = zero.StringFrom(o.Director)
	r.Synopsis = zero.StringFrom(o.Synopsis)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

func (r *groupRow) resolve() *models.Group {
	ret := &models.Group{
		ID:        r.ID,
		Name:      r.Name.String,
		Aliases:   r.Aliases.String,
		Duration:  nullIntPtr(r.Duration),
		Date:      r.Date.DatePtr(),
		Rating:    nullIntPtr(r.Rating),
		StudioID:  nullIntPtr(r.StudioID),
		Director:  r.Director.String,
		Synopsis:  r.Synopsis.String,
		CreatedAt: r.CreatedAt.Timestamp,
		UpdatedAt: r.UpdatedAt.Timestamp,
	}

	return ret
}

type groupRowRecord struct {
	updateRecord
}

func (r *groupRowRecord) fromPartial(o models.GroupPartial) {
	r.setNullString("name", o.Name)
	r.setNullString("aliases", o.Aliases)
	r.setNullInt("duration", o.Duration)
	r.setNullDate("date", o.Date)
	r.setNullInt("rating", o.Rating)
	r.setNullInt("studio_id", o.StudioID)
	r.setNullString("director", o.Director)
	r.setNullString("synopsis", o.Synopsis)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
}

type groupRepositoryType struct {
	repository
	scenes repository
	tags   joinRepository
}

var (
	groupRepository = groupRepositoryType{
		repository: repository{
			tableName: groupTable,
			idColumn:  idColumn,
		},
		scenes: repository{
			tableName: groupsScenesTable,
			idColumn:  groupIDColumn,
		},
		tags: joinRepository{
			repository: repository{
				tableName: groupsTagsTable,
				idColumn:  groupIDColumn,
			},
			fkColumn:     tagIDColumn,
			foreignTable: tagTable,
			orderBy:      "tags.name ASC",
		},
	}
)

type GroupStore struct {
	blobJoinQueryBuilder
	tagRelationshipStore

	tableMgr *table
}

func NewGroupStore(blobStore *BlobStore) *GroupStore {
	return &GroupStore{
		blobJoinQueryBuilder: blobJoinQueryBuilder{
			blobStore: blobStore,
			joinTable: groupTable,
		},
		tagRelationshipStore: tagRelationshipStore{
			idRelationshipStore: idRelationshipStore{
				joinTable: groupsTagsTableMgr,
			},
		},

		tableMgr: groupTableMgr,
	}
}

func (qb *GroupStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *GroupStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *GroupStore) Create(ctx context.Context, newObject *models.Group) error {
	var r groupRow
	r.fromGroup(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if newObject.URLs.Loaded() {
		const startPos = 0
		if err := groupsURLsTableMgr.insertJoins(ctx, id, startPos, newObject.URLs.List()); err != nil {
			return err
		}
	}

	if err := qb.tagRelationshipStore.createRelationships(ctx, id, newObject.TagIDs); err != nil {
		return err
	}

	updated, err := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject = *updated

	return nil
}

func (qb *GroupStore) UpdatePartial(ctx context.Context, id int, partial models.GroupPartial) (*models.Group, error) {
	r := groupRowRecord{
		updateRecord{
			Record: make(exp.Record),
		},
	}

	r.fromPartial(partial)

	if len(r.Record) > 0 {
		if err := qb.tableMgr.updateByID(ctx, id, r.Record); err != nil {
			return nil, err
		}
	}

	if partial.URLs != nil {
		if err := groupsURLsTableMgr.modifyJoins(ctx, id, partial.URLs.Values, partial.URLs.Mode); err != nil {
			return nil, err
		}
	}

	if err := qb.tagRelationshipStore.modifyRelationships(ctx, id, partial.TagIDs); err != nil {
		return nil, err
	}

	return qb.find(ctx, id)
}

func (qb *GroupStore) Update(ctx context.Context, updatedObject *models.Group) error {
	var r groupRow
	r.fromGroup(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if updatedObject.URLs.Loaded() {
		if err := groupsURLsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.URLs.List()); err != nil {
			return err
		}
	}

	if err := qb.tagRelationshipStore.replaceRelationships(ctx, updatedObject.ID, updatedObject.TagIDs); err != nil {
		return err
	}

	return nil
}

func (qb *GroupStore) Destroy(ctx context.Context, id int) error {
	// must handle image checksums manually
	if err := qb.destroyImages(ctx, id); err != nil {
		return err
	}

	return groupRepository.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *GroupStore) Find(ctx context.Context, id int) (*models.Group, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *GroupStore) FindMany(ctx context.Context, ids []int) ([]*models.Group, error) {
	ret := make([]*models.Group, len(ids))

	table := qb.table()
	if err := batchExec(ids, defaultBatchSize, func(batch []int) error {
		q := qb.selectDataset().Prepared(true).Where(table.Col(idColumn).In(batch))
		unsorted, err := qb.getMany(ctx, q)
		if err != nil {
			return err
		}

		for _, s := range unsorted {
			i := sliceutil.Index(ids, s.ID)
			ret[i] = s
		}

		return nil
	}); err != nil {
		return nil, err
	}

	for i := range ret {
		if ret[i] == nil {
			return nil, fmt.Errorf("group with id %d not found", ids[i])
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *GroupStore) find(ctx context.Context, id int) (*models.Group, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *GroupStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Group, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *GroupStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Group, error) {
	const single = false
	var ret []*models.Group
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f groupRow
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

func (qb *GroupStore) FindByName(ctx context.Context, name string, nocase bool) (*models.Group, error) {
	// query := "SELECT * FROM movies WHERE name = ?"
	// if nocase {
	// 	query += " COLLATE NOCASE"
	// }
	// query += " LIMIT 1"
	where := "name = ?"
	if nocase {
		where += " COLLATE NOCASE"
	}
	sq := qb.selectDataset().Prepared(true).Where(goqu.L(where, name)).Limit(1)
	ret, err := qb.get(ctx, sq)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return ret, nil
}

func (qb *GroupStore) FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Group, error) {
	// query := "SELECT * FROM movies WHERE name"
	// if nocase {
	// 	query += " COLLATE NOCASE"
	// }
	// query += " IN " + getInBinding(len(names))
	where := "name"
	if nocase {
		where += " COLLATE NOCASE"
	}
	where += " IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	sq := qb.selectDataset().Prepared(true).Where(goqu.L(where, args...))
	ret, err := qb.getMany(ctx, sq)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *GroupStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *GroupStore) All(ctx context.Context) ([]*models.Group, error) {
	table := qb.table()

	return qb.getMany(ctx, qb.selectDataset().Order(
		table.Col("name").Asc(),
		table.Col(idColumn).Asc(),
	))
}

func (qb *GroupStore) makeQuery(ctx context.Context, groupFilter *models.GroupFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}
	if groupFilter == nil {
		groupFilter = &models.GroupFilterType{}
	}

	query := groupRepository.newQuery()
	distinctIDs(&query, groupTable)

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"movies.name", "movies.aliases"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := filterBuilderFromHandler(ctx, &groupFilterHandler{
		groupFilter: groupFilter,
	})

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	var err error
	query.sortAndPagination, err = qb.getGroupSort(findFilter)
	if err != nil {
		return nil, err
	}

	query.sortAndPagination += getPagination(findFilter)

	return &query, nil
}

func (qb *GroupStore) Query(ctx context.Context, groupFilter *models.GroupFilterType, findFilter *models.FindFilterType) ([]*models.Group, int, error) {
	query, err := qb.makeQuery(ctx, groupFilter, findFilter)
	if err != nil {
		return nil, 0, err
	}

	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	groups, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return groups, countResult, nil
}

func (qb *GroupStore) QueryCount(ctx context.Context, groupFilter *models.GroupFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, groupFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

var groupSortOptions = sortOptions{
	"created_at",
	"date",
	"duration",
	"id",
	"name",
	"random",
	"rating",
	"scenes_count",
	"tag_count",
	"updated_at",
}

func (qb *GroupStore) getGroupSort(findFilter *models.FindFilterType) (string, error) {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}

	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
	if err := groupSortOptions.validateSort(sort); err != nil {
		return "", err
	}

	sortQuery := ""
	switch sort {
	case "tag_count":
		sortQuery += getCountSort(groupTable, groupsTagsTable, groupIDColumn, direction)
	case "scenes_count": // generic getSort won't work for this
		sortQuery += getCountSort(groupTable, groupsScenesTable, groupIDColumn, direction)
	default:
		sortQuery += getSort(sort, direction, "movies")
	}

	// Whatever the sorting, always use name/id as a final sort
	sortQuery += ", COALESCE(movies.name, movies.id) COLLATE NATURAL_CI ASC"
	return sortQuery, nil
}

func (qb *GroupStore) queryGroups(ctx context.Context, query string, args []interface{}) ([]*models.Group, error) {
	const single = false
	var ret []*models.Group
	if err := groupRepository.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f groupRow
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

func (qb *GroupStore) UpdateFrontImage(ctx context.Context, groupID int, frontImage []byte) error {
	return qb.UpdateImage(ctx, groupID, groupFrontImageBlobColumn, frontImage)
}

func (qb *GroupStore) UpdateBackImage(ctx context.Context, groupID int, backImage []byte) error {
	return qb.UpdateImage(ctx, groupID, groupBackImageBlobColumn, backImage)
}

func (qb *GroupStore) destroyImages(ctx context.Context, groupID int) error {
	if err := qb.DestroyImage(ctx, groupID, groupFrontImageBlobColumn); err != nil {
		return err
	}
	if err := qb.DestroyImage(ctx, groupID, groupBackImageBlobColumn); err != nil {
		return err
	}

	return nil
}

func (qb *GroupStore) GetFrontImage(ctx context.Context, groupID int) ([]byte, error) {
	return qb.GetImage(ctx, groupID, groupFrontImageBlobColumn)
}

func (qb *GroupStore) HasFrontImage(ctx context.Context, groupID int) (bool, error) {
	return qb.HasImage(ctx, groupID, groupFrontImageBlobColumn)
}

func (qb *GroupStore) GetBackImage(ctx context.Context, groupID int) ([]byte, error) {
	return qb.GetImage(ctx, groupID, groupBackImageBlobColumn)
}

func (qb *GroupStore) HasBackImage(ctx context.Context, groupID int) (bool, error) {
	return qb.HasImage(ctx, groupID, groupBackImageBlobColumn)
}

func (qb *GroupStore) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Group, error) {
	query := `SELECT DISTINCT movies.*
FROM movies
INNER JOIN movies_scenes ON movies.id = movies_scenes.movie_id
INNER JOIN performers_scenes ON performers_scenes.scene_id = movies_scenes.scene_id
WHERE performers_scenes.performer_id = ?
`
	args := []interface{}{performerID}
	return qb.queryGroups(ctx, query, args)
}

func (qb *GroupStore) CountByPerformerID(ctx context.Context, performerID int) (int, error) {
	query := `SELECT COUNT(DISTINCT movies_scenes.movie_id) AS count
FROM movies_scenes
INNER JOIN performers_scenes ON performers_scenes.scene_id = movies_scenes.scene_id
WHERE performers_scenes.performer_id = ?
`
	args := []interface{}{performerID}
	return groupRepository.runCountQuery(ctx, query, args)
}

func (qb *GroupStore) FindByStudioID(ctx context.Context, studioID int) ([]*models.Group, error) {
	query := `SELECT movies.*
FROM movies
WHERE movies.studio_id = ?
`
	args := []interface{}{studioID}
	return qb.queryGroups(ctx, query, args)
}

func (qb *GroupStore) CountByStudioID(ctx context.Context, studioID int) (int, error) {
	query := `SELECT COUNT(1) AS count
FROM movies
WHERE movies.studio_id = ?
`
	args := []interface{}{studioID}
	return groupRepository.runCountQuery(ctx, query, args)
}

func (qb *GroupStore) GetURLs(ctx context.Context, groupID int) ([]string, error) {
	return groupsURLsTableMgr.get(ctx, groupID)
}
