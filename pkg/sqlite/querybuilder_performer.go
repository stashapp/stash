package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const performerTable = "performers"

type PerformerQueryBuilder struct{}

func NewPerformerQueryBuilder() PerformerQueryBuilder {
	return PerformerQueryBuilder{}
}

func performerConstructor() interface{} {
	return &models.Performer{}
}

func (qb *PerformerQueryBuilder) repository(tx *sqlx.Tx) *repository {
	return &repository{
		tx:          tx,
		tableName:   performerTable,
		idColumn:    idColumn,
		constructor: performerConstructor,
	}
}

func (qb *PerformerQueryBuilder) Create(newObject models.Performer, tx *sqlx.Tx) (*models.Performer, error) {
	var ret models.Performer
	if err := qb.repository(tx).insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *PerformerQueryBuilder) Update(updatedObject models.PerformerPartial, tx *sqlx.Tx) (*models.Performer, error) {
	const partial = true
	if err := qb.repository(tx).update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	var ret models.Performer
	if err := qb.repository(tx).get(updatedObject.ID, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *PerformerQueryBuilder) UpdateFull(updatedObject models.Performer, tx *sqlx.Tx) (*models.Performer, error) {
	const partial = false
	if err := qb.repository(tx).update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	var ret models.Performer
	if err := qb.repository(tx).get(updatedObject.ID, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *PerformerQueryBuilder) Destroy(id int, tx *sqlx.Tx) error {
	// TODO - add on delete cascade to performers_scenes
	_, err := tx.Exec("DELETE FROM performers_scenes WHERE performer_id = ?", id)
	if err != nil {
		return err
	}

	return qb.repository(tx).destroyExisting([]int{id})
}

func (qb *PerformerQueryBuilder) Find(id int) (*models.Performer, error) {
	var ret models.Performer
	// TODO - this should accept a tx
	if err := qb.repository(nil).get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *PerformerQueryBuilder) FindMany(ids []int) ([]*models.Performer, error) {
	var performers []*models.Performer
	for _, id := range ids {
		performer, err := qb.Find(id)
		if err != nil {
			return nil, err
		}

		if performer == nil {
			return nil, fmt.Errorf("performer with id %d not found", id)
		}

		performers = append(performers, performer)
	}

	return performers, nil
}

func (qb *PerformerQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*models.Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		WHERE scenes_join.scene_id = ?
	`
	args := []interface{}{sceneID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByImageID(imageID int, tx *sqlx.Tx) ([]*models.Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_images as images_join on images_join.performer_id = performers.id
		WHERE images_join.image_id = ?
	`
	args := []interface{}{imageID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByGalleryID(galleryID int, tx *sqlx.Tx) ([]*models.Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_galleries as galleries_join on galleries_join.performer_id = performers.id
		WHERE galleries_join.gallery_id = ?
	`
	args := []interface{}{galleryID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindNameBySceneID(sceneID int, tx *sqlx.Tx) ([]*models.Performer, error) {
	query := `
		SELECT performers.name FROM performers
		LEFT JOIN performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		WHERE scenes_join.scene_id = ?
	`
	args := []interface{}{sceneID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByNames(names []string, tx *sqlx.Tx, nocase bool) ([]*models.Performer, error) {
	query := "SELECT * FROM performers WHERE name"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " IN " + getInBinding(len(names))

	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT performers.id FROM performers"), nil)
}

func (qb *PerformerQueryBuilder) All() ([]*models.Performer, error) {
	return qb.queryPerformers(selectAll("performers")+qb.getPerformerSort(nil), nil, nil)
}

func (qb *PerformerQueryBuilder) AllSlim() ([]*models.Performer, error) {
	return qb.queryPerformers("SELECT performers.id, performers.name, performers.gender FROM performers "+qb.getPerformerSort(nil), nil, nil)
}

func (qb *PerformerQueryBuilder) Query(performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, int) {
	if performerFilter == nil {
		performerFilter = &models.PerformerFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	tableName := "performers"
	query := queryBuilder{
		tableName: tableName,
	}

	query.body = selectDistinctIDs(tableName)
	query.body += `
		left join performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		left join scenes on scenes_join.scene_id = scenes.id
		left join performer_stash_ids on performer_stash_ids.performer_id = performers.id
	`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"performers.name", "performers.checksum", "performers.birthdate", "performers.ethnicity"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		query.addWhere(clause)
		query.addArg(thisArgs...)
	}

	if favoritesFilter := performerFilter.FilterFavorites; favoritesFilter != nil {
		var favStr string
		if *favoritesFilter == true {
			favStr = "1"
		} else {
			favStr = "0"
		}
		query.addWhere("performers.favorite = " + favStr)
	}

	if birthYear := performerFilter.BirthYear; birthYear != nil {
		clauses, thisArgs := getBirthYearFilterClause(birthYear.Modifier, birthYear.Value)
		query.addWhere(clauses...)
		query.addArg(thisArgs...)
	}

	if age := performerFilter.Age; age != nil {
		clauses, thisArgs := getAgeFilterClause(age.Modifier, age.Value)
		query.addWhere(clauses...)
		query.addArg(thisArgs...)
	}

	if gender := performerFilter.Gender; gender != nil {
		query.addWhere("performers.gender = ?")
		query.addArg(gender.Value.String())
	}

	if isMissingFilter := performerFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "scenes":
			query.addWhere("scenes_join.scene_id IS NULL")
		case "image":
			query.body += `left join performers_image on performers_image.performer_id = performers.id
			`
			query.addWhere("performers_image.performer_id IS NULL")
		case "stash_id":
			query.addWhere("performer_stash_ids.performer_id IS NULL")
		default:
			query.addWhere("(performers." + *isMissingFilter + " IS NULL OR TRIM(performers." + *isMissingFilter + ") = '')")
		}
	}

	if stashIDFilter := performerFilter.StashID; stashIDFilter != nil {
		query.addWhere("performer_stash_ids.stash_id = ?")
		query.addArg(stashIDFilter)
	}

	query.handleStringCriterionInput(performerFilter.Ethnicity, tableName+".ethnicity")
	query.handleStringCriterionInput(performerFilter.Country, tableName+".country")
	query.handleStringCriterionInput(performerFilter.EyeColor, tableName+".eye_color")
	query.handleStringCriterionInput(performerFilter.Height, tableName+".height")
	query.handleStringCriterionInput(performerFilter.Measurements, tableName+".measurements")
	query.handleStringCriterionInput(performerFilter.FakeTits, tableName+".fake_tits")
	query.handleStringCriterionInput(performerFilter.CareerLength, tableName+".career_length")
	query.handleStringCriterionInput(performerFilter.Tattoos, tableName+".tattoos")
	query.handleStringCriterionInput(performerFilter.Piercings, tableName+".piercings")

	// TODO - need better handling of aliases
	query.handleStringCriterionInput(performerFilter.Aliases, tableName+".aliases")

	query.sortAndPagination = qb.getPerformerSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := query.executeFind()

	var performers []*models.Performer
	for _, id := range idsResult {
		performer, _ := qb.Find(id)
		performers = append(performers, performer)
	}

	return performers, countResult
}

func getBirthYearFilterClause(criterionModifier models.CriterionModifier, value int) ([]string, []interface{}) {
	var clauses []string
	var args []interface{}

	yearStr := strconv.Itoa(value)
	startOfYear := yearStr + "-01-01"
	endOfYear := yearStr + "-12-31"

	if modifier := criterionModifier.String(); criterionModifier.IsValid() {
		switch modifier {
		case "EQUALS":
			// between yyyy-01-01 and yyyy-12-31
			clauses = append(clauses, "performers.birthdate >= ?")
			clauses = append(clauses, "performers.birthdate <= ?")
			args = append(args, startOfYear)
			args = append(args, endOfYear)
		case "NOT_EQUALS":
			// outside of yyyy-01-01 to yyyy-12-31
			clauses = append(clauses, "performers.birthdate < ? OR performers.birthdate > ?")
			args = append(args, startOfYear)
			args = append(args, endOfYear)
		case "GREATER_THAN":
			// > yyyy-12-31
			clauses = append(clauses, "performers.birthdate > ?")
			args = append(args, endOfYear)
		case "LESS_THAN":
			// < yyyy-01-01
			clauses = append(clauses, "performers.birthdate < ?")
			args = append(args, startOfYear)
		}
	}

	return clauses, args
}

func getAgeFilterClause(criterionModifier models.CriterionModifier, value int) ([]string, []interface{}) {
	var clauses []string
	var args []interface{}

	// get the date at which performer would turn the age specified
	dt := time.Now()
	birthDate := dt.AddDate(-value-1, 0, 0)
	yearAfter := birthDate.AddDate(1, 0, 0)

	if modifier := criterionModifier.String(); criterionModifier.IsValid() {
		switch modifier {
		case "EQUALS":
			// between birthDate and yearAfter
			clauses = append(clauses, "performers.birthdate >= ?")
			clauses = append(clauses, "performers.birthdate < ?")
			args = append(args, birthDate)
			args = append(args, yearAfter)
		case "NOT_EQUALS":
			// outside of birthDate and yearAfter
			clauses = append(clauses, "performers.birthdate < ? OR performers.birthdate >= ?")
			args = append(args, birthDate)
			args = append(args, yearAfter)
		case "GREATER_THAN":
			// < birthDate
			clauses = append(clauses, "performers.birthdate < ?")
			args = append(args, birthDate)
		case "LESS_THAN":
			// > yearAfter
			clauses = append(clauses, "performers.birthdate >= ?")
			args = append(args, yearAfter)
		}
	}

	return clauses, args
}

func (qb *PerformerQueryBuilder) getPerformerSort(findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}
	return getSort(sort, direction, "performers")
}

func (qb *PerformerQueryBuilder) queryPerformers(query string, args []interface{}, tx *sqlx.Tx) ([]*models.Performer, error) {
	var ret models.Performers
	if err := qb.repository(tx).query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Performer(ret), nil
}

func (qb *PerformerQueryBuilder) UpdatePerformerImage(performerID int, image []byte, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing cover and then create new
	if err := qb.DestroyPerformerImage(performerID, tx); err != nil {
		return err
	}

	_, err := tx.Exec(
		`INSERT INTO performers_image (performer_id, image) VALUES (?, ?)`,
		performerID,
		image,
	)

	return err
}

func (qb *PerformerQueryBuilder) DestroyPerformerImage(performerID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM performers_image WHERE performer_id = ?", performerID)
	if err != nil {
		return err
	}
	return err
}

func (qb *PerformerQueryBuilder) GetPerformerImage(performerID int, tx *sqlx.Tx) ([]byte, error) {
	query := `SELECT image from performers_image WHERE performer_id = ?`
	return getImage(tx, query, performerID)
}

func NewPerformerReaderWriter(tx *sqlx.Tx) *performerReaderWriter {
	return &performerReaderWriter{
		tx: tx,
		qb: NewPerformerQueryBuilder(),
	}
}

type performerReaderWriter struct {
	tx *sqlx.Tx
	qb PerformerQueryBuilder
}

func (t *performerReaderWriter) FindMany(ids []int) ([]*models.Performer, error) {
	return t.qb.FindMany(ids)
}

func (t *performerReaderWriter) FindByNames(names []string, nocase bool) ([]*models.Performer, error) {
	return t.qb.FindByNames(names, t.tx, nocase)
}

func (t *performerReaderWriter) All() ([]*models.Performer, error) {
	return t.qb.All()
}

func (t *performerReaderWriter) GetPerformerImage(performerID int) ([]byte, error) {
	return t.qb.GetPerformerImage(performerID, t.tx)
}

func (t *performerReaderWriter) FindBySceneID(id int) ([]*models.Performer, error) {
	return t.qb.FindBySceneID(id, t.tx)
}

func (t *performerReaderWriter) FindNamesBySceneID(sceneID int) ([]*models.Performer, error) {
	return t.qb.FindNameBySceneID(sceneID, t.tx)
}

func (t *performerReaderWriter) FindByImageID(id int) ([]*models.Performer, error) {
	return t.qb.FindByImageID(id, t.tx)
}

func (t *performerReaderWriter) FindByGalleryID(id int) ([]*models.Performer, error) {
	return t.qb.FindByGalleryID(id, t.tx)
}

func (t *performerReaderWriter) Create(newPerformer models.Performer) (*models.Performer, error) {
	return t.qb.Create(newPerformer, t.tx)
}

func (t *performerReaderWriter) Update(updatedPerformer models.PerformerPartial) (*models.Performer, error) {
	return t.qb.Update(updatedPerformer, t.tx)
}

func (t *performerReaderWriter) UpdateFull(updatedPerformer models.Performer) (*models.Performer, error) {
	return t.qb.UpdateFull(updatedPerformer, t.tx)
}

func (t *performerReaderWriter) UpdatePerformerImage(performerID int, image []byte) error {
	return t.qb.UpdatePerformerImage(performerID, image, t.tx)
}
