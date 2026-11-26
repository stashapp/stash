package migrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema88Migrator struct {
	migrator
}

func post88(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 88")

	m := schema88Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.migrateGroupAliases(ctx); err != nil {
		return fmt.Errorf("migrating group aliases: %w", err)
	}

	return nil
}

type groupAliasRow struct {
	GroupID int    `db:"group_id"`
	Alias   string `db:"alias"`
}

func (m *schema88Migrator) migrateGroupAliases(ctx context.Context) error {
	logger.Info("Migrating group aliases")

	const limit = 1000

	lastID := 0
	count := 0

	for {
		var rows []groupAliasRow

		if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
			// each group has at most one row here, copied from the old aliases column
			if err := tx.Select(&rows, "SELECT `group_id`, `alias` FROM `group_aliases` WHERE `group_id` > ? ORDER BY `group_id` LIMIT ?", lastID, limit); err != nil {
				return err
			}

			for _, r := range rows {
				if err := m.migrateGroupAlias(tx, r.GroupID, r.Alias); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}

		if len(rows) == 0 {
			break
		}

		lastID = rows[len(rows)-1].GroupID
		count += len(rows)
	}

	logger.Infof("Migrated aliases for %d groups", count)

	return nil
}

func (m *schema88Migrator) migrateGroupAlias(tx *sqlx.Tx, groupID int, aliases string) error {
	aliasList := splitGroupAliases(aliases)

	if len(aliasList) < 2 {
		// existing value is fine
		return nil
	}

	if _, err := tx.Exec("DELETE FROM `group_aliases` WHERE `group_id` = ?", groupID); err != nil {
		return err
	}

	for _, alias := range aliasList {
		if _, err := tx.Exec("INSERT INTO `group_aliases` (`group_id`, `alias`) VALUES (?, ?)", groupID, alias); err != nil {
			return err
		}
	}

	return nil
}

// splitGroupAliases splits a comma-delimited aliases string into trimmed,
// unique, non-empty aliases.
func splitGroupAliases(aliases string) []string {
	var ret []string
	for _, alias := range strings.Split(aliases, ",") {
		if alias = strings.TrimSpace(alias); alias != "" {
			ret = append(ret, alias)
		}
	}

	return sliceutil.AppendUniques(nil, ret)
}

func init() {
	sqlite.RegisterPostMigration(88, post88)
}
