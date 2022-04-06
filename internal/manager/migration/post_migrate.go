package migration

import (
	"context"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

type PostMigrator struct {
	TxnManager models.TransactionManager
	Config     *config.Instance
}

// PostMigrate is executed after migrations have been executed.
func (m *PostMigrator) PostMigrate(ctx context.Context, preVersion, postVersion uint) {
	if preVersion < 12 {
		m.migrate12(ctx)
	}
}
