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
)

const (
	groupTable    = "groups"
	groupIDColumn = "group_id"

	groupFrontImageBlobColumn = "front_image_blob"
	groupBackImageBlobColumn  = "back_image_blob"

	groupsTagsTable = "groups_tags"

	groupURLsTable = "group_urls"
	groupURLColumn = "url"

	groupRelationsTable = "groups_relations"
)

type groupRow struct {
	ID            int         `db:"id" goqu:"skipinsert"`
	Name          zero.String `db:"name"`
	Aliases       zero.String `db:"aliases"`
	Duration      null.Int    `db:"duration"`
	Date          NullDate    `db:"date"`
	DatePrecision null.Int    `db:"date_precision"`
	// expressed as 1-100
	Rating      null.Int    `db:"rating"`
	StudioID    null.Int    `db:"studio_id,omitempty"`
	Director    zero.String `db:"director"`
	Description zero.String `db:"description"`
	CreatedAt   Timestamp   `db:"created_at"`
	UpdatedAt   Timestamp   `db:"updated_at"`

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
	r.DatePrecision = datePrecisionFromDatePtr(o.Date)
	r.Rating = intFromPtr(o.Rating)
	r.StudioID = intFromPtr(o.StudioID)
	r.Director = zero.StringFrom(o.Director)
	r.Description = zero.StringFrom(o.Synopsis)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

func (r *groupRow) resolve() *models.Group {
	ret := &models.Group{
		ID:        r.ID,
		Name:      r.Name.String,
		Aliases:   r.Aliases.String,
		Duration:  nullIntPtr(r.Duration),
		Date:      r.Date.DatePtr(r.DatePrecision),
		Rating:    nullIntPtr(r.Rating),
		StudioID:  nullIntPtr(r.StudioID),
		Director:  r.Director.String,
		Synopsis:  r.Description.String,
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
	r.setNullDate("date", "date_precision", o.Date)
	r.setNullInt("rating", o.Rating)
	r.setNullInt("studio_id", o.StudioID)
	r.setNullString("director", o.Director)
	r.setNullString("description", o.Synopsis)
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
			orderBy:      tagTableSortSQL,
		},
	}
)

type GroupStore struct {
	blobJoinQueryBuilder
	tagRelationshipStore
	groupRelationshipStore

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
		groupRelationshipStore: groupRelationshipStore{
			table: groupRelationshipTableMgr,
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

	if err := qb.groupRelationshipStore.createContainingRelationships(ctx, id, newObject.ContainingGroups); err != nil {
		return err
	}

	if err := qb.groupRelationshipStore.createSubRelationships(ctx, id, newObject.SubGroups); err != nil {
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

	if err := qb.groupRelationshipStore.modifyContainingRelationships(ctx, id, partial.ContainingGroups); err != nil {
		return nil, err
	}

	if err := qb.groupRelationshipStore.modifySubRelationships(ctx, id, partial.SubGroups); err != nil {
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

	if err := qb.groupRelationshipStore.replaceContainingRelationships(ctx, updatedObject.ID, updatedObject.ContainingGroups); err != nil {
		return err
	}

	if err := qb.groupRelationshipStore.replaceSubRelationships(ctx, updatedObject.ID, updatedObject.SubGroups); err != nil {
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
			i := slices.Index(ids, s.ID)
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
	// query := "SELECT * FROM groups WHERE name = ?"
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
	// query := "SELECT * FROM groups WHERE name"
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
		searchColumns := []string{"groups.name", "groups.aliases"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := filterBuilderFromHandler(ctx, &groupFilterHandler{
		groupFilter: groupFilter,
	})

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	if err := qb.setGroupSort(&query, findFilter); err != nil {
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
	"o_counter",
	"sub_group_order",
	"tag_count",
	"updated_at",
}

func (qb *GroupStore) setGroupSort(query *queryBuilder, findFilter *models.FindFilterType) error {
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
		return err
	}

	switch sort {
	case "sub_group_order":
		// sub_group_order is a special sort that sorts by the order_index of the subgroups
		if query.hasJoin("groups_parents") {
			query.sortAndPagination += getSort("order_index", direction, "groups_parents")
		} else {
			// this will give unexpected results if the query is not filtered by a parent group and
			// the group has multiple parents and order indexes
			query.joinSort(groupRelationsTable, "", "groups.id = groups_relations.sub_id")
			query.sortAndPagination += getSort("order_index", direction, groupRelationsTable)
		}
	case "tag_count":
		query.sortAndPagination += getCountSort(groupTable, groupsTagsTable, groupIDColumn, direction)
	case "scenes_count": // generic getSort won't work for this
		query.sortAndPagination += getCountSort(groupTable, groupsScenesTable, groupIDColumn, direction)
	case "o_counter":
		query.sortAndPagination += qb.sortByOCounter(direction)
	default:
		query.sortAndPagination += getSort(sort, direction, "groups")
	}

	// Whatever the sorting, always use name/id as a final sort
	query.sortAndPagination += ", COALESCE(groups.name, groups.id) COLLATE NATURAL_CI ASC"
	return nil
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
	query := `SELECT DISTINCT groups.*
FROM groups
INNER JOIN groups_scenes ON groups.id = groups_scenes.group_id
INNER JOIN performers_scenes ON performers_scenes.scene_id = groups_scenes.scene_id
WHERE performers_scenes.performer_id = ?
`
	args := []interface{}{performerID}
	return qb.queryGroups(ctx, query, args)
}

func (qb *GroupStore) CountByPerformerID(ctx context.Context, performerID int) (int, error) {
	query := `SELECT COUNT(DISTINCT groups_scenes.group_id) AS count
FROM groups_scenes
INNER JOIN performers_scenes ON performers_scenes.scene_id = groups_scenes.scene_id
WHERE performers_scenes.performer_id = ?
`
	args := []interface{}{performerID}
	return groupRepository.runCountQuery(ctx, query, args)
}

func (qb *GroupStore) FindByStudioID(ctx context.Context, studioID int) ([]*models.Group, error) {
	query := `SELECT groups.*
FROM groups
WHERE groups.studio_id = ?
`
	args := []interface{}{studioID}
	return qb.queryGroups(ctx, query, args)
}

func (qb *GroupStore) CountByStudioID(ctx context.Context, studioID int) (int, error) {
	query := `SELECT COUNT(1) AS count
FROM groups
WHERE groups.studio_id = ?
`
	args := []interface{}{studioID}
	return groupRepository.runCountQuery(ctx, query, args)
}

func (qb *GroupStore) GetURLs(ctx context.Context, groupID int) ([]string, error) {
	return groupsURLsTableMgr.get(ctx, groupID)
}

// FindSubGroupIDs returns a list of group IDs where a group in the ids list is a sub-group of the parent group
func (qb *GroupStore) FindSubGroupIDs(ctx context.Context, containingID int, ids []int) ([]int, error) {
	/*
		SELECT gr.sub_id FROM groups_relations gr
		WHERE gr.containing_id = :parentID AND gr.sub_id IN (:ids);
	*/
	table := groupRelationshipTableMgr.table
	q := dialect.From(table).Prepared(true).
		Select(table.Col("sub_id")).Where(
		table.Col("containing_id").Eq(containingID),
		table.Col("sub_id").In(ids),
	)

	const single = false
	var ret []int
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var id int
		if err := r.Scan(&id); err != nil {
			return err
		}

		ret = append(ret, id)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

// FindInAscestors returns a list of group IDs where a group in the ids list is an ascestor of the ancestor group IDs
func (qb *GroupStore) FindInAncestors(ctx context.Context, ascestorIDs []int, ids []int) ([]int, error) {
	/*
		WITH RECURSIVE ascestors AS (
		 SELECT g.id AS parent_id FROM groups g WHERE g.id IN (:ascestorIDs)
		 UNION
		 SELECT gr.containing_id FROM groups_relations gr INNER JOIN ascestors a ON a.parent_id = gr.sub_id
		)
		SELECT p.parent_id FROM ascestors p WHERE p.parent_id IN (:ids);
	*/
	table := qb.table()
	const ascestors = "ancestors"
	const parentID = "parent_id"
	q := dialect.From(ascestors).Prepared(true).
		WithRecursive(ascestors,
			dialect.From(qb.table()).Select(table.Col(idColumn).As(parentID)).
				Where(table.Col(idColumn).In(ascestorIDs)).
				Union(
					dialect.From(groupRelationsJoinTable).InnerJoin(
						goqu.I(ascestors), goqu.On(goqu.I("parent_id").Eq(goqu.I("sub_id"))),
					).Select("containing_id"),
				),
		).Select(parentID).Where(goqu.I(parentID).In(ids))

	const single = false
	var ret []int
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var id int
		if err := r.Scan(&id); err != nil {
			return err
		}

		ret = append(ret, id)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *GroupStore) sortByOCounter(direction string) string {
	// need to sum the o_counter from scenes and images
	return " ORDER BY (" + selectGroupOCountSQL + ") " + direction
}
