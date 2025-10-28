package api

import (
	"net/http"

	"github.com/stashapp/stash/pkg/txn"
)

type routes struct {
	txnManager txn.Manager
}

func (rs routes) withReadTxn(r *http.Request, fn txn.TxnFunc) error {
	return txn.WithReadTxn(r.Context(), rs.txnManager, fn)
}
