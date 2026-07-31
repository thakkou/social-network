package groups

import (
	"testing"

	"01social/tests/setup"
)

func TestGroupJoinRequest_Accept(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 8 (Hugo) has a pending join request for Group 1 (Gophers United)
	// User 1 (Alice) is the admin of Group 1

	// Step 1: Admin views pending requests
	sessionAdmin, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/groups/1/requests").
		WithAuth(sessionAdmin).
		DoStatus(200)

	if err != nil {
		t.Fatalf("view requests: %v\nBody: %s", err, string(body))
	}

	// Step 2: Admin accepts user 8's request
	body, err = setup.POST(server, "/api/groups/1/requests/8/accept").
		WithAuth(sessionAdmin).
		DoStatus(200)

	if err != nil {
		t.Fatalf("accept request: %v\nBody: %s", err, string(body))
	}
}

func TestGroupJoinRequest_Reject(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 4 has a pending join request for Group 2 (Sports Club)
	// User 6 (Farid) is the admin of Group 2
	sessionAdmin, err := setup.CreateSession(db, 6)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/groups/2/requests/4/reject").
		WithAuth(sessionAdmin).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}

func TestGroupJoinRequest_NonAdmin(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 2 (Bob) is a member (not admin) of Group 1
	sessionID, err := setup.CreateSession(db, 2)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/groups/1/requests").
		WithAuth(sessionID).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	// Non-admin should get 403
	if body.StatusCode != 403 {
		t.Fatalf("expected 403 for non-admin, got %d", body.StatusCode)
	}
	body.Body.Close()
}
