package groups

import (
	"testing"

	"01social/tests/setup"
)

func TestGroupPostComment(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 is member of Group 1
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/groups/1/posts/1/comments", map[string]string{
		"text": "Great post!",
	}).WithAuth(sessionID).DoStatus(201)

	if err != nil {
		t.Fatalf("expected 201, got error: %v\nBody: %s", err, string(body))
	}
}

func TestGroupPostCommentReaction(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 is member of Group 1, comment 1 is on Group 1 post 1
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/groups/1/posts/1/comments/1/reaction", map[string]int{
		"is_like": 1,
	}).WithAuth(sessionID).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}
