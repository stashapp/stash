package manager

import "context"

// PostMigrate is executed after migrations have been executed.
func (s *Manager) PostMigrate(ctx context.Context) {
	setInitialMD5Config(ctx, s.Repository, s.Repository.Scene)
}
