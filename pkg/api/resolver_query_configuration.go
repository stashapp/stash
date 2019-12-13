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

func (r *queryResolver) Directories(ctx context.Context, path *string) ([]string, error) {
	var dirPath = ""
	if path != nil {
		dirPath = *path
	}
	return utils.ListDir(dirPath), nil
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

	return &models.ConfigGeneralResult{
		Stashes:                   config.GetStashPaths(),
		DatabasePath:              config.GetDatabasePath(),
		GeneratedPath:             config.GetGeneratedPath(),
		MaxTranscodeSize:          &maxTranscodeSize,
		MaxStreamingTranscodeSize: &maxStreamingTranscodeSize,
		Username:                  config.GetUsername(),
		Password:                  config.GetPasswordHash(),
		LogFile:                   &logFile,
		LogOut:                    config.GetLogOut(),
		LogLevel:                  config.GetLogLevel(),
		LogAccess:                 config.GetLogAccess(),
		Excludes:                  config.GetExcludes(),
	}
}

func makeConfigInterfaceResult() *models.ConfigInterfaceResult {
	soundOnPreview := config.GetSoundOnPreview()
	wallShowTitle := config.GetWallShowTitle()
	maximumLoopDuration := config.GetMaximumLoopDuration()
	autostartVideo := config.GetAutostartVideo()
	showStudioAsText := config.GetShowStudioAsText()
	css := config.GetCSS()
	cssEnabled := config.GetCSSEnabled()

	return &models.ConfigInterfaceResult{
		SoundOnPreview:      &soundOnPreview,
		WallShowTitle:       &wallShowTitle,
		MaximumLoopDuration: &maximumLoopDuration,
		AutostartVideo:      &autostartVideo,
		ShowStudioAsText:    &showStudioAsText,
		CSS:                 &css,
		CSSEnabled:          &cssEnabled,
	}
}
