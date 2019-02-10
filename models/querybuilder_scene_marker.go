package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/database"
	"strconv"
)

const sceneMarkersForTagQuery = `
SELECT scene_markers.* FROM scene_markers
LEFT JOIN scene_markers_tags as tags_join on tags_join.scene_marker_id = scene_markers.id
LEFT JOIN tags on tags_join.tag_id = tags.id
WHERE tags.id = ?
GROUP BY scene_markers.id
`

type sceneMarkerQueryBuilder struct {}

func NewSceneMarkerQueryBuilder() sceneMarkerQueryBuilder {
	return sceneMarkerQueryBuilder{}
}

func (qb *sceneMarkerQueryBuilder) Create(newSceneMarker SceneMarker, tx *sqlx.Tx) (*SceneMarker, error) {
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

func (qb *sceneMarkerQueryBuilder) Update(updatedSceneMarker SceneMarker, tx *sqlx.Tx) (*SceneMarker, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE scene_markers SET `+SqlGenKeys(updatedSceneMarker)+` WHERE scene_markers.id = :id`,
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

func (qb *sceneMarkerQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	return executeDeleteQuery("scene_markers", id, tx)
}

func (qb *sceneMarkerQueryBuilder) Find(id int) (*SceneMarker, error) {
	query := "SELECT * FROM scene_markers WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	results, err := qb.querySceneMarkers(query, args, nil)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return &results[0], nil
}

func (qb *sceneMarkerQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]SceneMarker, error) {
	query := `
		SELECT scene_markers.* FROM scene_markers
		JOIN scenes ON scenes.id = scene_markers.scene_id
		WHERE scenes.id = ?
		GROUP BY scene_markers.id
		ORDER BY scene_markers.seconds ASC
	`
	args := []interface{}{sceneID}
	return qb.querySceneMarkers(query, args, tx)
}

func (qb *sceneMarkerQueryBuilder) CountByTagID(tagID int) (int, error) {
	args := []interface{}{tagID}
	return runCountQuery(buildCountQuery(sceneMarkersForTagQuery), args)
}

func (qb *sceneMarkerQueryBuilder) GetMarkerStrings(q *string, sort *string) ([]*MarkerStringsResultType, error) {
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
	args := []interface{}{}
	return qb.queryMarkerStringsResultType(query, args)
}

func (qb *sceneMarkerQueryBuilder) Wall(q *string) ([]SceneMarker, error) {
	s := ""
	if q != nil {
		s = *q
	}
	query := "SELECT scene_markers.* FROM scene_markers WHERE scene_markers.title LIKE '%" + s + "%' ORDER BY RANDOM() LIMIT 80"
	return qb.querySceneMarkers(query, nil, nil)
}

func (qb *sceneMarkerQueryBuilder) Query(sceneMarkerFilter *SceneMarkerFilterType, findFilter *FindFilterType) ([]SceneMarker, int) {
	if sceneMarkerFilter == nil {
		sceneMarkerFilter = &SceneMarkerFilterType{}
	}
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	whereClauses := []string{}
	havingClauses := []string{}
	args := []interface{}{}
	body := selectDistinctIDs("scene_markers")
	body = body + `
		left join tags as primary_tag on primary_tag.id = scene_markers.primary_tag_id
		left join scenes as scene on scene.id = scene_markers.scene_id
		left join scene_markers_tags as tags_join on tags_join.scene_marker_id = scene_markers.id
		left join tags on tags_join.tag_id = tags.id
	`

	if tagIDs := sceneMarkerFilter.Tags; tagIDs != nil {
		//select `scene_markers`.* from `scene_markers`
		//left join `tags` as `primary_tags_join`
		//  on `primary_tags_join`.`id` = `scene_markers`.`primary_tag_id`
		//  and `primary_tags_join`.`id` in ('3', '37', '9', '89')
		//left join `scene_markers_tags` as `tags_join`
		//  on `tags_join`.`scene_marker_id` = `scene_markers`.`id`
		//  and `tags_join`.`tag_id` in ('3', '37', '9', '89')
		//group by `scene_markers`.`id`
		//having ((count(distinct `primary_tags_join`.`id`) + count(distinct `tags_join`.`tag_id`)) = 4)

		length := len(tagIDs)
		body += " LEFT JOIN tags AS ptj ON ptj.id = scene_markers.primary_tag_id AND ptj.id IN " + getInBinding(length)
		body += " LEFT JOIN scene_markers_tags AS tj ON tj.scene_marker_id = scene_markers.id AND tj.tag_id IN " + getInBinding(length)
		havingClauses = append(havingClauses, "((COUNT(DISTINCT ptj.id) + COUNT(DISTINCT tj.tag_id)) = " + strconv.Itoa(length) +")")
		for _, tagID := range tagIDs {
			args = append(args, tagID)
		}
		for _, tagID := range tagIDs {
			args = append(args, tagID)
		}
	}

	if sceneTagIDs := sceneMarkerFilter.SceneTags; sceneTagIDs != nil {
		length := len(sceneTagIDs)
		body += " LEFT JOIN scenes_tags AS scene_tags_join ON scene_tags_join.scene_id = scene.id AND scene_tags_join.tag_id IN " + getInBinding(length)
		havingClauses = append(havingClauses, "COUNT(DISTINCT scene_tags_join.tag_id) = " + strconv.Itoa(length))
		for _, tagID := range sceneTagIDs {
			args = append(args, tagID)
		}
	}

	if performerIDs := sceneMarkerFilter.Performers; performerIDs != nil {
		length := len(performerIDs)
		body += " LEFT JOIN performers_scenes as scene_performers ON scene.id = scene_performers.scene_id"
		whereClauses = append(whereClauses, "scene_performers.performer_id IN " + getInBinding(length))
		for _, performerID := range performerIDs {
			args = append(args, performerID)
		}
	}

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"scene_markers.title", "scene.title"}
		whereClauses = append(whereClauses, getSearch(searchColumns, *q))
	}

	if tagID := sceneMarkerFilter.TagID; tagID != nil {
		whereClauses = append(whereClauses, "(scene_markers.primary_tag_id = "+*tagID+" OR tags.id = "+*tagID+")")
	}

	sortAndPagination := qb.getSceneMarkerSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("scene_markers", body, args, sortAndPagination, whereClauses, havingClauses)

	var sceneMarkers []SceneMarker
	for _, id := range idsResult {
		sceneMarker, _ := qb.Find(id)
		sceneMarkers = append(sceneMarkers, *sceneMarker)
	}

	return sceneMarkers, countResult
}

func (qb *sceneMarkerQueryBuilder) getSceneMarkerSort(findFilter *FindFilterType) string {
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()
	return getSort(sort, direction, "scene_markers")
}

func (qb *sceneMarkerQueryBuilder) querySceneMarkers(query string, args []interface{}, tx *sqlx.Tx) ([]SceneMarker, error) {
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

	sceneMarkers := make([]SceneMarker, 0)
	sceneMarker := SceneMarker{}
	for rows.Next() {
		if err := rows.StructScan(&sceneMarker); err != nil {
			return nil, err
		}
		sceneMarkers = append(sceneMarkers, sceneMarker)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sceneMarkers, nil
}

func (qb *sceneMarkerQueryBuilder) queryMarkerStringsResultType(query string, args []interface{}) ([]*MarkerStringsResultType, error) {
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