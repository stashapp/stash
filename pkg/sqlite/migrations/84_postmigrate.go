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
	logger.Info("Running post-migration for schema version 84")

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

	if err := m.fixIncorrectParents(ctx, rootPaths); err != nil {
		return fmt.Errorf("fixing incorrect parent folders: %w", err)
	}

	if err := m.deduplicateFolders(ctx); err != nil {
		return fmt.Errorf("deduplicating folders: %w", err)
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

func (m *schema84Migrator) fixIncorrectParents(ctx context.Context, rootPaths []string) error {
	const (
		limit    = 1000
		logEvery = 10000
	)

	lastID := 0
	count := 0
	fixed := 0
	logged := false

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := "SELECT f.id, f.path, f.parent_folder_id, pf.path AS parent_path " +
				"FROM folders f " +
				"JOIN folders pf ON f.parent_folder_id = pf.id "

			if lastID != 0 {
				query += fmt.Sprintf("WHERE f.id > %d ", lastID)
			}

			query += fmt.Sprintf("ORDER BY f.id LIMIT %d", limit)

			rows, err := tx.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				var p string
				var parentFolderID int
				var parentPath string

				err := rows.Scan(&id, &p, &parentFolderID, &parentPath)
				if err != nil {
					return err
				}

				lastID = id
				gotSome = true
				count++

				expectedParent := filepath.Dir(p)
				if expectedParent == parentPath {
					continue
				}

				if !logged {
					logger.Info("Fixing folders with incorrect parent folder assignments...")
					logged = true
				}

				correctParentID, err := m.getOrCreateFolderHierarchy(tx, expectedParent, rootPaths)
				if err != nil {
					return fmt.Errorf("error getting/creating correct parent for folder %d %q: %w", id, p, err)
				}

				if correctParentID == nil {
					continue
				}

				logger.Debugf("Fixing folder %d %q: changing parent_folder_id from %d to %d", id, p, parentFolderID, *correctParentID)

				_, err = tx.Exec("UPDATE `folders` SET `parent_folder_id` = ? WHERE `id` = ?", *correctParentID, id)
				if err != nil {
					return fmt.Errorf("error fixing parent folder for folder %d %q: %w", id, p, err)
				}

				fixed++
			}

			return rows.Err()
		}); err != nil {
			return err
		}

		if !gotSome {
			break
		}

		if count%logEvery == 0 {
			logger.Infof("Checked %d folders", count)
		}
	}

	if fixed > 0 {
		logger.Infof("Fixed %d folders with incorrect parent assignments", fixed)
	}

	return nil
}

// deduplicateFolders finds folders that would have the same (parent_folder_id, basename) after
// migrateFolders sets basename = filepath.Base(path), and merges the duplicates.
// This can happen when the database contains entries for the same physical folder with different
// path representations (e.g., mixed separators like "\data/movies" vs "\data\movies" on Windows).
func (m *schema84Migrator) deduplicateFolders(ctx context.Context) error {
	for {
		n, err := m.deduplicateFoldersPass(ctx)
		if err != nil {
			return err
		}
		// repeat until no more duplicates are found, since merging child folders
		// from a duplicate parent into the canonical parent may create new conflicts
		if n == 0 {
			break
		}
	}
	return nil
}

func (m *schema84Migrator) deduplicateFoldersPass(ctx context.Context) (int, error) {
	type folderRow struct {
		ID             int    `db:"id"`
		Path           string `db:"path"`
		ParentFolderID int    `db:"parent_folder_id"`
	}

	var folders []folderRow
	if err := m.db.SelectContext(ctx, &folders,
		"SELECT id, path, parent_folder_id FROM folders WHERE parent_folder_id IS NOT NULL ORDER BY id"); err != nil {
		return 0, fmt.Errorf("loading folders: %w", err)
	}

	// group by (parent_folder_id, computed basename)
	type groupKey struct {
		parentID int
		basename string
	}
	groups := make(map[groupKey][]folderRow)
	for _, f := range folders {
		key := groupKey{
			parentID: f.ParentFolderID,
			basename: filepath.Base(f.Path),
		}
		groups[key] = append(groups[key], f)
	}

	deduped := 0
	for _, group := range groups {
		if len(group) <= 1 {
			continue
		}

		if deduped == 0 {
			logger.Info("Deduplicating folders with conflicting basenames...")
		}

		// prefer the folder whose path is already normalized for the current OS,
		// falling back to the newest entry (highest ID) since it's most likely
		// from the current filesystem
		keep := group[len(group)-1]
		for _, f := range group {
			if f.Path == filepath.Clean(f.Path) {
				keep = f
				break
			}
		}

		for _, dup := range group {
			if dup.ID == keep.ID {
				continue
			}

			logger.Infof("Merging duplicate folder %d %q into folder %d %q", dup.ID, dup.Path, keep.ID, keep.Path)

			if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
				return m.mergeFolder(tx, keep.ID, dup.ID)
			}); err != nil {
				return 0, fmt.Errorf("merging folder %d into %d: %w", dup.ID, keep.ID, err)
			}

			deduped++
		}
	}

	if deduped > 0 {
		logger.Infof("Deduplicated %d folder entries", deduped)
	}

	return deduped, nil
}

func (m *schema84Migrator) mergeFolder(tx *sqlx.Tx, keepID, dupID int) error {
	// Re-parent child folders from the duplicate to the canonical folder.
	// At this point basenames are still full paths (unique), so this won't cause
	// UNIQUE constraint violations on (parent_folder_id, basename).
	if _, err := tx.Exec("UPDATE folders SET parent_folder_id = ? WHERE parent_folder_id = ?", keepID, dupID); err != nil {
		return fmt.Errorf("re-parenting child folders: %w", err)
	}

	// Orphan the stale duplicate folder by clearing its parent so the UNIQUE
	// constraint on (parent_folder_id, basename) won't be violated when
	// migrateFolders sets basenames. Any stale file entries under it are left
	// untouched — the clean task will handle them on the next scan.
	if _, err := tx.Exec("UPDATE folders SET parent_folder_id = NULL WHERE id = ?", dupID); err != nil {
		return fmt.Errorf("orphaning duplicate folder: %w", err)
	}

	return nil
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
