package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
)

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		users, err = r.userService.AllUsers(ctx)
		return err
	}); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	// get current user
	return session.GetCurrentUser(ctx), nil
}
