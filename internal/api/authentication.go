package api

import (
	"errors"
	"fmt"
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

func httpError(w http.ResponseWriter, r *http.Request, text string, status int) {
	// if request accepts json, return json error response
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprintf(w, `{"error": "%s"}`, text)
	} else {
		http.Error(w, text, status)
	}
}

func authenticateHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := config.GetInstance()

			r = session.SetLocalRequest(r)
			ctx := r.Context()

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

			// reject connections from the public internet if authentication is not configured
			// don't apply for new systems
			if !c.IsNewSystem() && !c.HasCredentials() {
				requestIP, err := getRequestIPFromCtx(ctx)
				if err != nil {
					logger.Errorf("error getting request IP: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if err := checkAllowPublicWithoutAuth(c, requestIP); err != nil {
					httpError(w, r, "Access denied: Stash cannot be accessed from public IPs when authentication is not configured", http.StatusForbidden)
					return
				}
			}

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

func checkAllowPublicWithoutAuth(c *config.Config, requestIP net.IP) error {
	// reject connections from the public internet if authentication is not configured
	// don't apply for new systems
	if c.IsNewSystem() || c.HasCredentials() {
		return nil
	}

	if !isLocalIP(requestIP) && !matchIPWhitelist(c, requestIP) {
		return fmt.Errorf("stash accessed from external IP %s", requestIP.String())
	}

	return nil
}

func matchIPWhitelist(c *config.Config, requestIP net.IP) bool {
	nets, addrs := c.GetPublicWhitelist()

	for _, addr := range addrs {
		if addr.Equal(requestIP) {
			return true
		}
	}
	for _, net := range nets {
		if net.Contains(requestIP) {
			return true
		}
	}
	return false
}
