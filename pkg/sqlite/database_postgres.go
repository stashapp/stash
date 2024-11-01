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

const (
	// TODO: Optimize for these
	maxPGWriteConnections = 10
	maxPGReadConnections  = 25
)

func NewPostgresDatabase(dbConnector string, init bool) *PostgresDB {
	db := &PostgresDB{
		Database: Database{
			storeRepository: newDatabase(),
			dbConfig:        dbConnector,
		},
	}
	db.DBInterface = db

	if init {
		dialect = goqu.Dialect("postgres")
		dbWrapper.dbType = PostgresBackend
	}

	return db
}

// Does nothing
func (db *PostgresDB) lock()   {}
func (db *PostgresDB) unlock() {}

func (db *PostgresDB) openReadDB() error {
	const (
		disableForeignKeys = false
		writable           = false
	)
	var err error
	db.readDB, err = db.open(disableForeignKeys, writable)
	db.readDB.SetMaxOpenConns(maxPGReadConnections)
	db.readDB.SetMaxIdleConns(maxPGReadConnections)
	db.readDB.SetConnMaxIdleTime(dbConnTimeout)
	return err
}

func (db *PostgresDB) openWriteDB() error {
	const (
		disableForeignKeys = false
		writable           = true
	)
	var err error
	db.writeDB, err = db.open(disableForeignKeys, writable)
	db.writeDB.SetMaxOpenConns(maxPGWriteConnections)
	db.writeDB.SetMaxIdleConns(maxPGWriteConnections)
	db.writeDB.SetConnMaxIdleTime(dbConnTimeout)
	return err
}

// Ensure single connection for testing to avoid race conditions
func (db *PostgresDB) TestMode() {
	db.readDB.Close()
	db.readDB = db.writeDB
}

func (db *PostgresDB) DatabaseType() DatabaseType {
	return PostgresBackend
}

func (db *PostgresDB) AppSchemaVersion() uint {
	return uint(0 - (66 - int(appSchemaVersion)))
}

func (db *PostgresDB) DatabaseConnector() string {
	return db.dbConfig
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
DO $$
DECLARE 
    r record;
BEGIN
    FOR r IN SELECT quote_ident(tablename) AS tablename, quote_ident(schemaname) AS schemaname FROM pg_tables WHERE schemaname = 'public'
    LOOP
        RAISE INFO 'Dropping table %.%', r.schemaname, r.tablename;
        EXECUTE format('DROP TABLE IF EXISTS %I.%I CASCADE', r.schemaname, r.tablename);
    END LOOP;
END$$;
`)

	return err
}

func (db *PostgresDB) Backup(backupPath string) (err error) {
	logger.Warn("Postgres backend detected, ignoring Backup request")
	return nil
}

func (db *PostgresDB) RestoreFromBackup(backupPath string) (err error) {
	logger.Warn("Postgres backend detected, ignoring RestoreFromBackup request")
	return nil
}

func (db *PostgresDB) DatabaseBackupPath(backupDirectoryPath string) string {
	logger.Warn("Postgres backend detected, ignoring DatabaseBackupPath request")
	return ""
}
