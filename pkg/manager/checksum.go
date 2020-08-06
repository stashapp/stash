package manager

import (
	"errors"

	"github.com/spf13/viper"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

func setInitialMD5Config() {
	// if there are no scene files in the database, then default the
	// VideoFileNamingAlgorithm config setting to oshash and calculateMD5 to
	// false, otherwise set them to true for backwards compatibility purposes
	sqb := models.NewSceneQueryBuilder()
	count, err := sqb.Count()
	if err != nil {
		logger.Errorf("Error while counting scenes: %s", err.Error())
		return
	}

	usingMD5 := count != 0
	defaultAlgorithm := models.HashAlgorithmOshash

	if usingMD5 {
		defaultAlgorithm = models.HashAlgorithmMd5
	}

	viper.SetDefault(config.VideoFileNamingAlgorithm, defaultAlgorithm)
	viper.SetDefault(config.CalculateMD5, usingMD5)

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
func ValidateVideoFileNamingAlgorithm(newValue models.HashAlgorithm) error {
	// if algorithm is being set to MD5, then all checksums must be present
	qb := models.NewSceneQueryBuilder()
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
}
