package sqlite

import (
	"context"
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

var StudioReaderWriter = &studioQueryBuilder{
	repository{
		tableName: studioTable,
		idColumn:  idColumn,
	},
}

func (qb *studioQueryBuilder) Create(ctx context.Context, newObject models.Studio) (*models.Studio, error) {
	var ret models.Studio
	if err := qb.insertObject(ctx, newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *studioQueryBuilder) Update(ctx context.Context, updatedObject models.StudioPartial) (*models.Studio, error) {
	const partial = true
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(ctx, updatedObject.ID)
}

func (qb *studioQueryBuilder) UpdateFull(ctx context.Context, updatedObject models.Studio) (*models.Studio, error) {
	const partial = false
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(ctx, updatedObject.ID)
}

func (qb *studioQueryBuilder) Destroy(ctx context.Context, id int) error {
	// TODO - set null on foreign key in scraped items
	// remove studio from scraped items
	_, err := qb.tx.Exec(ctx, "UPDATE scraped_items SET studio_id = null WHERE studio_id = ?", id)
	if err != nil {
		return err
	}

	return qb.destroyExisting(ctx, []int{id})
}

func (qb *studioQueryBuilder) Find(ctx context.Context, id int) (*models.Studio, error) {
	var ret models.Studio
	if err := qb.getByID(ctx, id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *studioQueryBuilder) FindMany(ctx context.Context, ids []int) ([]*models.Studio, error) {
	var studios []*models.Studio
	for _, id := range ids {
		studio, err := qb.Find(ctx, id)
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

func (qb *studioQueryBuilder) FindChildren(ctx context.Context, id int) ([]*models.Studio, error) {
	query := "SELECT studios.* FROM studios WHERE studios.parent_id = ?"
	args := []interface{}{id}
	return qb.queryStudios(ctx, query, args)
}

func (qb *studioQueryBuilder) FindBySceneID(ctx context.Context, sceneID int) (*models.Studio, error) {
	query := "SELECT studios.* FROM studios JOIN scenes ON studios.id = scenes.studio_id WHERE scenes.id = ? LIMIT 1"
	args := []interface{}{sceneID}
	return qb.queryStudio(ctx, query, args)
}

func (qb *studioQueryBuilder) FindByName(ctx context.Context, name string, nocase bool) (*models.Studio, error) {
	query := "SELECT * FROM studios WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryStudio(ctx, query, args)
}

func (qb *studioQueryBuilder) FindByStashID(ctx context.Context, stashID models.StashID) ([]*models.Studio, error) {
	query := selectAll("studios") + `
		LEFT JOIN studio_stash_ids on studio_stash_ids.studio_id = studios.id
		WHERE studio_stash_ids.stash_id = ?
		AND studio_stash_ids.endpoint = ?
	`
	args := []interface{}{stashID.StashID, stashID.Endpoint}
	return qb.queryStudios(ctx, query, args)
}

func (qb *studioQueryBuilder) Count(ctx context.Context) (int, error) {
	return qb.runCountQuery(ctx, qb.buildCountQuery("SELECT studios.id FROM studios"), nil)
}

func (qb *studioQueryBuilder) All(ctx context.Context) ([]*models.Studio, error) {
	return qb.queryStudios(ctx, selectAll("studios")+qb.getStudioSort(nil), nil)
}

func (qb *studioQueryBuilder) QueryForAutoTag(ctx context.Context, words []string) ([]*models.Studio, error) {
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
	return qb.queryStudios(ctx, query+" WHERE "+where, args)
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

func (qb *studioQueryBuilder) makeFilter(ctx context.Context, studioFilter *models.StudioFilterType) *filterBuilder {
	query := &filterBuilder{}

	if studioFilter.And != nil {
		query.and(qb.makeFilter(ctx, studioFilter.And))
	}
	if studioFilter.Or != nil {
		query.or(qb.makeFilter(ctx, studioFilter.Or))
	}
	if studioFilter.Not != nil {
		query.not(qb.makeFilter(ctx, studioFilter.Not))
	}

	query.handleCriterion(ctx, stringCriterionHandler(studioFilter.Name, studioTable+".name"))
	query.handleCriterion(ctx, stringCriterionHandler(studioFilter.Details, studioTable+".details"))
	query.handleCriterion(ctx, stringCriterionHandler(studioFilter.URL, studioTable+".url"))
	query.handleCriterion(ctx, intCriterionHandler(studioFilter.Rating, studioTable+".rating"))
	query.handleCriterion(ctx, boolCriterionHandler(studioFilter.IgnoreAutoTag, studioTable+".ignore_auto_tag"))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if studioFilter.StashID != nil {
			qb.stashIDRepository().join(f, "studio_stash_ids", "studios.id")
			stringCriterionHandler(studioFilter.StashID, "studio_stash_ids.stash_id")(ctx, f)
		}
	}))

	query.handleCriterion(ctx, studioIsMissingCriterionHandler(qb, studioFilter.IsMissing))
	query.handleCriterion(ctx, studioSceneCountCriterionHandler(qb, studioFilter.SceneCount))
	query.handleCriterion(ctx, studioImageCountCriterionHandler(qb, studioFilter.ImageCount))
	query.handleCriterion(ctx, studioGalleryCountCriterionHandler(qb, studioFilter.GalleryCount))
	query.handleCriterion(ctx, studioParentCriterionHandler(qb, studioFilter.Parents))
	query.handleCriterion(ctx, studioAliasCriterionHandler(qb, studioFilter.Aliases))

	return query
}

func (qb *studioQueryBuilder) Query(ctx context.Context, studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, int, error) {
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
	filter := qb.makeFilter(ctx, studioFilter)

	query.addFilter(filter)

	query.sortAndPagination = qb.getStudioSort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	var studios []*models.Studio
	for _, id := range idsResult {
		studio, err := qb.Find(ctx, id)
		if err != nil {
			return nil, 0, err
		}

		studios = append(studios, studio)
	}

	return studios, countResult, nil
}

func studioIsMissingCriterionHandler(qb *studioQueryBuilder, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
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
	return func(ctx context.Context, f *filterBuilder) {
		if sceneCount != nil {
			f.addLeftJoin("scenes", "", "scenes.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct scenes.id)", *sceneCount)

			f.addHaving(clause, args...)
		}
	}
}

func studioImageCountCriterionHandler(qb *studioQueryBuilder, imageCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if imageCount != nil {
			f.addLeftJoin("images", "", "images.studio_id = studios.id")
			clause, args := getIntCriterionWhereClause("count(distinct images.id)", *imageCount)

			f.addHaving(clause, args...)
		}
	}
}

func studioGalleryCountCriterionHandler(qb *studioQueryBuilder, galleryCount *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
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

func (qb *studioQueryBuilder) queryStudio(ctx context.Context, query string, args []interface{}) (*models.Studio, error) {
	results, err := qb.queryStudios(ctx, query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *studioQueryBuilder) queryStudios(ctx context.Context, query string, args []interface{}) ([]*models.Studio, error) {
	var ret models.Studios
	if err := qb.query(ctx, query, args, &ret); err != nil {
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

func (qb *studioQueryBuilder) GetImage(ctx context.Context, studioID int) ([]byte, error) {
	return qb.imageRepository().get(ctx, studioID)
}

func (qb *studioQueryBuilder) HasImage(ctx context.Context, studioID int) (bool, error) {
	return qb.imageRepository().exists(ctx, studioID)
}

func (qb *studioQueryBuilder) UpdateImage(ctx context.Context, studioID int, image []byte) error {
	return qb.imageRepository().replace(ctx, studioID, image)
}

func (qb *studioQueryBuilder) DestroyImage(ctx context.Context, studioID int) error {
	return qb.imageRepository().destroy(ctx, []int{studioID})
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

func (qb *studioQueryBuilder) GetStashIDs(ctx context.Context, studioID int) ([]*models.StashID, error) {
	return qb.stashIDRepository().get(ctx, studioID)
}

func (qb *studioQueryBuilder) UpdateStashIDs(ctx context.Context, studioID int, stashIDs []models.StashID) error {
	return qb.stashIDRepository().replace(ctx, studioID, stashIDs)
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

func (qb *studioQueryBuilder) GetAliases(ctx context.Context, studioID int) ([]string, error) {
	return qb.aliasRepository().get(ctx, studioID)
}

func (qb *studioQueryBuilder) UpdateAliases(ctx context.Context, studioID int, aliases []string) error {
	return qb.aliasRepository().replace(ctx, studioID, aliases)
}
