package authservice

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stashapp/stash/internal/authz"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"
)

func TestConcurrentLegacyLoginsUpgradeWithoutFailure(t *testing.T) {
	config.InitializeEmpty()
	path := filepath.Join(t.TempDir(), "stash.sqlite")
	db := sqlite.NewDatabase()
	if err := db.Open(path); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var user *sqlite.User
	if err := txn.WithTxn(context.Background(), db, func(ctx context.Context) error {
		var err error
		user, err = db.User.Create(ctx, "legacy", "temporary", authz.RoleAdmin)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	legacy, err := bcrypt.GenerateFromPassword([]byte("legacy password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := raw.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, string(legacy), user.ID); err != nil {
		raw.Close()
		t.Fatal(err)
	}
	raw.Close()

	service := LoginService{Database: db}
	start := make(chan struct{})
	errs := make(chan error, 6)
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := service.AuthenticatePassword(context.Background(), "legacy", "legacy password")
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("valid concurrent login failed: %v", err)
		}
	}
	raw, err = sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	defer raw.Close()
	var stored string
	if err := raw.QueryRow(`SELECT password_hash FROM users WHERE id = ?`, user.ID).Scan(&stored); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(stored, "$argon2id$") {
		t.Fatal("legacy password hash was not upgraded to Argon2id")
	}
}

func TestLoginAuditIsRedactedAndEnumerationSafe(t *testing.T) {
	config.InitializeEmpty()
	db := sqlite.NewDatabase()
	if err := db.Open(filepath.Join(t.TempDir(), "login-audit.sqlite")); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := txn.WithTxn(context.Background(), db, func(ctx context.Context) error {
		_, err := db.User.Create(ctx, "known-user", "correct-password", authz.RoleUser)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	service := LoginService{Database: db}
	if _, _, err := service.Login(context.Background(), "known-user", "correct-password"); err != nil {
		t.Fatal(err)
	}
	_, _, _ = service.Login(context.Background(), "known-user", "wrong-password")
	_, _, _ = service.Login(context.Background(), "nonexistent-user", "wrong-password")
	var events []sqlite.AuditEvent
	if err := txn.WithReadTxn(context.Background(), db, func(ctx context.Context) error {
		var err error
		events, err = db.Audit.List(ctx)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	var successes, failures int
	for _, event := range events {
		if event.EventType != "login" {
			continue
		}
		if event.TargetType != nil || event.TargetID != nil || event.DetailsJSON != nil {
			t.Fatalf("login event exposed attempted identity or details: %+v", event)
		}
		switch event.Result {
		case "success":
			successes++
			if event.ActorUserID == nil {
				t.Fatal("successful login lacks persisted actor")
			}
		case "failure":
			failures++
			if event.ActorUserID != nil {
				t.Fatalf("failed login identifies account existence: %+v", event)
			}
		}
	}
	if successes != 1 || failures != 2 {
		t.Fatalf("login audit successes=%d failures=%d", successes, failures)
	}
}

func TestLoginAuditWriteFailureDoesNotBlockAuthentication(t *testing.T) {
	config.InitializeEmpty()
	path := filepath.Join(t.TempDir(), "login-audit-failure.sqlite")
	db := sqlite.NewDatabase()
	if err := db.Open(path); err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := txn.WithTxn(context.Background(), db, func(ctx context.Context) error {
		_, err := db.User.Create(ctx, "available-user", "correct-password", authz.RoleUser)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	raw, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := raw.Exec(`DROP TABLE user_audit_events`); err != nil {
		raw.Close()
		t.Fatal(err)
	}
	raw.Close()
	user, credentials, err := (LoginService{Database: db}).Login(context.Background(), "available-user", "correct-password")
	if err != nil || user == nil || credentials == nil {
		t.Fatalf("audit outage blocked valid authentication: user=%v credentials=%v err=%v", user, credentials, err)
	}
}
