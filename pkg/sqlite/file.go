package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

const (
	fileTable    = "files"
	fileIDColumn = "file_id"

	// we need to resolve the zip_path
	fileSelectQuery = "SELECT files.*, zipfile.path as zip_path FROM files LEFT JOIN files zipfile ON zipfile.id = files.zip_file_id"
)

type fileQueryBuilder struct {
	repository
}

func NewFileReaderWriter(tx dbi) *fileQueryBuilder {
	return &fileQueryBuilder{
		repository{
			tx:        tx,
			tableName: fileTable,
			idColumn:  idColumn,
		},
	}
}

func (qb *fileQueryBuilder) Create(newObject models.File) (*models.File, error) {
	var ret models.File
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *fileQueryBuilder) UpdateFull(updatedFile models.File) (*models.File, error) {
	const partial = false
	if err := qb.update(updatedFile.ID, updatedFile, partial); err != nil {
		return nil, err
	}

	return qb.find(updatedFile.ID)
}

func (qb *fileQueryBuilder) Destroy(id int) error {
	return qb.destroyExisting([]int{id})
}

func (qb *fileQueryBuilder) Find(ids []int) ([]*models.File, error) {
	var files []*models.File
	for _, id := range ids {
		file, err := qb.find(id)
		if err != nil {
			return nil, err
		}

		if file == nil {
			return nil, fmt.Errorf("file with id %d not found", id)
		}

		files = append(files, file)
	}

	return files, nil
}

func (qb *fileQueryBuilder) get(id int, dest interface{}) error {
	stmt := fileSelectQuery + " WHERE files.id = ? LIMIT 1"
	if err := qb.tx.Get(dest, stmt, id); err != nil {
		return fmt.Errorf("executing SQL %q: %w", stmt, err)
	}

	return nil
}

func (qb *fileQueryBuilder) find(id int) (*models.File, error) {
	var ret models.File
	if err := qb.get(id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &ret, nil
}

func (qb *fileQueryBuilder) FindByChecksum(checksum string) ([]*models.File, error) {
	query := fileSelectQuery + " WHERE files.checksum = ?"
	args := []interface{}{checksum}
	return qb.queryFiles(query, args)
}

func (qb *fileQueryBuilder) FindByOSHash(oshash string) ([]*models.File, error) {
	query := fileSelectQuery + " WHERE files.oshash = ?"
	args := []interface{}{oshash}
	return qb.queryFiles(query, args)
}

func (qb *fileQueryBuilder) FindByPath(path string) (*models.File, error) {
	query := fileSelectQuery + " WHERE files.path = ? "
	args := []interface{}{path}
	query += "LIMIT 1"

	return qb.queryFile(query, args)
}

// func (qb *fileQueryBuilder) AllOfType(fileType models.FileType) ([]*models.File, error) {
// 	return qb.queryFiles(selectAll(fileTable)+`
// 		WHERE files.type = ?`, []interface{}{fileType})
// }

func (qb *fileQueryBuilder) queryFile(query string, args []interface{}) (*models.File, error) {
	results, err := qb.queryFiles(query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *fileQueryBuilder) queryFiles(query string, args []interface{}) ([]*models.File, error) {
	var ret models.Files
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.File(ret), nil
}
