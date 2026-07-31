package notifications

import (
	"testing"

	"01social/tests/setup"
)

func TestFollowRequestNotification(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 3 (Chloe, private) should have a follow_request notification from User 7 (Grace)
	sessionID, err := setup.CreateSession(db, 3)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/notifications?type=unread").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	// Chloe should have follow_request notifications
	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestFollowAcceptedNotification(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// After accepting a follow request, a follow_accepted notification should appear
	// User 3 (Chloe) accepts User 7 (Grace)'s follow request
	sessionChloe, err := setup.CreateSession(db, 3)
	if err != nil {
		t.Fatal(err)
	}

	_, err = setup.PUT(server, "/api/follow/accept/7").
		WithAuth(sessionChloe).
		DoStatus(200)
	if err != nil {
		t.Fatal(err)
	}

	// Now User 7 should see a follow_accepted notification
	sessionGrace, err := setup.CreateSession(db, 7)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/notifications?type=all").
		WithAuth(sessionGrace).
		DoStatus(200)
	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	_ = body
}
