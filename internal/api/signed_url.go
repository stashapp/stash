package api

import (
	"context"
	"net/url"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/signedurl"
)

// userSigningKey returns the HMAC signing key for a given user.
func userSigningKey(c *config.Config, _ string) []byte {
	return c.GetJWTSignKey()
}

// signedParams generates signed URL query parameters for the given path prefix and user.
func signedParams(c *config.Config, username string, prefix string) url.Values {
	secret := userSigningKey(c, username)
	cid := signedurl.GenerateCredentialID(secret, username)
	expires := time.Now().Add(c.GetSignedURLExpiry())
	return signedurl.SignPrefix(prefix, secret, cid, expires)
}

// signingUser returns the username that media URLs should be signed for.
// Returns false if URLs do not need to be signed - when login is not required,
// or when the current user is the guest user, who can access media without credentials.
func signingUser(ctx context.Context) (string, bool) {
	if manager.GetInstance().UserService.IsSingleUserMode() {
		return "", false
	}

	u := session.GetCurrentUser(ctx)
	if u == nil || u.ID == 0 {
		return "", false
	}

	return u.Username, true
}

// signURL adds signed URL parameters for the current user to the provided URL, if required.
func signURL(ctx context.Context, c *config.Config, u *url.URL, prefix string) {
	if username, ok := signingUser(ctx); ok {
		u.RawQuery = signedParams(c, username, prefix).Encode()
	}
}

// resolveCredentialID maps a credential ID back to a user and their signing key.
// Locked users are not matched.
func resolveCredentialID(c *config.Config, users []*models.User, cid string) (*models.User, []byte, bool) {
	for _, u := range users {
		// users without roles are also considered locked
		if u.Locked || len(u.Roles) == 0 {
			continue
		}

		secret := userSigningKey(c, u.Username)
		if signedurl.GenerateCredentialID(secret, u.Username) == cid {
			return u, secret, true
		}
	}

	return nil, nil, false
}
