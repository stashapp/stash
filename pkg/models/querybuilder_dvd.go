package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

type DvdQueryBuilder struct{}

func NewDvdQueryBuilder() DvdQueryBuilder {
	return DvdQueryBuilder{}
}

func (qb *DvdQueryBuilder) Create(newDvd Dvd, tx *sqlx.Tx) (*Dvd, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO dvds (frontimage, backimage, checksum, name, aliases, durationdvd, year, director, synopsis, url, created_at, updated_at)
				VALUES (:frontimage, :backimage, :checksum, :name, :aliases, :durationdvd, :year, :director, :synopsis, :url, :created_at, :updated_at)
		`,
		newDvd,
	)
	if err != nil {
		return nil, err
	}
	dvdID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&newDvd, `SELECT * FROM dvds WHERE id = ? LIMIT 1`, dvdID); err != nil {
		return nil, err
	}
	return &newDvd, nil
}

func (qb *DvdQueryBuilder) Update(updatedDvd Dvd, tx *sqlx.Tx) (*Dvd, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE dvds SET `+SQLGenKeys(updatedDvd)+` WHERE dvds.id = :id`,
		updatedDvd,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedDvd, `SELECT * FROM dvds WHERE id = ? LIMIT 1`, updatedDvd.ID); err != nil {
		return nil, err
	}
	return &updatedDvd, nil
}

func (qb *DvdQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	// remove dvd from scenes
	_, err := tx.Exec("UPDATE scenes SET dvd_id = null WHERE dvd_id = ?", id)
	if err != nil {
		return err
	}

	// remove dvd from scraped items
	_, err = tx.Exec("UPDATE scraped_items SET dvd_id = null WHERE dvd_id = ?", id)
	if err != nil {
		return err
	}

	return executeDeleteQuery("dvds", id, tx)
}

func (qb *DvdQueryBuilder) Find(id int, tx *sqlx.Tx) (*Dvd, error) {
	query := "SELECT * FROM dvds WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryDvd(query, args, tx)
}

func (qb *DvdQueryBuilder) FindBySceneID(sceneID int) (*Dvd, error) {
	query := "SELECT dvds.* FROM dvds JOIN scenes ON dvds.id = scenes.dvd_id WHERE scenes.id = ? LIMIT 1"
	args := []interface{}{sceneID}
	return qb.queryDvd(query, args, nil)
}

func (qb *DvdQueryBuilder) FindByName(name string, tx *sqlx.Tx) (*Dvd, error) {
	query := "SELECT * FROM dvds WHERE name = ? LIMIT 1"
	args := []interface{}{name}
	return qb.queryDvd(query, args, tx)
}

func (qb *DvdQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT dvds.id FROM dvds"), nil)
}

func (qb *DvdQueryBuilder) All() ([]*Dvd, error) {
	return qb.queryDvds(selectAll("dvds")+qb.getDvdSort(nil), nil, nil)
}

func (qb *DvdQueryBuilder) Query(findFilter *FindFilterType) ([]*Dvd, int) {
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	var whereClauses []string
	var havingClauses []string
	var args []interface{}
	body := selectDistinctIDs("dvds")
	body += `
		left join scenes on dvds.id = scenes.dvd_id		
	`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"dvds.name"}
		whereClauses = append(whereClauses, getSearch(searchColumns, *q))
	}

	sortAndPagination := qb.getDvdSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("dvds", body, args, sortAndPagination, whereClauses, havingClauses)

	var dvds []*Dvd
	for _, id := range idsResult {
		dvd, _ := qb.Find(id, nil)
		dvds = append(dvds, dvd)
	}

	return dvds, countResult
}

func (qb *DvdQueryBuilder) getDvdSort(findFilter *FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}
	return getSort(sort, direction, "dvds")
}

func (qb *DvdQueryBuilder) queryDvd(query string, args []interface{}, tx *sqlx.Tx) (*Dvd, error) {
	results, err := qb.queryDvds(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *DvdQueryBuilder) queryDvds(query string, args []interface{}, tx *sqlx.Tx) ([]*Dvd, error) {
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

	dvds := make([]*Dvd, 0)
	for rows.Next() {
		dvd := Dvd{}
		if err := rows.StructScan(&dvd); err != nil {
			return nil, err
		}
		dvds = append(dvds, &dvd)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dvds, nil
}
