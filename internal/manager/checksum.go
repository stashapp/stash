package manager

import (
	"context"
	"errors"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

type SceneCounter interface {
	Count(ctx context.Context) (int, error)
}

func setInitialMD5Config(ctx context.Context, txnManager txn.Manager, counter SceneCounter) {
	// if there are no scene files in the database, then default the
	// VideoFileNamingAlgorithm config setting to oshash and calculateMD5 to
	// false, otherwise set them to true for backwards compatibility purposes
	var count int
	if err := txn.WithTxn(ctx, txnManager, func(ctx context.Context) error {
		var err error
		count, err = counter.Count(ctx)
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

	config := config.GetInstance()
	config.SetChecksumDefaultValues(defaultAlgorithm, usingMD5)
	if err := config.Write(); err != nil {
		logger.Errorf("Error while writing configuration file: %s", err.Error())
	}
}

type SceneMissingHashCounter interface {
	CountMissingChecksum(ctx context.Context) (int, error)
	CountMissingOSHash(ctx context.Context) (int, error)
}

// ValidateVideoFileNamingAlgorithm validates changing the
// VideoFileNamingAlgorithm configuration flag.
//
// If setting VideoFileNamingAlgorithm to MD5, then this function will ensure
// that all checksum values are set on all scenes.
//
// Likewise, if VideoFileNamingAlgorithm is set to oshash, then this function
// will ensure that all oshash values are set on all scenes.
func ValidateVideoFileNamingAlgorithm(ctx context.Context, qb SceneMissingHashCounter, newValue models.HashAlgorithm) error {
	// if algorithm is being set to MD5, then all checksums must be present
	if newValue == models.HashAlgorithmMd5 {
		missingMD5, err := qb.CountMissingChecksum(ctx)
		if err != nil {
			return err
		}

		if missingMD5 > 0 {
			return errors.New("some checksums are missing on scenes. Run Scan with calculateMD5 set to true")
		}
	} else if newValue == models.HashAlgorithmOshash {
		missingOSHash, err := qb.CountMissingOSHash(ctx)
		if err != nil {
			return err
		}

		if missingOSHash > 0 {
			return errors.New("some oshash values are missing on scenes. Run Scan to populate")
		}
	}

	return nil
}
