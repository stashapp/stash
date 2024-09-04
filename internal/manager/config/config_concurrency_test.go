package config

import (
	"sync"
	"testing"
	"time"
)

// should be run with -race
func TestConcurrentConfigAccess(t *testing.T) {
	i := InitializeEmpty()

	const workers = 8
	const loops = 200
	var wg sync.WaitGroup
	for k := 0; k < workers; k++ {
		wg.Add(1)
		go func(wk int) {
			for l := 0; l < loops; l++ {
				start := time.Now()
				if err := i.SetInitialConfig(); err != nil {
					t.Errorf("Failure setting initial configuration in worker %v iteration %v: %v", wk, l, err)
				}

				i.HasCredentials()
				i.ValidateCredentials("", "")
				i.GetConfigFile()
				i.GetConfigPath()
				i.GetDefaultDatabaseFilePath()
				i.SetInterface(BackupDirectoryPath, i.GetBackupDirectoryPath())
				i.GetStashPaths()
				_ = i.ValidateStashBoxes(nil)
				_ = i.Validate()
				_ = i.ActivatePublicAccessTripwire("")
				i.SetInterface(Cache, i.GetCachePath())
				i.SetInterface(Generated, i.GetGeneratedPath())
				i.SetInterface(Metadata, i.GetMetadataPath())
				i.SetInterface(Database, i.GetDatabasePath())

				// these must be set as strings since the original values are also strings
				// setting them as []byte will cause the returned string to be corrupted
				i.SetInterface(JWTSignKey, string(i.GetJWTSignKey()))
				i.SetInterface(SessionStoreKey, string(i.GetSessionStoreKey()))

				i.GetDefaultScrapersPath()
				i.SetInterface(Exclude, i.GetExcludes())
				i.SetInterface(ImageExclude, i.GetImageExcludes())
				i.SetInterface(VideoExtensions, i.GetVideoExtensions())
				i.SetInterface(ImageExtensions, i.GetImageExtensions())
				i.SetInterface(GalleryExtensions, i.GetGalleryExtensions())
				i.SetInterface(CreateGalleriesFromFolders, i.GetCreateGalleriesFromFolders())
				i.SetInterface(Language, i.GetLanguage())
				i.SetInterface(VideoFileNamingAlgorithm, i.GetVideoFileNamingAlgorithm())
				i.SetInterface(ScrapersPath, i.GetScrapersPath())
				i.SetInterface(ScraperUserAgent, i.GetScraperUserAgent())
				i.SetInterface(ScraperCDPPath, i.GetScraperCDPPath())
				i.SetInterface(ScraperCertCheck, i.GetScraperCertCheck())
				i.SetInterface(ScraperExcludeTagPatterns, i.GetScraperExcludeTagPatterns())
				i.SetInterface(StashBoxes, i.GetStashBoxes())
				i.GetDefaultPluginsPath()
				i.SetInterface(PluginsPath, i.GetPluginsPath())
				i.SetInterface(Host, i.GetHost())
				i.SetInterface(Port, i.GetPort())
				i.SetInterface(ExternalHost, i.GetExternalHost())
				i.SetInterface(PreviewSegmentDuration, i.GetPreviewSegmentDuration())
				i.SetInterface(ParallelTasks, i.GetParallelTasks())
				i.SetInterface(ParallelTasks, i.GetParallelTasksWithAutoDetection())
				i.SetInterface(PreviewAudio, i.GetPreviewAudio())
				i.SetInterface(PreviewSegments, i.GetPreviewSegments())
				i.SetInterface(PreviewExcludeStart, i.GetPreviewExcludeStart())
				i.SetInterface(PreviewExcludeEnd, i.GetPreviewExcludeEnd())
				i.SetInterface(PreviewPreset, i.GetPreviewPreset())
				i.SetInterface(MaxTranscodeSize, i.GetMaxTranscodeSize())
				i.SetInterface(MaxStreamingTranscodeSize, i.GetMaxStreamingTranscodeSize())
				i.SetInterface(ApiKey, i.GetAPIKey())
				i.SetInterface(Username, i.GetUsername())
				i.SetInterface(Password, i.GetPasswordHash())
				i.GetCredentials()
				i.SetInterface(MaxSessionAge, i.GetMaxSessionAge())
				i.SetInterface(CustomServedFolders, i.GetCustomServedFolders())
				i.SetInterface(LegacyCustomUILocation, i.GetUILocation())
				i.SetInterface(MenuItems, i.GetMenuItems())
				i.SetInterface(SoundOnPreview, i.GetSoundOnPreview())
				i.SetInterface(WallShowTitle, i.GetWallShowTitle())
				i.SetInterface(CustomPerformerImageLocation, i.GetCustomPerformerImageLocation())
				i.SetInterface(WallPlayback, i.GetWallPlayback())
				i.SetInterface(MaximumLoopDuration, i.GetMaximumLoopDuration())
				i.SetInterface(AutostartVideo, i.GetAutostartVideo())
				i.SetInterface(ShowStudioAsText, i.GetShowStudioAsText())
				i.SetInterface(legacyImageLightboxSlideshowDelay, *i.GetImageLightboxOptions().SlideshowDelay)
				i.SetInterface(ImageLightboxSlideshowDelay, *i.GetImageLightboxOptions().SlideshowDelay)
				i.GetCSSPath()
				i.GetCSS()
				i.GetJavascriptPath()
				i.GetJavascript()
				i.GetCustomLocalesPath()
				i.GetCustomLocales()
				i.SetInterface(CSSEnabled, i.GetCSSEnabled())
				i.SetInterface(CSSEnabled, i.GetCustomLocalesEnabled())
				i.SetInterface(HandyKey, i.GetHandyKey())
				i.SetInterface(UseStashHostedFunscript, i.GetUseStashHostedFunscript())
				i.SetInterface(DLNAServerName, i.GetDLNAServerName())
				i.SetInterface(DLNADefaultEnabled, i.GetDLNADefaultEnabled())
				i.SetInterface(DLNADefaultIPWhitelist, i.GetDLNADefaultIPWhitelist())
				i.SetInterface(DLNAInterfaces, i.GetDLNAInterfaces())
				i.SetInterface(DLNAPort, i.GetDLNAPort())
				i.SetInterface(LogFile, i.GetLogFile())
				i.SetInterface(LogOut, i.GetLogOut())
				i.SetInterface(LogLevel, i.GetLogLevel())
				i.SetInterface(LogAccess, i.GetLogAccess())
				i.SetInterface(MaxUploadSize, i.GetMaxUploadSize())
				i.SetInterface(FunscriptOffset, i.GetFunscriptOffset())
				i.SetInterface(DefaultIdentifySettings, i.GetDefaultIdentifySettings())
				i.SetInterface(DeleteGeneratedDefault, i.GetDeleteGeneratedDefault())
				i.SetInterface(DeleteFileDefault, i.GetDeleteFileDefault())
				i.SetInterface(dangerousAllowPublicWithoutAuth, i.GetDangerousAllowPublicWithoutAuth())
				i.SetInterface(SecurityTripwireAccessedFromPublicInternet, i.GetSecurityTripwireAccessedFromPublicInternet())
				i.SetInterface(DisableDropdownCreatePerformer, i.GetDisableDropdownCreate().Performer)
				i.SetInterface(DisableDropdownCreateStudio, i.GetDisableDropdownCreate().Studio)
				i.SetInterface(DisableDropdownCreateTag, i.GetDisableDropdownCreate().Tag)
				i.SetInterface(DisableDropdownCreateMovie, i.GetDisableDropdownCreate().Movie)
				i.SetInterface(AutostartVideoOnPlaySelected, i.GetAutostartVideoOnPlaySelected())
				i.SetInterface(ContinuePlaylistDefault, i.GetContinuePlaylistDefault())
				i.SetInterface(PythonPath, i.GetPythonPath())
				t.Logf("Worker %v iteration %v took %v", wk, l, time.Since(start))
			}
			wg.Done()
		}(k)
	}

	wg.Wait()
}
