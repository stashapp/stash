package migrations

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema65Migrator struct {
	migrator
}

func post65(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 65")

	m := schema65Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate()
}

func (m *schema65Migrator) migrate() error {
	if err := m.migrateConfig(); err != nil {
		return fmt.Errorf("failed to migrate config: %w", err)
	}

	return nil
}

func (m *schema65Migrator) migrateConfig() error {
	c := config.GetInstance()

	orgPath := c.GetConfigFile()

	if orgPath == "" {
		// no config file to migrate (usually in a test)
		return nil
	}

	items := c.GetMenuItems()
	replaced := false

	// replace "movies" with "groups" in the menu items
	for i, item := range items {
		if item == "movies" {
			items[i] = "groups"
			replaced = true
		}
	}

	if !replaced {
		return nil
	}

	// save a backup of the original config file
	backupPath := fmt.Sprintf("%s.64.%s", orgPath, time.Now().Format("20060102_150405"))

	data, err := c.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal backup config: %w", err)
	}

	logger.Infof("Backing up config to %s", backupPath)
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup config: %w", err)
	}

	c.SetInterface(config.MenuItems, items)

	if err := c.Write(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(65, post65)
}
