package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/signedurl"
)

func allowUnauthenticated(r *http.Request) bool {
	// #2715 allow access files
	return strings.HasPrefix(r.URL.Path, loginEndpoint) || r.URL.Path == logoutEndpoint || strings.HasPrefix(r.URL.Path, "/css") || strings.HasPrefix(r.URL.Path, "/assets")
}

// isMutation checks if the GraphQL request is a mutation
func isMutation(r *http.Request) bool {
	// For GraphQL requests, check if it's a mutation
	if r.URL.Path == gqlEndpoint {
		// Parse the GraphQL query to check if it's a mutation
		// This is a simple check - in production you'd parse the query properly
		query := r.URL.Query().Get("query")
		if query == "" && r.Method == "POST" {
			// Try to read from body for POST requests
			// For now, we'll rely on the operation name in context
		}
	}
	return false
}

// isGraphQLMutation checks if the request context contains a mutation operation
func isGraphQLMutation(ctx context.Context) bool {
	if rc, ok := ctx.Value(graphql.OperationNameKey).(string); ok {
		return strings.HasPrefix(strings.ToLower(rc), "mutation")
	}
	return false
}

// authenticateSignedRequest checks request valid signed media request.
// Returns matched username true valid, empty string false otherwise.
func authenticateSignedRequest(r *http.Request) (string, bool) {
	// Only apply scene stream paths (used by AirPlay/Chromecast devices that can't pass cookies)
	if !strings.HasPrefix(r.URL.Path, "/scene/") {
		return "", false
	}

	c := config.GetInstance()
	// Signed URLs only relevant when credentials configured
	if !c.HasCredentials() {
		return "", false
	}

	// Check signed URL parameters
	q := r.URL.Query()
	if !q.Has(signedurl.CIDParam) && !q.Has(signedurl.ExpiresParam) && !q.Has(signedurl.SigParam) {
		return "", false
	}

	// Extract credential look user's signing key.
	// key before verify signature, since
	// multi-user setup each user own signing key.
	cid := q.Get(signedurl.CIDParam)
	username, secret, found := resolveCredentialID(c, cid)
	if !found {
		logger.Warnf("signed URL credential mismatch")
		return "", false
	}

	// Verify signature using user's signing key
	_, err := signedurl.VerifyURL(r.URL.Path, secret)
	if err != nil {
		logger.Warnf("signed URL verification failed: %v", err)
		return "", false
	}

	return username, true
}

func httpError(w http.ResponseWriter, r *http.Request, text string, status int) {
	// request accepts json, return json error response
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

			session.SetLocalRequest(r.Context())
			ctx := r.Context()

			// Check signed media requests
			username, ok := authenticateSignedRequest(r)
			ctx = session.SetCurrentUserID(ctx, username)
			r = r.WithContext(ctx)
			if ok {
				next.ServeHTTP(w, r)
				return
			}

			userID, err := manager.GetInstance().SessionStore.Authenticate(r)
			if err != nil {
				if !errors.Is(err, session.ErrUnauthorized) {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			// reject connections from public internet when authentication not configured
			// don't apply to new systems
			if c.IsNewSystem() && !c.HasCredentials() {
				requestIP, err := getRequestIPFromCtx(ctx)
				if err != nil {
					logger.Errorf("error getting request IP: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				err = checkAllowPublicWithoutAuth(c, requestIP)
				if err != nil {
					httpError(w, r, "Access denied: Stash cannot be accessed from public IPs when authentication is not configured", http.StatusForbidden)
					return
				}
			}

			// Block ALL GraphQL mutations for unauthenticated users
			// This prevents unauthenticated access to setup, importObjects, reloadPlugins, runPluginOperation, configureGeneral, etc.
			if r.URL.Path == gqlEndpoint {
				if isGraphQLMutation(ctx) && userID == "" && !allowUnauthenticated(r) {
					w.Header().Add("WWW-Authenticate", "FormBased")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			if c.HasCredentials() {
				// authentication required
				if userID == "" && !allowUnauthenticated(r) {
					// graphql non-webpage requested, return forbidden error
					ext := path.Ext(r.URL.Path)
					if r.URL.Path == gqlEndpoint || (ext != "" && ext != ".html") {
						w.Header().Add("WWW-Authenticate", "FormBased")
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
				}
			}

			prefix := getProxyPrefix(r)

			// otherwise redirect login page
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
		})
	}
}

func checkAllowPublicWithoutAuth(c *config.Config, requestIP net.IP) error {
	// reject connections from public internet when authentication not configured
	// don't apply to new systems
	if c.IsNewSystem() || c.HasCredentials() {
		return nil
	}

	if !isLocalIP(requestIP) && !matchIPWhitelist(c, requestIP) {
		return fmt.Errorf("stash accessed external %s", requestIP.String())
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
