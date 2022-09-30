package migrations

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
)

func post12(ctx context.Context, db *sqlx.DB) error {
	m := schema12Migrator{
		migrator: migrator{
			db: db,
		},
	}

	return m.migrateConfig(ctx)
}

type schema12Migrator struct {
	migrator
}

func (m *schema12Migrator) migrateConfig(ctx context.Context) error {
	// if there are no scene files in the database, then default the
	// VideoFileNamingAlgorithm config setting to oshash and calculateMD5 to
	// false, otherwise set them to true for backwards compatibility purposes
	var count int
	if err := m.withTxn(ctx, func(tx *sqlx.Tx) error {
		query := "SELECT COUNT(*) from `scenes`"

		return tx.Get(&count, query)
	}); err != nil {
		return err
	}

	usingMD5 := count != 0
	defaultAlgorithm := models.HashAlgorithmOshash
	if usingMD5 {
		logger.Infof("Defaulting video file naming algorithm to %s", models.HashAlgorithmMd5)
		defaultAlgorithm = models.HashAlgorithmMd5
	}

	c := config.GetInstance()

	c.SetDefault(config.VideoFileNamingAlgorithm, defaultAlgorithm)
	c.SetDefault(config.CalculateMD5, usingMD5)
	if err := c.Write(); err != nil {
		logger.Errorf("Error while writing configuration file: %s", err.Error())
	}

	return nil
}

func init() {
	sqlite.RegisterPostMigration(12, post12)
}
