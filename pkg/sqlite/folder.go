package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
)

const folderTable = "folders"

type folderRow struct {
	ID             models.FolderID `db:"id" goqu:"skipinsert"`
	Path           string          `db:"path"`
	ZipFileID      null.Int        `db:"zip_file_id"`
	ParentFolderID null.Int        `db:"parent_folder_id"`
	ModTime        Timestamp       `db:"mod_time"`
	CreatedAt      Timestamp       `db:"created_at"`
	UpdatedAt      Timestamp       `db:"updated_at"`
}

func (r *folderRow) fromFolder(o models.Folder) {
	r.ID = o.ID
	r.Path = o.Path
	r.ZipFileID = nullIntFromFileIDPtr(o.ZipFileID)
	r.ParentFolderID = nullIntFromFolderIDPtr(o.ParentFolderID)
	r.ModTime = Timestamp{Timestamp: o.ModTime}
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

type folderQueryRow struct {
	folderRow

	ZipBasename   null.String `db:"zip_basename"`
	ZipFolderPath null.String `db:"zip_folder_path"`
	ZipSize       null.Int    `db:"zip_size"`
}

func (r *folderQueryRow) resolve() *models.Folder {
	ret := &models.Folder{
		ID: r.ID,
		DirEntry: models.DirEntry{
			ZipFileID: nullIntFileIDPtr(r.ZipFileID),
			ModTime:   r.ModTime.Timestamp,
		},
		Path:           string(r.Path),
		ParentFolderID: nullIntFolderIDPtr(r.ParentFolderID),
		CreatedAt:      r.CreatedAt.Timestamp,
		UpdatedAt:      r.UpdatedAt.Timestamp,
	}

	if ret.ZipFileID != nil && r.ZipFolderPath.Valid && r.ZipBasename.Valid {
		ret.ZipFile = &models.BaseFile{
			ID:       *ret.ZipFileID,
			Path:     filepath.Join(r.ZipFolderPath.String, r.ZipBasename.String),
			Basename: r.ZipBasename.String,
			Size:     r.ZipSize.Int64,
		}
	}

	return ret
}

type folderQueryRows []folderQueryRow

func (r folderQueryRows) resolve() []*models.Folder {
	var ret []*models.Folder

	for _, row := range r {
		f := row.resolve()
		ret = append(ret, f)
	}

	return ret
}

type FolderStore struct {
	repository

	tableMgr *table
}

func NewFolderStore() *FolderStore {
	return &FolderStore{
		repository: repository{
			tableName: sceneTable,
			idColumn:  idColumn,
		},

		tableMgr: folderTableMgr,
	}
}

func (qb *FolderStore) Create(ctx context.Context, f *models.Folder) error {
	var r folderRow
	r.fromFolder(*f)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	// only assign id once we are successful
	f.ID = models.FolderID(id)

	return nil
}

func (qb *FolderStore) Update(ctx context.Context, updatedObject *models.Folder) error {
	var r folderRow
	r.fromFolder(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *FolderStore) Destroy(ctx context.Context, id models.FolderID) error {
	return qb.tableMgr.destroyExisting(ctx, []int{int(id)})
}

func (qb *FolderStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *FolderStore) selectDataset() *goqu.SelectDataset {
	table := qb.table()
	fileTable := fileTableMgr.table

	zipFileTable := fileTable.As("zip_files")
	zipFolderTable := table.As("zip_files_folders")

	cols := []interface{}{
		table.Col("id"),
		table.Col("path"),
		table.Col("zip_file_id"),
		table.Col("parent_folder_id"),
		table.Col("mod_time"),
		table.Col("created_at"),
		table.Col("updated_at"),
		zipFileTable.Col("basename").As("zip_basename"),
		zipFolderTable.Col("path").As("zip_folder_path"),
		// size is needed to open containing zip files
		zipFileTable.Col("size").As("zip_size"),
	}

	ret := dialect.From(table).Select(cols...)

	return ret.LeftJoin(
		zipFileTable,
		goqu.On(table.Col("zip_file_id").Eq(zipFileTable.Col("id"))),
	).LeftJoin(
		zipFolderTable,
		goqu.On(zipFileTable.Col("parent_folder_id").Eq(zipFolderTable.Col(idColumn))),
	)
}

func (qb *FolderStore) countDataset() *goqu.SelectDataset {
	table := qb.table()
	fileTable := fileTableMgr.table

	zipFileTable := fileTable.As("zip_files")
	zipFolderTable := table.As("zip_files_folders")

	ret := dialect.From(table).Select(goqu.COUNT(goqu.DISTINCT(table.Col("id"))))

	return ret.LeftJoin(
		zipFileTable,
		goqu.On(table.Col("zip_file_id").Eq(zipFileTable.Col("id"))),
	).LeftJoin(
		zipFolderTable,
		goqu.On(zipFileTable.Col("parent_folder_id").Eq(zipFolderTable.Col(idColumn))),
	)
}

func (qb *FolderStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Folder, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *FolderStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Folder, error) {
	const single = false
	var rows folderQueryRows
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f folderQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		rows = append(rows, f)
		return nil
	}); err != nil {
		return nil, err
	}

	return rows.resolve(), nil
}

func (qb *FolderStore) Find(ctx context.Context, id models.FolderID) (*models.Folder, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting folder by id %d: %w", id, err)
	}

	return ret, nil
}

func (qb *FolderStore) FindByPath(ctx context.Context, p string) (*models.Folder, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("path").Eq(p))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting folder by path %s: %w", p, err)
	}

	return ret, nil
}

func (qb *FolderStore) FindByParentFolderID(ctx context.Context, parentFolderID models.FolderID) ([]*models.Folder, error) {
	q := qb.selectDataset().Where(qb.table().Col("parent_folder_id").Eq(int(parentFolderID)))

	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting folders by parent folder id %d: %w", parentFolderID, err)
	}

	return ret, nil
}

func (qb *FolderStore) allInPaths(q *goqu.SelectDataset, p []string) *goqu.SelectDataset {
	table := qb.table()

	var conds []exp.Expression
	for _, pp := range p {
		ppWildcard := pp + string(filepath.Separator) + "%"

		conds = append(conds, table.Col("path").Eq(pp), table.Col("path").Like(ppWildcard))
	}

	return q.Where(
		goqu.Or(conds...),
	)
}

// FindAllInPaths returns the all folders that are or are within any of the given paths.
// Returns all if limit is < 0.
// Returns all folders if p is empty.
func (qb *FolderStore) FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]*models.Folder, error) {
	q := qb.selectDataset().Prepared(true)
	q = qb.allInPaths(q, p)

	if limit > -1 {
		q = q.Limit(uint(limit))
	}

	q = q.Offset(uint(offset))

	ret, err := qb.getMany(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting folders in path %s: %w", p, err)
	}

	return ret, nil
}

// CountAllInPaths returns a count of all folders that are within any of the given paths.
// Returns count of all folders if p is empty.
func (qb *FolderStore) CountAllInPaths(ctx context.Context, p []string) (int, error) {
	q := qb.countDataset().Prepared(true)
	q = qb.allInPaths(q, p)

	return count(ctx, q)
}

// func (qb *FolderStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*file.Folder, error) {
// 	table := qb.table()

// 	q := qb.selectDataset().Prepared(true).Where(
// 		table.Col(idColumn).Eq(
// 			sq,
// 		),
// 	)

// 	return qb.getMany(ctx, q)
// }

func (qb *FolderStore) FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Folder, error) {
	table := qb.table()

	q := qb.selectDataset().Prepared(true).Where(
		table.Col("zip_file_id").Eq(zipFileID),
	)

	return qb.getMany(ctx, q)
}
