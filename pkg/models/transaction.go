package models

import (
	"context"

	"github.com/stashapp/stash/pkg/logger"
)

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
			rbErr := txn.Rollback()
			if rbErr != nil {
				logger.Warnf("error while trying to roll back transaction: %v", rbErr)
			}
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			rbErr := txn.Rollback()
			if rbErr != nil {
				logger.Warnf("error while trying to roll back transaction: %v", rbErr)
			}
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
			rbErr := txn.Rollback()
			if rbErr != nil {
				logger.Warnf("error while trying to roll back RO transaction: %v", rbErr)
			}
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			rbErr := txn.Rollback()
			if rbErr != nil {
				logger.Warnf("error while trying to roll back RO transaction: %v", rbErr)
			}
		} else {
			// all good, commit
			err = txn.Commit()
		}
	}()

	err = fn(txn.Repository())
	return err
}
