package api

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/user"
)

func HasRoleDirective(ctx context.Context, obj interface{}, next graphql.Resolver, role user.RoleEnum) (interface{}, error) {
	currentUser := session.GetCurrentUser(ctx)

	// if there is no current user, this is an anonymous request
	// we should not end up here unless there are no credentials required
	if currentUser == nil {
		return next(ctx)
	}

	if currentUser != nil && !user.IsRole(currentUser.Roles, role) {
		return nil, session.ErrUnauthorized
	}

	return next(ctx)
}
