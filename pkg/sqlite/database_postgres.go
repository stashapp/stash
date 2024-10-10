package sqlite

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
)

type PostgresDB Database

func NewPostgresDatabase(dbConnector string) *Database {
	dialect = goqu.Dialect("postgres")

	db := &PostgresDB{
		storeRepository: newDatabase(),
		dbConfig:        dbConnector,
	}
	db.dbInterface = db

	dbWrapper.dbType = PostgresBackend

	return (*Database)(db)
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
		logger.Warn("open with disableForeignKeys is not implemented.")
	}
	if !writable {
		_, err = conn.Exec("SET SESSION CHARACTERISTICS AS TRANSACTION READ ONLY;")

		if err != nil {
			return nil, fmt.Errorf("conn.Exec(): %w", err)
		}
	}

	return conn, nil
}

func (db *PostgresDB) Remove() error {
	logger.Warn("Postgres backend detected, ignoring Remove request")
	return nil
}

func (db *PostgresDB) Reset() error {
	logger.Warn("Postgres backend detected, ignoring Reset request")
	return nil
}

func (db *PostgresDB) Backup(backupPath string) (err error) {
	logger.Warn("Postgres backend detected, ignoring Backup request")
	return nil
}

func (db *PostgresDB) RestoreFromBackup(backupPath string) error {
	logger.Warn("Postgres backend detected, ignoring RestoreFromBackup request")
	return nil
}

func (db *PostgresDB) DatabaseBackupPath(backupDirectoryPath string) string {
	logger.Warn("Postgres backend detected, ignoring DatabaseBackupPath request")
	return ""
}
