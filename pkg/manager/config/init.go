package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"

	"github.com/stashapp/stash/pkg/utils"
)

var (
	initOnce     sync.Once
	instanceOnce sync.Once
)

type FlagStruct struct {
	ConfigFilePath string
	CpuProfilePath string
	NoBrowser      bool
}

func GetInstance() *Instance {
	instanceOnce.Do(func() {
		instance = &Instance{
			main:      viper.New(),
			overrides: viper.New(),
		}
	})
	return instance
}

func Initialize(flags FlagStruct, overrides *viper.Viper) (*Instance, error) {
	var err error
	initOnce.Do(func() {
		_ = GetInstance()
		instance.overrides = overrides
		instance.cpuProfilePath = flags.CpuProfilePath

		if err = initConfig(instance, flags); err != nil {
			return
		}

		if instance.isNewSystem {
			if instance.Validate() == nil {
				// system has been initialised by the environment
				instance.isNewSystem = false
			}
		}

		if !instance.isNewSystem {
			err = instance.setExistingSystemDefaults()
			if err == nil {
				err = instance.SetInitialConfig()
			}
		}
	})
	return instance, err
}

func initConfig(instance *Instance, flags FlagStruct) error {
	v := instance.main

	// The config file is called config.  Leave off the file extension.
	v.SetConfigName("config")

	v.AddConfigPath(".")                                // Look for config in the working directory
	v.AddConfigPath(filepath.FromSlash("$HOME/.stash")) // Look for the config in the home directory

	configFile := ""
	envConfigFile := os.Getenv("STASH_CONFIG_FILE")

	if flags.ConfigFilePath != "" {
		configFile = flags.ConfigFilePath
	} else if envConfigFile != "" {
		configFile = envConfigFile
	}

	if configFile != "" {
		v.SetConfigFile(configFile)

		// if file does not exist, assume it is a new system
		if exists, _ := utils.FileExists(configFile); !exists {
			instance.isNewSystem = true

			// ensure we can write to the file
			if err := utils.Touch(configFile); err != nil {
				return fmt.Errorf(`could not write to provided config path "%s": %s`, configFile, err.Error())
			} else {
				// remove the file
				os.Remove(configFile)
			}

			return nil
		}
	}

	err := v.ReadInConfig() // Find and read the config file
	// if not found, assume its a new system
	var notFoundErr viper.ConfigFileNotFoundError
	if errors.As(err, &notFoundErr) {
		instance.isNewSystem = true
		return nil
	} else if err != nil {
		return err
	}

	return nil
}
