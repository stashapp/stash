package authservice

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"
)

type LoginService struct {
	Database *sqlite.Database
}

const (
	DefaultSessionIdle     = 30 * time.Minute
	DefaultSessionAbsolute = 30 * 24 * time.Hour
)

func (s LoginService) AuthenticatePassword(ctx context.Context, username, password string) (*sqlite.User, error) {
	var user *sqlite.User
	retryer := txn.Retryer{
		Manager: s.Database,
		Retries: 5,
		OnFail: func(ctx context.Context, _ error, attempt int) error {
			timer := time.NewTimer(time.Duration(attempt) * 10 * time.Millisecond)
			defer timer.Stop()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				return nil
			}
		},
	}
	err := retryer.WithTxn(ctx, func(txCtx context.Context) error {
		var err error
		user, err = s.Database.User.AuthenticatePassword(txCtx, username, password)
		return err
	})
	return user, err
}

func (s LoginService) Login(ctx context.Context, username, password string) (*sqlite.User, *sqlite.SessionCredentials, error) {
	var user *sqlite.User
	var credentials *sqlite.SessionCredentials
	retryer := txn.Retryer{Manager: s.Database, Retries: 5, OnFail: loginRetryBackoff}
	err := retryer.WithTxn(ctx, func(txCtx context.Context) error {
		var err error
		user, err = s.Database.User.AuthenticatePassword(txCtx, username, password)
		if err != nil {
			return err
		}
		credentials, err = s.Database.Session.Create(txCtx, user.ID, DefaultSessionIdle, DefaultSessionAbsolute)
		return err
	})
	// Authentication auditing is deliberately best-effort and occurs after the
	// login transaction: audit storage failure must never block a valid login or
	// turn invalid-password traffic into an availability hazard. Failures store
	// no attempted username, address, credential, or request detail.
	var actorID *int64
	result := "failure"
	if err == nil && user != nil {
		actorID = &user.ID
		result = "success"
	}
	if auditErr := txn.WithTxn(context.WithoutCancel(ctx), s.Database, func(txCtx context.Context) error {
		return s.Database.Audit.RecordAuthentication(txCtx, actorID, "login", result)
	}); auditErr != nil {
		logger.Warnf("recording redacted login audit event failed: %v", auditErr)
	}
	return user, credentials, err
}

func loginRetryBackoff(ctx context.Context, _ error, attempt int) error {
	timer := time.NewTimer(time.Duration(attempt) * 10 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
