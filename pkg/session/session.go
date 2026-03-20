// Package session provides session authentication and management for the application.
package session

import (
	"context"
	"errors"
	"net/http"
	"time"

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
	loginTimeKey          = "loginTime"
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
	LoginRequired(ctx context.Context) (bool, error)
	ValidateCredentials(ctx context.Context, username string, password string) error
}

type Store struct {
	sessionStore  *sessions.CookieStore
	authenticator Authenticator
	txnManager    models.TxnManager
	config        SessionConfig
}

func NewStore(c SessionConfig, a Authenticator, txnMgr models.TxnManager) *Store {
	ret := &Store{
		sessionStore:  sessions.NewCookieStore(c.GetSessionStoreKey()),
		txnManager:    txnMgr,
		config:        c,
		authenticator: a,
	}

	ret.sessionStore.MaxAge(c.GetMaxSessionAge())
	ret.sessionStore.Options.SameSite = http.SameSiteLaxMode

	return ret
}

func (s *Store) LoginRequired(ctx context.Context) (bool, error) {
	var loginRequired bool
	if err := s.txnManager.WithReadTxn(ctx, func(ctx context.Context) error {
		var err error
		loginRequired, err = s.authenticator.LoginRequired(ctx)
		return err
	}); err != nil {
		return false, err
	}
	return loginRequired, nil
}

func (s *Store) Login(w http.ResponseWriter, r *http.Request) error {
	// ignore error - we want a new session regardless
	newSession, _ := s.sessionStore.Get(r, cookieName)

	username := r.FormValue(usernameFormKey)
	password := r.FormValue(passwordFormKey)

	// authenticate the user
	if err := s.txnManager.WithReadTxn(r.Context(), func(ctx context.Context) error {
		return s.authenticator.ValidateCredentials(ctx, username, password)
	}); err != nil {
		return &InvalidCredentialsError{Username: username}
	}

	logger.Infof("User %s logged in", username)

	newSession.Values[userIDKey] = username
	newSession.Values[loginTimeKey] = time.Now().Unix()

	err := newSession.Save(r, w)
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

type Session struct {
	UserID    string
	LoginTime time.Time
}

func (s *Store) GetSessionUserID(w http.ResponseWriter, r *http.Request) (*Session, error) {
	session, err := s.sessionStore.Get(r, cookieName)
	// ignore errors and treat as an empty user id, so that we handle expired
	// cookie
	if err != nil {
		return nil, nil
	}

	if !session.IsNew {
		// refresh the cookie
		err = session.Save(r, w)
		if err != nil {
			return nil, err
		}

		ret := &Session{}
		ret.UserID, _ = session.Values[userIDKey].(string)
		loginTimeUnix, _ := session.Values[loginTimeKey].(int64)
		ret.LoginTime = time.Unix(loginTimeUnix, 0)

		return ret, nil
	}

	return nil, nil
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

func GetRequestApiKey(r *http.Request) string {
	apiKey := r.Header.Get(ApiKeyHeader)

	// try getting the api key as a query parameter
	if apiKey == "" {
		apiKey = r.URL.Query().Get(ApiKeyParameter)
	}

	return apiKey
}
