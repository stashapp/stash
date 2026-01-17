package api

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
)

func HasRoleDirective(ctx context.Context, obj interface{}, next graphql.Resolver, role models.RoleEnum) (interface{}, error) {
	currentUser := session.GetCurrentUser(ctx)

	// if there is no current user, this is an anonymous request
	// we should not end up here unless there are no credentials required
	if currentUser == nil {
		return next(ctx)
	}

	if currentUser != nil && !currentUser.Roles.HasRole(role) {
		return nil, session.ErrUnauthorized
	}

	return next(ctx)
}

func IsUserOwnerDirective(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	currentUser := session.GetCurrentUser(ctx)

	// if there is no current user, this is an anonymous request
	// we should not end up here unless there are no credentials required
	if currentUser == nil {
		return next(ctx)
	}

	// get the user from the object
	userObj, ok := obj.(*models.User)
	if !ok {
		return nil, session.ErrUnauthorized
	}

	if currentUser.Username != userObj.Username {
		return nil, session.ErrUnauthorized
	}

	return next(ctx)
}
