package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const studioTable = "studios"
const studioIDColumn = "studio_id"

type studioQueryBuilder struct {
	repository
}

func NewStudioReaderWriter(tx *sqlx.Tx) *studioQueryBuilder {
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
		if err == sql.ErrNoRows {
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

func (qb *studioQueryBuilder) Count() (int, error) {
	return qb.runCountQuery(qb.buildCountQuery("SELECT studios.id FROM studios"), nil)
}

func (qb *studioQueryBuilder) All() ([]*models.Studio, error) {
	return qb.queryStudios(selectAll("studios")+qb.getStudioSort(nil), nil)
}

func (qb *studioQueryBuilder) AllSlim() ([]*models.Studio, error) {
	return qb.queryStudios("SELECT studios.id, studios.name, studios.parent_id FROM studios "+qb.getStudioSort(nil), nil)
}

func (qb *studioQueryBuilder) Query(studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, int, error) {
	if studioFilter == nil {
		studioFilter = &models.StudioFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	var whereClauses []string
	var havingClauses []string
	var args []interface{}
	body := selectDistinctIDs("studios")
	body += `
		left join scenes on studios.id = scenes.studio_id		
		left join studio_stash_ids on studio_stash_ids.studio_id = studios.id
	`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"studios.name"}

		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		whereClauses = append(whereClauses, clause)
		args = append(args, thisArgs...)
	}

	if parentsFilter := studioFilter.Parents; parentsFilter != nil && len(parentsFilter.Value) > 0 {
		body += `
			left join studios as parent_studio on parent_studio.id = studios.parent_id
		`

		for _, studioID := range parentsFilter.Value {
			args = append(args, studioID)
		}

		whereClause, havingClause := getMultiCriterionClause("studios", "parent_studio", "", "", "parent_id", parentsFilter)
		whereClauses = appendClause(whereClauses, whereClause)
		havingClauses = appendClause(havingClauses, havingClause)
	}

	if stashIDFilter := studioFilter.StashID; stashIDFilter != nil {
		whereClauses = append(whereClauses, "studio_stash_ids.stash_id = ?")
		args = append(args, stashIDFilter)
	}

	if isMissingFilter := studioFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "image":
			body += `left join studios_image on studios_image.studio_id = studios.id
			`
			whereClauses = appendClause(whereClauses, "studios_image.studio_id IS NULL")
		case "stash_id":
			whereClauses = appendClause(whereClauses, "studio_stash_ids.studio_id IS NULL")
		default:
			whereClauses = appendClause(whereClauses, "studios."+*isMissingFilter+" IS NULL")
		}
	}

	sortAndPagination := qb.getStudioSort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := qb.executeFindQuery(body, args, sortAndPagination, whereClauses, havingClauses)
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
	return getSort(sort, direction, "studios")
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
