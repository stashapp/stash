package migrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema40Migrator struct {
	migrator
}

func post40(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 40")

	m := schema40Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.migrate(ctx); err != nil {
		return fmt.Errorf("migrating performer aliases: %w", err)
	}

	return nil
}

func (m *schema40Migrator) migrate(ctx context.Context) error {
	logger.Info("Migrating performer aliases")

	const (
		limit    = 1000
		logEvery = 10000
	)

	lastID := 0
	count := 0

	for {
		gotSome := false

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			query := "SELECT `id`, `aliases` FROM `performers` WHERE `aliases` IS NOT NULL AND `aliases` != ''"

			if lastID != 0 {
				query += fmt.Sprintf(" AND `id` > %d ", lastID)
			}

			query += fmt.Sprintf(" ORDER BY `id` LIMIT %d", limit)

			rows, err := m.db.Query(query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var (
					id      int
					aliases string
				)

				err := rows.Scan(&id, &aliases)
				if err != nil {
					return err
				}

				lastID = id
				gotSome = true
				count++

				if err := m.migratePerformerAliases(id, aliases); err != nil {
					return err
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
			logger.Infof("Migrated %d rows", count)
		}
	}

	// drop the aliases column
	if _, err := m.db.Exec("ALTER TABLE `performers` DROP COLUMN `aliases`"); err != nil {
		return err
	}

	return nil
}

func (m *schema40Migrator) migratePerformerAliases(id int, aliases string) error {
	// split aliases by , or /
	aliasList := strings.FieldsFunc(aliases, func(r rune) bool {
		return strings.ContainsRune(",/", r)
	})

	// trim whitespace from each alias
	for i, alias := range aliasList {
		aliasList[i] = strings.TrimSpace(alias)
	}

	// remove duplicates
	aliasList = stringslice.StrAppendUniques(nil, aliasList)

	// insert aliases into table
	for _, alias := range aliasList {
		_, err := m.db.Exec("INSERT INTO `performer_aliases` (`performer_id`, `alias`) VALUES (?, ?)", id, alias)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(40, post40)
}
