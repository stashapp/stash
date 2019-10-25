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
	return &models.ConfigGeneralResult{
		Stashes:       config.GetStashPaths(),
		DatabasePath:  config.GetDatabasePath(),
		GeneratedPath: config.GetGeneratedPath(),
		Username:      config.GetUsername(),
		Password:      config.GetPasswordHash(),
		LogFile:       &logFile,
		LogOut:        config.GetLogOut(),
		LogLevel:      config.GetLogLevel(),
		LogAccess:     config.GetLogAccess(),
	}
}

func makeConfigInterfaceResult() *models.ConfigInterfaceResult {
	css := config.GetCSS()
	cssEnabled := config.GetCSSEnabled()
	return &models.ConfigInterfaceResult{
		CSS:        &css,
		CSSEnabled: &cssEnabled,
	}
}
