package sqlite

import (
	"context"

	"github.com/stashapp/stash/pkg/txn"
)

type hookManager struct {
	postCommitHooks   []txn.TxnFunc
	postRollbackHooks []txn.TxnFunc
}

func (m *hookManager) register(ctx context.Context) context.Context {
	return context.WithValue(ctx, hookManagerKey, m)
}

func (db *Database) hookManager(ctx context.Context) *hookManager {
	m, ok := ctx.Value(hookManagerKey).(*hookManager)
	if !ok {
		return nil
	}
	return m
}

func (db *Database) executePostCommitHooks(ctx context.Context) {
	m := db.hookManager(ctx)
	for _, h := range m.postCommitHooks {
		// ignore errors
		_ = h(ctx)
	}
}

func (db *Database) executePostRollbackHooks(ctx context.Context) {
	m := db.hookManager(ctx)
	for _, h := range m.postRollbackHooks {
		// ignore errors
		_ = h(ctx)
	}
}

func (db *Database) AddPostCommitHook(ctx context.Context, hook txn.TxnFunc) {
	m := db.hookManager(ctx)
	m.postCommitHooks = append(m.postCommitHooks, hook)
}

func (db *Database) AddPostRollbackHook(ctx context.Context, hook txn.TxnFunc) {
	m := db.hookManager(ctx)
	m.postRollbackHooks = append(m.postRollbackHooks, hook)
}
