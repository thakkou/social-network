package posts

import (
	"testing"

	"01social/tests/setup"
)

func TestCreatePost_Success(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/posts/create").
		WithAuth(sessionID).
		WithMultipart(
			map[string]string{
				"title":      "My Test Post",
				"text":       "This is a test post content",
				"privacy":    "public",
				"categories": "General",
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

func TestCreatePost_MissingTitle(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/posts/create").
		WithAuth(sessionID).
		WithMultipart(
			map[string]string{
				"title":   "",
				"text":    "Missing title",
				"privacy": "public",
			},
			nil,
		).
		DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestCreatePost_MissingCategories(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/posts/create").
		WithAuth(sessionID).
		WithMultipart(
			map[string]string{
				"title":   "No Category Post",
				"text":    "Post without categories",
				"privacy": "public",
			},
			nil,
		).
		DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestCreatePost_InvalidPrivacy(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/posts/create").
		WithAuth(sessionID).
		WithMultipart(
			map[string]string{
				"title":      "Invalid Privacy",
				"text":       "Post with invalid privacy",
				"privacy":    "invalid",
				"categories": "General",
			},
			nil,
		).
		DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestCreatePost_PrivatePostWithAllowedUsers(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 2)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.POST(server, "/api/posts/create").
		WithAuth(sessionID).
		WithMultipart(
			map[string]string{
				"title":            "Private Post",
				"text":             "Only visible to specific users",
				"privacy":          "private",
				"categories":       "General",
				"allowed_user_ids": "1,7",
			},
			nil,
		).
		DoStatus(201)

	if err != nil {
		t.Fatalf("expected 201, got error: %v\nBody: %s", err, string(body))
	}
}
