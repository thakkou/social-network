package groups

import (
	"testing"

	"01social/tests/setup"
)

func TestGroupContent(t *testing.T) {
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

	resp, err := setup.GET(server, "/api/groups/1/content").WithAuth(sessionID).Do()
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestGroupPostReaction(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 is a member of Group 1, Group 1 has post 1
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/groups/1/posts/1/reaction", map[string]int{
		"is_like": 1,
	}).WithAuth(sessionID).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}
