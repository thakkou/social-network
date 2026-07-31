package users

import (
	"testing"

	"01social/tests/setup"
)

func TestUpdatePrivacy_MakePrivate(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "PUT", "/api/profile/privacy", map[string]int{
		"is_private": 1,
	}).WithAuth(sessionID).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
	setup.AssertMessage(t, wrapper, "profile privacy updated")
}

func TestUpdatePrivacy_MakePublic(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 3 (Chloe) is private by default
	sessionID, err := setup.CreateSession(db, 3)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "PUT", "/api/profile/privacy", map[string]int{
		"is_private": 0,
	}).WithAuth(sessionID).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}

func TestUpdatePrivacy_InvalidValue(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "PUT", "/api/profile/privacy", map[string]int{
		"is_private": 2,
	}).WithAuth(sessionID).DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestUpdatePrivacy_Unauthenticated(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.JSON(server, "PUT", "/api/profile/privacy", map[string]int{
		"is_private": 1,
	}).Do()

	if err != nil {
		t.Fatal(err)
	}
	if body.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", body.StatusCode)
	}
	body.Body.Close()
}
