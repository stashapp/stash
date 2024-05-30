package manager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/internal/desktop"
	"github.com/stashapp/stash/internal/dlna"
	"github.com/stashapp/stash/internal/log"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scraper"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/utils"
	"github.com/stashapp/stash/ui"
)

// Called at startup
func Initialize(cfg *config.Config, l *log.Logger) (*Manager, error) {
	ctx := context.TODO()

	db := sqlite.NewDatabase()
	repo := db.Repository()

	// start with empty paths
	mgrPaths := &paths.Paths{}

	scraperRepository := scraper.NewRepository(repo)
	scraperCache := scraper.NewCache(cfg, scraperRepository)

	pluginCache := plugin.NewCache(cfg)

	sceneService := &scene.Service{
		File:             db.File,
		Repository:       db.Scene,
		MarkerRepository: db.SceneMarker,
		PluginCache:      pluginCache,
		Paths:            mgrPaths,
		Config:           cfg,
	}

	imageService := &image.Service{
		File:       db.File,
		Repository: db.Image,
	}

	galleryService := &gallery.Service{
		Repository:   db.Gallery,
		ImageFinder:  db.Image,
		ImageService: imageService,
		File:         db.File,
		Folder:       db.Folder,
	}

	sceneServer := &SceneServer{
		TxnManager:       repo.TxnManager,
		SceneCoverGetter: repo.Scene,
	}

	dlnaRepository := dlna.NewRepository(repo)
	dlnaService := dlna.NewService(dlnaRepository, cfg, sceneServer)

	mgr := &Manager{
		Config: cfg,
		Logger: l,

		Paths: mgrPaths,

		ImageThumbnailGenerateWaitGroup: sizedwaitgroup.New(1),

		JobManager:      initJobManager(cfg),
		ReadLockManager: fsutil.NewReadLockManager(),

		DownloadStore: NewDownloadStore(),

		PluginCache:  pluginCache,
		ScraperCache: scraperCache,

		DLNAService: dlnaService,

		Database:   db,
		Repository: repo,

		SceneService:   sceneService,
		ImageService:   imageService,
		GalleryService: galleryService,

		scanSubs: &subscriptionManager{},
	}

	if !cfg.IsNewSystem() {
		logger.Infof("using config file: %s", cfg.GetConfigFile())

		err := cfg.Validate()
		if err != nil {
			return nil, fmt.Errorf("invalid configuration: %w", err)
		}

		if err := mgr.postInit(ctx); err != nil {
			return nil, err
		}

		mgr.checkSecurityTripwire()
	} else {
		cfgFile := cfg.GetConfigFile()
		if cfgFile != "" {
			cfgFile += " "
		}

		// create temporary session store - this will be re-initialised
		// after config is complete
		mgr.SessionStore = session.NewStore(cfg)

		logger.Warnf("config file %snot found. Assuming new system...", cfgFile)
	}

	instance = mgr
	return mgr, nil
}

func formatDuration(t time.Duration) string {
	switch {
	case t >= time.Minute: // 1m23s or 2h45m12s
		t = t.Round(time.Second)
	case t >= time.Second: // 45.36s
		t = t.Round(10 * time.Millisecond)
	default: // 51ms
		t = t.Round(time.Millisecond)
	}

	return t.String()
}

func initJobManager(cfg *config.Config) *job.Manager {
	ret := job.NewManager()

	// desktop notifications
	ctx := context.Background()
	c := ret.Subscribe(context.Background())
	go func() {
		for {
			select {
			case j := <-c.RemovedJob:
				if cfg.GetNotificationsEnabled() {
					cleanDesc := strings.TrimRight(j.Description, ".")

					if j.StartTime == nil {
						// Task was never started
						return
					}

					timeElapsed := j.EndTime.Sub(*j.StartTime)
					msg := fmt.Sprintf("Task \"%s\" finished in %s.", cleanDesc, formatDuration(timeElapsed))
					desktop.SendNotification("Task Finished", msg)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ret
}

// postInit initialises the paths, caches and database after the initial
// configuration has been set. Should only be called if the configuration
// is valid.
func (s *Manager) postInit(ctx context.Context) error {
	s.RefreshConfig()

	s.SessionStore = session.NewStore(s.Config)
	s.PluginCache.RegisterSessionStore(s.SessionStore)

	s.RefreshPluginCache()
	s.RefreshPluginSourceManager()

	s.RefreshScraperCache()
	s.RefreshScraperSourceManager()

	s.RefreshDLNA()

	s.SetBlobStoreOptions()

	s.writeStashIcon()

	// clear the downloads and tmp directories
	// #1021 - only clear these directories if the generated folder is non-empty
	if s.Config.GetGeneratedPath() != "" {
		const deleteTimeout = 1 * time.Second

		utils.Timeout(func() {
			if err := fsutil.EmptyDir(s.Paths.Generated.Downloads); err != nil {
				logger.Warnf("could not empty downloads directory: %v", err)
			}
			if err := fsutil.EnsureDir(s.Paths.Generated.Tmp); err != nil {
				logger.Warnf("could not create temporary directory: %v", err)
			} else {
				if err := fsutil.EmptyDir(s.Paths.Generated.Tmp); err != nil {
					logger.Warnf("could not empty temporary directory: %v", err)
				}
			}
		}, deleteTimeout, func(done chan struct{}) {
			logger.Info("Please wait. Deleting temporary files...") // print
			<-done                                                  // and wait for deletion
			logger.Info("Temporary files deleted.")
		})
	}

	if err := s.Database.Open(s.Config.GetDatabasePath()); err != nil {
		var migrationNeededErr *sqlite.MigrationNeededError
		if errors.As(err, &migrationNeededErr) {
			logger.Warn(err)
		} else {
			return err
		}
	}

	// Set the proxy if defined in config
	if s.Config.GetProxy() != "" {
		os.Setenv("HTTP_PROXY", s.Config.GetProxy())
		os.Setenv("HTTPS_PROXY", s.Config.GetProxy())
		os.Setenv("NO_PROXY", s.Config.GetNoProxy())
		logger.Info("Using HTTP proxy")
	}

	s.RefreshFFMpeg(ctx)
	s.RefreshStreamManager()

	return nil
}

func (s *Manager) checkSecurityTripwire() {
	if err := session.CheckExternalAccessTripwire(s.Config); err != nil {
		session.LogExternalAccessError(*err)
	}
}

func (s *Manager) writeStashIcon() {
	iconPath := filepath.Join(s.Config.GetConfigPath(), "icon.png")
	err := os.WriteFile(iconPath, ui.FaviconProvider.GetFaviconPng(), 0644)
	if err != nil {
		logger.Errorf("Couldn't write icon file: %v", err)
	}
}

func (s *Manager) RefreshFFMpeg(ctx context.Context) {
	// use same directory as config path
	// executing binaries requires directory to be included
	// https://pkg.go.dev/os/exec#hdr-Executables_in_the_current_directory
	configDirectory := s.Config.GetConfigPathAbs()
	stashHomeDir := paths.GetStashHomeDirectory()

	// prefer the configured paths
	ffmpegPath := s.Config.GetFFMpegPath()
	ffprobePath := s.Config.GetFFProbePath()

	// ensure the paths are valid
	if ffmpegPath != "" {
		// path was set explicitly
		if err := ffmpeg.ValidateFFMpeg(ffmpegPath); err != nil {
			logger.Errorf("invalid ffmpeg path: %v", err)
			return
		}

		if err := ffmpeg.ValidateFFMpegCodecSupport(ffmpegPath); err != nil {
			logger.Warn(err)
		}
	} else {
		ffmpegPath = ffmpeg.ResolveFFMpeg(configDirectory, stashHomeDir)
	}

	if ffprobePath != "" {
		if err := ffmpeg.ValidateFFProbe(ffmpegPath); err != nil {
			logger.Errorf("invalid ffprobe path: %v", err)
			return
		}
	} else {
		ffprobePath = ffmpeg.ResolveFFProbe(configDirectory, stashHomeDir)
	}

	if ffmpegPath == "" {
		logger.Warn("Couldn't find FFmpeg")
	}
	if ffprobePath == "" {
		logger.Warn("Couldn't find FFProbe")
	}

	if ffmpegPath != "" && ffprobePath != "" {
		logger.Debugf("using ffmpeg: %s", ffmpegPath)
		logger.Debugf("using ffprobe: %s", ffprobePath)

		s.FFMpeg = ffmpeg.NewEncoder(ffmpegPath)
		s.FFProbe = ffmpeg.FFProbe(ffprobePath)

		s.FFMpeg.InitHWSupport(ctx)
	}
}
