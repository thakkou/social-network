package users

import (
	"testing"

	"01social/tests/setup"
)

func TestUpdateProfile_NicknameAndAbout(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.JSON(server, "PUT", "/api/profile/update", map[string]string{
		"nickname": "new_alice",
		"aboutme":  "Updated bio for Alice!",
	}).WithAuth(sessionID).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestUpdateProfile_Unauthenticated(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	resp, err := setup.JSON(server, "PUT", "/api/profile/update", map[string]string{
		"nickname": "hacker",
	}).Do()

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
