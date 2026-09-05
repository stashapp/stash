package migrations

import (
	"context"
	"database/sql"
	_ "embed"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed 87_multiuser_accounts.up.sql
var multiuserAccountsSQL string

func TestMultiuserAccountsMigration(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, multiuserAccountsSQL); err != nil {
		t.Fatal(err)
	}

	for _, table := range []string{"users", "user_sessions", "user_audit_events"} {
		var count int
		if err := db.QueryRowContext(ctx,
			`SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, table,
		).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("table %q was not created", table)
		}
	}

	if _, err := db.ExecContext(ctx, `
		INSERT INTO users (username, normalized_username, password_hash, role, status, created_at, updated_at)
		VALUES ('Admin', 'admin', 'hash', 'ADMIN', 'ACTIVE', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO user_sessions (id, user_id, secret_hash, csrf_hash, created_at, last_seen_at, idle_expires_at, absolute_expires_at)
		VALUES ('session', 1, 'secret-hash', 'csrf-hash', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP,
		        datetime('now', '+30 minutes'), datetime('now', '+1 day'))`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = 1`); err != nil {
		t.Fatal(err)
	}
	var sessions int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM user_sessions`).Scan(&sessions); err != nil {
		t.Fatal(err)
	}
	if sessions != 0 {
		t.Fatal("deleting a user did not cascade to sessions")
	}
}
