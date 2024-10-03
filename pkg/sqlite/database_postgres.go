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
	ret := NewDatabase()

	db := &PostgresDB{
		databaseFunctions: ret,
		storeRepository:   ret.storeRepository,
		lockChan:          ret.lockChan,
		dbType:            PostgresBackend,
		dbString:          dbConnector,
	}

	dbWrapper.dbType = PostgresBackend

	return (*Database)(db)
}

func (db *Database) open(disableForeignKeys bool, writable bool) (conn *sqlx.DB, err error) {
	conn, err = sqlx.Open("pgx", db.dbString)
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
