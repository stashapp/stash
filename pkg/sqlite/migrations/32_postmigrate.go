package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
	"gopkg.in/guregu/null.v4"
)

const legacyZipSeparator = "\x00"

func post32(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 32")

	m := schema32Migrator{
		migrator: migrator{
			db: db,
		},
		folderCache: make(map[string]folderInfo),
	}

	if err := m.migrateFolders(ctx); err != nil {
		return fmt.Errorf("migrating folders: %w", err)
	}

	if err := m.migrateFiles(ctx); err != nil {
		return fmt.Errorf("migrating files: %w", err)
	}

	if err := m.deletePlaceholderFolder(ctx); err != nil {
		return fmt.Errorf("deleting placeholder folder: %w", err)
	}

	return nil
}

type folderInfo struct {
	id    int
	zipID sql.NullInt64
}

type schema32Migrator struct {
	migrator
	folderCache map[string]folderInfo
}

func (m *schema32Migrator) migrateFolderSlashes(ctx context.Context) error {
	logger.Infof("Migrating folder slashes")
	const query = "SELECT `folders`.`id`, `folders`.`path` FROM `folders`"

	rows, err := m.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var p string

		err := rows.Scan(&id, &p)
		if err != nil {
			return err
		}

		convertedPath := filepath.ToSlash(p)

		_, err = m.db.Exec("UPDATE `folders` SET `path` = ? WHERE `id` = ?", convertedPath, id)
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (m *schema32Migrator) migrateFolders(ctx context.Context) error {
	if err := m.migrateFolderSlashes(ctx); err != nil {
		return err
	}

	logger.Infof("Migrating folders")

	const query = "SELECT `folders`.`id`, `folders`.`path` FROM `folders` INNER JOIN `galleries` ON `galleries`.`folder_id` = `folders`.`id`"

	rows, err := m.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var p string

		err := rows.Scan(&id, &p)
		if err != nil {
			return err
		}

		parent := path.Dir(p)
		parentID, zipFileID, err := m.createFolderHierarchy(parent)
		if err != nil {
			return err
		}

		_, err = m.db.Exec("UPDATE `folders` SET `parent_folder_id` = ?, `zip_file_id` = ? WHERE `id` = ?", parentID, zipFileID, id)
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (m *schema32Migrator) migrateFiles(ctx context.Context) error {
	const (
		limit    = 1000
		logEvery = 10000
	)
	offset := 0

	result := struct {
		Count int `db:"count"`
	}{0}

	if err := m.db.Get(&result, "SELECT COUNT(*) AS count FROM `files`"); err != nil {
		return err
	}

	logger.Infof("Migrating %d files...", result.Count)

	for {
		gotSome := false

		query := fmt.Sprintf("SELECT `id`, `basename` FROM `files` ORDER BY `id` LIMIT %d OFFSET %d", limit, offset)

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			rows, err := m.db.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				gotSome = true

				var id int
				var p string

				err := rows.Scan(&id, &p)
				if err != nil {
					return err
				}

				if strings.Contains(p, legacyZipSeparator) {
					// remove any null characters from the path
					p = strings.ReplaceAll(p, legacyZipSeparator, string(filepath.Separator))
				}

				convertedPath := filepath.ToSlash(p)
				parent := path.Dir(convertedPath)
				basename := path.Base(convertedPath)
				if parent != "." {
					parentID, zipFileID, err := m.createFolderHierarchy(parent)
					if err != nil {
						return err
					}

					_, err = m.db.Exec("UPDATE `files` SET `parent_folder_id` = ?, `zip_file_id` = ?, `basename` = ? WHERE `id` = ?", parentID, zipFileID, basename, id)
					if err != nil {
						return err
					}
				}
			}

			return rows.Err()
		}); err != nil {
			return err
		}

		if !gotSome {
			break
		}

		offset += limit

		if offset%logEvery == 0 {
			logger.Infof("Migrated %d files", offset)
		}
	}

	logger.Infof("Finished migrating files")

	return nil
}

func (m *schema32Migrator) deletePlaceholderFolder(ctx context.Context) error {
	// only delete the placeholder folder if no files/folders are attached to it
	result := struct {
		Count int `db:"count"`
	}{0}

	if err := m.db.Get(&result, "SELECT COUNT(*) AS count FROM `files` WHERE `parent_folder_id` = 1"); err != nil {
		return err
	}

	if result.Count > 0 {
		return fmt.Errorf("not deleting placeholder folder because it has %d files", result.Count)
	}

	result.Count = 0

	if err := m.db.Get(&result, "SELECT COUNT(*) AS count FROM `folders` WHERE `parent_folder_id` = 1"); err != nil {
		return err
	}

	if result.Count > 0 {
		return fmt.Errorf("not deleting placeholder folder because it has %d folders", result.Count)
	}

	_, err := m.db.Exec("DELETE FROM `folders` WHERE `id` = 1")
	return err
}

func (m *schema32Migrator) createFolderHierarchy(p string) (*int, sql.NullInt64, error) {
	parent := path.Dir(p)

	if parent == "." || parent == "/" {
		// get or create this folder
		return m.getOrCreateFolder(p, nil, sql.NullInt64{})
	}

	parentID, zipFileID, err := m.createFolderHierarchy(parent)
	if err != nil {
		return nil, sql.NullInt64{}, err
	}

	return m.getOrCreateFolder(p, parentID, zipFileID)
}

func (m *schema32Migrator) getOrCreateFolder(path string, parentID *int, zipFileID sql.NullInt64) (*int, sql.NullInt64, error) {
	foundEntry, ok := m.folderCache[path]
	if ok {
		return &foundEntry.id, foundEntry.zipID, nil
	}

	const query = "SELECT `id`, `zip_file_id` FROM `folders` WHERE `path` = ?"
	rows, err := m.db.Query(query, path)
	if err != nil {
		return nil, sql.NullInt64{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var zfid sql.NullInt64
		err := rows.Scan(&id, &zfid)
		if err != nil {
			return nil, sql.NullInt64{}, err
		}

		return &id, zfid, nil
	}

	if err := rows.Err(); err != nil {
		return nil, sql.NullInt64{}, err
	}

	const insertSQL = "INSERT INTO `folders` (`path`,`parent_folder_id`,`zip_file_id`,`mod_time`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)"

	var parentFolderID null.Int
	if parentID != nil {
		parentFolderID = null.IntFrom(int64(*parentID))
	}

	now := time.Now()
	result, err := m.db.Exec(insertSQL, path, parentFolderID, zipFileID, time.Time{}, now, now)
	if err != nil {
		return nil, sql.NullInt64{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, sql.NullInt64{}, err
	}

	idInt := int(id)

	m.folderCache[path] = folderInfo{id: idInt, zipID: zipFileID}

	return &idInt, zipFileID, nil
}

func init() {
	sqlite.RegisterPostMigration(32, post32)
}
