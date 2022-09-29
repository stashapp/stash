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
}

type DatabaseProvider interface {
	WithDatabase(ctx context.Context) (context.Context, error)
}

type TxnFunc func(ctx context.Context) error

// WithTxn executes fn in a transaction. If fn returns an error then
// the transaction is rolled back. Otherwise it is committed.
func WithTxn(ctx context.Context, m Manager, fn TxnFunc) error {
	const execComplete = true
	return withTxn(ctx, m, fn, execComplete)
}

func withTxn(ctx context.Context, m Manager, fn TxnFunc, execCompleteOnLocked bool) error {
	var err error
	ctx, err = begin(ctx, m)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			rollback(ctx, m)
			panic(p)
		}

		if err != nil {
			// something went wrong, rollback
			rollback(ctx, m)

			if execCompleteOnLocked || !m.IsLocked(err) {
				executePostCompleteHooks(ctx)
			}
		} else {
			// all good, commit
			err = commit(ctx, m)
			executePostCompleteHooks(ctx)
		}

	}()

	err = fn(ctx)
	return err
}

func begin(ctx context.Context, m Manager) (context.Context, error) {
	var err error
	ctx, err = m.Begin(ctx)
	if err != nil {
		return nil, err
	}

	hm := hookManager{}
	ctx = hm.register(ctx)

	return ctx, nil
}

func commit(ctx context.Context, m Manager) error {
	if err := m.Commit(ctx); err != nil {
		return err
	}

	executePostCommitHooks(ctx)
	return nil
}

func rollback(ctx context.Context, m Manager) {
	if err := m.Rollback(ctx); err != nil {
		return
	}

	executePostRollbackHooks(ctx)
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
	// use value < 0 to retry forever
	Retries int
	OnFail  func(ctx context.Context, err error, attempt int) error
}

func (r Retryer) WithTxn(ctx context.Context, fn TxnFunc) error {
	var attempt int
	var err error
	for attempt = 1; attempt <= r.Retries || r.Retries < 0; attempt++ {
		const execComplete = false
		err = withTxn(ctx, r.Manager, fn, execComplete)

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
