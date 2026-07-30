package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var (
	instance *sql.DB
	once     sync.Once
	initErr  error
)

// SetDB sets the global database instance for testing purposes.
// This allows tests to inject an in-memory database without going through Init().
func SetDB(db *sql.DB) {
	instance = db
}

// Init initializes the global database connection and runs migrations.
// Call this once, early in main().
func Init() error {
	once.Do(func() {
		dbPath := "pkg/db/sn.db?_foreign_keys=on"

		conn, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			initErr = fmt.Errorf("opening sqlite db: %w", err)
			return
		}

		if err = conn.Ping(); err != nil {
			initErr = fmt.Errorf("pinging sqlite db: %w", err)
			return
		}

		if err = runMigrations(conn); err != nil {
			initErr = fmt.Errorf("running migrations: %w", err)
			return
		}

		instance = conn
	})

	return initErr
}

// DB returns the global database connection.
// Panics if Init() hasn't been called yet — this is intentional:
// using the DB before it's initialized is a programming error, not
// a runtime condition to handle gracefully.
func DB() *sql.DB {
	if instance == nil {
		panic("sqlite: DB() called before Init()")
	}
	return instance
}

// Close closes the global database connection.
func Close() error {
	if instance == nil {
		return nil
	}
	return instance.Close()
}

func runMigrations(conn *sql.DB) error {
	driver, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite3 driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsSourceURL(),
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}

	log.Println("migrations applied successfully")
	return nil
}

func migrationsSourceURL() string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	dir = filepath.Join(dir, "..", "migrations", "sqlite")
	abs, _ := filepath.Abs(dir)
	return "file://" + filepath.ToSlash(abs)
}
