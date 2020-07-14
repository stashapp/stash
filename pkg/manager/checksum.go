package manager

import (
	"errors"

	"github.com/spf13/viper"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

func setInitialMD5Config() {
	// if there are no scene files in the database, then default the useMD5
	// and calculateMD5 config settings to false, otherwise set them to true
	// for backwards compatibility purposes
	sqb := models.NewSceneQueryBuilder()
	count, err := sqb.Count()
	if err != nil {
		logger.Errorf("Error while counting scenes: %s", err.Error())
		return
	}

	usingMD5 := count != 0

	viper.SetDefault(config.UseMD5, usingMD5)
	viper.SetDefault(config.CalculateMD5, usingMD5)

	if err := config.Write(); err != nil {
		logger.Errorf("Error while writing configuration file: %s", err.Error())
	}
}

// ValidateUseMD5 validates changing the UseMD5 configuration flag.
//
// If setting UseMD5 to true, then this function will ensure that all checksum
// values are set on all scenes.
//
// Likewise, if UseMD5 is set to false, then this function will ensure that all
// oshash values are set on all scenes.
func ValidateUseMD5(newValue bool) error {
	// if useMD5 is being set to true, then all checksums must be present
	qb := models.NewSceneQueryBuilder()
	if newValue {
		missingMD5, err := qb.CountMissingChecksum()
		if err != nil {
			return err
		}

		if missingMD5 > 0 {
			return errors.New("some checksums are missing on scenes. Run Scan with calculateMD5 set to true")
		}
	} else {
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
