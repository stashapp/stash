package manager

import "context"

// PostMigrate is executed after migrations have been executed.
func (s *singleton) PostMigrate(ctx context.Context, preVersion, postVersion uint) {

	// this should only be run on existing systems prior to schema version 12
	if preVersion < 12 {
		setInitialMD5Config(ctx, s.TxnManager)
	}
}
