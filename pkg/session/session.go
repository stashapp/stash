package session

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sliceutil"
)

type key int

const (
	contextUser key = iota
	contextVisitedPlugins
)

const (
	userIDKey         = "userID"
	visitedPluginsKey = "visitedPlugins"
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

type Store struct {
	sessionStore *sessions.CookieStore
	config       SessionConfig
}

func NewStore(c SessionConfig) *Store {
	ret := &Store{
		sessionStore: sessions.NewCookieStore(c.GetSessionStoreKey()),
		config:       c,
	}

	ret.sessionStore.MaxAge(c.GetMaxSessionAge())
	ret.sessionStore.Options.SameSite = http.SameSiteLaxMode

	return ret
}

func (s *Store) Login(w http.ResponseWriter, r *http.Request) error {
	// ignore error - we want a new session regardless
	newSession, _ := s.sessionStore.Get(r, cookieName)

	username := r.FormValue(usernameFormKey)
	password := r.FormValue(passwordFormKey)

	// authenticate the user
	if !s.config.ValidateCredentials(username, password) {
		return &InvalidCredentialsError{Username: username}
	}

	// since we only have one user, don't leak the name
	logger.Info("User logged in")

	newSession.Values[userIDKey] = username

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

	delete(session.Values, userIDKey)
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	// since we only have one user, don't leak the name
	logger.Infof("User logged out")

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

func SetCurrentUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextUser, userID)
}

// GetCurrentUserID gets the current user id from the provided context
func GetCurrentUserID(ctx context.Context) *string {
	userCtxVal := ctx.Value(contextUser)
	if userCtxVal != nil {
		currentUser := userCtxVal.(string)
		return &currentUser
	}

	return nil
}

func (s *Store) VisitedPluginHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get the visited plugins from the cookie and set in the context
			session, err := s.sessionStore.Get(r, cookieName)

			// ignore errors
			if err == nil {
				val := session.Values[visitedPluginsKey]

				visitedPlugins, _ := val.([]string)

				ctx := setVisitedPlugins(r.Context(), visitedPlugins)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetVisitedPlugins(ctx context.Context) []string {
	ctxVal := ctx.Value(contextVisitedPlugins)
	if ctxVal != nil {
		return ctxVal.([]string)
	}

	return nil
}

func AddVisitedPlugin(ctx context.Context, pluginID string) context.Context {
	curVal := GetVisitedPlugins(ctx)
	curVal = sliceutil.AppendUnique(curVal, pluginID)
	return setVisitedPlugins(ctx, curVal)
}

func setVisitedPlugins(ctx context.Context, visitedPlugins []string) context.Context {
	return context.WithValue(ctx, contextVisitedPlugins, visitedPlugins)
}

func (s *Store) MakePluginCookie(ctx context.Context) *http.Cookie {
	currentUser := GetCurrentUserID(ctx)
	visitedPlugins := GetVisitedPlugins(ctx)

	session := sessions.NewSession(s.sessionStore, cookieName)
	if currentUser != nil {
		session.Values[userIDKey] = *currentUser
	}

	session.Values[visitedPluginsKey] = visitedPlugins

	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.sessionStore.Codecs...)
	if err != nil {
		logger.Errorf("error creating session cookie: %s", err.Error())
		return nil
	}

	return sessions.NewCookie(session.Name(), encoded, session.Options)
}

func (s *Store) Authenticate(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	c := s.config

	// translate api key into current user, if present
	apiKey := r.Header.Get(ApiKeyHeader)

	// try getting the api key as a query parameter
	if apiKey == "" {
		apiKey = r.URL.Query().Get(ApiKeyParameter)
	}

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
