package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func (r *mutationResolver) RunOperation(ctx context.Context, operationID string, args []*models.OperationArgInput) (*string, error) {
	var resp common.PluginOutput
	err := plugin.RunOperationPlugin([]string{"TODO"}, common.PluginInput{}, &resp)
	return nil, err
}
