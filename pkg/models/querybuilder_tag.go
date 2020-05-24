package models

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

type TagQueryBuilder struct{}

func NewTagQueryBuilder() TagQueryBuilder {
	return TagQueryBuilder{}
}

func (qb *TagQueryBuilder) Create(newTag Tag, tx *sqlx.Tx) (*Tag, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO tags (name, created_at, updated_at)
				VALUES (:name, :created_at, :updated_at)
		`,
		newTag,
	)
	if err != nil {
		return nil, err
	}
	studioID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&newTag, `SELECT * FROM tags WHERE id = ? LIMIT 1`, studioID); err != nil {
		return nil, err
	}

	return &newTag, nil
}

func (qb *TagQueryBuilder) Update(updatedTag Tag, tx *sqlx.Tx) (*Tag, error) {
	ensureTx(tx)
	query := `UPDATE tags SET ` + SQLGenKeys(updatedTag) + ` WHERE tags.id = :id`
	_, err := tx.NamedExec(
		query,
		updatedTag,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedTag, `SELECT * FROM tags WHERE id = ? LIMIT 1`, updatedTag.ID); err != nil {
		return nil, err
	}
	return &updatedTag, nil
}

func (qb *TagQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	// delete tag from scenes and markers first
	_, err := tx.Exec("DELETE FROM scenes_tags WHERE tag_id = ?", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM scene_markers_tags WHERE tag_id = ?", id)
	if err != nil {
		return err
	}

	// cannot unset primary_tag_id in scene_markers because it is not nullable
	countQuery := "SELECT COUNT(*) as count FROM scene_markers where primary_tag_id = ?"
	args := []interface{}{id}
	primaryMarkers, err := runCountQuery(countQuery, args)
	if err != nil {
		return err
	}

	if primaryMarkers > 0 {
		return errors.New("Cannot delete tag used as a primary tag in scene markers")
	}

	return executeDeleteQuery("tags", id, tx)
}

func (qb *TagQueryBuilder) Find(id int, tx *sqlx.Tx) (*Tag, error) {
	query := "SELECT * FROM tags WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryTag(query, args, tx)
}

func (qb *TagQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN scenes_tags as scenes_join on scenes_join.tag_id = tags.id
		WHERE scenes_join.scene_id = ?
		GROUP BY tags.id
	`
	query += qb.getTagSort(nil)
	args := []interface{}{sceneID}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) FindBySceneMarkerID(sceneMarkerID int, tx *sqlx.Tx) ([]*Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN scene_markers_tags as scene_markers_join on scene_markers_join.tag_id = tags.id
		WHERE scene_markers_join.scene_marker_id = ?
		GROUP BY tags.id
	`
	query += qb.getTagSort(nil)
	args := []interface{}{sceneMarkerID}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) FindByName(name string, tx *sqlx.Tx, nocase bool) (*Tag, error) {
	query := "SELECT * FROM tags WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryTag(query, args, tx)
}

func (qb *TagQueryBuilder) FindByNames(names []string, tx *sqlx.Tx, nocase bool) ([]*Tag, error) {
	query := "SELECT * FROM tags WHERE name"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT tags.id FROM tags"), nil)
}

func (qb *TagQueryBuilder) All() ([]*Tag, error) {
	return qb.queryTags(selectAll("tags")+qb.getTagSort(nil), nil, nil)
}

func (qb *TagQueryBuilder) AllSlim() ([]*Tag, error) {
	return qb.queryTags("SELECT tags.id, tags.name FROM tags "+qb.getTagSort(nil), nil, nil)
}

func (qb *TagQueryBuilder) Query(findFilter *FindFilterType) ([]*Tag, int) {
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	var whereClauses []string
	var havingClauses []string
	var args []interface{}
	body := selectDistinctIDs("tags")

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"tags.name"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		whereClauses = append(whereClauses, clause)
		args = append(args, thisArgs...)
	}

	sortAndPagination := qb.getTagSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := executeFindQuery("tags", body, args, sortAndPagination, whereClauses, havingClauses)

	var tags []*Tag
	for _, id := range idsResult {
		tag, _ := qb.Find(id, nil)
		tags = append(tags, tag)
	}

	return tags, countResult
}

func (qb *TagQueryBuilder) getTagSort(findFilter *FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}
	return getSort(sort, direction, "tags")
}

func (qb *TagQueryBuilder) queryTag(query string, args []interface{}, tx *sqlx.Tx) (*Tag, error) {
	results, err := qb.queryTags(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *TagQueryBuilder) queryTags(query string, args []interface{}, tx *sqlx.Tx) ([]*Tag, error) {
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

	tags := make([]*Tag, 0)
	for rows.Next() {
		tag := Tag{}
		if err := rows.StructScan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
