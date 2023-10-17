package models

import (
	"net/http"

	"github.com/stashapp/stash/pkg/txn"
)

type TxnRoutes struct {
	TxnManager txn.Manager
}

func (rs TxnRoutes) WithReadTxn(r *http.Request, fn txn.TxnFunc) error {
	return txn.WithReadTxn(r.Context(), rs.TxnManager, fn)
}
