package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/stashapp/stash/internal/build"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

var (
	initOnce     sync.Once
	instanceOnce sync.Once
)

type flagStruct struct {
	configFilePath string
	cpuProfilePath string
	nobrowser      bool
	helpFlag       bool
	versionFlag    bool
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

func Initialize() (*Instance, error) {
	var err error
	initOnce.Do(func() {
		flags := initFlags()

		if flags.helpFlag {
			pflag.Usage()
			os.Exit(0)
		}

		if flags.versionFlag {
			fmt.Printf(build.VersionString() + "\n")
			os.Exit(0)
		}

		overrides := makeOverrideConfig()

		_ = GetInstance()
		instance.overrides = overrides
		instance.cpuProfilePath = flags.cpuProfilePath
		// instance.configUpdates = make(chan int)

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

func initConfig(instance *Instance, flags flagStruct) error {
	v := instance.main

	// The config file is called config.  Leave off the file extension.
	v.SetConfigName("config")

	v.AddConfigPath(".")                                // Look for config in the working directory
	v.AddConfigPath(filepath.FromSlash("$HOME/.stash")) // Look for the config in the home directory

	configFile := ""
	envConfigFile := os.Getenv("STASH_CONFIG_FILE")

	if flags.configFilePath != "" {
		configFile = flags.configFilePath
	} else if envConfigFile != "" {
		configFile = envConfigFile
	}

	if configFile != "" {
		v.SetConfigFile(configFile)

		// if file does not exist, assume it is a new system
		if exists, _ := fsutil.FileExists(configFile); !exists {
			instance.isNewSystem = true

			// ensure we can write to the file
			if err := fsutil.Touch(configFile); err != nil {
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

func initFlags() flagStruct {
	flags := flagStruct{}

	pflag.IP("host", net.IPv4(0, 0, 0, 0), "ip address for the host")
	pflag.Int("port", 9999, "port to serve from")
	pflag.StringVarP(&flags.configFilePath, "config", "c", "", "config file to use")
	pflag.StringVar(&flags.cpuProfilePath, "cpuprofile", "", "write cpu profile to file")
	pflag.BoolVar(&flags.nobrowser, "nobrowser", false, "Don't open a browser window after launch")
	pflag.BoolVarP(&flags.helpFlag, "help", "h", false, "show this help text and exit")
	pflag.BoolVarP(&flags.versionFlag, "version", "v", false, "show version number and exit")

	pflag.Parse()

	return flags
}

func initEnvs(viper *viper.Viper) {
	viper.SetEnvPrefix("stash")     // will be uppercased automatically
	bindEnv(viper, "host")          // STASH_HOST
	bindEnv(viper, "port")          // STASH_PORT
	bindEnv(viper, "external_host") // STASH_EXTERNAL_HOST
	bindEnv(viper, "generated")     // STASH_GENERATED
	bindEnv(viper, "metadata")      // STASH_METADATA
	bindEnv(viper, "cache")         // STASH_CACHE
	bindEnv(viper, "stash")         // STASH_STASH
}

func bindEnv(viper *viper.Viper, key string) {
	if err := viper.BindEnv(key); err != nil {
		panic(fmt.Sprintf("unable to set environment key (%v): %v", key, err))
	}
}

func makeOverrideConfig() *viper.Viper {
	viper := viper.New()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logger.Infof("failed to bind flags: %s", err.Error())
	}

	initEnvs(viper)

	return viper
}
