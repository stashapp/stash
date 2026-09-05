package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stashapp/stash/internal/authn"
	"github.com/stashapp/stash/internal/authz"
	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

var ErrLastActiveAdmin = errors.New("cannot remove or disable the last active Admin")
var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrBootstrapClosed = errors.New("first Admin bootstrap is closed")

var dummyPasswordHash = func() string {
	hash, err := authn.HashPassword("fixed internal timing equalizer; not a credential")
	if err != nil {
		panic(err)
	}
	return hash
}()

type User struct {
	ID                 int64               `db:"id"`
	Username           string              `db:"username"`
	NormalizedUsername string              `db:"normalized_username"`
	PasswordHash       *string             `db:"password_hash"`
	Role               authz.Role          `db:"role"`
	Status             authz.AccountStatus `db:"status"`
	CreatedAt          time.Time           `db:"created_at"`
	UpdatedAt          time.Time           `db:"updated_at"`
}

type UserStore struct{}

// BootstrapAdmin creates and audits the first account. Callers must execute it
// in a write transaction; SQLite's serialized writer lock makes the zero-user
// check and insert one atomic, single-winner operation.
func (s *UserStore) BootstrapAdmin(ctx context.Context, username, password string) (*User, error) {
	count, err := s.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count != 0 {
		return nil, ErrBootstrapClosed
	}
	user, err := s.Create(ctx, username, password, authz.RoleAdmin)
	if err != nil {
		return nil, err
	}
	_, err = dbWrapper.Exec(ctx, `INSERT INTO user_audit_events
		(occurred_at, actor_user_id, event_type, target_type, target_id, result)
		VALUES (?, ?, 'first_admin_bootstrapped', 'user', ?, 'success')`, time.Now().UTC(), user.ID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("auditing first Admin bootstrap: %w", err)
	}
	return user, nil
}

func normalizeUsername(username string) string {
	return cases.Fold().String(norm.NFKC.String(strings.TrimSpace(username)))
}

func (s *UserStore) Create(ctx context.Context, username, password string, role authz.Role) (*User, error) {
	username = strings.TrimSpace(username)
	normalized := normalizeUsername(username)
	if username == "" || normalized == "" {
		return nil, errors.New("username is required")
	}
	if role != authz.RoleUser && role != authz.RoleModerator && role != authz.RoleAdmin {
		return nil, fmt.Errorf("invalid role %q", role)
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	hash, err := authn.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}
	now := time.Now().UTC()
	result, err := dbWrapper.Exec(ctx, `INSERT INTO users
		(username, normalized_username, password_hash, role, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, username, normalized, hash, role, authz.StatusActive, now, now)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("reading created user id: %w", err)
	}
	return s.Find(ctx, id)
}

func (s *UserStore) Find(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := dbWrapper.Get(ctx, &user, `SELECT id, username, normalized_username, password_hash, role, status, created_at, updated_at FROM users WHERE id = ?`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := dbWrapper.Get(ctx, &user, `SELECT id, username, normalized_username, password_hash, role, status, created_at, updated_at
		FROM users WHERE normalized_username = ?`, normalizeUsername(username))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) Count(ctx context.Context) (int, error) {
	var count int
	if err := dbWrapper.Get(ctx, &count, `SELECT count(*) FROM users`); err != nil {
		return 0, err
	}
	return count, nil
}

// List returns user account metadata. Password hashes are deliberately not
// selected so callers cannot accidentally expose credential material.
func (s *UserStore) List(ctx context.Context) ([]User, error) {
	var users []User
	err := dbWrapper.Select(ctx, &users, `SELECT id, username, role, status, created_at, updated_at FROM users ORDER BY username, id`)
	return users, err
}

func (s *UserStore) ResetPassword(ctx context.Context, id int64, password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	hash, err := authn.HashPassword(password)
	if err != nil {
		return err
	}
	result, err := dbWrapper.Exec(ctx, `UPDATE users SET password_hash = ?, status = ?, updated_at = ?
		WHERE id = ? AND NOT (
			role = 'ADMIN' AND status = 'ACTIVE'
			AND (SELECT count(*) FROM users WHERE role = 'ADMIN' AND status = 'ACTIVE') = 1
		)`, hash, authz.StatusPasswordChangeRequired, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed != 1 {
		var exists int
		if err := dbWrapper.Get(ctx, &exists, `SELECT count(*) FROM users WHERE id = ?`, id); err != nil {
			return err
		}
		if exists == 0 {
			return sql.ErrNoRows
		}
		return ErrLastActiveAdmin
	}
	return nil
}

func (s *UserStore) AuthenticatePassword(ctx context.Context, username, password string) (*User, error) {
	user, err := s.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil || user.PasswordHash == nil {
		_, _ = authn.VerifyPassword(dummyPasswordHash, password)
		return nil, ErrInvalidCredentials
	}
	valid, upgrade, err := authn.VerifyPasswordWithUpgrade(*user.PasswordHash, password)
	if err != nil || !valid || user.Status == authz.StatusDisabled {
		return nil, ErrInvalidCredentials
	}
	if upgrade {
		hash, err := authn.HashPassword(password)
		if err != nil {
			return nil, err
		}
		if _, err := dbWrapper.Exec(ctx, `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ? AND password_hash = ?`, hash, time.Now().UTC(), user.ID, *user.PasswordHash); err != nil {
			return nil, fmt.Errorf("upgrading legacy password hash: %w", err)
		}
		user.PasswordHash = &hash
	}
	return user, nil
}

// SetAccess changes role/status while atomically refusing a transition that
// would leave the database without an active Admin. SQLite serializes the
// write and evaluates the count in the same statement.
func (s *UserStore) SetAccess(ctx context.Context, id int64, role authz.Role, status authz.AccountStatus) error {
	if role != authz.RoleUser && role != authz.RoleModerator && role != authz.RoleAdmin {
		return fmt.Errorf("invalid role %q", role)
	}
	if status != authz.StatusActive && status != authz.StatusDisabled && status != authz.StatusPasswordChangeRequired {
		return fmt.Errorf("invalid account status %q", status)
	}
	result, err := dbWrapper.Exec(ctx, `UPDATE users SET role = ?, status = ?, updated_at = ?
		WHERE id = ? AND NOT (
			role = 'ADMIN' AND status = 'ACTIVE'
			AND (? <> 'ADMIN' OR ? <> 'ACTIVE')
			AND (SELECT count(*) FROM users WHERE role = 'ADMIN' AND status = 'ACTIVE') = 1
		)`, role, status, time.Now().UTC(), id, role, status)
	if err != nil {
		return fmt.Errorf("changing user access: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 1 {
		return nil
	}
	var exists int
	if err := dbWrapper.Get(ctx, &exists, `SELECT count(*) FROM users WHERE id = ?`, id); err != nil {
		return err
	}
	if exists == 0 {
		return sql.ErrNoRows
	}
	return ErrLastActiveAdmin
}
