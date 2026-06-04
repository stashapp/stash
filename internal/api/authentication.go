package api

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/signedurl"
)

const (
	tripwireActivatedErrMsg = "Stash is exposed to the public internet without authentication, and is not serving any more content to protect your privacy. " +
		"More information and fixes are available at https://discourse.stashapp.cc/t/-/1658"

	externalAccessErrMsg = "You have attempted to access Stash over the internet, and authentication is not enabled. " +
		"This is extremely dangerous! The whole world can see your your stash page and browse your files! " +
		"Stash is not answering any other requests to protect your privacy. " +
		"Please read the log entry or visit https://discourse.stashapp.cc/t/-/1658"
)

func allowUnauthenticated(r *http.Request) bool {
	// #2715 - allow access to UI files
	return strings.HasPrefix(r.URL.Path, loginEndpoint) || r.URL.Path == logoutEndpoint || r.URL.Path == "/css" || strings.HasPrefix(r.URL.Path, "/assets")
}

// authenticateSignedRequest checks if the request is a valid signed media request.
// Returns the matched username and true if valid, or empty string and false otherwise.
func authenticateSignedRequest(r *http.Request) (string, bool) {
	// Only apply to scene stream paths (used by AirPlay/Chromecast devices that can't pass cookies)
	if !strings.HasPrefix(r.URL.Path, "/scene/") {
		return "", false
	}

	c := config.GetInstance()

	// Signed URLs are only relevant when credentials are configured
	if !c.HasCredentials() {
		return "", false
	}

	// Check for signed URL parameters
	q := r.URL.Query()
	if q.Get(signedurl.CIDParam) == "" || q.Get(signedurl.ExpiresParam) == "" || q.Get(signedurl.SigParam) == "" {
		return "", false
	}

	// Extract the credential ID and look up the user's signing key.
	// We need the key before we can verify the signature, since in a
	// multi-user setup each user could have their own signing key.
	cid := q.Get(signedurl.CIDParam)
	username, secret, found := resolveCredentialID(c, cid)
	if !found {
		logger.Warnf("signed URL credential ID mismatch")
		return "", false
	}

	// Verify the signature using the user's signing key
	if _, err := signedurl.VerifyURL(r.URL.Path, q, secret); err != nil {
		logger.Warnf("signed URL verification failed: %v", err)
		return "", false
	}

	return username, true
}

func authenticateHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := config.GetInstance()

			// error if external access tripwire activated
			if accessErr := session.CheckExternalAccessTripwire(c); accessErr != nil {
				http.Error(w, tripwireActivatedErrMsg, http.StatusForbidden)
				return
			}

			r = session.SetLocalRequest(r)

			// Check for signed media requests
			if username, ok := authenticateSignedRequest(r); ok {
				ctx := r.Context()
				ctx = session.SetCurrentUserID(ctx, username)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			userID, err := manager.GetInstance().SessionStore.Authenticate(w, r)
			if err != nil {
				if !errors.Is(err, session.ErrUnauthorized) {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// unauthorized error
				w.Header().Add("WWW-Authenticate", "FormBased")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if err := session.CheckAllowPublicWithoutAuth(c, r); err != nil {
				var accessErr session.ExternalAccessError
				if errors.As(err, &accessErr) {
					session.LogExternalAccessError(accessErr)

					err := c.ActivatePublicAccessTripwire(net.IP(accessErr).String())
					if err != nil {
						logger.Errorf("Error activating public access tripwire: %v", err)
					}

					http.Error(w, externalAccessErrMsg, http.StatusForbidden)
				} else {
					logger.Errorf("Error checking external access security: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			ctx := r.Context()

			if c.HasCredentials() {
				// authentication is required
				if userID == "" && !allowUnauthenticated(r) {
					// if graphql or a non-webpage was requested, we just return a forbidden error
					ext := path.Ext(r.URL.Path)
					if r.URL.Path == gqlEndpoint || (ext != "" && ext != ".html") {
						w.Header().Add("WWW-Authenticate", "FormBased")
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

					prefix := getProxyPrefix(r)

					// otherwise redirect to the login page
					returnURL := url.URL{
						Path:     prefix + r.URL.Path,
						RawQuery: r.URL.RawQuery,
					}
					q := make(url.Values)
					q.Set(returnURLParam, returnURL.String())
					u := url.URL{
						Path:     prefix + loginEndpoint,
						RawQuery: q.Encode(),
					}
					http.Redirect(w, r, u.String(), http.StatusFound)
					return
				}
			}

			ctx = session.SetCurrentUserID(ctx, userID)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
