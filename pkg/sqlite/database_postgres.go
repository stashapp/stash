package sqlite

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
)

type PostgresDB struct {
	Database
}

func NewPostgresDatabase(dbConnector string) *PostgresDB {
	dialect = goqu.Dialect("postgres")

	db := &PostgresDB{
		Database: Database{
			storeRepository: newDatabase(),
			dbConfig:        dbConnector,
		},
	}
	db.DBInterface = db

	dbWrapper.dbType = PostgresBackend

	return db
}

// Does nothing
func (db *PostgresDB) lock()   {}
func (db *PostgresDB) unlock() {}

func (db *PostgresDB) DatabaseType() DatabaseType {
	return PostgresBackend
}

func (db *PostgresDB) AppSchemaVersion() uint {
	return uint(0 - (66 - int(appSchemaVersion)))
}

func (db *PostgresDB) DatabaseConnector() string {
	return db.dbConfig.(string)
}

func (db *PostgresDB) open(disableForeignKeys bool, writable bool) (conn *sqlx.DB, err error) {
	conn, err = sqlx.Open("pgx", db.DatabaseConnector())

	if err != nil {
		return nil, fmt.Errorf("db.Open(): %w", err)
	}

	if disableForeignKeys {
		_, err = conn.Exec("SET session_replication_role = replica;")

		if err != nil {
			return nil, fmt.Errorf("conn.Exec(): %w", err)
		}
	}
	if !writable {
		_, err = conn.Exec("SET SESSION CHARACTERISTICS AS TRANSACTION READ ONLY;")

		if err != nil {
			return nil, fmt.Errorf("conn.Exec(): %w", err)
		}
	}

	return conn, nil
}

func (db *PostgresDB) Remove() (err error) {
	_, err = db.writeDB.Exec(`
DO $$ DECLARE
    r RECORD;
BEGIN
    -- Disable triggers to avoid foreign key constraint violations
    EXECUTE 'SET session_replication_role = replica';

    -- Drop all tables
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;

    -- Re-enable triggers
    EXECUTE 'SET session_replication_role = DEFAULT';
END $$;
`)

	return err
}

func (db *PostgresDB) Backup(backupPath string) (err error) {
	logger.Warn("Postgres backend detected, ignoring Backup request")
	return nil
}

// RestoreFromBackup restores the database from a backup file at the given path.
func (db *PostgresDB) RestoreFromBackup(backupPath string) (err error) {
	logger.Warn("Postgres backend detected, ignoring RestoreFromBackup request")
	return nil
}

// DatabaseBackupPath returns the path to a database backup file for the given directory.
func (db *PostgresDB) DatabaseBackupPath(backupDirectoryPath string) string {
	logger.Warn("Postgres backend detected, ignoring DatabaseBackupPath request")
	return ""
}
