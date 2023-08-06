package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migration/*.sql
var migrationFS embed.FS

// // Db represents the database connection.
type Db struct {
	*sql.DB
	DSN string
}

func NewDB(dsn string) (*Db, error) {
	db := &Db{
		DSN: dsn,
	}
	return db, nil
}

func NewOpenDB(dsn string) (*Db, error) {
	db := &Db{
		DSN: dsn,
	}
	err := db.Open()
	return db, err
}

func (sqdb *Db) Close() error {
	return sqdb.DB.Close()
}

// Open opens the database connection.
func (sqdb *Db) Open() (err error) {
	// Ensure a DSN is set before attempting to open the database.
	if sqdb.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	// Make the parent directory unless using an in-memory db.
	if sqdb.DSN != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(sqdb.DSN), 0700); err != nil {
			return err
		}
	}

	// Connect to the database.
	sqdb.DB, err = sql.Open("sqlite3", sqdb.DSN)
	if err != nil {
		return err
	}

	// Enable WAL. SQLite performs better with the WAL  because it allows
	// multiple readers to operate while data is being written.
	if _, err := sqdb.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	// Enable foreign key checks. For historical reasons, SQLite does not check
	// foreign key constraints by default... which is kinda insane. There's some
	// overhead on inserts to verify foreign key integrity but it's definitely
	// worth it.
	if _, err := sqdb.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	if err := sqdb.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

// migrate executes pending migration files.
//
// Migration files are embedded in the sqlite/migration folder and are executed
// in lexigraphical order.
//
// Once a migration is run, its name is stored in the 'migrations' table so it
// is not re-executed. Migrations run in a transaction to prevent partial
// migrations.
func (sqdb *Db) migrate() error {
	// Ensure the 'migrations' table exists so we don't duplicate migrations.
	if _, err := sqdb.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	// Read migration files from our embedded file system.
	// This uses Go 1.16's 'embed' package.
	names, err := fs.Glob(migrationFS, "migration/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	// Loop over all migration files and execute them in order.
	for _, name := range names {
		if err := sqdb.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}
	return nil
}

// migrate runs a single migration file within a transaction. On success, the
// migration file name is saved to the "migrations" table to prevent re-running.
func (sqdb *Db) migrateFile(name string) error {
	tx, err := sqdb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure migration has not already been run.
	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // already run migration, skip
	}

	// Read and execute migration file.
	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	// Insert record into migrations to prevent re-running migration.
	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES (?)`, name); err != nil {
		return err
	}

	return tx.Commit()
}
