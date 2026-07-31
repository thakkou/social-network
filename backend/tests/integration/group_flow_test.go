package integration

import (
	"testing"

	"01social/tests/setup"
)

func TestGroupFlow_CreateJoinPost(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// 1. Create a group as User 1 (Alice)
	sessionAdmin, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/groups/create").
		WithAuth(sessionAdmin).
		WithMultipart(map[string]string{
			"title":       "Integration Test Group",
			"description": "A group created during integration testing",
		}, nil).
		DoStatus(201)
	if err != nil {
		t.Fatalf("create group: %v\nBody: %s", err, string(body))
	}

	// 2. User 6 (Farid) requests to join
	sessionFarid, err := setup.CreateSession(db, 6)
	if err != nil {
		t.Fatal(err)
	}

	body, err = setup.POST(server, "/api/groups/5/join").
		WithAuth(sessionFarid).
		DoStatus(200)
	if err != nil {
		t.Fatalf("join request: %v\nBody: %s", err, string(body))
	}

	// 3. Admin accepts the join request
	body, err = setup.POST(server, "/api/groups/5/requests/6/accept").
		WithAuth(sessionAdmin).
		DoStatus(200)
	if err != nil {
		t.Fatalf("accept request: %v\nBody: %s", err, string(body))
	}

	// 4. User 6 sends a group message
	body, err = setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":     "group",
		"text":     "Thanks for adding me!",
		"group_id": 5,
	}).WithAuth(sessionFarid).DoStatus(200)
	if err != nil {
		t.Fatalf("group message: %v\nBody: %s", err, string(body))
	}
}

func TestGroupFlow_InviteAccept(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 (Alice) invites User 9 (Isabella) to Group 1
	sessionAlice, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/groups/1/invite", map[string]interface{}{
		"user_id": 9,
	}).WithAuth(sessionAlice).DoStatus(200)
	if err != nil {
		t.Fatalf("invite: %v\nBody: %s", err, string(body))
	}

	// User 9 accepts the invite
	sessionIsabella, err := setup.CreateSession(db, 9)
	if err != nil {
		t.Fatal(err)
	}

	body, err = setup.POST(server, "/api/groups/1/invites/accept").
		WithAuth(sessionIsabella).
		DoStatus(200)
	if err != nil {
		t.Fatalf("accept invite: %v\nBody: %s", err, string(body))
	}
}
