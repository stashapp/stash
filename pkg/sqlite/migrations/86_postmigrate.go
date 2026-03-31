package migrations

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/user"
	"gopkg.in/guregu/null.v4/zero"
)

func post86(ctx context.Context, db *sqlx.DB) error {
	// don't run this if we're in the setup context - it creates credentials itself
	if models.IsSetupContext(ctx) {
		logger.Debug("Skipping post-migration for schema version 86 because we're in the setup context")
		return nil
	}

	logger.Info("Running post-migration for schema version 86")

	m := schema86Migrator{
		migrator: migrator{
			db: db,
		},
	}

	if err := m.createInitialUser(ctx); err != nil {
		return fmt.Errorf("creating initial user: %w", err)
	}

	return nil
}

type schema86Migrator struct {
	migrator
}

func (m *schema86Migrator) createInitialUser(ctx context.Context) error {
	// if credentials are set in the config file, we need to create an initial user with those credentials, so that the user can log in after the migration
	cfg := config.GetInstance()

	username := cfg.GetLegacyUsername()
	password := cfg.GetLegacyPasswordHash()
	apiKey := cfg.GetLegacyAPIKey()

	if username == "" || password == "" {
		// create a default user with no password
		logger.Warn("No credentials found in config file. Creating default user with no password.")
		username = "admin"
		password = ""
	} else {
		logger.Info("Credentials found in config file. Creating initial user with those credentials.")
		encodedPassword, err := user.EncodeLegacyHash(password)
		if err != nil {
			return fmt.Errorf("error encoding legacy password: %w", err)
		}
		password = encodedPassword

		if apiKey != "" {
			apiKey, err = user.HashAPIKey(apiKey)
			if err != nil {
				return fmt.Errorf("error hashing legacy API key: %w", err)
			}
		}
	}

	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		const insertSQL = "INSERT INTO `users` (`id`,`username`,`password_hash`,`api_key`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)"

		const id = 1
		passwordHashSQL := zero.StringFrom(password)
		apiKeySQL := zero.StringFrom(apiKey)
		now := sqlite.Timestamp{Timestamp: time.Now()}

		_, err := tx.Exec(insertSQL, id, username, passwordHashSQL, apiKeySQL, now, now)
		if err != nil {
			return fmt.Errorf("error inserting initial user: %w", err)
		}

		// add admin role
		const insertRoleSQL = "INSERT INTO `user_roles` (`user_id`,`role`) VALUES (?,?)"
		_, err = tx.Exec(insertRoleSQL, id, models.RoleEnumAdmin)
		if err != nil {
			return fmt.Errorf("error inserting initial user role: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error creating initial user: %w", err)
	}

	// remove the legacy credentials from the config file, since they're now stored in the database
	if err := m.removeLegacyCredentials(); err != nil {
		return fmt.Errorf("error removing legacy credentials: %w", err)
	}

	return nil
}

func (m *schema86Migrator) removeLegacyCredentials() error {
	c := config.GetInstance()

	orgPath := c.GetConfigFile()

	hasCredentials := c.GetLegacyUsername() != "" || c.GetLegacyPasswordHash() != "" || c.GetLegacyAPIKey() != ""

	if orgPath == "" || !hasCredentials {
		return nil
	}

	logger.Info("Removing legacy credentials from config file")

	// save a backup of the original config file
	backupPath := fmt.Sprintf("%s.85.%s", orgPath, time.Now().Format("20060102_150405"))

	data, err := c.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal backup config: %w", err)
	}

	logger.Infof("Backing up config to %s", backupPath)
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup config: %w", err)
	}

	c.SetInterface(config.LegacyUsername, nil)
	c.SetInterface(config.LegacyPassword, nil)
	c.SetInterface(config.LegacyAPIKey, nil)

	if err := c.Write(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(86, post86)
}
