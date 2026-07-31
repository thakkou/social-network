package notifications

import (
	"testing"

	"01social/tests/setup"
)

func TestEventNotificationExists(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 3 (Chloe) should have a group_event notification (notification 13)
	sessionID, err := setup.CreateSession(db, 3)
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

func TestMultipleUsersHaveEventNotifications(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	users := []int{4, 6} // These users should have event notifications
	for _, uid := range users {
		sessionID, err := setup.CreateSession(db, uid)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := setup.GET(server, "/api/notifications?type=all").
			WithAuth(sessionID).
			Do()
		if err != nil {
			t.Fatalf("user %d: %v", uid, err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("user %d: expected 200, got %d", uid, resp.StatusCode)
		}
		resp.Body.Close()
	}
}
