package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

const tagTable = "tags"
const tagIDColumn = "tag_id"
const tagAliasesTable = "tag_aliases"
const tagAliasColumn = "alias"

type tagQueryBuilder struct {
	repository
}

var TagReaderWriter = &tagQueryBuilder{
	repository{
		tableName: tagTable,
		idColumn:  idColumn,
	},
}

func (qb *tagQueryBuilder) Create(ctx context.Context, newObject models.Tag) (*models.Tag, error) {
	var ret models.Tag
	if err := qb.insertObject(ctx, newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *tagQueryBuilder) Update(ctx context.Context, updatedObject models.TagPartial) (*models.Tag, error) {
	const partial = true
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(ctx, updatedObject.ID)
}

func (qb *tagQueryBuilder) UpdateFull(ctx context.Context, updatedObject models.Tag) (*models.Tag, error) {
	const partial = false
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(ctx, updatedObject.ID)
}

func (qb *tagQueryBuilder) Destroy(ctx context.Context, id int) error {
	// TODO - add delete cascade to foreign key
	// delete tag from scenes and markers first
	_, err := qb.tx.Exec(ctx, "DELETE FROM scenes_tags WHERE tag_id = ?", id)
	if err != nil {
		return err
	}

	// TODO - add delete cascade to foreign key
	_, err = qb.tx.Exec(ctx, "DELETE FROM scene_markers_tags WHERE tag_id = ?", id)
	if err != nil {
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

func (qb *tagQueryBuilder) Find(ctx context.Context, id int) (*models.Tag, error) {
	var ret models.Tag
	if err := qb.getByID(ctx, id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *tagQueryBuilder) FindMany(ctx context.Context, ids []int) ([]*models.Tag, error) {
	var tags []*models.Tag
	for _, id := range ids {
		tag, err := qb.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		if tag == nil {
			return nil, fmt.Errorf("tag with id %d not found", id)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (qb *tagQueryBuilder) FindBySceneID(ctx context.Context, sceneID int) ([]*models.Tag, error) {
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

func (qb *tagQueryBuilder) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Tag, error) {
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

func (qb *tagQueryBuilder) FindByImageID(ctx context.Context, imageID int) ([]*models.Tag, error) {
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

func (qb *tagQueryBuilder) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Tag, error) {
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

func (qb *tagQueryBuilder) FindBySceneMarkerID(ctx context.Context, sceneMarkerID int) ([]*models.Tag, error) {
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

func (qb *tagQueryBuilder) FindByName(ctx context.Context, name string, nocase bool) (*models.Tag, error) {
	query := "SELECT * FROM tags WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryTag(ctx, query, args)
}

func (qb *tagQueryBuilder) FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Tag, error) {
	query := "SELECT * FROM tags WHERE name"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryTags(ctx, query, args)
}

func (qb *tagQueryBuilder) FindByParentTagID(ctx context.Context, parentID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		INNER JOIN tags_relations ON tags_relations.child_id = tags.id
		WHERE tags_relations.parent_id = ?
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{parentID}
	return qb.queryTags(ctx, query, args)
}

func (qb *tagQueryBuilder) FindByChildTagID(ctx context.Context, parentID int) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		INNER JOIN tags_relations ON tags_relations.parent_id = tags.id
		WHERE tags_relations.child_id = ?
	`
	query += qb.getDefaultTagSort()
	args := []interface{}{parentID}
	return qb.queryTags(ctx, query, args)
}

func (qb *tagQueryBuilder) Count(ctx context.Context) (int, error) {
	return qb.runCountQuery(ctx, qb.buildCountQuery("SELECT tags.id FROM tags"), nil)
}

func (qb *tagQueryBuilder) All(ctx context.Context) ([]*models.Tag, error) {
	return qb.queryTags(ctx, selectAll("tags")+qb.getDefaultTagSort(), nil)
}

func (qb *tagQueryBuilder) QueryForAutoTag(ctx context.Context, words []string) ([]*models.Tag, error) {
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

func (qb *tagQueryBuilder) validateFilter(tagFilter *models.TagFilterType) error {
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

func (qb *tagQueryBuilder) makeFilter(ctx context.Context, tagFilter *models.TagFilterType) *filterBuilder {
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
	query.handleCriterion(ctx, boolCriterionHandler(tagFilter.IgnoreAutoTag, tagTable+".ignore_auto_tag"))

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

	return query
}

func (qb *tagQueryBuilder) Query(ctx context.Context, tagFilter *models.TagFilterType, findFilter *models.FindFilterType) ([]*models.Tag, int, error) {
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

	query.addFilter(filter)

	query.sortAndPagination = qb.getTagSort(&query, findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	var tags []*models.Tag
	for _, id := range idsResult {
		tag, err := qb.Find(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		tags = append(tags, tag)
	}

	return tags, countResult, nil
}

func tagAliasCriterionHandler(qb *tagQueryBuilder, alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    tagAliasesTable,
		stringColumn: tagAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			qb.aliasRepository().join(f, "", "tags.id")
		},
	}

	return h.handler(alias)
}

func tagIsMissingCriterionHandler(qb *tagQueryBuilder, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "image":
				qb.imageRepository().join(f, "", "tags.id")
				f.addWhere("tags_image.tag_id IS NULL")
			default:
				f.addWhere("(tags." + *isMissing + " IS NULL OR TRIM(tags." + *isMissing + ") = '')")
			}
		}
	}
}

func tagSceneCountCriterionHandler(qb *tagQueryBuilder, sceneCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes_tags", "", "scenes_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes_tags.scene_id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagImageCountCriterionHandler(qb *tagQueryBuilder, imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images_tags", "", "images_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct images_tags.image_id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagGalleryCountCriterionHandler(qb *tagQueryBuilder, galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries_tags", "", "galleries_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries_tags.gallery_id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagPerformerCountCriterionHandler(qb *tagQueryBuilder, performerCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerCount != nil {
			f.addLeftJoin("performers_tags", "", "performers_tags.tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct performers_tags.performer_id)", *performerCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagMarkerCountCriterionHandler(qb *tagQueryBuilder, markerCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if markerCount != nil {
			f.addLeftJoin("scene_markers_tags", "", "scene_markers_tags.tag_id = tags.id")
			f.addLeftJoin("scene_markers", "", "scene_markers_tags.scene_marker_id = scene_markers.id OR scene_markers.primary_tag_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct scene_markers.id)", *markerCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagParentsCriterionHandler(qb *tagQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tags != nil {
			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("tags_relations", "parent_relations", "tags.id = parent_relations.child_id")

				f.addWhere(fmt.Sprintf("parent_relations.parent_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 {
				return
			}

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
	}
}

func tagChildrenCriterionHandler(qb *tagQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tags != nil {
			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("tags_relations", "child_relations", "tags.id = child_relations.parent_id")

				f.addWhere(fmt.Sprintf("child_relations.child_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 {
				return
			}

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
	}
}

func tagParentCountCriterionHandler(qb *tagQueryBuilder, parentCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if parentCount != nil {
			f.addLeftJoin("tags_relations", "parents_count", "parents_count.child_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct parents_count.parent_id)", *parentCount)

			f.addHaving(clause, args...)
		}
	}
}

func tagChildCountCriterionHandler(qb *tagQueryBuilder, childCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if childCount != nil {
			f.addLeftJoin("tags_relations", "children_count", "children_count.parent_id = tags.id")
			clause, args := getIntCriterionWhereClause("count(distinct children_count.child_id)", *childCount)

			f.addHaving(clause, args...)
		}
	}
}

func (qb *tagQueryBuilder) getDefaultTagSort() string {
	return getSort("name", "ASC", "tags")
}

func (qb *tagQueryBuilder) getTagSort(query *queryBuilder, findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}

	if findFilter.Sort != nil {
		switch *findFilter.Sort {
		case "scenes_count":
			return getCountSort(tagTable, scenesTagsTable, tagIDColumn, direction)
		case "scene_markers_count":
			return getCountSort(tagTable, "scene_markers_tags", tagIDColumn, direction)
		case "images_count":
			return getCountSort(tagTable, imagesTagsTable, tagIDColumn, direction)
		case "galleries_count":
			return getCountSort(tagTable, galleriesTagsTable, tagIDColumn, direction)
		case "performers_count":
			return getCountSort(tagTable, performersTagsTable, tagIDColumn, direction)
		}
	}

	return getSort(sort, direction, "tags")
}

func (qb *tagQueryBuilder) queryTag(ctx context.Context, query string, args []interface{}) (*models.Tag, error) {
	results, err := qb.queryTags(ctx, query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *tagQueryBuilder) queryTags(ctx context.Context, query string, args []interface{}) ([]*models.Tag, error) {
	var ret models.Tags
	if err := qb.query(ctx, query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Tag(ret), nil
}

func (qb *tagQueryBuilder) imageRepository() *imageRepository {
	return &imageRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: "tags_image",
			idColumn:  tagIDColumn,
		},
		imageColumn: "image",
	}
}

func (qb *tagQueryBuilder) GetImage(ctx context.Context, tagID int) ([]byte, error) {
	return qb.imageRepository().get(ctx, tagID)
}

func (qb *tagQueryBuilder) HasImage(ctx context.Context, tagID int) (bool, error) {
	return qb.imageRepository().exists(ctx, tagID)
}

func (qb *tagQueryBuilder) UpdateImage(ctx context.Context, tagID int, image []byte) error {
	return qb.imageRepository().replace(ctx, tagID, image)
}

func (qb *tagQueryBuilder) DestroyImage(ctx context.Context, tagID int) error {
	return qb.imageRepository().destroy(ctx, []int{tagID})
}

func (qb *tagQueryBuilder) aliasRepository() *stringRepository {
	return &stringRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: tagAliasesTable,
			idColumn:  tagIDColumn,
		},
		stringColumn: tagAliasColumn,
	}
}

func (qb *tagQueryBuilder) GetAliases(ctx context.Context, tagID int) ([]string, error) {
	return qb.aliasRepository().get(ctx, tagID)
}

func (qb *tagQueryBuilder) UpdateAliases(ctx context.Context, tagID int, aliases []string) error {
	return qb.aliasRepository().replace(ctx, tagID, aliases)
}

func (qb *tagQueryBuilder) Merge(ctx context.Context, source []int, destination int) error {
	if len(source) == 0 {
		return nil
	}

	inBinding := getInBinding(len(source))

	args := []interface{}{destination}
	for _, id := range source {
		if id == destination {
			return errors.New("cannot merge where source == destination")
		}
		args = append(args, id)
	}

	tagTables := map[string]string{
		scenesTagsTable:      sceneIDColumn,
		"scene_markers_tags": "scene_marker_id",
		galleriesTagsTable:   galleryIDColumn,
		imagesTagsTable:      imageIDColumn,
		"performers_tags":    "performer_id",
	}

	args = append(args, destination)
	for table, idColumn := range tagTables {
		_, err := qb.tx.Exec(ctx, `UPDATE `+table+`
SET tag_id = ?
WHERE tag_id IN `+inBinding+`
AND NOT EXISTS(SELECT 1 FROM `+table+` o WHERE o.`+idColumn+` = `+table+`.`+idColumn+` AND o.tag_id = ?)`,
			args...,
		)
		if err != nil {
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

func (qb *tagQueryBuilder) UpdateParentTags(ctx context.Context, tagID int, parentIDs []int) error {
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

func (qb *tagQueryBuilder) UpdateChildTags(ctx context.Context, tagID int, childIDs []int) error {
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
func (qb *tagQueryBuilder) FindAllAncestors(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error) {
	inBinding := getInBinding(len(excludeIDs) + 1)

	query := `WITH RECURSIVE
parents AS (
	SELECT t.id AS parent_id, t.id AS child_id, t.name as path FROM tags t WHERE t.id = ?
	UNION
	SELECT tr.parent_id, tr.child_id, t.name || '->' || p.path as path FROM tags_relations tr INNER JOIN parents p ON p.parent_id = tr.child_id JOIN tags t ON t.id = tr.parent_id WHERE tr.parent_id NOT IN` + inBinding + `
)
SELECT t.*, p.path FROM tags t INNER JOIN parents p ON t.id = p.parent_id
`

	var ret models.TagPaths
	excludeArgs := []interface{}{tagID}
	for _, excludeID := range excludeIDs {
		excludeArgs = append(excludeArgs, excludeID)
	}
	args := []interface{}{tagID}
	args = append(args, append(append(excludeArgs, excludeArgs...), excludeArgs...)...)
	if err := qb.query(ctx, query, args, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// FindAllDescendants returns a slice of TagPath objects, representing all
// descendants of the tag with the provided id.
func (qb *tagQueryBuilder) FindAllDescendants(ctx context.Context, tagID int, excludeIDs []int) ([]*models.TagPath, error) {
	inBinding := getInBinding(len(excludeIDs) + 1)

	query := `WITH RECURSIVE
children AS (
	SELECT t.id AS parent_id, t.id AS child_id, t.name as path FROM tags t WHERE t.id = ?
	UNION
	SELECT tr.parent_id, tr.child_id, c.path || '->' || t.name as path FROM tags_relations tr INNER JOIN children c ON c.child_id = tr.parent_id JOIN tags t ON t.id = tr.child_id WHERE tr.child_id NOT IN` + inBinding + `
)
SELECT t.*, c.path FROM tags t INNER JOIN children c ON t.id = c.child_id
`

	var ret models.TagPaths
	excludeArgs := []interface{}{tagID}
	for _, excludeID := range excludeIDs {
		excludeArgs = append(excludeArgs, excludeID)
	}
	args := []interface{}{tagID}
	args = append(args, append(append(excludeArgs, excludeArgs...), excludeArgs...)...)
	if err := qb.query(ctx, query, args, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
