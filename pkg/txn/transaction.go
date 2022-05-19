package txn

import "context"

type Manager interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxnFunc func(ctx context.Context) error

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
