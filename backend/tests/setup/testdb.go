package setup

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// NewTestDB creates an in-memory SQLite database for testing.
// Each call returns a fresh, isolated database.
func NewTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("opening in-memory test db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging in-memory test db: %w", err)
	}

	// Enable foreign keys and WAL mode for better concurrent test behavior
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("enabling foreign keys: %w", err)
	}

	return db, nil
}
