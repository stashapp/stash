package manager

import (
	"context"
	"errors"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func setInitialMD5Config(ctx context.Context, r models.Repository) {
	// if there are no scene files in the database, then default the
	// VideoFileNamingAlgorithm config setting to oshash and calculateMD5 to
	// false, otherwise set them to true for backwards compatibility purposes
	var count int
	if err := r.WithTxn(ctx, func(ctx context.Context) error {
		var err error
		count, err = r.Scene.Count(ctx)
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

// ValidateVideoFileNamingAlgorithm validates changing the
// VideoFileNamingAlgorithm configuration flag.
//
// If setting VideoFileNamingAlgorithm to MD5, then this function will ensure
// that all checksum values are set on all scenes.
//
// Likewise, if VideoFileNamingAlgorithm is set to oshash, then this function
// will ensure that all oshash values are set on all scenes.
func ValidateVideoFileNamingAlgorithm(r models.Repository, newValue models.HashAlgorithm) error {
	// if algorithm is being set to MD5, then all checksums must be present
	return r.WithTxn(context.TODO(), func(ctx context.Context) error {
		qb := r.Scene
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
	})
}
