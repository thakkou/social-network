package setup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// migrationDir returns the absolute path to the SQLite migration files.
func migrationDir() string {
	// Walk up from the test binary location to find the project root
	// Relative path from tests/setup/ to backend/pkg/db/migrations/sqlite/
	candidates := []string{
		"backend/pkg/db/migrations/sqlite",
		"../backend/pkg/db/migrations/sqlite",
		"../../backend/pkg/db/migrations/sqlite",
		"../../../backend/pkg/db/migrations/sqlite",
		"pkg/db/migrations/sqlite",
		"../pkg/db/migrations/sqlite",
		"../../pkg/db/migrations/sqlite",
	}

	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if info, err := os.Stat(abs); err == nil && info.IsDir() {
			return abs
		}
	}

	return ""
}

// RunMigrations runs all .up.sql migration files (excluding the seeder)
// on the given database. Returns an error if any migration fails.
func RunMigrations(db *sql.DB) error {
	dir := migrationDir()
	if dir == "" {
		return fmt.Errorf("migration directory not found")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading migration directory %s: %w", dir, err)
	}

	// Collect and sort migration files by name (ensures order: 000001, 000002, ...)
	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			// Skip the seeder migration — we seed separately with seeder_test.sql
			if strings.Contains(entry.Name(), "000007") {
				continue
			}
			upFiles = append(upFiles, entry.Name())
		}
	}
	sort.Strings(upFiles)

	for _, fileName := range upFiles {
		filePath := filepath.Join(dir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", fileName, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("running migration %s: %w", fileName, err)
		}
		log.Printf("[TEST] Ran migration: %s", fileName)
	}

	return nil
}

// RunDownMigrations runs all .down.sql migration files in reverse order.
// Useful for cleaning up between test suites.
func RunDownMigrations(db *sql.DB) error {
	dir := migrationDir()
	if dir == "" {
		return fmt.Errorf("migration directory not found")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading migration directory: %w", err)
	}

	var downFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".down.sql") {
			downFiles = append(downFiles, entry.Name())
		}
	}
	// Reverse order for down migrations
	sort.Sort(sort.Reverse(sort.StringSlice(downFiles)))

	for _, fileName := range downFiles {
		filePath := filepath.Join(dir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("reading down migration %s: %w", fileName, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Printf("[TEST] Warning running down migration %s (tables may not exist): %v", fileName, err)
			// Don't return error — tables might not exist yet
		}
	}

	return nil
}
