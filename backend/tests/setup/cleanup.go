package setup

import (
	"database/sql"
	"log"
	"net/http/httptest"
)

// CleanupTestServer closes the test server and database.
// Use with defer for easy cleanup:
//
//	server, db, err := NewTestServerWithDefaults()
//	if err != nil { t.Fatal(err) }
//	defer CleanupTestServer(server, db)
func CleanupTestServer(server *httptest.Server, db *sql.DB) {
	if server != nil {
		server.Close()
	}
	if db != nil {
		CleanupTestDB(db)
	}
}

// CleanupTestDB drops all tables and closes the database connection.
func CleanupTestDB(db *sql.DB) {
	if db == nil {
		return
	}

	// Drop all tables to clean up between tests
	tables := []string{
		"GROUP_POST_COMMENT_REACTIONS",
		"GROUP_POST_COMMENTS",
		"GROUP_POST_REACTIONS",
		"GROUP_POSTS",
		"GROUP_MESSAGE_READS",
		"GROUP_MESSAGES",
		"GROUP_REQUESTS",
		"GROUP_INVITES",
		"GROUP_MEMBERS",
		"GROUP_EVENTS",
		"EVENT_RESPONSES",
		"NOTIFICATIONS",
		"MESSAGES",
		"CONVERSATIONS",
		"WS_TICKETS",
		"GROUP_POST_COMMENT_REACTIONS",
		"COMMENT_REACTIONS",
		"POST_REACTIONS",
		"POST_ALLOWED_USERS",
		"POST_CATEGORY",
		"CATEGORY",
		"COMMENTS",
		"POSTS",
		"FOLLOWS",
		"SESSIONS",
		"USERS",
		"rate_limits",
	}

	for _, table := range tables {
		_, err := db.Exec("DROP TABLE IF EXISTS " + table)
		if err != nil {
			log.Printf("[TEST] Warning dropping table %s: %v", table, err)
		}
	}

	db.Close()
}

