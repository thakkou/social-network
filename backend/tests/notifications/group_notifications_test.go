package notifications

import (
	"testing"

	"01social/tests/setup"
)

func TestGroupInviteNotification(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 7 (Grace) has a group invite notification from Group 1
	sessionID, err := setup.CreateSession(db, 7)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/notifications?type=all").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestGroupJoinRequestNotification(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 8 (Hugo) has a pending join request for Group 1
	// User 1 (Alice) should have a group_join_request notification
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/notifications?type=all").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}
