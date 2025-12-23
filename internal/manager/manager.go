// Package manager provides the core manager of the application.
// This consolidates all the services and managers into a single struct.
package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/internal/dlna"
	"github.com/stashapp/stash/internal/log"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/pkg"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/sqlite"

	// register custom migrations
	_ "github.com/stashapp/stash/pkg/sqlite/migrations"
)

type Manager struct {
	Config *config.Config
	Logger *log.Logger

	// ImageThumbnailGenerateWaitGroup is the global wait group image thumbnail generation
	// It uses the parallel tasks setting from the configuration.
	ImageThumbnailGenerateWaitGroup sizedwaitgroup.SizedWaitGroup

	Paths *paths.Paths

	FFMpeg        *ffmpeg.FFMpeg
	FFProbe       *ffmpeg.FFProbe
	StreamManager *ffmpeg.StreamManager

	JobManager      *job.Manager
	ReadLockManager *fsutil.ReadLockManager

	DownloadStore *DownloadStore
	SessionStore  *session.Store

	PluginCache  *plugin.Cache
	ScraperCache *scraper.Cache

	PluginPackageManager  *pkg.Manager
	ScraperPackageManager *pkg.Manager

	DLNAService *dlna.Service

	Database   *sqlite.Database
	Repository models.Repository

	SceneService   SceneService
	ImageService   ImageService
	GalleryService GalleryService
	GroupService   GroupService

	scanSubs *subscriptionManager
}

var instance *Manager

func GetInstance() *Manager {
	if instance == nil {
		panic("manager not initialized")
	}
	return instance
}

func (s *Manager) SetBlobStoreOptions() {
	storageType := s.Config.GetBlobsStorage()
	blobsPath := s.Config.GetBlobsPath()
	extraBlobsPaths := s.Config.GetExtraBlobsPaths()

	s.Database.SetBlobStoreOptions(sqlite.BlobStoreOptions{
		UseFilesystem:      storageType == config.BlobStorageTypeFilesystem,
		UseDatabase:        storageType == config.BlobStorageTypeDatabase,
		Path:               blobsPath,
		SupplementaryPaths: extraBlobsPaths,
	})
}

func (s *Manager) RefreshConfig() {
	cfg := s.Config
	*s.Paths = paths.NewPaths(cfg.GetGeneratedPath(), cfg.GetBlobsPath())
	if cfg.Validate() == nil {
		if err := fsutil.EnsureDir(s.Paths.Generated.Screenshots); err != nil {
			logger.Warnf("could not create screenshots directory: %v", err)
		}
		if err := fsutil.EnsureDir(s.Paths.Generated.Vtt); err != nil {
			logger.Warnf("could not create VTT directory: %v", err)
		}
		if err := fsutil.EnsureDir(s.Paths.Generated.Markers); err != nil {
			logger.Warnf("could not create markers directory: %v", err)
		}
		if err := fsutil.EnsureDir(s.Paths.Generated.Transcodes); err != nil {
			logger.Warnf("could not create transcodes directory: %v", err)
		}
		if err := fsutil.EnsureDir(s.Paths.Generated.Downloads); err != nil {
			logger.Warnf("could not create downloads directory: %v", err)
		}
		if err := fsutil.EnsureDir(s.Paths.Generated.InteractiveHeatmap); err != nil {
			logger.Warnf("could not create interactive heatmaps directory: %v", err)
		}

		s.ImageThumbnailGenerateWaitGroup.Size = cfg.GetParallelTasksWithAutoDetection()
	}
}

// RefreshPluginCache refreshes the plugin cache.
// Call this when the plugin configuration changes.
func (s *Manager) RefreshPluginCache() {
	s.PluginCache.ReloadPlugins()
}

// RefreshScraperCache refreshes the scraper cache.
// Call this when the scraper configuration changes.
func (s *Manager) RefreshScraperCache() {
	s.ScraperCache.ReloadScrapers()
}

// RefreshStreamManager refreshes the stream manager.
// Call this when the cache directory changes.
func (s *Manager) RefreshStreamManager() {
	// shutdown existing manager if needed
	if s.StreamManager != nil {
		s.StreamManager.Shutdown()
		s.StreamManager = nil
	}

	cfg := s.Config
	cacheDir := cfg.GetCachePath()
	s.StreamManager = ffmpeg.NewStreamManager(cacheDir, s.FFMpeg, s.FFProbe, cfg, s.ReadLockManager)
}

// RefreshDLNA starts/stops the DLNA service as needed.
func (s *Manager) RefreshDLNA() {
	dlnaService := s.DLNAService
	enabled := s.Config.GetDLNADefaultEnabled()
	if !enabled && dlnaService.IsRunning() {
		dlnaService.Stop(nil)
	} else if enabled && !dlnaService.IsRunning() {
		if err := dlnaService.Start(nil); err != nil {
			logger.Warnf("error starting DLNA service: %v", err)
		}
	}
}

func createPackageManager(localPath string, srcPathGetter pkg.SourcePathGetter) *pkg.Manager {
	const timeout = 10 * time.Second
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: timeout,
	}

	return &pkg.Manager{
		Local: &pkg.Store{
			BaseDir:      localPath,
			ManifestFile: pkg.ManifestFile,
		},
		PackagePathGetter: srcPathGetter,
		Client:            httpClient,
	}
}

func (s *Manager) RefreshScraperSourceManager() {
	s.ScraperPackageManager = createPackageManager(s.Config.GetScrapersPath(), s.Config.GetScraperPackagePathGetter())
}

func (s *Manager) RefreshPluginSourceManager() {
	s.PluginPackageManager = createPackageManager(s.Config.GetPluginsPath(), s.Config.GetPluginPackagePathGetter())
}

func setSetupDefaults(input *SetupInput) {
	if input.ConfigLocation == "" {
		input.ConfigLocation = filepath.Join(fsutil.GetHomeDirectory(), ".stash", "config.yml")
	}

	configDir := filepath.Dir(input.ConfigLocation)
	if input.GeneratedLocation == "" {
		input.GeneratedLocation = filepath.Join(configDir, "generated")
	}
	if input.CacheLocation == "" {
		input.CacheLocation = filepath.Join(configDir, "cache")
	}

	if input.DatabaseFile == "" {
		input.DatabaseFile = filepath.Join(configDir, "stash-go.sqlite")
	}

	if input.BlobsLocation == "" {
		input.BlobsLocation = filepath.Join(configDir, "blobs")
	}
}

func (s *Manager) Setup(ctx context.Context, input SetupInput) error {
	setSetupDefaults(&input)
	cfg := s.Config

	// create the config directory if it does not exist
	// don't do anything if config is already set in the environment
	if !config.FileEnvSet() {
		// #3304 - if config path is relative, it breaks the ffmpeg/ffprobe
		// paths since they must not be relative. The config file property is
		// resolved to an absolute path when stash is run normally, so convert
		// relative paths to absolute paths during setup.
		// #6287 - this should no longer be necessary since the ffmpeg code
		// converts to absolute paths. Converting the config location to
		// absolute means that scraper and plugin paths default to absolute
		// which we don't want.
		configFile := input.ConfigLocation
		configDir := filepath.Dir(configFile)

		if exists, _ := fsutil.DirExists(configDir); !exists {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return fmt.Errorf("error creating config directory: %v", err)
			}
		}

		if err := fsutil.Touch(configFile); err != nil {
			return fmt.Errorf("error creating config file: %v", err)
		}

		s.Config.SetConfigFile(configFile)
	}

	if err := cfg.SetInitialConfig(); err != nil {
		return fmt.Errorf("error setting initial configuration: %v", err)
	}

	// create the generated directory if it does not exist
	if !cfg.HasOverride(config.Generated) {
		if exists, _ := fsutil.DirExists(input.GeneratedLocation); !exists {
			if err := os.MkdirAll(input.GeneratedLocation, 0755); err != nil {
				return fmt.Errorf("error creating generated directory: %v", err)
			}
		}

		s.Config.SetString(config.Generated, input.GeneratedLocation)
	}

	// create the cache directory if it does not exist
	if !cfg.HasOverride(config.Cache) {
		if exists, _ := fsutil.DirExists(input.CacheLocation); !exists {
			if err := os.MkdirAll(input.CacheLocation, 0755); err != nil {
				return fmt.Errorf("error creating cache directory: %v", err)
			}
		}

		cfg.SetString(config.Cache, input.CacheLocation)
	}

	if input.SFWContentMode {
		cfg.SetBool(config.SFWContentMode, true)
	}

	if input.StoreBlobsInDatabase {
		cfg.SetInterface(config.BlobsStorage, config.BlobStorageTypeDatabase)
	} else {
		if !cfg.HasOverride(config.BlobsPath) {
			if exists, _ := fsutil.DirExists(input.BlobsLocation); !exists {
				if err := os.MkdirAll(input.BlobsLocation, 0755); err != nil {
					return fmt.Errorf("error creating blobs directory: %v", err)
				}
			}

			cfg.SetString(config.BlobsPath, input.BlobsLocation)
		}

		cfg.SetInterface(config.BlobsStorage, config.BlobStorageTypeFilesystem)
	}

	// set the configuration
	if !cfg.HasOverride(config.Database) {
		cfg.SetString(config.Database, input.DatabaseFile)
	}

	cfg.SetInterface(config.Stash, input.Stashes)

	if err := cfg.Write(); err != nil {
		return fmt.Errorf("error writing configuration file: %v", err)
	}

	// finish initialization
	if err := s.postInit(ctx); err != nil {
		return fmt.Errorf("error completing initialization: %v", err)
	}

	cfg.FinalizeSetup()

	return nil
}

func (s *Manager) validateFFmpeg() error {
	if s.FFMpeg == nil || s.FFProbe == nil {
		return errors.New("missing ffmpeg and/or ffprobe")
	}
	return nil
}

func (s *Manager) BackupDatabase(download bool) (string, string, error) {
	var backupPath string
	var backupName string
	if download {
		backupDir := s.Paths.Generated.Downloads
		if err := fsutil.EnsureDir(backupDir); err != nil {
			return "", "", fmt.Errorf("could not create backup directory %v: %w", backupDir, err)
		}
		f, err := os.CreateTemp(backupDir, "backup*.sqlite")
		if err != nil {
			return "", "", err
		}

		backupPath = f.Name()
		backupName = s.Database.DatabaseBackupPath("")
		f.Close()

		// delete the temp file so that the backup operation can create it
		if err := os.Remove(backupPath); err != nil {
			return "", "", fmt.Errorf("could not remove temporary backup file %v: %w", backupPath, err)
		}
	} else {
		backupDir := s.Config.GetBackupDirectoryPathOrDefault()
		if backupDir != "" {
			if err := fsutil.EnsureDir(backupDir); err != nil {
				return "", "", fmt.Errorf("could not create backup directory %v: %w", backupDir, err)
			}
		}
		backupPath = s.Database.DatabaseBackupPath(backupDir)
		backupName = filepath.Base(backupPath)
	}

	err := s.Database.Backup(backupPath)
	if err != nil {
		return "", "", err
	}

	return backupPath, backupName, nil
}

func (s *Manager) AnonymiseDatabase(download bool) (string, string, error) {
	var outPath string
	var outName string
	if download {
		outDir := s.Paths.Generated.Downloads
		if err := fsutil.EnsureDir(outDir); err != nil {
			return "", "", fmt.Errorf("could not create output directory %v: %w", outDir, err)
		}
		f, err := os.CreateTemp(outDir, "anonymous*.sqlite")
		if err != nil {
			return "", "", err
		}

		outPath = f.Name()
		outName = s.Database.AnonymousDatabasePath("")
		f.Close()
	} else {
		outDir := s.Config.GetBackupDirectoryPathOrDefault()
		if outDir != "" {
			if err := fsutil.EnsureDir(outDir); err != nil {
				return "", "", fmt.Errorf("could not create output directory %v: %w", outDir, err)
			}
		}
		outPath = s.Database.AnonymousDatabasePath(outDir)
		outName = filepath.Base(outPath)
	}

	err := s.Database.Anonymise(outPath)
	if err != nil {
		return "", "", err
	}

	return outPath, outName, nil
}

func (s *Manager) GetSystemStatus() *SystemStatus {
	workingDir := fsutil.GetWorkingDirectory()
	homeDir := fsutil.GetHomeDirectory()

	database := s.Database
	dbSchema := int(database.Version())
	dbPath := database.DatabasePath()
	appSchema := int(database.AppSchemaVersion())

	status := SystemStatusEnumOk
	if s.Config.IsNewSystem() {
		status = SystemStatusEnumSetup
	} else if dbSchema < appSchema {
		status = SystemStatusEnumNeedsMigration
	}

	configFile := s.Config.GetConfigFile()

	ffmpegPath := ""
	if s.FFMpeg != nil {
		ffmpegPath = s.FFMpeg.Path()
	}

	ffprobePath := ""
	if s.FFProbe != nil {
		ffprobePath = s.FFProbe.Path()
	}

	return &SystemStatus{
		Os:             runtime.GOOS,
		WorkingDir:     workingDir,
		HomeDir:        homeDir,
		DatabaseSchema: &dbSchema,
		DatabasePath:   &dbPath,
		AppSchema:      appSchema,
		Status:         status,
		ConfigPath:     &configFile,
		FfmpegPath:     &ffmpegPath,
		FfprobePath:    &ffprobePath,
	}
}

// Shutdown gracefully stops the manager
func (s *Manager) Shutdown() {
	// TODO: Each part of the manager needs to gracefully stop at some point

	if s.StreamManager != nil {
		s.StreamManager.Shutdown()
		s.StreamManager = nil
	}

	err := s.Database.Close()
	if err != nil {
		logger.Errorf("Error closing database: %s", err)
	}
}
