package sqlite

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/logger"
)

type PostgresDB Database

func NewPostgresDatabase(dbConnector string) *Database {
	dialect = goqu.Dialect("postgres")

	db := &PostgresDB{
		storeRepository: newDatabase(),
		lockChan:        make(chan struct{}, 1),
		dbConfig:        dbConnector,
	}
	db.dbInterface = db

	dbWrapper.dbType = PostgresBackend

	return (*Database)(db)
}

func (db *PostgresDB) DatabaseType() DatabaseType {
	return PostgresBackend
}

/*func (db *PostgresDB) AppSchemaVersion() uint {
	return uint(0 - (66 - int(appSchemaVersion)))
}*/

func (db *PostgresDB) DatabaseConnector() string {
	return db.dbConfig.(string)
}

func (db *PostgresDB) open(disableForeignKeys bool, writable bool) (conn *sqlx.DB, err error) {
	conn, err = sqlx.Open("pgx", db.DatabaseConnector())
	if err == nil {
		if disableForeignKeys {
			conn.Exec("SET session_replication_role = replica;")
		}
		if !writable {
			conn.Exec("SET default_transaction_read_only = ON;")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("db.Open(): %w", err)
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
