package api

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/internal/manager/task"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

var ErrOverriddenConfig = errors.New("cannot set overridden value")

func (r *mutationResolver) Setup(ctx context.Context, input manager.SetupInput) (bool, error) {
	err := manager.GetInstance().Setup(ctx, input)
	return err == nil, err
}

func (r *mutationResolver) DownloadFFMpeg(ctx context.Context) (string, error) {
	mgr := manager.GetInstance()
	configDir := mgr.Config.GetConfigPathAbs()

	// don't run if ffmpeg is already installed
	ffmpegPath := ffmpeg.FindFFMpeg(configDir)
	ffprobePath := ffmpeg.FindFFProbe(configDir)
	if ffmpegPath != "" && ffprobePath != "" {
		return "", fmt.Errorf("ffmpeg and ffprobe already installed at %s and %s", ffmpegPath, ffprobePath)
	}

	t := &task.DownloadFFmpegJob{
		ConfigDirectory: configDir,
		OnComplete: func(ctx context.Context) {
			// clear the ffmpeg and ffprobe paths
			logger.Infof("Clearing ffmpeg and ffprobe config paths so they are resolved from the config directory")
			mgr.Config.SetString(config.FFMpegPath, "")
			mgr.Config.SetString(config.FFProbePath, "")
			mgr.RefreshFFMpeg(ctx)
			mgr.RefreshStreamManager()
		},
	}

	jobID := mgr.JobManager.Add(ctx, "Downloading ffmpeg...", t)

	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) setConfigString(key string, value *string) {
	c := config.GetInstance()
	if value != nil {
		c.SetString(key, *value)
	}
}

func (r *mutationResolver) setConfigBool(key string, value *bool) {
	c := config.GetInstance()
	if value != nil {
		c.SetBool(key, *value)
	}
}

func (r *mutationResolver) setConfigInt(key string, value *int) {
	c := config.GetInstance()
	if value != nil {
		c.SetInt(key, *value)
	}
}

func (r *mutationResolver) setConfigFloat(key string, value *float64) {
	c := config.GetInstance()
	if value != nil {
		c.SetFloat(key, *value)
	}
}

func (r *mutationResolver) ConfigureGeneral(ctx context.Context, input ConfigGeneralInput) (*ConfigGeneralResult, error) {
	c := config.GetInstance()

	existingPaths := c.GetStashPaths()
	if input.Stashes != nil {
		for _, s := range input.Stashes {
			// Only validate existence of new paths
			isNew := true
			for _, path := range existingPaths {
				if path.Path == s.Path {
					isNew = false
					break
				}
			}
			if isNew {
				exists, err := fsutil.DirExists(s.Path)
				if !exists {
					return makeConfigGeneralResult(), err
				}
			}
		}
		c.SetInterface(config.Stash, input.Stashes)
	}

	checkConfigOverride := func(key string) error {
		if c.HasOverride(key) {
			return fmt.Errorf("%w: %s", ErrOverriddenConfig, key)
		}

		return nil
	}

	validateDir := func(key string, value string, optional bool) error {
		if err := checkConfigOverride(key); err != nil {
			return err
		}

		if !optional || value != "" {
			if err := fsutil.EnsureDir(value); err != nil {
				return err
			}
		}

		return nil
	}

	existingDBPath := c.GetDatabasePath()
	if input.DatabasePath != nil && existingDBPath != *input.DatabasePath {
		if err := checkConfigOverride(config.Database); err != nil {
			return makeConfigGeneralResult(), err
		}

		ext := filepath.Ext(*input.DatabasePath)
		if ext != ".db" && ext != ".sqlite" && ext != ".sqlite3" {
			return makeConfigGeneralResult(), fmt.Errorf("invalid database path, use extension db, sqlite, or sqlite3")
		}
		c.SetString(config.Database, *input.DatabasePath)
	}

	existingBackupDirectoryPath := c.GetBackupDirectoryPath()
	if input.BackupDirectoryPath != nil && existingBackupDirectoryPath != *input.BackupDirectoryPath {
		if err := validateDir(config.BackupDirectoryPath, *input.BackupDirectoryPath, true); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.SetString(config.BackupDirectoryPath, *input.BackupDirectoryPath)
	}

	existingGeneratedPath := c.GetGeneratedPath()
	if input.GeneratedPath != nil && existingGeneratedPath != *input.GeneratedPath {
		if err := validateDir(config.Generated, *input.GeneratedPath, false); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.SetString(config.Generated, *input.GeneratedPath)
	}

	refreshScraperCache := false
	refreshScraperSource := false
	existingScrapersPath := c.GetScrapersPath()
	if input.ScrapersPath != nil && existingScrapersPath != *input.ScrapersPath {
		if err := validateDir(config.ScrapersPath, *input.ScrapersPath, false); err != nil {
			return makeConfigGeneralResult(), err
		}

		refreshScraperCache = true
		refreshScraperSource = true
		c.SetString(config.ScrapersPath, *input.ScrapersPath)
	}

	refreshPluginCache := false
	refreshPluginSource := false
	existingPluginsPath := c.GetPluginsPath()
	if input.PluginsPath != nil && existingPluginsPath != *input.PluginsPath {
		if err := validateDir(config.PluginsPath, *input.PluginsPath, false); err != nil {
			return makeConfigGeneralResult(), err
		}

		refreshPluginCache = true
		refreshPluginSource = true
		c.SetString(config.PluginsPath, *input.PluginsPath)
	}

	existingMetadataPath := c.GetMetadataPath()
	if input.MetadataPath != nil && existingMetadataPath != *input.MetadataPath {
		if err := validateDir(config.Metadata, *input.MetadataPath, true); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.SetString(config.Metadata, *input.MetadataPath)
	}

	refreshStreamManager := false
	existingCachePath := c.GetCachePath()
	if input.CachePath != nil && existingCachePath != *input.CachePath {
		if err := validateDir(config.Cache, *input.CachePath, true); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.SetString(config.Cache, *input.CachePath)
		refreshStreamManager = true
	}

	refreshBlobStorage := false
	existingBlobsPath := c.GetBlobsPath()
	if input.BlobsPath != nil && existingBlobsPath != *input.BlobsPath {
		if err := validateDir(config.BlobsPath, *input.BlobsPath, true); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.SetString(config.BlobsPath, *input.BlobsPath)
		refreshBlobStorage = true
	}

	if input.BlobsStorage != nil && *input.BlobsStorage != c.GetBlobsStorage() {
		if *input.BlobsStorage == config.BlobStorageTypeFilesystem && c.GetBlobsPath() == "" {
			return makeConfigGeneralResult(), fmt.Errorf("blobs path must be set when using filesystem storage")
		}

		c.SetInterface(config.BlobsStorage, *input.BlobsStorage)

		refreshBlobStorage = true
	}

	refreshFfmpeg := false
	if input.FfmpegPath != nil && *input.FfmpegPath != c.GetFFMpegPath() {
		if *input.FfmpegPath != "" {
			if err := ffmpeg.ValidateFFMpeg(*input.FfmpegPath); err != nil {
				return makeConfigGeneralResult(), fmt.Errorf("invalid ffmpeg path: %w", err)
			}
		}

		c.SetString(config.FFMpegPath, *input.FfmpegPath)
		refreshFfmpeg = true
	}

	if input.FfprobePath != nil && *input.FfprobePath != c.GetFFProbePath() {
		if *input.FfprobePath != "" {
			if err := ffmpeg.ValidateFFProbe(*input.FfprobePath); err != nil {
				return makeConfigGeneralResult(), fmt.Errorf("invalid ffprobe path: %w", err)
			}
		}

		c.SetString(config.FFProbePath, *input.FfprobePath)
		refreshFfmpeg = true
	}

	if input.VideoFileNamingAlgorithm != nil && *input.VideoFileNamingAlgorithm != c.GetVideoFileNamingAlgorithm() {
		calculateMD5 := c.IsCalculateMD5()
		if input.CalculateMd5 != nil {
			calculateMD5 = *input.CalculateMd5
		}
		if !calculateMD5 && *input.VideoFileNamingAlgorithm == models.HashAlgorithmMd5 {
			return makeConfigGeneralResult(), errors.New("calculateMD5 must be true if using MD5")
		}

		// validate changing VideoFileNamingAlgorithm
		if err := r.withTxn(context.TODO(), func(ctx context.Context) error {
			return manager.ValidateVideoFileNamingAlgorithm(ctx, r.repository.Scene, *input.VideoFileNamingAlgorithm)
		}); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.SetInterface(config.VideoFileNamingAlgorithm, *input.VideoFileNamingAlgorithm)
	}

	r.setConfigBool(config.CalculateMD5, input.CalculateMd5)
	r.setConfigInt(config.ParallelTasks, input.ParallelTasks)
	r.setConfigBool(config.PreviewAudio, input.PreviewAudio)
	r.setConfigInt(config.PreviewSegments, input.PreviewSegments)
	r.setConfigFloat(config.PreviewSegmentDuration, input.PreviewSegmentDuration)
	r.setConfigString(config.PreviewExcludeStart, input.PreviewExcludeStart)
	r.setConfigString(config.PreviewExcludeEnd, input.PreviewExcludeEnd)
	if input.PreviewPreset != nil {
		c.SetString(config.PreviewPreset, input.PreviewPreset.String())
	}

	r.setConfigBool(config.TranscodeHardwareAcceleration, input.TranscodeHardwareAcceleration)
	if input.MaxTranscodeSize != nil {
		c.SetString(config.MaxTranscodeSize, input.MaxTranscodeSize.String())
	}

	if input.MaxStreamingTranscodeSize != nil {
		c.SetString(config.MaxStreamingTranscodeSize, input.MaxStreamingTranscodeSize.String())
	}
	r.setConfigBool(config.WriteImageThumbnails, input.WriteImageThumbnails)
	r.setConfigBool(config.CreateImageClipsFromVideos, input.CreateImageClipsFromVideos)

	if input.GalleryCoverRegex != nil {
		_, err := regexp.Compile(*input.GalleryCoverRegex)
		if err != nil {
			return makeConfigGeneralResult(), fmt.Errorf("Gallery cover regex '%v' invalid, '%v'", *input.GalleryCoverRegex, err.Error())
		}

		c.SetString(config.GalleryCoverRegex, *input.GalleryCoverRegex)
	}

	if input.Username != nil && *input.Username != c.GetUsername() {
		c.SetString(config.Username, *input.Username)
		if *input.Password == "" {
			logger.Info("Username cleared")
		} else {
			logger.Info("Username changed")
		}
	}

	if input.Password != nil {
		// bit of a hack - check if the passed in password is the same as the stored hash
		// and only set if they are different
		currentPWHash := c.GetPasswordHash()

		if *input.Password != currentPWHash {
			if *input.Password == "" {
				logger.Info("Password cleared")
			} else {
				logger.Info("Password changed")
			}
			c.SetPassword(*input.Password)
		}
	}

	r.setConfigInt(config.MaxSessionAge, input.MaxSessionAge)
	r.setConfigString(config.LogFile, input.LogFile)
	r.setConfigBool(config.LogOut, input.LogOut)
	r.setConfigBool(config.LogAccess, input.LogAccess)

	if input.LogLevel != nil && *input.LogLevel != c.GetLogLevel() {
		c.SetString(config.LogLevel, *input.LogLevel)
		logger := manager.GetInstance().Logger
		logger.SetLogLevel(*input.LogLevel)
	}

	if input.Excludes != nil {
		for _, r := range input.Excludes {
			_, err := regexp.Compile(r)
			if err != nil {
				return makeConfigGeneralResult(), fmt.Errorf("video exclusion pattern '%v' invalid: %w", r, err)
			}
		}
		c.SetInterface(config.Exclude, input.Excludes)
	}

	if input.ImageExcludes != nil {
		for _, r := range input.ImageExcludes {
			_, err := regexp.Compile(r)
			if err != nil {
				return makeConfigGeneralResult(), fmt.Errorf("image/gallery exclusion pattern '%v' invalid: %w", r, err)
			}
		}
		c.SetInterface(config.ImageExclude, input.ImageExcludes)
	}

	if input.VideoExtensions != nil {
		c.SetInterface(config.VideoExtensions, input.VideoExtensions)
	}

	if input.ImageExtensions != nil {
		c.SetInterface(config.ImageExtensions, input.ImageExtensions)
	}

	if input.GalleryExtensions != nil {
		c.SetInterface(config.GalleryExtensions, input.GalleryExtensions)
	}

	r.setConfigBool(config.CreateGalleriesFromFolders, input.CreateGalleriesFromFolders)

	if input.CustomPerformerImageLocation != nil {
		c.SetString(config.CustomPerformerImageLocation, *input.CustomPerformerImageLocation)
		initCustomPerformerImages(*input.CustomPerformerImageLocation)
	}

	if input.StashBoxes != nil {
		if err := c.ValidateStashBoxes(input.StashBoxes); err != nil {
			return nil, err
		}
		c.SetInterface(config.StashBoxes, input.StashBoxes)
	}

	if input.PythonPath != nil {
		r.setConfigString(config.PythonPath, input.PythonPath)
	}

	if input.TranscodeInputArgs != nil {
		c.SetInterface(config.TranscodeInputArgs, input.TranscodeInputArgs)
	}
	if input.TranscodeOutputArgs != nil {
		c.SetInterface(config.TranscodeOutputArgs, input.TranscodeOutputArgs)
	}
	if input.LiveTranscodeInputArgs != nil {
		c.SetInterface(config.LiveTranscodeInputArgs, input.LiveTranscodeInputArgs)
	}
	if input.LiveTranscodeOutputArgs != nil {
		c.SetInterface(config.LiveTranscodeOutputArgs, input.LiveTranscodeOutputArgs)
	}

	r.setConfigBool(config.DrawFunscriptHeatmapRange, input.DrawFunscriptHeatmapRange)

	if input.ScraperPackageSources != nil {
		c.SetInterface(config.ScraperPackageSources, input.ScraperPackageSources)
		refreshScraperSource = true
	}

	if input.PluginPackageSources != nil {
		c.SetInterface(config.PluginPackageSources, input.PluginPackageSources)
		refreshPluginSource = true
	}

	if err := c.Write(); err != nil {
		return makeConfigGeneralResult(), err
	}

	manager.GetInstance().RefreshConfig()
	if refreshScraperCache {
		manager.GetInstance().RefreshScraperCache()
	}
	if refreshPluginCache {
		manager.GetInstance().RefreshPluginCache()
	}
	if refreshFfmpeg {
		manager.GetInstance().RefreshFFMpeg(ctx)

		// refresh stream manager is required since ffmpeg changed
		refreshStreamManager = true
	}
	if refreshStreamManager {
		manager.GetInstance().RefreshStreamManager()
	}
	if refreshBlobStorage {
		manager.GetInstance().SetBlobStoreOptions()
	}
	if refreshScraperSource {
		manager.GetInstance().RefreshScraperSourceManager()
	}
	if refreshPluginSource {
		manager.GetInstance().RefreshPluginSourceManager()
	}

	return makeConfigGeneralResult(), nil
}

func (r *mutationResolver) ConfigureInterface(ctx context.Context, input ConfigInterfaceInput) (*ConfigInterfaceResult, error) {
	c := config.GetInstance()

	if input.MenuItems != nil {
		c.SetInterface(config.MenuItems, input.MenuItems)
	}

	r.setConfigBool(config.SoundOnPreview, input.SoundOnPreview)
	r.setConfigBool(config.WallShowTitle, input.WallShowTitle)

	r.setConfigBool(config.NoBrowser, input.NoBrowser)

	r.setConfigBool(config.NotificationsEnabled, input.NotificationsEnabled)

	r.setConfigBool(config.ShowScrubber, input.ShowScrubber)

	r.setConfigString(config.WallPlayback, input.WallPlayback)
	r.setConfigInt(config.MaximumLoopDuration, input.MaximumLoopDuration)
	r.setConfigBool(config.AutostartVideo, input.AutostartVideo)
	r.setConfigBool(config.ShowStudioAsText, input.ShowStudioAsText)
	r.setConfigBool(config.AutostartVideoOnPlaySelected, input.AutostartVideoOnPlaySelected)
	r.setConfigBool(config.ContinuePlaylistDefault, input.ContinuePlaylistDefault)

	r.setConfigString(config.Language, input.Language)

	if input.ImageLightbox != nil {
		options := input.ImageLightbox

		r.setConfigInt(config.ImageLightboxSlideshowDelay, options.SlideshowDelay)

		r.setConfigString(config.ImageLightboxDisplayModeKey, (*string)(options.DisplayMode))
		r.setConfigBool(config.ImageLightboxScaleUp, options.ScaleUp)
		r.setConfigBool(config.ImageLightboxResetZoomOnNav, options.ResetZoomOnNav)
		r.setConfigString(config.ImageLightboxScrollModeKey, (*string)(options.ScrollMode))

		r.setConfigInt(config.ImageLightboxScrollAttemptsBeforeChange, options.ScrollAttemptsBeforeChange)
	}

	if input.CSS != nil {
		c.SetCSS(*input.CSS)
	}

	r.setConfigBool(config.CSSEnabled, input.CSSEnabled)

	if input.Javascript != nil {
		c.SetJavascript(*input.Javascript)
	}

	r.setConfigBool(config.JavascriptEnabled, input.JavascriptEnabled)

	if input.CustomLocales != nil {
		c.SetCustomLocales(*input.CustomLocales)
	}

	r.setConfigBool(config.CustomLocalesEnabled, input.CustomLocalesEnabled)

	if input.DisableDropdownCreate != nil {
		ddc := input.DisableDropdownCreate
		r.setConfigBool(config.DisableDropdownCreatePerformer, ddc.Performer)
		r.setConfigBool(config.DisableDropdownCreateStudio, ddc.Studio)
		r.setConfigBool(config.DisableDropdownCreateTag, ddc.Tag)
		r.setConfigBool(config.DisableDropdownCreateMovie, ddc.Movie)
	}

	r.setConfigString(config.HandyKey, input.HandyKey)
	r.setConfigInt(config.FunscriptOffset, input.FunscriptOffset)
	r.setConfigBool(config.UseStashHostedFunscript, input.UseStashHostedFunscript)

	if err := c.Write(); err != nil {
		return makeConfigInterfaceResult(), err
	}

	return makeConfigInterfaceResult(), nil
}

func (r *mutationResolver) ConfigureDlna(ctx context.Context, input ConfigDLNAInput) (*ConfigDLNAResult, error) {
	c := config.GetInstance()

	r.setConfigString(config.DLNAServerName, input.ServerName)

	if input.WhitelistedIPs != nil {
		c.SetInterface(config.DLNADefaultIPWhitelist, input.WhitelistedIPs)
	}

	r.setConfigString(config.DLNAVideoSortOrder, input.VideoSortOrder)
	r.setConfigInt(config.DLNAPort, input.Port)

	refresh := false
	if input.Enabled != nil {
		c.SetBool(config.DLNADefaultEnabled, *input.Enabled)
		refresh = true
	}

	if input.Interfaces != nil {
		c.SetInterface(config.DLNAInterfaces, input.Interfaces)
	}

	if err := c.Write(); err != nil {
		return makeConfigDLNAResult(), err
	}

	if refresh {
		manager.GetInstance().RefreshDLNA()
	}

	return makeConfigDLNAResult(), nil
}

func (r *mutationResolver) ConfigureScraping(ctx context.Context, input ConfigScrapingInput) (*ConfigScrapingResult, error) {
	c := config.GetInstance()

	refreshScraperCache := false
	if input.ScraperUserAgent != nil {
		c.SetString(config.ScraperUserAgent, *input.ScraperUserAgent)
		refreshScraperCache = true
	}

	if input.ScraperCDPPath != nil {
		c.SetString(config.ScraperCDPPath, *input.ScraperCDPPath)
		refreshScraperCache = true
	}

	if input.ExcludeTagPatterns != nil {
		for _, r := range input.ExcludeTagPatterns {
			_, err := regexp.Compile(r)
			if err != nil {
				return makeConfigScrapingResult(), fmt.Errorf("tag exclusion pattern '%v' invalid: %w", r, err)
			}
		}
		c.SetInterface(config.ScraperExcludeTagPatterns, input.ExcludeTagPatterns)
	}

	r.setConfigBool(config.ScraperCertCheck, input.ScraperCertCheck)

	if refreshScraperCache {
		manager.GetInstance().RefreshScraperCache()
	}
	if err := c.Write(); err != nil {
		return makeConfigScrapingResult(), err
	}

	return makeConfigScrapingResult(), nil
}

func (r *mutationResolver) ConfigureDefaults(ctx context.Context, input ConfigDefaultSettingsInput) (*ConfigDefaultSettingsResult, error) {
	c := config.GetInstance()

	if input.Identify != nil {
		c.SetInterface(config.DefaultIdentifySettings, input.Identify)
	}

	if input.Scan != nil {
		// if input.Scan is used then ScanMetadataOptions is included in the config file
		// this causes the values to not be read correctly
		c.SetInterface(config.DefaultScanSettings, input.Scan.ScanMetadataOptions)
	}

	if input.AutoTag != nil {
		c.SetInterface(config.DefaultAutoTagSettings, input.AutoTag)
	}

	if input.Generate != nil {
		c.SetInterface(config.DefaultGenerateSettings, input.Generate)
	}

	r.setConfigBool(config.DeleteFileDefault, input.DeleteFile)
	r.setConfigBool(config.DeleteGeneratedDefault, input.DeleteGenerated)

	if err := c.Write(); err != nil {
		return makeConfigDefaultsResult(), err
	}

	return makeConfigDefaultsResult(), nil
}

func (r *mutationResolver) GenerateAPIKey(ctx context.Context, input GenerateAPIKeyInput) (string, error) {
	c := config.GetInstance()

	var newAPIKey string
	if input.Clear == nil || !*input.Clear {
		username := c.GetUsername()
		if username != "" {
			var err error
			newAPIKey, err = manager.GenerateAPIKey(username)
			if err != nil {
				return "", err
			}
		}
	}

	c.SetString(config.ApiKey, newAPIKey)
	if err := c.Write(); err != nil {
		return newAPIKey, err
	}

	return newAPIKey, nil
}

func (r *mutationResolver) ConfigureUI(ctx context.Context, input map[string]interface{}, partial map[string]interface{}) (map[string]interface{}, error) {
	c := config.GetInstance()

	if input != nil {
		c.SetUIConfiguration(input)
	}

	if partial != nil {
		// merge partial into existing config
		existing := c.GetUIConfiguration()
		utils.MergeMaps(existing, partial)
		c.SetUIConfiguration(existing)
	}

	if err := c.Write(); err != nil {
		return c.GetUIConfiguration(), err
	}

	return c.GetUIConfiguration(), nil
}

func (r *mutationResolver) ConfigureUISetting(ctx context.Context, key string, value interface{}) (map[string]interface{}, error) {
	c := config.GetInstance()

	cfg := utils.NestedMap(c.GetUIConfiguration())
	cfg.Set(key, value)

	return r.ConfigureUI(ctx, cfg, nil)
}

func (r *mutationResolver) ConfigurePlugin(ctx context.Context, pluginID string, input map[string]interface{}) (map[string]interface{}, error) {
	c := config.GetInstance()
	c.SetPluginConfiguration(pluginID, input)

	if err := c.Write(); err != nil {
		return c.GetPluginConfiguration(pluginID), err
	}

	return c.GetPluginConfiguration(pluginID), nil
}
