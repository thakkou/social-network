package users

import (
	"testing"

	"01social/tests/setup"
)

func TestGetProfile_OwnProfile(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/profile/1").
		WithAuth(sessionID).
		DoStatus(200)
	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestGetProfile_PublicUser(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 2)
	if err != nil {
		t.Fatal(err)
	}

	// User 2 (Bob) viewing User 1 (Alice) - they follow each other
	body, err := setup.GET(server, "/api/profile/1").
		WithAuth(sessionID).
		DoStatus(200)
	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}

func TestGetProfile_PrivateUserPendingFollow(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 3 (Chloe) is private
	sessionID, err := setup.CreateSession(db, 7)
	if err != nil {
		t.Fatal(err)
	}

	// User 7 (Grace) has a pending follow to User 3 (Chloe)
	// Should still see full profile since the follow exists?
	// Actually, pending means not accepted, so should see limited
	body, err := setup.GET(server, "/api/profile/3").
		WithAuth(sessionID).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	// Grace has a pending follow request to Chloe, so the status should be reflected
	if body.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", body.StatusCode)
	}
	body.Body.Close()
}

func TestGetProfile_NonExistentUser(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/profile/999").
		WithAuth(sessionID).
		Do()
	if err != nil {
		t.Fatal(err)
	}
	if body.StatusCode != 404 {
		t.Fatalf("expected 404 for non-existent user, got %d", body.StatusCode)
	}
	body.Body.Close()
}
