package chat

import (
	"testing"

	"01social/tests/setup"
)

func TestSendDirectMessage(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 (Alice) sends a message to User 2 (Bob) — they follow each other
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":        "direct",
		"text":        "Hello Bob! Testing direct messages.",
		"receiver_id": 2,
	}).WithAuth(sessionID).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestSendDirectMessage_EmptyText(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":        "direct",
		"text":        "",
		"receiver_id": 2,
	}).WithAuth(sessionID).DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestSendDirectMessage_ToSelf(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":        "direct",
		"text":        "Message to myself",
		"receiver_id": 1,
	}).WithAuth(sessionID).DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestSendDirectMessage_NoFollowRelationship(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 (Alice) tries to message User 9 (Isabella) — not following each other
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":        "direct",
		"text":        "Hello stranger!",
		"receiver_id": 9,
	}).WithAuth(sessionID).Do()

	if err != nil {
		t.Fatal(err)
	}
	// Should be forbidden (403) because they don't follow each other
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestGetConversations(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/conversations").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}
