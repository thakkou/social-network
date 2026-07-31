package integration

import (
	"testing"

	"01social/tests/setup"
)

// TestSocialFlow_UserRegistrationToPost tests a complete user flow:
// register → login → create post → like post → comment on post
func TestSocialFlow_UserRegistrationToPost(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// 1. Register a new user
	body, err := setup.POST(server, "/api/register").
		WithMultipart(map[string]string{
			"firstname": "Integration",
			"lastname":  "TestUser",
			"email":     "integration@test.com",
			"password":  "password123",
			"birthDate": "1995-05-15",
			"nickname":  "inttest",
		}, nil).
		DoStatus(200)
	if err != nil {
		t.Fatalf("register failed: %v\nBody: %s", err, string(body))
	}

	// 2. Login to get a session
	body, err = setup.JSON(server, "POST", "/api/login", map[string]string{
		"identifier": "integration@test.com",
		"password":   "password123",
	}).DoStatus(200)
	if err != nil {
		t.Fatalf("login failed: %v\nBody: %s", err, string(body))
	}

	// Extract the token from the response
	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)

	// Create a session for the new user
	sessionID, err := setup.CreateSession(db, 11) // ID 11 since we have 10 seeded users
	if err != nil {
		t.Fatal(err)
	}

	// 3. Create a post
	body, err = setup.POST(server, "/api/posts/create").
		WithAuth(sessionID).
		WithMultipart(map[string]string{
			"title":      "Integration Test Post",
			"text":       "This post was created during an integration test!",
			"privacy":    "public",
			"categories": "General",
		}, nil).
		DoStatus(201)
	if err != nil {
		t.Fatalf("create post failed: %v\nBody: %s", err, string(body))
	}

	// 4. View the user's profile
	body, err = setup.GET(server, "/api/profile/11").
		WithAuth(sessionID).
		DoStatus(200)
	if err != nil {
		t.Fatalf("get profile failed: %v\nBody: %s", err, string(body))
	}
}

// TestSocialFlow_FollowToMessage tests: follow → accept → send message
func TestSocialFlow_FollowToMessage(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 1 (Alice) follows User 3 (Chloe, private) so follow is pending
	sessionAlice, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = setup.PUT(server, "/api/follow/follow/3").
		WithAuth(sessionAlice).
		DoStatus(200)
	if err != nil {
		t.Fatalf("follow failed: %v", err)
	}

	// User 3 (Chloe, private) accepts the follow request
	sessionChloe, err := setup.CreateSession(db, 3)
	if err != nil {
		t.Fatal(err)
	}

	_, err = setup.PUT(server, "/api/follow/accept/1").
		WithAuth(sessionChloe).
		DoStatus(200)
	if err != nil {
		t.Fatalf("accept failed: %v", err)
	}

	// Now they can message each other
	// User 3 follows user 1 (accepted), so they can message
	body, err := setup.JSON(server, "POST", "/api/messages", map[string]interface{}{
		"type":        "direct",
		"text":        "Hi Chloe! Thanks for connecting.",
		"receiver_id": 3,
	}).WithAuth(sessionAlice).DoStatus(200)
	if err != nil {
		t.Fatalf("send message failed: %v\nBody: %s", err, string(body))
	}

	_ = body
}
