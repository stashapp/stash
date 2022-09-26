package manager

import (
	"context"
	"errors"

	"github.com/stashapp/stash/pkg/models"
)

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
