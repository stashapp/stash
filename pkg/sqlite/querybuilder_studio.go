package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const studioTable = "studios"

type StudioQueryBuilder struct{}

func NewStudioQueryBuilder() StudioQueryBuilder {
	return StudioQueryBuilder{}
}

func studioConstructor() interface{} {
	return &models.Studio{}
}

func (qb *StudioQueryBuilder) repository(tx *sqlx.Tx) *repository {
	return &repository{
		tx:          tx,
		tableName:   studioTable,
		idColumn:    idColumn,
		constructor: studioConstructor,
	}
}

func (qb *StudioQueryBuilder) Create(newObject models.Studio, tx *sqlx.Tx) (*models.Studio, error) {
	var ret models.Studio
	if err := qb.repository(tx).insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *StudioQueryBuilder) Update(updatedObject models.StudioPartial, tx *sqlx.Tx) (*models.Studio, error) {
	const partial = true
	if err := qb.repository(tx).update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID, tx)
}

func (qb *StudioQueryBuilder) UpdateFull(updatedObject models.Studio, tx *sqlx.Tx) (*models.Studio, error) {
	const partial = false
	if err := qb.repository(tx).update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID, tx)
}

func (qb *StudioQueryBuilder) Destroy(id int, tx *sqlx.Tx) error {
	// TODO - set null on foreign key in scraped items
	// remove studio from scraped items
	_, err := tx.Exec("UPDATE scraped_items SET studio_id = null WHERE studio_id = ?", id)
	if err != nil {
		return err
	}

	return qb.repository(tx).destroyExisting([]int{id})
}

func (qb *StudioQueryBuilder) Find(id int, tx *sqlx.Tx) (*models.Studio, error) {
	var ret models.Studio
	if err := qb.repository(tx).get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *StudioQueryBuilder) FindMany(ids []int) ([]*models.Studio, error) {
	var studios []*models.Studio
	for _, id := range ids {
		studio, err := qb.Find(id, nil)
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

func (qb *StudioQueryBuilder) FindChildren(id int, tx *sqlx.Tx) ([]*models.Studio, error) {
	query := "SELECT studios.* FROM studios WHERE studios.parent_id = ?"
	args := []interface{}{id}
	return qb.queryStudios(query, args, tx)
}

func (qb *StudioQueryBuilder) FindBySceneID(sceneID int) (*models.Studio, error) {
	query := "SELECT studios.* FROM studios JOIN scenes ON studios.id = scenes.studio_id WHERE scenes.id = ? LIMIT 1"
	args := []interface{}{sceneID}
	return qb.queryStudio(query, args, nil)
}

func (qb *StudioQueryBuilder) FindByName(name string, tx *sqlx.Tx, nocase bool) (*models.Studio, error) {
	query := "SELECT * FROM studios WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryStudio(query, args, tx)
}

func (qb *StudioQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT studios.id FROM studios"), nil)
}

func (qb *StudioQueryBuilder) All() ([]*models.Studio, error) {
	return qb.queryStudios(selectAll("studios")+qb.getStudioSort(nil), nil, nil)
}

func (qb *StudioQueryBuilder) AllSlim() ([]*models.Studio, error) {
	return qb.queryStudios("SELECT studios.id, studios.name, studios.parent_id FROM studios "+qb.getStudioSort(nil), nil, nil)
}

func (qb *StudioQueryBuilder) Query(studioFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, int) {
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
	idsResult, countResult := executeFindQuery("studios", body, args, sortAndPagination, whereClauses, havingClauses)

	var studios []*models.Studio
	for _, id := range idsResult {
		studio, _ := qb.Find(id, nil)
		studios = append(studios, studio)
	}

	return studios, countResult
}

func (qb *StudioQueryBuilder) getStudioSort(findFilter *models.FindFilterType) string {
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

func (qb *StudioQueryBuilder) queryStudio(query string, args []interface{}, tx *sqlx.Tx) (*models.Studio, error) {
	results, err := qb.queryStudios(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *StudioQueryBuilder) queryStudios(query string, args []interface{}, tx *sqlx.Tx) ([]*models.Studio, error) {
	var ret models.Studios
	if err := qb.repository(tx).query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Studio(ret), nil
}

func (qb *StudioQueryBuilder) UpdateStudioImage(studioID int, image []byte, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing cover and then create new
	if err := qb.DestroyStudioImage(studioID, tx); err != nil {
		return err
	}

	_, err := tx.Exec(
		`INSERT INTO studios_image (studio_id, image) VALUES (?, ?)`,
		studioID,
		image,
	)

	return err
}

func (qb *StudioQueryBuilder) DestroyStudioImage(studioID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM studios_image WHERE studio_id = ?", studioID)
	if err != nil {
		return err
	}
	return err
}

func (qb *StudioQueryBuilder) GetStudioImage(studioID int, tx *sqlx.Tx) ([]byte, error) {
	query := `SELECT image from studios_image WHERE studio_id = ?`
	return getImage(tx, query, studioID)
}

func (qb *StudioQueryBuilder) HasStudioImage(studioID int) (bool, error) {
	ret, err := runCountQuery(buildCountQuery("SELECT studio_id from studios_image WHERE studio_id = ?"), []interface{}{studioID})
	if err != nil {
		return false, err
	}

	return ret == 1, nil
}

func NewStudioReaderWriter(tx *sqlx.Tx) *studioReaderWriter {
	return &studioReaderWriter{
		tx: tx,
		qb: NewStudioQueryBuilder(),
	}
}

type studioReaderWriter struct {
	tx *sqlx.Tx
	qb StudioQueryBuilder
}

func (t *studioReaderWriter) Find(id int) (*models.Studio, error) {
	return t.qb.Find(id, t.tx)
}

func (t *studioReaderWriter) FindMany(ids []int) ([]*models.Studio, error) {
	return t.qb.FindMany(ids)
}

func (t *studioReaderWriter) FindByName(name string, nocase bool) (*models.Studio, error) {
	return t.qb.FindByName(name, t.tx, nocase)
}

func (t *studioReaderWriter) All() ([]*models.Studio, error) {
	return t.qb.All()
}

func (t *studioReaderWriter) GetStudioImage(studioID int) ([]byte, error) {
	return t.qb.GetStudioImage(studioID, t.tx)
}

func (t *studioReaderWriter) Create(newStudio models.Studio) (*models.Studio, error) {
	return t.qb.Create(newStudio, t.tx)
}

func (t *studioReaderWriter) Update(updatedStudio models.StudioPartial) (*models.Studio, error) {
	return t.qb.Update(updatedStudio, t.tx)
}

func (t *studioReaderWriter) UpdateFull(updatedStudio models.Studio) (*models.Studio, error) {
	return t.qb.UpdateFull(updatedStudio, t.tx)
}

func (t *studioReaderWriter) UpdateStudioImage(studioID int, image []byte) error {
	return t.qb.UpdateStudioImage(studioID, image, t.tx)
}
