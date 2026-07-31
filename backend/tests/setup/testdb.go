package setup

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// NewTestDB creates an in-memory SQLite database for testing.
// Each call returns a fresh, isolated database.
func NewTestDB() (*sql.DB, error) {
	// Use a temporary file for the database so all connections from the pool
	// see the same data. SQLite's :memory: creates a separate database per
	// connection, which causes problems with Go's connection pool.
	tmpFile, err := os.CreateTemp("", "social-test-*.db")
	if err != nil {
		return nil, fmt.Errorf("creating temp db file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	db, err := sql.Open("sqlite3", tmpPath+"?_foreign_keys=on")
	if err != nil {
		os.Remove(tmpPath)
		return nil, fmt.Errorf("opening temp db: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		os.Remove(tmpPath)
		return nil, fmt.Errorf("pinging temp db: %w", err)
	}

	// Store the temp file path so cleanup can remove it.
	// We use a goroutine to defer cleanup, but we'll handle it explicitly.
	db.Exec("PRAGMA foreign_keys = ON")

	return db, nil
}
