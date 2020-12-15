package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const sceneTable = "scenes"
const sceneIDColumn = "scene_id"
const performersScenesTable = "performers_scenes"
const scenesTagsTable = "scenes_tags"
const moviesScenesTable = "movies_scenes"

var scenesForPerformerQuery = selectAll(sceneTable) + `
LEFT JOIN performers_scenes as performers_join on performers_join.scene_id = scenes.id
WHERE performers_join.performer_id = ?
GROUP BY scenes.id
`

var countScenesForPerformerQuery = `
SELECT performer_id FROM performers_scenes as performers_join
WHERE performer_id = ?
GROUP BY scene_id
`

var scenesForStudioQuery = selectAll(sceneTable) + `
JOIN studios ON studios.id = scenes.studio_id
WHERE studios.id = ?
GROUP BY scenes.id
`
var scenesForMovieQuery = selectAll(sceneTable) + `
LEFT JOIN movies_scenes as movies_join on movies_join.scene_id = scenes.id
WHERE movies_join.movie_id = ?
GROUP BY scenes.id
`

var countScenesForTagQuery = `
SELECT tag_id AS id FROM scenes_tags
WHERE scenes_tags.tag_id = ?
GROUP BY scenes_tags.scene_id
`

var countScenesForMissingChecksumQuery = `
SELECT id FROM scenes
WHERE scenes.checksum is null
`

var countScenesForMissingOSHashQuery = `
SELECT id FROM scenes
WHERE scenes.oshash is null
`

type SceneQueryBuilder struct{}

func NewSceneQueryBuilder() SceneQueryBuilder {
	return SceneQueryBuilder{}
}

func sceneConstructor() interface{} {
	return &models.Scene{}
}

func (qb *SceneQueryBuilder) repository(tx *sqlx.Tx) *repository {
	return &repository{
		tx:          tx,
		tableName:   sceneTable,
		idColumn:    idColumn,
		constructor: sceneConstructor,
	}
}

func (qb *SceneQueryBuilder) Create(newObject models.Scene, tx *sqlx.Tx) (*models.Scene, error) {
	var ret models.Scene
	if err := qb.repository(tx).insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *SceneQueryBuilder) Update(updatedObject models.ScenePartial, tx *sqlx.Tx) (*models.Scene, error) {
	const partial = true
	if err := qb.repository(tx).update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.find(updatedObject.ID, tx)
}

func (qb *SceneQueryBuilder) UpdateFull(updatedObject models.Scene, tx *sqlx.Tx) (*models.Scene, error) {
	const partial = false
	if err := qb.repository(tx).update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.find(updatedObject.ID, tx)
}

func (qb *SceneQueryBuilder) UpdateFileModTime(id int, modTime models.NullSQLiteTimestamp, tx *sqlx.Tx) error {
	return qb.repository(tx).updateMap(id, map[string]interface{}{
		"file_mod_time": modTime,
	})
}

func (qb *SceneQueryBuilder) IncrementOCounter(id int, tx *sqlx.Tx) (int, error) {
	ensureTx(tx)
	_, err := tx.Exec(
		`UPDATE scenes SET o_counter = o_counter + 1 WHERE scenes.id = ?`,
		id,
	)
	if err != nil {
		return 0, err
	}

	scene, err := qb.find(id, tx)
	if err != nil {
		return 0, err
	}

	return scene.OCounter, nil
}

func (qb *SceneQueryBuilder) DecrementOCounter(id int, tx *sqlx.Tx) (int, error) {
	ensureTx(tx)
	_, err := tx.Exec(
		`UPDATE scenes SET o_counter = o_counter - 1 WHERE scenes.id = ? and scenes.o_counter > 0`,
		id,
	)
	if err != nil {
		return 0, err
	}

	scene, err := qb.find(id, tx)
	if err != nil {
		return 0, err
	}

	return scene.OCounter, nil
}

func (qb *SceneQueryBuilder) ResetOCounter(id int, tx *sqlx.Tx) (int, error) {
	ensureTx(tx)
	_, err := tx.Exec(
		`UPDATE scenes SET o_counter = 0 WHERE scenes.id = ?`,
		id,
	)
	if err != nil {
		return 0, err
	}

	scene, err := qb.find(id, tx)
	if err != nil {
		return 0, err
	}

	return scene.OCounter, nil
}

func (qb *SceneQueryBuilder) Destroy(id int, tx *sqlx.Tx) error {
	// delete all related table rows
	// TODO - this should be handled by a delete cascade
	if err := qb.performersRepository(tx).destroy([]int{id}); err != nil {
		return err
	}

	// scene markers should be handled prior to calling destroy
	// galleries should be handled prior to calling destroy

	return qb.repository(tx).destroyExisting([]int{id})
}

func (qb *SceneQueryBuilder) Find(id int) (*models.Scene, error) {
	return qb.find(id, nil)
}

func (qb *SceneQueryBuilder) FindMany(ids []int) ([]*models.Scene, error) {
	var scenes []*models.Scene
	for _, id := range ids {
		scene, err := qb.Find(id)
		if err != nil {
			return nil, err
		}

		if scene == nil {
			return nil, fmt.Errorf("scene with id %d not found", id)
		}

		scenes = append(scenes, scene)
	}

	return scenes, nil
}

func (qb *SceneQueryBuilder) find(id int, tx *sqlx.Tx) (*models.Scene, error) {
	var ret models.Scene
	if err := qb.repository(tx).get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *SceneQueryBuilder) FindByChecksum(checksum string) (*models.Scene, error) {
	query := "SELECT * FROM scenes WHERE checksum = ? LIMIT 1"
	args := []interface{}{checksum}
	return qb.queryScene(query, args, nil)
}

func (qb *SceneQueryBuilder) FindByOSHash(oshash string) (*models.Scene, error) {
	query := "SELECT * FROM scenes WHERE oshash = ? LIMIT 1"
	args := []interface{}{oshash}
	return qb.queryScene(query, args, nil)
}

func (qb *SceneQueryBuilder) FindByPath(path string) (*models.Scene, error) {
	query := selectAll(sceneTable) + "WHERE path = ? LIMIT 1"
	args := []interface{}{path}
	return qb.queryScene(query, args, nil)
}

func (qb *SceneQueryBuilder) FindByPerformerID(performerID int) ([]*models.Scene, error) {
	args := []interface{}{performerID}
	return qb.queryScenes(scenesForPerformerQuery, args, nil)
}

func (qb *SceneQueryBuilder) CountByPerformerID(performerID int) (int, error) {
	args := []interface{}{performerID}
	return runCountQuery(buildCountQuery(countScenesForPerformerQuery), args)
}

func (qb *SceneQueryBuilder) FindByStudioID(studioID int) ([]*models.Scene, error) {
	args := []interface{}{studioID}
	return qb.queryScenes(scenesForStudioQuery, args, nil)
}

func (qb *SceneQueryBuilder) FindByMovieID(movieID int) ([]*models.Scene, error) {
	args := []interface{}{movieID}
	return qb.queryScenes(scenesForMovieQuery, args, nil)
}

func (qb *SceneQueryBuilder) CountByMovieID(movieID int) (int, error) {
	args := []interface{}{movieID}
	return runCountQuery(buildCountQuery(scenesForMovieQuery), args)
}

func (qb *SceneQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT scenes.id FROM scenes"), nil)
}

func (qb *SceneQueryBuilder) Size() (float64, error) {
	return runSumQuery("SELECT SUM(cast(size as double)) as sum FROM scenes", nil)
}

func (qb *SceneQueryBuilder) CountByStudioID(studioID int) (int, error) {
	args := []interface{}{studioID}
	return runCountQuery(buildCountQuery(scenesForStudioQuery), args)
}

func (qb *SceneQueryBuilder) CountByTagID(tagID int) (int, error) {
	args := []interface{}{tagID}
	return runCountQuery(buildCountQuery(countScenesForTagQuery), args)
}

// CountMissingChecksum returns the number of scenes missing a checksum value.
func (qb *SceneQueryBuilder) CountMissingChecksum() (int, error) {
	return runCountQuery(buildCountQuery(countScenesForMissingChecksumQuery), []interface{}{})
}

// CountMissingOSHash returns the number of scenes missing an oshash value.
func (qb *SceneQueryBuilder) CountMissingOSHash() (int, error) {
	return runCountQuery(buildCountQuery(countScenesForMissingOSHashQuery), []interface{}{})
}

func (qb *SceneQueryBuilder) Wall(q *string) ([]*models.Scene, error) {
	s := ""
	if q != nil {
		s = *q
	}
	query := selectAll(sceneTable) + "WHERE scenes.details LIKE '%" + s + "%' ORDER BY RANDOM() LIMIT 80"
	return qb.queryScenes(query, nil, nil)
}

func (qb *SceneQueryBuilder) All() ([]*models.Scene, error) {
	return qb.queryScenes(selectAll(sceneTable)+qb.getSceneSort(nil), nil, nil)
}

func (qb *SceneQueryBuilder) Query(sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) ([]*models.Scene, int) {
	if sceneFilter == nil {
		sceneFilter = &models.SceneFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := queryBuilder{
		tableName: sceneTable,
	}

	query.body = selectDistinctIDs(sceneTable)
	query.body += `
		left join scene_markers on scene_markers.scene_id = scenes.id
		left join performers_scenes as performers_join on performers_join.scene_id = scenes.id
		left join movies_scenes as movies_join on movies_join.scene_id = scenes.id
		left join studios as studio on studio.id = scenes.studio_id
		left join galleries as gallery on gallery.scene_id = scenes.id
		left join scenes_tags as tags_join on tags_join.scene_id = scenes.id
		left join scene_stash_ids on scene_stash_ids.scene_id = scenes.id
	`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"scenes.title", "scenes.details", "scenes.path", "scenes.oshash", "scenes.checksum", "scene_markers.title"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		query.addWhere(clause)
		query.addArg(thisArgs...)
	}

	query.handleStringCriterionInput(sceneFilter.Path, "scenes.path")
	query.handleIntCriterionInput(sceneFilter.Rating, "scenes.rating")
	query.handleIntCriterionInput(sceneFilter.OCounter, "scenes.o_counter")

	if Organized := sceneFilter.Organized; Organized != nil {
		var organized string
		if *Organized == true {
			organized = "1"
		} else {
			organized = "0"
		}
		query.addWhere("scenes.organized = " + organized)
	}

	if durationFilter := sceneFilter.Duration; durationFilter != nil {
		clause, thisArgs := getDurationWhereClause(*durationFilter)
		query.addWhere(clause)
		query.addArg(thisArgs...)
	}

	if resolutionFilter := sceneFilter.Resolution; resolutionFilter != nil {
		if resolution := resolutionFilter.String(); resolutionFilter.IsValid() {
			switch resolution {
			case "LOW":
				query.addWhere("scenes.height < 480")
			case "STANDARD":
				query.addWhere("(scenes.height >= 480 AND scenes.height < 720)")
			case "STANDARD_HD":
				query.addWhere("(scenes.height >= 720 AND scenes.height < 1080)")
			case "FULL_HD":
				query.addWhere("(scenes.height >= 1080 AND scenes.height < 2160)")
			case "FOUR_K":
				query.addWhere("scenes.height >= 2160")
			}
		}
	}

	if hasMarkersFilter := sceneFilter.HasMarkers; hasMarkersFilter != nil {
		if strings.Compare(*hasMarkersFilter, "true") == 0 {
			query.addHaving("count(scene_markers.scene_id) > 0")
		} else {
			query.addWhere("scene_markers.id IS NULL")
		}
	}

	if isMissingFilter := sceneFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "gallery":
			query.addWhere("gallery.scene_id IS NULL")
		case "studio":
			query.addWhere("scenes.studio_id IS NULL")
		case "movie":
			query.addWhere("movies_join.scene_id IS NULL")
		case "performers":
			query.addWhere("performers_join.scene_id IS NULL")
		case "date":
			query.addWhere("scenes.date IS \"\" OR scenes.date IS \"0001-01-01\"")
		case "tags":
			query.addWhere("tags_join.scene_id IS NULL")
		case "stash_id":
			query.addWhere("scene_stash_ids.scene_id IS NULL")
		default:
			query.addWhere("(scenes." + *isMissingFilter + " IS NULL OR TRIM(scenes." + *isMissingFilter + ") = '')")
		}
	}

	if tagsFilter := sceneFilter.Tags; tagsFilter != nil && len(tagsFilter.Value) > 0 {
		for _, tagID := range tagsFilter.Value {
			query.addArg(tagID)
		}

		query.body += " LEFT JOIN tags on tags_join.tag_id = tags.id"
		whereClause, havingClause := getMultiCriterionClause("scenes", "tags", "scenes_tags", "scene_id", "tag_id", tagsFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if performersFilter := sceneFilter.Performers; performersFilter != nil && len(performersFilter.Value) > 0 {
		for _, performerID := range performersFilter.Value {
			query.addArg(performerID)
		}

		query.body += " LEFT JOIN performers ON performers_join.performer_id = performers.id"
		whereClause, havingClause := getMultiCriterionClause("scenes", "performers", "performers_scenes", "scene_id", "performer_id", performersFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if studiosFilter := sceneFilter.Studios; studiosFilter != nil && len(studiosFilter.Value) > 0 {
		for _, studioID := range studiosFilter.Value {
			query.addArg(studioID)
		}

		whereClause, havingClause := getMultiCriterionClause("scenes", "studio", "", "", "studio_id", studiosFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if moviesFilter := sceneFilter.Movies; moviesFilter != nil && len(moviesFilter.Value) > 0 {
		for _, movieID := range moviesFilter.Value {
			query.addArg(movieID)
		}

		query.body += " LEFT JOIN movies ON movies_join.movie_id = movies.id"
		whereClause, havingClause := getMultiCriterionClause("scenes", "movies", "movies_scenes", "scene_id", "movie_id", moviesFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if stashIDFilter := sceneFilter.StashID; stashIDFilter != nil {
		query.addWhere("scene_stash_ids.stash_id = ?")
		query.addArg(stashIDFilter)
	}

	query.sortAndPagination = qb.getSceneSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := query.executeFind()

	var scenes []*models.Scene
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

func getDurationWhereClause(durationFilter models.IntCriterionInput) (string, []interface{}) {
	// special case for duration. We accept duration as seconds as int but the
	// field is floating point. Change the equals filter to return a range
	// between x and x + 1
	// likewise, not equals needs to be duration < x OR duration >= x
	var clause string
	args := []interface{}{}

	value := durationFilter.Value
	if durationFilter.Modifier == models.CriterionModifierEquals {
		clause = "scenes.duration >= ? AND scenes.duration < ?"
		args = append(args, value)
		args = append(args, value+1)
	} else if durationFilter.Modifier == models.CriterionModifierNotEquals {
		clause = "(scenes.duration < ? OR scenes.duration >= ?)"
		args = append(args, value)
		args = append(args, value+1)
	} else {
		var count int
		clause, count = getIntCriterionWhereClause("scenes.duration", durationFilter)
		if count == 1 {
			args = append(args, value)
		}
	}

	return clause, args
}

func (qb *SceneQueryBuilder) QueryAllByPathRegex(regex string, ignoreOrganized bool) ([]*models.Scene, error) {
	var args []interface{}
	body := selectDistinctIDs("scenes") + " WHERE scenes.path regexp ?"

	if ignoreOrganized {
		body += " AND scenes.organized = 0"
	}

	args = append(args, "(?i)"+regex)

	idsResult, err := runIdsQuery(body, args)

	if err != nil {
		return nil, err
	}

	var scenes []*models.Scene
	for _, id := range idsResult {
		scene, err := qb.Find(id)

		if err != nil {
			return nil, err
		}

		scenes = append(scenes, scene)
	}

	return scenes, nil
}

func (qb *SceneQueryBuilder) QueryByPathRegex(findFilter *models.FindFilterType) ([]*models.Scene, int) {
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	var whereClauses []string
	var havingClauses []string
	var args []interface{}
	body := selectDistinctIDs("scenes")

	if q := findFilter.Q; q != nil && *q != "" {
		whereClauses = append(whereClauses, "scenes.path regexp ?")
		args = append(args, "(?i)"+*q)
	}

	sortAndPagination := qb.getSceneSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("scenes", body, args, sortAndPagination, whereClauses, havingClauses)

	var scenes []*models.Scene
	for _, id := range idsResult {
		scene, _ := qb.Find(id)
		scenes = append(scenes, scene)
	}

	return scenes, countResult
}

func (qb *SceneQueryBuilder) getSceneSort(findFilter *models.FindFilterType) string {
	if findFilter == nil {
		return " ORDER BY scenes.path, scenes.date ASC "
	}
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()
	return getSort(sort, direction, "scenes")
}

func (qb *SceneQueryBuilder) queryScene(query string, args []interface{}, tx *sqlx.Tx) (*models.Scene, error) {
	results, err := qb.queryScenes(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *SceneQueryBuilder) queryScenes(query string, args []interface{}, tx *sqlx.Tx) ([]*models.Scene, error) {
	var ret models.Scenes
	if err := qb.repository(tx).query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Scene(ret), nil
}

func (qb *SceneQueryBuilder) UpdateFormat(id int, format string, tx *sqlx.Tx) error {
	ensureTx(tx)

	return qb.repository(tx).updateMap(id, map[string]interface{}{
		"format": format,
	})
}

func (qb *SceneQueryBuilder) UpdateOSHash(id int, oshash string, tx *sqlx.Tx) error {
	ensureTx(tx)

	return qb.repository(tx).updateMap(id, map[string]interface{}{
		"oshash": oshash,
	})
}

func (qb *SceneQueryBuilder) UpdateChecksum(id int, checksum string, tx *sqlx.Tx) error {
	ensureTx(tx)

	return qb.repository(tx).updateMap(id, map[string]interface{}{
		"checksum": checksum,
	})
}

func (qb *SceneQueryBuilder) UpdateSceneCover(sceneID int, cover []byte, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing cover and then create new
	if err := qb.DestroySceneCover(sceneID, tx); err != nil {
		return err
	}

	_, err := tx.Exec(
		`INSERT INTO scenes_cover (scene_id, cover) VALUES (?, ?)`,
		sceneID,
		cover,
	)

	return err
}

func (qb *SceneQueryBuilder) DestroySceneCover(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM scenes_cover WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return err
}

func (qb *SceneQueryBuilder) GetSceneCover(sceneID int, tx *sqlx.Tx) ([]byte, error) {
	query := `SELECT cover from scenes_cover WHERE scene_id = ?`
	return getImage(tx, query, sceneID)
}

func (qb *SceneQueryBuilder) moviesRepository(tx *sqlx.Tx) *repository {
	return &repository{
		tx:        tx,
		tableName: moviesScenesTable,
		idColumn:  sceneIDColumn,
	}
}

func (qb *SceneQueryBuilder) GetMovies(id int, tx *sqlx.Tx) (ret []models.MoviesScenes, err error) {
	if err := qb.moviesRepository(tx).getAll(id, func(rows *sqlx.Rows) error {
		var ms models.MoviesScenes
		if err := rows.StructScan(&ms); err != nil {
			return err
		}

		ret = append(ret, ms)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *SceneQueryBuilder) UpdateMovies(sceneID int, movies []models.MoviesScenes, tx *sqlx.Tx) error {
	// destroy existing joins
	r := qb.moviesRepository(tx)
	if err := r.destroy([]int{sceneID}); err != nil {
		return err
	}

	for _, m := range movies {
		m.SceneID = sceneID
		if _, err := r.insert(m); err != nil {
			return err
		}
	}

	return nil
}

func (qb *SceneQueryBuilder) performersRepository(tx *sqlx.Tx) *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        tx,
			tableName: performersScenesTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn: performerIDColumn,
	}
}

func (qb *SceneQueryBuilder) GetPerformerIDs(id int, tx *sqlx.Tx) ([]int, error) {
	return qb.performersRepository(tx).getIDs(id)
}

func (qb *SceneQueryBuilder) UpdatePerformers(id int, performerIDs []int, tx *sqlx.Tx) error {
	// Delete the existing joins and then create new ones
	return qb.performersRepository(tx).replace(id, performerIDs)
}

func (qb *SceneQueryBuilder) tagsRepository(tx *sqlx.Tx) *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        tx,
			tableName: scenesTagsTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn: tagIDColumn,
	}
}

func (qb *SceneQueryBuilder) GetTagIDs(id int, tx *sqlx.Tx) ([]int, error) {
	return qb.tagsRepository(tx).getIDs(id)
}

func (qb *SceneQueryBuilder) UpdateTags(id int, tagIDs []int, tx *sqlx.Tx) error {
	// Delete the existing joins and then create new ones
	return qb.tagsRepository(tx).replace(id, tagIDs)
}

func (qb *SceneQueryBuilder) stashIDRepository(tx *sqlx.Tx) *stashIDRepository {
	return &stashIDRepository{
		repository{
			tx:        tx,
			tableName: "scene_stash_ids",
			idColumn:  sceneIDColumn,
		},
	}
}

func (qb *SceneQueryBuilder) GetStashIDs(sceneID int, tx *sqlx.Tx) ([]*models.StashID, error) {
	return qb.stashIDRepository(tx).get(sceneID)
}

func (qb *SceneQueryBuilder) UpdateStashIDs(sceneID int, stashIDs []models.StashID, tx *sqlx.Tx) error {
	return qb.stashIDRepository(tx).replace(sceneID, stashIDs)
}

func NewSceneReaderWriter(tx *sqlx.Tx) *sceneReaderWriter {
	return &sceneReaderWriter{
		tx: tx,
		qb: NewSceneQueryBuilder(),
	}
}

type sceneReaderWriter struct {
	tx *sqlx.Tx
	qb SceneQueryBuilder
}

func (t *sceneReaderWriter) Find(id int) (*models.Scene, error) {
	return t.qb.Find(id)
}

func (t *sceneReaderWriter) FindMany(ids []int) ([]*models.Scene, error) {
	return t.qb.FindMany(ids)
}

func (t *sceneReaderWriter) FindByChecksum(checksum string) (*models.Scene, error) {
	return t.qb.FindByChecksum(checksum)
}

func (t *sceneReaderWriter) FindByOSHash(oshash string) (*models.Scene, error) {
	return t.qb.FindByOSHash(oshash)
}

func (t *sceneReaderWriter) FindByMovieID(movieID int) ([]*models.Scene, error) {
	return t.qb.FindByMovieID(movieID)
}

func (t *sceneReaderWriter) CountByMovieID(movieID int) (int, error) {
	return t.qb.CountByMovieID(movieID)
}

func (t *sceneReaderWriter) FindByPerformerID(performerID int) ([]*models.Scene, error) {
	return t.qb.FindByPerformerID(performerID)
}

func (t *sceneReaderWriter) CountByPerformerID(performerID int) (int, error) {
	return t.qb.CountByPerformerID(performerID)
}

func (t *sceneReaderWriter) CountByStudioID(studioID int) (int, error) {
	return t.qb.CountByStudioID(studioID)
}

func (t *sceneReaderWriter) CountByTagID(tagID int) (int, error) {
	return t.qb.CountByTagID(tagID)
}

func (t *sceneReaderWriter) All() ([]*models.Scene, error) {
	return t.qb.All()
}

func (t *sceneReaderWriter) GetSceneCover(sceneID int) ([]byte, error) {
	return t.qb.GetSceneCover(sceneID, t.tx)
}

func (t *sceneReaderWriter) Create(newScene models.Scene) (*models.Scene, error) {
	return t.qb.Create(newScene, t.tx)
}

func (t *sceneReaderWriter) Update(updatedScene models.ScenePartial) (*models.Scene, error) {
	return t.qb.Update(updatedScene, t.tx)
}

func (t *sceneReaderWriter) UpdateFull(updatedScene models.Scene) (*models.Scene, error) {
	return t.qb.UpdateFull(updatedScene, t.tx)
}

func (t *sceneReaderWriter) Destroy(id int) error {
	return t.qb.Destroy(id, t.tx)
}

func (t *sceneReaderWriter) IncrementOCounter(id int) (int, error) {
	return t.qb.IncrementOCounter(id, t.tx)
}

func (t *sceneReaderWriter) DecrementOCounter(id int) (int, error) {
	return t.qb.DecrementOCounter(id, t.tx)
}

func (t *sceneReaderWriter) ResetOCounter(id int) (int, error) {
	return t.qb.ResetOCounter(id, t.tx)
}

func (t *sceneReaderWriter) UpdateSceneCover(sceneID int, cover []byte) error {
	return t.qb.UpdateSceneCover(sceneID, cover, t.tx)
}

func (t *sceneReaderWriter) GetMovies(sceneID int) ([]models.MoviesScenes, error) {
	return t.qb.GetMovies(sceneID, t.tx)
}

func (t *sceneReaderWriter) UpdateMovies(sceneID int, movies []models.MoviesScenes) error {
	return t.qb.UpdateMovies(sceneID, movies, t.tx)
}

func (t *sceneReaderWriter) GetPerformerIDs(sceneID int) ([]int, error) {
	return t.qb.GetPerformerIDs(sceneID, t.tx)
}

func (t *sceneReaderWriter) UpdatePerformers(sceneID int, performerIDs []int) error {
	return t.qb.UpdatePerformers(sceneID, performerIDs, t.tx)
}

func (t *sceneReaderWriter) GetTagIDs(sceneID int) ([]int, error) {
	return t.qb.GetTagIDs(sceneID, t.tx)
}

func (t *sceneReaderWriter) UpdateTags(sceneID int, tagIDs []int) error {
	return t.qb.UpdateTags(sceneID, tagIDs, t.tx)
}

func (t *sceneReaderWriter) GetStashIDs(sceneID int) ([]*models.StashID, error) {
	return t.qb.GetStashIDs(sceneID, t.tx)
}

func (t *sceneReaderWriter) UpdateStashIDs(sceneID int, stashIDs []models.StashID) error {
	return t.qb.UpdateStashIDs(sceneID, stashIDs, t.tx)
}
