package sqlite

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stashapp/stash/internal/authn"
	"github.com/stashapp/stash/internal/authz"
)

// argon2Params extracts the m,t,p cost parameters from a PHC-formatted
// Argon2id hash string, for structural work-factor comparisons that avoid
// brittle wall-clock timing assertions in CI.
func argon2Params(t *testing.T, encoded string) string {
	t.Helper()
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		t.Fatalf("not a well-formed Argon2id hash: %q", encoded)
	}
	return parts[3] // "m=...,t=...,p=..."
}

func userStoreTestContext(t *testing.T) (context.Context, *sqlx.Tx) {
	t.Helper()
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA foreign_keys=ON`); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE users (
		id integer primary key autoincrement, username text not null,
		normalized_username text not null unique, password_hash text,
		role text not null check(role in ('USER','MODERATOR','ADMIN')),
		status text not null check(status in ('ACTIVE','DISABLED','PASSWORD_CHANGE_REQUIRED')),
		created_at datetime not null, updated_at datetime not null);
		CREATE TABLE user_sessions (
		id text primary key, user_id integer not null references users(id) on delete cascade, secret_hash text not null unique,
		csrf_hash text not null, created_at datetime not null, last_seen_at datetime not null,
		idle_expires_at datetime not null, absolute_expires_at datetime not null, revoked_at datetime,
		user_agent text, remote_address text);
		CREATE TABLE user_api_tokens (
		id text primary key, user_id integer not null references users(id) on delete cascade, name text not null,
		secret_hash text not null unique, scopes_json text, created_at datetime not null,
		expires_at datetime not null, last_used_at datetime, revoked_at datetime);
		CREATE TABLE user_audit_events (
		id integer not null primary key autoincrement, occurred_at datetime not null,
		actor_user_id integer references users(id) on delete set null, event_type varchar(100) not null,
		target_type varchar(100), target_id varchar(255), result varchar(40) not null, details_json text)`)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	return context.WithValue(context.Background(), txnKey, tx), tx
}

func TestUserStoreCreateNormalizesAndHashes(t *testing.T) {
	ctx, _ := userStoreTestContext(t)
	store := &UserStore{}
	user, err := store.Create(ctx, "Ａdmin", "not stored plaintext", authz.RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if user.NormalizedUsername != "admin" || user.PasswordHash == nil || *user.PasswordHash == "not stored plaintext" {
		t.Fatalf("unexpected stored user: %+v", user)
	}
	if ok, err := authn.VerifyPassword(*user.PasswordHash, "not stored plaintext"); err != nil || !ok {
		t.Fatalf("stored password verification ok=%v err=%v", ok, err)
	}
	if _, err := store.Create(ctx, "admin", "different", authz.RoleUser); err == nil {
		t.Fatal("Unicode-equivalent normalized username was accepted")
	}
}

// TestAuthenticatePasswordUniformFailureAndWorkFactor is the permanent
// replacement for the throwaway timing repro used to verify the dummy-hash
// enumeration-resistance path: missing username, wrong password, nil
// password_hash, and a disabled user with the correct password must all
// return the identical ErrInvalidCredentials, and the dummy hash compared
// against for a missing user must carry the same Argon2id work factor as a
// real account hash (structural check, not wall-clock timing — CI machine
// speed makes timing thresholds brittle; an empirical timing run during
// audit separately confirmed real vs. missing vs. disabled paths cost the
// same wall-clock time to within noise).
func TestAuthenticatePasswordUniformFailureAndWorkFactor(t *testing.T) {
	ctx, _ := userStoreTestContext(t)
	store := &UserStore{}

	if _, err := store.Create(ctx, "active-user", "correct password", authz.RoleUser); err != nil {
		t.Fatal(err)
	}

	disabled, err := store.Create(ctx, "disabled-user", "correct password", authz.RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dbWrapper.Exec(ctx, "UPDATE users SET status = 'DISABLED' WHERE id = ?", disabled.ID); err != nil {
		t.Fatal(err)
	}

	noHash, err := store.Create(ctx, "no-password-user", "placeholder", authz.RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dbWrapper.Exec(ctx, "UPDATE users SET password_hash = NULL WHERE id = ?", noHash.ID); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name     string
		username string
		password string
	}{
		{"missing username", "does-not-exist", "correct password"},
		{"wrong password", "active-user", "wrong password"},
		{"nil password hash", "no-password-user", "anything"},
		{"disabled user, correct password", "disabled-user", "correct password"},
		{"disabled user, wrong password", "disabled-user", "wrong password"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := store.AuthenticatePassword(ctx, c.username, c.password); err != ErrInvalidCredentials {
				t.Fatalf("err=%v want ErrInvalidCredentials", err)
			}
		})
	}

	// Work-factor regression guard: the dummy hash used for the
	// missing-user/nil-hash path must carry the SAME Argon2id cost
	// parameters as a real account hash, not a cheap placeholder that
	// would silently reintroduce a timing side-channel.
	dummyOK, err := authn.VerifyPassword(dummyPasswordHash, "not the dummy password")
	if err != nil {
		t.Fatalf("dummy hash is not a well-formed Argon2id hash: %v", err)
	}
	if dummyOK {
		t.Fatal("dummy hash matched an arbitrary password")
	}
	realHash, err := authn.HashPassword("reference password")
	if err != nil {
		t.Fatal(err)
	}
	if dummyParams, realParams := argon2Params(t, dummyPasswordHash), argon2Params(t, realHash); dummyParams != realParams {
		t.Fatalf("dummy hash work factor %q does not match real account work factor %q", dummyParams, realParams)
	}
}

func TestUserStoreLastActiveAdmin(t *testing.T) {
	ctx, _ := userStoreTestContext(t)
	store := &UserStore{}
	first, err := store.Create(ctx, "first", "password", authz.RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SetAccess(ctx, first.ID, authz.RoleUser, authz.StatusActive); !errors.Is(err, ErrLastActiveAdmin) {
		t.Fatalf("sole Admin demotion err=%v", err)
	}
	second, err := store.Create(ctx, "second", "password", authz.RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SetAccess(ctx, first.ID, authz.RoleUser, authz.StatusActive); err != nil {
		t.Fatalf("Admin demotion with replacement: %v", err)
	}
	if err := store.SetAccess(ctx, second.ID, authz.RoleAdmin, authz.StatusDisabled); !errors.Is(err, ErrLastActiveAdmin) {
		t.Fatalf("last Admin disable err=%v", err)
	}
}

func TestUserStoreConcurrentLastActiveAdmin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.sqlite")
	db, err := sqlx.Open("sqlite3", fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL", path))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(2)
	_, err = db.Exec(`CREATE TABLE users (
		id integer primary key autoincrement, username text not null,
		normalized_username text not null unique, password_hash text,
		role text not null check(role in ('USER','MODERATOR','ADMIN')),
		status text not null check(status in ('ACTIVE','DISABLED','PASSWORD_CHANGE_REQUIRED')),
		created_at datetime not null, updated_at datetime not null);
		INSERT INTO users (username, normalized_username, password_hash, role, status, created_at, updated_at)
		VALUES ('a', 'a', 'test', 'ADMIN', 'ACTIVE', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		       ('b', 'b', 'test', 'ADMIN', 'ACTIVE', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`)
	if err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, id := range []int64{1, 2} {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			tx, err := db.Beginx()
			if err != nil {
				results <- err
				return
			}
			defer tx.Rollback()
			<-start
			ctx := context.WithValue(context.Background(), txnKey, tx)
			err = (&UserStore{}).SetAccess(ctx, id, authz.RoleUser, authz.StatusActive)
			if err == nil {
				err = tx.Commit()
			}
			results <- err
		}(id)
	}
	close(start)
	wg.Wait()
	close(results)

	succeeded, protected := 0, 0
	for err := range results {
		switch {
		case err == nil:
			succeeded++
		case errors.Is(err, ErrLastActiveAdmin):
			protected++
		default:
			t.Fatalf("unexpected concurrent result: %v", err)
		}
	}
	if succeeded != 1 || protected != 1 {
		t.Fatalf("succeeded=%d protected=%d want 1/1", succeeded, protected)
	}
	var activeAdmins int
	if err := db.Get(&activeAdmins, `SELECT count(*) FROM users WHERE role = 'ADMIN' AND status = 'ACTIVE'`); err != nil {
		t.Fatal(err)
	}
	if activeAdmins != 1 {
		t.Fatalf("active Admins=%d want=1", activeAdmins)
	}
}
