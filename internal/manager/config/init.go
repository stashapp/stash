package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/pflag"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type flagStruct struct {
	configFilePath string
	nobrowser      bool
}

var (
	flags flagStruct

	homeDir, _ = os.UserHomeDir()

	defaultConfigLocations = []string{
		"config.yml",
		filepath.Join(homeDir, ".stash", "config.yml"),
	}

	// map of env vars to config keys
	envBinds = map[string]string{
		"host":          Host,
		"port":          Port,
		"external_host": ExternalHost,
		"generated":     Generated,
		"metadata":      Metadata,
		"cache":         Cache,
		"stash":         Stash,
		"ui":            UILocation,
	}
)

var errConfigNotFound = errors.New("config file not found")

func init() {
	pflag.IP("host", net.IPv4(0, 0, 0, 0), "ip address for the host")
	pflag.Int("port", 9999, "port to serve from")
	pflag.StringVarP(&flags.configFilePath, "config", "c", "", "config file to use")
	pflag.BoolVar(&flags.nobrowser, "nobrowser", false, "Don't open a browser window after launch")
	pflag.StringP("ui-location", "u", "", "path to the webui")
}

// Called at startup
func Initialize() (*Config, error) {
	cfg := &Config{
		main:      koanf.New("."),
		overrides: koanf.New("."),
	}

	cfg.initOverrides()

	err := cfg.initConfig()
	if err != nil {
		return nil, err
	}

	if cfg.isNewSystem {
		if cfg.Validate() == nil {
			// system has been initialised by the environment
			cfg.isNewSystem = false
		}
	}

	if !cfg.isNewSystem {
		cfg.setExistingSystemDefaults()

		err := cfg.SetInitialConfig()
		if err != nil {
			return nil, err
		}

		err = cfg.Write()
		if err != nil {
			return nil, err
		}

		err = cfg.Validate()
		if err != nil {
			return nil, err
		}
	}

	instance = cfg
	return instance, nil
}

// Called by tests to initialize an empty config
func InitializeEmpty() *Config {
	cfg := &Config{
		main:      koanf.New("."),
		overrides: koanf.New("."),
	}
	instance = cfg
	return instance
}

func (i *Config) loadFromCommandLine() {
	v := i.overrides

	if err := v.Load(posflag.ProviderWithFlag(pflag.CommandLine, ".", v, func(f *pflag.Flag) (string, interface{}) {
		// ignore flags that have not been changed
		if !f.Changed {
			return "", nil
		}

		return f.Name, posflag.FlagVal(pflag.CommandLine, f)
	}), nil); err != nil {
		logger.Errorf("failed to load flags: %v", err)
	}
}

func (i *Config) loadFromEnv() {
	v := i.overrides

	if err := v.Load(env.ProviderWithValue("STASH_", ".", func(key, value string) (string, interface{}) {
		key = strings.ToLower(strings.TrimPrefix(key, "STASH_"))
		if newKey, ok := envBinds[key]; ok {
			return newKey, value
		}

		return "", nil
	}), nil); err != nil {
		logger.Errorf("failed to load envs: %v", err)
	}
}

func (i *Config) initOverrides() {
	i.loadFromCommandLine()
	i.loadFromEnv()
}

func (i *Config) initConfig() error {
	configFile := ""
	envConfigFile := os.Getenv("STASH_CONFIG_FILE")

	if flags.configFilePath != "" {
		configFile = flags.configFilePath
	} else if envConfigFile != "" {
		configFile = envConfigFile
	}

	if configFile != "" {
		// if file does not exist, assume it is a new system
		if exists, _ := fsutil.FileExists(configFile); !exists {
			i.isNewSystem = true
			i.SetConfigFile(configFile)

			// ensure we can write to the file
			if err := fsutil.Touch(configFile); err != nil {
				return fmt.Errorf(`could not write to provided config path "%s": %v`, configFile, err)
			} else {
				// remove the file
				os.Remove(configFile)
			}

			return nil
		} else {
			// load from provided config file
			if err := i.loadFirstFromFiles([]string{configFile}); err != nil {
				return err
			}
		}
	} else {
		// load from default locations
		if err := i.loadFirstFromFiles(defaultConfigLocations); err != nil {
			if errors.Is(err, errConfigNotFound) {
				i.isNewSystem = true
				return nil
			}

			return err
		}
	}

	return nil
}

func (i *Config) loadFirstFromFiles(f []string) error {
	for _, ff := range f {
		if exists, _ := fsutil.FileExists(ff); exists {
			return i.load(ff)
		}
	}

	return errConfigNotFound
}
