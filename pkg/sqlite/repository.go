package sqlite

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

const idColumn = "id"

type objectConstructor func() interface{}

type objectList interface {
	Append(o interface{})
}

type repository struct {
	tx          *sqlx.Tx
	tableName   string
	idColumn    string
	constructor objectConstructor
}

func (r *repository) get(id int, dest interface{}) error {
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1", r.tableName, r.idColumn)
	return r.tx.Get(dest, stmt, id)
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
	// Perform query and fetch result
	result := struct {
		Int int `db:"count"`
	}{0}
	if err := r.tx.Get(&result, query, args...); err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return result.Int, nil
}

func (r *repository) runIdsQuery(query string, args []interface{}) ([]int, error) {
	var result []struct {
		Int int `db:"id"`
	}
	if err := r.tx.Select(&result, query, args...); err != nil && err != sql.ErrNoRows {
		return []int{}, err
	}

	vsm := make([]int, len(result))
	for i, v := range result {
		vsm[i] = v.Int
	}
	return vsm, nil
}

func (r *repository) query(query string, args []interface{}, out objectList) error {
	var rows *sqlx.Rows
	var err error

	rows, err = r.tx.Queryx(query, args...)

	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		object := r.constructor()
		if err := rows.StructScan(object); err != nil {
			return err
		}
		out.Append(object)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (r *repository) executeFindQuery(body string, args []interface{}, sortAndPagination string, whereClauses []string, havingClauses []string) ([]int, int, error) {
	if len(whereClauses) > 0 {
		body = body + " WHERE " + strings.Join(whereClauses, " AND ") // TODO handle AND or OR
	}
	body = body + " GROUP BY " + r.tableName + ".id "
	if len(havingClauses) > 0 {
		body = body + " HAVING " + strings.Join(havingClauses, " AND ") // TODO handle AND or OR
	}

	countQuery := r.buildCountQuery(body)
	idsQuery := body + sortAndPagination

	// Perform query and fetch result
	//logger.Tracef("SQL: %s, args: %v", idsQuery, args)

	countResult, countErr := r.runCountQuery(countQuery, args)
	idsResult, idsErr := r.runIdsQuery(idsQuery, args)

	if countErr != nil {
		return nil, 0, fmt.Errorf("Error executing count query with SQL: %s, args: %v, error: %s", countQuery, args, countErr.Error())
	}
	if idsErr != nil {
		return nil, 0, fmt.Errorf("Error executing find query with SQL: %s, args: %v, error: %s", idsQuery, args, idsErr.Error())
	}

	return idsResult, countResult, nil
}

func listKeys(i interface{}, addPrefix bool) string {
	var query []string
	v := reflect.ValueOf(i)
	for i := 0; i < v.NumField(); i++ {
		//get key for struct tag
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
		//get key for struct tag
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

func updateSetPartial(i interface{}) string {
	var query []string
	v := reflect.ValueOf(i)
	for i := 0; i < v.NumField(); i++ {
		//get key for struct tag
		rawKey := v.Type().Field(i).Tag.Get("db")
		key := strings.Split(rawKey, ",")[0]
		if key == "id" {
			continue
		}

	}
	return strings.Join(query, ", ")
}
