package followers

import (
	"testing"

	"01social/tests/setup"
)

func TestAcceptFollow_Success(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 3 (Chloe, private) has a pending follow request from User 7 (Grace)
	// User 3 accepts User 7's follow request
	sessionID, err := setup.CreateSession(db, 3)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.PUT(server, "/api/follow/accept/7").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}

func TestAcceptFollow_NoPendingRequest(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 (Alice) tries to accept a non-existent follow request
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// There's no pending request from user 9 to user 1
	body, err := setup.PUT(server, "/api/follow/accept/9").
		WithAuth(sessionID).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if body.StatusCode != 404 {
		t.Fatalf("expected 404, got %d", body.StatusCode)
	}
	body.Body.Close()
}

func TestRejectFollow_Success(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 5 (Emma, private) rejects the follow request from User 8 (Hugo)
	sessionID, err := setup.CreateSession(db, 5)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.PUT(server, "/api/follow/reject/8").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}
