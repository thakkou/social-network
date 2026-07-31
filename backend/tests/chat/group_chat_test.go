package chat

import (
	"testing"

	"01social/tests/setup"
)

func TestSendGroupMessage(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 (Alice) is a member of Group 1 (Gophers United)
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":     "group",
		"text":     "Hello Gophers! Testing group messaging.",
		"group_id": 1,
	}).WithAuth(sessionID).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestSendGroupMessage_NonMember(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 9 (Isabella) is not a member of Group 1
	sessionID, err := setup.CreateSession(db, 9)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":     "group",
		"text":     "I shouldn't be able to post here!",
		"group_id": 1,
	}).WithAuth(sessionID).Do()

	if err != nil {
		t.Fatal(err)
	}
	// Non-members should get 403
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestGetGroupMessages(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 is a member of Group 1
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := setup.GET(server, "/api/groups/1/messages").
		WithAuth(sessionID).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
