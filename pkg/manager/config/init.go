package config

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

var once sync.Once

type flagStruct struct {
	configFilePath string
	cpuProfilePath string
}

func Initialize() (*Instance, error) {
	var err error
	once.Do(func() {
		flags := initFlags()
		instance = &Instance{
			cpuProfilePath: flags.cpuProfilePath,
		}

		if err = initConfig(flags); err != nil {
			return
		}
		initEnvs()

		if instance.isNewSystem {
			if instance.Validate() == nil {
				// system has been initialised by the environment
				instance.isNewSystem = false
			}
		}

		if !instance.isNewSystem {
			err = instance.SetInitialConfig()
		}
	})
	return instance, err
}

func initConfig(flags flagStruct) error {

	// The config file is called config.  Leave off the file extension.
	viper.SetConfigName("config")

	viper.AddConfigPath(".")            // Look for config in the working directory
	viper.AddConfigPath("$HOME/.stash") // Look for the config in the home directory

	configFile := ""
	envConfigFile := os.Getenv("STASH_CONFIG_FILE")

	if flags.configFilePath != "" {
		configFile = flags.configFilePath
	} else if envConfigFile != "" {
		configFile = envConfigFile
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)

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

	err := viper.ReadInConfig() // Find and read the config file
	// if not found, assume its a new system
	if _, isMissing := err.(viper.ConfigFileNotFoundError); isMissing {
		instance.isNewSystem = true
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func initFlags() flagStruct {
	flags := flagStruct{}

	pflag.IP("host", net.IPv4(0, 0, 0, 0), "ip address for the host")
	pflag.Int("port", 9999, "port to serve from")
	pflag.StringVarP(&flags.configFilePath, "config", "c", "", "config file to use")
	pflag.StringVar(&flags.cpuProfilePath, "cpuprofile", "", "write cpu profile to file")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logger.Infof("failed to bind flags: %s", err.Error())
	}

	return flags
}

func initEnvs() {
	viper.SetEnvPrefix("stash") // will be uppercased automatically
	bindEnv("host")             // STASH_HOST
	bindEnv("port")             // STASH_PORT
	bindEnv("external_host")    // STASH_EXTERNAL_HOST
	bindEnv("generated")        // STASH_GENERATED
	bindEnv("metadata")         // STASH_METADATA
	bindEnv("cache")            // STASH_CACHE

	// only set stash config flag if not already set
	if instance.GetStashPaths() == nil {
		bindEnv("stash") // STASH_STASH
	}
}

func bindEnv(key string) {
	if err := viper.BindEnv(key); err != nil {
		panic(fmt.Sprintf("unable to set environment key (%v): %v", key, err))
	}
}
