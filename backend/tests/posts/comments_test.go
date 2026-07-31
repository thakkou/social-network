package posts

import (
	"testing"

	"01social/tests/setup"
)

func TestCreateComment_Success(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Comment on post 2 (Bob's post)
	body, err := setup.JSON(server, "POST", "/api/comments/create", map[string]interface{}{
		"postId": 2,
		"text":   "Great post, Bob!",
	}).WithAuth(sessionID).DoStatus(201)

	if err != nil {
		t.Fatalf("expected 201, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestCreateComment_EmptyText(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "POST", "/api/comments/create", map[string]interface{}{
		"postId": 2,
		"text":   "",
	}).WithAuth(sessionID).DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestCreateComment_InvalidPost(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Invalid post ID - handler returns 500 (internal error) for FK constraint violation
	resp, err := setup.JSON(server, "POST", "/api/comments/create", map[string]interface{}{
		"postId": 9999,
		"text":   "Comment on non-existent post",
	}).WithAuth(sessionID).Do()

	if err != nil {
		t.Fatal(err)
	}
	// The handler returns 500 because the foreign key constraint fails
	// (comment references non-existent post_id)
	if resp.StatusCode != 500 {
		t.Fatalf("expected 500 for invalid post, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestLikeComment(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Like comment 1 (posted by Bob on Alice's post)
	body, err := setup.POST(server, "/api/comments/1/like").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}

func TestDeleteComment_OwnComment(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Comment 5 was posted by User 1 (Alice)
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.DELETE(server, "/api/comments/5/delete").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}
