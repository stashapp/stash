package api

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
)

// getUIConfig returns the UI configuration of the current user.
// Users without their own UI configuration use the configuration in the config file.
func (r *Resolver) getUIConfig(ctx context.Context) (map[string]interface{}, error) {
	var ret map[string]interface{}
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.UserUIConfig.Get(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	if ret == nil {
		ret = config.GetInstance().GetUIConfiguration()
	}

	return ret, nil
}

// setUIConfig sets the UI configuration of the current user.
// Configuration set by an administrator is also written to the config file, where it is used
// by users without their own configuration, and by server-side settings such as activity tracking.
func (r *Resolver) setUIConfig(ctx context.Context, uiConfig map[string]interface{}) error {
	_, hasUser := models.UserIDFromContext(ctx)

	cur := session.GetCurrentUser(ctx)
	isAdmin := cur != nil && cur.Roles.HasRole(models.RoleEnumAdmin)

	if !hasUser && !isAdmin {
		return fmt.Errorf("setting UI configuration: no current user")
	}

	if hasUser {
		if err := r.withTxn(ctx, func(ctx context.Context) error {
			return r.repository.UserUIConfig.Set(ctx, uiConfig)
		}); err != nil {
			return err
		}
	}

	if isAdmin {
		c := config.GetInstance()
		c.SetUIConfiguration(uiConfig)
		if err := c.Write(); err != nil {
			return err
		}
	}

	return nil
}
