package api

import (
	"net/url"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/signedurl"
)

// userSigningKey returns the HMAC signing key for a given user.
func userSigningKey(c *config.Config, _ string) []byte {
	return c.GetJWTSignKey()
}

// signedParams generates signed URL query parameters for the given path prefix and user.
func signedParams(c *config.Config, userID string, prefix string) url.Values {
	secret := userSigningKey(c, userID)
	cid := signedurl.GenerateCredentialID(secret, userID)
	expires := time.Now().Add(time.Duration(c.GetSignedURLExpiry()) * time.Second)
	return signedurl.SignPrefix(prefix, secret, cid, expires)
}

// resolveCredentialID maps a credential ID back to a username and their signing key.
func resolveCredentialID(c *config.Config, cid string) (string, []byte, bool) {
	username := c.GetUsername()
	secret := userSigningKey(c, username)
	if signedurl.GenerateCredentialID(secret, username) == cid {
		return username, secret, true
	}
	return "", nil, false
}
