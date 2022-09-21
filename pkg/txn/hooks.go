package txn

import (
	"context"
)

type key int

const (
	hookManagerKey key = iota + 1
)

type hookManager struct {
	postCommitHooks   []TxnFunc
	postRollbackHooks []TxnFunc
}

func (m *hookManager) register(ctx context.Context) context.Context {
	return context.WithValue(ctx, hookManagerKey, m)
}

func hookManagerCtx(ctx context.Context) *hookManager {
	m, ok := ctx.Value(hookManagerKey).(*hookManager)
	if !ok {
		return nil
	}
	return m
}

func executePostCommitHooks(ctx context.Context) {
	m := hookManagerCtx(ctx)
	for _, h := range m.postCommitHooks {
		// ignore errors
		_ = h(ctx)
	}
}

func executePostRollbackHooks(ctx context.Context) {
	m := hookManagerCtx(ctx)
	for _, h := range m.postRollbackHooks {
		// ignore errors
		_ = h(ctx)
	}
}

func AddPostCommitHook(ctx context.Context, hook TxnFunc) {
	m := hookManagerCtx(ctx)
	m.postCommitHooks = append(m.postCommitHooks, hook)
}

func AddPostRollbackHook(ctx context.Context, hook TxnFunc) {
	m := hookManagerCtx(ctx)
	m.postRollbackHooks = append(m.postRollbackHooks, hook)
}
