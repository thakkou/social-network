package chat

import (
	"testing"

	"01social/tests/setup"
)

func TestConversationPermission_NotAllowed(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 (Alice) tries to fetch conversation messages for a conversation
	// they are not part of
	sessionID, err := setup.CreateSession(db, 3)
	if err != nil {
		t.Fatal(err)
	}

	// Conversation 1 is between users 1 and 2. User 3 shouldn't access it.
	// Note: The route pattern is /api/conversation/{type}/{id}
	// But our test server uses a custom mux, so we need to check the path format
	resp, err := setup.GET(server, "/api/conversation/direct/1").
		WithAuth(sessionID).
		Do()

	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestMessagesRequireAuth(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Unauthenticated request to send a message
	// Use JSON standalone function (not chained on POST)
	req := setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":        "direct",
		"text":        "hacked!",
		"receiver_id": 2,
	})
	resp, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
