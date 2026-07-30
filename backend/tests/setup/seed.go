package setup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// SeedTestData executes the seeder_test.sql file on the given database.
// The seeder file contains 10 users, 63 posts, 88 comments, follows,
// groups, events, notifications, and all related data.
func SeedTestData(db *sql.DB) error {
	// Find the seeder_test.sql file
	seederPath := findSeederFile()
	if seederPath == "" {
		return fmt.Errorf("seeder_test.sql not found")
	}

	content, err := os.ReadFile(seederPath)
	if err != nil {
		return fmt.Errorf("reading seeder file %s: %w", seederPath, err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("executing seeder SQL: %w", err)
	}

	log.Printf("[TEST] Seeded test data from %s", seederPath)
	return nil
}

// findSeederFile searches for seeder_test.sql in likely locations.
func findSeederFile() string {
	candidates := []string{
		"backend/tests/setup/seeder_test.sql",
		"tests/setup/seeder_test.sql",
		"../tests/setup/seeder_test.sql",
		"../../tests/setup/seeder_test.sql",
		"../../../tests/setup/seeder_test.sql",
	}

	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if _, err := os.Stat(abs); err == nil {
			return abs
		}
	}

	return ""
}
