package models

import "context"

type Transaction interface {
	Begin() error
	Rollback() error
	Commit() error
	Repository() Repository
}

type ReadTransaction interface {
	Begin() error
	Rollback() error
	Commit() error
	Repository() ReaderRepository
}

type TransactionManager interface {
	WithTxn(ctx context.Context, fn func(r Repository) error) error
	WithReadTxn(ctx context.Context, fn func(r ReaderRepository) error) error
}

func WithTxn(txn Transaction, fn func(r Repository) error) error {
	err := txn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			txn.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			txn.Rollback()
		} else {
			// all good, commit
			err = txn.Commit()
		}
	}()

	err = fn(txn.Repository())
	return err
}

func WithROTxn(txn ReadTransaction, fn func(r ReaderRepository) error) error {
	err := txn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			txn.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			txn.Rollback()
		} else {
			// all good, commit
			err = txn.Commit()
		}
	}()

	err = fn(txn.Repository())
	return err
}
