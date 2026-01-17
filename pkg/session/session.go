// Package session provides session authentication and management for the application.
package session

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type key int

const (
	contextUser key = iota
	contextVisitedPlugins
	contextLocalRequest
)

const (
	userIDKey             = "userID"
	visitedPluginHooksKey = "visitedPluginsHooks"
)

const (
	ApiKeyHeader    = "ApiKey"
	ApiKeyParameter = "apikey"
)

const (
	cookieName      = "session"
	usernameFormKey = "username"
	passwordFormKey = "password"
)

type InvalidCredentialsError struct {
	Username string
}

func (e InvalidCredentialsError) Error() string {
	// don't leak the username
	return "invalid credentials"
}

var ErrUnauthorized = errors.New("unauthorized")

type Authenticator interface {
	LoginRequired(ctx context.Context) bool
	ValidateCredentials(ctx context.Context, username string, password string) error
}

type Store struct {
	sessionStore  *sessions.CookieStore
	authenticator Authenticator
	config        SessionConfig
}

func NewStore(c SessionConfig, a Authenticator) *Store {
	ret := &Store{
		sessionStore:  sessions.NewCookieStore(c.GetSessionStoreKey()),
		config:        c,
		authenticator: a,
	}

	ret.sessionStore.MaxAge(c.GetMaxSessionAge())
	ret.sessionStore.Options.SameSite = http.SameSiteLaxMode

	return ret
}

func (s *Store) LoginRequired(ctx context.Context) bool {
	return s.authenticator.LoginRequired(ctx)
}

func (s *Store) Login(w http.ResponseWriter, r *http.Request) error {
	// ignore error - we want a new session regardless
	newSession, _ := s.sessionStore.Get(r, cookieName)

	username := r.FormValue(usernameFormKey)
	password := r.FormValue(passwordFormKey)

	// authenticate the user
	err := s.authenticator.ValidateCredentials(r.Context(), username, password)
	if err != nil {
		return &InvalidCredentialsError{Username: username}
	}

	logger.Infof("User %s logged in", username)

	newSession.Values[userIDKey] = username

	err = newSession.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Logout(w http.ResponseWriter, r *http.Request) error {
	session, err := s.sessionStore.Get(r, cookieName)
	if err != nil {
		return err
	}

	userID, _ := session.Values[userIDKey].(string)

	delete(session.Values, userIDKey)
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	logger.Infof("User %s logged out", userID)

	return nil
}

func (s *Store) GetSessionUserID(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := s.sessionStore.Get(r, cookieName)
	// ignore errors and treat as an empty user id, so that we handle expired
	// cookie
	if err != nil {
		return "", nil
	}

	if !session.IsNew {
		val := session.Values[userIDKey]

		// refresh the cookie
		err = session.Save(r, w)
		if err != nil {
			return "", err
		}

		ret, _ := val.(string)

		return ret, nil
	}

	return "", nil
}

func SetCurrentUser(ctx context.Context, u models.User) context.Context {
	return context.WithValue(ctx, contextUser, u)
}

// GetCurrentUser gets the current user id from the provided context
func GetCurrentUser(ctx context.Context) *models.User {
	userCtxVal := ctx.Value(contextUser)
	if userCtxVal != nil {
		currentUser := userCtxVal.(models.User)
		return &currentUser
	}

	return nil
}

func (s *Store) Authenticate(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	c := s.config

	// translate api key into current user, if present
	apiKey := r.Header.Get(ApiKeyHeader)

	// try getting the api key as a query parameter
	if apiKey == "" {
		apiKey = r.URL.Query().Get(ApiKeyParameter)
	}

	// FIXME - handle this
	if apiKey != "" {
		// match against configured API and set userID to the
		// configured username. In future, we'll want to
		// get the username from the key.
		if c.GetAPIKey() != apiKey {
			return "", ErrUnauthorized
		}

		userID = c.GetUsername()
	} else {
		// handle session
		userID, err = s.GetSessionUserID(w, r)
	}

	if err != nil {
		return "", err
	}

	return
}
