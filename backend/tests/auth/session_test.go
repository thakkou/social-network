package auth

import (
	"testing"

	"01social/tests/setup"
)

func TestValidateSession_Valid(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// The /api/session/validate endpoint uses requiresAuth=false.
	// When a valid session IS found, the middleware returns 409 (Conflict)
	// because an authenticated user shouldn't be hitting a no-auth endpoint.
	// This is the expected behavior for the session validation flow.
	body, err := setup.GET(server, "/api/session/validate").
		WithAuth(sessionID).
		DoStatus(409)

	if err != nil {
		t.Fatalf("expected 409, got error: %v\nBody: %s", err, string(body))
	}
}

func TestValidateSession_NoCookie(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.GET(server, "/api/session/validate").DoStatus(401)
	if err != nil {
		t.Fatalf("expected 401, got error: %v\nBody: %s", err, string(body))
	}
}

func TestGetMe_Authenticated(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := setup.GET(server, "/api/me").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestGetMe_Unauthenticated(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.GET(server, "/api/me").Do()
	if err != nil {
		t.Fatal(err)
	}
	if body.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", body.StatusCode)
	}
	body.Body.Close()
}
