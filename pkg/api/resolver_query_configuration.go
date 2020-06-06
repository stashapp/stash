package api

import (
	"context"

	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *queryResolver) Configuration(ctx context.Context) (*models.ConfigResult, error) {
	return makeConfigResult(), nil
}

func (r *queryResolver) Directory(ctx context.Context, path *string) (*models.Directory, error) {
	var dirPath = ""
	if path != nil {
		dirPath = *path
	}
	currentDir := utils.GetDir(dirPath)

	return &models.Directory{
		Path:        currentDir,
		Parent:      utils.GetParent(currentDir),
		Directories: utils.ListDir(currentDir),
	}, nil
}

func makeConfigResult() *models.ConfigResult {
	return &models.ConfigResult{
		General:   makeConfigGeneralResult(),
		Interface: makeConfigInterfaceResult(),
	}
}

func makeConfigGeneralResult() *models.ConfigGeneralResult {
	logFile := config.GetLogFile()

	maxTranscodeSize := config.GetMaxTranscodeSize()
	maxStreamingTranscodeSize := config.GetMaxStreamingTranscodeSize()

	scraperUserAgent := config.GetScraperUserAgent()

	return &models.ConfigGeneralResult{
		Stashes:                   config.GetStashPaths(),
		DatabasePath:              config.GetDatabasePath(),
		GeneratedPath:             config.GetGeneratedPath(),
		CachePath:                 config.GetCachePath(),
		MaxTranscodeSize:          &maxTranscodeSize,
		MaxStreamingTranscodeSize: &maxStreamingTranscodeSize,
		ForceMkv:                  config.GetForceMKV(),
		ForceHevc:                 config.GetForceHEVC(),
		Username:                  config.GetUsername(),
		Password:                  config.GetPasswordHash(),
		MaxSessionAge:             config.GetMaxSessionAge(),
		LogFile:                   &logFile,
		LogOut:                    config.GetLogOut(),
		LogLevel:                  config.GetLogLevel(),
		LogAccess:                 config.GetLogAccess(),
		Excludes:                  config.GetExcludes(),
		ScraperUserAgent:          &scraperUserAgent,
	}
}

func makeConfigInterfaceResult() *models.ConfigInterfaceResult {
	soundOnPreview := config.GetSoundOnPreview()
	wallShowTitle := config.GetWallShowTitle()
	wallPlayback := config.GetWallPlayback()
	maximumLoopDuration := config.GetMaximumLoopDuration()
	autostartVideo := config.GetAutostartVideo()
	showStudioAsText := config.GetShowStudioAsText()
	css := config.GetCSS()
	cssEnabled := config.GetCSSEnabled()
	language := config.GetLanguage()

	return &models.ConfigInterfaceResult{
		SoundOnPreview:      &soundOnPreview,
		WallShowTitle:       &wallShowTitle,
		WallPlayback:        &wallPlayback,
		MaximumLoopDuration: &maximumLoopDuration,
		AutostartVideo:      &autostartVideo,
		ShowStudioAsText:    &showStudioAsText,
		CSS:                 &css,
		CSSEnabled:          &cssEnabled,
		Language:            &language,
	}
}
