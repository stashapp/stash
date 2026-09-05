package sqlite

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/stashapp/stash/internal/authn"
	"github.com/stashapp/stash/internal/authz"
)

type SessionCredentials struct {
	ID         string
	Secret     string
	CSRFSecret string
}

type Session struct {
	ID                string     `db:"id"`
	UserID            int64      `db:"user_id"`
	SecretHash        string     `db:"secret_hash"`
	CSRFHash          string     `db:"csrf_hash"`
	CreatedAt         time.Time  `db:"created_at"`
	LastSeenAt        time.Time  `db:"last_seen_at"`
	IdleExpiresAt     time.Time  `db:"idle_expires_at"`
	AbsoluteExpiresAt time.Time  `db:"absolute_expires_at"`
	RevokedAt         *time.Time `db:"revoked_at"`
}

type SessionStore struct{}

type sessionPrincipalRow struct {
	Session
	Role   authz.Role          `db:"role"`
	Status authz.AccountStatus `db:"status"`
}

func (s *SessionStore) Create(ctx context.Context, userID int64, idle, absolute time.Duration) (*SessionCredentials, error) {
	if idle <= 0 || absolute <= 0 || idle > absolute {
		return nil, errors.New("invalid session lifetime")
	}
	id, err := authn.NewOpaqueSecret(32)
	if err != nil {
		return nil, err
	}
	secret, err := authn.NewOpaqueSecret(32)
	if err != nil {
		return nil, err
	}
	csrf, err := authn.NewOpaqueSecret(32)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	_, err = dbWrapper.Exec(ctx, `INSERT INTO user_sessions
		(id, user_id, secret_hash, csrf_hash, created_at, last_seen_at, idle_expires_at, absolute_expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, id, userID, authn.HashOpaqueSecret(secret), authn.HashOpaqueSecret(csrf), now, now, now.Add(idle), now.Add(absolute))
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}
	return &SessionCredentials{ID: id, Secret: secret, CSRFSecret: csrf}, nil
}

func (s *SessionStore) Authenticate(ctx context.Context, id, secret string, now time.Time) (*Session, error) {
	var session Session
	err := dbWrapper.Get(ctx, &session, `SELECT id, user_id, secret_hash, csrf_hash, created_at, last_seen_at, idle_expires_at, absolute_expires_at, revoked_at
		FROM user_sessions WHERE id = ? AND secret_hash = ? AND revoked_at IS NULL
		AND idle_expires_at > ? AND absolute_expires_at > ?`, id, authn.HashOpaqueSecret(secret), now.UTC(), now.UTC())
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionStore) AuthenticatePrincipal(ctx context.Context, id, secret string, now time.Time) (*Session, authz.Principal, error) {
	var row sessionPrincipalRow
	err := dbWrapper.Get(ctx, &row, `SELECT s.id, s.user_id, s.secret_hash, s.csrf_hash, s.created_at, s.last_seen_at,
		s.idle_expires_at, s.absolute_expires_at, s.revoked_at, u.role, u.status
		FROM user_sessions s JOIN users u ON u.id = s.user_id
		WHERE s.id = ? AND s.secret_hash = ? AND s.revoked_at IS NULL
		AND s.idle_expires_at > ? AND s.absolute_expires_at > ? AND u.status <> 'DISABLED'`,
		id, authn.HashOpaqueSecret(secret), now.UTC(), now.UTC())
	if err != nil {
		return nil, authz.Principal{}, err
	}
	principal := authz.Principal{UserID: fmt.Sprint(row.UserID), Role: row.Role, Status: row.Status}
	return &row.Session, principal, nil
}

func (s *SessionStore) VerifyCSRF(session *Session, secret string) bool {
	if session == nil || session.CSRFHash == "" || secret == "" {
		return false
	}
	got := authn.HashOpaqueSecret(secret)
	return subtle.ConstantTimeCompare([]byte(session.CSRFHash), []byte(got)) == 1
}

// RequireActive verifies that a previously authenticated session still belongs
// to an enabled user and has not been revoked or expired. It deliberately does
// not accept a bearer secret: callers may retain the opaque session ID for
// connection-lifetime revalidation without retaining credentials in memory.
func (s *SessionStore) RequireActive(ctx context.Context, userID int64, id string, now time.Time) error {
	var found int
	err := dbWrapper.Get(ctx, &found, `SELECT 1
		FROM user_sessions s JOIN users u ON u.id = s.user_id
		WHERE s.id = ? AND s.user_id = ? AND s.revoked_at IS NULL
		AND s.idle_expires_at > ? AND s.absolute_expires_at > ? AND u.status <> 'DISABLED'`,
		id, userID, now.UTC(), now.UTC())
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionStore) Touch(ctx context.Context, userID int64, id string, now time.Time, idle time.Duration) error {
	if idle <= 0 {
		return errors.New("invalid idle lifetime")
	}
	now = now.UTC()
	requested := now.Add(idle)
	result, err := dbWrapper.Exec(ctx, `UPDATE user_sessions SET last_seen_at = ?,
		idle_expires_at = CASE WHEN absolute_expires_at < ? THEN absolute_expires_at ELSE ? END
		WHERE id = ? AND user_id = ? AND revoked_at IS NULL
		AND idle_expires_at > ? AND absolute_expires_at > ?`, now, requested, requested, id, userID, now, now)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed != 1 {
		return fmt.Errorf("active session not found: %w", sql.ErrNoRows)
	}
	return nil
}

func (s *SessionStore) Revoke(ctx context.Context, userID int64, id string) error {
	result, err := dbWrapper.Exec(ctx, `UPDATE user_sessions SET revoked_at = COALESCE(revoked_at, ?) WHERE id = ? AND user_id = ?`, time.Now().UTC(), id, userID)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed != 1 {
		return errors.New("session not found")
	}
	return nil
}

func (s *SessionStore) RevokeAllForUser(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}
	_, err := dbWrapper.Exec(ctx, `UPDATE user_sessions SET revoked_at = COALESCE(revoked_at, ?) WHERE user_id = ?`, time.Now().UTC(), userID)
	return err
}
