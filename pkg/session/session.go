package session

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
)

type key int

const (
	ContextUser key = iota
)

const (
	cookieName      = "session"
	usernameFormKey = "username"
	passwordFormKey = "password"
	userIDKey       = "userID"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type Store struct {
	sessionStore *sessions.CookieStore
	config       *config.Instance
}

func NewStore(c *config.Instance) *Store {
	ret := &Store{
		sessionStore: sessions.NewCookieStore(config.GetInstance().GetSessionStoreKey()),
		config:       c,
	}

	ret.sessionStore.MaxAge(config.GetInstance().GetMaxSessionAge())

	return ret
}

func (s *Store) Login(w http.ResponseWriter, r *http.Request) error {
	// ignore error - we want a new session regardless
	newSession, _ := s.sessionStore.Get(r, cookieName)

	username := r.FormValue(usernameFormKey)
	password := r.FormValue(passwordFormKey)

	// authenticate the user
	if !config.GetInstance().ValidateCredentials(username, password) {
		return ErrInvalidCredentials
	}

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

// GetCurrentUserID gets the current user id from the provided context
func GetCurrentUserID(ctx context.Context) *string {
	userCtxVal := ctx.Value(ContextUser)
	if userCtxVal != nil {
		currentUser := userCtxVal.(string)
		return &currentUser
	}

	return nil
}

func (s *Store) createSessionCookie(username string) (*http.Cookie, error) {
	session := sessions.NewSession(s.sessionStore, cookieName)
	session.Values[userIDKey] = username

	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.sessionStore.Codecs...)
	if err != nil {
		return nil, err
	}

	return sessions.NewCookie(session.Name(), encoded, session.Options), nil
}

func (s *Store) MakePluginCookie(ctx context.Context) *http.Cookie {
	currentUser := GetCurrentUserID(ctx)

	var cookie *http.Cookie
	var err error
	if currentUser != nil {
		cookie, err = s.createSessionCookie(*currentUser)
		if err != nil {
			logger.Errorf("error creating session cookie: %s", err.Error())
		}
	}

	return cookie
}
