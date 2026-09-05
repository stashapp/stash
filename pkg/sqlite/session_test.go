package sqlite

import (
	"strconv"
	"testing"
	"time"

	"github.com/stashapp/stash/internal/authn"
	"github.com/stashapp/stash/internal/authz"
)

func TestSessionAuthenticationCSRFAndRevocation(t *testing.T) {
	ctx, _ := userStoreTestContext(t)
	user, err := (&UserStore{}).Create(ctx, "session-user", "password", authz.RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	store := &SessionStore{}
	credentials, err := store.Create(ctx, user.ID, 30*time.Minute, 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	session, principal, err := store.AuthenticatePrincipal(ctx, credentials.ID, credentials.Secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if principal.UserID != strconv.FormatInt(user.ID, 10) || principal.Role != authz.RoleUser {
		t.Fatalf("principal=%+v", principal)
	}
	if session.SecretHash == credentials.Secret || session.SecretHash != authn.HashOpaqueSecret(credentials.Secret) {
		t.Fatal("session plaintext was stored or digest was incorrect")
	}
	if !store.VerifyCSRF(session, credentials.CSRFSecret) || store.VerifyCSRF(session, "wrong") {
		t.Fatal("CSRF verification accepted wrong or rejected correct secret")
	}
	if err := store.Touch(ctx, user.ID, credentials.ID, time.Now(), 48*time.Hour); err != nil {
		t.Fatal(err)
	}
	var idleExpiry, absoluteExpiry time.Time
	if err := dbWrapper.Get(ctx, &idleExpiry, `SELECT idle_expires_at FROM user_sessions WHERE id = ?`, credentials.ID); err != nil {
		t.Fatal(err)
	}
	if err := dbWrapper.Get(ctx, &absoluteExpiry, `SELECT absolute_expires_at FROM user_sessions WHERE id = ?`, credentials.ID); err != nil {
		t.Fatal(err)
	}
	if !idleExpiry.Equal(absoluteExpiry) {
		t.Fatalf("idle expiry %s extended past absolute %s", idleExpiry, absoluteExpiry)
	}
	if _, _, err := store.AuthenticatePrincipal(ctx, credentials.ID, "wrong", time.Now()); err == nil {
		t.Fatal("wrong session secret authenticated")
	}
	if err := store.Revoke(ctx, user.ID+1, credentials.ID); err == nil {
		t.Fatal("cross-user session revocation succeeded")
	}
	if err := store.Revoke(ctx, user.ID, credentials.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.RequireActive(ctx, user.ID, credentials.ID, time.Now()); err == nil {
		t.Fatal("revoked session remained active")
	}
	if err := store.Revoke(ctx, user.ID, credentials.ID); err != nil {
		t.Fatalf("owner retry was not idempotent: %v", err)
	}
	if _, _, err := store.AuthenticatePrincipal(ctx, credentials.ID, credentials.Secret, time.Now()); err == nil {
		t.Fatal("revoked session authenticated")
	}
}

func TestSessionRequireActive(t *testing.T) {
	ctx, _ := userStoreTestContext(t)
	store := &SessionStore{}

	user, err := (&UserStore{}).Create(ctx, "require-active-user", "password", authz.RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	other, err := (&UserStore{}).Create(ctx, "require-active-other", "password", authz.RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	credentials, err := store.Create(ctx, user.ID, 30*time.Minute, 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	if err := store.RequireActive(ctx, user.ID, credentials.ID, time.Now()); err != nil {
		t.Fatalf("active session should be active: %v", err)
	}
	if err := store.RequireActive(ctx, other.ID, credentials.ID, time.Now()); err == nil {
		t.Fatal("RequireActive succeeded for a session belonging to a different user")
	}

	if _, err := dbWrapper.Exec(ctx, `UPDATE users SET status = 'DISABLED' WHERE id = ?`, user.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.RequireActive(ctx, user.ID, credentials.ID, time.Now()); err == nil {
		t.Fatal("RequireActive succeeded for a disabled user")
	}
	if _, err := dbWrapper.Exec(ctx, `UPDATE users SET status = 'ACTIVE' WHERE id = ?`, user.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := dbWrapper.Exec(ctx, `UPDATE user_sessions SET idle_expires_at = ? WHERE id = ?`, time.Now().Add(-time.Minute), credentials.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.RequireActive(ctx, user.ID, credentials.ID, time.Now()); err == nil {
		t.Fatal("RequireActive succeeded for an idle-expired session")
	}
	if _, err := dbWrapper.Exec(ctx, `UPDATE user_sessions SET idle_expires_at = ? WHERE id = ?`, time.Now().Add(time.Hour), credentials.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := dbWrapper.Exec(ctx, `UPDATE user_sessions SET absolute_expires_at = ? WHERE id = ?`, time.Now().Add(-time.Minute), credentials.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.RequireActive(ctx, user.ID, credentials.ID, time.Now()); err == nil {
		t.Fatal("RequireActive succeeded for an absolute-expired session")
	}
	if _, err := dbWrapper.Exec(ctx, `UPDATE user_sessions SET absolute_expires_at = ? WHERE id = ?`, time.Now().Add(24*time.Hour), credentials.ID); err != nil {
		t.Fatal(err)
	}

	if err := store.Revoke(ctx, user.ID, credentials.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.RequireActive(ctx, user.ID, credentials.ID, time.Now()); err == nil {
		t.Fatal("RequireActive succeeded for a revoked session")
	}
}

func TestSessionStatusAndExpiry(t *testing.T) {
	ctx, _ := userStoreTestContext(t)
	user, err := (&UserStore{}).Create(ctx, "status-user", "password", authz.RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	store := &SessionStore{}
	credentials, err := store.Create(ctx, user.ID, 30*time.Minute, 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dbWrapper.Exec(ctx, `UPDATE users SET status = 'PASSWORD_CHANGE_REQUIRED' WHERE id = ?`, user.ID); err != nil {
		t.Fatal(err)
	}
	_, principal, err := store.AuthenticatePrincipal(ctx, credentials.ID, credentials.Secret, time.Now())
	if err != nil || !principal.CanChangeRequiredPassword() || principal.Allows(authz.LibraryRead) {
		t.Fatalf("password-change principal=%+v err=%v", principal, err)
	}
	if _, err := dbWrapper.Exec(ctx, `UPDATE users SET status = 'DISABLED' WHERE id = ?`, user.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.AuthenticatePrincipal(ctx, credentials.ID, credentials.Secret, time.Now()); err == nil {
		t.Fatal("disabled user session authenticated")
	}
	if _, err := dbWrapper.Exec(ctx, `UPDATE users SET status = 'ACTIVE' WHERE id = ?`, user.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := dbWrapper.Exec(ctx, `UPDATE user_sessions SET idle_expires_at = ? WHERE id = ?`, time.Now().Add(-time.Minute), credentials.ID); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.AuthenticatePrincipal(ctx, credentials.ID, credentials.Secret, time.Now()); err == nil {
		t.Fatal("idle-expired session authenticated")
	}
}
