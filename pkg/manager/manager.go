package manager

import (
	"sync"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
)

type singleton struct {
	Config *config.Instance

	Status TaskStatus
	Paths  *paths.Paths

	FFMPEGPath  string
	FFProbePath string

	PluginCache  *plugin.Cache
	ScraperCache *scraper.Cache

	DownloadStore *DownloadStore

	TxnManager models.TransactionManager
}

var instance *singleton
var once sync.Once

type flagStruct struct {
	configFilePath string
}

var flags = flagStruct{}

func GetInstance() *singleton {
	Initialize()
	return instance
}

func Initialize() *singleton {
	once.Do(func() {
		_ = utils.EnsureDir(paths.GetStashHomeDirectory())
		cfg := config.Initialize()
		initLog()

		instance = &singleton{
			Config:        cfg,
			Status:        TaskStatus{Status: Idle, Progress: -1},
			DownloadStore: NewDownloadStore(),

			TxnManager: sqlite.NewTransactionManager(),
		}

		cfgFile := cfg.GetConfigFile()
		if cfgFile != "" {
			logger.Infof("using config file: %s", cfg.GetConfigFile())

			if err := cfg.Validate(); err != nil {
				logger.Warnf("error initializing configuration: %s", err.Error())
			} else {
				instance.PostInit()
			}
		} else {
			logger.Warn("config file not found. Assuming new system...")
		}

		initFFMPEG()
	})

	return instance
}

func initFFMPEG() {
	configDirectory := paths.GetStashHomeDirectory()
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
		} else {
			// After download get new paths for ffmpeg and ffprobe
			ffmpegPath, ffprobePath = ffmpeg.GetPaths(configDirectory)
		}
	}

	instance.FFMPEGPath = ffmpegPath
	instance.FFProbePath = ffprobePath
}

func initLog() {
	config := config.GetInstance()
	logger.Init(config.GetLogFile(), config.GetLogOut(), config.GetLogLevel())
}

func initPluginCache() *plugin.Cache {
	config := config.GetInstance()
	ret, err := plugin.NewCache(config.GetPluginsPath())

	if err != nil {
		logger.Errorf("Error reading plugin configs: %s", err.Error())
	}

	return ret
}

// PostInit initialises the paths, caches and txnManager after the initial
// configuration has been set. Should only be called if the configuration
// is valid.
func (s *singleton) PostInit() {
	s.Paths = paths.NewPaths(s.Config.GetGeneratedPath())
	s.PluginCache = initPluginCache()
	s.ScraperCache = instance.initScraperCache()

	s.RefreshConfig()

	// clear the downloads and tmp directories
	// #1021 - only clear these directories if the generated folder is non-empty
	if s.Config.GetGeneratedPath() != "" {
		utils.EmptyDir(instance.Paths.Generated.Downloads)
		utils.EmptyDir(instance.Paths.Generated.Tmp)
	}

	// perform the post-migration for new databases
	if database.Initialize(s.Config.GetDatabasePath()) {
		s.PostMigrate()
	}
}

// initScraperCache initializes a new scraper cache and returns it.
func (s *singleton) initScraperCache() *scraper.Cache {
	ret, err := scraper.NewCache(config.GetInstance(), s.TxnManager)

	if err != nil {
		logger.Errorf("Error reading scraper configs: %s", err.Error())
	}

	return ret
}

func (s *singleton) RefreshConfig() {
	s.Paths = paths.NewPaths(s.Config.GetGeneratedPath())
	config := s.Config
	if config.Validate() == nil {
		utils.EnsureDir(s.Paths.Generated.Screenshots)
		utils.EnsureDir(s.Paths.Generated.Vtt)
		utils.EnsureDir(s.Paths.Generated.Markers)
		utils.EnsureDir(s.Paths.Generated.Transcodes)
		utils.EnsureDir(s.Paths.Generated.Downloads)
		paths.EnsureJSONDirs(config.GetMetadataPath())
	}
}

// RefreshScraperCache refreshes the scraper cache. Call this when scraper
// configuration changes.
func (s *singleton) RefreshScraperCache() {
	s.ScraperCache = s.initScraperCache()
}
