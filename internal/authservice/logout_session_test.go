package authservice

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stashapp/stash/internal/authz"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"
)

// Permanent coverage of internal/api/session.go's handleLogout core
// logic (readDBSessionCookie -> AuthenticatePrincipal -> Revoke), since
// internal/api itself cannot be compiled in this environment (pre-existing
// missing gqlgen codegen, unrelated to this change). Verifies the actual
// pkg/sqlite behavior the handler depends on.

func logoutLikeHandler(db *sqlite.Database, sessionID, sessionSecret string) error {
	retryer := txn.Retryer{Manager: db, Retries: 5}
	return retryer.WithTxn(context.Background(), func(ctx context.Context) error {
		sessionRecord, _, err := db.Session.AuthenticatePrincipal(ctx, sessionID, sessionSecret, time.Now())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}
		return db.Session.Revoke(ctx, sessionRecord.UserID, sessionRecord.ID)
	})
}

func newRealDB(t *testing.T) (*sqlite.Database, string) {
	t.Helper()
	config.InitializeEmpty()
	path := filepath.Join(t.TempDir(), "stash.sqlite")
	db := sqlite.NewDatabase()
	if err := db.Open(path); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db, path
}

func createUserAndSession(t *testing.T, db *sqlite.Database, username string) (*sqlite.User, *sqlite.SessionCredentials) {
	t.Helper()
	var user *sqlite.User
	var creds *sqlite.SessionCredentials
	if err := txn.WithTxn(context.Background(), db, func(ctx context.Context) error {
		var err error
		user, err = db.User.Create(ctx, username, "password12345", authz.RoleUser)
		if err != nil {
			return err
		}
		creds, err = db.Session.Create(ctx, user.ID, 30*time.Minute, 24*time.Hour)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	return user, creds
}

func TestLogoutRevokesAndReplayFails(t *testing.T) {
	db, _ := newRealDB(t)
	_, creds := createUserAndSession(t, db, "logout-user")

	// Confirm the session works before logout.
	if err := txn.WithReadTxn(context.Background(), db, func(ctx context.Context) error {
		_, _, err := db.Session.AuthenticatePrincipal(ctx, creds.ID, creds.Secret, time.Now())
		return err
	}); err != nil {
		t.Fatalf("session should authenticate before logout: %v", err)
	}

	if err := logoutLikeHandler(db, creds.ID, creds.Secret); err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	// Replay of the OLD cookie must now fail.
	err := txn.WithReadTxn(context.Background(), db, func(ctx context.Context) error {
		_, _, err := db.Session.AuthenticatePrincipal(ctx, creds.ID, creds.Secret, time.Now())
		return err
	})
	if err == nil {
		t.Fatal("old session cookie replayed successfully after logout — session was not actually revoked")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows-classified failure after logout, got: %v", err)
	}
}

func TestLogoutIdempotentOnAlreadyRevoked(t *testing.T) {
	db, _ := newRealDB(t)
	_, creds := createUserAndSession(t, db, "idempotent-user")

	if err := logoutLikeHandler(db, creds.ID, creds.Secret); err != nil {
		t.Fatalf("first logout failed: %v", err)
	}
	// Second logout call with the same (now-revoked) cookie must NOT error.
	if err := logoutLikeHandler(db, creds.ID, creds.Secret); err != nil {
		t.Fatalf("second logout on already-revoked session should be idempotent, got error: %v", err)
	}
}

func TestLogoutOnInvalidCookieIsIdempotentNoError(t *testing.T) {
	db, _ := newRealDB(t)
	// No session ever created for these credentials.
	if err := logoutLikeHandler(db, "nonexistent-id", "nonexistent-secret"); err != nil {
		t.Fatalf("logout with an invalid/never-existed cookie should not error, got: %v", err)
	}
}

func TestLogoutCrossUserRevokeImpossible(t *testing.T) {
	db, _ := newRealDB(t)
	_, credsA := createUserAndSession(t, db, "user-a")
	_, credsB := createUserAndSession(t, db, "user-b")

	// Logging out with user A's cookie must not be able to touch user B's
	// session under any circumstance — Revoke's userID argument comes only
	// from the session record resolved FROM credsA, never from anything
	// caller-suppliable, so this is really testing there's no way to smuggle
	// user B's session ID through this path at all.
	if err := logoutLikeHandler(db, credsA.ID, credsA.Secret); err != nil {
		t.Fatalf("logout for user A failed: %v", err)
	}

	// User B's session must still be fully valid.
	if err := txn.WithReadTxn(context.Background(), db, func(ctx context.Context) error {
		_, _, err := db.Session.AuthenticatePrincipal(ctx, credsB.ID, credsB.Secret, time.Now())
		return err
	}); err != nil {
		t.Fatalf("user B's session was affected by user A's logout: %v", err)
	}
}

func TestTouchExpiryClassifiesAsErrNoRows(t *testing.T) {
	db, path := newRealDB(t)
	user, creds := createUserAndSession(t, db, "touch-user")

	// Force the session to already be idle-expired, via a raw connection to
	// the same file (pkg/sqlite's dbWrapper is unexported outside its package).
	raw, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := raw.Exec(`UPDATE user_sessions SET idle_expires_at = ? WHERE id = ?`, time.Now().Add(-time.Minute), creds.ID); err != nil {
		raw.Close()
		t.Fatal(err)
	}
	raw.Close()

	err = txn.WithTxn(context.Background(), db, func(ctx context.Context) error {
		return db.Session.Touch(ctx, user.ID, creds.ID, time.Now(), 30*time.Minute)
	})
	if err == nil {
		t.Fatal("Touch on an idle-expired session should fail")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected Touch failure to be sql.ErrNoRows-classified (so authenticateHandler treats it as 401, not 500), got: %v", err)
	}
}
