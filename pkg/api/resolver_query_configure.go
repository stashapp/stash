package api

import (
	"context"
	"fmt"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *queryResolver) ConfigureGeneral(ctx context.Context, input *models.ConfigGeneralInput) (models.ConfigGeneralResult, error) {
	if input == nil {
		return makeConfigGeneralResult(), fmt.Errorf("nil input")
	}

	if len(input.Stashes) > 0 {
		for _, stashPath := range input.Stashes {
			exists, err := utils.DirExists(stashPath)
			if !exists {
				return makeConfigGeneralResult(), err
			}
		}
		config.Set(config.Stash, input.Stashes)
	}

	if err := config.Write(); err != nil {
		return makeConfigGeneralResult(), err
	}

	return makeConfigGeneralResult(), nil
}

func makeConfigGeneralResult() models.ConfigGeneralResult {
	return models.ConfigGeneralResult{
		Stashes: config.GetStashPaths(),
	}
}
