package api

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
)

func (r *userResolver) Name(ctx context.Context, obj *models.User) (string, error) {
	return obj.Username, nil
}

func (r *userResolver) Roles(ctx context.Context, obj *models.User) ([]models.RoleEnum, error) {
	ret := make([]models.RoleEnum, len(obj.Roles))
	for i, role := range obj.Roles {
		ret[i] = models.RoleEnum(role)
	}
	return ret, nil
}

func (r *userResolver) APIKey(ctx context.Context, obj *models.User) (*string, error) {
	return nil, nil
}
