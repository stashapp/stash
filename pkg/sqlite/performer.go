package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

const performerTable = "performers"
const performerIDColumn = "performer_id"
const performersTagsTable = "performers_tags"

var countPerformersForTagQuery = `
SELECT tag_id AS id FROM performers_tags
WHERE performers_tags.tag_id = ?
GROUP BY performers_tags.performer_id
`

type performerQueryBuilder struct {
	repository
}

func NewPerformerReaderWriter(tx dbi) *performerQueryBuilder {
	return &performerQueryBuilder{
		repository{
			tx:        tx,
			tableName: performerTable,
			idColumn:  idColumn,
		},
	}
}

func (qb *performerQueryBuilder) Create(newObject models.Performer) (*models.Performer, error) {
	var ret models.Performer
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *performerQueryBuilder) Update(updatedObject models.PerformerPartial) (*models.Performer, error) {
	const partial = true
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	var ret models.Performer
	if err := qb.get(updatedObject.ID, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *performerQueryBuilder) UpdateFull(updatedObject models.Performer) (*models.Performer, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	var ret models.Performer
	if err := qb.get(updatedObject.ID, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *performerQueryBuilder) Destroy(id int) error {
	// TODO - add on delete cascade to performers_scenes
	_, err := qb.tx.Exec("DELETE FROM performers_scenes WHERE performer_id = ?", id)
	if err != nil {
		return err
	}

	return qb.destroyExisting([]int{id})
}

func (qb *performerQueryBuilder) Find(id int) (*models.Performer, error) {
	var ret models.Performer
	if err := qb.get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *performerQueryBuilder) FindMany(ids []int) ([]*models.Performer, error) {
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

func (qb *performerQueryBuilder) FindBySceneID(sceneID int) ([]*models.Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		WHERE scenes_join.scene_id = ?
	`
	args := []interface{}{sceneID}
	return qb.queryPerformers(query, args)
}

func (qb *performerQueryBuilder) FindByImageID(imageID int) ([]*models.Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_images as images_join on images_join.performer_id = performers.id
		WHERE images_join.image_id = ?
	`
	args := []interface{}{imageID}
	return qb.queryPerformers(query, args)
}

func (qb *performerQueryBuilder) FindByGalleryID(galleryID int) ([]*models.Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_galleries as galleries_join on galleries_join.performer_id = performers.id
		WHERE galleries_join.gallery_id = ?
	`
	args := []interface{}{galleryID}
	return qb.queryPerformers(query, args)
}

func (qb *performerQueryBuilder) FindNamesBySceneID(sceneID int) ([]*models.Performer, error) {
	query := `
		SELECT performers.name FROM performers
		LEFT JOIN performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		WHERE scenes_join.scene_id = ?
	`
	args := []interface{}{sceneID}
	return qb.queryPerformers(query, args)
}

func (qb *performerQueryBuilder) FindByNames(names []string, nocase bool) ([]*models.Performer, error) {
	query := "SELECT * FROM performers WHERE name"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " IN " + getInBinding(len(names))

	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryPerformers(query, args)
}

func (qb *performerQueryBuilder) CountByTagID(tagID int) (int, error) {
	args := []interface{}{tagID}
	return qb.runCountQuery(qb.buildCountQuery(countPerformersForTagQuery), args)
}

func (qb *performerQueryBuilder) Count() (int, error) {
	return qb.runCountQuery(qb.buildCountQuery("SELECT performers.id FROM performers"), nil)
}

func (qb *performerQueryBuilder) All() ([]*models.Performer, error) {
	return qb.queryPerformers(selectAll("performers")+qb.getPerformerSort(nil), nil)
}

func (qb *performerQueryBuilder) AllSlim() ([]*models.Performer, error) {
	return qb.queryPerformers("SELECT performers.id, performers.name, performers.gender FROM performers "+qb.getPerformerSort(nil), nil)
}

func (qb *performerQueryBuilder) Query(performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, int, error) {
	if performerFilter == nil {
		performerFilter = &models.PerformerFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	tableName := "performers"
	query := qb.newQuery()

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

	if tagsFilter := performerFilter.Tags; tagsFilter != nil && len(tagsFilter.Value) > 0 {
		for _, tagID := range tagsFilter.Value {
			query.addArg(tagID)
		}

		query.body += ` left join performers_tags as tags_join on tags_join.performer_id = performers.id
			LEFT JOIN tags on tags_join.tag_id = tags.id`
		whereClause, havingClause := getMultiCriterionClause("performers", "tags", "performers_tags", "performer_id", "tag_id", tagsFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	query.sortAndPagination = qb.getPerformerSort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind()
	if err != nil {
		return nil, 0, err
	}

	var performers []*models.Performer
	for _, id := range idsResult {
		performer, err := qb.Find(id)
		if err != nil {
			return nil, 0, err
		}
		performers = append(performers, performer)
	}

	return performers, countResult, nil
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

func (qb *performerQueryBuilder) getPerformerSort(findFilter *models.FindFilterType) string {
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

func (qb *performerQueryBuilder) queryPerformers(query string, args []interface{}) ([]*models.Performer, error) {
	var ret models.Performers
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Performer(ret), nil
}

func (qb *performerQueryBuilder) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersTagsTable,
			idColumn:  performerIDColumn,
		},
		fkColumn: tagIDColumn,
	}
}

func (qb *performerQueryBuilder) GetTagIDs(id int) ([]int, error) {
	return qb.tagsRepository().getIDs(id)
}

func (qb *performerQueryBuilder) UpdateTags(id int, tagIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.tagsRepository().replace(id, tagIDs)
}

func (qb *performerQueryBuilder) imageRepository() *imageRepository {
	return &imageRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: "performers_image",
			idColumn:  performerIDColumn,
		},
		imageColumn: "image",
	}
}

func (qb *performerQueryBuilder) GetImage(performerID int) ([]byte, error) {
	return qb.imageRepository().get(performerID)
}

func (qb *performerQueryBuilder) UpdateImage(performerID int, image []byte) error {
	return qb.imageRepository().replace(performerID, image)
}

func (qb *performerQueryBuilder) DestroyImage(performerID int) error {
	return qb.imageRepository().destroy([]int{performerID})
}

func (qb *performerQueryBuilder) stashIDRepository() *stashIDRepository {
	return &stashIDRepository{
		repository{
			tx:        qb.tx,
			tableName: "performer_stash_ids",
			idColumn:  performerIDColumn,
		},
	}
}

func (qb *performerQueryBuilder) GetStashIDs(performerID int) ([]*models.StashID, error) {
	return qb.stashIDRepository().get(performerID)
}

func (qb *performerQueryBuilder) UpdateStashIDs(performerID int, stashIDs []models.StashID) error {
	return qb.stashIDRepository().replace(performerID, stashIDs)
}
