package manager

import (
	"context"
	"errors"

	"github.com/stashapp/stash/pkg/models"
)

// ValidateVideoFileNamingAlgorithm validates changing the
// VideoFileNamingAlgorithm configuration flag.
//
// If setting VideoFileNamingAlgorithm to MD5, then this function will ensure
// that all checksum values are set on all scenes.
//
// Likewise, if VideoFileNamingAlgorithm is set to oshash, then this function
// will ensure that all oshash values are set on all scenes.
func ValidateVideoFileNamingAlgorithm(txnManager models.TransactionManager, newValue models.HashAlgorithm) error {
	// if algorithm is being set to MD5, then all checksums must be present
	return txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		qb := r.Scene()
		if newValue == models.HashAlgorithmMd5 {
			missingMD5, err := qb.CountMissingChecksum()
			if err != nil {
				return err
			}

			if missingMD5 > 0 {
				return errors.New("some checksums are missing on scenes. Run Scan with calculateMD5 set to true")
			}
		} else if newValue == models.HashAlgorithmOshash {
			missingOSHash, err := qb.CountMissingOSHash()
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
