package api

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
)

func (r *mutationResolver) UserCreate(ctx context.Context, input UserCreateInput) (user *models.User, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		err := r.userService.CreateUser(ctx, models.User{
			Username: input.Name,
			Roles:    models.Roles(input.Roles),
		}, input.Password)
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}

		user, err = r.userService.GetUser(ctx, input.Name)
		if err != nil {
			return fmt.Errorf("error getting user after creation: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) UserUpdate(ctx context.Context, input UserUpdateInput) (user *models.User, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		err := r.userService.UpdateUser(ctx, input.ExistingName, models.User{
			Username: input.Name,
			Roles:    models.Roles(input.Roles),
		})
		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}

		user, err = r.userService.GetUser(ctx, input.Name)
		if err != nil {
			return fmt.Errorf("error getting user after update: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mutationResolver) UserDestroy(ctx context.Context, input UserDestroyInput) (bool, error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		err := r.userService.DeleteUser(ctx, input.Name)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input UserChangePasswordInput) (bool, error) {
	// get current user
	u := session.GetCurrentUser(ctx)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.userService.ChangePassword(ctx, u.Username, input.ExistingPassword, input.NewPassword)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ChangeUserPassword(ctx context.Context, input ChangeUserPasswordInput) (bool, error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		return r.userService.ChangeUserPassword(ctx, input.Name, input.NewPassword)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) GenerateAPIKey(ctx context.Context, input GenerateAPIKeyInput) (string, error) {
	u := session.GetCurrentUser(ctx)

	if u == nil {
		return "", fmt.Errorf("no current user in context")
	}

	if input.Clear != nil && *input.Clear {
		err := r.withTxn(ctx, func(ctx context.Context) error {
			return r.userService.ClearAPIKey(ctx, u.Username)
		})
		return "", err
	}

	var newAPIKey string
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var err error
		newAPIKey, err = r.userService.GenerateAPIKey(ctx, u.Username)
		return err
	}); err != nil {
		return "", err
	}

	return newAPIKey, nil
}
