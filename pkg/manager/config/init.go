package config

import (
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
}

func Initialize() (*Instance, error) {
	var err error
	once.Do(func() {
		instance = &Instance{}

		flags := initFlags()
		err = initConfig(flags)
		initEnvs()
	})
	return instance, err
}

func initConfig(flags flagStruct) error {
	// The config file is called config.  Leave off the file extension.
	viper.SetConfigName("config")

	if flagConfigFileExists, _ := utils.FileExists(flags.configFilePath); flagConfigFileExists {
		viper.SetConfigFile(flags.configFilePath)
	}
	viper.AddConfigPath(".")            // Look for config in the working directory
	viper.AddConfigPath("$HOME/.stash") // Look for the config in the home directory

	// for Docker compatibility, if STASH_CONFIG_FILE is set, then touch the
	// given filename
	envConfigFile := os.Getenv("STASH_CONFIG_FILE")
	if envConfigFile != "" {
		utils.Touch(envConfigFile)
		viper.SetConfigFile(envConfigFile)
	}

	err := viper.ReadInConfig() // Find and read the config file
	// continue, but set an error to be handled by caller

	instance.SetInitialConfig()

	return err
}

func initFlags() flagStruct {
	flags := flagStruct{}

	pflag.IP("host", net.IPv4(0, 0, 0, 0), "ip address for the host")
	pflag.Int("port", 9999, "port to serve from")
	pflag.StringVarP(&flags.configFilePath, "config", "c", "", "config file to use")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logger.Infof("failed to bind flags: %s", err.Error())
	}

	return flags
}

func initEnvs() {
	viper.SetEnvPrefix("stash")    // will be uppercased automatically
	viper.BindEnv("host")          // STASH_HOST
	viper.BindEnv("port")          // STASH_PORT
	viper.BindEnv("external_host") // STASH_EXTERNAL_HOST
	viper.BindEnv("generated")     // STASH_GENERATED
	viper.BindEnv("metadata")      // STASH_METADATA
	viper.BindEnv("cache")         // STASH_CACHE

	// only set stash config flag if not already set
	if instance.GetStashPaths() == nil {
		viper.BindEnv("stash") // STASH_STASH
	}
}
