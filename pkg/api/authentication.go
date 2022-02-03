package api

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/session"
)

const loginEndPoint = "/login"

const (
	tripwireActivatedErrMsg = "Stash is exposed to the public internet without authentication, and is not serving any more content to protect your privacy. " +
		"More information and fixes are available at https://github.com/stashapp/stash/wiki/Authentication-Required-When-Accessing-Stash-From-the-Internet"

	externalAccessErrMsg = "You have attempted to access Stash over the internet, and authentication is not enabled. " +
		"This is extremely dangerous! The whole world can see your your stash page and browse your files! " +
		"Stash is not answering any other requests to protect your privacy. " +
		"Please read the log entry or visit https://github.com/stashapp/stash/wiki/Authentication-Required-When-Accessing-Stash-From-the-Internet"
)

func allowUnauthenticated(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, loginEndPoint) || r.URL.Path == "/css"
}

func authenticateHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := config.GetInstance()

			if !checkSecurityTripwireActivated(c, w) {
				return
			}

			userID, err := manager.GetInstance().SessionStore.Authenticate(w, r)
			if err != nil {
				if errors.Is(err, session.ErrUnauthorized) {
					w.WriteHeader(http.StatusInternalServerError)
					_, err = w.Write([]byte(err.Error()))
					if err != nil {
						logger.Error(err)
					}
					return
				}

				// unauthorized error
				w.Header().Add("WWW-Authenticate", `FormBased`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if err := session.CheckAllowPublicWithoutAuth(c, r); err != nil {
				var externalAccess session.ExternalAccessError
				switch {
				case errors.As(err, &externalAccess):
					securityActivateTripwireAccessedFromInternetWithoutAuth(c, externalAccess, w)
					return
				default:
					logger.Errorf("Error checking external access security: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			ctx := r.Context()

			if c.HasCredentials() {
				// authentication is required
				if userID == "" && !allowUnauthenticated(r) {
					// authentication was not received, redirect
					// if graphql was requested, we just return a forbidden error
					if r.URL.Path == "/graphql" {
						w.Header().Add("WWW-Authenticate", `FormBased`)
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

					prefix := getProxyPrefix(r.Header)

					// otherwise redirect to the login page
					u := url.URL{
						Path: prefix + "/login",
					}
					q := u.Query()
					q.Set(returnURLParam, prefix+r.URL.Path)
					u.RawQuery = q.Encode()
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

func checkSecurityTripwireActivated(c *config.Instance, w http.ResponseWriter) bool {
	if accessErr := session.CheckExternalAccessTripwire(c); accessErr != nil {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(tripwireActivatedErrMsg))
		if err != nil {
			logger.Error(err)
		}
		return false
	}

	return true
}

func securityActivateTripwireAccessedFromInternetWithoutAuth(c *config.Instance, accessErr session.ExternalAccessError, w http.ResponseWriter) {
	session.LogExternalAccessError(accessErr)

	err := c.ActivatePublicAccessTripwire(net.IP(accessErr).String())
	if err != nil {
		logger.Error(err)
	}

	w.WriteHeader(http.StatusForbidden)
	_, err = w.Write([]byte(externalAccessErrMsg))
	if err != nil {
		logger.Error(err)
	}
}
