package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
	"strconv"
)

const countSceneMarkersForTagQuery = `
SELECT scene_markers.id FROM scene_markers
LEFT JOIN scene_markers_tags as tags_join on tags_join.scene_marker_id = scene_markers.id
WHERE tags_join.tag_id = ? OR scene_markers.primary_tag_id = ?
GROUP BY scene_markers.id
`

type SceneMarkerQueryBuilder struct{}

func NewSceneMarkerQueryBuilder() SceneMarkerQueryBuilder {
	return SceneMarkerQueryBuilder{}
}

func (qb *SceneMarkerQueryBuilder) Create(newSceneMarker SceneMarker, tx *sqlx.Tx) (*SceneMarker, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO scene_markers (title, seconds, primary_tag_id, scene_id, created_at, updated_at)
				VALUES (:title, :seconds, :primary_tag_id, :scene_id, :created_at, :updated_at)
		`,
		newSceneMarker,
	)
	if err != nil {
		return nil, err
	}
	sceneMarkerID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&newSceneMarker, `SELECT * FROM scene_markers WHERE id = ? LIMIT 1`, sceneMarkerID); err != nil {
		return nil, err
	}
	return &newSceneMarker, nil
}

func (qb *SceneMarkerQueryBuilder) Update(updatedSceneMarker SceneMarker, tx *sqlx.Tx) (*SceneMarker, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE scene_markers SET `+SQLGenKeys(updatedSceneMarker)+` WHERE scene_markers.id = :id`,
		updatedSceneMarker,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedSceneMarker, `SELECT * FROM scene_markers WHERE id = ? LIMIT 1`, updatedSceneMarker.ID); err != nil {
		return nil, err
	}
	return &updatedSceneMarker, nil
}

func (qb *SceneMarkerQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	return executeDeleteQuery("scene_markers", id, tx)
}

func (qb *SceneMarkerQueryBuilder) Find(id int) (*SceneMarker, error) {
	query := "SELECT * FROM scene_markers WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	results, err := qb.querySceneMarkers(query, args, nil)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *SceneMarkerQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*SceneMarker, error) {
	query := `
		SELECT scene_markers.* FROM scene_markers
		WHERE scene_markers.scene_id = ?
		GROUP BY scene_markers.id
		ORDER BY scene_markers.seconds ASC
	`
	args := []interface{}{sceneID}
	return qb.querySceneMarkers(query, args, tx)
}

func (qb *SceneMarkerQueryBuilder) CountByTagID(tagID int) (int, error) {
	args := []interface{}{tagID, tagID}
	return runCountQuery(buildCountQuery(countSceneMarkersForTagQuery), args)
}

func (qb *SceneMarkerQueryBuilder) GetMarkerStrings(q *string, sort *string) ([]*MarkerStringsResultType, error) {
	query := "SELECT count(*) as `count`, scene_markers.id as id, scene_markers.title as title FROM scene_markers"
	if q != nil {
		query = query + " WHERE title LIKE '%" + *q + "%'"
	}
	query = query + " GROUP BY title"
	if sort != nil && *sort == "count" {
		query = query + " ORDER BY `count` DESC"
	} else {
		query = query + " ORDER BY title ASC"
	}
	var args []interface{}
	return qb.queryMarkerStringsResultType(query, args)
}

func (qb *SceneMarkerQueryBuilder) Wall(q *string) ([]*SceneMarker, error) {
	s := ""
	if q != nil {
		s = *q
	}
	query := "SELECT scene_markers.* FROM scene_markers WHERE scene_markers.title LIKE '%" + s + "%' ORDER BY RANDOM() LIMIT 80"
	return qb.querySceneMarkers(query, nil, nil)
}

func (qb *SceneMarkerQueryBuilder) Query(sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) ([]*SceneMarker, int) {
	if sceneMarkerFilter == nil {
		sceneMarkerFilter = &SceneMarkerFilterType{}
	}
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	var whereClauses []string
	var havingClauses []string
	var args []interface{}
	body := selectDistinctIDs("scene_markers")
	body = body + `
		left join tags as primary_tag on primary_tag.id = scene_markers.primary_tag_id
		left join scenes as scene on scene.id = scene_markers.scene_id
		left join scene_markers_tags as tags_join on tags_join.scene_marker_id = scene_markers.id
		left join tags on tags_join.tag_id = tags.id
	`

	if tagsFilter := sceneMarkerFilter.Tags; tagsFilter != nil && len(tagsFilter.Value) > 0 {
		//select `scene_markers`.* from `scene_markers`
		//left join `tags` as `primary_tags_join`
		//  on `primary_tags_join`.`id` = `scene_markers`.`primary_tag_id`
		//  and `primary_tags_join`.`id` in ('3', '37', '9', '89')
		//left join `scene_markers_tags` as `tags_join`
		//  on `tags_join`.`scene_marker_id` = `scene_markers`.`id`
		//  and `tags_join`.`tag_id` in ('3', '37', '9', '89')
		//group by `scene_markers`.`id`
		//having ((count(distinct `primary_tags_join`.`id`) + count(distinct `tags_join`.`tag_id`)) = 4)

		length := len(tagsFilter.Value)

		if tagsFilter.Modifier == CriterionModifierIncludes || tagsFilter.Modifier == CriterionModifierIncludesAll {
			body += " LEFT JOIN tags AS ptj ON ptj.id = scene_markers.primary_tag_id AND ptj.id IN " + getInBinding(length)
			body += " LEFT JOIN scene_markers_tags AS tj ON tj.scene_marker_id = scene_markers.id AND tj.tag_id IN " + getInBinding(length)

			// only one required for include any
			requiredCount := 1

			// all required for include all
			if tagsFilter.Modifier == CriterionModifierIncludesAll {
				requiredCount = length
			}

			havingClauses = append(havingClauses, "((COUNT(DISTINCT ptj.id) + COUNT(DISTINCT tj.tag_id)) >= "+strconv.Itoa(requiredCount)+")")
		} else if tagsFilter.Modifier == CriterionModifierExcludes {
			// excludes all of the provided ids
			whereClauses = append(whereClauses, "scene_markers.primary_tag_id not in "+getInBinding(length))
			whereClauses = append(whereClauses, "not exists (select smt.scene_marker_id from scene_markers_tags as smt where smt.scene_marker_id = scene_markers.id and smt.tag_id in "+getInBinding(length)+")")
		}

		for _, tagID := range tagsFilter.Value {
			args = append(args, tagID)
		}
		for _, tagID := range tagsFilter.Value {
			args = append(args, tagID)
		}
	}

	if sceneTagsFilter := sceneMarkerFilter.SceneTags; sceneTagsFilter != nil && len(sceneTagsFilter.Value) > 0 {
		length := len(sceneTagsFilter.Value)

		if sceneTagsFilter.Modifier == CriterionModifierIncludes || sceneTagsFilter.Modifier == CriterionModifierIncludesAll {
			body += " LEFT JOIN scenes_tags AS scene_tags_join ON scene_tags_join.scene_id = scene.id AND scene_tags_join.tag_id IN " + getInBinding(length)

			// only one required for include any
			requiredCount := 1

			// all required for include all
			if sceneTagsFilter.Modifier == CriterionModifierIncludesAll {
				requiredCount = length
			}

			havingClauses = append(havingClauses, "COUNT(DISTINCT scene_tags_join.tag_id) >= "+strconv.Itoa(requiredCount))
		} else if sceneTagsFilter.Modifier == CriterionModifierExcludes {
			// excludes all of the provided ids
			whereClauses = append(whereClauses, "not exists (select st.scene_id from scenes_tags as st where st.scene_id = scene.id AND st.tag_id IN "+getInBinding(length)+")")
		}

		for _, tagID := range sceneTagsFilter.Value {
			args = append(args, tagID)
		}
	}

	if performersFilter := sceneMarkerFilter.Performers; performersFilter != nil && len(performersFilter.Value) > 0 {
		length := len(performersFilter.Value)

		if performersFilter.Modifier == CriterionModifierIncludes || performersFilter.Modifier == CriterionModifierIncludesAll {
			body += " LEFT JOIN performers_scenes as scene_performers ON scene.id = scene_performers.scene_id"
			whereClauses = append(whereClauses, "scene_performers.performer_id IN "+getInBinding(length))

			// only one required for include any
			requiredCount := 1

			// all required for include all
			if performersFilter.Modifier == CriterionModifierIncludesAll {
				requiredCount = length
			}

			havingClauses = append(havingClauses, "COUNT(DISTINCT scene_performers.performer_id) >= "+strconv.Itoa(requiredCount))
		} else if performersFilter.Modifier == CriterionModifierExcludes {
			// excludes all of the provided ids
			whereClauses = append(whereClauses, "not exists (select sp.scene_id from performers_scenes as sp where sp.scene_id = scene.id AND sp.performer_id IN "+getInBinding(length)+")")
		}

		for _, performerID := range performersFilter.Value {
			args = append(args, performerID)
		}
	}

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"scene_markers.title", "scene.title"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		whereClauses = append(whereClauses, clause)
		args = append(args, thisArgs...)
	}

	if tagID := sceneMarkerFilter.TagID; tagID != nil {
		whereClauses = append(whereClauses, "(scene_markers.primary_tag_id = "+*tagID+" OR tags.id = "+*tagID+")")
	}

	sortAndPagination := qb.getSceneMarkerSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("scene_markers", body, args, sortAndPagination, whereClauses, havingClauses)

	var sceneMarkers []*SceneMarker
	for _, id := range idsResult {
		sceneMarker, _ := qb.Find(id)
		sceneMarkers = append(sceneMarkers, sceneMarker)
	}

	return sceneMarkers, countResult
}

func (qb *SceneMarkerQueryBuilder) getSceneMarkerSort(findFilter *FindFilterType) string {
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()
	tableName := "scene_markers"
	if sort == "scenes_updated_at" {
		sort = "updated_at"
		tableName = "scene"
	}
	return getSort(sort, direction, tableName)
}

func (qb *SceneMarkerQueryBuilder) querySceneMarkers(query string, args []interface{}, tx *sqlx.Tx) ([]*SceneMarker, error) {
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, args...)
	} else {
		rows, err = database.DB.Queryx(query, args...)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	sceneMarkers := make([]*SceneMarker, 0)
	for rows.Next() {
		sceneMarker := SceneMarker{}
		if err := rows.StructScan(&sceneMarker); err != nil {
			return nil, err
		}
		sceneMarkers = append(sceneMarkers, &sceneMarker)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sceneMarkers, nil
}

func (qb *SceneMarkerQueryBuilder) queryMarkerStringsResultType(query string, args []interface{}) ([]*MarkerStringsResultType, error) {
	rows, err := database.DB.Queryx(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	markerStrings := make([]*MarkerStringsResultType, 0)
	for rows.Next() {
		markerString := MarkerStringsResultType{}
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
