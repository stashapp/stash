package manager

import "context"

// PostMigrate is executed after migrations have been executed.
func (s *singleton) PostMigrate(ctx context.Context) {
	setInitialMD5Config(ctx, s.TxnManager)
}
