package api

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
	"golang.org/x/text/collate"
)

func (r *queryResolver) Configuration(ctx context.Context) (*ConfigResult, error) {
	return makeConfigResult(), nil
}

func (r *queryResolver) Directory(ctx context.Context, path, locale *string) (*Directory, error) {

	directory := &Directory{}
	var err error

	col := newCollator(locale, collate.IgnoreCase, collate.Numeric)

	var dirPath = ""
	if path != nil {
		dirPath = *path
	}
	currentDir := getDir(dirPath)
	directories, err := listDir(col, currentDir)
	if err != nil {
		return directory, err
	}

	directory.Path = currentDir
	directory.Parent = getParent(currentDir)
	directory.Directories = directories

	return directory, err
}

func getDir(path string) string {
	if path == "" {
		path = fsutil.GetHomeDirectory()
	}

	return path
}

func getParent(path string) *string {
	isRoot := path[len(path)-1:] == "/"
	if isRoot {
		return nil
	} else {
		parentPath := filepath.Clean(path + "/..")
		return &parentPath
	}
}

func makeConfigResult() *ConfigResult {
	return &ConfigResult{
		General:   makeConfigGeneralResult(),
		Interface: makeConfigInterfaceResult(),
		Dlna:      makeConfigDLNAResult(),
		Scraping:  makeConfigScrapingResult(),
		Defaults:  makeConfigDefaultsResult(),
		UI:        makeConfigUIResult(),
	}
}

func makeConfigGeneralResult() *ConfigGeneralResult {
	config := config.GetInstance()
	logFile := config.GetLogFile()

	maxTranscodeSize := config.GetMaxTranscodeSize()
	maxStreamingTranscodeSize := config.GetMaxStreamingTranscodeSize()

	customPerformerImageLocation := config.GetCustomPerformerImageLocation()

	scraperUserAgent := config.GetScraperUserAgent()
	scraperCDPPath := config.GetScraperCDPPath()

	return &ConfigGeneralResult{
		Stashes:                      config.GetStashPaths(),
		DatabasePath:                 config.GetDatabasePath(),
		GeneratedPath:                config.GetGeneratedPath(),
		MetadataPath:                 config.GetMetadataPath(),
		ConfigFilePath:               config.GetConfigFile(),
		ScrapersPath:                 config.GetScrapersPath(),
		CachePath:                    config.GetCachePath(),
		CalculateMd5:                 config.IsCalculateMD5(),
		VideoFileNamingAlgorithm:     config.GetVideoFileNamingAlgorithm(),
		ParallelTasks:                config.GetParallelTasks(),
		PreviewAudio:                 config.GetPreviewAudio(),
		PreviewSegments:              config.GetPreviewSegments(),
		PreviewSegmentDuration:       config.GetPreviewSegmentDuration(),
		PreviewExcludeStart:          config.GetPreviewExcludeStart(),
		PreviewExcludeEnd:            config.GetPreviewExcludeEnd(),
		PreviewPreset:                config.GetPreviewPreset(),
		MaxTranscodeSize:             &maxTranscodeSize,
		MaxStreamingTranscodeSize:    &maxStreamingTranscodeSize,
		WriteImageThumbnails:         config.IsWriteImageThumbnails(),
		APIKey:                       config.GetAPIKey(),
		Username:                     config.GetUsername(),
		Password:                     config.GetPasswordHash(),
		MaxSessionAge:                config.GetMaxSessionAge(),
		LogFile:                      &logFile,
		LogOut:                       config.GetLogOut(),
		LogLevel:                     config.GetLogLevel(),
		LogAccess:                    config.GetLogAccess(),
		VideoExtensions:              config.GetVideoExtensions(),
		ImageExtensions:              config.GetImageExtensions(),
		GalleryExtensions:            config.GetGalleryExtensions(),
		CreateGalleriesFromFolders:   config.GetCreateGalleriesFromFolders(),
		Excludes:                     config.GetExcludes(),
		ImageExcludes:                config.GetImageExcludes(),
		CustomPerformerImageLocation: &customPerformerImageLocation,
		ScraperUserAgent:             &scraperUserAgent,
		ScraperCertCheck:             config.GetScraperCertCheck(),
		ScraperCDPPath:               &scraperCDPPath,
		StashBoxes:                   config.GetStashBoxes(),
		PythonPath:                   config.GetPythonPath(),
	}
}

func makeConfigInterfaceResult() *ConfigInterfaceResult {
	config := config.GetInstance()
	menuItems := config.GetMenuItems()
	soundOnPreview := config.GetSoundOnPreview()
	wallShowTitle := config.GetWallShowTitle()
	showScrubber := config.GetShowScrubber()
	wallPlayback := config.GetWallPlayback()
	noBrowser := config.GetNoBrowser()
	notificationsEnabled := config.GetNotificationsEnabled()
	maximumLoopDuration := config.GetMaximumLoopDuration()
	autostartVideo := config.GetAutostartVideo()
	autostartVideoOnPlaySelected := config.GetAutostartVideoOnPlaySelected()
	continuePlaylistDefault := config.GetContinuePlaylistDefault()
	showStudioAsText := config.GetShowStudioAsText()
	css := config.GetCSS()
	cssEnabled := config.GetCSSEnabled()
	language := config.GetLanguage()
	handyKey := config.GetHandyKey()
	scriptOffset := config.GetFunscriptOffset()
	imageLightboxOptions := config.GetImageLightboxOptions()

	// FIXME - misnamed output field means we have redundant fields
	disableDropdownCreate := config.GetDisableDropdownCreate()

	return &ConfigInterfaceResult{
		MenuItems:                    menuItems,
		SoundOnPreview:               &soundOnPreview,
		WallShowTitle:                &wallShowTitle,
		WallPlayback:                 &wallPlayback,
		ShowScrubber:                 &showScrubber,
		MaximumLoopDuration:          &maximumLoopDuration,
		NoBrowser:                    &noBrowser,
		NotificationsEnabled:         &notificationsEnabled,
		AutostartVideo:               &autostartVideo,
		ShowStudioAsText:             &showStudioAsText,
		AutostartVideoOnPlaySelected: &autostartVideoOnPlaySelected,
		ContinuePlaylistDefault:      &continuePlaylistDefault,
		CSS:                          &css,
		CSSEnabled:                   &cssEnabled,
		Language:                     &language,

		ImageLightbox: &imageLightboxOptions,

		// FIXME - see above
		DisabledDropdownCreate: disableDropdownCreate,
		DisableDropdownCreate:  disableDropdownCreate,

		HandyKey:        &handyKey,
		FunscriptOffset: &scriptOffset,
	}
}

func makeConfigDLNAResult() *ConfigDLNAResult {
	config := config.GetInstance()

	return &ConfigDLNAResult{
		ServerName:     config.GetDLNAServerName(),
		Enabled:        config.GetDLNADefaultEnabled(),
		WhitelistedIPs: config.GetDLNADefaultIPWhitelist(),
		Interfaces:     config.GetDLNAInterfaces(),
	}
}

func makeConfigScrapingResult() *ConfigScrapingResult {
	config := config.GetInstance()

	scraperUserAgent := config.GetScraperUserAgent()
	scraperCDPPath := config.GetScraperCDPPath()

	return &ConfigScrapingResult{
		ScraperUserAgent:   &scraperUserAgent,
		ScraperCertCheck:   config.GetScraperCertCheck(),
		ScraperCDPPath:     &scraperCDPPath,
		ExcludeTagPatterns: config.GetScraperExcludeTagPatterns(),
	}
}

func makeConfigDefaultsResult() *ConfigDefaultSettingsResult {
	config := config.GetInstance()
	deleteFileDefault := config.GetDeleteFileDefault()
	deleteGeneratedDefault := config.GetDeleteGeneratedDefault()

	return &ConfigDefaultSettingsResult{
		Identify:        config.GetDefaultIdentifySettings(),
		Scan:            config.GetDefaultScanSettings(),
		AutoTag:         config.GetDefaultAutoTagSettings(),
		Generate:        config.GetDefaultGenerateSettings(),
		DeleteFile:      &deleteFileDefault,
		DeleteGenerated: &deleteGeneratedDefault,
	}
}

func makeConfigUIResult() map[string]interface{} {
	return config.GetInstance().GetUIConfiguration()
}

func (r *queryResolver) ValidateStashBoxCredentials(ctx context.Context, input config.StashBoxInput) (*StashBoxValidationResult, error) {
	client := stashbox.NewClient(models.StashBox{Endpoint: input.Endpoint, APIKey: input.APIKey}, r.txnManager)
	user, err := client.GetUser(ctx)

	valid := user != nil && user.Me != nil
	var status string
	if valid {
		status = fmt.Sprintf("Successfully authenticated as %s", user.Me.Name)
	} else {
		switch {
		case strings.Contains(strings.ToLower(err.Error()), "doctype"):
			// Index file returned rather than graphql
			status = "Invalid endpoint"
		case strings.Contains(err.Error(), "request failed"):
			status = "No response from server"
		case strings.HasPrefix(err.Error(), "invalid character") ||
			strings.HasPrefix(err.Error(), "illegal base64 data") ||
			err.Error() == "unexpected end of JSON input" ||
			err.Error() == "token contains an invalid number of segments":
			status = "Malformed API key."
		case err.Error() == "" || err.Error() == "signature is invalid":
			status = "Invalid or expired API key."
		default:
			status = fmt.Sprintf("Unknown error: %s", err)
		}
	}

	result := StashBoxValidationResult{
		Valid:  valid,
		Status: status,
	}

	return &result, nil
}
