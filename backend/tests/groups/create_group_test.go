package groups

import (
	"testing"

	"01social/tests/setup"
)

func TestCreateGroup_Success(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/groups/create").
		WithAuth(sessionID).
		WithMultipart(
			map[string]string{
				"title":       "Test Group",
				"description": "A test group created during testing",
			},
			nil,
		).
		DoStatus(201)

	if err != nil {
		t.Fatalf("expected 201, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestCreateGroup_MissingTitle(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/groups/create").
		WithAuth(sessionID).
		WithMultipart(
			map[string]string{
				"title":       "",
				"description": "No title group",
			},
			nil,
		).
		DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestCreateGroup_Unauthenticated(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	resp, err := setup.POST(server, "/api/groups/create").Do()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
