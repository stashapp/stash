package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/database"
	"path/filepath"
)

type galleryQueryBuilder struct {}

func NewGalleryQueryBuilder() galleryQueryBuilder {
	return galleryQueryBuilder{}
}

func (qb *galleryQueryBuilder) Create(newGallery Gallery, tx *sqlx.Tx) (*Gallery, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO galleries (path, checksum, scene_id, created_at, updated_at)
				VALUES (:path, :checksum, :scene_id, :created_at, :updated_at)
		`,
		newGallery,
	)
	if err != nil {
		return nil, err
	}
	galleryID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err := tx.Get(&newGallery, `SELECT * FROM galleries WHERE id = ? LIMIT 1`, galleryID); err != nil {
		return nil, err
	}
	return &newGallery, nil
}

func (qb *galleryQueryBuilder) Update(updatedGallery Gallery, tx *sqlx.Tx) (*Gallery, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE galleries SET `+SqlGenKeys(updatedGallery)+` WHERE galleries.id = :id`,
		updatedGallery,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedGallery, `SELECT * FROM galleries WHERE id = ? LIMIT 1`, updatedGallery.ID); err != nil {
		return nil, err
	}
	return &updatedGallery, nil
}

func (qb *galleryQueryBuilder) Find(id int) (*Gallery, error) {
	query := "SELECT * FROM galleries WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryGallery(query, args, nil)
}

func (qb *galleryQueryBuilder) FindByChecksum(checksum string, tx *sqlx.Tx) (*Gallery, error) {
	query := "SELECT * FROM galleries WHERE checksum = ? LIMIT 1"
	args := []interface{}{checksum}
	return qb.queryGallery(query, args, tx)
}

func (qb *galleryQueryBuilder) FindByPath(path string) (*Gallery, error) {
	query := "SELECT * FROM galleries WHERE path = ? LIMIT 1"
	args := []interface{}{path}
	return qb.queryGallery(query, args, nil)
}

func (qb *galleryQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) (*Gallery, error) {
	query := "SELECT galleries.* FROM galleries JOIN scenes ON scenes.id = galleries.scene_id WHERE scenes.id = ? LIMIT 1"
	args := []interface{}{sceneID}
	return qb.queryGallery(query, args, tx)
}

func (qb *galleryQueryBuilder) ValidGalleriesForScenePath(scenePath string) ([]Gallery, error) {
	sceneDirPath := filepath.Dir(scenePath)
	query := "SELECT galleries.* FROM galleries WHERE galleries.scene_id IS NULL AND galleries.path LIKE '" + sceneDirPath + "%' ORDER BY path ASC"
	return qb.queryGalleries(query, nil, nil)
}

func (qb *galleryQueryBuilder) All() ([]Gallery, error) {
	return qb.queryGalleries(selectAll("galleries") + qb.getGallerySort(nil), nil, nil)
}

func (qb *galleryQueryBuilder) Query(findFilter *FindFilterType) ([]Gallery, int) {
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	whereClauses := []string{}
	havingClauses := []string{}
	args := []interface{}{}
	body := selectDistinctIDs("galleries")

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"galleries.path", "galleries.checksum"}
		whereClauses = append(whereClauses, getSearch(searchColumns, *q))
	}

	sortAndPagination := qb.getGallerySort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("galleries", body, args, sortAndPagination, whereClauses, havingClauses)

	var galleries []Gallery
	for _, id := range idsResult {
		gallery, _ := qb.Find(id)
		galleries = append(galleries, *gallery)
	}

	return galleries, countResult
}

func (qb *galleryQueryBuilder) getGallerySort(findFilter *FindFilterType) string {
	var sort string
	var direction string
	//if findFilter == nil { // TODO temp until title is removed from schema and UI
		sort = "path"
		direction = "ASC"
	//} else {
	//	sort = findFilter.getSort("path")
	//	direction = findFilter.getDirection()
	//}
	return getSort(sort, direction, "galleries")
}

func (qb *galleryQueryBuilder) queryGallery(query string, args []interface{}, tx *sqlx.Tx) (*Gallery, error) {
	results, err := qb.queryGalleries(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return &results[0], nil
}

func (qb *galleryQueryBuilder) queryGalleries(query string, args []interface{}, tx *sqlx.Tx) ([]Gallery, error) {
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

	galleries := make([]Gallery, 0)
	gallery := Gallery{}
	for rows.Next() {
		if err := rows.StructScan(&gallery); err != nil {
			return nil, err
		}
		galleries = append(galleries, gallery)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return galleries, nil
}