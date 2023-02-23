package txn

import (
	"context"
	"fmt"
)

type Manager interface {
	Begin(ctx context.Context, exclusive bool) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	IsLocked(err error) bool
}

type DatabaseProvider interface {
	WithDatabase(ctx context.Context) (context.Context, error)
}

type DatabaseProviderManager interface {
	DatabaseProvider
	Manager
}

type TxnFunc func(ctx context.Context) error

// WithTxn executes fn in a transaction. If fn returns an error then
// the transaction is rolled back. Otherwise it is committed.
// Transaction is exclusive. Only one thread may run a transaction
// using this function at a time. This function will wait until the
// lock is available before executing.
// This function should be used for making changes to the database.
func WithTxn(ctx context.Context, m Manager, fn TxnFunc) error {
	const (
		execComplete = true
		exclusive    = true
	)
	return withTxn(ctx, m, fn, exclusive, execComplete)
}

// WithReadTxn executes fn in a transaction. If fn returns an error then
// the transaction is rolled back. Otherwise it is committed.
// Transaction is not exclusive and does not enforce read-only restrictions.
// Multiple threads can run transactions using this function concurrently,
// but concurrent writes may result in locked database error.
func WithReadTxn(ctx context.Context, m Manager, fn TxnFunc) error {
	const (
		execComplete = true
		exclusive    = false
	)
	return withTxn(ctx, m, fn, exclusive, execComplete)
}

func withTxn(outerCtx context.Context, m Manager, fn TxnFunc, exclusive bool, execCompleteOnLocked bool) error {
	ctx, err := begin(outerCtx, m, exclusive)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			rollback(ctx, outerCtx, m)
			panic(p)
		}

		if err != nil {
			// something went wrong, rollback
			rollback(ctx, outerCtx, m)

			if execCompleteOnLocked || !m.IsLocked(err) {
				executePostCompleteHooks(ctx, outerCtx)
			}
		} else {
			// all good, commit
			err = commit(ctx, outerCtx, m)
			executePostCompleteHooks(ctx, outerCtx)
		}

	}()

	err = fn(ctx)
	return err
}

func begin(ctx context.Context, m Manager, exclusive bool) (context.Context, error) {
	var err error
	ctx, err = m.Begin(ctx, exclusive)
	if err != nil {
		return nil, err
	}

	hm := hookManager{}
	ctx = hm.register(ctx)

	return ctx, nil
}

func commit(ctx context.Context, outerCtx context.Context, m Manager) error {
	if err := m.Commit(ctx); err != nil {
		return err
	}

	executePostCommitHooks(ctx, outerCtx)
	return nil
}

func rollback(ctx context.Context, outerCtx context.Context, m Manager) {
	if err := m.Rollback(ctx); err != nil {
		return
	}

	executePostRollbackHooks(ctx, outerCtx)
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

// Retryer is a provides WithTxn function that retries the transaction
// if it fails with a locked database error.
// Transactions are run in exclusive mode.
type Retryer struct {
	Manager Manager
	// use value < 0 to retry forever
	Retries int
	OnFail  func(ctx context.Context, err error, attempt int) error
}

func (r Retryer) WithTxn(ctx context.Context, fn TxnFunc) error {
	var attempt int
	var err error
	for attempt = 1; attempt <= r.Retries || r.Retries < 0; attempt++ {
		const (
			execComplete = false
			exclusive    = true
		)
		err = withTxn(ctx, r.Manager, fn, exclusive, execComplete)

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
