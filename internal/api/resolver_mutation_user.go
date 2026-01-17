package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
)

func (r *mutationResolver) UserCreate(ctx context.Context, input UserCreateInput) (*models.User, error) {
	err := r.userService.CreateUser(ctx, models.User{
		Username: input.Name,
		Roles:    models.Roles(input.Roles),
	}, input.Password)
	if err != nil {
		return nil, err
	}

	return r.userService.GetUser(ctx, input.Name)
}

func (r *mutationResolver) UserUpdate(ctx context.Context, input UserUpdateInput) (*models.User, error) {
	err := r.userService.UpdateUser(ctx, input.ExistingName, models.User{
		Username: input.Name,
		Roles:    models.Roles(input.Roles),
	})
	if err != nil {
		return nil, err
	}

	return r.userService.GetUser(ctx, input.Name)
}

func (r *mutationResolver) UserDestroy(ctx context.Context, input UserDestroyInput) (bool, error) {
	err := r.userService.DeleteUser(ctx, input.Name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input UserChangePasswordInput) (bool, error) {
	// get current user
	u := session.GetCurrentUser(ctx)

	err := r.userService.ChangePassword(ctx, u.Username, input.ExistingPassword, input.NewPassword)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ChangeUserPassword(ctx context.Context, input ChangeUserPasswordInput) (bool, error) {
	err := r.userService.ChangeUserPassword(ctx, input.Name, input.NewPassword)
	if err != nil {
		return false, err
	}

	return true, nil
}
