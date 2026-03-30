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
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/user"
)

// const (
// 	tripwireActivatedErrMsg = "Stash is exposed to the public internet without authentication, and is not serving any more content to protect your privacy. " +
// 		"More information and fixes are available at https://discourse.stashapp.cc/t/-/1658"

// 	externalAccessErrMsg = "You have attempted to access Stash over the internet, and authentication is not enabled. " +
// 		"This is extremely dangerous! The whole world can see your your stash page and browse your files! " +
// 		"Stash is not answering any other requests to protect your privacy. " +
// 		"Please read the log entry or visit https://discourse.stashapp.cc/t/-/1658"
// )

func allowUnauthenticated(r *http.Request) bool {
	// #2715 - allow access to UI files
	return strings.HasPrefix(r.URL.Path, loginEndpoint) || r.URL.Path == logoutEndpoint || r.URL.Path == "/css" || strings.HasPrefix(r.URL.Path, "/assets")
}

func handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	// if graphql or a non-webpage was requested, we just return a forbidden error
	ext := path.Ext(r.URL.Path)
	if r.URL.Path == gqlEndpoint || (ext != "" && ext != ".html") {
		w.Header().Add("WWW-Authenticate", "FormBased")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	prefix := getProxyPrefix(r)

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
}

type UserAuthenticator interface {
	GetGuestUser(ctx context.Context) *models.User

	IsSingleUserMode() bool
	GetSingleUser(ctx context.Context) (*models.User, error)

	AuthenticateByAPIKey(ctx context.Context, apiKey string) (*models.User, error)
	AuthenticateSession(ctx context.Context, username string, loginTime time.Time) (*models.User, error)
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

type PublicIPWhitelistGetter interface {
	GetPublicWhitelist() ([]net.IPNet, []net.IP)
}

func matchIPWhitelist(g PublicIPWhitelistGetter, requestIP net.IP) bool {
	nets, addrs := g.GetPublicWhitelist()

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

type authenticationConfig interface {
	IsNewSystem() bool
	GetNewSystemCredentials() (string, string, time.Time)

	GetPublicAccess() bool
	PublicIPWhitelistGetter
}

func authenticateHandler(txnMgr models.TxnManager, g UserAuthenticator, cfg authenticationConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			singleUserMode := g.IsSingleUserMode()
			publicAccess := cfg.GetPublicAccess()

			if singleUserMode && publicAccess {
				// mis-configuration - cannot enable single user mode and public access at the same time
				httpError(w, r, "Server misconfiguration: single user mode and public access cannot both be enabled", http.StatusServiceUnavailable)
				return
			}

			if !publicAccess {
				requestIP, err := getRequestIPFromCtx(ctx)
				if err != nil {
					logger.Errorf("error getting request IP: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if !isLocalIP(requestIP) && !matchIPWhitelist(cfg, requestIP) {
					httpError(w, r, "Access denied: Stash is not configured to allow access from external IPs", http.StatusForbidden)
					return
				}
			}

			if cfg.IsNewSystem() {
				expectedUsername, _, startupTime := cfg.GetNewSystemCredentials()

				session, _ := manager.GetInstance().SessionStore.GetSession(w, r)
				if !allowUnauthenticated(r) && expectedUsername != "" && (session == nil || session.UserID != expectedUsername || session.LoginTime.Before(startupTime)) {
					handleUnauthorized(w, r)
					return
				}

				next.ServeHTTP(w, r)
				return
			}

			var defaultUser *models.User
			if singleUserMode {
				if err := txn.WithReadTxn(ctx, txnMgr, func(ctx context.Context) error {
					var err error
					defaultUser, err = g.GetSingleUser(ctx)
					if err != nil {
						return fmt.Errorf("error retrieving default user: %w", err)
					}

					return nil
				}); err != nil {
					logger.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			// // error if external access tripwire activated
			// if accessErr := session.CheckExternalAccessTripwire(loginRequired, c); accessErr != nil {
			// 	httpError(w, r, tripwireActivatedErrMsg, http.StatusForbidden)
			// 	return
			// }

			r = session.SetLocalRequest(r)

			var u *models.User
			var err error

			if defaultUser != nil {
				u = defaultUser
				u.Roles = models.Roles{models.RoleEnumAdmin}
			}

			if !singleUserMode {
				// try to authenticate using api key first
				apiKey := session.GetRequestApiKey(r)
				if apiKey != "" {
					requestIP, err := keyByIP(r)
					if err != nil {
						logger.Errorf("Error getting IP for API key rate limiting: %v", err)
						httpError(w, r, "An unexpected error occurred. See logs", http.StatusInternalServerError)
						return
					}

					canIncr, err := loginRateLimiter.canIncrement(requestIP)
					if !canIncr {
						loginRateLimiter.OnLimit(w, r, requestIP)
						httpError(w, r, "Too many attempts using bad API key. Please try again later.", http.StatusTooManyRequests)
						return
					}

					if err := txn.WithReadTxn(ctx, txnMgr, func(ctx context.Context) error {
						u, err = g.AuthenticateByAPIKey(ctx, apiKey)
						return err
					}); err != nil {
						if errors.Is(err, user.ErrInternalError) {
							logger.Errorf("error authenticating by API key: %v", err)
							httpError(w, r, "internal server error", http.StatusInternalServerError)
							return
						}

						// must be authentication failure
						// rate limit failed API key attempts to prevent brute forcing
						_ = loginRateLimiter.Increment(requestIP)

						// don't forward to login page if api key auth fails
						httpError(w, r, "Invalid API key", http.StatusUnauthorized)
						return
					}
				} else {
					session, getErr := manager.GetInstance().SessionStore.GetSession(w, r)
					if getErr != nil {
						logger.Errorf("error getting session user ID: %v", getErr)
						httpError(w, r, "internal server error", http.StatusInternalServerError)
						return
					}

					if session != nil {
						if err := txn.WithReadTxn(ctx, txnMgr, func(ctx context.Context) error {
							u, err = g.AuthenticateSession(ctx, session.UserID, session.LoginTime)
							return err
						}); err != nil {
							if errors.Is(err, user.ErrInternalError) {
								logger.Errorf("error authenticating user by ID: %v", err)
								httpError(w, r, "internal server error", http.StatusInternalServerError)
								return
							}
						}
					}
				}

				if err != nil {
					if errors.Is(err, user.ErrInternalError) {
						httpError(w, r, "internal server error", http.StatusInternalServerError)
						return
					}

					// fall through and treat as unauthenticated
				}
			}

			// TODO remove this in favour of ip whitelist
			// if err := session.CheckAllowPublicWithoutAuth(loginRequired, c, r); err != nil {
			// 	var accessErr session.ExternalAccessError
			// 	if errors.As(err, &accessErr) {
			// 		session.LogExternalAccessError(accessErr)

			// 		err := c.ActivatePublicAccessTripwire(net.IP(accessErr).String())
			// 		if err != nil {
			// 			logger.Errorf("Error activating public access tripwire: %v", err)
			// 		}

			// 		httpError(w, r, externalAccessErrMsg, http.StatusForbidden)
			// 	} else {
			// 		logger.Errorf("Error checking external access security: %v", err)
			// 		w.WriteHeader(http.StatusInternalServerError)
			// 	}
			// 	return
			// }

			if u == nil && !singleUserMode {
				u = g.GetGuestUser(ctx)
			}

			// authentication is required
			if u == nil && !allowUnauthenticated(r) {
				handleUnauthorized(w, r)
				return
			}

			if u != nil {
				// set the user object in the context
				ctx = session.SetCurrentUser(ctx, *u)
			}

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
