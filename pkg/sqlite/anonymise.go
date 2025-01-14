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
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	hex     = "0123456789abcdef"
)

type Anonymiser struct {
	*SQLiteDB
	sourceDb DBInterface
}

var anon_dialect = goqu.Dialect("sqlite3")

func NewAnonymiser(db DBInterface, outPath string) (*Anonymiser, error) {
	newDB := NewSQLiteDatabase(NewDatabase(), outPath, false)
	if err := newDB.Open(); err != nil {
		return nil, fmt.Errorf("opening %s: %w", outPath, err)
	}

	return &Anonymiser{SQLiteDB: newDB, sourceDb: db}, nil
}

func (db *Anonymiser) Anonymise(ctx context.Context) error {
	if err := func() error {
		defer db.Close()

		return utils.Do([]func() error{
			func() error { return db.fetch(ctx) },
			func() error { return db.deleteBlobs() },
			func() error { return db.deleteStashIDs() },
			func() error { return db.clearOHistory() },
			func() error { return db.clearWatchHistory() },
			func() error { return db.anonymiseFolders(ctx) },
			func() error { return db.anonymiseFiles(ctx) },
			func() error { return db.anonymiseCaptions(ctx) },
			func() error { return db.anonymiseFingerprints(ctx) },
			func() error { return db.anonymiseScenes(ctx) },
			func() error { return db.anonymiseMarkers(ctx) },
			func() error { return db.anonymiseImages(ctx) },
			func() error { return db.anonymiseGalleries(ctx) },
			func() error { return db.anonymisePerformers(ctx) },
			func() error { return db.anonymiseStudios(ctx) },
			func() error { return db.anonymiseTags(ctx) },
			func() error { return db.anonymiseGroups(ctx) },
			func() error { return db.anonymiseSavedFilters(ctx) },
			func() error { return db.Optimise(ctx) },
		})
	}(); err != nil {
		// delete the database
		_ = db.Remove()

		return err
	}

	return nil
}

func (db *Anonymiser) fetch(ctx context.Context) error {
	const disableForeignKeys = true
	const writable = true
	const batchSize = 5000

	conn, err := db.open(disableForeignKeys, writable)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	for _, table := range []exp.IdentifierExpression{
		goqu.I(fileTable),
		goqu.I(fingerprintTable),
		goqu.I(folderTable),
		goqu.I(galleryTable),
		goqu.I(galleriesChaptersTable),
		goqu.I(galleriesFilesTable),
		goqu.I(galleriesImagesTable),
		goqu.I(galleriesTagsTable),
		goqu.I(galleriesURLsTable),
		goqu.I(groupURLsTable),
		goqu.I(groupTable),
		goqu.I(groupRelationsTable),
		goqu.I(groupsScenesTable),
		goqu.I(groupsTagsTable),
		goqu.I(imageFileTable),
		goqu.I(imagesURLsTable),
		goqu.I(imageTable),
		goqu.I(imagesFilesTable),
		goqu.I(imagesTagsTable),
		goqu.I(performersAliasesTable),
		goqu.I("performer_stash_ids"),
		goqu.I(performerURLsTable),
		goqu.I(performerTable),
		goqu.I(performersGalleriesTable),
		goqu.I(performersImagesTable),
		goqu.I(performersScenesTable),
		goqu.I(performersTagsTable),
		goqu.I(savedFilterTable),
		goqu.I(sceneMarkerTable),
		goqu.I("scene_markers_tags"),
		goqu.I(scenesURLsTable),
		goqu.I(sceneTable),
		goqu.I(scenesFilesTable),
		goqu.I(scenesGalleriesTable),
		goqu.I(scenesODatesTable),
		goqu.I(scenesTagsTable),
		goqu.I(scenesViewDatesTable),
		goqu.I(studioAliasesTable),
		goqu.I("studio_stash_ids"),
		goqu.I(studioTable),
		goqu.I(studiosTagsTable),
		goqu.I(tagAliasesTable),
		goqu.I(tagTable),
		goqu.I(tagRelationsTable),
		goqu.I(videoCaptionsTable),
		goqu.I(videoFileTable),
	} {
		offset := 0
		for {
			q := dialect.From(table).Select(table.All()).Limit(uint(batchSize)).Offset(uint(offset))
			var rowsSlice []map[string]interface{}

			// Fetch
			if err := txn.WithTxn(ctx, db.sourceDb, func(ctx context.Context) error {
				if err := queryFunc(ctx, q, false, func(r *sqlx.Rows) error {
					for r.Next() {
						row := make(map[string]interface{})
						if err := r.MapScan(row); err != nil {
							return fmt.Errorf("failed structscan: %w", err)
						}
						rowsSlice = append(rowsSlice, row)
					}

					return nil
				}); err != nil {
					return fmt.Errorf("querying %s: %w", table, err)
				}

				return nil
			}); err != nil {
				return fmt.Errorf("failed fetch transaction: %w", err)
			}

			if len(rowsSlice) == 0 {
				break
			}

			// Insert
			txn, err := conn.BeginTx(ctx, nil)
			if err != nil {
				return fmt.Errorf("begin tx: %w", err)
			}

			i := anon_dialect.Insert(table).Rows(rowsSlice)
			sql, args, err := i.ToSQL()
			if err != nil {
				return fmt.Errorf("failed tosql: %w", err)
			}

			_, err = txn.ExecContext(ctx, sql, args...)
			if err != nil {
				return fmt.Errorf("exec `%s` [%v]: %w", sql, args, err)
			}

			if err := txn.Commit(); err != nil {
				return fmt.Errorf("commit: %w", err)
			}

			// Move to the next batch
			offset += batchSize
		}
	}

	if err := conn.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

func (db *Anonymiser) truncateColumn(tableName string, column string) error {
	_, err := db.writeDB.Exec("UPDATE " + tableName + " SET " + column + " = NULL")
	return err
}

func (db *Anonymiser) truncateTable(tableName string) error {
	_, err := db.writeDB.Exec("DELETE FROM " + tableName)
	return err
}

func (db *Anonymiser) deleteBlobs() error {
	return utils.Do([]func() error{
		func() error { return db.truncateColumn(tagTable, tagImageBlobColumn) },
		func() error { return db.truncateColumn(studioTable, studioImageBlobColumn) },
		func() error { return db.truncateColumn(performerTable, performerImageBlobColumn) },
		func() error { return db.truncateColumn(sceneTable, sceneCoverBlobColumn) },
		func() error { return db.truncateColumn(groupTable, groupFrontImageBlobColumn) },
		func() error { return db.truncateColumn(groupTable, groupBackImageBlobColumn) },

		func() error { return db.truncateTable(blobTable) },
	})
}

func (db *Anonymiser) deleteStashIDs() error {
	return utils.Do([]func() error{
		func() error { return db.truncateTable("scene_stash_ids") },
		func() error { return db.truncateTable("studio_stash_ids") },
		func() error { return db.truncateTable("performer_stash_ids") },
	})
}

func (db *Anonymiser) clearOHistory() error {
	return utils.Do([]func() error{
		func() error { return db.truncateTable(scenesODatesTable) },
	})
}

func (db *Anonymiser) clearWatchHistory() error {
	return utils.Do([]func() error{
		func() error { return db.truncateTable(scenesViewDatesTable) },
	})
}

func (db *Anonymiser) anonymiseFolders(ctx context.Context) error {
	logger.Infof("Anonymising folders")
	return txn.WithTxn(ctx, db, func(ctx context.Context) error {
		return db.anonymiseFoldersRecurse(ctx, 0, "")
	})
}

func (db *Anonymiser) anonymiseFoldersRecurse(ctx context.Context, parentFolderID int, parentPath string) error {
	table := folderTableMgr.table

	stmt := anon_dialect.Update(table)

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
	query := anon_dialect.From(table).Select(table.Col(idColumn), table.Col("path"))
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
	logger.Infof("Anonymising files")
	return txn.WithTxn(ctx, db, func(ctx context.Context) error {
		table := fileTableMgr.table
		stmt := anon_dialect.Update(table).Set(goqu.Record{"basename": goqu.Cast(table.Col(idColumn), "VARCHAR")})

		if _, err := exec(ctx, stmt); err != nil {
			return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
		}

		return nil
	})
}

func (db *Anonymiser) anonymiseCaptions(ctx context.Context) error {
	logger.Infof("Anonymising captions")
	return txn.WithTxn(ctx, db, func(ctx context.Context) error {
		table := goqu.T(videoCaptionsTable)
		stmt := anon_dialect.Update(table).Set(goqu.Record{"filename": goqu.Cast(table.Col("file_id"), "VARCHAR")})

		if _, err := exec(ctx, stmt); err != nil {
			return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
		}

		return nil
	})
}

func (db *Anonymiser) anonymiseFingerprints(ctx context.Context) error {
	logger.Infof("Anonymising fingerprints")
	table := fingerprintTableMgr.table
	lastID := 0
	lastType := ""
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(fileIDColumn),
				table.Col("type"),
				table.Col("fingerprint"),
			).Where(goqu.L("(file_id, type)").Gt(goqu.L("(?, ?)", lastID, lastType))).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id          int
					typ         string
					fingerprint string
				)

				if err := rows.Scan(
					&id,
					&typ,
					&fingerprint,
				); err != nil {
					return err
				}

				if err := db.anonymiseFingerprint(ctx, table, "fingerprint", fingerprint); err != nil {
					return err
				}

				lastID = id
				lastType = typ

				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d fingerprints", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	return nil
}

func (db *Anonymiser) anonymiseScenes(ctx context.Context) error {
	logger.Infof("Anonymising scenes")
	table := sceneTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("title"),
				table.Col("details"),
				table.Col("code"),
				table.Col("director"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id       int
					title    sql.NullString
					details  sql.NullString
					code     sql.NullString
					director sql.NullString
				)

				if err := rows.Scan(
					&id,
					&title,
					&details,
					&code,
					&director,
				); err != nil {
					return err
				}

				set := goqu.Record{}

				// if title set set new title
				db.obfuscateNullString(set, "title", title)
				db.obfuscateNullString(set, "details", details)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(table.Col(idColumn).Eq(id))

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				if code.Valid {
					if err := db.anonymiseText(ctx, table, "code", code.String); err != nil {
						return err
					}
				}

				if director.Valid {
					if err := db.anonymiseText(ctx, table, "director", director.String); err != nil {
						return err
					}
				}

				lastID = id
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d scenes", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	if err := db.anonymiseURLs(ctx, goqu.T(scenesURLsTable), "scene_id"); err != nil {
		return err
	}

	return nil
}

func (db *Anonymiser) anonymiseMarkers(ctx context.Context) error {
	logger.Infof("Anonymising scene markers")
	table := sceneMarkerTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("title"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id    int
					title string
				)

				if err := rows.Scan(
					&id,
					&title,
				); err != nil {
					return err
				}

				if err := db.anonymiseText(ctx, table, "title", title); err != nil {
					return err
				}

				lastID = id
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d scene markers", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	return nil
}

func (db *Anonymiser) anonymiseImages(ctx context.Context) error {
	logger.Infof("Anonymising images")
	table := imageTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("title"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id    int
					title sql.NullString
				)

				if err := rows.Scan(
					&id,
					&title,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "title", title)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(table.Col(idColumn).Eq(id))

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d images", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	if err := db.anonymiseURLs(ctx, goqu.T(imagesURLsTable), "image_id"); err != nil {
		return err
	}

	return nil
}

func (db *Anonymiser) anonymiseGalleries(ctx context.Context) error {
	logger.Infof("Anonymising galleries")
	table := galleryTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("title"),
				table.Col("details"),
				table.Col("photographer"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id           int
					title        sql.NullString
					details      sql.NullString
					photographer sql.NullString
				)

				if err := rows.Scan(
					&id,
					&title,
					&details,
					&photographer,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "title", title)
				db.obfuscateNullString(set, "details", details)
				db.obfuscateNullString(set, "photographer", photographer)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(table.Col(idColumn).Eq(id))

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d galleries", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	if err := db.anonymiseURLs(ctx, goqu.T(galleriesURLsTable), "gallery_id"); err != nil {
		return err
	}

	return nil
}

func (db *Anonymiser) anonymisePerformers(ctx context.Context) error {
	logger.Infof("Anonymising performers")
	table := performerTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("name"),
				table.Col("disambiguation"),
				table.Col("details"),
				table.Col("tattoos"),
				table.Col("piercings"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id             int
					name           sql.NullString
					disambiguation sql.NullString
					details        sql.NullString
					tattoos        sql.NullString
					piercings      sql.NullString
				)

				if err := rows.Scan(
					&id,
					&name,
					&disambiguation,
					&details,
					&tattoos,
					&piercings,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "name", name)
				db.obfuscateNullString(set, "disambiguation", disambiguation)
				db.obfuscateNullString(set, "details", details)
				db.obfuscateNullString(set, "tattoos", tattoos)
				db.obfuscateNullString(set, "piercings", piercings)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(table.Col(idColumn).Eq(id))

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d performers", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	if err := db.anonymiseAliases(ctx, goqu.T(performersAliasesTable), "performer_id"); err != nil {
		return err
	}

	if err := db.anonymiseURLs(ctx, goqu.T(performerURLsTable), "performer_id"); err != nil {
		return err
	}

	if err := db.anonymiseCustomFields(ctx, goqu.T(performersCustomFieldsTable.GetTable()), "performer_id"); err != nil {
		return err
	}

	return nil
}

func (db *Anonymiser) anonymiseStudios(ctx context.Context) error {
	logger.Infof("Anonymising studios")
	table := studioTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("name"),
				table.Col("url"),
				table.Col("details"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id      int
					name    sql.NullString
					url     sql.NullString
					details sql.NullString
				)

				if err := rows.Scan(
					&id,
					&name,
					&url,
					&details,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "name", name)
				db.obfuscateNullString(set, "url", url)
				db.obfuscateNullString(set, "details", details)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(table.Col(idColumn).Eq(id))

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				gotSome = true
				total++

				// TODO - anonymise studio aliases

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d studios", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	if err := db.anonymiseAliases(ctx, goqu.T(studioAliasesTable), "studio_id"); err != nil {
		return err
	}

	return nil
}

func (db *Anonymiser) anonymiseAliases(ctx context.Context, table exp.IdentifierExpression, idColumn string) error {
	lastID := 0
	lastAlias := ""
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("alias"),
			).Where(goqu.L("(" + idColumn + ", alias)").Gt(goqu.L("(?, ?)", lastID, lastAlias))).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id    int
					alias sql.NullString
				)

				if err := rows.Scan(
					&id,
					&alias,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "alias", alias)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(
						table.Col(idColumn).Eq(id),
						table.Col("alias").Eq(alias),
					)

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				lastAlias = alias.String
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d %s aliases", total, table.GetTable())
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	return nil
}

func (db *Anonymiser) anonymiseURLs(ctx context.Context, table exp.IdentifierExpression, idColumn string) error {
	lastID := 0
	lastURL := ""
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("url"),
			).Where(goqu.L("(" + idColumn + ", url)").Gt(goqu.L("(?, ?)", lastID, lastURL))).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id  int
					url sql.NullString
				)

				if err := rows.Scan(
					&id,
					&url,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "url", url)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(
						table.Col(idColumn).Eq(id),
						table.Col("url").Eq(url),
					)

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				lastURL = url.String
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d %s URLs", total, table.GetTable())
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	return nil
}

func (db *Anonymiser) anonymiseTags(ctx context.Context) error {
	logger.Infof("Anonymising tags")
	table := tagTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("name"),
				table.Col("description"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id          int
					name        sql.NullString
					description sql.NullString
				)

				if err := rows.Scan(
					&id,
					&name,
					&description,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "name", name)
				db.obfuscateNullString(set, "description", description)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(table.Col(idColumn).Eq(id))

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d tags", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	if err := db.anonymiseAliases(ctx, goqu.T(tagAliasesTable), "tag_id"); err != nil {
		return err
	}

	return nil
}

func (db *Anonymiser) anonymiseGroups(ctx context.Context) error {
	logger.Infof("Anonymising groups")
	table := groupTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("name"),
				table.Col("aliases"),
				table.Col("description"),
				table.Col("director"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id          int
					name        sql.NullString
					aliases     sql.NullString
					description sql.NullString
					director    sql.NullString
				)

				if err := rows.Scan(
					&id,
					&name,
					&aliases,
					&description,
					&director,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "name", name)
				db.obfuscateNullString(set, "aliases", aliases)
				db.obfuscateNullString(set, "description", description)
				db.obfuscateNullString(set, "director", director)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(table.Col(idColumn).Eq(id))

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d groups", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	if err := db.anonymiseURLs(ctx, goqu.T(groupURLsTable), "group_id"); err != nil {
		return err
	}

	return nil
}

func (db *Anonymiser) anonymiseSavedFilters(ctx context.Context) error {
	logger.Infof("Anonymising saved filters")
	table := savedFilterTableMgr.table
	lastID := 0
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := anon_dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("name"),
			).Where(table.Col(idColumn).Gt(lastID)).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id   int
					name sql.NullString
				)

				if err := rows.Scan(
					&id,
					&name,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				db.obfuscateNullString(set, "name", name)

				if len(set) > 0 {
					stmt := anon_dialect.Update(table).Set(set).Where(table.Col(idColumn).Eq(id))

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d saved filters", total)
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	return nil
}

func (db *Anonymiser) anonymiseText(ctx context.Context, table exp.IdentifierExpression, column string, value string) error {
	set := goqu.Record{}
	set[column] = db.obfuscateString(value, letters)

	stmt := anon_dialect.Update(table).Set(set).Where(table.Col(column).Eq(value))

	if _, err := exec(ctx, stmt); err != nil {
		return fmt.Errorf("anonymising %s: %w", column, err)
	}

	return nil
}

func (db *Anonymiser) anonymiseFingerprint(ctx context.Context, table exp.IdentifierExpression, column string, value string) error {
	set := goqu.Record{}
	set[column] = db.obfuscateString(value, hex)

	stmt := anon_dialect.Update(table).Set(set).Where(table.Col(column).Eq(value))

	if _, err := exec(ctx, stmt); err != nil {
		return fmt.Errorf("anonymising %s: %w", column, err)
	}

	return nil
}

func (db *Anonymiser) obfuscateNullString(out goqu.Record, column string, in sql.NullString) {
	if in.Valid {
		out[column] = db.obfuscateString(in.String, letters)
	}
}

func (db *Anonymiser) obfuscateString(in string, dict string) string {
	out := strings.Builder{}
	for _, c := range in {
		if unicode.IsSpace(c) {
			out.WriteRune(c)
		} else {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(dict))))
			if err != nil {
				panic("error generating random number")
			}

			out.WriteByte(dict[num.Int64()])
		}
	}

	return out.String()
}

func (db *Anonymiser) anonymiseCustomFields(ctx context.Context, table exp.IdentifierExpression, idColumn string) error {
	lastID := 0
	lastField := ""
	total := 0
	const logEvery = 10000

	for gotSome := true; gotSome; {
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			query := dialect.From(table).Select(
				table.Col(idColumn),
				table.Col("field"),
				table.Col("value"),
			).Where(
				goqu.L("("+idColumn+", field)").Gt(goqu.L("(?, ?)", lastID, lastField)),
			).Order(
				table.Col(idColumn).Asc(), table.Col("field").Asc(),
			).Limit(1000)

			gotSome = false

			const single = false
			return queryFunc(ctx, query, single, func(rows *sqlx.Rows) error {
				var (
					id    int
					field string
					value string
				)

				if err := rows.Scan(
					&id,
					&field,
					&value,
				); err != nil {
					return err
				}

				set := goqu.Record{}
				set["field"] = db.obfuscateString(field, letters)
				set["value"] = db.obfuscateString(value, letters)

				if len(set) > 0 {
					stmt := dialect.Update(table).Set(set).Where(
						table.Col(idColumn).Eq(id),
						table.Col("field").Eq(field),
					)

					if _, err := exec(ctx, stmt); err != nil {
						return fmt.Errorf("anonymising %s: %w", table.GetTable(), err)
					}
				}

				lastID = id
				lastField = field
				gotSome = true
				total++

				if total%logEvery == 0 {
					logger.Infof("Anonymised %d %s custom fields", total, table.GetTable())
				}

				return nil
			})
		}); err != nil {
			return err
		}
	}

	return nil
}
