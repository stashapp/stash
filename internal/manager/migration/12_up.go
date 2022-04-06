package migration

import (
	"context"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func (m *PostMigrator) migrate12(ctx context.Context) {
	// if there are no scene files in the database, then default the
	// VideoFileNamingAlgorithm config setting to oshash and calculateMD5 to
	// false, otherwise set them to true for backwards compatibility purposes
	var count int
	if err := m.TxnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		var err error
		count, err = r.Scene().Count()
		return err
	}); err != nil {
		logger.Errorf("Error while counting scenes: %s", err.Error())
		return
	}

	usingMD5 := count != 0
	defaultAlgorithm := models.HashAlgorithmOshash
	if usingMD5 {
		defaultAlgorithm = models.HashAlgorithmMd5
	}

	m.Config.Set(config.VideoFileNamingAlgorithm, defaultAlgorithm)
	m.Config.Set(config.CalculateMD5, usingMD5)
	if err := m.Config.Write(); err != nil {
		logger.Errorf("Error while writing configuration file: %s", err.Error())
	}
}
