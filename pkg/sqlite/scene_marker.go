package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

const sceneMarkerTable = "scene_markers"

const countSceneMarkersForTagQuery = `
SELECT scene_markers.id FROM scene_markers
LEFT JOIN scene_markers_tags as tags_join on tags_join.scene_marker_id = scene_markers.id
WHERE tags_join.tag_id = ? OR scene_markers.primary_tag_id = ?
GROUP BY scene_markers.id
`

type sceneMarkerQueryBuilder struct {
	repository
}

var SceneMarkerReaderWriter = &sceneMarkerQueryBuilder{
	repository{
		tableName: sceneMarkerTable,
		idColumn:  idColumn,
	},
}

func (qb *sceneMarkerQueryBuilder) Create(ctx context.Context, newObject models.SceneMarker) (*models.SceneMarker, error) {
	var ret models.SceneMarker
	if err := qb.insertObject(ctx, newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *sceneMarkerQueryBuilder) Update(ctx context.Context, updatedObject models.SceneMarker) (*models.SceneMarker, error) {
	const partial = false
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	var ret models.SceneMarker
	if err := qb.getByID(ctx, updatedObject.ID, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *sceneMarkerQueryBuilder) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

func (qb *sceneMarkerQueryBuilder) Find(ctx context.Context, id int) (*models.SceneMarker, error) {
	query := "SELECT * FROM scene_markers WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	results, err := qb.querySceneMarkers(ctx, query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *sceneMarkerQueryBuilder) FindMany(ctx context.Context, ids []int) ([]*models.SceneMarker, error) {
	var markers []*models.SceneMarker
	for _, id := range ids {
		marker, err := qb.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		if marker == nil {
			return nil, fmt.Errorf("scene marker with id %d not found", id)
		}

		markers = append(markers, marker)
	}

	return markers, nil
}

func (qb *sceneMarkerQueryBuilder) FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneMarker, error) {
	query := `
		SELECT scene_markers.* FROM scene_markers
		WHERE scene_markers.scene_id = ?
		GROUP BY scene_markers.id
		ORDER BY scene_markers.seconds ASC
	`
	args := []interface{}{sceneID}
	return qb.querySceneMarkers(ctx, query, args)
}

func (qb *sceneMarkerQueryBuilder) CountByTagID(ctx context.Context, tagID int) (int, error) {
	args := []interface{}{tagID, tagID}
	return qb.runCountQuery(ctx, qb.buildCountQuery(countSceneMarkersForTagQuery), args)
}

func (qb *sceneMarkerQueryBuilder) GetMarkerStrings(ctx context.Context, q *string, sort *string) ([]*models.MarkerStringsResultType, error) {
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

func (qb *sceneMarkerQueryBuilder) Wall(ctx context.Context, q *string) ([]*models.SceneMarker, error) {
	s := ""
	if q != nil {
		s = *q
	}
	query := "SELECT scene_markers.* FROM scene_markers WHERE scene_markers.title LIKE '%" + s + "%' ORDER BY RANDOM() LIMIT 80"
	return qb.querySceneMarkers(ctx, query, nil)
}

func (qb *sceneMarkerQueryBuilder) makeFilter(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType) *filterBuilder {
	query := &filterBuilder{}

	query.handleCriterion(ctx, sceneMarkerTagIDCriterionHandler(qb, sceneMarkerFilter.TagID))
	query.handleCriterion(ctx, sceneMarkerTagsCriterionHandler(qb, sceneMarkerFilter.Tags))
	query.handleCriterion(ctx, sceneMarkerSceneTagsCriterionHandler(qb, sceneMarkerFilter.SceneTags))
	query.handleCriterion(ctx, sceneMarkerPerformersCriterionHandler(qb, sceneMarkerFilter.Performers))

	return query
}

func (qb *sceneMarkerQueryBuilder) Query(ctx context.Context, sceneMarkerFilter *models.SceneMarkerFilterType, findFilter *models.FindFilterType) ([]*models.SceneMarker, int, error) {
	if sceneMarkerFilter == nil {
		sceneMarkerFilter = &models.SceneMarkerFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, sceneMarkerTable)

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"scene_markers.title", "scenes.title"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := qb.makeFilter(ctx, sceneMarkerFilter)

	query.addFilter(filter)

	query.sortAndPagination = qb.getSceneMarkerSort(&query, findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	var sceneMarkers []*models.SceneMarker
	for _, id := range idsResult {
		sceneMarker, err := qb.Find(ctx, id)
		if err != nil {
			return nil, 0, err
		}

		sceneMarkers = append(sceneMarkers, sceneMarker)
	}

	return sceneMarkers, countResult, nil
}

func sceneMarkerTagIDCriterionHandler(qb *sceneMarkerQueryBuilder, tagID *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tagID != nil {
			f.addLeftJoin("scene_markers_tags", "", "scene_markers_tags.scene_marker_id = scene_markers.id")

			f.addWhere("(scene_markers.primary_tag_id = ? OR scene_markers_tags.tag_id = ?)", *tagID, *tagID)
		}
	}
}

func sceneMarkerTagsCriterionHandler(qb *sceneMarkerQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tags != nil {
			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("scene_markers_tags", "", "scene_markers.id = scene_markers_tags.scene_marker_id")

				f.addWhere(fmt.Sprintf("%s scene_markers_tags.tag_id IS NULL", notClause))
				return
			}

			if len(tags.Value) == 0 {
				return
			}
			valuesClause := getHierarchicalValues(ctx, qb.tx, tags.Value, tagTable, "tags_relations", "", tags.Depth)

			f.addWith(`marker_tags AS (
SELECT mt.scene_marker_id, t.column1 AS root_tag_id FROM scene_markers_tags mt
INNER JOIN (` + valuesClause + `) t ON t.column2 = mt.tag_id
UNION
SELECT m.id, t.column1 FROM scene_markers m
INNER JOIN (` + valuesClause + `) t ON t.column2 = m.primary_tag_id
)`)

			f.addLeftJoin("marker_tags", "", "marker_tags.scene_marker_id = scene_markers.id")

			addHierarchicalConditionClauses(f, tags, "marker_tags", "root_tag_id")
		}
	}
}

func sceneMarkerSceneTagsCriterionHandler(qb *sceneMarkerQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tags != nil {
			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("scenes_tags", "", "scene_markers.scene_id = scenes_tags.scene_id")

				f.addWhere(fmt.Sprintf("scenes_tags.tag_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 {
				return
			}

			valuesClause := getHierarchicalValues(ctx, qb.tx, tags.Value, tagTable, "tags_relations", "", tags.Depth)

			f.addWith(`scene_tags AS (
SELECT st.scene_id, t.column1 AS root_tag_id FROM scenes_tags st
INNER JOIN (` + valuesClause + `) t ON t.column2 = st.tag_id
)`)

			f.addLeftJoin("scene_tags", "", "scene_tags.scene_id = scene_markers.scene_id")

			addHierarchicalConditionClauses(f, tags, "scene_tags", "root_tag_id")
		}
	}
}

func sceneMarkerPerformersCriterionHandler(qb *sceneMarkerQueryBuilder, performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    performersScenesTable,
		joinAs:       "performers_join",
		primaryFK:    sceneIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			f.addLeftJoin(performersScenesTable, "performers_join", "performers_join.scene_id = scene_markers.scene_id")
		},
	}

	handler := h.handler(performers)
	return func(ctx context.Context, f *filterBuilder) {
		// Make sure scenes is included, otherwise excludes filter fails
		f.addLeftJoin(sceneTable, "", "scenes.id = scene_markers.scene_id")
		handler(ctx, f)
	}
}

func (qb *sceneMarkerQueryBuilder) getSceneMarkerSort(query *queryBuilder, findFilter *models.FindFilterType) string {
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()
	tableName := "scene_markers"
	if sort == "scenes_updated_at" {
		// ensure scene table is joined
		query.join(sceneTable, "", "scenes.id = scene_markers.scene_id")
		sort = "updated_at"
		tableName = "scenes"
	}
	return getSort(sort, direction, tableName)
}

func (qb *sceneMarkerQueryBuilder) querySceneMarkers(ctx context.Context, query string, args []interface{}) ([]*models.SceneMarker, error) {
	var ret models.SceneMarkers
	if err := qb.query(ctx, query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.SceneMarker(ret), nil
}

func (qb *sceneMarkerQueryBuilder) queryMarkerStringsResultType(ctx context.Context, query string, args []interface{}) ([]*models.MarkerStringsResultType, error) {
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

func (qb *sceneMarkerQueryBuilder) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: "scene_markers_tags",
			idColumn:  "scene_marker_id",
		},
		fkColumn: tagIDColumn,
	}
}

func (qb *sceneMarkerQueryBuilder) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return qb.tagsRepository().getIDs(ctx, id)
}

func (qb *sceneMarkerQueryBuilder) UpdateTags(ctx context.Context, id int, tagIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.tagsRepository().replace(ctx, id, tagIDs)
}
