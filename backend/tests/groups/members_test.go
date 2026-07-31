package groups

import (
	"testing"

	"01social/tests/setup"
)

func TestListGroups(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/groups").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestGetMyGroups(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// User 1 is a member of Groups 1, 2, 4
	body, err := setup.GET(server, "/api/users/groups").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}

func TestGroupJoinAndLeave(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 9 (Isabella) is not in Group 1 (Gophers United)
	sessionID, err := setup.CreateSession(db, 9)
	if err != nil {
		t.Fatal(err)
	}

	// Join request
	body, err := setup.POST(server, "/api/groups/1/join").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("join failed: %v\nBody: %s", err, string(body))
	}

	// Leave group (user 9 is not a member yet, only requested)
	resp, err := setup.POST(server, "/api/groups/1/leave").
		WithAuth(sessionID).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestGroupMemberKick(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 (Alice) is admin of Group 1, can kick User 2 (Bob)
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/groups/1/members/2/kick").
		WithAuth(sessionID).
		Do()

	if err != nil {
		t.Fatal(err)
	}
	// Should succeed since Alice is the creator of Group 1
	if body.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", body.StatusCode)
	}
	body.Body.Close()
}
