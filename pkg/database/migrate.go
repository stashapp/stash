package database

import "context"

type MigrateStore interface {
	Close()
	CurrentSchemaVersion() uint
	PostMigrate(ctx context.Context) error
	RequiredSchemaVersion() uint
	RunMigration(ctx context.Context, newVersion uint) error
}
