package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/user"
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

type UserAuthenticator interface {
	LoginRequired(ctx context.Context) (bool, error)
	AuthenticateByAPIKey(ctx context.Context, apiKey string) (*models.User, error)
	AuthenticateSession(ctx context.Context, username string, loginTime time.Time) (*models.User, error)
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

func authenticateHandler(txnMgr models.TxnManager, g UserAuthenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := config.GetInstance()
			ctx := r.Context()

			// TODO - check ip whitelist here

			var loginRequired bool
			if err := txn.WithReadTxn(ctx, txnMgr, func(ctx context.Context) error {
				var err error
				loginRequired, err = g.LoginRequired(ctx)
				return err
			}); err != nil {
				logger.Errorf("Error checking if login is required: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// error if external access tripwire activated
			if accessErr := session.CheckExternalAccessTripwire(loginRequired, c); accessErr != nil {
				http.Error(w, tripwireActivatedErrMsg, http.StatusForbidden)
				return
			}

			r = session.SetLocalRequest(r)

			// try to authenticate using api key first
			var u *models.User
			var err error

			apiKey := session.GetRequestApiKey(r)
			if apiKey != "" {
				if err := txn.WithReadTxn(ctx, txnMgr, func(ctx context.Context) error {
					u, err = g.AuthenticateByAPIKey(ctx, apiKey)
					return err
				}); err != nil {
					if errors.Is(err, user.ErrInternalError) {
						logger.Errorf("error authenticating by API key: %v", err)
						http.Error(w, "internal server error", http.StatusInternalServerError)
						return
					}
				}
			} else {
				session, getErr := manager.GetInstance().SessionStore.GetSessionUserID(w, r)
				if getErr != nil {
					logger.Errorf("error getting session user ID: %v", getErr)
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}

				if session != nil {
					if err := txn.WithReadTxn(ctx, txnMgr, func(ctx context.Context) error {
						u, err = g.AuthenticateSession(ctx, session.UserID, session.LoginTime)
						return err
					}); err != nil {
						if errors.Is(err, user.ErrInternalError) {
							logger.Errorf("error authenticating user by ID: %v", err)
							http.Error(w, "internal server error", http.StatusInternalServerError)
							return
						}
					}
				}
			}

			if err != nil {
				if errors.Is(err, user.ErrInternalError) {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}

				// fall through and treat as unauthenticated
			}

			// TODO remove this in favour of ip whitelist
			if err := session.CheckAllowPublicWithoutAuth(loginRequired, c, r); err != nil {
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

			if loginRequired {
				// authentication is required
				if u == nil && !allowUnauthenticated(r) {
					handleUnauthorized(w, r)
					return
				}
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
