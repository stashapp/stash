package manager

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/utils"
	"net"
	"sync"
)

type singleton struct {
	Status JobStatus
	Paths  *paths.Paths
	JSON   *jsonUtils

	FFMPEGPath  string
	FFProbePath string
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	Initialize()
	return instance
}

func Initialize() *singleton {
	once.Do(func() {
		_ = utils.EnsureDir(paths.GetConfigDirectory())
		initConfig()
		initFlags()
		initEnvs()
		instance = &singleton{
			Status: Idle,
			Paths:  paths.NewPaths(),
			JSON:   &jsonUtils{},
		}

		instance.refreshConfig()

		initFFMPEG()
	})

	return instance
}

func initConfig() {
	// The config file is called config.  Leave off the file extension.
	viper.SetConfigName("config")

	viper.AddConfigPath("$HOME/.stash") // Look for the config in the home directory
	viper.AddConfigPath(".")            // Look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		_ = utils.Touch(paths.GetDefaultConfigFilePath())
		if err = viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	viper.SetDefault(config.Database, paths.GetDefaultDatabaseFilePath())

	// Set generated to the metadata path for backwards compat
	viper.SetDefault(config.Generated, viper.GetString(config.Metadata))

	// Watch for changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		instance.refreshConfig()
	})

	//viper.Set("stash", []string{"/", "/stuff"})
	//viper.WriteConfig()
}

func initFlags() {
	pflag.IP("host", net.IPv4(0, 0, 0, 0), "ip address for the host")
	pflag.Int("port", 9999, "port to serve from")
	pflag.Int("verbose", 0, "verbosity level")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logger.Infof("failed to bind flags: %s", err.Error())
	}
}

func initEnvs() {
	viper.SetEnvPrefix("stash") // will be uppercased automatically
	viper.BindEnv("host")       // STASH_HOST
	viper.BindEnv("port")       // STASH_PORT
	viper.BindEnv("stash")      // STASH_STASH
	viper.BindEnv("generated")  // STASH_GENERATED
	viper.BindEnv("metadata")   // STASH_METADATA
	viper.BindEnv("cache")      // STASH_CACHE
}

func initFFMPEG() {
	configDirectory := paths.GetConfigDirectory()
	ffmpegPath, ffprobePath := ffmpeg.GetPaths(configDirectory)
	if ffmpegPath == "" || ffprobePath == "" {
		logger.Infof("couldn't find FFMPEG, attempting to download it")
		if err := ffmpeg.Download(configDirectory); err != nil {
			msg := `Unable to locate / automatically download FFMPEG

Check the readme for download links.
The FFMPEG and FFProbe binaries should be placed in %s

The error was: %s
`
			logger.Fatalf(msg, configDirectory, err)
		}
	}

	// TODO: is this valid after download?
	instance.FFMPEGPath = ffmpegPath
	instance.FFProbePath = ffprobePath
}

func (s *singleton) refreshConfig() {
	s.Paths = paths.NewPaths()
	if config.IsValid() {
		_ = utils.EnsureDir(s.Paths.Generated.Screenshots)
		_ = utils.EnsureDir(s.Paths.Generated.Vtt)
		_ = utils.EnsureDir(s.Paths.Generated.Markers)
		_ = utils.EnsureDir(s.Paths.Generated.Transcodes)

		_ = utils.EnsureDir(s.Paths.JSON.Performers)
		_ = utils.EnsureDir(s.Paths.JSON.Scenes)
		_ = utils.EnsureDir(s.Paths.JSON.Galleries)
		_ = utils.EnsureDir(s.Paths.JSON.Studios)
	}
}
