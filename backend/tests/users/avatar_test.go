package users

import (
	"testing"

	"01social/tests/setup"
)

func TestUpdateAvatar_WithMultipart(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Send a minimal image-like payload
	body, err := setup.PUT(server, "/api/profile/update").
		WithAuth(sessionID).
		WithMultipart(
			map[string]string{
				"nickname": "alice_new",
				"aboutme":  "Updated!",
			},
			map[string][]byte{
				"avatar": []byte("fake-image-data"),
			},
		).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}
}
