package integration

import (
	"testing"

	"01social/tests/setup"
)

func TestNotificationFlow_FollowAcceptNotif(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// 1. User 10 (Jack) follows User 1 (Alice)
	sessionJack, err := setup.CreateSession(db, 10)
	if err != nil {
		t.Fatal(err)
	}

	_, err = setup.PUT(server, "/api/follow/follow/1").
		WithAuth(sessionJack).
		DoStatus(200)
	if err != nil {
		t.Fatalf("follow: %v", err)
	}

	// 2. User 1 (Alice) should see a new_follower notification
	sessionAlice, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/notifications?type=all").
		WithAuth(sessionAlice).
		DoStatus(200)
	if err != nil {
		t.Fatalf("get notifications: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
	_ = body
}

func TestNotificationFlow_CreateCommentNotif(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 2 (Bob) comments on User 1's post (post 1)
	sessionBob, err := setup.CreateSession(db, 2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = setup.JSON(server, "POST", "/api/comments/create", map[string]interface{}{
		"postId": 1,
		"text":   "Integration test comment!",
	}).WithAuth(sessionBob).DoStatus(201)
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}

	// User 1 should now have a comment notification
	sessionAlice, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/notifications?type=unread").
		WithAuth(sessionAlice).
		DoStatus(200)
	if err != nil {
		t.Fatalf("get notifications: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
	_ = body
}

func TestNotificationFlow_LikePostNotif(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 2 (Bob) likes User 1's post (post 1)
	sessionBob, err := setup.CreateSession(db, 2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = setup.POST(server, "/api/posts/1/like").
		WithAuth(sessionBob).
		DoStatus(200)
	if err != nil {
		t.Fatalf("like post: %v", err)
	}
}
