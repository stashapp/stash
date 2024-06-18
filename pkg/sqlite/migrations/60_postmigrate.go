package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
)

type schema60Migrator struct {
	migrator
}

func post60(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 60")

	m := schema60Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate(ctx)
}

func (m *schema60Migrator) decodeJSON(s string, v interface{}) {
	if s == "" {
		return
	}

	if err := json.Unmarshal([]byte(s), v); err != nil {
		logger.Errorf("error decoding json %q: %v", s, err)
	}
}

type schema60DefaultFilters map[string]interface{}

func (m *schema60Migrator) migrate(ctx context.Context) error {

	// save default filters into the UI config
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		query := "SELECT id, mode, find_filter, object_filter, ui_options FROM `saved_filters` WHERE `name` = ''"

		rows, err := m.db.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		defaultFilters := make(schema60DefaultFilters)

		for rows.Next() {
			var (
				id              int
				mode            string
				findFilterStr   string
				objectFilterStr string
				uiOptionsStr    string
			)

			if err := rows.Scan(&id, &mode, &findFilterStr, &objectFilterStr, &uiOptionsStr); err != nil {
				return err
			}

			// convert the filters to the correct format
			findFilter := make(map[string]interface{})
			objectFilter := make(map[string]interface{})
			uiOptions := make(map[string]interface{})

			m.decodeJSON(findFilterStr, &findFilter)
			m.decodeJSON(objectFilterStr, &objectFilter)
			m.decodeJSON(uiOptionsStr, &uiOptions)

			o := map[string]interface{}{
				"mode":          mode,
				"find_filter":   findFilter,
				"object_filter": objectFilter,
				"ui_options":    uiOptions,
			}

			defaultFilters[strings.ToLower(mode)] = o
		}

		if err := rows.Err(); err != nil {
			return err
		}

		if err := m.saveDefaultFilters(defaultFilters); err != nil {
			return fmt.Errorf("saving default filters: %w", err)
		}

		// remove the default filters from the database
		query = "DELETE FROM `saved_filters` WHERE `name` = ''"
		if _, err := m.db.Exec(query); err != nil {
			return fmt.Errorf("deleting default filters: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (m *schema60Migrator) saveDefaultFilters(defaultFilters schema60DefaultFilters) error {
	if len(defaultFilters) == 0 {
		logger.Debugf("no default filters to save")
		return nil
	}

	// save the default filters into the UI config
	config := config.GetInstance()

	orgPath := config.GetConfigFile()

	if orgPath == "" {
		// no config file to migrate (usually in a test or new system)
		logger.Debugf("no config file to migrate")
		return nil
	}

	uiConfig := config.GetUIConfiguration()
	if uiConfig == nil {
		uiConfig = make(map[string]interface{})
	}

	// if the defaultFilters key already exists, don't overwrite them
	if _, found := uiConfig["defaultFilters"]; found {
		logger.Warn("defaultFilters already exists in the UI config, skipping migration")
		return nil
	}

	if err := m.backupConfig(orgPath); err != nil {
		return fmt.Errorf("backing up config: %w", err)
	}

	uiConfig["defaultFilters"] = map[string]interface{}(defaultFilters)
	config.SetUIConfiguration(uiConfig)

	if err := config.Write(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (m *schema60Migrator) backupConfig(orgPath string) error {
	c := config.GetInstance()

	// save a backup of the original config file
	backupPath := fmt.Sprintf("%s.59.%s", orgPath, time.Now().Format("20060102_150405"))

	data, err := c.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal backup config: %w", err)
	}

	logger.Infof("Backing up config to %s", backupPath)
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup config: %w", err)
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(60, post60)
}
