package migrations

import (
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

func post32(db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 32")

	m := schema32Migrator{db: db}

	if err := m.migrateFolders(); err != nil {
		return fmt.Errorf("migrating folders: %w", err)
	}

	if err := m.migrateFiles(); err != nil {
		return fmt.Errorf("migrating files: %w", err)
	}

	if err := m.deletePlaceholderFolder(); err != nil {
		return fmt.Errorf("deleting placeholder folder: %w", err)
	}

	return nil
}

type schema32Migrator struct {
	db *sqlx.DB
}

func (m *schema32Migrator) migrateFolders() error {
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

		convertedPath := filepath.ToSlash(p)
		parent := path.Dir(convertedPath)
		parentID, zipFileID, err := m.createFolderHierarchy(parent)
		if err != nil {
			return err
		}

		_, err = m.db.Exec("UPDATE `folders` SET `parent_folder_id` = ?, `zip_file_id` = ?, `path` = ? WHERE `id` = ?", parentID, zipFileID, convertedPath, id)
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (m *schema32Migrator) migrateFiles() error {
	const limit = 1000
	offset := 0
	for {
		query := fmt.Sprintf("SELECT `id`, `basename` FROM `files` ORDER BY `id` LIMIT %d OFFSET %d", limit, offset)
		offset += limit

		rows, err := m.db.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		gotSome := false
		for rows.Next() {
			gotSome = true

			var id int
			var p string

			err := rows.Scan(&id, &p)
			if err != nil {
				return err
			}

			if strings.Contains(p, legacyZipSeparator) {
				// if err := m.migrateFileInZip(id, p); err != nil {
				// 	return err
				// }

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

		if err := rows.Err(); err != nil {
			return err
		}

		if !gotSome {
			break
		}
	}

	return nil
}

func (m *schema32Migrator) deletePlaceholderFolder() error {
	// _, err := m.db.Exec("DELETE FROM `folders` WHERE `id` = 1")
	// return err
	return nil
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

	const insertSQL = "INSERT INTO `folders` (`path`,`parent_folder_id`,`zip_file_id`,`mod_time`,`last_scanned`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)"

	var parentFolderID null.Int
	if parentID != nil {
		parentFolderID = null.IntFrom(int64(*parentID))
	}

	now := time.Now()
	result, err := m.db.Exec(insertSQL, path, parentFolderID, zipFileID, time.Time{}, time.Time{}, now, now)
	if err != nil {
		return nil, sql.NullInt64{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, sql.NullInt64{}, err
	}

	idInt := int(id)
	return &idInt, zipFileID, nil
}

func init() {
	sqlite.RegisterCustomMigration(32, post32)
}
