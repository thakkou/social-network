package groups

import (
	"testing"

	"01social/tests/setup"
)

func TestGroupInvite_Accept(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 7 (Grace) has a pending invite to Group 1 from User 1 (Alice)
	sessionID, err := setup.CreateSession(db, 7)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/groups/1/invites/accept").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestGroupInvite_Reject(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 5 (Emma) has a pending invite to Group 2 from User 6 (Farid)
	sessionID, err := setup.CreateSession(db, 5)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/groups/2/invites/reject").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}
