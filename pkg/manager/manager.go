package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
				panic(fmt.Sprintf("error initializing configuration: %s", err.Error()))
			} else {
				if err := instance.PostInit(); err != nil {
					panic(err)
				}
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
func (s *singleton) PostInit() error {
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

	if err := database.Initialize(s.Config.GetDatabasePath()); err != nil {
		return err
	}

	if database.Ready() == nil {
		s.PostMigrate()
	}

	return nil
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
	}
}

// RefreshScraperCache refreshes the scraper cache. Call this when scraper
// configuration changes.
func (s *singleton) RefreshScraperCache() {
	s.ScraperCache = s.initScraperCache()
}

func setSetupDefaults(input *models.SetupInput) {
	if input.ConfigLocation == "" {
		input.ConfigLocation = filepath.Join(utils.GetHomeDirectory(), "config.yml")
	}

	configDir := filepath.Dir(input.ConfigLocation)
	if input.GeneratedLocation == "" {
		input.GeneratedLocation = filepath.Join(configDir, "generated")
	}

	if input.DatabaseFile == "" {
		input.DatabaseFile = filepath.Join(configDir, "stash-go.sqlite")
	}
}

func (s *singleton) Setup(input models.SetupInput) error {
	setSetupDefaults(&input)

	// create the generated directory if it does not exist
	if exists, _ := utils.DirExists(input.GeneratedLocation); !exists {
		if err := os.Mkdir(input.GeneratedLocation, 0755); err != nil {
			return fmt.Errorf("error creating generated directory: %s", err.Error())
		}
	}

	if err := utils.Touch(input.ConfigLocation); err != nil {
		return fmt.Errorf("error creating config file: %s", err.Error())
	}

	s.Config.SetConfigFile(input.ConfigLocation)

	// set the configuration
	s.Config.Set(config.Generated, input.GeneratedLocation)
	s.Config.Set(config.Database, input.DatabaseFile)
	s.Config.Set(config.Stash, input.Stashes)
	if err := s.Config.Write(); err != nil {
		return fmt.Errorf("error writing configuration file: %s", err.Error())
	}

	// initialise the database
	if err := s.PostInit(); err != nil {
		return fmt.Errorf("error initializing the database: %s", err.Error())
	}

	return nil
}

func (s *singleton) Migrate(input models.MigrateInput) error {
	// always backup so that we can roll back to the previous version if
	// migration fails
	backupPath := input.BackupPath
	if backupPath == "" {
		backupPath = database.DatabaseBackupPath()
	}

	// perform database backup
	if err := database.Backup(database.DB, backupPath); err != nil {
		return fmt.Errorf("error backing up database: %s", err)
	}

	if err := database.RunMigrations(); err != nil {
		errStr := fmt.Sprintf("error performing migration: %s", err)

		// roll back to the backed up version
		restoreErr := database.RestoreFromBackup(backupPath)
		if restoreErr != nil {
			errStr = fmt.Sprintf("ERROR: unable to restore database from backup after migration failure: %s\n%s", restoreErr.Error(), errStr)
		} else {
			errStr = "An error occurred migrating the database to the latest schema version. The backup database file was automatically renamed to restore the database.\n" + errStr
		}

		return errors.New(errStr)
	}

	// perform post-migration operations
	s.PostMigrate()

	// if no backup path was provided, then delete the created backup
	if input.BackupPath == "" {
		if err := os.Remove(backupPath); err != nil {
			logger.Warnf("error removing unwanted database backup (%s): %s", backupPath, err.Error())
		}
	}

	return nil
}

func (s *singleton) GetSystemStatus() *models.SystemStatus {
	status := models.SystemStatusEnumOk
	dbSchema := int(database.Version())
	dbPath := database.DatabasePath()
	appSchema := int(database.AppSchemaVersion())

	if s.Config.GetConfigFile() == "" {
		status = models.SystemStatusEnumSetup
	} else if dbSchema < appSchema {
		status = models.SystemStatusEnumNeedsMigration
	}

	return &models.SystemStatus{
		DatabaseSchema: &dbSchema,
		DatabasePath:   &dbPath,
		AppSchema:      appSchema,
		Status:         status,
	}
}
