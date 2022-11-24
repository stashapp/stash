package sqlite

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type Anonymiser struct {
	*Database
}

func NewAnonymiser(db *Database, outPath string) (*Anonymiser, error) {
	if _, err := db.db.Exec(fmt.Sprintf(`VACUUM INTO "%s"`, outPath)); err != nil {
		return nil, fmt.Errorf("vacuuming into %s: %w", outPath, err)
	}

	newDB := NewDatabase()
	if err := newDB.Open(outPath); err != nil {
		return nil, fmt.Errorf("opening %s: %w", outPath, err)
	}

	return &Anonymiser{Database: newDB}, nil
}

func (db *Anonymiser) Anonymise(ctx context.Context) error {
	defer db.Close()

	return utils.Do([]func() error{
		func() error { return db.deleteBlobs() },
		func() error { return db.anonymiseFolders(ctx) },
		func() error { return db.anonymiseFiles(ctx) },
		func() error { return db.anonymiseScenes(ctx) },
	})

	// anonymise fingerprints
	// anonymise scenes
	// anonymise images
	// anonymise galleries
	// anonymise performers
	// anonymise studios
	// anonymise tags
	// anonymise movies
}

func (db *Anonymiser) truncateTable(tableName string) error {
	_, err := db.db.Exec("DELETE FROM " + tableName)
	return err
}

func (db *Anonymiser) deleteBlobs() error {
	return utils.Do([]func() error{
		func() error { return db.truncateTable("scenes_cover") },
		func() error { return db.truncateTable("movies_images") },
		func() error { return db.truncateTable("performers_image") },
		func() error { return db.truncateTable("studios_image") },
		func() error { return db.truncateTable("tags_image") },
	})
}

func (db *Anonymiser) anonymiseFolders(ctx context.Context) error {
	return txn.WithTxn(ctx, db, func(ctx context.Context) error {
		return db.anonymiseFoldersRecurse(ctx, 0, "")
	})
}

func (db *Anonymiser) anonymiseFoldersRecurse(ctx context.Context, parentFolderID int, parentPath string) error {
	table := folderTableMgr.table

	stmt := dialect.Update(table)

	if parentFolderID == 0 {
		stmt = stmt.Set(goqu.Record{"path": goqu.Cast(table.Col(idColumn), "VARCHAR")}).Where(table.Col("parent_folder_id").IsNull())
	} else {
		stmt = stmt.Prepared(true).Set(goqu.Record{
			"path": goqu.L("? || ? || id", parentPath, string(filepath.Separator)),
		}).Where(table.Col("parent_folder_id").Eq(parentFolderID))
	}

	if _, err := exec(ctx, stmt); err != nil {
		return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
	}

	// now recurse to sub-folders
	query := dialect.From(table).Select(table.Col(idColumn), table.Col("path"))
	if parentFolderID == 0 {
		query = query.Where(table.Col("parent_folder_id").IsNull())
	} else {
		query = query.Where(table.Col("parent_folder_id").Eq(parentFolderID))
	}

	const single = false
	return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
		var id int
		var path string
		if err := rows.Scan(&id, &path); err != nil {
			return err
		}

		return db.anonymiseFoldersRecurse(ctx, id, path)
	})
}

func (db *Anonymiser) anonymiseFiles(ctx context.Context) error {
	return txn.WithTxn(ctx, db, func(ctx context.Context) error {
		table := fileTableMgr.table
		stmt := dialect.Update(table).Set(goqu.Record{"basename": goqu.Cast(table.Col(idColumn), "VARCHAR")})

		if _, err := exec(ctx, stmt); err != nil {
			return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
		}

		return nil
	})
}

func (db *Anonymiser) anonymiseScenes(ctx context.Context) error {
	table := sceneTableMgr.table
	lastID := 0
	total := 0

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("title"),
				table.Col("details"),
				table.Col("url"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id      int
					title   sql.NullInt64
					details sql.NullInt64
					url     sql.NullInt64
				)

				gotSome = true
				total++

				if err := rows.Scan(&id, &title, &details, &url); err != nil {
					return err
				}

				lastID = id

				// if title set set new title

				stmt := dialect.Update(table).Set(goqu.Record{
					"title":   db.obfuscateString(title),
					"details": db.obfuscateString(details),
					"url":     db.obfuscateString(url),
				}).Where(table.Col(idColumn).Eq(id))

				if _, err := exec(ctx, stmt); err != nil {
					return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	// anonymise distinct code and director
	if err := utils.Do([]func() error{
		func() error { return db.anonymiseText(ctx, table, "code") },
		func() error { return db.anonymiseText(ctx, table, "director") },
	}); err != nil {
		return err
	}

	return nil
}

func (db *Anonymiser) anonymiseText(ctx context.Context, table exp.IdentifierExpression, column string) error {
	lastID := 0
	total := 0

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			gotSome = false

			query := dialect.From(column).Select(
				goqu.ROW_NUMBER(),
				goqu.DISTINCT(column),
			).Where(goqu.ROW_NUMBER().Gt(lastID)).Limit(1000)

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					rowNumber int
					value     string
				)

				gotSome = true
				total++

				if err := rows.Scan(&rowNumber, &value); err != nil {
					return err
				}

				lastID = rowNumber

				set := goqu.Record{}
				set[column] = db.obfuscateString(value)

				stmt := dialect.Update(table).Set(set).Where(table.Col(column).Eq(value))

				if _, err := exec(ctx, stmt); err != nil {
					return fmt.Errorf("anonymising %s: %w", column, err)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	return nil
}

func (db *Anonymiser) obfuscateString(in string) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	out := strings.Builder{}
	for _, c := range in {
		if unicode.IsSpace(c) {
			out.WriteRune(c)
		} else {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
			if err != nil {
				panic("error generating random number")
			}

			out.WriteByte(letters[num.Int64()])
		}
	}

	return out.String()
}
