package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
	"gopkg.in/guregu/null.v4"
)

func post84(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 76")

	m := schema84Migrator{
		migrator: migrator{
			db: db,
		},
		folderCache: make(map[string]folderInfo),
	}

	rootPaths := config.GetInstance().GetStashPaths().Paths()

	if err := m.createMissingFolderHierarchies(ctx, rootPaths); err != nil {
		return fmt.Errorf("creating missing folder hierarchies: %w", err)
	}

	if err := m.migrateFolders(ctx); err != nil {
		return fmt.Errorf("migrating folders: %w", err)
	}

	return nil
}

type schema84Migrator struct {
	migrator
	folderCache map[string]folderInfo
}

func (m *schema84Migrator) createMissingFolderHierarchies(ctx context.Context, rootPaths []string) error {
	// before we set the basenames, we need to address any folders that are missing their
	// parent folders.
	const (
		limit    = 1000
		logEvery = 10000
	)

	lastID := 0
	count := 0
	logged := false

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := "SELECT `folders`.`id`, `folders`.`path` FROM `folders` WHERE `folders`.`parent_folder_id` IS NULL "

			if lastID != 0 {
				query += fmt.Sprintf("AND `folders`.`id` > %d ", lastID)
			}

			query += fmt.Sprintf("ORDER BY `folders`.`id` LIMIT %d", limit)

			rows, err := tx.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				// log once if we find any folders with missing parent folders
				if !logged {
					logger.Info("Migrating folders with missing parents...")
					logged = true
				}

				var id int
				var p string

				err := rows.Scan(&id, &p)
				if err != nil {
					return err
				}

				lastID = id
				gotSome = true
				count++

				// don't try to create parent folders for root paths
				if slices.Contains(rootPaths, p) {
					continue
				}

				parentDir := filepath.Dir(p)
				if parentDir == p {
					// this can happen if the path is something like "C:\", where the parent directory is the same as the current directory
					continue
				}

				parentID, err := m.getOrCreateFolderHierarchy(tx, parentDir, rootPaths)
				if err != nil {
					return fmt.Errorf("error creating parent folder for folder %d %q: %w", id, p, err)
				}

				if parentID == nil {
					continue
				}

				// now set the parent folder ID for the current folder
				logger.Debugf("Migrating folder %d %q: setting parent folder ID to %d", id, p, *parentID)

				_, err = tx.Exec("UPDATE `folders` SET `parent_folder_id` = ? WHERE `id` = ?", *parentID, id)
				if err != nil {
					return fmt.Errorf("error setting parent folder for folder %d %q: %w", id, p, err)
				}
			}

			return rows.Err()
		}); err != nil {
			return err
		}

		if !gotSome {
			break
		}

		if count%logEvery == 0 {
			logger.Infof("Migrated %d folders", count)
		}
	}

	return nil
}

func (m *schema84Migrator) findFolderByPath(tx *sqlx.Tx, path string) (*int, error) {
	query := "SELECT `folders`.`id` FROM `folders` WHERE `folders`.`path` = ?"

	var id int
	if err := tx.Get(&id, query, path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &id, nil
}

// this is a copy of the GetOrCreateFolderHierarchy function from pkg/file/folder.go,
// but modified to use low-level SQL queries instead of the models.FolderFinderCreator interface, to avoid
func (m *schema84Migrator) getOrCreateFolderHierarchy(tx *sqlx.Tx, path string, rootPaths []string) (*int, error) {
	// get or create folder hierarchy
	folderID, err := m.findFolderByPath(tx, path)
	if err != nil {
		return nil, err
	}

	if folderID == nil {
		var parentID *int

		if !slices.Contains(rootPaths, path) {
			parentPath := filepath.Dir(path)

			// it's possible that the parent path is the same as the current path, if there are folders outside
			// of the root paths. In that case, we should just return nil for the parent ID.
			if parentPath == path {
				return nil, nil
			}

			parentID, err = m.getOrCreateFolderHierarchy(tx, parentPath, rootPaths)
			if err != nil {
				return nil, err
			}
		}

		logger.Debugf("%s doesn't exist. Creating new folder entry...", path)

		// we need to set basename to path, which will be addressed in the next step
		const insertSQL = "INSERT INTO `folders` (`path`,`basename`,`parent_folder_id`,`mod_time`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)"

		var parentFolderID null.Int
		if parentID != nil {
			parentFolderID = null.IntFrom(int64(*parentID))
		}

		now := time.Now()
		result, err := tx.Exec(insertSQL, path, path, parentFolderID, time.Time{}, now, now)
		if err != nil {
			return nil, fmt.Errorf("creating folder %s: %w", path, err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("creating folder %s: %w", path, err)
		}

		idInt := int(id)
		folderID = &idInt
	}

	return folderID, nil
}

func (m *schema84Migrator) migrateFolders(ctx context.Context) error {
	const (
		limit    = 1000
		logEvery = 10000
	)

	lastID := 0
	count := 0
	logged := false

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := "SELECT `folders`.`id`, `folders`.`path` FROM `folders` "

			if lastID != 0 {
				query += fmt.Sprintf("WHERE `folders`.`id` > %d ", lastID)
			}

			query += fmt.Sprintf("ORDER BY `folders`.`id` LIMIT %d", limit)

			rows, err := tx.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				if !logged {
					logger.Infof("Migrating folders to set basenames...")
					logged = true
				}

				var id int
				var p string

				err := rows.Scan(&id, &p)
				if err != nil {
					return err
				}

				lastID = id
				gotSome = true
				count++

				basename := filepath.Base(p)
				logger.Debugf("Migrating folder %d %q: setting basename to %q", id, p, basename)
				_, err = tx.Exec("UPDATE `folders` SET `basename` = ? WHERE `id` = ?", basename, id)
				if err != nil {
					return fmt.Errorf("error migrating folder %d %q: %w", id, p, err)
				}
			}

			return rows.Err()
		}); err != nil {
			return err
		}

		if !gotSome {
			break
		}

		if count%logEvery == 0 {
			logger.Infof("Migrated %d folders", count)
		}
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(84, post84)
}
