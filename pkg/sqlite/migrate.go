package sqlite

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	sqlite3mig "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
)

func (db *Database) needsMigration() bool {
	return db.schemaVersion != appSchemaVersion
}

type Migrator struct {
	db   *Database
	conn *sqlx.DB
	m    *migrate.Migrate
}

func NewMigrator(db *Database) (*Migrator, error) {
	m := &Migrator{
		db: db,
	}

	const disableForeignKeys = true
	const writable = true
	var err error
	m.conn, err = m.db.open(disableForeignKeys, writable)
	if err != nil {
		return nil, err
	}

	m.conn.SetMaxOpenConns(maxReadConnections)
	m.conn.SetMaxIdleConns(maxReadConnections)
	m.conn.SetConnMaxIdleTime(dbConnTimeout)

	m.m, err = m.getMigrate()

	// if error encountered, close the connection
	if err != nil {
		m.Close()
	}

	return m, err
}

func (m *Migrator) Close() {
	if m.m != nil {
		m.m.Close()
		m.m = nil
	}
}

func (m *Migrator) CurrentSchemaVersion() uint {
	databaseSchemaVersion, _, _ := m.m.Version()
	return databaseSchemaVersion
}

func (m *Migrator) RequiredSchemaVersion() uint {
	return appSchemaVersion
}

func (m *Migrator) getMigrate() (*migrate.Migrate, error) {
	migrations, err := iofs.New(migrationsBox, "migrations")
	if err != nil {
		return nil, err
	}

	driver, err := sqlite3mig.WithInstance(m.conn.DB, &sqlite3mig.Config{})
	if err != nil {
		return nil, err
	}

	// use sqlite3Driver so that migration has access to durationToTinyInt
	return migrate.NewWithInstance(
		"iofs",
		migrations,
		m.db.dbPath,
		driver,
	)
}

func (m *Migrator) RunMigration(ctx context.Context, newVersion uint) error {
	databaseSchemaVersion, _, _ := m.m.Version()

	if newVersion != databaseSchemaVersion+1 {
		return fmt.Errorf("invalid migration version %d, expected %d", newVersion, databaseSchemaVersion+1)
	}

	// run pre migrations as needed
	if err := m.runCustomMigrations(ctx, preMigrations[newVersion]); err != nil {
		return fmt.Errorf("running pre migrations for schema version %d: %w", newVersion, err)
	}

	if err := m.m.Steps(1); err != nil {
		// migration failed
		return err
	}

	// run post migrations as needed
	if err := m.runCustomMigrations(ctx, postMigrations[newVersion]); err != nil {
		return fmt.Errorf("running post migrations for schema version %d: %w", newVersion, err)
	}

	// update the schema version
	m.db.schemaVersion, _, _ = m.m.Version()

	return nil
}

func (m *Migrator) runCustomMigrations(ctx context.Context, fns []customMigrationFunc) error {
	for _, fn := range fns {
		if err := m.runCustomMigration(ctx, fn); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) runCustomMigration(ctx context.Context, fn customMigrationFunc) error {
	if err := fn(ctx, m.conn); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) PostMigrate(ctx context.Context) error {
	// optimise the database
	var err error
	logger.Info("Running database analyze")

	// don't use Optimize/vacuum as this adds a significant amount of time
	// to the migration
	err = analyze(ctx, m.conn)

	if err == nil {
		logger.Debug("Flushing WAL")
		err = flushWAL(ctx, m.conn)
	}

	if err != nil {
		return fmt.Errorf("error optimising database: %s", err)
	}

	return nil
}

func (db *Database) getDatabaseSchemaVersion() (uint, error) {
	m, err := NewMigrator(db)
	if err != nil {
		return 0, err
	}
	defer m.Close()

	ret, _, _ := m.m.Version()
	return ret, nil
}

func (db *Database) ReInitialise() error {
	return db.initialise()
}

// RunAllMigrations runs all migrations to bring the database up to the current schema version
func (db *Database) RunAllMigrations() error {
	ctx := context.Background()

	m, err := NewMigrator(db)
	if err != nil {
		return err
	}
	defer m.Close()

	databaseSchemaVersion, _, _ := m.m.Version()
	stepNumber := appSchemaVersion - databaseSchemaVersion
	if stepNumber != 0 {
		logger.Infof("Migrating database from version %d to %d", databaseSchemaVersion, appSchemaVersion)

		// run each migration individually, and run custom migrations as needed
		var i uint = 1
		for ; i <= stepNumber; i++ {
			newVersion := databaseSchemaVersion + i
			if err := m.RunMigration(ctx, newVersion); err != nil {
				return err
			}
		}
	}

	return nil
}
