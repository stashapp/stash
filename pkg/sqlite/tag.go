package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

const (
	tagTable        = "tags"
	tagIDColumn     = "tag_id"
	tagAliasesTable = "tag_aliases"
	tagAliasColumn  = "alias"

	tagImageBlobColumn = "image_blob"
)

type tagRow struct {
	ID            int         `db:"id" goqu:"skipinsert"`
	Name          null.String `db:"name"` // TODO: make schema non-nullable
	Description   zero.String `db:"description"`
	IgnoreAutoTag bool        `db:"ignore_auto_tag"`
	CreatedAt     Timestamp   `db:"created_at"`
	UpdatedAt     Timestamp   `db:"updated_at"`

	// not used in resolutions or updates
	ImageBlob zero.String `db:"image_blob"`
}

func (r *tagRow) fromTag(o models.Tag) {
	r.ID = o.ID
	r.Name = null.StringFrom(o.Name)
	r.Description = zero.StringFrom(o.Description)
	r.IgnoreAutoTag = o.IgnoreAutoTag
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

func (r *tagRow) resolve() *models.Tag {
	ret := &models.Tag{
		ID:            r.ID,
		Name:          r.Name.String,
		Description:   r.Description.String,
		IgnoreAutoTag: r.IgnoreAutoTag,
		CreatedAt:     r.CreatedAt.Timestamp,
		UpdatedAt:     r.UpdatedAt.Timestamp,
	}

	return ret
}

type tagPathRow struct {
	tagRow
	Path string `db:"path"`
}

func (r *tagPathRow) resolve() *models.TagPath {
	ret := &models.TagPath{
		Tag:  *r.tagRow.resolve(),
		Path: r.Path,
	}

	return ret
}

type tagRowRecord struct {
	updateRecord
}

func (r *tagRowRecord) fromPartial(o models.TagPartial) {
	r.setString("name", o.Name)
	r.setNullString("description", o.Description)
	r.setBool("ignore_auto_tag", o.IgnoreAutoTag)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
}

type TagStore struct {
	repository
	blobJoinQueryBuilder

	tableMgr *table
}

func NewTagStore(blobStore *BlobStore) *TagStore {
	return &TagStore{
		repository: repository{
			tableName: tagTable,
			idColumn:  idColumn,
		},
		blobJoinQueryBuilder: blobJoinQueryBuilder{
			blobStore: blobStore,
			joinTable: tagTable,
		},
		tableMgr: tagTableMgr,
	}
}

func (qb *TagStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *TagStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table()).Select(qb.table().All())
}

func (qb *TagStore) Create(ctx context.Context, newObject *models.Tag) error {
	var r tagRow
	r.fromTag(*newObject)

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

func (qb *TagStore) UpdatePartial(ctx context.Context, id int, partial models.TagPartial) (*models.Tag, error) {
	r := tagRowRecord{
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

	return qb.find(ctx, id)
}

func (qb *TagStore) Update(ctx context.Context, updatedObject *models.Tag) error {
	var r tagRow
	r.fromTag(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *TagStore) Destroy(ctx context.Context, id int) error {
	// must handle image checksums manually
	if err := qb.destroyImage(ctx, id); err != nil {
		return err
	}

	// cannot unset primary_tag_id in scene_markers because it is not nullable
	countQuery := "SELECT COUNT(*) as count FROM scene_markers where primary_tag_id = ?"
	args := []interface{}{id}
	primaryMarkers, err := qb.runCountQuery(ctx, countQuery, args)
	if err != nil {
		return err
	}

	if primaryMarkers > 0 {
		return errors.New("cannot delete tag used as a primary tag in scene markers")
	}

	return qb.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *TagStore) Find(ctx context.Context, id int) (*models.Tag, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *TagStore) FindMany(ctx context.Context, ids []int) ([]*models.Tag, error) {
	ret := make([]*models.Tag, len(ids))

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
			return nil, fmt.Errorf("tag with id %d not found", ids[i])
		}
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *TagStore) find(ctx context.Context, id int) (*models.Tag, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *TagStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Tag, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *TagStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Tag, error) {
	const single = false
	var ret []*models.Tag
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f tagRow
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

func (qb *TagStore) FindBySceneID(ctx context.Context, sceneID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN scenes_tags as scenes_join on scenes_join.tag_id = tags.id
		WHERE scenes_join.scene_id = ?
		GROUP BY tags.id
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{sceneID}
	return qb.queryTags(ctx, query, args)
}

func (qb *TagStore) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN performers_tags as performers_join on performers_join.tag_id = tags.id
		WHERE performers_join.performer_id = ?
		GROUP BY tags.id
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{performerID}
	return qb.queryTags(ctx, query, args)
}

func (qb *TagStore) FindByImageID(ctx context.Context, imageID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN images_tags as images_join on images_join.tag_id = tags.id
		WHERE images_join.image_id = ?
		GROUP BY tags.id
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{imageID}
	return qb.queryTags(ctx, query, args)
}

func (qb *TagStore) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN galleries_tags as galleries_join on galleries_join.tag_id = tags.id
		WHERE galleries_join.gallery_id = ?
		GROUP BY tags.id
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{galleryID}
	return qb.queryTags(ctx, query, args)
}

func (qb *TagStore) FindBySceneMarkerID(ctx context.Context, sceneMarkerID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN scene_markers_tags as scene_markers_join on scene_markers_join.tag_id = tags.id
		WHERE scene_markers_join.scene_marker_id = ?
		GROUP BY tags.id
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{sceneMarkerID}
	return qb.queryTags(ctx, query, args)
}

func (qb *TagStore) FindByName(ctx context.Context, name string, nocase bool) (*models.Tag, error) {
	// query := "SELECT * FROM tags WHERE name = ?"
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

func (qb *TagStore) FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Tag, error) {
	// query := "SELECT * FROM tags WHERE name"
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

func (qb *TagStore) FindByParentTagID(ctx context.Context, parentID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		INNER JOIN tags_relations ON tags_relations.child_id = tags.id
		WHERE tags_relations.parent_id = ?
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{parentID}
	return qb.queryTags(ctx, query, args)
}

func (qb *TagStore) FindByChildTagID(ctx context.Context, parentID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		INNER JOIN tags_relations ON tags_relations.parent_id = tags.id
		WHERE tags_relations.child_id = ?
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{parentID}
	return qb.queryTags(ctx, query, args)
}

func (qb *TagStore) CountByParentTagID(ctx context.Context, parentID int) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(goqu.T("tags")).
		InnerJoin(goqu.T("tags_relations"), goqu.On(goqu.I("tags_relations.parent_id").Eq(goqu.I("tags.id")))).
		Where(goqu.I("tags_relations.child_id").Eq(goqu.V(parentID))) // Pass the parentID here
	return count(ctx, q)
}

func (qb *TagStore) CountByChildTagID(ctx context.Context, childID int) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(goqu.T("tags")).
		InnerJoin(goqu.T("tags_relations"), goqu.On(goqu.I("tags_relations.child_id").Eq(goqu.I("tags.id")))).
		Where(goqu.I("tags_relations.parent_id").Eq(goqu.V(childID))) // Pass the childID here
	return count(ctx, q)
}

func (qb *TagStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *TagStore) All(ctx context.Context) ([]*models.Tag, error) {
	table := qb.table()

	return qb.getMany(ctx, qb.selectDataset().Order(
		table.Col("name").Asc(),
		table.Col(idColumn).Asc(),
	))
}

func (qb *TagStore) QueryForAutoTag(ctx context.Context, words []string) ([]*models.Tag, error) {
	// TODO - Query needs to be changed to support queries of this type, and
	// this method should be removed
	query := selectAll(tagTable)
	query += " LEFT JOIN tag_aliases ON tag_aliases.tag_id = tags.id"

	var whereClauses []string
	var args []interface{}

	for _, w := range words {
		ww := w + "%"
		whereClauses = append(whereClauses, "tags.name like ?")
		args = append(args, ww)

		// include aliases
		whereClauses = append(whereClauses, "tag_aliases.alias like ?")
		args = append(args, ww)
	}

	whereOr := "(" + strings.Join(whereClauses, " OR ") + ")"
	where := strings.Join([]string{
		"tags.ignore_auto_tag = 0",
		whereOr,
	}, " AND ")
	return qb.queryTags(ctx, query+" WHERE "+where, args)
}

func (qb *TagStore) validateFilter(tagFilter *models.TagFilterType) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if tagFilter.And != nil {
		if tagFilter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if tagFilter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(tagFilter.And)
	}

	if tagFilter.Or != nil {
		if tagFilter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(tagFilter.Or)
	}

	if tagFilter.Not != nil {
		return qb.validateFilter(tagFilter.Not)
	}

	return nil
}

func (qb *TagStore) makeFilter(ctx context.Context, tagFilter *models.TagFilterType) *filterBuilder {
	query := &filterBuilder{}

	if tagFilter.And != nil {
		query.and(qb.makeFilter(ctx, tagFilter.And))
	}
	if tagFilter.Or != nil {
		query.or(qb.makeFilter(ctx, tagFilter.Or))
	}
	if tagFilter.Not != nil {
		query.not(qb.makeFilter(ctx, tagFilter.Not))
	}

	query.handleCriterion(ctx, stringCriterionHandler(tagFilter.Name, tagTable+".name"))
	query.handleCriterion(ctx, tagAliasCriterionHandler(qb, tagFilter.Aliases))

	query.handleCriterion(ctx, stringCriterionHandler(tagFilter.Description, tagTable+".description"))
	query.handleCriterion(ctx, boolCriterionHandler(tagFilter.IgnoreAutoTag, tagTable+".ignore_auto_tag", nil))

	query.handleCriterion(ctx, tagIsMissingCriterionHandler(qb, tagFilter.IsMissing))
	query.handleCriterion(ctx, tagSceneCountCriterionHandler(qb, tagFilter.SceneCount))
	query.handleCriterion(ctx, tagImageCountCriterionHandler(qb, tagFilter.ImageCount))
	query.handleCriterion(ctx, tagGalleryCountCriterionHandler(qb, tagFilter.GalleryCount))
	query.handleCriterion(ctx, tagPerformerCountCriterionHandler(qb, tagFilter.PerformerCount))
	query.handleCriterion(ctx, tagMarkerCountCriterionHandler(qb, tagFilter.MarkerCount))
	query.handleCriterion(ctx, tagParentsCriterionHandler(qb, tagFilter.Parents))
	query.handleCriterion(ctx, tagChildrenCriterionHandler(qb, tagFilter.Children))
	query.handleCriterion(ctx, tagParentCountCriterionHandler(qb, tagFilter.ParentCount))
	query.handleCriterion(ctx, tagChildCountCriterionHandler(qb, tagFilter.ChildCount))
	query.handleCriterion(ctx, timestampCriterionHandler(tagFilter.CreatedAt, "tags.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(tagFilter.UpdatedAt, "tags.updated_at"))

	return query
}

func (qb *TagStore) Query(ctx context.Context, tagFilter *models.TagFilterType, findFilter *models.FindFilterType) ([]*models.Tag, int, error) {
	if tagFilter == nil {
		tagFilter = &models.TagFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, tagTable)

	if q := findFilter.Q; q != nil && *q != "" {
		query.join(tagAliasesTable, "", "tag_aliases.tag_id = tags.id")
		searchColumns := []string{"tags.name", "tag_aliases.alias"}
		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(tagFilter); err != nil {
		return nil, 0, err
	}
	filter := qb.makeFilter(ctx, tagFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, 0, err
	}

	query.sortAndPagination = qb.getTagSort(&query, findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	tags, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return tags, countResult, nil
}

func tagAliasCriterionHandler(qb *TagStore, alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    tagAliasesTable,
		stringColumn: tagAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			qb.aliasRepository().join(f, "", "tags.id")
		},
	}

	return h.handler(alias)
}

func tagIsMissingCriterionHandler(qb *TagStore, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "image":
				f.addWhere("tags.image_blob IS NULL")
			default:
				f.addWhere("(tags." + *isMissing + " IS NULL OR TRIM(tags." + *isMissing + ") = '')")
			}
		}
	}
}

func tagSceneCountCriterionHandler(qb *TagStore, sceneCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes_tags", "", "scenes_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes_tags.scene_id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagImageCountCriterionHandler(qb *TagStore, imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images_tags", "", "images_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct images_tags.image_id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagGalleryCountCriterionHandler(qb *TagStore, galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries_tags", "", "galleries_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries_tags.gallery_id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagPerformerCountCriterionHandler(qb *TagStore, performerCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerCount != nil {
			f.addLeftJoin("performers_tags", "", "performers_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct performers_tags.performer_id)", *performerCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagMarkerCountCriterionHandler(qb *TagStore, markerCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if markerCount != nil {
			f.addLeftJoin("scene_markers_tags", "", "scene_markers_tags.tag_id = tags.id")
			f.addLeftJoin("scene_markers", "", "scene_markers_tags.scene_marker_id = scene_markers.id OR scene_markers.primary_tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct scene_markers.id)", *markerCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagParentsCriterionHandler(qb *TagStore, criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			tags := criterion.CombineExcludes()

			// validate the modifier
			switch tags.Modifier {
			case models.CriterionModifierIncludesAll, models.CriterionModifierIncludes, models.CriterionModifierExcludes, models.CriterionModifierIsNull, models.CriterionModifierNotNull:
				// valid
			default:
				f.setError(fmt.Errorf("invalid modifier %s for tag parent/children", criterion.Modifier))
			}

			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("tags_relations", "parent_relations", "tags.id = parent_relations.child_id")

				f.addWhere(fmt.Sprintf("parent_relations.parent_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 && len(tags.Excludes) == 0 {
				return
			}

			if len(tags.Value) > 0 {
				var args []interface{}
				for _, val := range tags.Value {
					args = append(args, val)
				}

				depthVal := 0
				if tags.Depth != nil {
					depthVal = *tags.Depth
				}

				var depthCondition string
				if depthVal != -1 {
					depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
				}

				query := `parents AS (
		SELECT parent_id AS root_id, child_id AS item_id, 0 AS depth FROM tags_relations WHERE parent_id IN` + getInBinding(len(tags.Value)) + `
		UNION
		SELECT root_id, child_id, depth + 1 FROM tags_relations INNER JOIN parents ON item_id = parent_id ` + depthCondition + `
	)`

				f.addRecursiveWith(query, args...)

				f.addLeftJoin("parents", "", "parents.item_id = tags.id")

				addHierarchicalConditionClauses(f, tags, "parents", "root_id")
			}

			if len(tags.Excludes) > 0 {
				var args []interface{}
				for _, val := range tags.Excludes {
					args = append(args, val)
				}

				depthVal := 0
				if tags.Depth != nil {
					depthVal = *tags.Depth
				}

				var depthCondition string
				if depthVal != -1 {
					depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
				}

				query := `parents2 AS (
		SELECT parent_id AS root_id, child_id AS item_id, 0 AS depth FROM tags_relations WHERE parent_id IN` + getInBinding(len(tags.Excludes)) + `
		UNION
		SELECT root_id, child_id, depth + 1 FROM tags_relations INNER JOIN parents2 ON item_id = parent_id ` + depthCondition + `
	)`

				f.addRecursiveWith(query, args...)

				f.addLeftJoin("parents2", "", "parents2.item_id = tags.id")

				addHierarchicalConditionClauses(f, models.HierarchicalMultiCriterionInput{
					Value:    tags.Excludes,
					Depth:    tags.Depth,
					Modifier: models.CriterionModifierExcludes,
				}, "parents2", "root_id")
			}
		}
	}
}

func tagChildrenCriterionHandler(qb *TagStore, criterion *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if criterion != nil {
			tags := criterion.CombineExcludes()

			// validate the modifier
			switch tags.Modifier {
			case models.CriterionModifierIncludesAll, models.CriterionModifierIncludes, models.CriterionModifierExcludes, models.CriterionModifierIsNull, models.CriterionModifierNotNull:
				// valid
			default:
				f.setError(fmt.Errorf("invalid modifier %s for tag parent/children", criterion.Modifier))
			}

			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("tags_relations", "child_relations", "tags.id = child_relations.parent_id")

				f.addWhere(fmt.Sprintf("child_relations.child_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 && len(tags.Excludes) == 0 {
				return
			}

			if len(tags.Value) > 0 {
				var args []interface{}
				for _, val := range tags.Value {
					args = append(args, val)
				}

				depthVal := 0
				if tags.Depth != nil {
					depthVal = *tags.Depth
				}

				var depthCondition string
				if depthVal != -1 {
					depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
				}

				query := `children AS (
		SELECT child_id AS root_id, parent_id AS item_id, 0 AS depth FROM tags_relations WHERE child_id IN` + getInBinding(len(tags.Value)) + `
		UNION
		SELECT root_id, parent_id, depth + 1 FROM tags_relations INNER JOIN children ON item_id = child_id ` + depthCondition + `
	)`

				f.addRecursiveWith(query, args...)

				f.addLeftJoin("children", "", "children.item_id = tags.id")

				addHierarchicalConditionClauses(f, tags, "children", "root_id")
			}

			if len(tags.Excludes) > 0 {
				var args []interface{}
				for _, val := range tags.Excludes {
					args = append(args, val)
				}

				depthVal := 0
				if tags.Depth != nil {
					depthVal = *tags.Depth
				}

				var depthCondition string
				if depthVal != -1 {
					depthCondition = fmt.Sprintf("WHERE depth < %d", depthVal)
				}

				query := `children2 AS (
		SELECT child_id AS root_id, parent_id AS item_id, 0 AS depth FROM tags_relations WHERE child_id IN` + getInBinding(len(tags.Excludes)) + `
		UNION
		SELECT root_id, parent_id, depth + 1 FROM tags_relations INNER JOIN children2 ON item_id = child_id ` + depthCondition + `
	)`

				f.addRecursiveWith(query, args...)

				f.addLeftJoin("children2", "", "children2.item_id = tags.id")

				addHierarchicalConditionClauses(f, models.HierarchicalMultiCriterionInput{
					Value:    tags.Excludes,
					Depth:    tags.Depth,
					Modifier: models.CriterionModifierExcludes,
				}, "children2", "root_id")
			}
		}
	}
}

func tagParentCountCriterionHandler(qb *TagStore, parentCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if parentCount != nil {
			f.addLeftJoin("tags_relations", "parents_count", "parents_count.child_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct parents_count.parent_id)", *parentCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagChildCountCriterionHandler(qb *TagStore, childCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if childCount != nil {
			f.addLeftJoin("tags_relations", "children_count", "children_count.parent_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct children_count.child_id)", *childCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *TagStore) getDefaultTagSort() string {
	return getSort("name", "ASC", "tags")
}

func (qb *TagStore) getTagSort(query *queryBuilder, findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}

	sortQuery := ""
	switch sort {
	case "scenes_count":
		sortQuery += getCountSort(tagTable, scenesTagsTable, tagIDColumn, direction)
	case "scene_markers_count":
		sortQuery += fmt.Sprintf(" ORDER BY (SELECT COUNT(*) FROM scene_markers_tags WHERE tags.id = scene_markers_tags.tag_id)+(SELECT COUNT(*) FROM scene_markers WHERE tags.id = scene_markers.primary_tag_id) %s", getSortDirection(direction))
	case "images_count":
		sortQuery += getCountSort(tagTable, imagesTagsTable, tagIDColumn, direction)
	case "galleries_count":
		sortQuery += getCountSort(tagTable, galleriesTagsTable, tagIDColumn, direction)
	case "performers_count":
		sortQuery += getCountSort(tagTable, performersTagsTable, tagIDColumn, direction)
	default:
		sortQuery += getSort(sort, direction, "tags")
	}

	// Whatever the sorting, always use name/id as a final sort
	sortQuery += ", COALESCE(tags.name, tags.id) COLLATE NATURAL_CI ASC"
	return sortQuery
}

func (qb *TagStore) queryTags(ctx context.Context, query string, args []interface{}) ([]*models.Tag, error) {
	const single = false
	var ret []*models.Tag
	if err := qb.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f tagRow
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

func (qb *TagStore) queryTagPaths(ctx context.Context, query string, args []interface{}) ([]*models.TagPath, error) {
	const single = false
	var ret []*models.TagPath
	if err := qb.queryFunc(ctx, query, args, single, func(r *sqlx.Rows) error {
		var f tagPathRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		t := f.resolve()

		ret = append(ret, t)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *TagStore) GetImage(ctx context.Context, tagID int) ([]byte, error) {
	return qb.blobJoinQueryBuilder.GetImage(ctx, tagID, tagImageBlobColumn)
}

func (qb *TagStore) HasImage(ctx context.Context, tagID int) (bool, error) {
	return qb.blobJoinQueryBuilder.HasImage(ctx, tagID, tagImageBlobColumn)
}

func (qb *TagStore) UpdateImage(ctx context.Context, tagID int, image []byte) error {
	return qb.blobJoinQueryBuilder.UpdateImage(ctx, tagID, tagImageBlobColumn, image)
}

func (qb *TagStore) destroyImage(ctx context.Context, tagID int) error {
	return qb.blobJoinQueryBuilder.DestroyImage(ctx, tagID, tagImageBlobColumn)
}

func (qb *TagStore) aliasRepository() *stringRepository {
	return &stringRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: tagAliasesTable,
			idColumn:  tagIDColumn,
		},
		stringColumn: tagAliasColumn,
	}
}

func (qb *TagStore) GetAliases(ctx context.Context, tagID int) ([]string, error) {
	return qb.aliasRepository().get(ctx, tagID)
}

func (qb *TagStore) UpdateAliases(ctx context.Context, tagID int, aliases []string) error {
	return qb.aliasRepository().replace(ctx, tagID, aliases)
}

func (qb *TagStore) Merge(ctx context.Context, source []int, destination int) error {
	if len(source) == 0 {
		return nil
	}

	inBinding := getInBinding(len(source))

	args := []interface{}{destination}
	srcArgs := make([]interface{}, len(source))
	for i, id := range source {
		if id == destination {
			return errors.New("cannot merge where source == destination")
		}
		srcArgs[i] = id
	}

	args = append(args, srcArgs...)

	tagTables := map[string]string{
		scenesTagsTable:      sceneIDColumn,
		"scene_markers_tags": "scene_marker_id",
		galleriesTagsTable:   galleryIDColumn,
		imagesTagsTable:      imageIDColumn,
		"performers_tags":    "performer_id",
	}

	args = append(args, destination)
	for table, idColumn := range tagTables {
		_, err := qb.tx.Exec(ctx, `UPDATE OR IGNORE `+table+`
SET tag_id = ?
WHERE tag_id IN `+inBinding+`
AND NOT EXISTS(SELECT 1 FROM `+table+` o WHERE o.`+idColumn+` = `+table+`.`+idColumn+` AND o.tag_id = ?)`,
			args...,
		)
		if err != nil {
			return err
		}

		// delete source tag ids from the table where they couldn't be set
		if _, err := qb.tx.Exec(ctx, `DELETE FROM `+table+` WHERE tag_id IN `+inBinding, srcArgs...); err != nil {
			return err
		}
	}

	_, err := qb.tx.Exec(ctx, "UPDATE "+sceneMarkerTable+" SET primary_tag_id = ? WHERE primary_tag_id IN "+inBinding, args...)
	if err != nil {
		return err
	}

	_, err = qb.tx.Exec(ctx, "INSERT INTO "+tagAliasesTable+" (tag_id, alias) SELECT ?, name FROM "+tagTable+" WHERE id IN "+inBinding, args...)
	if err != nil {
		return err
	}

	_, err = qb.tx.Exec(ctx, "UPDATE "+tagAliasesTable+" SET tag_id = ? WHERE tag_id IN "+inBinding, args...)
	if err != nil {
		return err
	}

	for _, id := range source {
		err = qb.Destroy(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (qb *TagStore) UpdateParentTags(ctx context.Context, tagID int, parentIDs []int) error {
	tx := qb.tx
	if _, err := tx.Exec(ctx, "DELETE FROM tags_relations WHERE child_id = ?", tagID); err != nil {
		return err
	}

	if len(parentIDs) > 0 {
		var args []interface{}
		var values []string
		for _, parentID := range parentIDs {
			values = append(values, "(? , ?)")
			args = append(args, parentID, tagID)
		}

		query := "INSERT INTO tags_relations (parent_id, child_id) VALUES " + strings.Join(values, ", ")
		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return err
		}
	}

	return nil
}

func (qb *TagStore) UpdateChildTags(ctx context.Context, tagID int, childIDs []int) error {
	tx := qb.tx
	if _, err := tx.Exec(ctx, "DELETE FROM tags_relations WHERE parent_id = ?", tagID); err != nil {
		return err
	}

	if len(childIDs) > 0 {
		var args []interface{}
		var values []string
		for _, childID := range childIDs {
			values = append(values, "(? , ?)")
			args = append(args, tagID, childID)
		}

		query := "INSERT INTO tags_relations (parent_id, child_id) VALUES " + strings.Join(values, ", ")
		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return err
		}
	}

	return nil
}

// FindAllAncestors returns a slice of TagPath objects, representing all
// ancestors of the tag with the provided id.
func (qb *TagStore) FindAllAncestors(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error) {
	inBinding := getInBinding(len(excludeIDs) + 1)

	query := `WITH RECURSIVE
parents AS (
	SELECT t.id AS parent_id, t.id AS child_id, t.name as path FROM tags t WHERE t.id = ?
	UNION
	SELECT tr.parent_id, tr.child_id, t.name || '->' || p.path as path FROM tags_relations tr INNER JOIN parents p ON p.parent_id = tr.child_id JOIN tags t ON t.id = tr.parent_id WHERE tr.parent_id NOT IN` + inBinding + `
)
SELECT t.*, p.path FROM tags t INNER JOIN parents p ON t.id = p.parent_id
`

	excludeArgs := []interface{}{tagID}
	for _, excludeID := range excludeIDs {
		excludeArgs = append(excludeArgs, excludeID)
	}
	args := []interface{}{tagID}
	args = append(args, append(append(excludeArgs, excludeArgs...), excludeArgs...)...)

	return qb.queryTagPaths(ctx, query, args)
}

// FindAllDescendants returns a slice of TagPath objects, representing all
// descendants of the tag with the provided id.
func (qb *TagStore) FindAllDescendants(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error) {
	inBinding := getInBinding(len(excludeIDs) + 1)

	query := `WITH RECURSIVE
children AS (
	SELECT t.id AS parent_id, t.id AS child_id, t.name as path FROM tags t WHERE t.id = ?
	UNION
	SELECT tr.parent_id, tr.child_id, c.path || '->' || t.name as path FROM tags_relations tr INNER JOIN children c ON c.child_id = tr.parent_id JOIN tags t ON t.id = tr.child_id WHERE tr.child_id NOT IN` + inBinding + `
)
SELECT t.*, c.path FROM tags t INNER JOIN children c ON t.id = c.child_id
`

	excludeArgs := []interface{}{tagID}
	for _, excludeID := range excludeIDs {
		excludeArgs = append(excludeArgs, excludeID)
	}
	args := []interface{}{tagID}
	args = append(args, append(append(excludeArgs, excludeArgs...), excludeArgs...)...)

	return qb.queryTagPaths(ctx, query, args)
}
