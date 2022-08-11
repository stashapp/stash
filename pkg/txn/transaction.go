package txn

import (
	"context"
	"fmt"
)

type Manager interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	IsLocked(err error) bool

	AddPostCommitHook(ctx context.Context, hook TxnFunc)
	AddPostRollbackHook(ctx context.Context, hook TxnFunc)
}

type DatabaseProvider interface {
	WithDatabase(ctx context.Context) (context.Context, error)
}

type TxnFunc func(ctx context.Context) error

// WithTxn executes fn in a transaction. If fn returns an error then
// the transaction is rolled back. Otherwise it is committed.
func WithTxn(ctx context.Context, m Manager, fn TxnFunc) error {
	var err error
	ctx, err = m.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = m.Rollback(ctx)
			panic(p)
		}

		if err != nil {
			// something went wrong, rollback
			_ = m.Rollback(ctx)
		} else {
			// all good, commit
			err = m.Commit(ctx)
		}
	}()

	err = fn(ctx)
	return err
}

// WithDatabase executes fn with the context provided by p.WithDatabase.
// It does not run inside a transaction, so all database operations will be
// executed in their own transaction.
func WithDatabase(ctx context.Context, p DatabaseProvider, fn TxnFunc) error {
	var err error
	ctx, err = p.WithDatabase(ctx)
	if err != nil {
		return err
	}

	return fn(ctx)
}

type Retryer struct {
	Manager Manager
	Retries int
	OnFail  func(ctx context.Context, err error, attempt int) error
}

func (r Retryer) WithTxn(ctx context.Context, fn TxnFunc) error {
	var attempt int
	var err error
	for attempt = 1; attempt <= r.Retries; attempt++ {
		err = WithTxn(ctx, r.Manager, fn)

		if err == nil {
			return nil
		}

		if !r.Manager.IsLocked(err) {
			return err
		}

		if r.OnFail != nil {
			if err := r.OnFail(ctx, err, attempt); err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", attempt, err)
}
