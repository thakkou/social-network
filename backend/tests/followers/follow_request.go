package followers

import (
	"testing"

	"01social/tests/setup"
)

func TestFollowRequest_Flow(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Step 1: User 9 (Isabella) follows User 5 (Emma, private)
	session9, err := setup.CreateSession(db, 9)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.PUT(server, "/api/follow/follow/5").
		WithAuth(session9).
		DoStatus(200)
	if err != nil {
		t.Fatalf("follow failed: %v\nBody: %s", err, string(body))
	}

	// Step 2: User 5 (Emma) accepts the follow
	session5, err := setup.CreateSession(db, 5)
	if err != nil {
		t.Fatal(err)
	}

	body, err = setup.PUT(server, "/api/follow/accept/9").
		WithAuth(session5).
		DoStatus(200)
	if err != nil {
		t.Fatalf("accept failed: %v\nBody: %s", err, string(body))
	}

	// Step 3: User 9 unfollows User 5
	body, err = setup.PUT(server, "/api/follow/unfollow/5").
		WithAuth(session9).
		DoStatus(200)
	if err != nil {
		t.Fatalf("unfollow failed: %v\nBody: %s", err, string(body))
	}
}

func TestFollow_Twice_Same(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Follow -> Accept -> Follow again should work (upsert)
	session7, err := setup.CreateSession(db, 7)
	if err != nil {
		t.Fatal(err)
	}

	// User 7 follows User 1 (already following through seed data? No, only 7->3)
	body, err := setup.PUT(server, "/api/follow/follow/1").
		WithAuth(session7).
		DoStatus(200)
	if err != nil {
		t.Fatalf("first follow: %v\nBody: %s", err, string(body))
	}

	// Follow again - should be idempotent
	body, err = setup.PUT(server, "/api/follow/follow/1").
		WithAuth(session7).
		DoStatus(200)
	if err != nil {
		t.Fatalf("second follow: %v\nBody: %s", err, string(body))
	}
}
