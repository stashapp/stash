package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const idColumn = "id"

type objectList interface {
	Append(o interface{})
	New() interface{}
}

type repository struct {
	tx        dbi
	tableName string
	idColumn  string
}

func (r *repository) get(id int, dest interface{}) error {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1", r.tableName, r.idColumn)
	return r.tx.Get(dest, stmt, id)
}

func (r *repository) getAll(id int, f func(rows *sqlx.Rows) error) error {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", r.tableName, r.idColumn)
	return r.queryFunc(stmt, []interface{}{id}, false, f)
}

func (r *repository) insert(obj interface{}) (sql.Result, error) {
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", r.tableName, listKeys(obj, false), listKeys(obj, true))
	return r.tx.NamedExec(stmt, obj)
}

func (r *repository) insertObject(obj interface{}, out interface{}) error {
	result, err := r.insert(obj)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	return r.get(int(id), out)
}

func (r *repository) update(id int, obj interface{}, partial bool) error {
	exists, err := r.exists(id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("%s %d does not exist in %s", r.idColumn, id, r.tableName)
	}

	stmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s.%s = :id", r.tableName, updateSet(obj, partial), r.tableName, r.idColumn)
	_, err = r.tx.NamedExec(stmt, obj)

	return err
}

func (r *repository) updateMap(id int, m map[string]interface{}) error {
	exists, err := r.exists(id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("%s %d does not exist in %s", r.idColumn, id, r.tableName)
	}

	stmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s.%s = :id", r.tableName, updateSetMap(m), r.tableName, r.idColumn)
	_, err = r.tx.NamedExec(stmt, m)

	return err
}

func (r *repository) destroyExisting(ids []int) error {
	for _, id := range ids {
		exists, err := r.exists(id)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("%s %d does not exist in %s", r.idColumn, id, r.tableName)
		}
	}

	return r.destroy(ids)
}

func (r *repository) destroy(ids []int) error {
	for _, id := range ids {
		stmt := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", r.tableName, r.idColumn)
		if _, err := r.tx.Exec(stmt, id); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) exists(id int) (bool, error) {
	stmt := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? LIMIT 1", r.idColumn, r.tableName, r.idColumn)
	stmt = r.buildCountQuery(stmt)

	c, err := r.runCountQuery(stmt, []interface{}{id})
	if err != nil {
		return false, err
	}

	return c == 1, nil
}

func (r *repository) buildCountQuery(query string) string {
	return "SELECT COUNT(*) as count FROM (" + query + ") as temp"
}

func (r *repository) runCountQuery(query string, args []interface{}) (int, error) {
	result := struct {
		Int int `db:"count"`
	}{0}

	// Perform query and fetch result
	if err := r.tx.Get(&result, query, args...); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	return result.Int, nil
}

func (r *repository) runIdsQuery(query string, args []interface{}) ([]int, error) {
	var result []struct {
		Int int `db:"id"`
	}

	if err := r.tx.Select(&result, query, args...); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return []int{}, err
	}

	vsm := make([]int, len(result))
	for i, v := range result {
		vsm[i] = v.Int
	}
	return vsm, nil
}

func (r *repository) runSumQuery(query string, args []interface{}) (float64, error) {
	// Perform query and fetch result
	result := struct {
		Float64 float64 `db:"sum"`
	}{0}

	// Perform query and fetch result
	if err := r.tx.Get(&result, query, args...); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	return result.Float64, nil
}

func (r *repository) queryFunc(query string, args []interface{}, single bool, f func(rows *sqlx.Rows) error) error {
	logger.Tracef("SQL: %s, args: %v", query, args)

	rows, err := r.tx.Queryx(query, args...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := f(rows); err != nil {
			return err
		}
		if single {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (r *repository) query(query string, args []interface{}, out objectList) error {
	return r.queryFunc(query, args, false, func(rows *sqlx.Rows) error {
		object := out.New()
		if err := rows.StructScan(object); err != nil {
			return err
		}
		out.Append(object)
		return nil
	})
}

func (r *repository) queryStruct(query string, args []interface{}, out interface{}) error {
	return r.queryFunc(query, args, true, func(rows *sqlx.Rows) error {
		if err := rows.StructScan(out); err != nil {
			return err
		}
		return nil
	})
}

func (r *repository) querySimple(query string, args []interface{}, out interface{}) error {
	rows, err := r.tx.Queryx(query, args...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(out); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (r *repository) buildQueryBody(body string, whereClauses []string, havingClauses []string) string {
	if len(whereClauses) > 0 {
		body = body + " WHERE " + strings.Join(whereClauses, " AND ") // TODO handle AND or OR
	}
	if len(havingClauses) > 0 {
		body = body + " GROUP BY " + r.tableName + ".id "
		body = body + " HAVING " + strings.Join(havingClauses, " AND ") // TODO handle AND or OR
	}

	return body
}

func (r *repository) executeFindQuery(body string, args []interface{}, sortAndPagination string, whereClauses []string, havingClauses []string, withClauses []string, recursiveWith bool) ([]int, int, error) {
	body = r.buildQueryBody(body, whereClauses, havingClauses)

	withClause := ""
	if len(withClauses) > 0 {
		var recursive string
		if recursiveWith {
			recursive = " RECURSIVE "
		}
		withClause = "WITH " + recursive + strings.Join(withClauses, ", ") + " "
	}

	countQuery := withClause + r.buildCountQuery(body)
	idsQuery := withClause + body + sortAndPagination

	// Perform query and fetch result
	logger.Tracef("SQL: %s, args: %v", idsQuery, args)

	var countResult int
	var countErr error
	var idsResult []int
	var idsErr error

	countResult, countErr = r.runCountQuery(countQuery, args)
	idsResult, idsErr = r.runIdsQuery(idsQuery, args)

	if countErr != nil {
		return nil, 0, fmt.Errorf("error executing count query with SQL: %s, args: %v, error: %s", countQuery, args, countErr.Error())
	}
	if idsErr != nil {
		return nil, 0, fmt.Errorf("error executing find query with SQL: %s, args: %v, error: %s", idsQuery, args, idsErr.Error())
	}

	return idsResult, countResult, nil
}

func (r *repository) newQuery() queryBuilder {
	return queryBuilder{
		repository: r,
	}
}

func (r *repository) join(j joiner, as string, parentIDCol string) {
	t := r.tableName
	if as != "" {
		t = as
	}
	j.addLeftJoin(r.tableName, as, fmt.Sprintf("%s.%s = %s", t, r.idColumn, parentIDCol))
}

//nolint:golint,unused
func (r *repository) innerJoin(j joiner, as string, parentIDCol string) {
	t := r.tableName
	if as != "" {
		t = as
	}
	j.addInnerJoin(r.tableName, as, fmt.Sprintf("%s.%s = %s", t, r.idColumn, parentIDCol))
}

type joiner interface {
	addLeftJoin(table, as, onClause string)
	addInnerJoin(table, as, onClause string)
}

type joinRepository struct {
	repository
	fkColumn string
}

func (r *joinRepository) getIDs(id int) ([]int, error) {
	query := fmt.Sprintf(`SELECT %s as id from %s WHERE %s = ?`, r.fkColumn, r.tableName, r.idColumn)
	return r.runIdsQuery(query, []interface{}{id})
}

func (r *joinRepository) insert(id, foreignID int) (sql.Result, error) {
	stmt := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)", r.tableName, r.idColumn, r.fkColumn)
	return r.tx.Exec(stmt, id, foreignID)
}

func (r *joinRepository) replace(id int, foreignIDs []int) error {
	if err := r.destroy([]int{id}); err != nil {
		return err
	}

	for _, fk := range foreignIDs {
		if _, err := r.insert(id, fk); err != nil {
			return err
		}
	}

	return nil
}

type imageRepository struct {
	repository
	imageColumn string
}

func (r *imageRepository) get(id int) ([]byte, error) {
	query := fmt.Sprintf("SELECT %s from %s WHERE %s = ?", r.imageColumn, r.tableName, r.idColumn)
	var ret []byte
	err := r.querySimple(query, []interface{}{id}, &ret)
	return ret, err
}

func (r *imageRepository) replace(id int, image []byte) error {
	if err := r.destroy([]int{id}); err != nil {
		return err
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)", r.tableName, r.idColumn, r.imageColumn)
	_, err := r.tx.Exec(stmt, id, image)

	return err
}

type captionRepository struct {
	repository
}

func (r *captionRepository) get(id int) ([]*models.SceneCaption, error) {
	query := fmt.Sprintf("SELECT %s, %s, %s from %s WHERE %s = ?", sceneCaptionCodeColumn, sceneCaptionFilenameColumn, sceneCaptionTypeColumn, r.tableName, r.idColumn)
	var ret []*models.SceneCaption
	err := r.queryFunc(query, []interface{}{id}, false, func(rows *sqlx.Rows) error {
		var captionCode string
		var captionFilename string
		var captionType string

		if err := rows.Scan(&captionCode, &captionFilename, &captionType); err != nil {
			return err
		}

		caption := &models.SceneCaption{
			LanguageCode: captionCode,
			Filename:     captionFilename,
			CaptionType:  captionType,
		}
		ret = append(ret, caption)
		return nil
	})
	return ret, err
}

func (r *captionRepository) insert(id int, caption *models.SceneCaption) (sql.Result, error) {
	stmt := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s) VALUES (?, ?, ?, ?)", r.tableName, r.idColumn, sceneCaptionCodeColumn, sceneCaptionFilenameColumn, sceneCaptionTypeColumn)
	return r.tx.Exec(stmt, id, caption.LanguageCode, caption.Filename, caption.CaptionType)
}

func (r *captionRepository) replace(id int, captions []*models.SceneCaption) error {
	if err := r.destroy([]int{id}); err != nil {
		return err
	}

	for _, caption := range captions {
		if _, err := r.insert(id, caption); err != nil {
			return err
		}
	}

	return nil
}

type stringRepository struct {
	repository
	stringColumn string
}

func (r *stringRepository) get(id int) ([]string, error) {
	query := fmt.Sprintf("SELECT %s from %s WHERE %s = ?", r.stringColumn, r.tableName, r.idColumn)
	var ret []string
	err := r.queryFunc(query, []interface{}{id}, false, func(rows *sqlx.Rows) error {
		var out string
		if err := rows.Scan(&out); err != nil {
			return err
		}

		ret = append(ret, out)
		return nil
	})
	return ret, err
}

func (r *stringRepository) insert(id int, s string) (sql.Result, error) {
	stmt := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)", r.tableName, r.idColumn, r.stringColumn)
	return r.tx.Exec(stmt, id, s)
}

func (r *stringRepository) replace(id int, newStrings []string) error {
	if err := r.destroy([]int{id}); err != nil {
		return err
	}

	for _, s := range newStrings {
		if _, err := r.insert(id, s); err != nil {
			return err
		}
	}

	return nil
}

type stashIDRepository struct {
	repository
}

type stashIDs []*models.StashID

func (s *stashIDs) Append(o interface{}) {
	*s = append(*s, o.(*models.StashID))
}

func (s *stashIDs) New() interface{} {
	return &models.StashID{}
}

func (r *stashIDRepository) get(id int) ([]*models.StashID, error) {
	query := fmt.Sprintf("SELECT stash_id, endpoint from %s WHERE %s = ?", r.tableName, r.idColumn)
	var ret stashIDs
	err := r.query(query, []interface{}{id}, &ret)
	return []*models.StashID(ret), err
}

func (r *stashIDRepository) replace(id int, newIDs []models.StashID) error {
	if err := r.destroy([]int{id}); err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s, endpoint, stash_id) VALUES (?, ?, ?)", r.tableName, r.idColumn)
	for _, stashID := range newIDs {
		_, err := r.tx.Exec(query, id, stashID.Endpoint, stashID.StashID)
		if err != nil {
			return err
		}
	}
	return nil
}

func listKeys(i interface{}, addPrefix bool) string {
	var query []string
	v := reflect.ValueOf(i)
	for i := 0; i < v.NumField(); i++ {
		// Get key for struct tag
		rawKey := v.Type().Field(i).Tag.Get("db")
		key := strings.Split(rawKey, ",")[0]
		if key == "id" {
			continue
		}
		if addPrefix {
			key = ":" + key
		}
		query = append(query, key)
	}
	return strings.Join(query, ", ")
}

func updateSet(i interface{}, partial bool) string {
	var query []string
	v := reflect.ValueOf(i)
	for i := 0; i < v.NumField(); i++ {
		// Get key for struct tag
		rawKey := v.Type().Field(i).Tag.Get("db")
		key := strings.Split(rawKey, ",")[0]
		if key == "id" {
			continue
		}

		add := true
		if partial {
			reflectValue := reflect.ValueOf(v.Field(i).Interface())
			add = !reflectValue.IsNil()
		}

		if add {
			query = append(query, fmt.Sprintf("%s=:%s", key, key))
		}
	}
	return strings.Join(query, ", ")
}

func updateSetMap(m map[string]interface{}) string {
	var query []string
	for k := range m {
		query = append(query, fmt.Sprintf("%s=:%s", k, k))
	}
	return strings.Join(query, ", ")
}
