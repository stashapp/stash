package models

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

const scenesForPerformerQuery = `
SELECT scenes.* FROM scenes
LEFT JOIN performers_scenes as performers_join on performers_join.scene_id = scenes.id
LEFT JOIN performers on performers_join.performer_id = performers.id
WHERE performers.id = ?
GROUP BY scenes.id
`

const scenesForStudioQuery = `
SELECT scenes.* FROM scenes
JOIN studios ON studios.id = scenes.studio_id
WHERE studios.id = ?
GROUP BY scenes.id
`

const scenesForTagQuery = `
SELECT scenes.* FROM scenes
LEFT JOIN scenes_tags as tags_join on tags_join.scene_id = scenes.id
LEFT JOIN tags on tags_join.tag_id = tags.id
WHERE tags.id = ?
GROUP BY scenes.id
`

type SceneQueryBuilder struct{}

func NewSceneQueryBuilder() SceneQueryBuilder {
	return SceneQueryBuilder{}
}

func (qb *SceneQueryBuilder) Create(newScene Scene, tx *sqlx.Tx) (*Scene, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO scenes (checksum, path, title, details, url, date, rating, size, duration, video_codec,
                    			    audio_codec, width, height, framerate, bitrate, studio_id, created_at, updated_at)
				VALUES (:checksum, :path, :title, :details, :url, :date, :rating, :size, :duration, :video_codec,
				        :audio_codec, :width, :height, :framerate, :bitrate, :studio_id, :created_at, :updated_at)
		`,
		newScene,
	)
	if err != nil {
		return nil, err
	}
	sceneID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err := tx.Get(&newScene, `SELECT * FROM scenes WHERE id = ? LIMIT 1`, sceneID); err != nil {
		return nil, err
	}
	return &newScene, nil
}

func (qb *SceneQueryBuilder) Update(updatedScene ScenePartial, tx *sqlx.Tx) (*Scene, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE scenes SET `+SQLGenKeysPartial(updatedScene)+` WHERE scenes.id = :id`,
		updatedScene,
	)
	if err != nil {
		return nil, err
	}

	return qb.find(updatedScene.ID, tx)
}

func (qb *SceneQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	return executeDeleteQuery("scenes", id, tx)
}
func (qb *SceneQueryBuilder) Find(id int) (*Scene, error) {
	return qb.find(id, nil)
}

func (qb *SceneQueryBuilder) find(id int, tx *sqlx.Tx) (*Scene, error) {
	query := "SELECT * FROM scenes WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryScene(query, args, tx)
}

func (qb *SceneQueryBuilder) FindByChecksum(checksum string) (*Scene, error) {
	query := "SELECT * FROM scenes WHERE checksum = ? LIMIT 1"
	args := []interface{}{checksum}
	return qb.queryScene(query, args, nil)
}

func (qb *SceneQueryBuilder) FindByPath(path string) (*Scene, error) {
	query := "SELECT * FROM scenes WHERE path = ? LIMIT 1"
	args := []interface{}{path}
	return qb.queryScene(query, args, nil)
}

func (qb *SceneQueryBuilder) FindByPerformerID(performerID int) ([]*Scene, error) {
	args := []interface{}{performerID}
	return qb.queryScenes(scenesForPerformerQuery, args, nil)
}

func (qb *SceneQueryBuilder) CountByPerformerID(performerID int) (int, error) {
	args := []interface{}{performerID}
	return runCountQuery(buildCountQuery(scenesForPerformerQuery), args)
}

func (qb *SceneQueryBuilder) FindByStudioID(studioID int) ([]*Scene, error) {
	args := []interface{}{studioID}
	return qb.queryScenes(scenesForStudioQuery, args, nil)
}

func (qb *SceneQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT scenes.id FROM scenes"), nil)
}

func (qb *SceneQueryBuilder) CountByStudioID(studioID int) (int, error) {
	args := []interface{}{studioID}
	return runCountQuery(buildCountQuery(scenesForStudioQuery), args)
}

func (qb *SceneQueryBuilder) CountByTagID(tagID int) (int, error) {
	args := []interface{}{tagID}
	return runCountQuery(buildCountQuery(scenesForTagQuery), args)
}

func (qb *SceneQueryBuilder) Wall(q *string) ([]*Scene, error) {
	s := ""
	if q != nil {
		s = *q
	}
	query := "SELECT scenes.* FROM scenes WHERE scenes.details LIKE '%" + s + "%' ORDER BY RANDOM() LIMIT 80"
	return qb.queryScenes(query, nil, nil)
}

func (qb *SceneQueryBuilder) All() ([]*Scene, error) {
	return qb.queryScenes(selectAll("scenes")+qb.getSceneSort(nil), nil, nil)
}

func (qb *SceneQueryBuilder) Query(sceneFilter *SceneFilterType, findFilter *FindFilterType) ([]*Scene, int) {
	if sceneFilter == nil {
		sceneFilter = &SceneFilterType{}
	}
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	var whereClauses []string
	var havingClauses []string
	var args []interface{}
	body := selectDistinctIDs("scenes")
	body = body + `
		left join scene_markers on scene_markers.scene_id = scenes.id
		left join performers_scenes as performers_join on performers_join.scene_id = scenes.id
		left join performers on performers_join.performer_id = performers.id
		left join studios as studio on studio.id = scenes.studio_id
		left join galleries as gallery on gallery.scene_id = scenes.id
		left join scenes_tags as tags_join on tags_join.scene_id = scenes.id
		left join tags on tags_join.tag_id = tags.id
	`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"scenes.title", "scenes.details", "scenes.path", "scenes.checksum", "scene_markers.title"}
		whereClauses = append(whereClauses, getSearch(searchColumns, *q))
	}

	if rating := sceneFilter.Rating; rating != nil {
		clause, count := getIntCriterionWhereClause("scenes.rating", *sceneFilter.Rating)
		whereClauses = append(whereClauses, clause)
		if count == 1 {
			args = append(args, sceneFilter.Rating.Value)
		}
	}

	if durationFilter := sceneFilter.Duration; durationFilter != nil {
		clause, count := getIntCriterionWhereClause("scenes.duration", *durationFilter)
		whereClauses = append(whereClauses, clause)
		if count == 1 {
			args = append(args, durationFilter.Value)
		}
	}

	if resolutionFilter := sceneFilter.Resolution; resolutionFilter != nil {
		if resolution := resolutionFilter.String(); resolutionFilter.IsValid() {
			switch resolution {
			case "LOW":
				whereClauses = append(whereClauses, "(scenes.height >= 240 AND scenes.height < 480)")
			case "STANDARD":
				whereClauses = append(whereClauses, "(scenes.height >= 480 AND scenes.height < 720)")
			case "STANDARD_HD":
				whereClauses = append(whereClauses, "(scenes.height >= 720 AND scenes.height < 1080)")
			case "FULL_HD":
				whereClauses = append(whereClauses, "(scenes.height >= 1080 AND scenes.height < 2160)")
			case "FOUR_K":
				whereClauses = append(whereClauses, "scenes.height >= 2160")
			default:
				whereClauses = append(whereClauses, "scenes.height < 240")
			}
		}
	}

	if hasMarkersFilter := sceneFilter.HasMarkers; hasMarkersFilter != nil {
		if strings.Compare(*hasMarkersFilter, "true") == 0 {
			havingClauses = append(havingClauses, "count(scene_markers.scene_id) > 0")
		} else {
			whereClauses = append(whereClauses, "scene_markers.id IS NULL")
		}
	}

	if isMissingFilter := sceneFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "gallery":
			whereClauses = append(whereClauses, "gallery.scene_id IS NULL")
		case "studio":
			whereClauses = append(whereClauses, "scenes.studio_id IS NULL")
		case "performers":
			whereClauses = append(whereClauses, "performers_join.scene_id IS NULL")
		case "date":
			whereClauses = append(whereClauses, "scenes.date IS \"\" OR scenes.date IS \"0001-01-01\"")
		default:
			whereClauses = append(whereClauses, "scenes."+*isMissingFilter+" IS NULL")
		}
	}

	if tagsFilter := sceneFilter.Tags; tagsFilter != nil && len(tagsFilter.Value) > 0 {
		for _, tagID := range tagsFilter.Value {
			args = append(args, tagID)
		}

		whereClause, havingClause := getMultiCriterionClause("tags", "scenes_tags", "tag_id", tagsFilter)
		whereClauses = appendClause(whereClauses, whereClause)
		havingClauses = appendClause(havingClauses, havingClause)
	}

	if performersFilter := sceneFilter.Performers; performersFilter != nil && len(performersFilter.Value) > 0 {
		for _, performerID := range performersFilter.Value {
			args = append(args, performerID)
		}

		whereClause, havingClause := getMultiCriterionClause("performers", "performers_scenes", "performer_id", performersFilter)
		whereClauses = appendClause(whereClauses, whereClause)
		havingClauses = appendClause(havingClauses, havingClause)
	}

	if studiosFilter := sceneFilter.Studios; studiosFilter != nil && len(studiosFilter.Value) > 0 {
		for _, studioID := range studiosFilter.Value {
			args = append(args, studioID)
		}

		whereClause, havingClause := getMultiCriterionClause("studio", "", "studio_id", studiosFilter)
		whereClauses = appendClause(whereClauses, whereClause)
		havingClauses = appendClause(havingClauses, havingClause)
	}

	sortAndPagination := qb.getSceneSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("scenes", body, args, sortAndPagination, whereClauses, havingClauses)

	var scenes []*Scene
	for _, id := range idsResult {
		scene, _ := qb.Find(id)
		scenes = append(scenes, scene)
	}

	return scenes, countResult
}

func appendClause(clauses []string, clause string) []string {
	if clause != "" {
		return append(clauses, clause)
	}

	return clauses
}

// returns where clause and having clause
func getMultiCriterionClause(table string, joinTable string, joinTableField string, criterion *MultiCriterionInput) (string, string) {
	whereClause := ""
	havingClause := ""
	if criterion.Modifier == CriterionModifierIncludes {
		// includes any of the provided ids
		whereClause = table + ".id IN " + getInBinding(len(criterion.Value))
	} else if criterion.Modifier == CriterionModifierIncludesAll {
		// includes all of the provided ids
		whereClause = table + ".id IN " + getInBinding(len(criterion.Value))
		havingClause = "count(distinct " + table + ".id) IS " + strconv.Itoa(len(criterion.Value))
	} else if criterion.Modifier == CriterionModifierExcludes {
		// excludes all of the provided ids
		if joinTable != "" {
			whereClause = "not exists (select " + joinTable + ".scene_id from " + joinTable + " where " + joinTable + ".scene_id = scenes.id and " + joinTable + "." + joinTableField + " in " + getInBinding(len(criterion.Value)) + ")"
		} else {
			whereClause = "not exists (select s.id from scenes as s where s.id = scenes.id and s." + joinTableField + " in " + getInBinding(len(criterion.Value)) + ")"
		}
	}

	return whereClause, havingClause
}

func (qb *SceneQueryBuilder) QueryAllByPathRegex(regex string) ([]*Scene, error) {
	var args []interface{}
	body := selectDistinctIDs("scenes") + " WHERE scenes.path regexp '(?i)" + regex + "'"

	idsResult, err := runIdsQuery(body, args)

	if err != nil {
		return nil, err
	}

	var scenes []*Scene
	for _, id := range idsResult {
		scene, err := qb.Find(id)

		if err != nil {
			return nil, err
		}

		scenes = append(scenes, scene)
	}

	return scenes, nil
}

func (qb *SceneQueryBuilder) QueryByPathRegex(findFilter *FindFilterType) ([]*Scene, int) {
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	var whereClauses []string
	var havingClauses []string
	var args []interface{}
	body := selectDistinctIDs("scenes")

	if q := findFilter.Q; q != nil && *q != "" {
		whereClauses = append(whereClauses, "scenes.path regexp '(?i)"+*q+"'")
	}

	sortAndPagination := qb.getSceneSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("scenes", body, args, sortAndPagination, whereClauses, havingClauses)

	var scenes []*Scene
	for _, id := range idsResult {
		scene, _ := qb.Find(id)
		scenes = append(scenes, scene)
	}

	return scenes, countResult
}

func (qb *SceneQueryBuilder) getSceneSort(findFilter *FindFilterType) string {
	if findFilter == nil {
		return " ORDER BY scenes.path, scenes.date ASC "
	}
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()
	return getSort(sort, direction, "scenes")
}

func (qb *SceneQueryBuilder) queryScene(query string, args []interface{}, tx *sqlx.Tx) (*Scene, error) {
	results, err := qb.queryScenes(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *SceneQueryBuilder) queryScenes(query string, args []interface{}, tx *sqlx.Tx) ([]*Scene, error) {
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

	scenes := make([]*Scene, 0)
	for rows.Next() {
		scene := Scene{}
		if err := rows.StructScan(&scene); err != nil {
			return nil, err
		}
		scenes = append(scenes, &scene)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scenes, nil
}
