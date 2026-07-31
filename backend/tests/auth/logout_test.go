package auth

import (
	"testing"

	"01social/tests/setup"
)

func TestLogout_Success(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Create a session for user 1 (Alice)
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/logout").
		WithAuth(sessionID).
		DoStatus(201)

	if err != nil {
		t.Fatalf("expected 201, got error: %v\nBody: %s", err, string(body))
	}
}

func TestLogout_WithoutAuth(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.POST(server, "/api/logout").Do()
	if err != nil {
		t.Fatal(err)
	}
	if body.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", body.StatusCode)
	}
	body.Body.Close()
}

func TestLogout_InvalidSession(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Use a fake session ID — middleware rejects it because logout requires auth
	// and the session ID doesn't exist in the database.
	resp, err := setup.POST(server, "/api/logout").
		WithAuth("invalid-session-id").
		Do()
	if err != nil {
		t.Fatal(err)
	}
	// Middleware looks up the session in DB and fails, returning 401
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
