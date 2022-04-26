package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/file"
	"gopkg.in/guregu/null.v4"
)

const folderTable = "folders"

type folderRow struct {
	ID             file.FolderID `db:"id" goqu:"skipinsert"`
	Path           string        `db:"path"`
	ZipFileID      null.Int      `db:"zip_file_id"`
	ParentFolderID null.Int      `db:"parent_folder_id"`
	ModTime        time.Time     `db:"mod_time"`
	MissingSince   null.Time     `db:"missing_since"`
	LastScanned    time.Time     `db:"last_scanned"`
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at"`
}

func (r *folderRow) fromFolder(o file.Folder) {
	r.ID = o.ID
	r.Path = o.Path
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
			Path:           r.Path,
			ZipFileID:      nullIntFileIDPtr(r.ZipFileID),
			ParentFolderID: nullIntFolderIDPtr(r.ParentFolderID),
			ModTime:        r.ModTime,
			MissingSince:   r.MissingSince.Ptr(),
			LastScanned:    r.LastScanned,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
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

type folderQueryBuilder struct {
	repository

	tableMgr *table
}

var FolderReaderWriter = &folderQueryBuilder{
	repository: repository{
		tableName: sceneTable,
		idColumn:  idColumn,
	},

	tableMgr: folderTableMgr,
}

func (qb *folderQueryBuilder) Create(ctx context.Context, f *file.Folder) error {
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

func (qb *folderQueryBuilder) Update(ctx context.Context, updatedObject *file.Folder) error {
	var r folderRow
	r.fromFolder(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	return nil
}

func (qb *folderQueryBuilder) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *folderQueryBuilder) selectDataset() *goqu.SelectDataset {
	table := qb.table()
	return dialect.From(table)
}

func (qb *folderQueryBuilder) get(ctx context.Context, q *goqu.SelectDataset) (*file.Folder, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *folderQueryBuilder) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*file.Folder, error) {
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

func (qb *folderQueryBuilder) FindByPath(ctx context.Context, path string) (*file.Folder, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("path").Eq(path))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting folder by path %s: %w", path, err)
	}

	return ret, nil
}
