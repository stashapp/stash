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
	postCompleteHooks []TxnFunc
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

func executeHooks(ctx context.Context, hooks []TxnFunc) {
	for _, h := range hooks {
		// ignore errors
		_ = h(ctx)
	}
}

func executePostCommitHooks(ctx context.Context) {
	m := hookManagerCtx(ctx)
	executeHooks(ctx, m.postCommitHooks)
}

func executePostRollbackHooks(ctx context.Context) {
	m := hookManagerCtx(ctx)
	executeHooks(ctx, m.postRollbackHooks)
}

func executePostCompleteHooks(ctx context.Context) {
	m := hookManagerCtx(ctx)
	executeHooks(ctx, m.postCompleteHooks)
}

func AddPostCommitHook(ctx context.Context, hook TxnFunc) {
	m := hookManagerCtx(ctx)
	m.postCommitHooks = append(m.postCommitHooks, hook)
}

func AddPostRollbackHook(ctx context.Context, hook TxnFunc) {
	m := hookManagerCtx(ctx)
	m.postRollbackHooks = append(m.postRollbackHooks, hook)
}

func AddPostCompleteHook(ctx context.Context, hook TxnFunc) {
	m := hookManagerCtx(ctx)
	m.postCompleteHooks = append(m.postCompleteHooks, hook)
}
