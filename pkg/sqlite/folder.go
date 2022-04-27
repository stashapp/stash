package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/file"
	"gopkg.in/guregu/null.v4"
)

const folderTable = "folders"

// path stores file paths in a platform-agnostic format and converts to platform-specific format for actual use.
type path string

func (p *path) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid path type %T", value)
	}

	*p = path(filepath.FromSlash(v))
	return nil
}

func (p path) Value() (driver.Value, error) {
	return string(filepath.ToSlash(string(p))), nil
}

type folderRow struct {
	ID file.FolderID `db:"id" goqu:"skipinsert"`
	// Path is stored in the OS-agnostic slash format
	Path           path      `db:"path"`
	ZipFileID      null.Int  `db:"zip_file_id"`
	ParentFolderID null.Int  `db:"parent_folder_id"`
	ModTime        time.Time `db:"mod_time"`
	MissingSince   null.Time `db:"missing_since"`
	LastScanned    time.Time `db:"last_scanned"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (r *folderRow) fromFolder(o file.Folder) {
	r.ID = o.ID
	r.Path = path(o.Path)
	r.ZipFileID = nullIntFromFileIDPtr(o.ZipFileID)
	r.ParentFolderID = nullIntFromFolderIDPtr(o.ParentFolderID)
	r.ModTime = o.ModTime
	r.MissingSince = null.TimeFromPtr(o.MissingSince)
	r.LastScanned = o.LastScanned
	r.CreatedAt = o.CreatedAt
	r.UpdatedAt = o.UpdatedAt
}

type folderQueryRow struct {
	folderRow
}

func (r *folderQueryRow) resolve() *file.Folder {
	ret := &file.Folder{
		ID: r.ID,
		DirEntry: file.DirEntry{
			ZipFileID:    nullIntFileIDPtr(r.ZipFileID),
			ModTime:      r.ModTime,
			MissingSince: r.MissingSince.Ptr(),
			LastScanned:  r.LastScanned,
		},
		Path:           string(r.Path),
		ParentFolderID: nullIntFolderIDPtr(r.ParentFolderID),
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}

	return ret
}

type folderQueryRows []folderQueryRow

func (r folderQueryRows) resolve() []*file.Folder {
	var ret []*file.Folder

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

func (qb *FolderStore) Create(ctx context.Context, f *file.Folder) error {
	var r folderRow
	r.fromFolder(*f)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	// only assign id once we are successful
	f.ID = file.FolderID(id)

	return nil
}

func (qb *FolderStore) Update(ctx context.Context, updatedObject *file.Folder) error {
	var r folderRow
	r.fromFolder(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *FolderStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *FolderStore) selectDataset() *goqu.SelectDataset {
	table := qb.table()
	return dialect.From(table)
}

func (qb *FolderStore) get(ctx context.Context, q *goqu.SelectDataset) (*file.Folder, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *FolderStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*file.Folder, error) {
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

func (qb *FolderStore) FindByPath(ctx context.Context, path string) (*file.Folder, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("path").Eq(path))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting folder by path %s: %w", path, err)
	}

	return ret, nil
}
