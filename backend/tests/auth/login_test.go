package auth

import (
	"testing"

	"01social/tests/setup"
)

func TestLogin_Success_WithEmail(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// The seeder_test.sql now uses bcrypt hash for "password123"
	body, err := setup.JSON(server, "POST", "/api/login", map[string]string{
		"identifier": "alice@example.com",
		"password":   "password123",
	}).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
	setup.AssertMessage(t, wrapper, "Login Success")
}

func TestLogin_Success_WithNickname(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.JSON(server, "POST", "/api/login", map[string]string{
		"identifier": "ali_m",
		"password":   "password123",
	}).DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)
}

func TestLogin_WrongPassword(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.JSON(server, "POST", "/api/login", map[string]string{
		"identifier": "alice@example.com",
		"password":   "wrongpassword",
	}).DoStatus(401)

	if err != nil {
		t.Fatalf("expected 401, got error: %v\nBody: %s", err, string(body))
	}
}

func TestLogin_NonExistentUser(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.JSON(server, "POST", "/api/login", map[string]string{
		"identifier": "nonexistent@example.com",
		"password":   "password123",
	}).DoStatus(401)

	if err != nil {
		t.Fatalf("expected 401, got error: %v\nBody: %s", err, string(body))
	}
}

func TestLogin_WrongMethod(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.GET(server, "/api/login").DoStatus(405)
	if err != nil {
		t.Fatalf("expected 405, got error: %v\nBody: %s", err, string(body))
	}
}
