package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

const studioTable = "studios"
const studioIDColumn = "studio_id"
const studioAliasesTable = "studio_aliases"
const studioAliasColumn = "alias"

type studioQueryBuilder struct {
	repository
}

func NewStudioReaderWriter(tx dbi) *studioQueryBuilder {
	return &studioQueryBuilder{
		repository{
			tx:        tx,
			tableName: studioTable,
			idColumn:  idColumn,
		},
	}
}

func (qb *studioQueryBuilder) Create(newObject models.Studio) (*models.Studio, error) {
	var ret models.Studio
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *studioQueryBuilder) Update(updatedObject models.StudioPartial) (*models.Studio, error) {
	const partial = true
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID)
}

func (qb *studioQueryBuilder) UpdateFull(updatedObject models.Studio) (*models.Studio, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID)
}

func (qb *studioQueryBuilder) Destroy(id int) error {
	// TODO - set null on foreign key in scraped items
	// remove studio from scraped items
	_, err := qb.tx.Exec("UPDATE scraped_items SET studio_id = null WHERE studio_id = ?", id)
	if err != nil {
		return err
	}

	return qb.destroyExisting([]int{id})
}

func (qb *studioQueryBuilder) Find(id int) (*models.Studio, error) {
	var ret models.Studio
	if err := qb.get(id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *studioQueryBuilder) FindMany(ids []int) ([]*models.Studio, error) {
	var studios []*models.Studio
	for _, id := range ids {
		studio, err := qb.Find(id)
		if err != nil {
			return nil, err
		}

		if studio == nil {
			return nil, fmt.Errorf("studio with id %d not found", id)
		}

		studios = append(studios, studio)
	}

	return studios, nil
}

func (qb *studioQueryBuilder) FindChildren(id int) ([]*models.Studio, error) {
	query := "SELECT studios.* FROM studios WHERE studios.parent_id = ?"
	args := []interface{}{id}
	return qb.queryStudios(query, args)
}

func (qb *studioQueryBuilder) FindBySceneID(sceneID int) (*models.Studio, error) {
	query := "SELECT studios.* FROM studios JOIN scenes ON studios.id = scenes.studio_id WHERE scenes.id = ? LIMIT 1"
	args := []interface{}{sceneID}
	return qb.queryStudio(query, args)
}

func (qb *studioQueryBuilder) FindByName(name string, nocase bool) (*models.Studio, error) {
	query := "SELECT * FROM studios WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryStudio(query, args)
}

func (qb *studioQueryBuilder) FindByStashID(stashID models.StashID) ([]*models.Studio, error) {
	query := selectAll("studios") + `
		LEFT JOIN studio_stash_ids on studio_stash_ids.studio_id = studios.id
		WHERE studio_stash_ids.stash_id = ?
		AND studio_stash_ids.endpoint = ?
	`
	args := []interface{}{stashID.StashID, stashID.Endpoint}
	return qb.queryStudios(query, args)
}

func (qb *studioQueryBuilder) Count() (int, error) {
	return qb.runCountQuery(qb.buildCountQuery("SELECT studios.id FROM studios"), nil)
}

func (qb *studioQueryBuilder) All() ([]*models.Studio, error) {
	return qb.queryStudios(selectAll("studios")+qb.getStudioSort(nil), nil)
}

func (qb *studioQueryBuilder) QueryForAutoTag(words []string) ([]*models.Studio, error) {
	// TODO - Query needs to be changed to support queries of this type, and
	// this method should be removed
	query := selectAll(studioTable)
	query += " LEFT JOIN studio_aliases ON studio_aliases.studio_id = studios.id"

	var whereClauses []string
	var args []interface{}

	for _, w := range words {
		ww := w + "%"
		whereClauses = append(whereClauses, "studios.name like ?")
		args = append(args, ww)

		// include aliases
		whereClauses = append(whereClauses, "studio_aliases.alias like ?")
		args = append(args, ww)
	}

	whereOr := "(" + strings.Join(whereClauses, " OR ") + ")"
	where := strings.Join([]string{
		"studios.ignore_auto_tag = 0",
		whereOr,
	}, " AND ")
	return qb.queryStudios(query+" WHERE "+where, args)
}

func (qb *studioQueryBuilder) validateFilter(filter *models.StudioFilterType) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if filter.And != nil {
		if filter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if filter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(filter.And)
	}

	if filter.Or != nil {
		if filter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(filter.Or)
	}

	if filter.Not != nil {
		return qb.validateFilter(filter.Not)
	}

	return nil
}

func (qb *studioQueryBuilder) makeFilter(studioFilter *models.StudioFilterType) *filterBuilder {
	query := &filterBuilder{}

	if studioFilter.And != nil {
		query.and(qb.makeFilter(studioFilter.And))
	}
	if studioFilter.Or != nil {
		query.or(qb.makeFilter(studioFilter.Or))
	}
	if studioFilter.Not != nil {
		query.not(qb.makeFilter(studioFilter.Not))
	}

	query.handleCriterion(stringCriterionHandler(studioFilter.Name, studioTable+".name"))
	query.handleCriterion(stringCriterionHandler(studioFilter.Details, studioTable+".details"))
	query.handleCriterion(stringCriterionHandler(studioFilter.URL, studioTable+".url"))
	query.handleCriterion(intCriterionHandler(studioFilter.Rating, studioTable+".rating"))
	query.handleCriterion(boolCriterionHandler(studioFilter.IgnoreAutoTag, studioTable+".ignore_auto_tag"))

	query.handleCriterion(criterionHandlerFunc(func(f *filterBuilder) {
		if studioFilter.StashID != nil {
			qb.stashIDRepository().join(f, "studio_stash_ids", "studios.id")
			stringCriterionHandler(studioFilter.StashID, "studio_stash_ids.stash_id")(f)
		}
	}))

	query.handleCriterion(studioIsMissingCriterionHandler(qb, studioFilter.IsMissing))
	query.handleCriterion(studioSceneCountCriterionHandler(qb, studioFilter.SceneCount))
	query.handleCriterion(studioImageCountCriterionHandler(qb, studioFilter.ImageCount))
	query.handleCriterion(studioGalleryCountCriterionHandler(qb, studioFilter.GalleryCount))
	query.handleCriterion(studioParentCriterionHandler(qb, studioFilter.Parents))
	query.handleCriterion(studioAliasCriterionHandler(qb, studioFilter.Aliases))

	return query
}

func (qb *studioQueryBuilder) Query(studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, int, error) {
	if studioFilter == nil {
		studioFilter = &models.StudioFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, studioTable)

	if q := findFilter.Q; q != nil && *q != "" {
		query.join(studioAliasesTable, "", "studio_aliases.studio_id = studios.id")
		searchColumns := []string{"studios.name", "studio_aliases.alias"}

		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(studioFilter); err != nil {
		return nil, 0, err
	}
	filter := qb.makeFilter(studioFilter)

	query.addFilter(filter)

	query.sortAndPagination = qb.getStudioSort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind()
	if err != nil {
		return nil, 0, err
	}

	var studios []*models.Studio
	for _, id := range idsResult {
		studio, err := qb.Find(id)
		if err != nil {
			return nil, 0, err
		}

		studios = append(studios, studio)
	}

	return studios, countResult, nil
}

func studioIsMissingCriterionHandler(qb *studioQueryBuilder, isMissing *string) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "image":
				f.addLeftJoin("studios_image", "", "studios_image.studio_id = studios.id")
				f.addWhere("studios_image.studio_id IS NULL")
			case "stash_id":
				qb.stashIDRepository().join(f, "studio_stash_ids", "studios.id")
				f.addWhere("studio_stash_ids.studio_id IS NULL")
			default:
				f.addWhere("(studios." + *isMissing + " IS NULL OR TRIM(studios." + *isMissing + ") = '')")
			}
		}
	}
}

func studioSceneCountCriterionHandler(qb *studioQueryBuilder, sceneCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes", "", "scenes.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes.id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func studioImageCountCriterionHandler(qb *studioQueryBuilder, imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images", "", "images.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct images.id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func studioGalleryCountCriterionHandler(qb *studioQueryBuilder, galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(f *filterBuilder) {
		if galleryCount != nil {
			f.addLeftJoin("galleries", "", "galleries.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct galleries.id)", *galleryCount)

			f.addHaving(clause, args...)
		}
	}
}

func studioParentCriterionHandler(qb *studioQueryBuilder, parents *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		f.addLeftJoin("studios", "parent_studio", "parent_studio.id = studios.parent_id")
	}
	h := multiCriterionHandlerBuilder{
		primaryTable: studioTable,
		foreignTable: "parent_studio",
		joinTable:    "",
		primaryFK:    studioIDColumn,
		foreignFK:    "parent_id",
		addJoinsFunc: addJoinsFunc,
	}
	return h.handler(parents)
}

func studioAliasCriterionHandler(qb *studioQueryBuilder, alias *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    studioAliasesTable,
		stringColumn: studioAliasColumn,
		addJoinTable: func(f *filterBuilder) {
			qb.aliasRepository().join(f, "", "studios.id")
		},
	}

	return h.handler(alias)
}

func (qb *studioQueryBuilder) getStudioSort(findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}

	switch sort {
	case "scenes_count":
		return getCountSort(studioTable, sceneTable, studioIDColumn, direction)
	case "images_count":
		return getCountSort(studioTable, imageTable, studioIDColumn, direction)
	case "galleries_count":
		return getCountSort(studioTable, galleryTable, studioIDColumn, direction)
	default:
		return getSort(sort, direction, "studios")
	}
}

func (qb *studioQueryBuilder) queryStudio(query string, args []interface{}) (*models.Studio, error) {
	results, err := qb.queryStudios(query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *studioQueryBuilder) queryStudios(query string, args []interface{}) ([]*models.Studio, error) {
	var ret models.Studios
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Studio(ret), nil
}

func (qb *studioQueryBuilder) imageRepository() *imageRepository {
	return &imageRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: "studios_image",
			idColumn:  studioIDColumn,
		},
		imageColumn: "image",
	}
}

func (qb *studioQueryBuilder) GetImage(studioID int) ([]byte, error) {
	return qb.imageRepository().get(studioID)
}

func (qb *studioQueryBuilder) HasImage(studioID int) (bool, error) {
	return qb.imageRepository().exists(studioID)
}

func (qb *studioQueryBuilder) UpdateImage(studioID int, image []byte) error {
	return qb.imageRepository().replace(studioID, image)
}

func (qb *studioQueryBuilder) DestroyImage(studioID int) error {
	return qb.imageRepository().destroy([]int{studioID})
}

func (qb *studioQueryBuilder) stashIDRepository() *stashIDRepository {
	return &stashIDRepository{
		repository{
			tx:        qb.tx,
			tableName: "studio_stash_ids",
			idColumn:  studioIDColumn,
		},
	}
}

func (qb *studioQueryBuilder) GetStashIDs(studioID int) ([]*models.StashID, error) {
	return qb.stashIDRepository().get(studioID)
}

func (qb *studioQueryBuilder) UpdateStashIDs(studioID int, stashIDs []models.StashID) error {
	return qb.stashIDRepository().replace(studioID, stashIDs)
}

func (qb *studioQueryBuilder) aliasRepository() *stringRepository {
	return &stringRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: studioAliasesTable,
			idColumn:  studioIDColumn,
		},
		stringColumn: studioAliasColumn,
	}
}

func (qb *studioQueryBuilder) GetAliases(studioID int) ([]string, error) {
	return qb.aliasRepository().get(studioID)
}

func (qb *studioQueryBuilder) UpdateAliases(studioID int, aliases []string) error {
	return qb.aliasRepository().replace(studioID, aliases)
}
