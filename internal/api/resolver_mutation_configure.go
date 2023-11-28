package api

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

var ErrOverriddenConfig = errors.New("cannot set overridden value")

func (r *mutationResolver) Setup(ctx context.Context, input manager.SetupInput) (bool, error) {
	err := manager.GetInstance().Setup(ctx, input)
	return err == nil, err
}

func (r *mutationResolver) Migrate(ctx context.Context, input manager.MigrateInput) (bool, error) {
	err := manager.GetInstance().Migrate(ctx, input)
	return err == nil, err
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
		c.Set(config.Stash, input.Stashes)
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
		c.Set(config.Database, input.DatabasePath)
	}

	existingBackupDirectoryPath := c.GetBackupDirectoryPath()
	if input.BackupDirectoryPath != nil && existingBackupDirectoryPath != *input.BackupDirectoryPath {
		if err := validateDir(config.BackupDirectoryPath, *input.BackupDirectoryPath, true); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.Set(config.BackupDirectoryPath, input.BackupDirectoryPath)
	}

	existingGeneratedPath := c.GetGeneratedPath()
	if input.GeneratedPath != nil && existingGeneratedPath != *input.GeneratedPath {
		if err := validateDir(config.Generated, *input.GeneratedPath, false); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.Set(config.Generated, input.GeneratedPath)
	}

	refreshScraperCache := false
	existingScrapersPath := c.GetScrapersPath()
	if input.ScrapersPath != nil && existingScrapersPath != *input.ScrapersPath {
		if err := validateDir(config.ScrapersPath, *input.ScrapersPath, false); err != nil {
			return makeConfigGeneralResult(), err
		}

		refreshScraperCache = true
		c.Set(config.ScrapersPath, input.ScrapersPath)
	}

	existingMetadataPath := c.GetMetadataPath()
	if input.MetadataPath != nil && existingMetadataPath != *input.MetadataPath {
		if err := validateDir(config.Metadata, *input.MetadataPath, true); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.Set(config.Metadata, input.MetadataPath)
	}

	refreshStreamManager := false
	existingCachePath := c.GetCachePath()
	if input.CachePath != nil && existingCachePath != *input.CachePath {
		if err := validateDir(config.Cache, *input.CachePath, true); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.Set(config.Cache, input.CachePath)
		refreshStreamManager = true
	}

	refreshBlobStorage := false
	existingBlobsPath := c.GetBlobsPath()
	if input.BlobsPath != nil && existingBlobsPath != *input.BlobsPath {
		if err := validateDir(config.BlobsPath, *input.BlobsPath, true); err != nil {
			return makeConfigGeneralResult(), err
		}

		c.Set(config.BlobsPath, input.BlobsPath)
		refreshBlobStorage = true
	}

	if input.BlobsStorage != nil && *input.BlobsStorage != c.GetBlobsStorage() {
		if *input.BlobsStorage == config.BlobStorageTypeFilesystem && c.GetBlobsPath() == "" {
			return makeConfigGeneralResult(), fmt.Errorf("blobs path must be set when using filesystem storage")
		}

		// TODO - migrate between systems
		c.Set(config.BlobsStorage, input.BlobsStorage)

		refreshBlobStorage = true
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

		c.Set(config.VideoFileNamingAlgorithm, *input.VideoFileNamingAlgorithm)
	}

	if input.CalculateMd5 != nil {
		c.Set(config.CalculateMD5, *input.CalculateMd5)
	}

	if input.ParallelTasks != nil {
		c.Set(config.ParallelTasks, *input.ParallelTasks)
	}

	if input.PreviewAudio != nil {
		c.Set(config.PreviewAudio, *input.PreviewAudio)
	}

	if input.PreviewSegments != nil {
		c.Set(config.PreviewSegments, *input.PreviewSegments)
	}
	if input.PreviewSegmentDuration != nil {
		c.Set(config.PreviewSegmentDuration, *input.PreviewSegmentDuration)
	}
	if input.PreviewExcludeStart != nil {
		c.Set(config.PreviewExcludeStart, *input.PreviewExcludeStart)
	}
	if input.PreviewExcludeEnd != nil {
		c.Set(config.PreviewExcludeEnd, *input.PreviewExcludeEnd)
	}
	if input.PreviewPreset != nil {
		c.Set(config.PreviewPreset, input.PreviewPreset.String())
	}

	if input.TranscodeHardwareAcceleration != nil {
		c.Set(config.TranscodeHardwareAcceleration, *input.TranscodeHardwareAcceleration)
	}
	if input.MaxTranscodeSize != nil {
		c.Set(config.MaxTranscodeSize, input.MaxTranscodeSize.String())
	}

	if input.MaxStreamingTranscodeSize != nil {
		c.Set(config.MaxStreamingTranscodeSize, input.MaxStreamingTranscodeSize.String())
	}

	if input.WriteImageThumbnails != nil {
		c.Set(config.WriteImageThumbnails, *input.WriteImageThumbnails)
	}

	if input.CreateImageClipsFromVideos != nil {
		c.Set(config.CreateImageClipsFromVideos, *input.CreateImageClipsFromVideos)
	}

	if input.GalleryCoverRegex != nil {

		_, err := regexp.Compile(*input.GalleryCoverRegex)
		if err != nil {
			return makeConfigGeneralResult(), fmt.Errorf("Gallery cover regex '%v' invalid, '%v'", *input.GalleryCoverRegex, err.Error())
		}

		c.Set(config.GalleryCoverRegex, *input.GalleryCoverRegex)
	}

	if input.Username != nil && *input.Username != c.GetUsername() {
		c.Set(config.Username, input.Username)
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

	if input.MaxSessionAge != nil {
		c.Set(config.MaxSessionAge, *input.MaxSessionAge)
	}

	if input.LogFile != nil {
		c.Set(config.LogFile, input.LogFile)
	}

	if input.LogOut != nil {
		c.Set(config.LogOut, *input.LogOut)
	}

	if input.LogAccess != nil {
		c.Set(config.LogAccess, *input.LogAccess)
	}

	if input.LogLevel != nil && *input.LogLevel != c.GetLogLevel() {
		c.Set(config.LogLevel, input.LogLevel)
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
		c.Set(config.Exclude, input.Excludes)
	}

	if input.ImageExcludes != nil {
		for _, r := range input.ImageExcludes {
			_, err := regexp.Compile(r)
			if err != nil {
				return makeConfigGeneralResult(), fmt.Errorf("image/gallery exclusion pattern '%v' invalid: %w", r, err)
			}
		}
		c.Set(config.ImageExclude, input.ImageExcludes)
	}

	if input.VideoExtensions != nil {
		c.Set(config.VideoExtensions, input.VideoExtensions)
	}

	if input.ImageExtensions != nil {
		c.Set(config.ImageExtensions, input.ImageExtensions)
	}

	if input.GalleryExtensions != nil {
		c.Set(config.GalleryExtensions, input.GalleryExtensions)
	}

	if input.CreateGalleriesFromFolders != nil {
		c.Set(config.CreateGalleriesFromFolders, input.CreateGalleriesFromFolders)
	}

	if input.CustomPerformerImageLocation != nil {
		c.Set(config.CustomPerformerImageLocation, *input.CustomPerformerImageLocation)
		initCustomPerformerImages(*input.CustomPerformerImageLocation)
	}

	if input.StashBoxes != nil {
		if err := c.ValidateStashBoxes(input.StashBoxes); err != nil {
			return nil, err
		}
		c.Set(config.StashBoxes, input.StashBoxes)
	}

	if input.PythonPath != nil {
		c.Set(config.PythonPath, input.PythonPath)
	}

	if input.TranscodeInputArgs != nil {
		c.Set(config.TranscodeInputArgs, input.TranscodeInputArgs)
	}
	if input.TranscodeOutputArgs != nil {
		c.Set(config.TranscodeOutputArgs, input.TranscodeOutputArgs)
	}
	if input.LiveTranscodeInputArgs != nil {
		c.Set(config.LiveTranscodeInputArgs, input.LiveTranscodeInputArgs)
	}
	if input.LiveTranscodeOutputArgs != nil {
		c.Set(config.LiveTranscodeOutputArgs, input.LiveTranscodeOutputArgs)
	}

	if input.DrawFunscriptHeatmapRange != nil {
		c.Set(config.DrawFunscriptHeatmapRange, input.DrawFunscriptHeatmapRange)
	}

	refreshScraperSource := false
	if input.ScraperPackageSources != nil {
		c.Set(config.ScraperPackageSources, input.ScraperPackageSources)
		refreshScraperSource = true
	}

	refreshPluginSource := false
	if input.PluginPackageSources != nil {
		c.Set(config.PluginPackageSources, input.PluginPackageSources)
		refreshPluginSource = true
	}

	if err := c.Write(); err != nil {
		return makeConfigGeneralResult(), err
	}

	manager.GetInstance().RefreshConfig()
	if refreshScraperCache {
		manager.GetInstance().RefreshScraperCache()
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

	setBool := func(key string, v *bool) {
		if v != nil {
			c.Set(key, *v)
		}
	}

	setString := func(key string, v *string) {
		if v != nil {
			c.Set(key, *v)
		}
	}

	if input.MenuItems != nil {
		c.Set(config.MenuItems, input.MenuItems)
	}

	setBool(config.SoundOnPreview, input.SoundOnPreview)
	setBool(config.WallShowTitle, input.WallShowTitle)

	setBool(config.NoBrowser, input.NoBrowser)

	setBool(config.NotificationsEnabled, input.NotificationsEnabled)

	setBool(config.ShowScrubber, input.ShowScrubber)

	if input.WallPlayback != nil {
		c.Set(config.WallPlayback, *input.WallPlayback)
	}

	if input.MaximumLoopDuration != nil {
		c.Set(config.MaximumLoopDuration, *input.MaximumLoopDuration)
	}

	setBool(config.AutostartVideo, input.AutostartVideo)
	setBool(config.ShowStudioAsText, input.ShowStudioAsText)
	setBool(config.AutostartVideoOnPlaySelected, input.AutostartVideoOnPlaySelected)
	setBool(config.ContinuePlaylistDefault, input.ContinuePlaylistDefault)

	if input.Language != nil {
		c.Set(config.Language, *input.Language)
	}

	if input.ImageLightbox != nil {
		options := input.ImageLightbox

		if options.SlideshowDelay != nil {
			c.Set(config.ImageLightboxSlideshowDelay, *options.SlideshowDelay)
		}

		setString(config.ImageLightboxDisplayModeKey, (*string)(options.DisplayMode))
		setBool(config.ImageLightboxScaleUp, options.ScaleUp)
		setBool(config.ImageLightboxResetZoomOnNav, options.ResetZoomOnNav)
		setString(config.ImageLightboxScrollModeKey, (*string)(options.ScrollMode))

		if options.ScrollAttemptsBeforeChange != nil {
			c.Set(config.ImageLightboxScrollAttemptsBeforeChange, *options.ScrollAttemptsBeforeChange)
		}
	}

	if input.CSS != nil {
		c.SetCSS(*input.CSS)
	}

	setBool(config.CSSEnabled, input.CSSEnabled)

	if input.Javascript != nil {
		c.SetJavascript(*input.Javascript)
	}

	setBool(config.JavascriptEnabled, input.JavascriptEnabled)

	if input.CustomLocales != nil {
		c.SetCustomLocales(*input.CustomLocales)
	}

	setBool(config.CustomLocalesEnabled, input.CustomLocalesEnabled)

	if input.DisableDropdownCreate != nil {
		ddc := input.DisableDropdownCreate
		setBool(config.DisableDropdownCreatePerformer, ddc.Performer)
		setBool(config.DisableDropdownCreateStudio, ddc.Studio)
		setBool(config.DisableDropdownCreateTag, ddc.Tag)
		setBool(config.DisableDropdownCreateMovie, ddc.Movie)
	}

	if input.HandyKey != nil {
		c.Set(config.HandyKey, *input.HandyKey)
	}

	if input.FunscriptOffset != nil {
		c.Set(config.FunscriptOffset, *input.FunscriptOffset)
	}

	if input.UseStashHostedFunscript != nil {
		c.Set(config.UseStashHostedFunscript, *input.UseStashHostedFunscript)
	}

	if err := c.Write(); err != nil {
		return makeConfigInterfaceResult(), err
	}

	return makeConfigInterfaceResult(), nil
}

func (r *mutationResolver) ConfigureDlna(ctx context.Context, input ConfigDLNAInput) (*ConfigDLNAResult, error) {
	c := config.GetInstance()

	if input.ServerName != nil {
		c.Set(config.DLNAServerName, *input.ServerName)
	}

	if input.WhitelistedIPs != nil {
		c.Set(config.DLNADefaultIPWhitelist, input.WhitelistedIPs)
	}

	if input.VideoSortOrder != nil {
		c.Set(config.DLNAVideoSortOrder, input.VideoSortOrder)
	}

	refresh := false
	if input.Enabled != nil {
		c.Set(config.DLNADefaultEnabled, *input.Enabled)
		refresh = true
	}

	if input.Interfaces != nil {
		c.Set(config.DLNAInterfaces, input.Interfaces)
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
		c.Set(config.ScraperUserAgent, input.ScraperUserAgent)
		refreshScraperCache = true
	}

	if input.ScraperCDPPath != nil {
		c.Set(config.ScraperCDPPath, input.ScraperCDPPath)
		refreshScraperCache = true
	}

	if input.ExcludeTagPatterns != nil {
		for _, r := range input.ExcludeTagPatterns {
			_, err := regexp.Compile(r)
			if err != nil {
				return makeConfigScrapingResult(), fmt.Errorf("tag exclusion pattern '%v' invalid: %w", r, err)
			}
		}
		c.Set(config.ScraperExcludeTagPatterns, input.ExcludeTagPatterns)
	}

	if input.ScraperCertCheck != nil {
		c.Set(config.ScraperCertCheck, input.ScraperCertCheck)
	}

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
		c.Set(config.DefaultIdentifySettings, input.Identify)
	}

	if input.Scan != nil {
		// if input.Scan is used then ScanMetadataOptions is included in the config file
		// this causes the values to not be read correctly
		c.Set(config.DefaultScanSettings, input.Scan.ScanMetadataOptions)
	}

	if input.AutoTag != nil {
		c.Set(config.DefaultAutoTagSettings, input.AutoTag)
	}

	if input.Generate != nil {
		c.Set(config.DefaultGenerateSettings, input.Generate)
	}

	if input.DeleteFile != nil {
		c.Set(config.DeleteFileDefault, *input.DeleteFile)
	}

	if input.DeleteGenerated != nil {
		c.Set(config.DeleteGeneratedDefault, *input.DeleteGenerated)
	}

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

	c.Set(config.ApiKey, newAPIKey)
	if err := c.Write(); err != nil {
		return newAPIKey, err
	}

	return newAPIKey, nil
}

func (r *mutationResolver) ConfigureUI(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	c := config.GetInstance()
	c.SetUIConfiguration(input)

	if err := c.Write(); err != nil {
		return c.GetUIConfiguration(), err
	}

	return c.GetUIConfiguration(), nil
}

func (r *mutationResolver) ConfigureUISetting(ctx context.Context, key string, value interface{}) (map[string]interface{}, error) {
	c := config.GetInstance()

	cfg := c.GetUIConfiguration()
	cfg[key] = value

	return r.ConfigureUI(ctx, cfg)
}

func (r *mutationResolver) ConfigurePlugin(ctx context.Context, pluginID string, input map[string]interface{}) (map[string]interface{}, error) {
	c := config.GetInstance()
	c.SetPluginConfiguration(pluginID, input)

	if err := c.Write(); err != nil {
		return c.GetPluginConfiguration(pluginID), err
	}

	return c.GetPluginConfiguration(pluginID), nil
}
