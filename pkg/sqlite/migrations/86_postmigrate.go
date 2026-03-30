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

type schema86Migrator struct {
	migrator
}

func post86(ctx context.Context, db *sqlx.DB) error {
	logger.Info("Running post-migration for schema version 86")

	m := schema86Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrate(ctx)
}

func (m *schema86Migrator) migrate(ctx context.Context) error {
	return m.migrateFingerprintQueue(ctx)
}

func (m *schema86Migrator) migrateFingerprintQueue(ctx context.Context) error {
	c := config.GetInstance()

	orgPath := c.GetConfigFile()
	if orgPath == "" {
		// no config file to migrate (usually in a test or new system)
		logger.Debugf("no config file to migrate")
		return nil
	}

	uiConfig := c.GetUIConfiguration()
	if uiConfig == nil {
		logger.Debugf("no UI config to migrate")
		return nil
	}

	taggerConfig, ok := uiConfig["taggerConfig"].(map[string]any)
	if !ok {
		logger.Debugf("no taggerConfig in UI config")
		return nil
	}

	fingerprintQueue, ok := taggerConfig["fingerprintQueue"].(map[string]any)
	if !ok {
		logger.Debugf("no fingerprintQueue in taggerConfig")
		return nil
	}

	if len(fingerprintQueue) == 0 {
		logger.Debugf("fingerprintQueue is empty")
		return nil
	}

	// Backup config before modifying
	if err := m.backupConfig(orgPath); err != nil {
		return fmt.Errorf("backing up config: %w", err)
	}

	// Migrate each endpoint's queue to the database
	// Legacy format: fingerprintQueue[endpoint] = ["sceneId1", "sceneId2", ...]
	// We need to look up the stash-box scene ID from scene_stash_ids table
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		for endpoint, queueData := range fingerprintQueue {
			queue, ok := queueData.([]any)
			if !ok {
				logger.Warnf("fingerprintQueue[%s] is not an array, skipping", endpoint)
				continue
			}

			for _, entryData := range queue {
				// Legacy format: entries are just scene ID strings
				sceneID, ok := entryData.(string)
				if !ok {
					// Try parsing as float64 (JSON numbers)
					if f, ok := entryData.(float64); ok {
						sceneID = fmt.Sprintf("%d", int(f))
					} else {
						logger.Warnf("fingerprintQueue entry is not a string or number, skipping: %T", entryData)
						continue
					}
				}

				if sceneID == "" {
					logger.Warnf("fingerprintQueue entry is empty, skipping")
					continue
				}

				// Look up the stash-box scene ID from scene_stash_ids
				var stashBoxSceneID string
				err := tx.QueryRow(`
					SELECT stash_id FROM scene_stash_ids
					WHERE scene_id = ? AND endpoint = ?
				`, sceneID, endpoint).Scan(&stashBoxSceneID)
				if err != nil {
					logger.Warnf("Could not find stash_id for scene %s endpoint %s, skipping: %v", sceneID, endpoint, err)
					continue
				}

				// Insert into the new table, ignore conflicts (entry already exists)
				_, err = tx.Exec(`
					INSERT OR IGNORE INTO fingerprint_submissions (endpoint, stash_id, scene_id, vote)
					VALUES (?, ?, ?, ?)
				`, endpoint, stashBoxSceneID, sceneID, "VALID")
				if err != nil {
					return fmt.Errorf("inserting fingerprint submission: %w", err)
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// Remove fingerprintQueue from taggerConfig
	delete(taggerConfig, "fingerprintQueue")
	uiConfig["taggerConfig"] = taggerConfig
	c.SetUIConfiguration(uiConfig)

	if err := c.Write(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	logger.Info("Migrated fingerprintQueue to database")
	return nil
}

func (m *schema86Migrator) backupConfig(orgPath string) error {
	c := config.GetInstance()

	backupPath := fmt.Sprintf("%s.85.%s", orgPath, time.Now().Format("20060102_150405"))

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
	sqlite.RegisterPostMigration(86, post86)
}
